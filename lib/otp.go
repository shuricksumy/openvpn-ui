package lib

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/shuricksumy/openvpn-ui/models"
)

func Get2FA(clientID string) (*otp.Key, *string, bool, error) {

	client, err_get_client := models.GetClientDetailsById(clientID)
	var userOTP *otp.Key
	var isOTPNew bool

	if err_get_client != nil {
		return nil, nil, true, err_get_client
	}

	if client.OTPIsEnabled || client.StaticPassIsUsed {
		//USE IT
		isOTPNew = false
		var hashInBytes []byte

		if client.OTPKey == nil {
			hashInBytes, _ = GenerateHashInByte()
		} else {
			hashInBytes, _ = hex.DecodeString(NilStringToString(client.OTPKey))
		}

		userOTP, _ = totp.Generate(totp.GenerateOpts{
			Issuer:      "Example.com",
			AccountName: NilStringToString(client.OTPUserName),
			Algorithm:   otp.AlgorithmSHA256,
			Secret:      hashInBytes,
		})
		return userOTP, client.OTPKey, isOTPNew, nil
	} else {
		//Create new
		isOTPNew = true
		hashInBytes, key := GenerateHashInByte()

		userOTP, _ = totp.Generate(totp.GenerateOpts{
			Issuer:      "Example.com",
			AccountName: client.ClientName,
			Algorithm:   otp.AlgorithmSHA256,
			Secret:      hashInBytes,
		})
		return userOTP, key, isOTPNew, nil
	}
}

func GenerateHashInByte() ([]byte, *string) {
	randomBytes := make([]byte, 32)
	rand.Read(randomBytes)
	hash := sha256.New()
	hash.Write(randomBytes)
	hashInBytes := hash.Sum(nil) // KEY byte

	strKey := hex.EncodeToString(hashInBytes)
	hexKey := StringToNilString(strKey)

	return hashInBytes, hexKey
}
