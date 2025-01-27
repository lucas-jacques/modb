package modb

import "strconv"

type ParamsSet interface {
	Next() string
}

type Numbered struct {
	count int
}

func (n *Numbered) Next() string {
	n.count++
	return "$" + strconv.Itoa(n.count)
}

// "?" placeholder
type QuestionMark struct{}

func (q *QuestionMark) Next() string {
	return "?"
}
