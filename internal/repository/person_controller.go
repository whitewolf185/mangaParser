package repository

import (
	"context"
	"fmt"

	"github.com/whitewolf185/mangaparser/api/domain"
	"github.com/whitewolf185/mangaparser/internal/gen/manga_parser/public/table"
	customerrors "github.com/whitewolf185/mangaparser/pkg/custom_errors"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type personController struct {
	db *pgx.Conn
}

func NewPersonController(db *pgx.Conn) personController {
	return personController{
		db: db,
	}
}

func (perC personController) GetEmailByID(ctx context.Context, person domain.PersonInfo) (string, error){
	selectStmt := table.Persons.SELECT(table.Persons.Email)
	var whereStmt postgres.BoolExpression

	personID, _ := uuid.Parse(person.PersonID)
	switch {
	case personID != uuid.UUID{} && personID != uuid.Nil:
		whereStmt = table.Persons.ID.EQ(postgres.UUID(personID))
	case person.TelegramID > 0:
		whereStmt = table.Persons.TelegramID.EQ(postgres.Int64(person.TelegramID))
	default:
		return "", customerrors.ErrEmailsNotFound
	}
	stmt, args := selectStmt.WHERE(whereStmt).Sql()

	var email string
	err := perC.db.QueryRow(ctx, stmt, args...).Scan(&email)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return "", errors.Wrap(customerrors.ErrEmailsNotFound, fmt.Sprintf("person not found %v", person))
	case err != nil:
		return "", err
	case email == "":
		return "", errors.Wrap(customerrors.ErrEmailsNotFound, fmt.Sprintf("person %v", person))
	}

	return email, nil
}