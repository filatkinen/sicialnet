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
	Id         string     `json:"id,omitempty"`
	FirstName  string     `json:"first_name,omitempty"`
	SecondName *string    `json:"second_name,omitempty"`
	BirthDate  *time.Time `json:"birthdate,omitempty"`
	Sex        *string    `json:"sex,omitempty"`
	Biography  *string    `json:"biography,omitempty"`
	City       *string    `json:"city,omitempty"`
}

func UUID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return strings.ToLower(fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])), nil
}
