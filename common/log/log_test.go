package log

import (
	"reflect"
	"testing"

	"go.uber.org/zap"
)

func TestSetFieldString(t *testing.T) {
	type args struct {
		param map[string]string
	}
	tests := []struct {
		name       string
		args       args
		wantFields []Field
	}{
		{
			name: "success",
			args: args{param: map[string]string{"code": "200"}},
			wantFields: []Field{
				{Field: zap.String("code", "200")},
			},
		},
		{
			name: "empty",
			args: args{param: map[string]string{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if gotFields := SetFieldString(tt.args.param); !reflect.DeepEqual(gotFields, tt.wantFields) {
				t.Errorf("SetFieldString() = %v, want %v", gotFields, tt.wantFields)
			}
		})
	}
}
