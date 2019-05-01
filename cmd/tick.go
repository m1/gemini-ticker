package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/m1/gemini-ticker/pkg/ticker"
)

var (
	tickCmd = cobra.Command{
		Use:   "tick [symbol, DEFAULT=btcusd]",
		Short: "tick",
		Long:  "tick is for ticking over the market for the symbol passed",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return nil
			}

			if isValidSymbol(args[0]) {
				return nil
			}

			return fmt.Errorf("invalid symbol specified: %s", args[0])

		},
		Run: tick,
	}

	validSymbols = []string{
		"btcusd",
		"ethusd",
		"ethbtc",
		"zecusd",
		"zecbtc",
		"zeceth",
		"zecbch",
		"zecltc",
		"bchusd",
		"bchbtc",
		"bcheth",
		"ltcusd",
		"ltcbtc",
		"ltceth",
		"ltcbch",
	}
)

func tick(_ *cobra.Command, args []string) {
	// default symbol
	symbol := validSymbols[0]
	if len(args) > 0 {
		symbol = args[0]
	}

	ch := make(chan ticker.TickFormat)
	go ticker.Tick(symbol, ch)
	for {
		v, ok := <-ch
		if !ok {
			return
		}

		fmt.Printf("%v %v - %v %v \n",
			v.Bid,
			v.BidRemaining,
			v.Ask,
			v.AskRemaining,
		)
	}
}

func isValidSymbol(s string) bool {
	for _, symbol := range validSymbols {
		if symbol == strings.ToLower(s) {
			return true
		}
	}

	return false
}
