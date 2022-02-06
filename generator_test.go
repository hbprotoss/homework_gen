package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestGen(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	//config := initYamlConfig()
	//w := config[0].WorkConfigs[1].Gen
	w := &Work{
		Min:              10,
		Max:              99,
		MaxResult:        99,
		Ops:              []Op{'+'},
		OpCounts:         []int8{1},
		UpgradeChecker:   NoUpgrade,
		UseSpecialNumber: true,
		SpecialNumber:    OneNumber,
	}
	for i := 0; i < 100; i++ {
		result := w.Gen()
		fmt.Printf("%s = %d\n", result.Question, result.Answer)
	}
}
