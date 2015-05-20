package jsonapi;

import ("net/http";"runtime";"fmt");

type TaskContext struct {
    Context chan Task
}

func NewTaskContext(a *API, r *http.Request, w http.ResponseWriter, session Session) *TaskContext {
    res := &TaskContext{
        Context: make(chan Task, 10), // the buffer is to prevent some context switching when a lot of tasks are pushed at once
    };
    go func() {
        has_paniced := false;
        for worker := range res.Context {
            tworker := worker; // range will reuse the same worker object since it is not a pointer... we do not want it to overwrite the last one before the go func() has a chance to start -- removing this could create inconsistent behavior
            defer tworker.Cleanup(a,r);
            go func() {
                defer func() {
                    if raw := recover(); !has_paniced && raw != nil {
                        has_paniced = true;
                        _, should_print_stack := a.CatchResponses(w,r,raw);
                        if(should_print_stack) {
                            const size = 64 << 10
                            buf := make([]byte, size)
                            buf = buf[:runtime.Stack(buf, false)]
                            fmt.Printf("jsonapi: panic %v\n%s", raw, buf);
                        }
                    }
                    tworker.ResponseWorker(has_paniced);
                }();
                tworker.Work(a,session,res,r);
            }();
        }
    }()
    return res;
}

func(w *TaskContext) Cleanup() {
    close(w.Context);
}

func(w *TaskContext) Push(t_list ...Task) {
    for _,t := range t_list {
        w.Context <- t;
    }
}
