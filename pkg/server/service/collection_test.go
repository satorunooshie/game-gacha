package service

import (
	"errors"
	"reflect"
	"testing"

	"game-gacha/pkg/server/model"

	"github.com/golang/mock/gomock"
)

func Test_collectionService_CollectionList(t *testing.T) {
	type args struct {
		userID string
	}
	tests := []struct {
		name    string
		args    args
		before  func(mock *mockRepository, args args)
		want    *collectionListResponse
		wantErr bool
	}{
		{
			name: "[正常系]コレクション複数",
			args: args{
				userID: "aaa",
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
				mock.userCollectionItemRepository.EXPECT().SelectUserCollectionItems(args.userID).Return(
					[]*model.UserCollectionItem{
						{
							UserID:           "aaa",
							CollectionItemID: "1",
							CreatedAt:        nil,
						},
						{
							UserID:           "aaa",
							CollectionItemID: "2",
							CreatedAt:        nil,
						},
					}, nil)
				mock.collectionItemRepository.EXPECT().SelectAllCollectionItems().Return(
					[]*model.CollectionItem{
						{
							ID:        "1",
							Name:      "No.1",
							Rarity:    1,
							CreatedAt: nil,
						},
						{
							ID:        "2",
							Name:      "No.2",
							Rarity:    1,
							CreatedAt: nil,
						},
						{
							ID:        "3",
							Name:      "No.3",
							Rarity:    1,
							CreatedAt: nil,
						},
					}, nil)
			},
			want: &collectionListResponse{
				Collections: []*collection{
					{
						CollectionID: "1",
						Name:         "No.1",
						Rarity:       1,
						HasItem:      true,
					},
					{
						CollectionID: "2",
						Name:         "No.2",
						Rarity:       1,
						HasItem:      true,
					},
					{
						CollectionID: "3",
						Name:         "No.3",
						Rarity:       1,
						HasItem:      false,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "[正常系]コレクション0",
			args: args{
				userID: "aaa",
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
				mock.userCollectionItemRepository.EXPECT().SelectUserCollectionItems(args.userID).Return(
					[]*model.UserCollectionItem{
						{},
					}, nil)
				mock.collectionItemRepository.EXPECT().SelectAllCollectionItems().Return(
					[]*model.CollectionItem{
						{
							ID:        "1",
							Name:      "No.1",
							Rarity:    1,
							CreatedAt: nil,
						},
						{
							ID:        "2",
							Name:      "No.2",
							Rarity:    1,
							CreatedAt: nil,
						},
						{
							ID:        "3",
							Name:      "No.3",
							Rarity:    1,
							CreatedAt: nil,
						},
					}, nil)
			},
			want: &collectionListResponse{
				Collections: []*collection{
					{
						CollectionID: "1",
						Name:         "No.1",
						Rarity:       1,
						HasItem:      false,
					},
					{
						CollectionID: "2",
						Name:         "No.2",
						Rarity:       1,
						HasItem:      false,
					},
					{
						CollectionID: "3",
						Name:         "No.3",
						Rarity:       1,
						HasItem:      false,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "[異常系]ログインユーザー取得エラー",
			args: args{
				userID: "aaa",
			},
			before: func(mock *mockRepository, args args) {
				mock.userRepository.EXPECT().SelectUserByPK(args.userID).Return(
					&model.User{}, errors.New("user not found"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "[正常系]ユーザーコレクション取得エラー",
			args: args{
				userID: "aaa",
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
				mock.userCollectionItemRepository.EXPECT().SelectUserCollectionItems(args.userID).Return(
					[]*model.UserCollectionItem{
						nil,
					}, errors.New("user collection not found"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "[異常系]マスターコレクション取得エラー",
			args: args{
				userID: "aaa",
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
				mock.userCollectionItemRepository.EXPECT().SelectUserCollectionItems(args.userID).Return(
					[]*model.UserCollectionItem{
						{
							UserID:           "aaa",
							CollectionItemID: "1",
							CreatedAt:        nil,
						},
						{
							UserID:           "aaa",
							CollectionItemID: "2",
							CreatedAt:        nil,
						},
					}, nil)
				mock.collectionItemRepository.EXPECT().SelectAllCollectionItems().Return(
					[]*model.CollectionItem{
						nil,
					}, errors.New("master collection not found"))
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
			s := NewCollectionService(
				mock.userRepository,
				mock.userCollectionItemRepository,
				mock.collectionItemRepository,
			)
			got, err := s.CollectionList(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("CollectionList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CollectionList() got = %v, want %v", got, tt.want)
			}
		})
	}
}
