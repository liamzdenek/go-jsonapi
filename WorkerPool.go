package jsonapi;

import ("net/http";"fmt");

type WorkerPool struct {
    Queue chan InternalWorkerWrapper
}

type InternalWorkerWrapper struct {
    A *API
    R *http.Request
    Worker Worker
    Panics chan interface{}
}

func(iww *InternalWorkerWrapper) Work() {
    defer func() {
        if r := recover(); r != nil {
            iww.Panics <- r;
        }
    }();
    iww.Worker.Work(iww.A, iww.R);
}

func NewWorkerPool() *WorkerPool {
    p := &WorkerPool{
        Queue: make(chan InternalWorkerWrapper, 50),
    };
    // TODO: make this configurable
    for i := 0; i < 4; i++ {
        StartWorker(p.Queue);
    }
    return p;
}

func(wp *WorkerPool) NewWorkContext(a *API, r *http.Request) chan Worker {
    res := make(chan Worker);
    go func() {
        panics := make(chan interface{})
        defer close(panics);
        has_paniced := false;
        OUTER: for {
            select {
            case worker, ok := <-res:
                if(!ok) {
                    break OUTER; // chan closed
                }
                if(has_paniced) {
                    continue;
                }
                defer worker.Cleanup(a,r);
                wp.Queue <- InternalWorkerWrapper{
                    Worker: worker,
                    Panics: panics,
                    A: a,
                    R: r,
                }
            case caught := <-panics:
                has_paniced = true;
                fmt.Printf("CAUGHT PANIC: %#v\n", caught);
            }
        }
        fmt.Printf("CONTEXT CLEANUP\n");
    }()
    return res;
}

func StartWorker(input chan InternalWorkerWrapper) {
    go func() {
        for work := range input {
            work.Work();
        }
    }();
}
