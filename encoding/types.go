package encoding

import "fmt"

type ErrorIncorrectTag struct {
	Expected uint8
	Given    uint8
}

func (e *ErrorIncorrectTag) Error() string {
	return fmt.Sprintf("Incorrect tag %d, expected %d.", e.Given, e.Expected)
}
