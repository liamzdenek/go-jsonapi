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
    BaseURIPath string
}

type MountedResource struct{
    R Resource
    P Permissions
}

type MountedLinkage struct{
    SrcR string
    DstR string
    LinkageBehavior LinkageBehavior
}

func(mr *MountedLinkage) Resolve(a *API, src HasId, r *http.Request) (res Linkage, included []interface{}) {
    resource := a.Resources[mr.DstR];
    resource.P.Check(mr.SrcR+".linkto."+mr.DstR+".FindMany", "", r);
    switch lb := mr.LinkageBehavior.(type) {
        case IdLinkageBehavior:
            ids := lb.Link(src);
            for _, id := range ids {
                res.Linkage = append(res.Linkage, LinkageIdentifier{
                    Type: mr.DstR,
                    Id: id,
                });
            }
            linkdata, err := resource.R.FindMany(ids, r);
            Check(err);
            for _, link := range linkdata {
                fixedlink,_ := a.AddLinkages(link, mr.DstR, r, false);
                included = append(included, fixedlink);
            }
        case HasIdLinkageBehavior:
            panic("TODO");
        default:
            panic("Attempted to resolve a linkage behavior that is neither an Id or HasId LinkageBehavior.. This should never happen");
    }
    return;
}

func NewAPI() *API {
    api := &API{
        Router: pat.New(),
        Resources: make(map[string]*MountedResource),
        Linkages: make(map[string]map[string]*MountedLinkage),
        BaseURIPath: "/",
    };
    api.InitRouter();
    return api;
}

func(a *API) MountResource(name string, r Resource, p Permissions) {
    a.Resources[name] = &MountedResource{R: r, P: p};
}

func(a *API) MountLinkage(name, srcR, dstR string, behavior LinkageBehavior) {
    if(a.Resources[srcR] == nil) {
        panic("Source resource "+srcR+" for linkage does not exist");
    }
    if(a.Resources[dstR] == nil) {
        panic("Destination resource "+dstR+" for linkage does not exist");
    }
    if(a.Linkages[srcR] == nil) {
        a.Linkages[srcR] = make(map[string]*MountedLinkage);
    }
    if(!VerifyLinkageBehavior(behavior)) {
        panic("Linkage provided cannot be used as an Id or HasId LinkageBehavior");
    }
    a.Linkages[srcR][name] = &MountedLinkage{
        SrcR: srcR,
        DstR: dstR,
        LinkageBehavior: behavior,
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

    resource.P.Check(resource_str+".FindAll", id_str, r);

    data, err := resource.R.FindOne(id_str, r);
    Check(err);

    Reply(a.PrepareResponse(data, resource_str,r))
}

func(a *API) PrepareResponse(data HasId, resource_str string, r *http.Request) interface{} {
    res := map[string]interface{}{}
    fmt.Printf("Adding linkages %#v\n", data);
    if r, included := a.AddLinkages(data,resource_str,r,true); r != nil {
        res["data"] = r;
        if included != nil {
            res["included"] = included;
        }
    } else {
        res["data"] = nil;
    }
    return res;
}

func(a *API) AddLinkages(data HasId, resource_str string, r *http.Request, recursive bool) (interface{}, interface{}) {
    if(data == nil) {
        return nil, nil;
    }
    res := DenatureObject(data);
    var included interface{};
    var links interface{};

    delete(res, "Id");
    delete(res, "ID");
    delete(res, "iD");
    res["id"] = data.GetId();
    res["type"] = resource_str;
    if(recursive) {
        links, included = a.GenerateLinkages(data, resource_str, r, true);
        res["links"] = links;
    }
    fmt.Printf("Res: %#v\n", res);

    return res, included;
}

func(a *API) GenerateLinkages(data HasId, resource_str string, r *http.Request, getIncluded bool) (interface{}, interface{}) {
    res := map[string]interface{}{};
    included := []interface{}{};
    if linkages := a.Linkages[resource_str]; len(linkages) > 0 {
        for linkname, linkage := range linkages {
            linkdata, incl := linkage.Resolve(a, data, r)
            linkdata.Self = a.GetBaseURL(r)+resource_str+"/"+data.GetId()+"/links/"+linkname;
            linkdata.Related = a.GetBaseURL(r)+resource_str+"/"+data.GetId()+"/"+linkname;
            res[linkname] = linkdata
            for _, hasid := range incl {
                included = append(included, hasid);
            }
        }
    }
    res["self"] = a.GetBaseURL(r)+resource_str+"/"+data.GetId();
    if len(included) == 0 {
        return res, nil;
    }
    return res, included;
}

func(a *API) GetBaseURL(r *http.Request) string {
    fmt.Printf("URL: %#v\n", r.URL);
    if r.URL.Scheme == "" {
        r.URL.Scheme = "http";
    }
    return r.URL.Scheme+"://"+r.Host+a.BaseURIPath;
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
