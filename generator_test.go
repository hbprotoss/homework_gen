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
		Min:              100,
		Max:              999,
		Ops:              []Op{'÷'},
		OpCounts:         []int8{1},
		UpgradeChecker:   Ignore,
		UseSpecialNumber: true,
		SpecialNumber:    OneNumber,
		ResultChecker:    ZeroInMiddle,
	}
	for i := 0; i < 100; i++ {
		result := w.Gen()
		fmt.Printf("%s = %s\n", result.Question, result.Answer)
	}
}
