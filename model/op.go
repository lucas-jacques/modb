package model

type Operation uint8

const (
	OpInsert Operation = iota
	OpUpdate
	OpSelect
	OpDelete
)
