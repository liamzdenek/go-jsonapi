package resource

import (
	. ".."
	"encoding/json"
	"fmt"
	"reflect"
)

func init() {
	var t Resource
	t = &RAM{}
	_ = t
}

type RAM struct {
	Data map[string][]byte
	Type reflect.Type
}

func NewRAM(datatype interface{}) *RAM {
	return &RAM{
		Type: reflect.TypeOf(datatype),
		Data: map[string][]byte{},
	}
}

func (r *RAM) Push(key string, value interface{}) {
	if reflect.TypeOf(value) != r.Type {
		panic("Unsuitable type passed to RAM")
	}
	SetId(value, key)
	data, err := json.Marshal(value)
	Check(err)
	r.Data[key] = data
}

func (r *RAM) Get(key string) interface{} {
	data, ok := r.Data[key]
	if !ok {
		return nil
	}
	v := reflect.New(reflect.PtrTo(r.Type)).Interface()
	if json.Unmarshal(data, v) != nil {
		return nil
	}
	return v
}

func (r *RAM) GetFuture() Future {
	return &FutureRAM{
		RAM: r,
	}
}

type FutureRAM struct {
	RAM *RAM
}

// commands to FutureRAM should never be combined as it would not be more efficient
func (fr *FutureRAM) ShouldCombine(f Future) bool { return false }
func (fr *FutureRAM) Combine(f Future) error      { panic(TODO()) }

func (fr *FutureRAM) Work(ef *ExecutableFuture) {
	for {
		req := ef.GetRequest()
		switch k := req.Kind.(type) {
		case *FutureRequestKindFindByIds:
			fr.WorkFindByIds(ef, req, k)
		default:
			panic(fmt.Sprintf("FutureRAM got unsupported query kind %T: %#v\n", req.Kind, req.Kind))
		}
	}
}

func (fr *FutureRAM) WorkFindByIds(ef *ExecutableFuture, req *FutureRequest, k *FutureRequestKindFindByIds) {
	list := map[Field][]*Record{}
	for _, id := range k.Ids {
		data := fr.RAM.Get(id)
		if data == nil {
			continue
		}
		record := &Record{
			Id:         GetId(data),
			Attributes: data,
		}
		field := Field{Field: "id", Value: record.Id}
		list[field] = append(list[field], record)
	}
	req.SendResponse(&FutureResponse{
		IsSuccess: true,
		Success: map[Future]FutureResponseKind{
			fr: &FutureResponseKindByFields{
				IsSingle: len(k.Ids) == 1,
				Records:  list,
			},
		},
	})
}
