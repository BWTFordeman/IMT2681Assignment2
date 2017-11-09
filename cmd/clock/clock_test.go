package main

import (
	"testing"
)

//TODO...

func TestGetFixerData(t *testing.T) {
	f := getFixerData()
	if f.BaseCurrency != "EUR" && f.BaseCurrency != "" {
		t.Errorf("Error getting fixerdata, expected EUR got:%v:", f.BaseCurrency)
	}
}
