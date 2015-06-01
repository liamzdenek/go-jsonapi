package jsonapi;


func init() {
    var t Resource = &ResourceTypeSetter{};
    _ = t;
}

type ResourceTypeSetter struct {
    Parent Resource
    Name string
    SetEmptyOnly bool
}

func NewResourceTypeSetter(parent Resource, name string, setEmptyOnly bool) *ResourceTypeSetter {
    return &ResourceTypeSetter{
        Parent: parent,
        Name: name,
        SetEmptyOnly: setEmptyOnly,
    }
}

func(rts *ResourceTypeSetter) Set(records *[]*Record) {
    for _, record := range *records {
        if !rts.SetEmptyOnly || record.Type == "" {
            record.Type = rts.Name;
        }
    }
}

func(rts *ResourceTypeSetter) FindDefault(r *Request, rp RequestParams) ([]*Record, error) {
    records, err := rts.Parent.FindDefault(r,rp);
    rts.Set(&records);
    return records, err;
}

func(rts *ResourceTypeSetter) FindOne(r *Request, rp RequestParams, id string) (*Record, error) {
    record, err := rts.Parent.FindOne(r,rp,id);
    rts.Set(&[]*Record{record});
    return record, err;
}

func(rts *ResourceTypeSetter) FindMany(r *Request, rp RequestParams, ids []string) ([]*Record, error) {
    records, err := rts.Parent.FindMany(r,rp,ids);
    rts.Set(&records);
    return records, err;
}
func(rts *ResourceTypeSetter) FindManyByField(r *Request, rp RequestParams, field, value string) ([]*Record, error) {
    records, err := rts.Parent.FindManyByField(r,rp,field,value);
    rts.Set(&records);
    return records, err;
}
