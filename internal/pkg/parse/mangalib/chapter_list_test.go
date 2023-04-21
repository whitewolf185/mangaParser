package mangalib

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"

	"github.com/whitewolf185/mangaparser/internal/config/flags"
	"github.com/whitewolf185/mangaparser/internal/pkg/parse/mangalib/mock"
)

func Test_mangaLibController_GetChapterListUrl(t *testing.T) {
	flags.InitServiceFlags()
	ctx := context.Background()
	type args struct {
		mangaID uuid.UUID
	}
	tests := []struct {
		name    string
		mockGen func(ctrl *gomock.Controller) UrlGetter
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "good chapter parse",
			mockGen: func(ctrl *gomock.Controller) UrlGetter {
				mocker := mock.NewMockUrlGetter(ctrl)
				mocker.EXPECT().GetUrlByID(gomock.Any(), gomock.Any(), gomock.Any()).
					Return("https://mangalib.me/dugeundugeun-gonglyaggi?section=chapters", nil)
				return mocker
			},
			args: args{
				mangaID: uuid.New(),
			},
			want: []string{
				"https://mangalib.me/dugeundugeun-gonglyaggi/v1/c1",
				"https://mangalib.me/dugeundugeun-gonglyaggi/v1/c2",
				"https://mangalib.me/dugeundugeun-gonglyaggi/v1/c3",
				"https://mangalib.me/dugeundugeun-gonglyaggi/v1/c4",
				"https://mangalib.me/dugeundugeun-gonglyaggi/v1/c5",
				"https://mangalib.me/dugeundugeun-gonglyaggi/v1/c6",
				"https://mangalib.me/dugeundugeun-gonglyaggi/v1/c7",
				"https://mangalib.me/dugeundugeun-gonglyaggi/v1/c8",
				"https://mangalib.me/dugeundugeun-gonglyaggi/v1/c9",
				"https://mangalib.me/dugeundugeun-gonglyaggi/v1/c10",
				"https://mangalib.me/dugeundugeun-gonglyaggi/v1/c11",
				"https://mangalib.me/dugeundugeun-gonglyaggi/v1/c12",
				"https://mangalib.me/dugeundugeun-gonglyaggi/v1/c13",
				"https://mangalib.me/dugeundugeun-gonglyaggi/v1/c14",
				"https://mangalib.me/dugeundugeun-gonglyaggi/v1/c15",
				"https://mangalib.me/dugeundugeun-gonglyaggi/v1/c16",
				"https://mangalib.me/dugeundugeun-gonglyaggi/v1/c17",
				"https://mangalib.me/dugeundugeun-gonglyaggi/v1/c18",
				"https://mangalib.me/dugeundugeun-gonglyaggi/v1/c19",
				"https://mangalib.me/dugeundugeun-gonglyaggi/v1/c20",
				"https://mangalib.me/dugeundugeun-gonglyaggi/v1/c21",
				"https://mangalib.me/dugeundugeun-gonglyaggi/v1/c22",
				"https://mangalib.me/dugeundugeun-gonglyaggi/v1/c23",
				"https://mangalib.me/dugeundugeun-gonglyaggi/v1/c24",
				"https://mangalib.me/dugeundugeun-gonglyaggi/v1/c25",
				"https://mangalib.me/dugeundugeun-gonglyaggi/v1/c26",
				"https://mangalib.me/dugeundugeun-gonglyaggi/v1/c27",
				"https://mangalib.me/dugeundugeun-gonglyaggi/v1/c28",
				"https://mangalib.me/dugeundugeun-gonglyaggi/v1/c29",
				"https://mangalib.me/dugeundugeun-gonglyaggi/v1/c30",
				"https://mangalib.me/dugeundugeun-gonglyaggi/v1/c30.5",
			},
		},
	}

	asserter := func(got []string, want []string) bool {
		for i := range want {
			if want[i] != got[i] {
				return false
			}
		}
		return true
	}

	ctrl := gomock.NewController(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mlc, err := NewMangaLibController(tt.mockGen(ctrl))
			if err != nil {
				t.Fatal(err)
			}
			got, err := mlc.GetChapterListUrlByMangaID(ctx, tt.args.mangaID)
			if (err != nil) != tt.wantErr {
				t.Errorf("mangaLibController.GetChapterListUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !asserter(got, tt.want) {
				t.Errorf("mangaLibController.GetChapterListUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
