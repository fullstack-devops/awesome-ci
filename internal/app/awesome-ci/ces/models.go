package ces

import "github.com/fullstack-devops/awesome-ci/internal/pkg/rcpersist"

type CES struct {
	Type    rcpersist.CESType // required
	EnvFile string            // required
	OutFile *string
}

type KeyValue struct {
	Name  string
	Value string
}
