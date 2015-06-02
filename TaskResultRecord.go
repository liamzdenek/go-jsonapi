package jsonapi;

type TaskResultRecords interface {
    Task
    GetResult() *TaskResultRecordData
}
type TaskResultRecordData struct {
    Records []*Record
    Paginator *Paginator
    IsSingle bool
}

