package producersetting

import (
	"encoding/json"
	"fmt"
	cluster "github.com/bsm/sarama-cluster"
	log "github.com/sirupsen/logrus"
	"gitlab.sz.sensetime.com/rd-platform/public/strategy-service/utils/kafka"
	"time"
)

var (
	MaxRetryTimes = 9
)

func NewProducer(config *ProducerConfig) (*kafka.Producer, error) {
	kafkaConfig := cluster.NewConfig()
	producer, err := kafka.NewProducer(config.Brokers, kafkaConfig)
	if err != nil {
		return nil, err
	}
	return producer, err
}

func ChannelToKafka(p *kafka.Producer, c chan interface{}, topic string, retry int) error {
	var err error
	for timestamp := range c {
		respite := time.Duration(1) * time.Second
		for i := 0; i < retry; i++ {
			jsontime, jsonerr := json.Marshal(timestamp)
			if jsonerr != nil {
				fmt.Println(jsonerr)
			}
			if err = p.WriteMessage(topic, jsontime); err != nil {
				log.Warnf("retry %v times, data coordinate error: %v", retry, err)
				time.Sleep(respite)
				respite *= 2
			} else {
				break
			}
		}
	}
	return err
}
