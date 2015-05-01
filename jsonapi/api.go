package jsonapi;

type API struct{
    Resources map[string]MountedResource;
    Linkages map[string]map[string]MountedLinkage
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
    return &API{
        Resources: make(map[string]MountedResource),
        Linkages: make(map[string]map[string]MountedLinkage),
    };
}

func(a *API) MountResource(name string, r Resource, p Permissions) {
    a.Resources[name] = MountedResource{R: r, P: p};
}

func(a *API) MountLinkage(name, srcR, dstR string, behavior Behavior) {
    if(a.Linkages[srcR] == nil) {
        a.Linkages[srcR] = make(map[string]MountedLinkage);
    }
    a.Linkages[srcR][name] = MountedLinkage{
        DstR: dstR,
        Behavior: behavior,
    };
}
