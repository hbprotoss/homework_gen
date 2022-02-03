package main

import (
	"fmt"
	"testing"
)

func TestGen(t *testing.T) {
	config := initYamlConfig()
	w := config[0].WorkConfigs[9].Gen
	questions := make(map[string]bool)
	for i := 0; i < 10; i++ {
		result := w.Gen(questions)
		fmt.Printf("%s = %d\n", result.Question, result.Answer)
	}
}