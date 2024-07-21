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

type Filter string

var (
	Equal    Filter = "="
	Required Filter = "required"
	File     Filter = "file"
)
