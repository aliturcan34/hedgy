package main

import (
	"context"
	g "levelzero/generation"
	"time"
)

type Market struct {
	Stocks    []*g.Stock
	marketTime time.Time
}

func MakeMarket(stockIds []string, startTime time.Time) *Market {
	stocks := make([]*g.Stock, 0)
	for _, id := range stockIds {
		stock := g.MakeStock(id, startTime)
		stocks = append(stocks, stock)
	}

	return &Market{stocks, startTime}
}

func (m *Market) Run(ctx context.Context) <-chan g.Stock {
	ch := make(chan g.Stock)

	go func() {
		defer close(ch)

		// generate data from startTime to Now
		clockTime := time.Now()
		// We generate new points as long as market time < clock time which is when the request is made. We advance by 1 second every iteration
		for ; m.marketTime.Compare(clockTime) == -1; m.marketTime = m.marketTime.Add(time.Second * 1) {
			for _, stock := range m.Stocks {
				stock.Advance(m.marketTime)
				ch <- *stock
			}
		}

		// generated past points keep publishing live points
		for {
			m.marketTime = time.Now()
			select {
			case <-ctx.Done():
				return
			default:
				for _, stock := range m.Stocks {
					stock.Advance(m.marketTime)
					ch <- *stock
				}
				time.Sleep(1 * time.Second)
			}
		}
	}()

	return ch
}
