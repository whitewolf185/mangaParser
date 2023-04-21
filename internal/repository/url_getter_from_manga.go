package repository

import (
	"context"
	"fmt"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"

	"github.com/whitewolf185/mangaparser/internal/config"
	"github.com/whitewolf185/mangaparser/internal/gen/manga_parser/public/enum"
	"github.com/whitewolf185/mangaparser/internal/gen/manga_parser/public/table"
	customerrors "github.com/whitewolf185/mangaparser/pkg/custom_errors"
)

// urlRepo структура получения ссылка на главу для конкретной манги. В данном случае идет речь про mangalib, так как у mangadex будет другой подход
type urlRepo struct {
	db *pgx.Conn
}

func NewUrlRepo(db *pgx.Conn) urlRepo {
	return urlRepo{db: db}
}

func (ur urlRepo) GetUrlByID(ctx context.Context, mangaID uuid.UUID, sourceType config.MangaSourceType) (string, error) {
	stmt, args := table.MangaInfo.SELECT(
		table.MangaInfo.URL,
	).WHERE(
		table.MangaInfo.ID.EQ(postgres.UUID(mangaID)).
			AND(table.MangaInfo.SourceType.EQ(ur.getSourceTypeEnum(sourceType))),
	).Sql()

	var url string
	err := ur.db.QueryRow(ctx, stmt, args...).Scan(&url)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return "", errors.Wrap(customerrors.ErrUrlIsEmpty, fmt.Sprintf("manga not found ID = %s", mangaID.String()))
	case err != nil:
		return "", err
	case url == "":
		return "", errors.Wrap(customerrors.ErrUrlIsEmpty, fmt.Sprintf("manga ID = %s", mangaID.String()))
	}

	return url, nil
}

func (ur urlRepo) getSourceTypeEnum(sourceType config.MangaSourceType) postgres.StringExpression {
	switch sourceType {
	case config.MangaLib:
		return enum.MangaSourceType.Mangalib
	case config.MangaDex:
		return enum.MangaSourceType.Mangadex
	}
	return postgres.String("")
}
