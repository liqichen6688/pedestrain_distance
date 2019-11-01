package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"pedestrain_distance/data_struct"
	"pedestrain_distance/producersetting"
)

const filepath = "video-process-OBJECT_FACE_PEDESTRIAN_ESTATE-1572343893823932889/pach_120s.mp4/output.json"

func getPeople(file string) data_struct.People {
	f, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
	}
	s := bufio.NewScanner(f)
	People := make(map[string][]data_struct.Person)
	for s.Scan() {
		var Person data_struct.Person
		if err := json.Unmarshal(s.Bytes(), &Person); err != nil {
			fmt.Println(err)
		}
		if Person.Objects.Pedestrian != nil {
			if val, ok := People[Person.Time]; ok {
				People[Person.Time] = append(val, Person)
			} else {
				People[Person.Time] = append(People[Person.Time], Person)
			}
		}
	}
	if s.Err() != nil {
		fmt.Println(s.Err())
	}
	return People
}

func getPairDistance(FirstPerson, SecondPerson data_struct.Person) data_struct.PairDistance {
	firstVertices := FirstPerson.Objects.Pedestrian.Rectangle.Vertices
	secondVertices := SecondPerson.Objects.Pedestrian.Rectangle.Vertices
	distance := math.Pow(float64(firstVertices[0].X-secondVertices[0].X), 2)
	distance += math.Pow(float64(firstVertices[0].Y-secondVertices[0].Y), 2)
	distance = math.Sqrt(distance)
	averageLength := math.Abs(float64(firstVertices[0].X-firstVertices[1].X)) + math.Abs(float64(secondVertices[0].X-secondVertices[1].X))
	averageLength = averageLength / 2
	relativeDistance := distance / averageLength
	var PairDistance data_struct.PairDistance
	PairDistance.First_ID = FirstPerson.Objects.Pedestrian.ID
	PairDistance.Second_ID = SecondPerson.Objects.Pedestrian.ID
	PairDistance.Distance = relativeDistance
	return PairDistance
}

func getDistanceTimeStamp(People data_struct.People, DisTimeChan chan data_struct.DistanceTimeStamp) {
	for key, val := range People {
		var DistanceTimeStamp data_struct.DistanceTimeStamp
		length := len(val)
		for i := 0; i < length; i++ {
			for j := i + 1; j < length; j++ {
				DistanceTimeStamp.Time = key
				DistanceTimeStamp.PairDistances = append(DistanceTimeStamp.PairDistances, getPairDistance(val[i], val[j]))
			}
		}
		DisTimeChan <- DistanceTimeStamp
	}
	close(DisTimeChan)
}

func main() {
	var ProducerConfig *producersetting.ProducerConfig
	ProducerConfig = new(producersetting.ProducerConfig)
	ProducerConfig.Topic = "testTopic"
	ProducerConfig.Brokers = append(ProducerConfig.Brokers, "localhost:9092")
	People := getPeople(filepath)
	DisTimeChan := make(chan data_struct.DistanceTimeStamp)
	go getDistanceTimeStamp(People, DisTimeChan)
	producer, err := producersetting.NewProducer(ProducerConfig)
	if err != nil {
		fmt.Println(err)
	}
	err = producersetting.ChannelToKafka(producer, DisTimeChan, ProducerConfig.Topic, producersetting.MaxRetryTimes)
	if err == nil {
		fmt.Println(err)
	}
}
