package jsonapi;

import "sync";
import "fmt"

type FutureList struct{
    Unabridged []*PreparedFuture
    Optimized []*PreparedFuture
    Requests map[*PreparedFuture][]*FutureRequest
    IsOptimized bool
}

type PreparedFuture struct {
    Parents []*PreparedFuture
    Future Future
    Input chan *FutureRequest
    Relationship Relationship
    CombinedWith []*PreparedFuture
    //IsIncluded, IsPrimaryData bool
}

func NewFutureList() *FutureList {
    return &FutureList{
        Unabridged: []*PreparedFuture{},
        Optimized: []*PreparedFuture{},
        Requests: map[*PreparedFuture][]*FutureRequest{},
    }
}

func (fl *FutureList) PushFuture(pf *PreparedFuture) {
    fl.Unabridged = append(fl.Unabridged, pf);
}

func (fl *FutureList) PushRequest(pf *PreparedFuture, req *FutureRequest) {
    if _, ok := fl.Requests[pf]; !ok {
        fl.Requests[pf] = []*FutureRequest{};
    }
    fl.Requests[pf] = append(fl.Requests[pf], req);
}

func(fl *FutureList) Build(r *Request, amr *APIMountedResource, should_include_root bool) *FutureOutput{
    fo := &FutureOutput{};
    output := &PreparedFuture{
        Future: fo,
        Relationship: &RelationshipIdentity{},
    }
    for node, _ := range fl.Requests {
        fl.recurseBuild(r, node, amr, r.IncludeInstructions, output, should_include_root);
    }
    fl.PushFuture(output);
    return fo;
}

func(fl *FutureList) recurseBuild(r *Request, node *PreparedFuture, amr *APIMountedResource, ii *IncludeInstructions, output *PreparedFuture, should_include_root bool) {
    rels := r.API.GetRelationshipsByResource(amr.Name);
    if should_include_root {
        output.Parents = append(output.Parents, node);
    }
    for _, rel := range rels {
        should_fetch := ii.ShouldFetch(rel.Name);
        should_include := ii.ShouldInclude(rel.Name);
        if should_fetch || should_include {
            target_future := rel.Relationship.GetTargetFuture();
            prepared := &PreparedFuture{
                Parents: []*PreparedFuture{node},
                Future: target_future,
                Relationship: rel,
            }
            if should_include {
                panic("SHOULD INCLUDE");
            }
            fl.recurseBuild(r, prepared, r.API.GetResource(rel.DstResourceName), ii.GetChild(rel.Name), output, should_include);
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

func(fl *FutureList) Run(r *Request) {
    for _,prepfuture := range fl.Optimized {
        prepfuture.Input = make(chan *FutureRequest);
        go func(prepfuture *PreparedFuture) {
            defer r.CatchPanic();
            fmt.Printf("RUNNING: %#v\n", prepfuture);
            fmt.Printf("RUNNING: %T\n", prepfuture.Future);
            prepfuture.Future.Work(prepfuture);
        }(prepfuture);
    }
}

func(fl *FutureList) Defer() {
    for _,prepfuture := range fl.Optimized {
        close(prepfuture.Input);
    }
}

func(fl *FutureList) Takeover(r *Request) {
    fl.Optimize();
    fl.Run(r);
    wg := &sync.WaitGroup{};
    for pf,requests := range fl.Requests {
        for _,request := range requests {
            fl.HandleInput(r, wg, pf, request);
        }
    }
    wg.Wait();
}

func(fl *FutureList) HandleInput(r *Request, wg *sync.WaitGroup, pf *PreparedFuture, req *FutureRequest) {
    wg.Add(1);
    go func() {
        defer r.CatchPanic();
        defer wg.Done();
        fmt.Printf("Sending Request: %#v\n", req);
        pf.Input <- req;
        fmt.Printf("Getting response...\n");
        res := <-req.Response;
        fmt.Printf("Got Response: %#v\n", res);
        if !res.IsSuccess {
            // TODO()
            panic(res.Failure);
        }
        OUTER:for _,prepfuture := range fl.Optimized {
            for future,_ := range res.Success {
                for _, parent := range prepfuture.Parents {
                    if parent.Future == future {
                        reqs := prepfuture.Relationship.Link(r, res);
                        for _,rawreq := range reqs {
                            req := &FutureRequest{
                                Request: r,
                                Response: make(chan *FutureResponse),
                                Kind: rawreq,
                            };
                            fl.HandleInput(r,wg,prepfuture, req);
                        }
                        continue OUTER;
                    }
                }
            }
        }
    }()
}
