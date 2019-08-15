package wallet

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/crypto"
)

// Special error returned when signature verification fails
var ErrInvalidSignature = errors.New("login.error.invalidSignature")

// Returns an hex string representation of a message to be used for login challenge.
// The challenge is prefixed by a human-readable message so that it can be displayed on Metamask.
// The challenge itself is an hex string of 32 random bytes.
func CreateSignInChallenge(i18nMessage string) string {
	// generate array from random 32 bytes
	challenge := make([]byte, 32)
	rand.Read(challenge)
	challengeHex := "0x" + hex.EncodeToString(challenge)

	result := append([]byte(i18nMessage), []byte(challengeHex)...)

	return "0x" + hex.EncodeToString(result)
}

// Verifies if the given signature matches the provided challenge and returns the address
// of the wallet that made the signature.
func VerifySignInChallenge(challengeHex, signatureHex string) (addressHex string, err error) {
	minimumChallengeSize := 2
	if len(challengeHex) < minimumChallengeSize {
		// challenge is stored in memory, so it can happed that it will be empty after server restart
		err = fmt.Errorf("wrong challenge size: Expected more than %d characters, but %d found", minimumChallengeSize, len(challengeHex))
		return
	}
	if len(signatureHex) != 132 {
		err = fmt.Errorf("wrong signature size: Expected 132 characters, but %d found", len(signatureHex))
		return
	}

	// for some reason hex decode wants an hex string with no hex prefix("0x"), hence [2:]
	signature, err := hex.DecodeString(signatureHex[2:])
	if err != nil {
		return
	}

	// need to subtract 27 from V, reason here: https://github.com/ethereum/wiki/wiki/JavaScript-API#returns-45
	signature[64] -= 27

	// get the hash of the challenge
	challenge, err := hex.DecodeString(challengeHex[2:])
	if err != nil {
		return
	}

	// pre-pend the eth_sign RPC prefix https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_sign
	challenge = append([]byte("\x19Ethereum Signed Message:\n"+strconv.Itoa(len(challenge))), challenge...)
	challengeHash := crypto.Keccak256(challenge)

	// get the public key from the signature and challenge hash
	pubKeyECDSA, err := crypto.SigToPub(challengeHash, signature)

	// if any of the decodes raised errors, return it
	if err != nil {
		return
	}
	pubKey := crypto.FromECDSAPub(pubKeyECDSA)
	addressHex = crypto.PubkeyToAddress(*pubKeyECDSA).String()

	// build the byte array with R and S
	signatureRS := signature[0:64]

	// verify the signature and return result.
	// TODO check if this verify signature is really needed. Since we are obtaining the public key from the signature, it will always verify to true.
	if !crypto.VerifySignature(pubKey, challengeHash, signatureRS) {
		err = ErrInvalidSignature
	}
	return
}
