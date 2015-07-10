package jsonapi;

import "fmt"
import "sync"

// ExecutableFuture.Children and ExecutableFuture.Relationships must be disjoint sets
type ExecutableFuture struct {
    Future
    Request *Request
    Children map[Relationship]*ExecutableFuture
    Input chan *FutureRequest

    //Relationships map[Relationship]*ExecutableFuture
    ResponsibleFor []*ExecutableFuture
}

func NewExecutableFuture(r *Request, f Future) *ExecutableFuture {
    return &ExecutableFuture{
        Request: r,
        Future: f,
        Children: make(map[Relationship]*ExecutableFuture),
        //Relationships: make(map[Relationship]*ExecutableFuture),
    }
}

func(ef *ExecutableFuture) GetRequest() *FutureRequest {
    var req *FutureRequest;
    var should_break bool;
    select {
    case <-ef.Request.Done.Wait():
        panic(&FutureRequestedPanic{});
    case req, should_break = <-ef.Input:
    }
    if should_break && req == nil {
        panic(&FutureRequestedPanic{});
    }
    return req;
}

func(ef *ExecutableFuture) PushChild(r Relationship, f *ExecutableFuture) {
    ef.Children[r] = f;
}

func(ef *ExecutableFuture) Build(amr *APIMountedResource) *FutureOutput {
    fo := &FutureOutput{};
    efo := &ExecutableFuture{
        Request: ef.Request,
        Future: fo,
    };
    ii := ef.Request.IncludeInstructions;
    ef.internalBuild(amr,ii,efo,true,false);
    return fo;
}

func (ef *ExecutableFuture) internalBuild(amr *APIMountedResource, ii *IncludeInstructions, efo *ExecutableFuture, is_primary, is_included bool) {
    fo := efo.Future.(*FutureOutput);
    if is_primary || is_included {
        rel := &RelationshipIdentity{}
        if is_primary {
            rel.IsPrimary = true;
        }
        ef.PushChild(rel, efo);
        fo.PushParent(ef);
    }
    rels := ef.Request.API.GetRelationshipsByResource(amr.Name);
    for _, rel := range rels {
        should_fetch := ii.ShouldFetch(rel.Name);
        should_include := ii.ShouldInclude(rel.Name);
        if should_fetch || should_include {
            target_resource := ef.Request.API.GetResource(rel.DstResourceName)
            target_future := target_resource.GetFuture();
            tef := NewExecutableFuture(ef.Request, target_future);
            tef.internalBuild(target_resource, ii.GetChild(rel.Name), efo, false, should_include);
            ef.PushChild(rel, tef);
        }
    }
}

func(ef *ExecutableFuture) Defer() {
    if ef.Input != nil {
        close(ef.Input);
    }
    for _,tef := range ef.Children {
        tef.Defer();
    }
}

func(ef *ExecutableFuture) Optimize() {
    //panic("TODO");
}

func(ef *ExecutableFuture) Execute() {
    if ef.Input == nil {
        ef.Request.API.Logger.Debugf("Running: %#v\n", ef.Future);
        ef.Input = make(chan *FutureRequest);
        ef.Go(func() {
            ef.Future.Work(ef);
        });
        for _,tef := range ef.Children {
            tef.Execute();
        }
    }
}

func(ef *ExecutableFuture) Takeover(fr *FutureRequest) {
    ef.Optimize();
    ef.Execute();
    ef.HandleRequest(fr, nil);
    <-ef.Request.Done.Wait()
}

func(ef *ExecutableFuture) HandleRequest(req *FutureRequest, cb func(*FutureResponse)) {
    ef.Go(func() {
        fmt.Printf("Sending Request: %#v %#v\n", ef.Input, req);
        ef.Input <- req;
        fmt.Printf("Getting response...\n");
        res := req.GetResponse();
        fmt.Printf("Got Response: %#v\n", res);
        if cb != nil {
            cb(res);
        }
        ef.HandleResponse(res);
    });
}

func(ef *ExecutableFuture) HandleResponse(res *FutureResponse) {
    wg := &sync.WaitGroup{};
    res.WaitForComplete = make(chan bool);
    defer func() {
        wg.Wait()
        close(res.WaitForComplete);
    }();
    if !res.IsSuccess {
        panic(res.Failure);
    }
    fmt.Printf("HANDLERESPONSE TAKING OVER\n");
    efs := append(ef.ResponsibleFor, ef);
    // for each future that this is responsible for...
    for _, cef := range efs {
        wg.Add(len(cef.Children));
        // iterate through each of the children of that future...
        for rel, tef := range cef.Children {
            if cef == tef { // simple inf loop check
                wg.Done();
                continue;
            }
            // take the output data for that future
            relres := res.Success[cef.Future];
            // convert it into the request for that child
            reqkind := rel.Link(ef.Request, cef, tef, relres)
            req := &FutureRequest{
                Request: ef.Request,
                Response: make(chan *FutureResponse),
                Kind: reqkind,
            };
            fmt.Printf("SENDING HANDLEREQUEST TO TEF\n");
            tef.HandleRequest(req, func(tefres *FutureResponse) {
                defer wg.Done();
                fmt.Printf("GOT RESPONSE FROM TEF: %#v\n", tefres);
                if tefres.IsSuccess {
                    if modifier, ok := tefres.Success[tef.Future].(FutureResponseModifier); ok {
                        modifier.Modify(relres);
                    }
                }
            });
        }
    }
}

func(ef *ExecutableFuture) Go(f func()) {
    go func() {
        defer ef.CatchPanic();
        f();
    }();
}

func(ef *ExecutableFuture) CatchPanic() {
    if raw := recover(); raw != nil {
        if _, ok := raw.(*FutureRequestedPanic); !ok {
            ef.Request.HandlePanic(raw);
        }
    }
}

/*
type FutureList struct{
    Request *Request
    Unabridged []*PreparedFuture
    Optimized []*PreparedFuture
    Requests map[*PreparedFuture][]*FutureRequest
    IsOptimized bool
}

type PreparedFuture struct {
    FutureList *FutureList
    Parents []*PreparedFuture
    Future Future
    Input chan *FutureRequest
    Relationship Relationship
    CombinedWith []*PreparedFuture
    Current *FutureRequest
    //IsIncluded, IsPrimaryData bool
}

func(pf *PreparedFuture) GetNext() (req *FutureRequest, should_break bool) {
    select {
    case <-pf.FutureList.Request.Done.Wait():
        panic(&FutureRequestedPanic{});
    case req, should_break = <-pf.Input:
    }
    should_break = !should_break;
    pf.Current = req;
    return
}

func NewFutureList(r *Request) *FutureList {
    return &FutureList{
        Request: r,
        Unabridged: []*PreparedFuture{},
        Optimized: []*PreparedFuture{},
        Requests: map[*PreparedFuture][]*FutureRequest{},
    }
}

func (fl *FutureList) PushFuture(pf *PreparedFuture) {
    pf.FutureList = fl;
    fl.Unabridged = append(fl.Unabridged, pf);
}

func (fl *FutureList) PushRequest(pf *PreparedFuture, req *FutureRequest) {
    if _, ok := fl.Requests[pf]; !ok {
        fl.Requests[pf] = []*FutureRequest{};
    }
    fl.Requests[pf] = append(fl.Requests[pf], req);
}

func(fl *FutureList) Build(node *PreparedFuture, amr *APIMountedResource, should_include_root bool) *FutureOutput{
    fo := &FutureOutput{};
    output := &PreparedFuture{
        Future: fo,
        Relationship: &RelationshipIdentity{},
    }
    fl.recurseBuild(node, amr, fl.Request.IncludeInstructions, output, should_include_root);
    fl.PushFuture(output);
    return fo;
}

func(fl *FutureList) recurseBuild(node *PreparedFuture, amr *APIMountedResource, ii *IncludeInstructions, output *PreparedFuture, should_include_root bool) {
    rels := fl.Request.API.GetRelationshipsByResource(amr.Name);
    if should_include_root {
        output.Parents = append(output.Parents, node);
    }
    for _, rel := range rels {
        should_fetch := ii.ShouldFetch(rel.Name);
        should_include := ii.ShouldInclude(rel.Name);
        if should_fetch || should_include {
            target_future := fl.Request.API.GetResource(rel.DstResourceName).GetFuture();
            prepared := &PreparedFuture{
                Parents: []*PreparedFuture{node},
                Future: target_future,
                Relationship: rel,
            }
            if should_include {
                panic("SHOULD INCLUDE "+amr.Name);
            }
            fl.recurseBuild(prepared, fl.Request.API.GetResource(rel.DstResourceName), ii.GetChild(rel.Name), output, should_include);
        }
    }
}

func(fl *FutureList) Optimize() {
    if fl.IsOptimized {
        return;
    }
    fl.IsOptimized = true;
    //fmt.Printf("OPTIMIZE: %#v\n", fl);
    fl.Optimized = fl.Unabridged;
}

func(fl *FutureList) Run() {
    for _,tprepfuture := range fl.Optimized {
        prepfuture := tprepfuture;
        prepfuture.Input = make(chan *FutureRequest);
        fl.Go(func() {
            prepfuture.Future.Work(prepfuture);
        });
    }
}

func(fl *FutureList) Defer() {
    for _,prepfuture := range fl.Optimized {
        close(prepfuture.Input);
    }
}

func(fl *FutureList) Takeover() {
    fl.Optimize();
    fl.Run();
    for pf,requests := range fl.Requests {
        for _,request := range requests {
            fl.HandleInput(pf, request);
        }
    }
    <-fl.Request.Done.Wait()
}

func(fl *FutureList) CatchPanic() {
    if raw := recover(); raw != nil {
        if _, ok := raw.(*FutureRequestedPanic); !ok {
            fl.Request.HandlePanic(raw);
        }
    }
}


func(fl *FutureList) HandleInput(pf *PreparedFuture, req *FutureRequest) {
    fl.Go(func() {
        fmt.Printf("Sending Request: %#v\n", req);
        pf.Input <- req;
        fmt.Printf("Getting response...\n");
        res := req.GetResponse();
        fmt.Printf("Got Response: %#v\n", res);
        if !res.IsSuccess {
            panic(res.Failure);
        }
        OUTER:for _,prepfuture := range fl.Optimized {
            for future,future_output := range res.Success {
                for _, parent := range prepfuture.Parents {
                    if parent.Future == future {
                        rawreq := prepfuture.Relationship.Link(fl.Request, pf, prepfuture, future_output);
                        req := &FutureRequest{
                            Request: fl.Request,
                            Response: make(chan *FutureResponse),
                            Kind: rawreq,
                        };
                        fl.HandleInput(prepfuture, req);
                        continue OUTER;
                    }
                }
            }
        }
    });
}
*/
