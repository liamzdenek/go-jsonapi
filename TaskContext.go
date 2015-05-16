package jsonapi;

import ("net/http";);

type TaskContext struct {
    Context chan Task
}

func NewTaskContext(a *API, r *http.Request, w http.ResponseWriter) *TaskContext {
    res := &TaskContext{
        Context: make(chan Task),
    };
    go func() {
        has_paniced := false;
        for worker := range res.Context {
            tworker := worker; // range will reuse the same worker object since it is not a pointer... we do not want it to overwrite the last one before the go func() has a chance to start -- removing this could create inconsistent behavior
            //fmt.Printf("OUTER WORKER %#v\n", tworker);
            defer tworker.Cleanup(a,r);
            go func() {
                defer func() {
                    if raw := recover(); !has_paniced && raw != nil {
                        //fmt.Printf("\nPANIC: %#v\n\n", raw);
                        has_paniced = true;
                        a.CatchResponses(w,r,raw);
                    }
                    tworker.ResponseWorker(has_paniced);
                }();
                //fmt.Printf("INNER WORKER: %#v\n", tworker);
                tworker.Work(a,r);
            }();
        }
        //fmt.Printf("CONTEXT CLEANUP\n");
    }()
    return res;
}

func(w *TaskContext) Cleanup() {
    close(w.Context);
}

func(w *TaskContext) Push(t Task) {
    w.Context <- t;
}
