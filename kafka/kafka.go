//  This file is part of the eliona project.
//  Copyright © 2022 LEICOM iTEC AG. All Rights Reserved.
//  ______ _ _
// |  ____| (_)
// | |__  | |_  ___  _ __   __ _
// |  __| | | |/ _ \| '_ \ / _` |
// | |____| | | (_) | | | | (_| |
// |______|_|_|\___/|_| |_|\__,_|
//
//  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING
//  BUT NOT LIMITED  TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
//  NON INFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
//  DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package kafka

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/eliona-smart-building-assistant/go-eliona/common"
	"github.com/eliona-smart-building-assistant/go-eliona/log"
	"github.com/google/uuid"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// Brokers reads the kafka brokers from environment variable BROKERS
func Brokers() string {
	brokers := common.Getenv("BROKERS", "kafka:29092")
	brokers = strings.Replace(brokers, "|", ",", -1)
	if len(brokers) > 0 {
		log.Debug("Kafka", "Boostrap servers %s.", brokers)
	} else {
		log.Debug("Kafka", "Explicitly no bootstrap servers defined.")
	}
	return brokers
}

func NewProducer() *kafka.Producer {

	configMap := kafka.ConfigMap{
		"client.id": uuid.NewString(),
	}
	if brokers := Brokers(); len(brokers) > 0 {
		_ = configMap.SetKey("bootstrap.servers", brokers)
	}

	// Creates a new producer
	producer, err := kafka.NewProducer(&configMap)
	if err != nil {
		log.Fatal("Kafka", "Failed to create producer: %s", err)
	}

	// Set report handler for produced messages
	go func() {
		for e := range producer.Events() {
			switch event := e.(type) {
			case *kafka.Message:
				if event.TopicPartition.Error != nil {
					log.Error("Kafka", "Delivery failed: %v\n", event.TopicPartition)
				} else {
					log.Debug("Kafka", "Delivered message to %v\n", event.TopicPartition)
				}
			}
		}
	}()

	log.Debug("Kafka", "New producer created")
	return producer
}

// Produce sends a typed message to a specific kafka topic
func Produce(producer *kafka.Producer, topic string, value any) error {
	payload, err := json.Marshal(value)
	if err != nil {
		log.Error("Kafka", "Failed to marshal value: %s", err.Error())
		return err
	}

	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          payload,
	}, nil)

	if err != nil {
		log.Error("Kafka", "Failed to deliver message: %s", err.Error())
		return err
	}

	producer.Flush(1000 * 10)
	return nil
}

func NewConsumer() *kafka.Consumer {
	configMap := kafka.ConfigMap{
		"group.id": uuid.NewString(),
	}
	if brokers := Brokers(); len(brokers) > 0 {
		_ = configMap.SetKey("bootstrap.servers", brokers)
	}

	consumer, err := kafka.NewConsumer(&configMap)
	if err != nil {
		log.Fatal("Kafka", "Failed to create consumer: %s", err)
	}
	log.Debug("Kafka", "New consumer created")
	return consumer
}

func Subscribe(consumer *kafka.Consumer, topic string) {
	// Subscribe to topic
	err := consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		log.Fatal("Kafka", "Failed to subscribe topic %s: %s", topic, err)
	}
}

func Read[T any](consumer *kafka.Consumer, values chan T) {

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	run := true

	for run {
		select {
		case sig := <-sig:
			log.Debug("Kafka", "Caught signal %v: terminating", sig)
			run = false
		default:
			e := consumer.Poll((int)(time.Second * 1000.0))
			if e == nil {
				continue
			}
			switch event := e.(type) {
			case *kafka.Message:
				log.Debug("Kafka", "Message on %s: %s", event.TopicPartition, string(event.Value))
				var value T
				err := json.Unmarshal(event.Value, &value)
				if err != nil {
					log.Error("Kafka", "Unmarshal error: %v (%v)", err, event)
				}
				values <- value
			case kafka.Error:
				log.Error("Kafka", "Consumer error: %v (%v)", event, event.Code())
				if event.Code() == kafka.ErrAllBrokersDown {
					run = false
				}
			default:
				log.Debug("Kafka", "Ignore message: %v", event)
			}
		}
	}

	close(values)
}
