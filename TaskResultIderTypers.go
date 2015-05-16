package jsonapi;

type TaskResultIderTypers interface{
    Task
    GetResult() *TaskFindByIdsResult
}

type TaskFindByIdsResult struct {
    Result []IderTyper
    IsSingle bool
}
