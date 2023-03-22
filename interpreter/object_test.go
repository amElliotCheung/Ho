package interpreter

import (
	"reflect"
	"testing"
)

func Test_randomObject(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want Object
	}{
		{"int test",
			"int",
			&Integer{Value: 1},
		},
		{"string test",
			"string",
			&Integer{Value: 2},
		},
		{"bool test",
			"bool",
			&Integer{Value: 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := randomObject(tt.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("randomObject() = %v, want %v", got, tt.want)
			}
		})
	}
}
