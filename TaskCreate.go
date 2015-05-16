package jsonapi;

import("net/http";);

type TaskCreate struct {

}

func NewTaskCreate(resource, id string) *TaskCreate {
    return &TaskCreate{

    }
}

func(t *TaskCreate) Work(a *API, r *http.Request) {
    
}

func(t *TaskCreate) ResponseWorker(has_paniced bool) {

}

func(t *TaskCreate) Cleanup(a *API, r *http.Request) {

}
