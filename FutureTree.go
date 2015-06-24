package jsonapi;

type FutureTree struct{
    Root Future
    RootResource *APIMountedResource
    List map[*APIMountedRelationship]Future
    //Tree map[Future][]Future
    Data map[Future]FutureTreeData
    RunList []Future
    ReqChans map[Future]chan *FutureRequest
    Optimized bool
}

type FutureTreeData struct {
    IsIncluded, IsPrimaryData bool
}

func NewFutureTree(root Future, root_resource *APIMountedResource) *FutureTree {
    return &FutureTree{
        Root: root,
        RootResource: root_resource,
        Data: map[Future]FutureTreeData{},
        List: map[*APIMountedRelationship]Future{},
        //Tree: map[Future][]Future{},
        RunList: []Future{},
    }
}

func (ft *FutureTree) Optimize() {
    if(ft.Optimized) {
        return;
    }
    ft.Optimized = true;
}

func (ft *FutureTree) Run(r *Request, fr *FutureRequest) {
    ft.Optimize();

    reqchans := map[Future]chan *FutureRequest{};
    for _,future := range ft.RunList {
        reqchans[future] = make(chan *FutureRequest);
    }
    ft.ReqChans = reqchans;

    for _,future := range ft.RunList {
        go func() {
            future.Work(reqchans[future]);
        }()
    }

    reqchans[ft.Root] <- fr;
}

func (ft *FutureTree) Defer() {
    for _,c := range ft.ReqChans {
        close(c);
    }
}

func(ft *FutureTree) BuildIncludeInstructions(r *Request) {
    ft.Push(ft.Root, nil, nil, FutureTreeData{
        IsIncluded: false,
        IsPrimaryData: true,
    });
    ft.buildIncludeInstructions(r,r.IncludeInstructions,nil);
}

func(ft *FutureTree) Push(f Future, parent Future, amr *APIMountedRelationship, data FutureTreeData) {
    ft.Data[f] = data;
    ft.List[amr] = f;
    ft.RunList = append(ft.RunList, f);
   
    /*
    if parent != nil {
        if _,ok := ft.Tree[parent]; !ok {
            ft.Tree[parent] = []Future{};
        }
        ft.Tree[parent] = append(ft.Tree[parent], f);
    }
    */
}

func(ft *FutureTree) DontRun(f Future) {
    for i, future := range ft.RunList {
        if future == f {
            ft.RunList = append(ft.RunList[:i], ft.RunList[i+1:]...);
        }
    }
}

func(ft *FutureTree) buildIncludeInstructions(r *Request, ii *IncludeInstructions, parent_amr *APIMountedRelationship) {
    parent_future := ft.List[parent_amr];
    var parent_resource *APIMountedResource;
    if parent_amr == nil {
        parent_resource = ft.RootResource
    } else {
        parent_resource = r.API.GetResource(parent_amr.SrcResourceName);
    }

    for relname, rel := range r.API.GetRelationshipsByResource(parent_resource.Name) {
        res := r.API.GetResource(rel.SrcResourceName)
        new_future := rel.Link(r, res, rel, parent_future);

        ft.Push(new_future, parent_future, rel, FutureTreeData{
            IsPrimaryData: false,
            IsIncluded: ii.ShouldInclude(relname),
        });

        if ii.ShouldFetch(relname) {
            r.API.Logger.Debugf("SHOULD FETCH: %s\n", relname);
            ft.buildIncludeInstructions(r, ii.GetChild(relname), rel);
        }
    }
}

func(ft *FutureTree) Takeover(r *Request, req *FutureRequest) {
    r.API.Logger.Debugf("RUNNING\n");
    req.Response = make(chan *FutureResponse);
    
    ft.Run(r, req);
    defer ft.Defer();

    r.API.Logger.Debugf("Got response in main: %#v\n", <-req.Response);
}
