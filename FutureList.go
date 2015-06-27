package jsonapi;


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
            target_future := rel.Relationship.GetTargetFuture();
            prepared := &PreparedFuture{
                Parents: []*PreparedFuture{node},
                Future: target_future,
                Relationship: rel,
            }
            if should_include {
                panic("SHOULD INCLUDE");
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

func(fl *FutureList) Go(f func()) {
    go func() {
        defer fl.CatchPanic();
        f();
    }();
}

func(fl *FutureList) HandleInput(pf *PreparedFuture, req *FutureRequest) {
    fl.Go(func() {
        //fmt.Printf("Sending Request: %#v\n", req);
        pf.Input <- req;
        //fmt.Printf("Getting response...\n");
        res := req.GetResponse();
        //fmt.Printf("Got Response: %#v\n", res);
        OUTER:for _,prepfuture := range fl.Optimized {
            for future,_ := range res.Success {
                for _, parent := range prepfuture.Parents {
                    if parent.Future == future {
                        rawreq := prepfuture.Relationship.Link(fl.Request, res);
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
