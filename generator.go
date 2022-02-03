package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

type Generator interface {
	Gen(questions map[string]bool) WorkResult
}

type Work struct {
	Max            int8           // xx以内算数
	Ops            []Op           // 运算
	OpCounts       []int8         // 操作符数量
	UpgradeChecker UpgradeChecker // 是否检查进位
}

type WorkResult struct {
	Question string
	Answer   int16
}

func (w *Work) Gen(questions map[string]bool) WorkResult {
	var opCount int8
	if len(w.OpCounts) == 1 {
		opCount = w.OpCounts[0]
	} else {
		opCount = w.OpCounts[rand.Intn(len(w.OpCounts))]
	}
	question, answer, _ := w.calc(1, opCount+1, w.Max, 0, "", questions)
	return WorkResult{
		question,
		answer,
	}
}

func (w *Work) calc(layer, maxLayer, maxElement int8, sum int16, questionNow string, questions map[string]bool) (question string, answer int16, retry bool) {
	if layer > maxLayer {
		_, ok := questions[questionNow]
		if !ok {
			questions[questionNow] = true
		}
		return questionNow, sum, ok // 已经计算过了，需要重新计算
	}
	seq := randSeq(maxElement)
	opSeq := randOpSeq(w.Ops)
	// 每一层，随机尝试每个数
	for _, next := range seq {
		if layer == 1 {
			question, answer, retry = w.calc(layer+1, maxLayer, maxElement, int16(next), fmt.Sprintf("%d", next), questions)
			if retry {
				continue
			} else {
				return
			}
		}
		// 每一个选定数字下，尝试随机操作
		for _, nextOp := range opSeq {
			tmpSum := nextOp.Calc(sum, int16(next))
			if w.ShouldRetry(sum, int16(next), tmpSum, nextOp) {
				continue
			} else {
				sum = tmpSum
				question, answer, retry = w.calc(layer+1, maxLayer, maxElement, sum,
					fmt.Sprintf("%s %s %d", questionNow, string(nextOp), next),
					questions)
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
	sumStr := strconv.Itoa(int(sum))
	nextStr := strconv.Itoa(int(next))
	if len(sumStr) > len(nextStr) {
		nextStr = padString(nextStr, "0", len(sumStr))
	} else {
		sumStr = padString(sumStr, "0", len(nextStr))
	}

	return w.UpgradeChecker.ShouldRetry(sumStr, nextStr, op)
}

func padString(s, pad string, toLen int) string {
	return strings.Repeat(pad, toLen-len(s)) + s
}

func randSeq(max int8) (seq []int8) {
	seq = make([]int8, max)
	for i := int8(1); i <= max; i++ {
		seq[i-1] = i
	}
	rand.Shuffle(len(seq), func(i, j int) {
		seq[i], seq[j] = seq[j], seq[i]
	})
	return
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

type UpgradeChecker struct {
	ShouldRetry func(sum, next string, op Op) bool
}
