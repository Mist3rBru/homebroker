package main

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/Mist3rBru/homebroker/internal/infra/kafka"
	"github.com/Mist3rBru/homebroker/internal/market/dto"
	"github.com/Mist3rBru/homebroker/internal/market/entity"
	"github.com/Mist3rBru/homebroker/internal/market/transformer"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	orderIn := make(chan *entity.Order)
	orderOut := make(chan *entity.Order)

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	kafkaMsgChan := make(chan *ckafka.Message)
	configMap := &ckafka.ConfigMap{
		"bootstrap.servers": "host.docker.internal:9094",
		"group.id":          "homebroker",
		"auto.offset.reset": "earliest",
	}

	consumer := kafka.NewConsumer(configMap, []string{"input"})
	go consumer.Consume(kafkaMsgChan)

	book := entity.NewBook(orderIn, orderOut, wg)
	go book.Trade()

	go func() {
		for msg := range kafkaMsgChan {
			wg.Add(1)
			fmt.Println(string(msg.Value))
			trandeInput := dto.TradeInput{}
			err := json.Unmarshal(msg.Value, &trandeInput)
			if err != nil {
				panic(err)
			}
			order := transformer.TransformInput(trandeInput)
			orderIn <- order
		}
	}()

	producer := kafka.NewProducer(configMap)
	for res := range orderOut {
		output := transformer.TransformOutput(res)
		jsonOutput, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(jsonOutput)
		producer.Publish(jsonOutput, []byte("orders"), "output")
	}
}
