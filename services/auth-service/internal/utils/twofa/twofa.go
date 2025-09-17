package twofa

import (
	"errors"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

const (
	DefaultPeriod     = 30
	DefaultSkew       = 0
	DefaultDigits     = otp.DigitsSix
	DefaultSecretSize = 32
)

var (
	ErrMissingCode = errors.New("missing 2FA code")
	ErrInvalidCode = errors.New("invalid 2FA code")
)

type TwoFAUtil struct {
	Issuer string
}

func NewTwoFAUtil(issuer string) *TwoFAUtil {
	return &TwoFAUtil{Issuer: issuer}
}

type SetupResult struct {
	Secret string
	OTPURL string
}

// GenerateSecret generates a new TOTP secret and returns the setup result containing the secret and OTP URL.
func (t *TwoFAUtil) GenerateSecret(accountName string) (*SetupResult, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      t.Issuer,
		AccountName: accountName,
		Period:      DefaultPeriod,
		SecretSize:  DefaultSecretSize,
		Algorithm:   otp.AlgorithmSHA1,
	})
	if err != nil {
		return nil, err
	}
	return &SetupResult{
		Secret: key.Secret(),
		OTPURL: key.URL(),
	}, nil
}

// VerifyCode verifies a TOTP code against the user's secret.
func (t *TwoFAUtil) VerifyCode(secret, code string) error {
	if code == "" {
		return ErrMissingCode
	}

	opts := totp.ValidateOpts{
		Period:    DefaultPeriod,
		Skew:      DefaultSkew,
		Digits:    DefaultDigits,
		Algorithm: otp.AlgorithmSHA1,
	}

	valid, err := totp.ValidateCustom(code, secret, time.Now().UTC(), opts)
	if err != nil || !valid {
		return ErrInvalidCode
	}

	return nil
}
