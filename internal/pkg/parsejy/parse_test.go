package parsejy_test

import (
	"bytes"
	"testing"

	"github.com/fullstack-devops/awesome-ci/internal/pkg/parsejy"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		input      []byte
		syntax     parsejy.Syntax
		wantErr    bool
		wantOutput string
	}{
		{
			name:    "invalid JSON input with valid query",
			query:   ".name",
			input:   []byte(`{name: John`),
			syntax:  parsejy.JSONSyntax,
			wantErr: true,
		},
		{
			name:    "invalid YAML input with valid query",
			query:   ".name",
			input:   []byte(`name= ohn`),
			syntax:  parsejy.YamlSyntax,
			wantErr: true,
		},
		{
			name:    "empty input with valid query",
			query:   ".name",
			input:   []byte{},
			syntax:  parsejy.JSONSyntax,
			wantErr: true,
		},
		{
			name:    "empty query with valid input",
			query:   "",
			input:   []byte(`{"name": "John"}`),
			syntax:  parsejy.JSONSyntax,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			/* origPrintf := fmt.Printf
			myPrintf := func(format string, a ...interface{}) {
				buf.WriteString(fmt.Sprintf(format, a...))
			}
			defer func() {
				panic(origPrintf)
			}() */

			err := parsejy.Parse(tt.query, tt.input, tt.syntax)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if buf.String() != tt.wantOutput {
				t.Errorf("Parse() output = %v, wantOutput %v", buf.String(), tt.wantOutput)
			}
		})
	}
}
