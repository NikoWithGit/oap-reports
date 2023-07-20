package consumer

import (
	"encoding/json"
	"oap-reposts/iface"
	"oap-reposts/model"

	"github.com/IBM/sarama"
)

type KafkaConsumer struct {
	sarama.Consumer
	rs     iface.ReportService
	logger iface.Ilogger
}

func NewKafkaConsumer(brokersUrl []string, rs iface.ReportService, l iface.Ilogger) (*KafkaConsumer, error) {
	config := getConfig()
	consumer, err := sarama.NewConsumer(brokersUrl, config)
	if err != nil {
		return nil, err
	}
	return &KafkaConsumer{consumer, rs, l}, nil
}

func getConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	return config
}

func (kc *KafkaConsumer) ListenTopic(topic string) error {
	consumer, err := kc.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		return err
	}
	kc.logger.Info("Consumer started")

	for {
		select {
		case err = <-consumer.Errors():
			kc.logger.Error(err.Error())

		case msg := <-consumer.Messages():
			var order model.Order
			err := json.Unmarshal(msg.Value, &order)
			if err != nil {
				kc.logger.Error(err.Error())
			}
			err = kc.rs.AddOrder(&order)
			if err != nil {
				kc.logger.Error(err.Error())
			}
			kc.logger.Infof("Received message: | Topic(%s) | Message(%s) \n", msg.Topic, string(msg.Value))
		}
	}
}
