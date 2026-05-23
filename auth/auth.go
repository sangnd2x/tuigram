package auth

import (
	"context"
	"fmt"

	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
	"golang.org/x/term"
	"os"
)

type Authenticator struct {
	phone string
}

func New(phone string) *Authenticator {
	return &Authenticator{phone: phone}
}

func PromptPhone() (string, error) {
	fmt.Print("Enter your phone number (with country code, e.g. +1234567890): ")
	var phone string
	_, err := fmt.Scan(&phone)
	return phone, err
}

func (a *Authenticator) Phone(_ context.Context) (string, error) {
	return a.phone, nil
}

func (a *Authenticator) Password(_ context.Context) (string, error) {
	fmt.Fprint(os.Stderr, "Enter 2FA password: ")
	pwd, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Fprintln(os.Stderr)
	if err != nil {
		return "", err
	}
	return string(pwd), nil
}

func (a *Authenticator) Code(_ context.Context, _ *tg.AuthSentCode) (string, error) {
	fmt.Print("Enter OTP code: ")
	var code string
	_, err := fmt.Scan(&code)
	return code, err
}

func (a *Authenticator) AcceptTermsOfService(_ context.Context, _ tg.HelpTermsOfService) error {
	return nil
}

func (a *Authenticator) SignUp(_ context.Context) (auth.UserInfo, error) {
	return auth.UserInfo{}, fmt.Errorf("sign-up not supported")
}
