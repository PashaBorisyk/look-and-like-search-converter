package queue

import (
	cluster "github.com/bsm/sarama-cluster"
	"log"
	"os"
	"os/signal"
	"sync"
)

var consumer *cluster.Consumer
var signals chan os.Signal

func InitConsumer(wg *sync.WaitGroup, doOnNext func(id string)) {
	defer wg.Done()

	log.Println("Configuring kafka...")

	// init (custom) config, enable errors and notifications
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true

	serverUrl := os.Getenv("KAFKA_SERVER_URL")
	if serverUrl == "" {
		panic("KAFKA_SERVER_URL environment variable must be provided")
	}
	brokers := []string{serverUrl}
	topics := []string{"ready_to_index_id"}
	var err error
	consumer, err = cluster.NewConsumer(brokers, "search-converter-group", topics, config)

	if err != nil {
		log.Println("Can not create kafka consumer: ", err)
		return
	}
	defer consumer.Close()

	// trap SIGINT to trigger a shutdown.
	signals = make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// consume errors
	go func() {
		for err := range consumer.Errors() {
			log.Printf("Error: %s\n", err.Error())
		}
	}()

	// consume notifications
	go func() {
		for ntf := range consumer.Notifications() {
			log.Printf("Rebalanced: %+v\n", ntf)
		}
	}()

	log.Println("Kafka configuring done")

	for {
		select {
		case msg, ok := <-consumer.Messages():
			if ok {
				doOnNext(string(msg.Value))
				consumer.MarkOffset(msg, "") // mark message as processed
			}
		case <-signals:
			return
		}
	}

	// consume messages, watch signals

}
