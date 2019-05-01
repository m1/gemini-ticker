# gemini-ticker

A CLI ticker for the Gemini currency exchange. Can also be used 
as a Go package

## Usage

### CLI

```
➜  gemini-ticker tick --help
tick is for ticking over the market for the symbol passed

Usage:
   tick [symbol, DEFAULT=btcusd] [flags]

Flags:
  -h, --help   help for tick
```

For example:

```
➜  gemini-ticker tick
querying wss://api.gemini.com/v1/marketdata/btcusd?top_of_book=false
5287.26 1 - 5288.19 2 
5287.25 8.00320915 - 5288.19 2 
5287.26 0.0097 - 5288.19 2 
```

Another example:

```
➜  gemini-ticker tick ethusd        
querying wss://api.gemini.com/v1/marketdata/ethusd?top_of_book=false
158.27 242.1952 - 158.52 50 
158.27 198 - 158.52 50 
158.27 167 - 158.52 50 

```

### Package

```go

package main

import (
	"fmt"
	
	"github.com/m1/gemini-ticker/pkg/ticker"
)

func main() {
	ch := make(chan ticker.TickFormat)
    	go ticker.Tick("btcusd", ch)
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
```