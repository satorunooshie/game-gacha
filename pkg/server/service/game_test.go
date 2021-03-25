package service

import (
	"errors"
	"game-gacha/pkg/server/model"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

func Test_gameService_GameFinish(t *testing.T) {
	type args struct {
		userID string
		score  int
	}
	tests := []struct {
		name    string
		args    args
		before  func(mock *mockRepository, args args)
		want    *gameFinishResponse
		wantErr bool
	}{
		{
			name: "[正常系]ゲーム終了",
			args: args{
				userID: "aaa",
				score:  1000,
			},
			before: func(mock *mockRepository, args args) {
				mock.userRepository.EXPECT().SelectUserByPK(args.userID).Return(
					&model.User{
						ID:        "aaa",
						AuthToken: "bbb",
						Name:      "ccc",
						HighScore: 0,
						Coin:      0,
					}, nil)
				mock.userRepository.EXPECT().UpdateUserByPK(
					&model.User{
						ID:        "aaa",
						AuthToken: "bbb",
						Name:      "ccc",
						HighScore: 1000,
						Coin:      100,
					}).Return(nil)
			},
			want: &gameFinishResponse{
				GottenCoin: 100,
			},
			wantErr: false,
		},
		{
			name: "[異常系]ログインユーザー取得エラー",
			args: args{
				userID: "aaa",
				score:  1000,
			},
			before: func(mock *mockRepository, args args) {
				mock.userRepository.EXPECT().SelectUserByPK(args.userID).Return(
					&model.User{}, errors.New("user not found"))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := newMockRepository(ctrl)
			tt.before(mock, tt.args)
			s := NewGameService(mock.userRepository)
			got, err := s.GameFinish(tt.args.userID, tt.args.score)
			if (err != nil) != tt.wantErr {
				t.Errorf("GameFinish() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GameFinish() got = %v, want %v", got, tt.want)
			}
		})
	}
}
