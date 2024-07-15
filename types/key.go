package types

type key int

const (
	ProcIdKey key = iota
	TxKey
	TokenKey
)

type CodeName string

var (
	ValidationErr CodeName = "VALIDATION_ERR"
)
