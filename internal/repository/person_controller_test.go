package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/whitewolf185/mangaparser/api/domain"
	"github.com/whitewolf185/mangaparser/internal/config"
	"github.com/whitewolf185/mangaparser/internal/gen/manga_parser/public/table"
	customerrors "github.com/whitewolf185/mangaparser/pkg/custom_errors"
)

func Test_personController_GetEmailByID(t *testing.T) {
	db := config.TestPostgres(t)
	perC := NewPersonController(db)
	ctx := context.Background()

	type prepare struct {
		data func(person domain.PersonInfo) uuid.UUID
		cleanUp func(personID uuid.UUID)
	}

	tests := []struct {
		name    string
		prepare prepare
		person  domain.PersonInfo
		want    string
		wantErr error
	}{
		{
			name: "person and email has found",
			prepare: prepare{
				data: func(person domain.PersonInfo) uuid.UUID {
					personID := uuid.MustParse(person.PersonID)
					stmt, args := table.Persons.INSERT(table.Persons.ID, table.Persons.Email).
					VALUES(postgres.UUID(personID), postgres.String("mail.email@mail.ru")).Sql()
					_, err := db.Exec(ctx, stmt, args...)
					if err != nil {
						t.Fatal(err)
					}
					return personID
				},
				cleanUp: func(personID uuid.UUID) {
					stmt, args := table.Persons.DELETE().WHERE(table.Persons.ID.EQ(postgres.UUID(personID))).Sql()
					_, err := db.Exec(ctx, stmt, args...)
					if err != nil {
						t.Fatal(err)
					}
				},
			},
			person: domain.PersonInfo{
				PersonID: uuid.New().String(),
			},
			want: "mail.email@mail.ru",
		},
		{
			name: "person ID not found",
			prepare: prepare{},
			person: domain.PersonInfo{
				PersonID: uuid.New().String(),
			},
			want: "",
			wantErr: customerrors.ErrEmailsNotFound,
		},
		{
			name: "email is empty",
			prepare: prepare{
				data: func(person domain.PersonInfo) uuid.UUID {
					personID := uuid.MustParse(person.PersonID)
					stmt, args := table.Persons.INSERT(table.Persons.ID, table.Persons.Email).
					VALUES(postgres.UUID(personID), postgres.String("")).Sql()
					_, err := db.Exec(ctx, stmt, args...)
					if err != nil {
						t.Fatal(err)
					}
					return personID
				},
				cleanUp: func(personID uuid.UUID) {
					stmt, args := table.Persons.DELETE().WHERE(table.Persons.ID.EQ(postgres.UUID(personID))).Sql()
					_, err := db.Exec(ctx, stmt, args...)
					if err != nil {
						t.Fatal(err)
					}
				},
			},
			person: domain.PersonInfo{
				PersonID: uuid.New().String(),
			},
			want: "",
			wantErr: customerrors.ErrEmailsNotFound,
		},
		{
			name: "found by telegram ID",
			prepare: prepare{
				data: func(person domain.PersonInfo) uuid.UUID {
					personID := uuid.New()
					stmt, args := table.Persons.INSERT(table.Persons.ID, table.Persons.Email, table.Persons.TelegramID).
					VALUES(postgres.UUID(personID), postgres.String("mail.email@mail.ru"), postgres.Int64(person.TelegramID)).Sql()
					_, err := db.Exec(ctx, stmt, args...)
					if err != nil {
						t.Fatal(err)
					}
					return personID
				},
				cleanUp: func(personID uuid.UUID) {
					stmt, args := table.Persons.DELETE().WHERE(table.Persons.ID.EQ(postgres.UUID(personID))).Sql()
					_, err := db.Exec(ctx, stmt, args...)
					if err != nil {
						t.Fatal(err)
					}
				},
			},
			person: domain.PersonInfo{
				TelegramID: 466328545,
			},
			want: "mail.email@mail.ru",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare.data != nil {
				IDToClean := tt.prepare.data(tt.person)
				t.Cleanup(func() {
					tt.prepare.cleanUp(IDToClean)
				})
			}
			got, err := perC.GetEmailByID(ctx, tt.person)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("personController.GetEmailByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("personController.GetEmailByID() = %v, want %v", got, tt.want)
			}
		})
	}
}
