package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/nats-io/nats.go"
)

var resultMap = make(map[string]interface{})
var resultMapSlice = make([]map[string]interface{}, 0)

func getHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if len(resultMapSlice) > 0 {
		json.NewEncoder(w).Encode(resultMapSlice)
		resultMapSlice = make([]map[string]interface{}, 0)
	} else {
		w.Write([]byte("{\"error\": \"there is no information in queue\"}"))
	}

}

func recvNats() {
	nc, err := nats.Connect("nats://nats:4222")
	// nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalln(err)
	}
	defer nc.Close()

	count := 0
	nc.Subscribe("cbr", func(msg *nats.Msg) {
		count++
		json.Unmarshal(msg.Data, &resultMap)
		resultMapSlice = append(resultMapSlice, resultMap)
		resultMap = make(map[string]interface{})
	})

	for {
		old := count
		time.Sleep(15 * time.Second)
		if old == count {
			break
		}
	}
}

func main() {
	go recvNats()

	http.HandleFunc("/get/", getHandler)
	http.ListenAndServe(":8080", nil)
}
