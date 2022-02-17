package model

import (
	"errors"
	"strings"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

const ApiKeyLength = 40 //example: f235122f9a1e4884123456788a2126f8dd76996b

type ApiKey struct {
	Name string
	Key  string
}

func NewApiKey(name, uid string) (*ApiKey, error) {
	if len(uid) < 8 {
		return nil, errors.New("id must be set")
	}
	u := uuid.NewV4()
	k := strings.Replace(u.String(), "-", "", -1)
	//mark the center of the api key with the user prefix id
	return &ApiKey{Name: name, Key: k[0:16] + uid[0:8] + k[16:]}, nil
}

//HideKey hashes the key
func (me *ApiKey) HideKey() {
	me.Key, _ = hashString(me.Key)
}

func MatchesApiKey(hiddenApiKey, apiKey string) bool {
	return checkStringHash(apiKey, hiddenApiKey)
}

func hashString(str string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(str), 14)
	return string(bytes), err
}

func checkStringHash(str, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(str))
	return err == nil
}
