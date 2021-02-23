package depseeker

import (
	"context"
	"reflect"
	"testing"
	"time"
)

func TestDepseeker_Run(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()
	type fields struct {
		options Options
	}
	type args struct {
		ctx context.Context
		url string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []Dependency
		wantErr bool
	}{
		{
			name:   "normal case",
			fields: fields{},
			args: args{
				ctx: ctx,
				url: "http://localhost:5000",
			},
			want:    []Dependency{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Depseeker{
				options: tt.fields.options,
			}
			got, err := d.Run(tt.args.ctx, tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Depseeker.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Depseeker.Run() = %v, want %v", got, tt.want)
			}
		})
	}
}
