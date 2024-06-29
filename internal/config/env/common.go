package env

import (
	"fmt"
)

func notFoundError(variableName string) error {
	return fmt.Errorf("`%s` is not set in env", variableName)
}
