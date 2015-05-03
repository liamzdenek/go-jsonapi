package jsonapi;

type ResourceManager struct {
    Resources map[string]*ResourceManagerResource;
    Relationship map[string]map[string]*ResourceManagerRelationship
}

func NewResourceManager() *ResourceManager {
    return &ResourceManager{
        Resources: make(map[string]*ResourceManagerResource),
        Relationship: make(map[string]map[string]*ResourceManagerRelationship),
    }
}

func(rm *ResourceManager) MountResource(name string, r Resource, a Authenticator) {
    rm.Resources[name] = &ResourceManagerResource{R: r, A: a};
}

func(rm *ResourceManager) MountRelationship(name, srcR, dstR string, behavior RelationshipBehavior, auth Authenticator) {
    if(rm.Resources[srcR] == nil) {
        panic("Source resource "+srcR+" for linkage does not exist");
    }
    if(rm.Resources[dstR] == nil) {
        panic("Destination resource "+dstR+" for linkage does not exist");
    }
    if(rm.Relationship[srcR] == nil) {
        rm.Relationship[srcR] = make(map[string]*ResourceManagerRelationship);
    }
    if(!VerifyRelationshipBehavior(behavior)) {
        panic("Linkage provided cannot be used as an Id or HasId LinkageBehavior");
    }
    rm.Relationship[srcR][name] = &ResourceManagerRelationship{
        SrcR: srcR,
        DstR: dstR,
        B: behavior,
        A: auth,
    };
}
