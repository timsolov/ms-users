package profile

import (
	"context"
	"reflect"
	"testing"

	"github.com/timsolov/ms-users/app/domain/entity"
)

func TestProfileQuery_Do(t *testing.T) {
	type fields struct {
		repo Repository
	}
	type args struct {
		ctx   context.Context
		query *Profile
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantUser entity.User
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := ProfileQuery{
				repo: tt.fields.repo,
			}
			gotUser, err := uc.Do(tt.args.ctx, tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProfileQuery.Do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUser, tt.wantUser) {
				t.Errorf("ProfileQuery.Do() = %v, want %v", gotUser, tt.wantUser)
			}
		})
	}
}
