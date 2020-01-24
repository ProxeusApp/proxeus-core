package model

import (
	"errors"
	"strings"

	uuid "github.com/satori/go.uuid"
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

//HideKey turns f235122f9a1e4884123456788a2126f8dd76996b into f235...996b
func (me *ApiKey) HideKey() {
	if me.IsNew() {
		me.Key = me.Key[0:4] + "..." + me.Key[len(me.Key)-4:]
	}
}

//IsNew returns true on keys like f235122f9a1e4884123456788a2126f8dd76996b
//but hidden keys like f235...996b are hidden and therefore not new anymore
func (me *ApiKey) IsNew() bool {
	//if not hidden then new
	return len(me.Key) > 11
}

func MatchesApiKey(hiddenApiKey, apiKey string) bool {
	return strings.HasPrefix(apiKey, hiddenApiKey[0:4]) && strings.HasSuffix(apiKey, hiddenApiKey[len(hiddenApiKey)-4:])
}
