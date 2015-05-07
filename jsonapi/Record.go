package jsonapi;

type Record interface {
    Ider
    Linker
    Typer
    Include() bool
}
