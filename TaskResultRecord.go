package jsonapi;

type TaskResultRecords interface {
    Task
    GetResult() *TaskResultRecordData
}
type TaskResultRecordData struct {
    Result []Record
    IsSingle bool
}
