package jsonapi;

import ("encoding/json";"errors";"fmt");

type OutputData struct { // data
    Data []*OutputDatum
    LinkageSet *OutputLinkageSet
    Linkage *OutputLinkage
    Target OutputDataType
    Included *[]Record `json:"-"`
}

type OutputDataType int;

const (
    SingleResource OutputDataType = iota
    ManyResources
    Relationship
    OneToOneLinkage
    OneToManyLinkage
);

func NewOutputDataResources(isSingle bool, data []Record) *OutputData {
    wrapped := []*OutputDatum{};
    for _, datum := range data {
        wrapped = append(wrapped, &OutputDatum{
            Datum: datum,
        });
    }
    t := ManyResources;
    if(isSingle) {
        t = SingleResource;
    }
    return &OutputData{Data: wrapped, Target: t};
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

func(o *OutputData) Prepare() {
    for i,datum := range o.Data {
        if(datum.Datum == nil) {
            fmt.Printf("DAtum is null\n");
            o.Data = append(o.Data[:i],o.Data[i+1:]...);
            continue;
        }
        datum.Prepare();
    }
}

func (o OutputData) MarshalJSON() ([]byte, error) {
    //Primary data MUST be either:
    //* a single resource object or null, for requests that target single resources
    //* an array of resource objects or an empty array ([]), for requests that target resource collections
    //* resource linkage, for requests that target a resource's relationship
    // A logical collection of resources (e.g., the target of a to-many relationship) MUST be represented as an array, even if it only contains one item.
    o.Included = &[]Record{};
    if(o.Target == ManyResources || o.Target == SingleResource) {
        if(o.Target == ManyResources) {
            return json.Marshal(o.Data);
        }
        if(o.Target == SingleResource) {
            if(len(o.Data) == 0) {
                return json.Marshal(nil);
            }
            return json.Marshal(o.Data[0]);
        }
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
    panic(NewResponderBaseErrors(500, errors.New("Unknown data type sent to OutputData")));
}

type OutputDatum struct { // data[i]
    Datum Record
    res map[string]interface{}
}

func (o *OutputDatum) Prepare() {
    //a.Logger.Printf("Denatre object: %#v\n", o.Datum);
    res := DenatureObject(o.Datum);
    delete(res, "ID");
    delete(res, "Id");
    delete(res, "iD");
    res = map[string]interface{}{"attributes":res};
    res["id"] = GetId(o.Datum);
    links := o.Datum.Data().Links;
    if(len(links.Linkages) > 0) {
        res["relationships"] = links
    }
    res["type"] = o.Datum.Type();
    o.res = res;
}

func (o OutputDatum) MarshalJSON() ([]byte, error) {
    return json.Marshal(o.res);
}

