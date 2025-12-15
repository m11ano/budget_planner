package entity

import (
	"net"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	appErrors "github.com/m11ano/budget_planner/backend/auth/internal/app/errors"
	"github.com/m11ano/budget_planner/backend/auth/pkg/emailnormalize"
	"github.com/samber/lo"
	"golang.org/x/crypto/bcrypt"
)

var ErrAccountInvalidEmail = appErrors.ErrBadRequest.Extend("invalid email").WithTextCode("INVALID_EMAIL")

var ErrAccountInvalidPassword = appErrors.ErrBadRequest.Extend("invalid password").WithTextCode("INVALID_PASSWORD")

var ErrAccountProfileInvalidName = appErrors.ErrBadRequest.Extend("invalid name").WithTextCode("INVALID_NAME")

var ErrAccountProfileInvalidSurname = appErrors.ErrBadRequest.Extend("invalid surname").WithTextCode("INVALID_SURNAME")

type Account struct {
	ID             uuid.UUID
	Email          string
	PasswordHash   string
	IsConfirmed    bool
	IsBlocked      bool
	LastLoginAt    *time.Time
	LastRequestAt  *time.Time
	LastRequestIP  *net.IP
	ProfileName    string
	ProfileSurname string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (item *Account) Version() int64 {
	return item.UpdatedAt.UnixMicro()
}

func (item *Account) SetPassword(password string, skipStrongCheck bool) error {
	if !skipStrongCheck {

		if len(password) < 8 {
			return ErrAccountInvalidPassword
		}

		var hasLetter, hasDigit, hasSpecial bool
		for _, r := range password {
			switch {
			case unicode.IsLetter(r):
				hasLetter = true
			case unicode.IsDigit(r):
				hasDigit = true
			default:
				hasSpecial = true
			}
		}

		count := 0
		if hasLetter {
			count++
		}
		if hasDigit {
			count++
		}
		if hasSpecial {
			count++
		}

		if count < 2 {
			return ErrAccountInvalidPassword
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	item.PasswordHash = string(hash)

	return nil
}

func (item *Account) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(item.PasswordHash), []byte(password))
	return err == nil
}

func (item *Account) SetEmail(email string) error {
	email = strings.TrimSpace(email)

	err := validate.Var(email, "required,email")
	if err != nil {
		return ErrAccountInvalidEmail
	}

	res, err := emailnormalize.Normalize(email)
	if err != nil {
		return ErrAccountInvalidEmail
	}

	item.Email = res.NormalizedAddress

	return nil
}

func (item *Account) SetLastLoginAt(value *time.Time) {
	if value != nil {
		value = lo.ToPtr(value.Truncate(time.Microsecond))
	}

	item.LastLoginAt = value
}

func (item *Account) SetLastRequestAt(value *time.Time) {
	if value != nil {
		value = lo.ToPtr(value.Truncate(time.Microsecond))
	}

	item.LastRequestAt = value
}

func (item *Account) SetLastRequestIP(value *net.IP) {
	item.LastRequestIP = value
}

func (item *Account) SetProfileName(value string) error {
	value = strings.TrimSpace(value)

	if value == "" {
		return ErrAccountProfileInvalidName
	}

	item.ProfileName = value

	return nil
}

func (item *Account) SetProfileSurname(value string) error {
	value = strings.TrimSpace(value)

	if value == "" {
		return ErrAccountProfileInvalidSurname
	}

	item.ProfileSurname = value

	return nil
}

func NewAccount(
	email string,
	password string,
	skipPasswordCheck bool,
	isConfirmed bool,
) (*Account, error) {
	timeNow := time.Now().Truncate(time.Microsecond)

	account := &Account{
		ID:          uuid.New(),
		IsConfirmed: isConfirmed,
		CreatedAt:   timeNow,
		UpdatedAt:   timeNow,
	}

	err := account.SetEmail(email)
	if err != nil {
		return nil, err
	}

	err = account.SetPassword(password, skipPasswordCheck)
	if err != nil {
		return nil, err
	}

	return account, nil
}
