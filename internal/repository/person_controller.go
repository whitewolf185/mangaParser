package repository

import (
	"context"
	"fmt"

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

func (perC personController) GetEmailByID(ctx context.Context, personID uuid.UUID) (string, error){
	stmt, args := table.Persons.SELECT(table.Persons.Email).
		WHERE(table.Persons.ID.EQ(postgres.UUID(personID))).Sql()

	var email string
	err := perC.db.QueryRow(ctx, stmt, args...).Scan(&email)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return "", errors.Wrap(customerrors.ErrEmailsNotFound, fmt.Sprintf("person not found ID = %s", personID.String()))
	case err != nil:
		return "", err
	case email == "":
		return "", errors.Wrap(customerrors.ErrEmailsNotFound, fmt.Sprintf("person ID = %s", personID.String()))
	}

	return email, nil
}