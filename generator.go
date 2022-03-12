package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

type Generator interface {
	Gen() WorkResult
}

type Work struct {
	Min              int16          // 最小值
	Max              int16          // xx以内算数
	MaxResult        int16          // 结果最大值
	Ops              []Op           // 运算
	OpCounts         []int8         // 操作符数量
	UpgradeChecker   UpgradeChecker // 是否检查进位
	UseSpecialNumber bool           // 是否使用特殊数字
	SpecialNumber    SpecialNumber  // 是否必须包含特殊数字
	ResultChecker    ResultChecker  // 结果检查器
}

func NewWork(min, max, maxResult int16, ops []string, opCounts []int8, upgradeChecker, resultChecker, specialNumber string) *Work {
	if min == 0 {
		min = 1
	}
	w := &Work{
		Min:            min,
		Max:            max,
		MaxResult:      maxResult,
		Ops:            ToOps(ops),
		OpCounts:       opCounts,
		UpgradeChecker: UpgradeCheckerMap[upgradeChecker],
		ResultChecker:  ResultCheckerMap[resultChecker],
	}
	if specialNumber != "" {
		w.UseSpecialNumber = true
		w.SpecialNumber = SpecialNumberMap[specialNumber]
	}
	return w
}

type WorkResult struct {
	Question string
	Answer   string
}

func (w *Work) Gen() WorkResult {
	var opCount int8
	if len(w.OpCounts) == 1 {
		opCount = w.OpCounts[0]
	} else {
		opCount = w.OpCounts[rand.Intn(len(w.OpCounts))]
	}
	maxLayer := opCount + 1
	var specialNumberLayer int8
	if w.UseSpecialNumber {
		// fixme: 有减法时，特殊要求的数字总生成在最后，防止死循环，简单处理了
		if contains(w.Ops, Minus) || contains(w.Ops, Divide) {
			specialNumberLayer = maxLayer
		} else {
			specialNumberLayer = randRange(1, maxLayer+1)
		}
	} else {
		specialNumberLayer = 0
	}
	var question string
	var answer int16
	var retry bool
	for {
		question, answer, retry = w.calc(1, maxLayer, specialNumberLayer, 0, "")
		if !retry {
			break
		}
	}

	// fixme: 除法余数结果展示偷懒了
	sp := strings.Split(question, " ÷ ")
	if len(sp) == 2 {
		op1, _ := strconv.Atoi(sp[0])
		op2, _ := strconv.Atoi(sp[1])
		result := op1 / op2
		mod := op1 % op2
		var sAnswer string
		if mod == 0 {
			sAnswer = strconv.Itoa(result)
		} else {
			sAnswer = fmt.Sprintf("%d...%d", result, mod)
		}
		return WorkResult{
			question,
			sAnswer,
		}
	} else {
		return WorkResult{
			question,
			strconv.Itoa(int(answer)),
		}
	}
}

func (w *Work) calc(layer, maxLayer, specialNumberLayer int8, sum int16, questionNow string) (question string, answer int16, retry bool) {
	if layer > maxLayer {
		return questionNow, sum, false
	}
	var seq []int16
	if layer == specialNumberLayer {
		seq = w.SpecialNumber.RandSeq(w.Min, w.Max)
	} else {
		seq = Normal.RandSeq(w.Min, w.Max)
	}
	opSeq := randOpSeq(w.Ops)
	// 每一层，随机尝试每个数
	for _, next := range seq {
		if layer == 1 {
			question, answer, retry = w.calc(layer+1, maxLayer, specialNumberLayer, next, fmt.Sprintf("%d", next))
			if retry {
				continue
			} else {
				return
			}
		}
		// 每一个选定数字下，尝试随机操作
		for _, nextOp := range opSeq {
			tmpSum := nextOp.Calc(sum, next)
			if w.ShouldRetry(sum, next, tmpSum, nextOp) {
				continue
			} else {
				sum = tmpSum
				question, answer, retry = w.calc(layer+1, maxLayer, specialNumberLayer, sum, fmt.Sprintf("%s %s %d", questionNow, string(nextOp), next))
				if retry {
					continue
				} else {
					return
				}
			}
		}
	}
	return question, sum, true
}

func (w *Work) ShouldRetry(sum, next, sumResult int16, op Op) bool {
	if sumResult <= 0 {
		return true
	}
	if w.MaxResult > 0 && sumResult > w.MaxResult {
		return true
	}
	sumStr := strconv.Itoa(int(sum))
	nextStr := strconv.Itoa(int(next))
	if len(sumStr) > len(nextStr) {
		nextStr = padString(nextStr, "0", len(sumStr))
	} else {
		sumStr = padString(sumStr, "0", len(nextStr))
	}

	return w.UpgradeChecker.ShouldRetry(sumStr, nextStr, op) ||
		w.ResultChecker.ShouldRetry(sum, next, sumResult, op)
}

func padString(s, pad string, toLen int) string {
	return strings.Repeat(pad, toLen-len(s)) + s
}

func randOpSeq(opCandidates []Op) (seq []Op) {
	if len(opCandidates) == 1 {
		return opCandidates
	}
	seq = make([]Op, len(opCandidates))
	copy(seq, opCandidates)
	rand.Shuffle(len(seq), func(i, j int) {
		seq[i], seq[j] = seq[j], seq[i]
	})
	return
}

// [min, max)
func randRange(min, max int8) int8 {
	return int8(rand.Intn(int(max-min)) + int(min))
}

type UpgradeChecker struct {
	ShouldRetry func(sum, next string, op Op) bool
}

type ResultChecker struct {
	ShouldRetry func(sum, next, result int16, op Op) bool
}

type SpecialNumber struct {
	RandSeq func(min, max int16) []int16
}
