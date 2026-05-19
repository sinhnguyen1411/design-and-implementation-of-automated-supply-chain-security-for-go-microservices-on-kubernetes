package datafeed

import "time"

type Tick struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
	Volume float64 `json:"volume"`
	TS     int64   `json:"ts"`
}

type Bar struct {
	Symbol string  `json:"symbol"`
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume float64 `json:"volume"`
}

func NewTick(symbol string, price, volume float64) Tick {
	return Tick{Symbol: symbol, Price: price, Volume: volume, TS: time.Now().UnixMilli()}
}

func AggregateToBar(ticks []Tick) Bar {
	if len(ticks) == 0 { return Bar{} }
	b := Bar{Symbol: ticks[0].Symbol, Open: ticks[0].Price, Close: ticks[len(ticks)-1].Price, High: ticks[0].Price, Low: ticks[0].Price}
	for _, t := range ticks {
		if t.Price > b.High { b.High = t.Price }
		if t.Price < b.Low { b.Low = t.Price }
		b.Volume += t.Volume
	}
	return b
}
