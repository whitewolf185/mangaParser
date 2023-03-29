package mailer

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/whitewolf185/mangaparser/internal/config"
	"github.com/whitewolf185/mangaparser/internal/config/flags"
	"github.com/whitewolf185/mangaparser/internal/pkg/mailer/mock"
)

func TestEbookMailer_SendManga(t *testing.T) {
	flags.InitServiceFlags()
	ctx := context.Background()
	type args struct {
		personID      uuid.UUID
		mangaFilePath string
	}
	tests := []struct {
		name    string
		mockGen func(ctrl *gomock.Controller) EmailGetter
		args    args
		wantErr bool
	}{
		{
			name: "send test file to me",
			mockGen: func(ctrl *gomock.Controller) EmailGetter {
				mm := mock.NewMockEmailGetter(ctrl)
				from := config.GetValue(config.EmailAccount)
				mm.EXPECT().GetEmail(gomock.Any(), gomock.Any()).Return(from, nil)
				return mm
			},
			args: args{
				personID: uuid.New(),
				mangaFilePath: "file_test/hello_test.txt",
			},
		},
	}
	ctrl := gomock.NewController(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			em := NewEbookMailer(tt.mockGen(ctrl))
			if err := em.SendManga(ctx, tt.args.personID, tt.args.mangaFilePath); (err != nil) != tt.wantErr {
				t.Errorf("EbookMailer.SendManga() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
