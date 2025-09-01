package generation

import (
	"math"
	"math/rand/v2"
)

const (
	priceVolatility  = 0.000082
	volumeVolatility = 50
	volumeMean       = 100

	// min and max starting prices. we randomly generate a starting price also
	minPrice = 100
	maxPrice = 1000
)

type Generator struct {
	random *rand.Rand
}

// MakeGenerator creates a Generator with a random number generator
func MakeGenerator() Generator {
	return Generator{
		random: rand.New(rand.NewPCG(1, 2)),
	}
}

// GetInitialPrice gets an initial price => minPrice <= price <= maxPrice
func (g Generator) GetInitialPrice() float64 {
	return minPrice + rand.Float64()*(maxPrice-minPrice)
}

// GetRandomPriceChange creates a random normal number and scales it with the preset volatility
func (g Generator) GetRandomPriceChange() float64 {
	return g.random.NormFloat64() * priceVolatility
}

// GetRandomVolume  generates a volume similar to a price.
// We can't have 0 volume so if the generated number is less than 0 we give it the lowest volume possible = 1
func (g Generator) GetRandomVolume() int {
	vol := int(math.Floor(g.random.NormFloat64()*float64(volumeVolatility) + volumeMean))
	// low chance that vol < 0, put a small volume if this is the case vol can't be negative
	if vol <= 0 {
		vol = 1
	}
	return vol
}
