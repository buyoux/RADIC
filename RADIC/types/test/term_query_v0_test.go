package test

import (
	"fmt"
	"strings"
	"testing"
)

// （(A|B|C)&D)|E&((F|G)&H)

// | 或运算
func should(s ...string) string {
	if len(s) == 0 {
		return ""
	}

	strBuilder := strings.Builder{}
	strBuilder.WriteString("(")
	for _, ele := range s {
		if len(ele) > 0 {
			strBuilder.WriteString(ele + "|")
		}
	}
	rect := strBuilder.String()
	return rect[0:len(rect)-1] + ")"

	// return "(" + strings.Join(s, "|") + ")"
}

// & 与运算
func must(s ...string) string {
	if len(s) == 0 {
		return ""
	}
	return "(" + strings.Join(s, "&") + ")"
}

func TestN(t *testing.T) {
	fmt.Println(must(should(must(should("A", "B", "C"), "D"), "E"), must(should("F", "G"), "H")))
}
