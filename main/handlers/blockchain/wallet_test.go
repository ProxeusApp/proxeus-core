package blockchain

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifySignInChallenge(t *testing.T) {

	t.Run("when empty strings are given", func(t *testing.T) {
		_, err := VerifySignInChallenge("", "")

		assert.NotNil(t, err, "should return error")
	})

	t.Run("when wrong signatureHex given", func(t *testing.T) {
		_, err := VerifySignInChallenge("0x000", "not valid signature")

		assert.NotNil(t, err, "should return error")
	})

	t.Run("when valid challenge and signature are given", func(t *testing.T) {
		challenge, err := VerifySignInChallenge(
			"0x4163636f756e74207369676e206d657373616765307834366434356564653539663133346265393863633632396461616537633935323361636665306637323266633738376565343130633434393838653339663833",
			"0x0f76fb55e04c2ea5be34bd7dcf9c96244d2bae795bfa2043e031e6fc37802e2f6678a661708036c3bedaea1da71989eccec94326390f1d66847c78ace67714be1c")

		assert.Nil(t, err, "should not return error")
		assert.Equal(t, "0x8d03c4AB63911d57EDDdfB98268D220f3080F51D", challenge, "should return a correctly generated challenge")
	})
}

func TestCreateSignInChallenge(t *testing.T) {

	generatedChallenges := make(map[string]bool, 10)
	for i := 0; i < 10; i++ {
		challenge := CreateSignInChallenge("test")

		assert.False(t, generatedChallenges[challenge], "should never generate the same challenge")
		assert.True(t, strings.HasPrefix(challenge, "0x"))

		generatedChallenges[challenge] = true
	}
}
