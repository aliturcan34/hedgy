package main

import (
	"context"
	g "levelzero/generation"
	"time"
)

type Market struct {
	Stocks []*g.Stock
}

func MakeMarket(stockIds []string) *Market {
	stocks := make([]*g.Stock, 0)
	for _, id := range stockIds {
		stock := g.MakeStock(id)
		stocks = append(stocks, stock)
	}

	return &Market{stocks}
}

func (m *Market) Run(ctx context.Context) <-chan g.Stock {
	ch := make(chan g.Stock)

	go func() {
		defer close(ch)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				timeStamp := time.Now()
				for _, stock := range m.Stocks {
					stock.Advance(timeStamp)
					ch <- *stock
				}
				time.Sleep(1 * time.Second)
			}
		}
	}()

	return ch
}
