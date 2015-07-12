package jsonapi;

import "fmt"
import "time"

var FutureOutputTimeout time.Duration = 10 * time.Second;

type FutureOutput struct {
    Parents []*ExecutableFuture
    PrimaryData Future
    PrimaryDataType OutputType
}

func(fo *FutureOutput) PushParent(ef *ExecutableFuture) { fo.Parents = append(fo.Parents, ef); }
func(fo *FutureOutput) ShouldCombine(f Future) bool { return false; }
func(fo *FutureOutput) Combine(f Future) error { panic("This should never be called"); }

func(fo *FutureOutput) Work(pf *ExecutableFuture) {
    reply_todo := []*FutureRequest{};
    responses := map[Future]FutureResponseKind{};
    need := fo.Parents;
    OUTER:for {
        var rawreq *FutureRequest;
        rawreq = pf.GetRequest();
        fmt.Printf("GOT RAW REQ: %#v\n", rawreq);
        reply_todo = append(reply_todo, rawreq);
        switch req := rawreq.Kind.(type) {
        case *FutureRequestKindIdentity:
            responses[req.Future] = req.Response;
            for i, child := range need {
                if(child.Future == req.Future) {
                    need = append(need[:i], need[i+1:]...);
                }
            }
            //fmt.Printf("STILL NEED: %#v\n", need);
            if len(need) == 0 {
                break OUTER;
            }
        default:
            panic(TODO());
        }
    }

    //fmt.Printf("OUTPUT GOT DATA: %#v\n", responses);
    
    var output_relationship *ORelationship = nil;
    output := NewOutput();
    output_primary := []*Record{};
    output_included := []*Record{};
    is_single := false;
    //var output_relationship *ORelationship = nil;

    for future, reskind := range responses {
        switch res := reskind.(type) {
        case FutureResponseKindWithRecords:
            if future == fo.PrimaryData {
                is_single = res.GetIsSingle();
                output_primary = append(output_primary, res.GetRecords()...);
            } else {
                output_included = append(output_included, res.GetRecords()...);
            }
        default:
            panic(fmt.Sprintf("Future got unsupported reskind %T", reskind));
        }
    }
    if(fo.PrimaryDataType == OutputTypeResources) {
        output.Data = &ORecords{
            IsSingle: is_single,
            Records: output_primary,
        }
    } else {
        output.Data = output_relationship;
    }
    output.Included = output_included;
    panic(output);
}
