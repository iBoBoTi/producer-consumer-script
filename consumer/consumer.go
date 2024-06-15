// consumer.go
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/iBoBoTi/woodcore-task/models"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
)


func main() {
	queueConn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("error connecting to rabbitmq")
	}
	defer queueConn.Close()


	db, err:=initDB()
	if err != nil {
		log.Fatal(err)
	}
    defer db.Close()

    ch, err := queueConn.Channel()
	if err != nil {
		log.Println(fmt.Errorf("error opening channel: %v", err))
		return
	}
    defer ch.Close()

    q, err := ch.QueueDeclare(
        "payload_queue", // name
        false,           // durable
        false,           // delete when unused
        false,           // exclusive
        false,           // no-wait
        nil,             // arguments
    )
	if err != nil {
		log.Println(fmt.Errorf("error declaring queue: %v", err))
		return
	}

    msgs, err := ch.Consume(
        q.Name, // queue
        "",     // consumer
        true,   // auto-ack
        false,  // exclusive
        false,  // no-local
        false,  // no-wait
        nil,    // args
    )
	if err != nil{
		log.Println(fmt.Errorf("error creating consumer: %v", err))
		return
	}

    channel := make(chan bool)

    go func() {
        for d := range msgs {
            var payload models.Payload
            err := json.Unmarshal(d.Body, &payload)
			
            if err != nil {
                log.Printf("Failed to decode message: %s", err)
                continue
            }
            consumePayload(db, &payload)
        }
    }()

    log.Printf("Waiting for messages. To exit press CTRL+C")
    <-channel
}

func consumePayload(db *sql.DB, payload *models.Payload) {
    switch payload.Action {
    case models.CreateAction:
		err := CreateUser(db, payload)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("user created successfully")
		
    case models.DeleteAction:
		err := DeleteUser(db, payload)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("user deleted successfully")
		
	default:
		log.Println("Unknown Payload Action")

    }
}

func initDB() (*sql.DB, error){
	err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

	host := os.Getenv("DB_HOST")
	port := 5432
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
	host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(fmt.Errorf("error connecting to db: %v", err))
	}
	return db, nil
}

func CreateUser(db *sql.DB, payload *models.Payload) error {
		stmt, err := db.Prepare("INSERT INTO users(name, age) VALUES($1, $2)")
		if err != nil {
			return fmt.Errorf("error preparing insert statement: %v", err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(payload.Data.Name, payload.Data.Age)
		if err != nil {
			return fmt.Errorf("error inserting data: %v", err)
		}
		return nil
}

func DeleteUser(db *sql.DB, payload *models.Payload) error {
	_, err := db.Exec("DELETE FROM users WHERE name=$1", payload.Data.Name)
	if err != nil {
		return fmt.Errorf("failed to delete user: %s", err)
	} 
	return nil
}