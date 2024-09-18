package jwtManager

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

type JwtManager struct {
	keysProvider KeysProvider
}

type KeysProvider interface {
	GetPrivateKey() ([]byte, error)
}

func New(keysProvider KeysProvider) *JwtManager {
	return &JwtManager{
		keysProvider: keysProvider,
	}
}

const (
	tokenFieldUserUid     = "user_uid"
	tokenFieldEmail       = "email"
	tokenFieldPermissions = "permissions"
	tokenFieldRole        = "role"
	tokenFieldSessionUid  = "session_uid"
	tokenFieldExpiredAt   = "expired_at"
)

const permissionsSeparator = ","

func (p *JwtManager) NewToken(info entity.TokenInfo, ttl time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims[tokenFieldUserUid] = info.UserUid.String()
	claims[tokenFieldEmail] = info.Email
	claims[tokenFieldPermissions] = strings.Join(info.Permissions, permissionsSeparator)
	claims[tokenFieldRole] = info.Role
	claims[tokenFieldSessionUid] = info.SessionUid.String()
	claims[tokenFieldExpiredAt] = strconv.FormatInt(int64(ttl), 10)

	pk, err := p.keysProvider.GetPrivateKey()
	if err != nil {
		return "", errors.WithStack(err)
	}

	signedToken, err := token.SignedString(pk)
	if err != nil {
		return "", errors.Wrapf(err, "cant sign token")
	}

	return signedToken, nil
}

func (p *JwtManager) ParseToken(tokenStr string) (entity.TokenInfo, error) {
	tk, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		pk, err := p.keysProvider.GetPrivateKey()
		if err != nil {
			return nil, errors.WithStack(err)
		}

		return pk, nil
	})
	if err != nil {
		return entity.TokenInfo{}, errors.Wrap(err, "cant parse jwt token")
	}

	if !tk.Valid {
		return entity.TokenInfo{}, errors.New("token is invalid")
	}

	claims, ok := tk.Claims.(jwt.MapClaims)
	if !ok {
		return entity.TokenInfo{}, errors.New("invalid claims")
	}

	tokenInfo := entity.TokenInfo{}

	userUid, ok := claims[tokenFieldUserUid]
	if !ok {
		return entity.TokenInfo{}, errors.New("empty uid in token")
	}

	tokenInfo.UserUid = uuid.FromStringOrNil(userUid.(string))

	userEmail, ok := claims[tokenFieldEmail]
	if !ok {
		return entity.TokenInfo{}, errors.New("empty userEmail in token")
	}
	tokenInfo.Email, ok = userEmail.(string)
	if !ok {
		return entity.TokenInfo{}, errors.New("invalid token user email")
	}

	perms, ok := claims[tokenFieldPermissions]
	if !ok {
		return entity.TokenInfo{}, errors.New("empty permissions in token")
	}
	permsStr, ok := perms.(string)
	if !ok {
		return entity.TokenInfo{}, errors.New("invalid permissions")
	}
	tokenInfo.Permissions = strings.Split(permsStr, permissionsSeparator)

	role, ok := claims[tokenFieldRole]
	if !ok {
		return entity.TokenInfo{}, errors.New("invalid role")
	}
	tokenInfo.Role, ok = role.(string)
	if !ok {
		return entity.TokenInfo{}, errors.New("invalid role")
	}

	sessionUid, ok := claims[tokenFieldSessionUid]
	if !ok {
		return entity.TokenInfo{}, errors.New("invalid session uid")
	}

	tokenInfo.SessionUid = uuid.FromStringOrNil(sessionUid.(string))

	expiredAt, ok := claims[tokenFieldExpiredAt]
	if !ok {
		return entity.TokenInfo{}, errors.New("empty expired_at in token")
	}

	exp, ok := expiredAt.(string)
	if !ok {
		return entity.TokenInfo{}, errors.New("invalid expired_at to string")
	}
	expInt, err := strconv.ParseInt(exp, 10, 0)
	if err != nil {
		return entity.TokenInfo{}, errors.Errorf("invalid expired_at: %v", err)
	}

	tokenInfo.ExpiredAt = int64(expInt)

	return tokenInfo, nil
}
