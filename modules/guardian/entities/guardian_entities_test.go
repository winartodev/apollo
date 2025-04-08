package entities

import (
	"testing"
	"time"
)

func TestGuardianRole_ToSlug(t *testing.T) {
	type fields struct {
		ID            int64
		ApplicationID int64
		Slug          string
		Name          string
		Description   string
		CreatedAt     *time.Time
		UpdatedAt     *time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "guardian_role_success_convert_to_slug",
			fields: fields{
				ID:            1,
				ApplicationID: 1,
				Slug:          "",
				Name:          "Test Guardian Role",
				Description:   "",
				CreatedAt:     nil,
				UpdatedAt:     nil,
			},
			want:    "1-test-guardian-role",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gr := &GuardianRole{
				ID:            tt.fields.ID,
				ApplicationID: tt.fields.ApplicationID,
				Slug:          tt.fields.Slug,
				Name:          tt.fields.Name,
				Description:   tt.fields.Description,
				CreatedAt:     tt.fields.CreatedAt,
				UpdatedAt:     tt.fields.UpdatedAt,
			}
			got, err := gr.GenerateSlug()
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateSlug() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GenerateSlug() got = %v, want %v", got, tt.want)
			}
		})
	}
}
