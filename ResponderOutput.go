package jsonapi;

import ("net/http";)

type ResponderOutput struct{
    Output *Output
}

func NewResponderOutput(o *Output) *ResponderOutput {
    return &ResponderOutput{
        Output: o,
    }
}

func(ro *ResponderOutput) Respond(a *API, w http.ResponseWriter, r *http.Request) error {
    ro.Output.Prepare();
    a.Send(ro.Output,w)
    return nil;
}
