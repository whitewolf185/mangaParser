package pdfmerger

import (
	"os"
	"reflect"
	"testing"

	"github.com/pkg/errors"
	customerrors "github.com/whitewolf185/mangaparser/pkg/custom_errors"
)

func Test_getImagesPathStr(t *testing.T) {
	emptyDirPath := "empty_dir_test"
	tests := []struct {
		name          string
		imagesDirPath string
		prepare       func()
		cleanUp       func()
		want          []string
		wantErr       error
	}{
		{
			name:          "all good",
			imagesDirPath: "imgs_1test",
			want:          []string{"imgs_1test/1.jpg", "imgs_1test/2.jpg", "imgs_1test/12.jpg"},
			wantErr:       nil,
		},
		{
			name:          "empty string err",
			imagesDirPath: "",
			want:          nil,
			wantErr:       customerrors.ErrEmptyStr,
		},
		{
			name:          "empty files err",
			imagesDirPath: emptyDirPath,
			prepare: func() {
				err := os.Mkdir(emptyDirPath, os.ModePerm)
				if err != nil {
					t.Fatal(err)
				}
			},
			cleanUp: func() {
				err := os.RemoveAll(emptyDirPath)
				if err != nil {
					t.Fatal(err)
				}
			},
			want:    nil,
			wantErr: customerrors.ErrEmptyDir,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare()
			}
			if tt.cleanUp != nil {
				t.Cleanup(tt.cleanUp)
			}
			got, err := GetImagesPathStr(tt.imagesDirPath)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("getImagesPathStr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getImagesPathStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreatePDFFromImagesDir(t *testing.T) {
	type args struct {
		imagesDirPath string
		outputPath    string
	}
	tests := []struct {
		name    string
		args    args
		cleanUp func()
		wantErr bool
	}{
		{
			name: "all good",
			args: args{
				imagesDirPath: "imgs_1test",
				outputPath:    "out.pdf",
			},
			cleanUp: func() {
				os.RemoveAll("out.pdf")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreatePDFFromImagesDir(tt.args.imagesDirPath, tt.args.outputPath); (err != nil) != tt.wantErr {
				t.Errorf("CreatePDFFromImagesDir() error = %v, wantErr %v", err, tt.wantErr)
			}
			if _, err := os.Stat(tt.args.outputPath); err != nil {
				t.Fatalf("out put file does not exists %s", err.Error())
			}
			if tt.cleanUp != nil {
				t.Cleanup(tt.cleanUp)
			}
		})
	}
}
