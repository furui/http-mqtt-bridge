package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	authKey := os.Getenv("AUTH_KEY")
	mqttHost := os.Getenv("MQTT_HOST")
	mqttUser := os.Getenv("MQTT_USER")
	mqttPass := os.Getenv("MQTT_PASS")

	if authKey == "" || mqttHost == "" {
		log.Fatal("AUTH_KEY and MQTT_HOST environment variables are required.")
	}

	opts := mqtt.NewClientOptions().AddBroker(mqttHost)
	opts.SetClientID("http-mqtt-bridge")
	if mqttUser != "" {
		opts.SetUsername(mqttUser)
	}
	if mqttPass != "" {
		opts.SetPassword(mqttPass)
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect to MQTT broker: %v", token.Error())
	}
	log.Printf("Connected to MQTT broker at %s", mqttHost)

	http.HandleFunc("/publish", publishHandler(client, authKey))

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting HTTP server on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}

func publishHandler(client mqtt.Client, validAuthKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse parameters", http.StatusBadRequest)
			return
		}

		providedKey := r.FormValue("auth_key")
		topic := r.FormValue("topic")
		payload := r.FormValue("payload")

		if providedKey != validAuthKey {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if topic == "" || payload == "" {
			http.Error(w, "Missing 'topic' or 'payload' parameter", http.StatusBadRequest)
			return
		}

		token := client.Publish(topic, 1, false, payload)
		token.Wait()

		if token.Error() != nil {
			log.Printf("Publish error: %v", token.Error())
			http.Error(w, "Failed to publish message", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Message published to %s\n", topic)
	}
}
