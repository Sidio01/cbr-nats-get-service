package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/nats-io/nats.go"
	"golang.org/x/text/encoding/charmap"
)

type currency struct {
	CharCode string `xml:"CharCode"`
	Nominal  int    `xml:"Nominal"`
	Value    string `xml:"Value"`
}

type currenciesList struct {
	Valute []currency `xml:"Valute"`
	Date   string     `xml:"Date,attr"`
}

var sendCh = make(chan interface{})

func main() {
	nc, _ := nats.Connect("nats://nats:4222")
	// nc, err := nats.Connect(nats.DefaultURL)
	ec, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	defer ec.Close()
	ec.BindSendChan("cbr", sendCh)

	for {
		time.Sleep(5 * time.Second)

		resp, err := http.Get("https://www.cbr.ru/scripts/XML_daily.asp")
		if err != nil {
			log.Fatalln(err)
		}

		data := xml.NewDecoder(resp.Body)
		data.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
			switch charset {
			case "windows-1251":
				return charmap.Windows1251.NewDecoder().Reader(input), nil
			default:
				return nil, fmt.Errorf("unknown charset: %s", charset)
			}
		}

		currencies := new(currenciesList)
		err = data.Decode(&currencies)
		if err != nil {
			log.Fatalln(err)
		}

		currenciesResult := make([]currency, 0)
		for _, curr := range currencies.Valute {
			if curr.CharCode == "USD" || curr.CharCode == "EUR" {
				currenciesResult = append(currenciesResult, curr)
			}
		}

		result := map[string]interface{}{"parse_time": time.Now().Format("2006-01-02 15:04:05"), "currencies": currenciesResult}

		sendCh <- result
		resp.Body.Close()
	}
}
