package jsonapi;

type TaskReplyer struct {
    TaskResultOutput TaskResultOutput
    Output chan chan bool
}

func NewTaskReplyer(wo TaskResultOutput) *TaskReplyer {
    return &TaskReplyer{
        TaskResultOutput: wo,
        Output: make(chan chan bool),
    }
}

func(w *TaskReplyer) Work(r *Request) {
    //a.Logger.Printf("Waiting for final result\n");
    res := w.TaskResultOutput.GetResult();
    //a.Logger.Printf("FINAL RESULT: %#v\n", res);
    Reply(res)
}

func(w *TaskReplyer) ResponseWorker(has_paniced bool) {
    go func() {
        for req := range w.Output {
            req <- true;
        }
    }();
}


func(w *TaskReplyer) Cleanup(r *Request) {
    defer close(w.Output);
}

func (w *TaskReplyer) Wait() {
    r := make(chan bool);
    defer close(r);
    w.Output <- r;
    <-r;
}
