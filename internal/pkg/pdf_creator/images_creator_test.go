package pdf_creator

import (
	"context"
	"testing"

	"github.com/whitewolf185/mangaparser/internal/pkg/err_controller"
)

func Test_imageController_GetImagesFromURLs(t *testing.T) {
	type fields struct {
		imageGetter   ImageGetter
		errController *err_controller.ErrController
	}
	type args struct {
		ctx              context.Context
		folderPathToSave string
		urlImages        []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ic := imageController{
				imageGetter:   tt.fields.imageGetter,
				errController: tt.fields.errController,
			}
			if err := ic.GetImagesFromURLs(tt.args.ctx, tt.args.folderPathToSave, tt.args.urlImages); (err != nil) != tt.wantErr {
				t.Errorf("GetImagesFromURLs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
