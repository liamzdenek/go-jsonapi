package jsonapi;

import (
    "github.com/gorilla/pat"
    "net/http"
    "encoding/json"
    "fmt"
    "reflect"
    "strings"
);

type API struct{
    Resources map[string]*MountedResource;
    Linkages map[string]map[string]*MountedLinkage
    Router *pat.Router
}

type MountedResource struct{
    R Resource
    P Permissions
}

type MountedLinkage struct{
    DstR string
    Behavior Behavior
}

func NewAPI() *API {
    api := &API{
        Router: pat.New(),
        Resources: make(map[string]*MountedResource),
        Linkages: make(map[string]map[string]*MountedLinkage),
    };
    api.InitRouter();
    return api;
}

func(a *API) MountResource(name string, r Resource, p Permissions) {
    a.Resources[name] = &MountedResource{R: r, P: p};
}

func(a *API) MountLinkage(name, srcR, dstR string, behavior Behavior) {
    if(a.Linkages[srcR] == nil) {
        a.Linkages[srcR] = make(map[string]*MountedLinkage);
    }
    a.Linkages[srcR][name] = &MountedLinkage{
        DstR: dstR,
        Behavior: behavior,
    };
}

func(a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    defer a.CatchResponses(w,r);
    a.Router.ServeHTTP(w, r);
}

func(a *API) CatchResponses(w http.ResponseWriter, r *http.Request) {
    if raw := recover(); raw != nil {
        switch r := raw.(type) {
        case error:
            res := &TopLevel{};
            res.Errors = []error{r};
            a.Send(res, w);
            panic(r);
        case *ResponseReply:
            a.Send(r.Reply, w);
        default:
            w.Write([]byte(fmt.Sprintf("Error handling request: %#v\n", r)));
        }
    }
}

func(a *API) Send(obj interface{}, w http.ResponseWriter) {
    str, err := json.Marshal(obj);
    Check(err);
    w.Write(str);
}

func(a *API) InitRouter() {
    //a.Router.Get("/:resource")
    a.Router.Get("/{resource}/{id}", a.FindOne)
}

func(a *API) FindOne(w http.ResponseWriter, r *http.Request) {
    resource_str := r.URL.Query().Get(":resource");
    id_str := r.URL.Query().Get(":id");
    
    resource := a.Resources[resource_str];

    if(resource == nil) {
        panic(&ErrorResourceDoesNotExist{Resource:resource_str});
    }

    resource.P.Check(resource_str+".FindAll", id_str, w, r);

    data, err := resource.R.FindOne(id_str, r);
    Check(err);

    Reply(a.PrepareResponse(data, resource_str))
}

func(a *API) PrepareResponse(data HasId, resource_str string) *TopLevel {
    res := &TopLevel{};
    fmt.Printf("Adding linkages %#v\n", data);
    res.Data = a.AddLinkages(data, resource_str);
    return res;
}

func(a *API) AddLinkages(data HasId, resource_str string) interface{} {
    // this probably causes performance problems
    res := DenatureObject(data);

    delete(res, "Id");
    delete(res, "ID");
    delete(res, "iD");
    res["id"] = data.GetId();


    fmt.Printf("Res: %#v\n", res);

    return res;
}

func DenatureObject(data interface{}) map[string]interface{} {
    v := reflect.Indirect(reflect.ValueOf(data));
    t := v.Type();

    values := make(map[string]interface{}, t.NumField());

    for i := 0; i < t.NumField(); i++ {
        tag := strings.Split(t.Field(i).Tag.Get("json"), ",");
        if len(tag[0]) == 0 { 
            tag[0] = t.Field(i).Name
        }
        if len(tag) > 1 && len(tag[1]) > 0 {
            if(tag[1] == "omitempty") {
                if(IsZeroOfUnderlyingType(v.Field(i).Interface())) {
                    continue;
                }
            }
        }
        values[tag[0]] = v.Field(i).Interface();
    }

    return values;
}

func IsZeroOfUnderlyingType(x interface{}) bool {
    return x == reflect.Zero(reflect.TypeOf(x)).Interface()
}
