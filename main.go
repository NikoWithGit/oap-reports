package main

import (
	"database/sql"
	"oap-reposts/consumer"
	"oap-reposts/env"
	"oap-reposts/logger"
	repoimpl "oap-reposts/repo-impl"

	_ "github.com/lib/pq"
)

func main() {

	zaplogger, err := logger.NewZapLogger()
	if err != nil {
		panic(err)
	}

	dsn := env.GetDbDsn()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		zaplogger.Panic(err.Error())
	}

	perortRepo := repoimpl.NewReportRepoImpl(db)

	consumer, err := consumer.NewKafkaConsumer([]string{"order-and-pay-kafka-1:9092"}, perortRepo, zaplogger)
	if err != nil {
		zaplogger.Panic(err.Error())
	}
	defer consumer.Close()
	consumer.ListenTopic("order-complete")
}
