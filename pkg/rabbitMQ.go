package pkg

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"time"
)

func GetRabbitMQConnection() (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial("amqp://user:password@localhost:5672/")
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, err
	}

	return conn, ch, nil
}

type TransactionStatus struct {
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
}

func PublishTransactionStatus(noKontrak, status string) {
	conn, ch, err := GetRabbitMQConnection()
	if err != nil {
		log.Printf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	defer ch.Close()

	queueName := "transaction_status"
	_, err = ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Failed to declare queue: %v", err)
	}

	// Simulate interval execution
	ticker := time.NewTicker(1 * time.Minute) // Set interval ke 1 menit
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Data yang akan dikirim
			data := TransactionStatus{
				TransactionID: noKontrak,
				Status:        status,
			}

			body, err := json.Marshal(data)
			if err != nil {
				log.Printf("Failed to marshal data: %v", err)
				continue
			}

			err = ch.Publish(
				"",        // exchange
				queueName, // routing key
				false,     // mandatory
				false,     // immediate
				amqp.Publishing{
					ContentType: "application/json",
					Body:        body,
				},
			)
			if err != nil {
				log.Printf("Failed to publish message: %v", err)
			} else {
				log.Printf("Published message: %s", string(body))
			}
		}
	}
}

func ConsumeTransactionStatus() {
	conn, ch, err := GetRabbitMQConnection()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	defer ch.Close()

	queueName := "transaction_status"
	msgs, err := ch.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register consumer: %v", err)
	}

	for msg := range msgs {
		var status TransactionStatus
		if err := json.Unmarshal(msg.Body, &status); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		// Proses status pengajuan
		log.Printf("Received message: %+v", status)

		// Update ke database atau panggil API pengecekan
	}
}
