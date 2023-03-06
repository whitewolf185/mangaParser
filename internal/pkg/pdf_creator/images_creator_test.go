package pdf_creator

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/whitewolf185/mangaparser/internal/pkg/err_controller"
	"github.com/whitewolf185/mangaparser/internal/pkg/pdf_creator/mock"
)

func Test_imageController_GetImagesFromURLs(t *testing.T) {
	ctx := context.Background()
	type mockPrepare struct {
		imageGetter   func(ctrl *gomock.Controller) ImageGetter
		errController *err_controller.ErrController
	}
	type args struct {
		folderPathToSave string
		urlImages        []string
	}
	tests := []struct {
		name    string
		mockPrepare  mockPrepare
		args    args
		wantErr bool
	}{
		{
			name: "one url succesufully downloaded",
			mockPrepare: mockPrepare{
				imageGetter: func(ctrl *gomock.Controller) ImageGetter {
					ig := mock.NewMockImageGetter(ctrl)
					ig.EXPECT().GetImageAndSave(gomock.Any(), gomock.Any(),gomock.Any(),gomock.Any(),gomock.Any(), gomock.Any()).
					Do(func (_ context.Context, wg *sync.WaitGroup, _ *err_controller.ErrController, _ string, _ int, _ string)  {
						wg.Add(1)
						wg.Done()
					})
					return ig
				},
				errController: err_controller.NewErrController(),
			},
			args: args{
				folderPathToSave: "./images_creator_test",
				urlImages: []string{"url.com"},
			},
		},
	}
	ctrl := gomock.NewController(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ec := tt.mockPrepare.errController
			ic := imageController{
				imageGetter:   tt.mockPrepare.imageGetter(ctrl),
				errController: ec,
			}
			if err := ic.GetImagesFromURLs(ctx, tt.args.folderPathToSave, tt.args.urlImages); (err != nil) != tt.wantErr {
				t.Errorf("GetImagesFromURLs() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Cleanup(func() {
				err := os.RemoveAll(tt.args.folderPathToSave)
				if err != nil {
					t.Fatal(err)
				}
			})
		})
	}
}
