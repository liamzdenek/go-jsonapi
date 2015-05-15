package jsonapi;

import ("net/http";"fmt");

type WorkReplyer struct {
    WorkerOutput WorkerOutput
    Output chan chan bool
}

func NewWorkReplyer(wo WorkerOutput) *WorkReplyer {
    return &WorkReplyer{
        WorkerOutput: wo,
        Output: make(chan chan bool),
    }
}

func(w *WorkReplyer) Work(a *API, r *http.Request) {
    fmt.Printf("Waiting for final result\n");
    res := w.WorkerOutput.GetResult();
    fmt.Printf("FINAL RESULT: %#v\n", res);
    go func() {
        for req := range w.Output {
            req <- true;
        }
    }();
    Reply(res)
}
func(w *WorkReplyer) Cleanup(a *API, r *http.Request) {
    defer close(w.Output);
}

func (w *WorkReplyer) Wait() {
    r := make(chan bool);
    defer close(r);
    w.Output <- r;
    <-r;
}
