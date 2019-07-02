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
	log.Println("InitConsumer kafka subscriber")

	// init (custom) config, enable errors and notifications
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true

	// init consumer
	brokers := []string{"localhost:9092"}
	topics := []string{"ready_to_index_id"}
	var err error
	consumer, err = cluster.NewConsumer(brokers, "akka_streams_group", topics, config)

	if err != nil {
		panic(err)
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

	log.Println("Subscribe called")

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

func Subscribe(doOnNext func(id string)) {

}
