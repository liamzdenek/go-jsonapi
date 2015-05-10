package jsonapi;

type RecordCreatedStatus uint;

const (
    StatusCreated = 1 << iota
    StatusFailed
    StatusHasUserProvidedId
    StatusUserProvidedIdConflict
    StatusUnsupported
);
