package jsonapi;

import "fmt";

type ErrorResourceDoesNotExist struct {
    ResourceName string
};

func(e *ErrorResourceDoesNotExist) Error() string {
    return fmt.Sprintf("The provided resource \"%s\" does not exist.", e.ResourceName);
}
