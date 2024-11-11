
package main

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestFilterAdsByPriceRange(t *testing.T) {
    ads := []Ad{
        {Ad_ID: 1, SellPrice: 500000, City: "Tehran"},
        {Ad_ID: 2, SellPrice: 1500000, City: "Mashhad"},
        {Ad_ID: 3, SellPrice: 1000000, City: "Tehran"},
    }

    filter := Filter{
        StartPurchasePrice: 400000,
        EndPurchasePrice:   1200000,
        City:               "Tehran",
    }

    result := ApplyFilter(ads, filter)

    assert.Equal(t, 2, len(result))
    for _, ad := range result {
        assert.GreaterOrEqual(t, ad.SellPrice, filter.StartPurchasePrice)
        assert.LessOrEqual(t, ad.SellPrice, filter.EndPurchasePrice)
        assert.Equal(t, ad.City, filter.City)
    }
}
