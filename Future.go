package jsonapi;

type Future interface{
    ShouldCombine(Future) bool
    Combine(Future) error
    Work(ef *ExecutableFuture)
}
