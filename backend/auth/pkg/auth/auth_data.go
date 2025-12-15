package auth

import "github.com/google/uuid"

type AuthData struct {
	SessionID uuid.UUID
	AccountID uuid.UUID
}

func ClaimsToAuthData(claims *SessionAccessClaims) (*AuthData, error) {
	sessionID, err := claims.GetSessionID()
	if err != nil {
		return nil, err
	}

	accountID, err := claims.GetAccountID()
	if err != nil {
		return nil, err
	}

	data := &AuthData{
		SessionID: sessionID,
		AccountID: accountID,
	}

	return data, nil
}
