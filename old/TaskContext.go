package jsonapi;

/**
TaskContext's primary responsibility is managing a set of Tasks. It is also responsible for capturing panics and ensuring that every other task is informed of this such that as little computational power is wasted as possible once we are certain of our response.
*/
type TaskContext struct {
    Context chan Task
}

/**
NewTaskContext() creates and initializes a *TaskContext object, and initializes its worker goroutine
*/
func NewTaskContext(r *Request) *TaskContext {
    res := &TaskContext{
        Context: make(chan Task, 10), // the buffer is to prevent some context switching when a lot of tasks are pushed at once
    };
    go func() {
        has_paniced := false;
        for tworker := range res.Context {
            worker := tworker; // range will reuse the same worker object since it is not a pointer... we do not want it to overwrite the last one before the go func() has a chance to start -- removing this could create inconsistent behavior
            defer worker.Cleanup(r);
            go func() {
                defer func() {
                    if raw := recover(); !has_paniced && raw != nil {
                        has_paniced = true;
                        r.HandlePanic(raw);
                    }
                    worker.ResponseWorker(has_paniced);
                }();
                worker.Work(r);
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
