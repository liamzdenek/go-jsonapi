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
    Resource *APIMountedResource
    Relationship *APIMountedRelationship
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
    ef.Resource = amr;
    rels := ef.Request.API.GetRelationshipsByResource(amr.Name);
    for _, rel := range rels {
        should_fetch := ii.ShouldFetch(rel.Name);
        should_include := ii.ShouldInclude(rel.Name);
        if should_fetch || should_include {
            target_resource := ef.Request.API.GetResource(rel.DstResourceName)
            target_future := target_resource.GetFuture();
            tef := NewExecutableFuture(ef.Request, target_future);
            tef.Relationship = rel;
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
        fmt.Printf("Got Response: %#v\n", res)
        if res.IsSuccess {
            for _, field := range res.Success {
                if hasr, ok := field.(FutureResponseKindWithRecords); ok {
                    records := hasr.GetRecords()
                    for _, record := range records {
                        if record.Type == "" {
                            record.Type = ef.Resource.Name;
                        }
                    }
                    for _, rel := range ef.Request.API.GetRelationshipsByResource(ef.Resource.Name) {
                        for _, record := range records {
                            newrel := &ORelationship{
                                IsSingle: rel.Relationship.IsSingle(),
                                RelationshipName: rel.Name,
                                RelatedBase: ef.Request.GetBaseURL()+record.Type+"/"+record.Id,
                            }
                            record.PushRelationship(newrel);
                        }
                    }
                }
            };
        }
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
                if !tefres.IsSuccess {
                    return;
                }

                rel.PushBackRelationships(ef.Request, ef, tef, relres, tefres.Success[tef.Future]);
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
