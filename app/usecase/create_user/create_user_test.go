package create_user

import (
	"context"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/timsolov/ms-users/app/domain/repository"
)

func TestCreateUserCommand_Do(t *testing.T) {
	type fields struct {
		repo repository.UserRepository
	}
	type args struct {
		ctx context.Context
		cmd *CreateUser
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantUserID uuid.UUID
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := CreateUserCommand{
				repo: tt.fields.repo,
			}
			gotUserID, err := uc.Do(tt.args.ctx, tt.args.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUserCommand.Do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUserID, tt.wantUserID) {
				t.Errorf("CreateUserCommand.Do() = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}
