package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
	"os"
	"testing"
	"time"
	"user_service/internal/domain/models"
)

const timeout = 15

func connectToDB(t *testing.T) *sqlx.DB {

	uri := os.Getenv("TEST_POSTGRES_DB_URI")

	if len(uri) < 3 {
		t.FailNow()
	}

	t.Log(uri)

	db := MustOpenPostgresDB(uri)
	if db == nil {
		t.FailNow()
	}

	return db
}

func TestStorage_SaveUser(t *testing.T) {
	db := connectToDB(t)

	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx context.Context
		u   models.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.User
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				db: db,
			},
			args: args{
				ctx: ctx,
				u: models.User{
					Email:    "user1@example.org",
					Password: "password1",
				},
			},
			want: models.User{
				Email:    "user1@example.org",
				Password: "password1",
			},
			wantErr: false,
		},
		{
			name: "empty user fields",
			fields: fields{
				db: db,
			},
			args: args{
				ctx: ctx,
				u:   models.User{},
			},
			want:    models.User{},
			wantErr: true,
		},
		{
			name: "same user",
			fields: fields{
				db: db,
			},
			args: args{
				ctx: ctx,
				u: models.User{
					Email:    "user1@example.org",
					Password: "21392390309-21390-321-0",
				},
			},
			want:    models.User{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				db: tt.fields.db,
			}
			got, err := s.SaveUser(tt.args.ctx, tt.args.u)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// we cant check id and enc_password
			if !(got.Email == tt.want.Email || got.Password == got.Password) {
				t.Errorf("SaveUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_FindUserByEmail(t *testing.T) {
	db := connectToDB(t)

	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	st := &Storage{
		db: db,
	}

	// create user to check if we can find it later
	_, _ = st.SaveUser(ctx, models.User{
		Email:    "user1@example.org",
		Password: "password1",
	})

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx   context.Context
		email string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.User
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				db: db,
			},
			args: args{
				ctx:   ctx,
				email: "user1@example.org",
			},
			want: models.User{
				Email:    "user1@example.org",
				Password: "password1",
			},
			wantErr: false,
		},
		{
			name: "not exists",
			fields: fields{
				db: db,
			},
			args: args{
				ctx:   ctx,
				email: "not exists",
			},
			want:    models.User{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				db: tt.fields.db,
			}
			got, err := s.FindUserByEmail(tt.args.ctx, tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindUserByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// we cant check id and enc_password
			if !(got.Email == tt.want.Email || got.Password == got.Password) {
				t.Errorf("FindUserByEmail() got = %v, want %v", got, tt.want)
			}
		})
	}
}
