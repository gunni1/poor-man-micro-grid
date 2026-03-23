package shared

import (
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func Connect(broker string, clientId string) mqtt.Client {
	opts := mqtt.NewClientOptions().
		AddBroker(broker).
		SetClientID(clientId)

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	return client
}

func TelemetryTopic(assetType string, assetId string) string {
	return fmt.Sprintf("microgrid/%s/%s/telemetry", assetType, assetId)
}
