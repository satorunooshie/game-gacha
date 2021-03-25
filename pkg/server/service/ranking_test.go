package service

import (
	"errors"
	"reflect"
	"testing"

	"game-gacha/pkg/server/model"

	"github.com/golang/mock/gomock"
)

func Test_rankingService_RankingList(t *testing.T) {
	type args struct {
		userID        string
		startPosition int
		limit         int
	}
	tests := []struct {
		name    string
		args    args
		before  func(mock *mockRepository, args args)
		want    *rankingListResponse
		wantErr bool
	}{
		{
			name: "[正常]ランクインユーザー複数",
			args: args{
				userID:        "aaa",
				startPosition: 1,
				limit:         10,
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
				mock.userRepository.EXPECT().SelectUsersOrderByHighScore(
					args.startPosition-1, args.limit).Return(
					[]*model.User{
						{
							ID:        "UserId1",
							AuthToken: "User1AuthToken",
							Name:      "User1",
							HighScore: 10000,
							Coin:      1000,
						},
						{
							ID:        "UserId2",
							AuthToken: "User2AuthToken",
							Name:      "User2",
							HighScore: 10,
							Coin:      1000,
						},
					}, nil)
			},
			want: &rankingListResponse{
				Ranks: []*rank{
					{
						UserID:   "UserId1",
						UserName: "User1",
						Rank:     1,
						Score:    10000,
					},
					{
						UserID:   "UserId2",
						UserName: "User2",
						Rank:     2,
						Score:    10,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "[正常]ランクインユーザー0",
			args: args{
				userID:        "aaa",
				startPosition: 1,
				limit:         10,
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
				mock.userRepository.EXPECT().SelectUsersOrderByHighScore(
					args.startPosition-1, args.limit).Return(
					[]*model.User{}, nil)
			},
			want: &rankingListResponse{
				Ranks: []*rank{},
			},
			wantErr: false,
		},
		{
			name: "[異常]ログインユーザー取得エラー",
			args: args{
				userID:        "aaa",
				startPosition: 1,
				limit:         10,
			},
			before: func(mock *mockRepository, args args) {
				mock.userRepository.EXPECT().SelectUserByPK(args.userID).Return(
					&model.User{}, errors.New("user not found"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "[異常]ランクインユーザー取得エラー",
			args: args{
				userID:        "aaa",
				startPosition: 1,
				limit:         10,
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
				mock.userRepository.EXPECT().SelectUsersOrderByHighScore(
					args.startPosition-1, args.limit).Return(
					[]*model.User{
						nil,
					}, errors.New("users not found"))
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
			s := NewRankingService(mock.userRepository)
			got, err := s.RankingList(tt.args.userID, tt.args.startPosition, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRankInfoList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRankInfoList() got = %v, want %v", got, tt.want)
			}
		})
	}
}
