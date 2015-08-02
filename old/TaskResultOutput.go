package jsonapi

type TaskResultOutput interface {
	Task
	GetResult() *Output
}
