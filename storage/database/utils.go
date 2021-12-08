package database

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/ProxeusApp/proxeus-core/storage"

	"github.com/asdine/storm/q"

	"github.com/ProxeusApp/proxeus-core/sys/model"
)

type (
	simpleQuery struct {
		metaOnly bool
		rawIndex int
		index    int
		limit    int
		exclude  []interface{}
		include  []interface{}
	}
)

func containsCaseInsensitiveReg(contains string) string {
	contains = strings.TrimSpace(contains)
	if contains == "" {
		return ""
	}
	return "(?i)" + regexp.QuoteMeta(contains)
}

func makeSimpleQuery(o storage.Options) *simpleQuery {
	const Limit = 1000
	sq := &simpleQuery{metaOnly: o.MetaOnly, index: o.Index, limit: o.Limit}

	if l := len(o.Exclude); l > 0 {
		sq.exclude = make([]interface{}, l)
		i := 0
		for k := range o.Exclude {
			sq.exclude[i] = k
			i++
		}
	}

	if l := len(o.Include); l > 0 {
		sq.include = make([]interface{}, l)
		i := 0
		for k := range o.Include {
			sq.include[i] = k
			i++
		}
	}

	if sq.limit == 0 || sq.limit > Limit {
		sq.limit = Limit
	}
	sq.rawIndex = sq.index
	sq.index = sq.index * sq.limit
	return sq
}

func defaultMatcher(auth model.Auth, contains string, params *simpleQuery, includeReadGranted bool) []q.Matcher {
	matchers := commonMatcher(auth, contains, params)
	matchers = append(matchers, q.And(IsReadGrantedFor(auth, includeReadGranted)))
	return matchers
}

func publishedMatcher(auth model.Auth, contains string, params *simpleQuery) []q.Matcher {
	matchers := commonMatcher(auth, contains, params)
	var m q.Matcher
	if auth == nil {
		m = q.Eq("Published", true)
	} else {
		m = q.Or(
			q.Eq("Owner", auth.UserID()),
			q.Eq("Published", true),
		)
	}
	matchers = append(matchers, q.And(m))
	return matchers
}

func commonMatcher(auth model.Auth, contains string, params *simpleQuery) []q.Matcher {
	contains = containsCaseInsensitiveReg(contains)
	matchers := make([]q.Matcher, 0)
	if contains != "" {
		matchers = append(matchers,
			q.And(
				q.Or(
					q.Re("Name", contains),
					q.Re("Detail", contains),
				),
			),
		)
	}
	if params != nil {
		if len(params.exclude) > 0 {
			matchers = append(matchers,
				q.And(
					q.Not(q.In("ID", params.exclude)),
				),
			)
		}
		if len(params.include) > 0 {
			matchers = append(matchers,
				q.And(
					q.In("ID", params.include),
				),
			)
		}
	}
	return matchers
}

func EncryptWithAES(secret, stringToEncrypt string) (string, error) {

	key := []byte(secret)
	plaintext := []byte(stringToEncrypt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext), nil
}

func DecryptWithAES(secret, encryptedString string) (string, error) {

	key := []byte(secret)
	enc, _ := hex.DecodeString(encryptedString)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := aesGCM.NonceSize()
	if len(enc) < nonceSize {
		return "", errors.New("decrypted key is corrupted")
	}
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	pt := string(plaintext)
	return pt, nil
}
