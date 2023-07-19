package internalhttp_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	internalhttp "github.com/filatkinen/socialnet/internal/server/http"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
	"time"
)

func newStringPointer(s string) *string {
	return &s
}

func toJSON(value any) string {
	data, err := json.Marshal(value)
	if err != nil {
		return ""
	}
	return string(data)
}

var (
	users = []*internalhttp.UserRegisterBody{
		{
			FirstName:  "Ivan1",
			SecondName: newStringPointer("Frolov"),
			Birthdate:  newStringPointer("2002-02-11"),
			Biography:  newStringPointer("Hokkey"),
			City:       newStringPointer("Moskva"),
			Password:   "passI",
		},
		{
			FirstName:  "Masha1",
			SecondName: newStringPointer("Frolova"),
			Birthdate:  newStringPointer("2003-02-11"),
			Biography:  newStringPointer("Dance"),
			City:       newStringPointer("Piter"),
			Password:   "passM",
		},
	}
)

func testHTTPStatus(t *testing.T, servicePort string) { //nolint
	URL := fmt.Sprintf("http://localhost:%s/unknown", servicePort)
	resp, err := http.Get(URL) //nolint
	defer resp.Body.Close()    //nolint
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	URL = fmt.Sprintf("http://localhost:%s/", servicePort)
	resp, err = http.Get(URL) //nolint
	defer resp.Body.Close()   //nolint
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func testHTTPAPI(t *testing.T, servicePort string) { //nolint

	client := http.Client{
		Timeout: time.Second * 2,
	}
	URL := fmt.Sprintf("http://localhost:%s", servicePort)

	usersID := []string{}
	usersTokens := []string{}

	// add users
	for i := range users {
		fullURL := URL + "/user/register"
		userID, err := testUserRegister(t, client, users[i], fullURL)
		require.NoError(t, err)
		require.NotEqual(t, userID, "")
		usersID = append(usersID, userID)
	}

	// login users
	for i := range users {
		fullURL := URL + "/login"
		lb := internalhttp.LoginBody{
			Id:       usersID[i],
			Password: users[i].Password,
		}
		token, err := testUserLogin(t, client, &lb, fullURL)
		require.NoError(t, err)
		require.NotEqual(t, token, "")
		//fmt.Println(token)
		usersTokens = append(usersTokens, token) //nolint
	}
	// get users
	for i := range usersID {
		fullURL := URL + "/user/get/" + usersID[i]
		userGet, err := testUserGet(t, client, fullURL)
		require.NoError(t, err)
		require.Equal(t, users[i].FirstName, userGet.FirstName)
		//fmt.Println(toJSON(userGet))
	}
}

func testUserRegister(t *testing.T, client http.Client, user *internalhttp.UserRegisterBody, url string) (string, error) { //nolint
	data, err := json.Marshal(user)
	require.NoError(t, err)
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	require.NoError(t, err)
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var userGet internalhttp.UserCreateResponse
	jsDecoder := json.NewDecoder(resp.Body)
	err = jsDecoder.Decode(&userGet)
	return userGet.UserId, err
}

func testUserLogin(t *testing.T, client http.Client, lb *internalhttp.LoginBody, url string) (string, error) { //nolint
	data, err := json.Marshal(lb)
	require.NoError(t, err)
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	require.NoError(t, err)
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var token internalhttp.UserLoginResponse
	jsDecoder := json.NewDecoder(resp.Body)
	err = jsDecoder.Decode(&token)
	return token.Token, err
}

func testUserGet(t *testing.T, client http.Client, url string) (*internalhttp.User, error) { //nolint
	resp, err := http.Get(url) //nolint
	defer resp.Body.Close()    //nolint
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var user internalhttp.User
	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	err = json.Unmarshal(data, &user)
	fmt.Println(toJSON(user))
	require.NoError(t, err)

	return &user, err

}
