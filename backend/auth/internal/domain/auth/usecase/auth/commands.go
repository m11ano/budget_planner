package auth

import (
	"context"
	"errors"
	"net"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	appErrors "github.com/m11ano/budget_planner/backend/auth/internal/app/errors"
	"github.com/m11ano/budget_planner/backend/auth/internal/common/uctypes"
	"github.com/m11ano/budget_planner/backend/auth/internal/domain/auth/entity"
	"github.com/m11ano/budget_planner/backend/auth/internal/domain/auth/usecase"
	"github.com/m11ano/budget_planner/backend/auth/pkg/auth"
	"github.com/m11ano/budget_planner/backend/auth/pkg/emailnormalize"
)

func (uc *UsecaseImpl) LoginByEmail(
	ctx context.Context,
	email string,
	password string,
	ip net.IP,
) (*usecase.AuthSessionDTO, error) {
	const op = "LoginByEmail"

	res, err := emailnormalize.Normalize(strings.TrimSpace(email))
	if err != nil {
		return nil, appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", uc.pkg, op)
	}

	email = res.NormalizedAddress

	out := &usecase.AuthSessionDTO{}

	timeNow := time.Now()

	err = uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		var err error

		out.AccountDTO, err = uc.accountUC.FindOneByEmail(ctx, email, &uctypes.QueryGetOneParams{
			ForUpdate: true,
		})
		if err != nil {
			if errors.Is(err, appErrors.ErrNotFound) {
				return appErrors.ErrUnauthorized
			}
			return err
		}

		if !out.AccountDTO.Account.VerifyPassword(password) {
			return usecase.ErrPasswordIncorrect
		}

		if !out.AccountDTO.Account.IsConfirmed {
			return usecase.ErrAccountNotConfirmed
		}

		if out.AccountDTO.Account.IsBlocked {
			return usecase.ErrAccountBlocked
		}

		out.Session = entity.NewSession(
			out.AccountDTO.Account.ID,
			ip,
			timeNow,
			time.Duration(uc.cfg.Auth.RefreshTokenLifetimeHrs)*time.Hour,
		)
		if err := uc.sessionUC.Create(ctx, out.Session); err != nil {
			return err
		}

		out.AccountDTO.Account.SetLastRequestAt(&timeNow)
		out.AccountDTO.Account.SetLastLoginAt(&timeNow)
		out.AccountDTO.Account.SetLastRequestIP(&ip)

		err = uc.accountUC.UpdateAccount(ctx, out.AccountDTO.Account)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	out.RefreshClaims = uc.makeRefreshClaims(out.Session)
	out.AccessClaims = uc.makeAccessClaims(
		out.Session.ID,
		out.AccountDTO.Account.ID,
		nil,
		timeNow,
		time.Duration(uc.cfg.Auth.AccessTokenLifetimeSec)*time.Second,
	)

	return out, nil
}

func (uc *UsecaseImpl) IssueAccessJWT(access *auth.SessionAccessClaims) (string, error) {
	const op = "IssueAccessJWT"

	signed, err := auth.IssueAccessJWT(access, []byte(uc.cfg.Auth.JwtAccessSecret))
	if err != nil {
		return "", appErrors.Chainf(appErrors.ErrInternal.WithWrap(err), "%s.%s", uc.pkg, op)
	}

	return signed, nil
}

func (uc *UsecaseImpl) IssueRefreshJWT(refresh *entity.SessionRefreshClaims) (string, error) {
	const op = "IssueRefreshJWT"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, refresh)

	signed, err := token.SignedString([]byte(uc.cfg.Auth.JwtRefreshSecret))
	if err != nil {
		return "", appErrors.Chainf(appErrors.ErrInternal.WithWrap(err), "%s.%s", uc.pkg, op)
	}

	return signed, nil
}

// func (uc *UsecaseImpl) ParseAccessToken(tokenStr string, validate bool) (*auth.SessionAccessClaims, error) {
// 	const op = "ParseAccessToken"

// 	claims, err := auth.ParseAccessToken(tokenStr, validate, []byte(uc.cfg.Auth.JwtAccessSecret))
// 	if err != nil {
// 		if errors.Is(err, auth.ErrInvalidToken) {
// 			return nil, appErrors.Chainf(usecase.ErrInvalidToken, "%s.%s", uc.pkg, op)
// 		}
// 		return nil, appErrors.Chainf(appErrors.ErrInternal.WithWrap(err), "%s.%s", uc.pkg, op)
// 	}

// 	return claims, nil
// }

func (uc *UsecaseImpl) ParseRefreshToken(tokenStr string, validate bool) (*entity.SessionRefreshClaims, error) {
	const op = "ParseRefreshToken"

	ops := []jwt.ParserOption{
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	}

	if !validate {
		ops = append(ops, jwt.WithoutClaimsValidation())
	}

	token, err := jwt.ParseWithClaims(tokenStr, &entity.SessionRefreshClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(uc.cfg.Auth.JwtRefreshSecret), nil
	}, ops...)
	if err != nil || !token.Valid {
		return nil, appErrors.Chainf(usecase.ErrInvalidToken, "%s.%s", uc.pkg, op)
	}

	claims, ok := token.Claims.(*entity.SessionRefreshClaims)
	if !ok {
		return nil, appErrors.Chainf(usecase.ErrInvalidToken, "%s.%s", uc.pkg, op)
	}

	return claims, nil
}

func (uc *UsecaseImpl) GenerateNewClaims(
	ctx context.Context,
	refreshClaims *entity.SessionRefreshClaims,
	ip net.IP,
) (*usecase.AuthSessionDTO, error) {
	const op = "GenerateNewClaims"

	out := &usecase.AuthSessionDTO{}

	timeNow := time.Now()

	err := uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		var err error

		out.Session, err = uc.sessionUC.FindOneByID(ctx, refreshClaims.SessionID, &uctypes.QueryGetOneParams{
			ForUpdate: true,
		})
		if err != nil {
			if errors.Is(err, appErrors.ErrNotFound) {
				return appErrors.ErrUnauthorized
			}
			return err
		}

		if out.Session.RefreshToken != refreshClaims.RefreshKey {
			return appErrors.ErrUnauthorized
		}

		if out.Session.RefreshVersion != refreshClaims.RefreshVersion {
			return appErrors.ErrUnauthorized
		}

		if out.Session.RefreshTokenExpiresAt.Before(timeNow) {
			return appErrors.ErrUnauthorized
		}

		out.AccountDTO, err = uc.accountUC.FindOneByID(ctx, out.Session.AccountID, &uctypes.QueryGetOneParams{
			ForUpdate: true,
		})
		if err != nil {
			if errors.Is(err, appErrors.ErrNotFound) {
				return appErrors.ErrUnauthorized
			}
			return err
		}

		if !out.AccountDTO.Account.IsConfirmed {
			return usecase.ErrAccountNotConfirmed
		}

		if out.AccountDTO.Account.IsBlocked {
			return usecase.ErrAccountBlocked
		}

		out.Session.GenerateNewRefresh(timeNow, time.Duration(uc.cfg.Auth.RefreshTokenLifetimeHrs)*time.Hour, ip)

		err = uc.sessionUC.Update(ctx, out.Session)
		if err != nil {
			return err
		}

		out.AccountDTO.Account.SetLastRequestAt(&timeNow)
		out.AccountDTO.Account.SetLastLoginAt(&timeNow)
		out.AccountDTO.Account.SetLastRequestIP(&ip)

		err = uc.accountUC.UpdateAccount(ctx, out.AccountDTO.Account)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	out.RefreshClaims = uc.makeRefreshClaims(out.Session)
	out.AccessClaims = uc.makeAccessClaims(
		out.Session.ID,
		out.AccountDTO.Account.ID,
		nil,
		timeNow,
		time.Duration(uc.cfg.Auth.AccessTokenLifetimeSec)*time.Second,
	)

	return out, nil
}
