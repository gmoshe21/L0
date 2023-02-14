package main

import (
	"L0/internal/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/nats-io/stan.go"
)

func main() {
	natsConn, err := stan.Connect(
		"test-cluster",
		"test")
	if err != nil {
		log.Fatalf("NewNatsConnect: %+v", err)
	}
	var order []models.Order

	data, err := ioutil.ReadFile("./DATA.json")
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(data, &order)
	if err != nil {
		fmt.Println(err)
	}

	for key, val := range order {
		orderBytes, _ := json.Marshal(val)
		err := natsConn.Publish("create", orderBytes)

		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("push:", key)
		time.Sleep(3000 * time.Millisecond)
	}
}
