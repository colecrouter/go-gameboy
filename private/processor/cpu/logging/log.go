package logging

import (
	"fmt"
	"os"
)

type LogState struct {
	pc       uint16
	sp       uint16
	mnemonic string
}

func (ls LogState) String() string {
	return fmt.Sprintf("PC: 0x%04X, SP: 0x%04X, OP: %s", ls.pc, ls.sp, ls.mnemonic)
}

type Logger struct {
	logs []LogState
	f    *os.File
}

func (l *Logger) Log(pc, sp uint16, mnemonic string) {
	l.logs = append(l.logs, LogState{pc: pc, sp: sp, mnemonic: mnemonic})
}

// --- Node types for hierarchical representation ---
type Node interface{}
type LineNode struct{ Line string }
type RepeatNode struct {
	Times int
	Block []Node
}

// helper to check if a block of logs starting at offsetA equals that at offsetB
func sameBlock(input []LogState, offsetA, offsetB, k int) bool {
	for i := 0; i < k; i++ {
		if input[offsetA+i].mnemonic != input[offsetB+i].mnemonic {
			return false
		}
	}
	return true
}

// flushNodes builds a node tree with repeated blocks collapsed
func flushNodes(input []LogState) []Node {
	var nodes []Node
	n, i := len(input), 0
	for i < n {
		var found bool
		var candidateK, candidateCount int
		// maximum candidate block is half the remaining slice
		for k := (n - i) / 2; k >= 1; k-- {
			count := 1
			for i+count*k+k <= n && sameBlock(input, i, i+count*k, k) {
				count++
			}
			if count > 1 {
				found, candidateK, candidateCount = true, k, count
				break
			}
		}
		if found {
			block := flushNodes(input[i : i+candidateK])
			nodes = append(nodes, RepeatNode{Times: candidateCount, Block: block})
			i += candidateCount * candidateK
		} else {
			nodes = append(nodes, LineNode{Line: input[i].String()})
			i++
		}
	}
	return nodes
}

// flattenNodes merges nested RepeatNodes where possible
func flattenNodes(nodes []Node) []Node {
	var out []Node
	for _, n := range nodes {
		switch node := n.(type) {
		case RepeatNode:
			flat := flattenNodes(node.Block)
			// If the inner block is exactly one RepeatNode, combine counts.
			if len(flat) == 1 {
				if inner, ok := flat[0].(RepeatNode); ok {
					out = append(out, RepeatNode{Times: node.Times * inner.Times, Block: inner.Block})
					continue
				}
			}
			out = append(out, RepeatNode{Times: node.Times, Block: flat})
		default:
			out = append(out, n)
		}
	}
	return out
}

// renderNodes returns the final output lines.
func renderNodes(nodes []Node) []string {
	var out []string
	for _, n := range nodes {
		switch node := n.(type) {
		case LineNode:
			out = append(out, node.Line)
		case RepeatNode:
			sub := renderNodes(node.Block)
			total := len(sub) * node.Times
			header := fmt.Sprintf("START REPEAT: block repeated %d times (%d instructions)", node.Times, total)
			out = append(out, header)
			out = append(out, sub...)
			out = append(out, "END REPEAT")
		}
	}
	return out
}

// Flush groups by SP and renders nodes to file.
func (l *Logger) Flush() {
	if len(l.logs) == 0 {
		return
	}
	if l.f == nil {
		f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return
		}
		l.f = f
	}

	currentSP := l.logs[0].sp
	var group []LogState
	for _, entry := range l.logs {
		if entry.sp == currentSP {
			group = append(group, entry)
		} else {
			nodes := flattenNodes(flushNodes(group))
			for _, line := range renderNodes(nodes) {
				fmt.Fprintln(l.f, line)
			}
			currentSP = entry.sp
			group = []LogState{entry}
		}
	}
	nodes := flattenNodes(flushNodes(group))
	for _, line := range renderNodes(nodes) {
		fmt.Fprintln(l.f, line)
	}
	l.logs = l.logs[:0]
}
