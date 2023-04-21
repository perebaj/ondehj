package event

import (
	"context"
	"database/sql"
	"reflect"
	"testing"
)

func TestSQLRepository_Create(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		ctx   context.Context
		event Event
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Event
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &SQLRepository{
				db: tt.fields.db,
			}
			got, err := r.Create(tt.args.ctx, tt.args.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLRepository.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SQLRepository.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}
