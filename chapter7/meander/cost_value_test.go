package meander_test

import (
	"github.com/Sharykhin/blueprint/chapter7/meander"
	"github.com/cheekybits/is"
	"testing"
)

func TestParseCost(t *testing.T) {
	is := is.New(t)
	is.Equal(meander.Cost1, meander.ParseCost("$"))
	is.Equal(meander.Cost2, meander.ParseCost("$$"))
	is.Equal(meander.Cost3, meander.ParseCost("$$$"))
	is.Equal(meander.Cost4, meander.ParseCost("$$$$"))
	is.Equal(meander.Cost5, meander.ParseCost("$$$$$"))
}

func TestParseCostRange(t *testing.T) {
	is := is.New(t)
	var l meander.CostRange
	var err error
	l, err = meander.ParseCostRange("$$...$$$")
	is.NoErr(err)
	is.Equal(l.From, meander.Cost2)
	is.Equal(l.To, meander.Cost3)
	l, err = meander.ParseCostRange("$...$$$$$")
	is.NoErr(err)
	is.Equal(l.From, meander.Cost1)
	is.Equal(l.To, meander.Cost5)
}

func TestCostRangeString(t *testing.T) {
	is := is.New(t)
	r := meander.CostRange{
		From: meander.Cost2,
		To:   meander.Cost4,
	}
	is.Equal("$$...$$$$", r.String())
}
