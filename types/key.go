package types

type key int

const (
	ProcIdKey key = iota
	TxKey
)

type CodeName string

var (
	ValidationErr CodeName = "VALIDATION_ERR"
)
