package service

import (
	"errors"
	"reflect"
	"testing"

	"game-gacha/pkg/server/model"

	"github.com/golang/mock/gomock"
)

// TODO: Test_userService_UserCreate

func Test_userService_UserGet(t *testing.T) {
	type args struct {
		userID string
	}
	tests := []struct {
		name    string
		args    args
		before  func(mock *mockRepository, args args)
		want    *userGetResponse
		wantErr bool
	}{
		{
			name: "[正常系]ユーザ取得",
			args: args{
				userID: "aaa",
			},
			before: func(mock *mockRepository, args args) {
				mock.userRepository.EXPECT().SelectUserByPK(args.userID).Return(
					&model.User{
						ID:        args.userID,
						Name:      "ccc",
						HighScore: 0,
						Coin:      0,
					}, nil)
			},
			want: &userGetResponse{
				ID:        "aaa",
				Name:      "ccc",
				HighScore: 0,
				Coin:      0,
			},
			wantErr: false,
		},
		{
			name: "[異常系]ユーザ取得エラー",
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := newMockRepository(ctrl)
			tt.before(mock, tt.args)
			s := NewUserService(mock.userRepository)
			got, err := s.UserGet(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserGet() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userService_UserUpdate(t *testing.T) {
	type args struct {
		userID string
		name   string
	}
	tests := []struct {
		name    string
		args    args
		before  func(mock *mockRepository, args args)
		wantErr bool
	}{
		{
			name: "[正常系]ユーザ名更新",
			args: args{
				userID: "aaa",
				name:   "updated",
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
						Name:      "updated",
						HighScore: 0,
						Coin:      0,
					}).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "[異常系]ログインユーザ取得エラー",
			args: args{
				userID: "aaa",
				name:   "updated",
			},
			before: func(mock *mockRepository, args args) {
				mock.userRepository.EXPECT().SelectUserByPK(args.userID).Return(
					&model.User{}, errors.New("user not found"))
			},
			wantErr: true,
		},
		{
			name: "[異常系]ユーザ名更新エラー",
			args: args{
				userID: "aaa",
				name:   "updated",
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
						Name:      "updated",
						HighScore: 0,
						Coin:      0,
					}).Return(errors.New("failed to update user"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := newMockRepository(ctrl)
			tt.before(mock, tt.args)
			s := NewUserService(mock.userRepository)
			if err := s.UserUpdate(tt.args.userID, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("UserUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
