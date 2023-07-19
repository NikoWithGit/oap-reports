package consumer

import (
	"oap-reposts/iface"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
)

type KafkaConsumer struct {
	sarama.Consumer
	reportRepo iface.ReportRepo
	logger     iface.Ilogger
}

func NewKafkaConsumer(brokersUrl []string, rr iface.ReportRepo, l iface.Ilogger) (*KafkaConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokersUrl, config)
	if err != nil {
		return nil, err
	}
	return &KafkaConsumer{consumer, rr, l}, nil
}

func (kConsumer *KafkaConsumer) ListenTopic(topic string) {
	consumer, err := kConsumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}
	kConsumer.logger.Info("Consumer started")
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	msgCount := 0

	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case err = <-consumer.Errors():
				kConsumer.logger.Error(err.Error())
			case msg := <-consumer.Messages():
				msgCount++
				message := string(msg.Value)
				err = kConsumer.reportRepo.AddReport(message)
				if err != nil {
					kConsumer.logger.Error(err.Error())
				}
				kConsumer.logger.Infof("Received message Count %d: | Topic(%s) | Message(%s) \n", msgCount, msg.Topic, message)
			case <-sigchan:
				kConsumer.logger.Info("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()

	<-doneCh
	kConsumer.logger.Infof("Processed %d msgCount", msgCount)

	if err := kConsumer.Close(); err != nil {
		kConsumer.logger.Panic(err.Error())
	}
}
