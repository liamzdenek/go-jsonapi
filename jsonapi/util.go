package jsonapi;

func Check(e error) {
    if e != nil {
        panic(e);
    }
}

// Reply is a syntax sugar panic button
func Reply(a interface{}) {
    panic(a);
}
