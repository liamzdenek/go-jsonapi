package jsonapi;

type TaskResultIderTypers interface{
    GetResult() *TaskFindByIdsResult
}

type TaskFindByIdsResult struct {
    Result []IderTyper
    IsSingle bool
}
