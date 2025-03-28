package authentication

// Signs and verifies messages that can be sent to vehicles asynchronously.

import (
	"encoding/base64"

	"github.com/golang-jwt/jwt/v5"
	"github.com/teslamotors/vehicle-command/internal/schnorr"
)

const TeslaSchnorrSHA256 = "Tesla.SS256"

// SigningMethodSchnorrP256 implements jwt.SigningMethod using Schnorr signatures over the NIST
// P-256 curve.
type SigningMethodSchnorrP256 struct{}

var tss256 SigningMethodSchnorrP256 // Singleton used for RegisterSigningMethod

func init() {
	jwt.RegisterSigningMethod(TeslaSchnorrSHA256, func() jwt.SigningMethod { return &tss256 })
}

func (s *SigningMethodSchnorrP256) Verify(signingString string, signature []byte, key interface{}) error {
	pkeyBytes, ok := key.([]byte)
	if !ok {
		return jwt.ErrInvalidKeyType
	}
	return schnorr.Verify(pkeyBytes, []byte(signingString), signature)
}

func (s *SigningMethodSchnorrP256) Sign(signingString string, key interface{}) ([]byte, error) {
	skey, ok := key.(ECDHPrivateKey)
	if !ok {
		return nil, jwt.ErrInvalidKeyType
	}
	return skey.SchnorrSignature([]byte(signingString))
}

func (s *SigningMethodSchnorrP256) Alg() string {
	return TeslaSchnorrSHA256
}

func SignMessage(privateKey ECDHPrivateKey, message jwt.MapClaims, audience string) (string, error) {
	message["iss"] = base64.StdEncoding.EncodeToString(privateKey.PublicBytes())
	message["aud"] = audience
	token := jwt.New(&tss256)
	token.Claims = message
	return token.SignedString(privateKey)
}
