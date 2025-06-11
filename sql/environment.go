package sql

import (
	"fmt"

	"example.com/termquery/utils"
)

func GetDatabricksEnvVar(variableName string, envFunc utils.GetEnvFunc) string {
	envVar := envFunc(variableName)
	if envVar == "" {
		panic(fmt.Errorf("%s must be set", variableName))
	} else {
		return envVar
	}
}
