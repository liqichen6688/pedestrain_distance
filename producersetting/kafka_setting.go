package producersetting

import (
	cluster "github.com/bsm/sarama-cluster"
	"../../../Documents/pedestrain_distance/Main"
	"time"
	log "github.com/sirupsen/logrus"
)

var(
	MaxRetryTimes = 9
)

func NewProducer(config *ProducerConfig) (*kafka.Producer, error){
	kafkaConfig := cluster.NewConfig()
	producer, err := kafka.NewProducer(config.Brokers, kafkaConfig)
	if err != nil {
		return nil, err
	}
	return *producer, err
}

func ChannelToKafka(p *kafka.Producer, c chan Main.DistanceTimeStamp, topic string, retry int) error {
	var err error
	for timestamp := range c{
		respite := time.Duration(1) * time.Second
		for i := 0; i < retry; i++ {
			if err = p.WriteMessage(topic, timestamp); err != nil {
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

