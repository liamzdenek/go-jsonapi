package jsonapi;

type Record interface {
    Ider
    Typer
    Data() *WorkFindLinksByRecordResult
    Include() bool
}
