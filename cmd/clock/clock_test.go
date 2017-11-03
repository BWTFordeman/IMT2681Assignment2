package main

import (
	"testing"
	"time"
)

//TODO...

func TestGetFixerData(t *testing.T) {
	k := time.Now()
	f := getFixerData()
	if f.BaseCurrency != "EUR" || f.BaseCurrency == "" {
		t.Errorf("Error getting fixerdata, expected EUR got:%v", f.BaseCurrency)
	}
	if f.Date != k.Format("2006-01-02") {
		t.Errorf("Error getting fixerdata, expected %v got %v", k.Format("2006-01-02"), f.Date)
	}
}
