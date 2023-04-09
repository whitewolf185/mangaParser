package mangalib

import (
	"context"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/whitewolf185/mangaparser/internal/config/flags"
	"github.com/whitewolf185/mangaparser/internal/pkg/parse/mangalib/mock"
)

func Test_mangaLibController_GetPicsUrlInChapter(t *testing.T) {
	flags.InitServiceFlags()
	ctx := context.Background()
	type args struct {
		chapterUrl string
	}
	tests := []struct {
		name    string
		mockPrepare func (ctrl *gomock.Controller) UrlGetter
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "Проверяем мангу",
			mockPrepare: func(ctrl *gomock.Controller) UrlGetter {
				mocker := mock.NewMockUrlGetter(ctrl)
				return mocker
			},
			args: args{
				chapterUrl: "https://mangalib.me/agi-shinsuro-hwansaenghassspnida/v1/c11?bid=12649&page=1",
			},
			want: []string{
				"https://img33.imgslib.link//manga/agi-shinsuro-hwansaenghassspnida/chapters/2261954/01-result-waifu2x-1x-3n-png_ojSP.png",
				"https://img33.imgslib.link//manga/agi-shinsuro-hwansaenghassspnida/chapters/2261954/02-result-waifu2x-1x-3n-png_zouv.png", 
				"https://img33.imgslib.link//manga/agi-shinsuro-hwansaenghassspnida/chapters/2261954/03-result-waifu2x-1x-3n-png_196x.png", 
				"https://img33.imgslib.link//manga/agi-shinsuro-hwansaenghassspnida/chapters/2261954/04-result-waifu2x-1x-3n-png_h78m.png", 
				"https://img33.imgslib.link//manga/agi-shinsuro-hwansaenghassspnida/chapters/2261954/05-result-waifu2x-1x-3n-png_XfE5.png", 
				"https://img33.imgslib.link//manga/agi-shinsuro-hwansaenghassspnida/chapters/2261954/06-result-waifu2x-1x-3n-png_uBte.png",
			},
		},
	}
	ctrl := gomock.NewController(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mlb, err := NewMangaLibController(tt.mockPrepare(ctrl))
			if err != nil {
				t.Fatal(err)
			}
			got, err := mlb.GetPicsUrlInChapter(ctx, tt.args.chapterUrl)
			if (err != nil) != tt.wantErr {
				t.Errorf("mangaLibController.GetPicsUrlInChapter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mangaLibController.GetPicsUrlInChapter() = %v, want %v", got, tt.want)
			}
		})
	}
}
