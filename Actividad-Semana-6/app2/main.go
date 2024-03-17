package main

import (
	"log"

	"github.com/streadway/amqp"
)

func main() {
	conn, err := amqp.Dial("amqp://joseph:1192948@localhost:5673/")
	if err != nil {
		log.Fatalf("Error al conectar con RabbitMQ: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error al abrir un canal: %s", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"testQueue", // nombre
		true,        // durable: ahora se establece como true para coincidir
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // argumentos
	)

	if err != nil {
		log.Fatalf("Error al declarar una cola: %s", err)
	}

	msgs, err := ch.Consume(
		q.Name, // cola
		"",     // consumidor
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // argumentos
	)
	if err != nil {
		log.Fatalf("Error al registrar un consumidor: %s", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Mensaje recibido: %s", d.Body)
		}
	}()

	log.Printf(" [*] Esperando mensajes. Para salir presiona CTRL+C")
	<-forever
}
