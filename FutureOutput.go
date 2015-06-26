package jsonapi;

import "fmt"
import "time"

var FutureOutputTimeout time.Duration = 10 * time.Second;

type FutureOutput struct {
    PrimaryData Future
}

func(fo *FutureOutput) ShouldCombine(f Future) bool { return false; }
func(fo *FutureOutput) Combine(f Future) error { panic("This should never be called"); }

func(fo *FutureOutput) Work(pf *PreparedFuture) {
    reply_todo := []*FutureRequest{};
    should_break := false;
    responses := map[Future]FutureResponseKind{};
    need := pf.Parents;
    OUTER:for {
        var rawreq *FutureRequest;
        rawreq, should_break = pf.GetNext();
        if should_break {
            return;
        }
        //fmt.Printf("GOT RAW REQ: %#v\n", rawreq);
        reply_todo = append(reply_todo, rawreq);
        switch req := rawreq.Kind.(type) {
        case *FutureRequestKindIdentity:
            if !req.Response.IsSuccess {
                panic(TODO());
            }
            for future, response := range req.Response.Success {
                responses[future] = response;
                for i, child := range need {
                    if(child.Future == future) {
                        need = append(need[:i], need[i+1:]...);
                    }
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
    
    output := NewOutput();
    output_primary := []*Record{};
    output_included := []*Record{};
    is_single := false;
    //var output_relationship *ORelationship = nil;

    for future, reskind := range responses {
        switch res := reskind.(type) {
        case FutureResponseKindRecords:
            if future == fo.PrimaryData {
                is_single = res.IsSingle;
                output_primary = append(output_primary, res.Records...);
            } else {
                output_included = append(output_included, res.Records...);
            }
        default:
            panic(fmt.Sprintf("Unsupported reskind %#v", reskind));
        }
    }
    output.Data = &ORecords{
        IsSingle: is_single,
        Records: output_primary,
    }
    output.Included = output_included;
    panic(output);
}
