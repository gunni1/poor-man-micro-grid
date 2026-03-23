package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"poor-man-micro-grid.local/shared"
)

const assetType = "wind"

func main() {

	broker := shared.GetEnv("MQTT_BROKER", "tcp://localhost:1883")
	assetId := shared.GetEnv("ASSET_ID", fmt.Sprintf("%s-%s", assetType, uuid.NewString()))
	pNominalInKw := shared.GetEnvAsFloat("P_NOMINAL_KW")

	clientId := fmt.Sprintf("%s-%s", assetType, assetId)
	client := shared.Connect(broker, clientId)

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

		msg := shared.Telemetry{
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
