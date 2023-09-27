package storage

import (
	"crypto/rand"
	"fmt"
	"strings"
	"time"
)

type Token struct {
	Hash   string
	Expire time.Time
	UserID string
}

type UserCredential struct {
	UserID     string
	PassBcrypt string
}

type User struct {
	Id            string     `json:"id,omitempty"`
	FirstName     string     `json:"firstName,omitempty"`
	SecondName    *string    `json:"secondName,omitempty"`
	BirthDate     *time.Time `json:"birthdate,omitempty"`
	Sex           *string    `json:"sex,omitempty"`
	Biography     *string    `json:"biography,omitempty"`
	City          *string    `json:"city,omitempty"`
	DialogShardId *int
}

type Post struct {
	PostId   string    `json:"postId,omitempty"`
	UserId   string    `json:"friendId,omitempty"`
	PostText string    `json:"post,omitempty"`
	PostDate time.Time `json:"postDate,omitempty"`
}

type DialogMessage struct {
	From string `json:"from"`
	To   string `json:"to"`
	Text string `json:"text"`
}

func UUID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return strings.ToLower(fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])), nil
}
