package main

import (
	"math/rand"
	"strconv"
	"strings"
)

const OutDir = "works"

type Op byte

const (
	Plus     Op = '+'
	Minus    Op = '-'
	Multiply Op = '*'
	Divide   Op = '÷'
)

func (op Op) Calc(op1, op2 int16) int16 {
	switch op {
	case Plus:
		return op1 + op2
	case Minus:
		return op1 - op2
	case Multiply:
		return op1 * op2
	case Divide:
		return op1 / op2
	default:
		panic("invalid op")
	}
}

func ToOps(ops []string) []Op {
	var result []Op
	for _, op := range ops {
		result = append(result, toOp(op))
	}
	return result
}

func toOp(op string) Op {
	switch op {
	case "+":
		return Plus
	case "-":
		return Minus
	case "*":
		return Multiply
	case "÷":
		return Divide
	}
	panic("invalid op")
}

func contains(s []Op, e Op) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

var (
	NoUpgrade = UpgradeChecker{ShouldRetry: func(sum, next string, op Op) bool {
		for i := range sum {
			n1, _ := strconv.Atoi(string(sum[i]))
			n2, _ := strconv.Atoi(string(next[i]))
			switch op {
			case Plus:
				if n1+n2 >= 10 {
					return true
				}
			case Minus:
				if n1-n2 < 0 {
					return true
				}
			}
		}
		return false
	}}
	Upgrade = UpgradeChecker{ShouldRetry: func(sum, next string, op Op) bool {
		var ret = true
		for i := range sum {
			n1, _ := strconv.Atoi(string(sum[i]))
			n2, _ := strconv.Atoi(string(next[i]))
			switch op {
			case Plus:
				if n1+n2 >= 10 {
					ret = false
					break
				}
			case Minus:
				if n1-n2 < 0 {
					ret = false
					break
				}
			}
		}
		return ret
	}}
	Ignore = UpgradeChecker{ShouldRetry: func(sum, next string, op Op) bool {
		return false
	}}
)

var UpgradeCheckerMap = map[string]UpgradeChecker{
	"noUpgrade": NoUpgrade,
	"upgrade":   Upgrade,
	"ignore":    Ignore,
}

var (
	ZeroInMiddle = ResultChecker{ShouldRetry: func(sum, next, result int16, op Op) bool {
		r := strconv.Itoa(int(result))
		if len(r) < 3 {
			return true
		}
		return !strings.Contains(r[1:len(r)-1], "0")
	}}
)

var ResultCheckerMap = map[string]ResultChecker{
	"zeroInMiddle": ZeroInMiddle,
}

var (
	Normal = SpecialNumber{RandSeq: func(min, max int16) []int16 {
		seq := make([]int16, max-min+1)
		for i := min; i <= max; i++ {
			seq[i-min] = i
		}
		rand.Shuffle(len(seq), func(i, j int) {
			seq[i], seq[j] = seq[j], seq[i]
		})
		return seq
	}}
	TenMultiple = SpecialNumber{RandSeq: func(min, max int16) []int16 {
		multiple := int16(10)
		seq := make([]int16, max/multiple)
		for i := int16(1); i <= max/multiple; i++ {
			seq[i-1] = i * multiple
		}
		rand.Shuffle(len(seq), func(i, j int) {
			seq[i], seq[j] = seq[j], seq[i]
		})
		return seq
	}}
	OneNumber = SpecialNumber{RandSeq: func(min, max int16) []int16 {
		min = 1
		max = 9
		seq := make([]int16, max-min+1)
		for i := min; i <= max; i++ {
			seq[i-min] = i
		}
		rand.Shuffle(len(seq), func(i, j int) {
			seq[i], seq[j] = seq[j], seq[i]
		})
		return seq
	}}
	CompositeNumber = SpecialNumber{RandSeq: func(min, max int16) []int16 {
		if rand.Intn(2) == 0 {
			return TenMultiple.RandSeq(min, max)
		} else {
			return OneNumber.RandSeq(min, max)
		}
	}}
)

var SpecialNumberMap = map[string]SpecialNumber{
	"normal":      Normal,
	"tenMultiple": TenMultiple,
	"oneNumber":   OneNumber,
	"composite":   CompositeNumber,
}
