package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

const assetType = "wind"

type Telemetry struct {
	AssetId   string  `json:"asset_id"`
	AssetType string  `json:"asset_type"`
	PkW       float64 `json:"p_kw"`
}

func main() {

	broker := getEnv("MQTT_BROKER", "tcp://localhost:1883")
	assetId := getEnv("ASSET_ID", fmt.Sprintf("%s-%s", assetType, uuid.NewString()))
	pNominalInKw := getEnvAsFloat("P_NOMINAL_KW")

	// MQTT Setup
	opts := mqtt.NewClientOptions().
		AddBroker(broker).
		SetClientID(fmt.Sprintf("%s-%s", assetType, assetId))

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	defer client.Disconnect(250)

	topic := fmt.Sprintf("microgrid/%s/%s/telemetry", assetType, assetId)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	log.Printf("Starting wind sim for asset: %s with pNom: %.2f", assetId, pNominalInKw)

	//AR Model initialization
	model := WindModel{
		V:     7.5,
		VMean: 8.0,
		Alpha: 0.005,
		Sigma: 0.6,
	}

	for range ticker.C {
		v := model.Step()
		p := windPower(v, pNominalInKw)

		msg := Telemetry{
			AssetId:   assetId,
			AssetType: assetType,
			PkW:       p,
		}

		payload, _ := json.Marshal(msg)
		token := client.Publish(topic, 0, false, payload)
		token.Wait()
		log.Printf("Published to topic %s: %s", topic, string(payload))
	}

}

type WindModel struct {
	V     float64
	VMean float64
	Alpha float64 // mean reversion rate
	Sigma float64 // standard deviation of noise
}

// Calculates wind speed for next step
func (model *WindModel) Step() float64 {
	noise := rand.NormFloat64() * model.Sigma
	model.V = model.V + model.Alpha*(model.VMean-model.V) + noise

	if model.V < 0 {
		model.V = 0
	}
	return model.V
}

func windPower(v float64, pRated float64) float64 {
	const (
		cutIn  = 3.0  // m/s
		rated  = 12.0 // m/s
		cutOut = 25.0 // m/s
	)

	switch {
	case v < cutIn:
		return 0
	case v >= cutIn && v < rated:
		return pRated * math.Pow(v/cutIn, 3) / math.Pow(rated/cutIn, 3)
	case v >= rated && v < cutOut:
		return pRated
	default:
		return 0
	}
}

func getEnv(key string, fallback string) string {
	envValue, isPresent := os.LookupEnv(key)
	if !isPresent {
		return fallback
	}
	return envValue
}

func getEnvAsFloat(key string) float64 {
	asStr := envMandatory(key)
	val, err := strconv.ParseFloat(asStr, 64)
	if err != nil {
		log.Fatalf("Error parsing %s as float: %v", key, err)
	}
	return val
}

func envMandatory(key string) string {
	envValue, isPresent := os.LookupEnv(key)
	if !isPresent {
		log.Fatalf("environment var %s is missing", key)
	}
	return envValue
}
