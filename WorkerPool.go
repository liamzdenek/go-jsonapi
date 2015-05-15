package jsonapi;

import ("net/http";"fmt");

type WorkerContext struct {
    Context chan Worker
}

func NewWorkerContext(a *API, r *http.Request) *WorkerContext {
    res := &WorkerContext{
        Context: make(chan Worker),
    };
    go func() {
        panics := make(chan interface{})
        defer close(panics);
        has_paniced := false;
        OUTER: for {
            select {
            case worker, ok := <-res.Context:
                if(!ok) {
                    break OUTER; // chan closed
                }
                if(has_paniced) {
                    continue;
                }
                defer worker.Cleanup(a,r);
                go func() {
                    defer func() {
                        if r := recover(); r != nil {
                            fmt.Printf("PANICS: %#v\n", r);
                            panics <- r;
                        }
                    }();
                    worker.Work(a,r);
                }();
            case caught := <-panics:
                has_paniced = true;
                fmt.Printf("CAUGHT PANIC: %#v\n", caught);
                fmt.Printf("CAUGHT PANIC: %s\n", caught);
            }
        }
        fmt.Printf("CONTEXT CLEANUP\n");
    }()
    return res;
}

func(w *WorkerContext) Cleanup() {
    close(w.Context);
}

func PushWork(wctx *WorkerContext, w Worker) {
    wctx.Context <- w;
}
