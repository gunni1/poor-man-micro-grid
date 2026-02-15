package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand/v2"
	"os"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

const assetType = "pv"

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
	log.Printf("Starting PV sim for asset: %s with pNom: %.2f", assetId, pNominalInKw)

	for range ticker.C {
		p := simPVPower(time.Now(), pNominalInKw)

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

func simPVPower(now time.Time, pNominal float64) float64 {
	decimalHour := float64(now.Hour()) + float64(now.Minute())/60.0
	//Sun only between hour 6 and 18
	if decimalHour < 6 || decimalHour > 18 {
		return 0
	}
	// use sin for sun cycle...
	x := (decimalHour - 6) / 12.0 * math.Pi
	p := math.Sin(x) * pNominal

	// Random cloud factor between 0.9 and 1.1
	cloudFactor := 0.9 + rand.Float64()*0.2
	return math.Max(0, p*cloudFactor)
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
