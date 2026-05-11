package ces

import "github.com/layer87-labs/relctl/internal/pkg/rcpersist"

type CES struct {
	Type    rcpersist.CESType // required
	EnvFile string            // required
	OutFile *string
}

type KeyValue struct {
	Name  string
	Value string
}
