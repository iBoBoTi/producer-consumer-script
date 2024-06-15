package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/iBoBoTi/woodcore-task/models"
	amqp "github.com/rabbitmq/amqp091-go"
)



func main(){
	queueConn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("error connecting to rabbitmq")
	}
	defer queueConn.Close()

	jsonFile, err := os.Open("data.json")
    if err != nil {
        log.Println(fmt.Errorf("error opening json file: %v", err))
        return
    }
    defer jsonFile.Close()

    // Read the opened file into a byte slice
    byteValue, err := io.ReadAll(jsonFile)
    if err != nil {
        log.Println(fmt.Errorf("error reading json file: %v", err))
        return
    }

    var payload models.Payload
    if err := json.Unmarshal(byteValue, &payload); err != nil {
        log.Println(fmt.Errorf("error unmarshalling json: %v", err))
        return
    }
	produce(queueConn, &payload)
}

func produce(conn *amqp.Connection, payload *models.Payload){
	ch, err :=  conn.Channel()
	if err != nil {
		log.Println(fmt.Errorf("error opening channel: %v", err))
		return
	}
	defer ch.Close()

	q, err :=  ch.QueueDeclare(
		"payload_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println(fmt.Errorf("error declaring queue: %v", err))
		return
	}

	body, err := json.Marshal(payload)
	if err != nil {
		log.Println(fmt.Errorf("error marshalling payload: %v", err))
		return
	}

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body: body,
		},
	)
	if err != nil{
		log.Println(fmt.Errorf("error publishing body: %v", err))
		return
	}

	log.Println("payload published")
}