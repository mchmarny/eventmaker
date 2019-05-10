package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	uuid "github.com/satori/go.uuid"
)

const (
	host     = "mqtt.googleapis.com"
	port     = "8883"
	idPrefix = "eid"
)

var (
	deviceID                = flag.String("device", "", "Cloud IoT Core Device ID")
	projectID               = flag.String("project", "", "GCP Project ID")
	registryID              = flag.String("registry", "", "Cloud IoT Registry ID (short form)")
	region                  = flag.String("region", "us-central1", "GCP Region")
	eventSrc                = flag.String("src", "iot-demo-client", "Name of the event source")
	certsCA                 = flag.String("ca", "root-ca.pem", "Download https://pki.google.com/roots.pem")
	privateKey              = flag.String("key", "device.key.pem", "Path to private key file")
	metricLabel             = flag.String("metric", "my-label", "Name of the metric label")
	metricRange             = flag.String("range", "1-10", "Numeric metric range [1-10]")
	eventFreq               = flag.String("freq", "5s", "Event frequency [5s]")
	errorInvalidMetricRange = errors.New("Invalid metric range format. Expected min-max (e.g. 1-10)")
)

func main() {

	flag.Parse()

	log.Println("Loading Google's roots...")
	certpool := x509.NewCertPool()
	pemCerts, err := ioutil.ReadFile(*certsCA)
	if err == nil {
		certpool.AppendCertsFromPEM(pemCerts)
	}

	config := &tls.Config{
		RootCAs:            certpool,
		ClientAuth:         tls.NoClientCert,
		ClientCAs:          nil,
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{},
		MinVersion:         tls.VersionTLS12,
	}

	clientID := fmt.Sprintf("projects/%v/locations/%v/registries/%v/devices/%v",
		*projectID,
		*region,
		*registryID,
		*deviceID,
	)

	log.Println("Creating MQTT client options...")
	opts := MQTT.NewClientOptions()

	broker := fmt.Sprintf("ssl://%v:%v", host, port)
	log.Printf("Broker '%v'", broker)

	opts.AddBroker(broker)
	opts.SetClientID(clientID).SetTLSConfig(config)
	opts.SetUsername("using-keys-instead")

	token := jwt.New(jwt.SigningMethodRS256)
	token.Claims = jwt.StandardClaims{
		Audience:  *projectID,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	}

	log.Println("Loading private key...")
	keyBytes, err := ioutil.ReadFile(*privateKey)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Parsing private key...")
	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyBytes)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Signing token")
	tokenString, err := token.SignedString(key)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Setting password...")
	opts.SetPassword(tokenString)

	opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
		fmt.Printf("[handler] Topic: %v\n", msg.Topic())
		fmt.Printf("[handler] Payload: %v\n", msg.Payload())
	})

	log.Println("Connecting...")
	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	log.Println("Publishing messages...")

	// event freq
	freq, err := time.ParseDuration(*eventFreq)
	if err != nil {
		log.Fatal(err)
	}

	rangeParts := strings.Split(*metricRange, "-")
	if len(rangeParts) != 2 {
		log.Fatal(errorInvalidMetricRange)
	}
	min, minErr := strconv.ParseFloat(rangeParts[0], 64)
	max, maxErr := strconv.ParseFloat(rangeParts[1], 64)
	if minErr != nil || maxErr != nil {
		log.Fatal(errorInvalidMetricRange)
	}

	for {
		data := makeEvent(min, max)
		log.Printf("Publishing: %v", data)
		token := client.Publish(
			fmt.Sprintf("/devices/%v/config", *deviceID),
			0,
			false,
			data)
		if !token.WaitTimeout(5 * time.Second) {
			fmt.Printf("Publish timed-out: %v\n", data)
		}
		time.Sleep(freq)
	}

}

func makeEvent(min, max float64) string {

	event := struct {
		SourceID    string    `json:"source_id"`
		EventID     string    `json:"event_id"`
		EventTs     time.Time `json:"event_ts"`
		MetricLabel string    `json:"label"`
		MetricValue float64   `json:"metric"`
	}{
		SourceID:    *eventSrc,
		EventID:     fmt.Sprintf("%s-%s", idPrefix, uuid.NewV4().String()),
		EventTs:     time.Now().UTC(),
		MetricLabel: *metricLabel,
		MetricValue: min + rand.Float64()*(max-min),
	}

	data, _ := json.Marshal(event)

	return string(data)

}
