package ticker

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

const (
	url          = "wss://api.gemini.com/v1/marketdata"

	eventSideAsk = "ask"
	eventSideBid = "bid"

	eventReasonCancel = "cancel"
	eventReasonPlace  = "place"
)

type marketData struct {
	Type           string        `json:"type"`
	SocketSequence int           `json:"socket_sequence"`
	Events         []marketEvent `json:"events"`
}

type marketEvent struct {
	Type      string  `json:"type"`
	Side      string  `json:"side"`
	Price     float64 `json:"price,string"`
	Remaining float64 `json:"remaining,string"`
	Delta     float64 `json:"delta,string"`
	Reason    string  `json:"reason"`
}

// TickFormat is the format that we pass back the updated
// tick results
type TickFormat struct {
	Bid          float64 `json:"bid"`
	BidRemaining float64 `json:"bid_remaining"`
	Ask          float64 `json:"ask"`
	AskRemaining float64 `json:"ask_remaining"`
}

// Tick is the command that starts the ticking process,
// passes back a `TickFormat` struct on the passed channel
func Tick(symbol string, ch chan TickFormat) error {
	var bestBid float64
	var bestAsk float64

	bids := make(map[float64]float64)
	asks := make(map[float64]float64)

	wsUrl := fmt.Sprintf(
		"%v/%v?top_of_book=false",
		url,
		symbol,
	)

	fmt.Printf("querying %v\n", wsUrl)
	conn, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	if err != nil {
		return err
	}

	defer conn.Close()
	for {
		var marketData marketData
		isUpdated := false
		_, data, err := conn.ReadMessage()
		if err != nil {
			return err
		}

		err = json.Unmarshal(data, &marketData)
		if err != nil {
			return err
		}

		if marketData.SocketSequence == 0 {
			continue
		}

		for _, event := range marketData.Events {
			if event.Side == eventSideAsk {
				isUpdated = updateValues(event, &bestAsk, asks)
			} else {
				isUpdated = updateValues(event, &bestBid, bids)
			}
		}

		if isUpdated {
			if bestBid != 0 && bestAsk != 0 {
				ch <- TickFormat{
					Bid:          bestBid,
					BidRemaining: bids[bestBid],
					Ask:          bestAsk,
					AskRemaining: asks[bestAsk],
				}
			}
		}
	}
}

func updateValues(event marketEvent, best *float64, prices map[float64]float64) bool {
	updated := false
	if event.Reason == eventReasonPlace {
		if priceIsBest(event, *best) {
			updated = true

			*best = event.Price
			prices[event.Price] = event.Remaining
		}
	} else if event.Reason == eventReasonCancel {
		if event.Remaining == 0 {
			delete(prices, event.Price)

			// need to find new best
			if event.Price == *best {
				var newBest float64
				for price := range prices {
					if priceIsBest(event, newBest) {
						newBest = price
					}
				}
				*best = newBest
				updated = true
			}

		} else {
			// update the remaining price, if that price
			// is the best price then signify an update
			if event.Price == *best &&
				prices[event.Price] != event.Remaining {
				updated = true
			}

			prices[event.Price] = event.Remaining
		}
	}

	return updated
}

func priceIsBest(event marketEvent, best float64) bool {
	return (best == 0) ||
		(event.Side == eventSideAsk && event.Price < best) ||
		(event.Side == eventSideBid && event.Price > best)
}
