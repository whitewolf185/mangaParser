package mailer

import (
	"context"
	"testing"

	"github.com/whitewolf185/mangaparser/internal/config"
	"github.com/whitewolf185/mangaparser/internal/config/flags"
)

func TestEbookMailer_SendManga(t *testing.T) {
	flags.InitServiceFlags()
	ctx := context.Background()
	type args struct {
		email      string
		mangaFilePath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "send test file to me",
			args: args{
				email: config.GetValue(config.EmailAccount),
				mangaFilePath: "file_test/hello_test.txt",
			},
		},
		{
			name: "empty email",
			args: args{
				email: "",
				mangaFilePath: "file_test/hello_test.txt",
			},
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			em := NewEbookMailer()
			if err := em.SendManga(ctx, tt.args.email, tt.args.mangaFilePath); (err != nil) != tt.wantErr {
				t.Errorf("EbookMailer.SendManga() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
