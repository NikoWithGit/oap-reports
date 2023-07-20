package main

import (
	"oap-reposts/consumer"
	"oap-reposts/controller"
	"oap-reposts/db"
	"oap-reposts/iface"
	"oap-reposts/logger"
	repoimpl "oap-reposts/repo-impl"
	"oap-reposts/server"
	"oap-reposts/service"

	_ "github.com/lib/pq"
)

func main() {

	zaplogger, err := logger.NewZapLogger()
	if err != nil {
		panic(err)
	}

	db, err := db.NewSqlDb()
	if err != nil {
		zaplogger.Panic(err.Error())
	}
	defer db.Close()

	reportRepo := repoimpl.NewReportRepoImpl(db)
	reportService := service.NewReportService(reportRepo)

	errChan := make(chan error, 2)

	go func() {
		errChan <- kafkaInitAndListen(reportService, zaplogger)
	}()

	go func() {
		errChan <- serverInitAndListen(reportService, zaplogger)
	}()

	for i := 0; i < cap(errChan); i++ {
		if err = <-errChan; err != nil {
			zaplogger.Error(err.Error())
		}
	}
}

func kafkaInitAndListen(rs iface.ReportService, logger iface.Ilogger) error {
	consumer, consumerErr := consumer.NewKafkaConsumer([]string{"order-and-pay-kafka-1:9092"}, rs, logger)
	if consumerErr != nil {
		return consumerErr
	}
	defer consumer.Close()

	consumerErr = consumer.ListenTopic("order-complete")
	return consumerErr
}

func serverInitAndListen(rs iface.ReportService, logger iface.Ilogger) error {
	controller := controller.NewReportController(rs, logger)
	server := server.NewServer()
	server.RegisterRoutes(controller)
	serverErr := server.Start()
	return serverErr
}
