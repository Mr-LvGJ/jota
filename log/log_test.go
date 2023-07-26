package log

import (
	"context"
	"testing"
)

func TestDebug(t *testing.T) {
	type args struct {
		ctx    context.Context
		msg    string
		fields []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewGlobal(&Config{})
			Debug(tt.args.ctx, tt.args.msg, tt.args.fields...)
		})
	}
}
