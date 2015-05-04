package jsonapi;

import ("encoding/json";);

type OutputData struct { // data
    Data []*OutputDatum
    LinkageSet *OutputLinkageSet
    Linkage *OutputLinkage
    Target OutputDataType
}

type OutputDataType int;

const (
    SingleResource OutputDataType = iota
    ManyResources
    Relationship
    OneToOneLinkage
    OneToManyLinkage
);

func NewOutputDataResources(isSingle bool, data []*OutputDatum) *OutputData {
    t := ManyResources;
    if(isSingle) {
        t = SingleResource;
    }
    return &OutputData{Data: data, Target: t};
}

func NewOutputDataRelationship(links *OutputLinkageSet) *OutputData {
    t := Relationship;
    return &OutputData{LinkageSet:links, Target: t};
}

func NewOutputDataLinkage(isSingle bool, l *OutputLinkage) *OutputData {
    t := OneToManyLinkage;
    if(isSingle) {
        t = OneToOneLinkage;
    }
    return &OutputData{Linkage:l, Target: t}
}

func (o OutputData) MarshalJSON() ([]byte, error) {
    //Primary data MUST be either:
    //* a single resource object or null, for requests that target single resources
    //* an array of resource objects or an empty array ([]), for requests that target resource collections
    //* resource linkage, for requests that target a resource's relationship
    // A logical collection of resources (e.g., the target of a to-many relationship) MUST be represented as an array, even if it only contains one item.
    if(o.Target == ManyResources) {
        return json.Marshal(o.Data);
    }
    if(o.Target == SingleResource) {
        if(len(o.Data) == 0) {
            return json.Marshal(nil);
        }
        return json.Marshal(o.Data[0]);
    }
    if(o.Target == Relationship) {
        return json.Marshal(o.LinkageSet);
    }
    if(o.Target == OneToOneLinkage) {
        if(len(o.Linkage.Links) == 0) {
            return json.Marshal(nil);
        }
        return json.Marshal(o.Linkage.Links[0]);
    }
    if(o.Target == OneToManyLinkage) {
        return json.Marshal(o.Linkage.Links);
    }
    panic("Unknown data type sent to OutputData");
}

type OutputDatum struct { // data[i]
    Datum IderLinkerTyper
}

func (o OutputDatum) MarshalJSON() ([]byte, error) {
    res := DenatureObject(o.Datum);
    delete(res, "ID");
    delete(res, "Id");
    delete(res, "iD");
    res["id"] = o.Datum.Id();
    res["links"] = o.Datum.Link();
    res["type"] = o.Datum.Type();
    return json.Marshal(res);
}

