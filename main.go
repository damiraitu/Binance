package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)

type Response struct {
	LastUpdateId int32      `json:"lastUpdateId"`
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
}

func home(w http.ResponseWriter, r *http.Request) {
	conn, res, err := websocket.DefaultDialer.Dial("wss://stream.binance.com/ws/manabtc@depth20", nil)

	if res.StatusCode == http.StatusSwitchingProtocols{
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		return
	}
	// fmt.Println(conn)
	var sumOfBids float64
	var sumOfAsks float64

	for {
		sumOfBids = 0
		sumOfAsks = 0
		resp := &Response{}
		readErr := conn.ReadJSON(resp)
		if readErr != nil {
			log.Println(readErr)
			return
		}

		w.Write([]byte("------------------------ NEW RESPONSE ------------------------- \n"))

		for i := 0; i < 15; i++ {
			price, err := strconv.ParseFloat(resp.Bids[i][0], 64)
			if err != nil {
				log.Println(err)
				return
			}
			qty, err := strconv.ParseFloat(resp.Bids[i][1], 64)
			if err != nil {
				log.Println(err)
				return
			}
			totalPrice := price * qty
			sumOfBids += totalPrice

			priceAsk, err := strconv.ParseFloat(resp.Asks[i][0], 64)
			if err != nil {
				log.Println(err)
				return
			}
			qtyAsk, err := strconv.ParseFloat(resp.Asks[i][1], 64)
			if err != nil {
				log.Println(err)
				return
			}
			totalPriceAsk := priceAsk * qtyAsk
			sumOfAsks += totalPriceAsk
		}
		fmt.Fprintf(w, "Sum of BIDS: %f\n", sumOfBids)
		fmt.Fprintf(w, "Sum of ASKS: %f\n", sumOfAsks)

		w.Write([]byte("------------------------ END OF RESPONSE ------------------------- \n"))
	}

}


func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
