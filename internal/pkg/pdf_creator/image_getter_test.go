package pdf_creator

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/whitewolf185/mangaparser/internal/pkg/err_controller"
)

func Test_imageGetter_GetImageAndSave(t *testing.T) {
	ctx := context.Background()
	ec := err_controller.NewErrController()
	folderController := func(folderPathToSave string) error {
		err := os.MkdirAll(folderPathToSave, os.ModePerm)
		if err != nil {
			return err
		}
		return nil
	}
	type args struct {
		wg               *sync.WaitGroup
		folderPathToSave string
		page             int
		url              string
	}
	tests := []struct {
		name    string
		args    args
		cleanUp func(ttt *testing.T, folderToCleanUp string)
	}{
		{
			name: "файл успешно сохранен",
			args: args{
				wg:               &sync.WaitGroup{},
				folderPathToSave: "./test_pic",
				page:             1,
				url:              "https://upload.wikimedia.org/wikipedia/ru/c/cd/JoJos_Bizarre_Adventure.jpg",
			},
			cleanUp: func(ttt *testing.T, folderToCleanUp string) {
				err := os.RemoveAll(folderToCleanUp)
				if err != nil {
					t.Fatal(err)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := folderController(tt.args.folderPathToSave)
			if err != nil {
				t.Error(err)
			}
			t.Cleanup(func() {
				if tt.cleanUp != nil {
					tt.cleanUp(t, tt.args.folderPathToSave)
				}
			})
			ig := imageGetter{}
			ig.GetImageAndSave(ctx, tt.args.wg, ec, tt.args.folderPathToSave, tt.args.page, tt.args.url)
			tt.args.wg.Wait()
			if erCont := ec.IsNul(); erCont != nil {
				t.Errorf("something goes wrong. Errors:\n%v", erCont)
			}
			formatWithoutExt := strings.Replace(formatToSaveFile, ".%s", "", 1)
			formatWithoutExt = formatWithoutExt + ".*"

			matches, err := filepath.Glob(fmt.Sprintf(formatWithoutExt, tt.args.folderPathToSave, tt.args.page))
			if err != nil || matches == nil || len(matches) == 0 {
				t.Errorf("file does not exists%s", err.Error())
			}
		})
	}
}

func Test_imageGetter_extensionParser(t *testing.T) {

	tests := []struct {
		name    string
		url     string
		want    string
		wantErr bool
	}{
		{
			name:    "this is actually an image",
			url:     "filepath.png",
			want:    ".png",
			wantErr: false,
		},
		{
			name:    "this is NOT an image",
			url:     "filepath.pdf",
			want:    "",
			wantErr: true,
		},
		{
			name:    "this is file without extension",
			url:     "filepath",
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ig := imageGetter{}
			got, err := ig.extensionParser(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("extensionParser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("extensionParser() got = %v, want %v", got, tt.want)
			}
		})
	}
}
