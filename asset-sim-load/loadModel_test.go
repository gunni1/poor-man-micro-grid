package main

import (
	"fmt"
	"testing"
	"time"
)

func TestLoadModel(t *testing.T) {
	model := LoadModel{
		P:     350,
		Base:  300,
		Peak:  700,
		Alpha: 0.08,
		Sigma: 20,
	}

	start := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 24*12; i++ {
		ts := start.Add(time.Duration(i) * 5 * time.Minute)
		p := model.Step(ts)
		fmt.Printf("%s: %.2f kW\n", ts.Format(time.RFC3339), p)
		if p < 0 {
			t.Errorf("Power should not be negative at %s, got %f", ts.Format(time.RFC3339), p)
		}
	}
}
