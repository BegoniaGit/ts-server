package handler

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"yan.site/ts_server/config"
	"yan.site/ts_server/dao"
	"yan.site/ts_server/model"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("MQ: %s: %s", msg, err)
	}
}

type ReceiveMQ struct {
	mysqlStorage *dao.MysqlStorage
	crawlManager *CrawlManager
}

func NewReceiveMQ(mysqlStorageP *dao.MysqlStorage, crawlManagerP *CrawlManager) *ReceiveMQ {
	return &ReceiveMQ{
		mysqlStorage: mysqlStorageP,
		crawlManager: crawlManagerP,
	}
}

func (r ReceiveMQ) Start() {
	mqConfig := config.GetConf().TsServerConfig.MQ
	if mqConfig.Enable == false {
		log.Println("MQ: MQ has not start")
		return
	}
	log.Println("MQ: receive message started")
	conn, err := amqp.Dial(mqConfig.Url)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		mqConfig.Queue, // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var records []model.Record
			err := json.Unmarshal([]byte(d.Body), &records)
			r.crawlManager.SaveData(records...)
			failOnError(err, "resolve json error")
			log.Printf("MQ: Received a message: %s", d.Body)
		}
	}()
	log.Printf("MQ:  [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
