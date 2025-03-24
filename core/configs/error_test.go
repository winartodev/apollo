package configs

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"testing"
)

func Test_isNonNilAndNotExpectedMigrationError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "success_error_is_no_migration_and_no_change",
			args: args{
				err: errors.New("some error"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isNonNilAndNotExpectedMigrationError(tt.args.err); got != tt.want {
				t.Errorf("isNonNilAndNotExpectedMigrationError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isMigrationNoChange(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "success_error_is_no_migration_and_no_change",
			args: args{
				err: migrate.ErrNoChange,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isErrorNoChange(tt.args.err); got != tt.want {
				t.Errorf("isErrorNoChange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isNoMigration(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "success_error_is_no_migration",
			args: args{
				err: migrate.ErrNilVersion,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isErrorNoMigration(tt.args.err); got != tt.want {
				t.Errorf("isErrorNoMigration() = %v, want %v", got, tt.want)
			}
		})
	}
}
