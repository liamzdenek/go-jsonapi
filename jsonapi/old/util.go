package jsonapi;

func Check(err error) {
    if err != nil {
        panic(err);
    }
}

func Reply(resp interface{}) {
    panic(&ResponseReply{Reply:resp});
}
