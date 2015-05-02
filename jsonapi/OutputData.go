package jsonapi;

import ("encoding/json";);

type OutputData struct { // data
    Data []*OutputDatum
    Target OutputDataType
}

type OutputDataType int;

const (
    SingleResource OutputDataType = iota
    ManyResources
    ResourceLinkage
);

func (o OutputData) MarshalJSON() ([]byte, error) {
    //Primary data MUST be either:
    //* a single resource object or null, for requests that target single resources
    //* an array of resource objects or an empty array ([]), for requests that target resource collections
    //* resource linkage, for requests that target a resource's relationship
    if(o.Target == SingleResource) {
        if(len(o.Data) == 0) {
            return json.Marshal(nil);
        }
        return json.Marshal(o.Data[0]);
    }
    if(o.Target == ResourceLinkage) {
        panic("TODO");
        //return json.Marshal(
    }
    // ManyResources
    // A logical collection of resources (e.g., the target of a to-many relationship) MUST be represented as an array, even if it only contains one item.
    return json.Marshal(o.Data);
}

type OutputDatum struct { // data[i]
    Datum IderLinkerTyper
}

func (o OutputDatum) MarshalJSON() ([]byte, error) {
    res := DenatureObject(o.Datum);
    res["id"] = o.Datum.Id();
    res["links"] = o.Datum.Link();
    res["type"] = o.Datum.Type();
    return json.Marshal(res);
}

