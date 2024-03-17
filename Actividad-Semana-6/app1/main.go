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
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // argumentos
	)
	if err != nil {
		log.Fatalf("Error al declarar una cola: %s", err)
	}

	cuerpo := "¡¡¡Probando que funciona correctamente la Aplicación 1!!!"
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(cuerpo),
		})
	if err != nil {
		log.Fatalf("Error al publicar un mensaje: %s", err)
	}
	log.Printf(" [x] Enviado %s", cuerpo)
}
