package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/shyam81992/Inventory-Management-Job/models"

	"github.com/shyam81992/Inventory-Management-Job/config"
	"github.com/shyam81992/Inventory-Management-Job/db"
	"github.com/streadway/amqp"
)

func initialize() {

	forever := make(chan bool)

	c := make(chan *amqp.Error)
	go func() {
		err := <-c
		log.Println("reconnect: " + err.Error())
		time.Sleep(time.Duration(60) * time.Second)
		initialize()
		forever <- true

	}()

	conn, err := amqp.Dial(config.RabbitConfig["uri"])
	if err != nil {
		fmt.Println(err, "Failed to connect to RabbitMQ")
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err, "Failed to register a consumer")
		return
	}
	defer ch.Close()

	conn.NotifyClose(c)
	q, err := ch.QueueDeclare(
		config.RabbitConfig["queuename"], // name
		true,                             // durable
		false,                            // delete when unused
		false,                            // exclusive
		false,                            // no-wait
		nil,                              // arguments
	)
	if err != nil {
		fmt.Println(err, "Failed to register a consumer")
		return
	}
	QueueListen(ch, q.Name)
	<-forever
}

func QueueListen(ch *amqp.Channel, name string) {
	msgs, err := ch.Consume(
		name,  // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		fmt.Println(err, "Failed to register a consumer")
		// close the channel and handle the channel close event
		return
	}

	go func() {
		for d := range msgs {
			fmt.Println("Received a message:", string(d.Body))
			processmsg(d.Body)
			log.Printf("Done")
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

}

func processmsg(msg []byte) {

	var shipment models.Shipment

	if err := json.Unmarshal(msg, &shipment); err != nil {
		fmt.Println(err, "Invalid msg", string(msg))
		return
	}
	fmt.Println("shipment msg", shipment)

	// sendmail functionality
}

func main() {
	config.LoadConfig()
	db.InitDb()

	initialize()

}
