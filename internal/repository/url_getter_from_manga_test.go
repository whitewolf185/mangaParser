package repository

import (
	"context"
	"testing"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/whitewolf185/mangaparser/internal/config"
	"github.com/whitewolf185/mangaparser/internal/gen/manga_parser/public/enum"
	"github.com/whitewolf185/mangaparser/internal/gen/manga_parser/public/table"
	customerrors "github.com/whitewolf185/mangaparser/pkg/custom_errors"
)

func Test_urlRepo_GetUrlByID(t *testing.T) {
	db := config.TestPostgres(t)
	ur := NewUrlRepo(db)
	ctx := context.Background()

	type args struct {
		mangaID    uuid.UUID
		sourceType config.MangaSourceType
	}
	type prepare struct {
		data    func(mangaID uuid.UUID)
		cleanUp func(mangaID uuid.UUID)
	}

	tests := []struct {
		name    string
		prepare prepare
		args    args
		want    string
		wantErr error
	}{
		{
			name: "url has found",
			prepare: prepare{
				data: func(mangaID uuid.UUID) {
					stmt, args := table.MangaInfo.INSERT(
						table.MangaInfo.ID,
						table.MangaInfo.URL,
						table.MangaInfo.SourceType,
						table.MangaInfo.MangaName,
					).VALUES(
						postgres.UUID(mangaID),
						postgres.String("url.uuufdal/jfoore?chapters=true"),
						enum.MangaSourceType.Mangalib,
						postgres.String("manga about merried gay pair"),
					).Sql()
					_, err := db.Exec(ctx, stmt, args...)
					if err != nil {
						t.Fatal(err)
					}
				},
				cleanUp: func(mangaID uuid.UUID) {
					stmt, args := table.MangaInfo.DELETE().WHERE(table.MangaInfo.ID.EQ(postgres.UUID(mangaID))).Sql()
					_, err := db.Exec(ctx, stmt, args...)
					if err != nil {
						t.Fatal(err)
					}
				},
			},
			args: args{
				mangaID:    uuid.New(),
				sourceType: config.MangaLib,
			},
			want: "url.uuufdal/jfoore?chapters=true",
		},
		{
			name: "manga not found",
			args: args{
				mangaID:    uuid.New(),
				sourceType: config.MangaLib,
			},
			wantErr: customerrors.ErrUrlIsEmpty,
		},
		{
			name: "url is empty",
			prepare: prepare{
				data: func(mangaID uuid.UUID) {
					stmt, args := table.MangaInfo.INSERT(
						table.MangaInfo.ID,
						table.MangaInfo.URL,
						table.MangaInfo.SourceType,
						table.MangaInfo.MangaName,
					).VALUES(
						postgres.UUID(mangaID),
						postgres.String(""),
						enum.MangaSourceType.Mangalib,
						postgres.String("manga about merried gay pair"),
					).Sql()
					_, err := db.Exec(ctx, stmt, args...)
					if err != nil {
						t.Fatal(err)
					}
				},
				cleanUp: func(mangaID uuid.UUID) {
					stmt, args := table.MangaInfo.DELETE().WHERE(table.MangaInfo.ID.EQ(postgres.UUID(mangaID))).Sql()
					_, err := db.Exec(ctx, stmt, args...)
					if err != nil {
						t.Fatal(err)
					}
				},
			},
			args: args{
				mangaID:    uuid.New(),
				sourceType: config.MangaLib,
			},
			want:    "",
			wantErr: customerrors.ErrUrlIsEmpty,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare.data != nil {
				tt.prepare.data(tt.args.mangaID)
				t.Cleanup(func() {
					tt.prepare.cleanUp(tt.args.mangaID)
				})
			}
			got, err := ur.GetUrlByID(ctx, tt.args.mangaID, tt.args.sourceType)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("urlRepo.GetUrlByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("urlRepo.GetUrlByID() = %v, want %v", got, tt.want)
			}
		})
	}
}
