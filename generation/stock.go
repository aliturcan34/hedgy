// Package generation is responsible for generating stock prices
package generation

import (
	"fmt"
	"levelzero/protos/stockpb"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type Stock struct {
	Data *stockpb.Stock
	gen  Generator
}

// TODO: make volatility an arg in the future
func MakeStock(id string) *Stock {
	data := stockpb.Stock{
		Id: id, TimeStamp: timestamppb.New(time.Now()),
		Last: 0, Volume: 0, TotalVolume: 0, Volatility: 0.000082,
	}
	stock := Stock{Data: &data, gen: MakeGenerator()}
	stock.Data.Last = stock.gen.GetInitialPrice()
	return &stock
}

func (s *Stock) Advance(timeStamp time.Time) {
	// our change percentage will be a normal distribution
	changeP := s.gen.GetRandomPriceChange()
	vol := s.gen.GetRandomVolume()

	s.Data.Last = s.Data.Last * (1 + changeP)
	s.Data.Volume = int32(vol)
	s.Data.TotalVolume += int32(vol)
	s.Data.TimeStamp = timestamppb.New(timeStamp)
}

func (s Stock) String() string {
	return fmt.Sprintf("%s | Price: %f | Volume: %d | TotalVolume | %d", s.Data.Id, s.Data.Last, s.Data.Volume, s.Data.TotalVolume)
}
