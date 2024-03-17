package main

import (
	"database/sql"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/streadway/amqp"
)

func main() {
	conn, err := amqp.Dial("amqp://joseph:1192948@localhost:5673/")

	if err != nil {
		log.Fatal("Error al conectar a RabbitMQ: ", err)
	}
	defer conn.Close()

	connString := "server=JosephAlmonte\\SQLEXPRESS;database=App3Db;trusted_connection=true"
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error al abrir la conexión a la base de datos: ", err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Error al hacer ping a la base de datos: ", err.Error())
	}

	consumeMessages(conn, db)

	log.Println("Aplicación 3 inicializada correctamente")
}

func consumeMessages(conn *amqp.Connection, db *sql.DB) {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Error al abrir un canal: ", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"testQueue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Error al declarar una cola: ", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Error al registrar un consumidor: ", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Mensaje recibido: %s", d.Body)
			saveMessageToDB(db, d.Body)
			sendMessageBack(conn, "Mensaje procesado correctamente")
		}
	}()

	log.Printf("Esperando mensajes. Para salir presione CTRL+C")
	<-forever
}

func saveMessageToDB(db *sql.DB, message []byte) {
	_, err := db.Exec("INSERT INTO Messages (Message) VALUES (@p1)", string(message))
	if err != nil {
		log.Fatal("Error al insertar el mensaje en la base de datos: ", err)
	}
}

func sendMessageBack(conn *amqp.Connection, responseMessage string) {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Error al abrir un canal: ", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"responseQueue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Error al declarar una cola: ", err)
	}

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(responseMessage),
		})
	if err != nil {
		log.Fatal("Error al publicar un mensaje: ", err)
	}
}
