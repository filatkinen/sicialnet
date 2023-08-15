package storage

import (
	"context"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Storage interface {
	UserAdd(ctx context.Context, user *User) error
	UserDelete(ctx context.Context, userID string) error
	UserUpdate(ctx context.Context, user *User) error
	UserGet(ctx context.Context, userID string) (*User, error)
	UserGetRandom(ctx context.Context) (*User, error)
	UserSearch(ctx context.Context, firstNameMask string, secondNameMask string) ([]*User, error)

	UserAddPost(ctx context.Context, post *Post) (string, error)
	UserAddFriend(ctx context.Context, userID string, friendID string) error
	UserGetFriends(ctx context.Context, userID string) ([]string, error)
	UserGetFriendsPosts(ctx context.Context, userID string, offset int, limit int) ([]*Post, error)

	TokenAdd(ctx context.Context, token *Token) error
	TokenDelete(ctx context.Context, hash string) error
	TokenGet(ctx context.Context, hash string) (*Token, error)
	TokenDeleteAllUser(ctx context.Context, userID string) error

	UserCredentialSet(ctx context.Context, cred *UserCredential) error
	UserCredentialDelete(ctx context.Context, userID string) error
	UserCredentialGet(ctx context.Context, userID string) (*UserCredential, error)

	Close(ctx context.Context) error
}
