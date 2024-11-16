package types

import "strings"

type TermQueryV0 struct {
	Must    []TermQueryV0
	Should  []TermQueryV0
	Keyword string
}

func (exp TermQueryV0) Empty() bool {
	return len(exp.Keyword) == 0 && len(exp.Must) == 0 && len(exp.Should) == 0
}

func KeywordExpression(keyword string) TermQueryV0 {
	return TermQueryV0{Keyword: keyword}
}

func MustExpression(exps ...TermQueryV0) TermQueryV0 {
	if len(exps) == 0 {
		return TermQueryV0{}
	}
	must := make([]TermQueryV0, 0, len(exps))
	for _, exp := range exps {
		if !exp.Empty() {
			must = append(must, exp)
		}
	}
	return TermQueryV0{Must: must}
}

func ShouldExpression(exps ...TermQueryV0) TermQueryV0 {
	if len(exps) == 0 {
		return TermQueryV0{}
	}
	should := make([]TermQueryV0, 0, len(exps))
	for _, exp := range exps {
		if !exp.Empty() {
			should = append(should, exp)
		}
	}
	return TermQueryV0{Should: should}
}

func (exp TermQueryV0) String() string {
	if len(exp.Keyword) > 0 {
		return exp.Keyword
	} else if len(exp.Must) > 0 {
		if len(exp.Must) == 1 {
			return exp.Must[0].String()
		} else {
			strBuilder := strings.Builder{}
			strBuilder.WriteByte('(')
			for _, ele := range exp.Must {
				strBuilder.WriteString(ele.String())
				strBuilder.WriteByte('&')
			}
			str := strBuilder.String()
			str = str[0:len(str)-1] + ")"
			return str
		}
	} else {
		if len(exp.Should) == 1 {
			return exp.Should[0].String()
		} else {
			strBuilder := strings.Builder{}
			strBuilder.WriteByte('(')
			for _, ele := range exp.Should {
				strBuilder.WriteString(ele.String())
				strBuilder.WriteByte('|')
			}
			str := strBuilder.String()
			str = str[0:len(str)-1] + ")"
			return str
		}
	}
}
