package types

import "strings"

// type TermQuery struct {
// 	Must    []*TermQuery
// 	Should  []*TermQuery
// 	Keyword *Keyword
// }

func NewTermQuery(field, keyword string) *TermQuery {
	return &TermQuery{Keyword: &Keyword{Field: field, Word: keyword}} //TermQuery的一级成员里只有Field-keyword非空，Must和Should都为空
}

func (q TermQuery) Empty() bool {
	return q.Keyword == nil && len(q.Must) == 0 && len(q.Should) == 0
}

// Builder生成器模式，方法返回结构体本身，所以可以通过链式调用一直调用
func (q *TermQuery) And(querys ...*TermQuery) *TermQuery {
	if len(querys) == 0 {
		return q
	}
	must := make([]*TermQuery, 0, len(querys)+1)
	if !q.Empty() {
		must = append(must, q)
	}
	for _, query := range querys {
		if !query.Empty() {
			must = append(must, query)
		}
	}
	// TermQuery的一级成员只有Must非空，Should和KeyWord都为空
	return &TermQuery{Must: must}
}

func (q *TermQuery) Or(querys ...*TermQuery) *TermQuery {
	if len(querys) == 0 {
		return q
	}
	should := make([]*TermQuery, 0, len(querys)+1)
	if !q.Empty() {
		should = append(should, q)
	}
	for _, query := range querys {
		if !query.Empty() {
			should = append(should, query)
		}
	}
	return &TermQuery{Should: should}
}

// print函数会自动调用变量的ToString()方法
func (q TermQuery) ToString() string {
	if q.Keyword != nil {
		return q.Keyword.ToString()
	} else if len(q.Must) > 0 {
		if len(q.Must) == 1 {
			return q.Must[0].ToString()
		} else {
			sb := strings.Builder{}
			sb.WriteByte('(')
			for _, e := range q.Must {
				s := e.ToString()
				if len(s) > 0 {
					sb.WriteString(s)
					sb.WriteByte('&')
				}
			}
			s := sb.String()
			s = s[0:len(s)-1] + ")"
			return s
		}
	} else if len(q.Should) > 0 {
		if len(q.Should) == 1 {
			return q.Should[0].ToString()
		} else {
			sb := strings.Builder{}
			sb.WriteByte('(')
			for _, e := range q.Should {
				s := e.ToString()
				if len(s) > 0 {
					sb.WriteString(s)
					sb.WriteByte('|')
				}
			}
			s := sb.String()
			s = s[0:len(s)-1] + ")"
			return s
		}

	}
	return ""
}
