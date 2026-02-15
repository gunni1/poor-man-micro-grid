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

const assetType = "load"

type Telemetry struct {
	AssetId   string  `json:"asset_id"`
	AssetType string  `json:"asset_type"`
	PkW       float64 `json:"p_kw"`
}

type LoadModel struct {
	P     float64 // aktuelle Leistung
	Base  float64 // Grundlast
	Peak  float64 // Zusatzlast tagsüber
	Alpha float64 // Trägheit
	Sigma float64 // Zufall
}

func main() {
	broker := getEnv("MQTT_BROKER", "tcp://localhost:1883")
	assetId := getEnv("ASSET_ID", fmt.Sprintf("%s-%s", assetType, uuid.NewString()))
	baseLoad := getEnvAsFloat("BASE_LOAD")
	peakLoad := getEnvAsFloat("PEAK_LOAD")

	model := LoadModel{
		P:     baseLoad + rand.NormFloat64()*50,
		Base:  baseLoad,
		Peak:  peakLoad,
		Alpha: 0.08,
		Sigma: 20,
	}

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

	for range ticker.C {
		p := model.Step(time.Now())
		telemetry := Telemetry{
			AssetId:   assetId,
			AssetType: assetType,
			PkW:       p,
		}
		payload, err := json.Marshal(telemetry)
		if err != nil {
			log.Printf("Error marshaling telemetry: %v", err)
			continue
		}
		token := client.Publish(topic, 0, false, payload)
		token.Wait()
		if token.Error() != nil {
			log.Printf("Error publishing telemetry: %v", token.Error())
		} else {
			log.Printf("Published telemetry: %s", string(payload))
		}
	}

}

func (model *LoadModel) Step(t time.Time) float64 {
	hour := float64(t.Hour()) + float64(t.Minute())/60.0
	dayFactor := math.Sin((hour - 6) / 12.0 * math.Pi)
	if dayFactor < 0 {
		dayFactor = 0
	}
	target := model.Base + model.Peak*dayFactor
	noise := rand.NormFloat64() * model.Sigma
	model.P += model.Alpha*(target-model.P) + noise
	if model.P < 0 {
		model.P = 0
	}
	return model.P
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

func getEnv(key string, fallback string) string {
	envValue, isPresent := os.LookupEnv(key)
	if !isPresent {
		return fallback
	}
	return envValue
}
