package jwtManager

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

type JwtManager struct {
	keysProvider KeysProvider
}

type KeysProvider interface {
	Key() []byte
}

func New(keysProvider KeysProvider) *JwtManager {
	return &JwtManager{
		keysProvider: keysProvider,
	}
}

type TokenInfo struct {
	UserUid     uuid.UUID
	Email       string
	Permissions []string
	Role        string
	ExpiredAt   int64
}

const (
	tokenFieldUserUid     = "user_uid"
	tokenFieldEmail       = "email"
	tokenFieldPermissions = "permissions"
	tokenFieldRole        = "role"
	tokenFieldExpiredAt   = "expired_at"
)

const permissionsSeparator = ","

func (p *JwtManager) NewToken(info TokenInfo, ttl time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims[tokenFieldUserUid] = info.UserUid
	claims[tokenFieldEmail] = info.Email
	claims[tokenFieldPermissions] = strings.Join(info.Permissions, permissionsSeparator)
	claims[tokenFieldRole] = info.Role
	claims[tokenFieldExpiredAt] = info.ExpiredAt

	signedToken, err := token.SignedString(p.keysProvider.Key())
	if err != nil {
		return "", errors.Wrapf(err, "cant sign token")
	}

	return signedToken, nil
}

func (p *JwtManager) ParseToken(tokenStr string) (TokenInfo, error) {
	tk, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return p.keysProvider.Key(), nil
	})
	if err != nil {
		return TokenInfo{}, errors.Wrap(err, "cant parse jwt token")
	}

	if !tk.Valid {
		return TokenInfo{}, errors.New("invalid token")
	}

	claims, ok := tk.Claims.(jwt.MapClaims)
	if !ok {
		return TokenInfo{}, errors.New("invalid claims")
	}

	tokenInfo := TokenInfo{}

	userUid, ok := claims[tokenFieldUserUid]
	if !ok {
		return TokenInfo{}, errors.New("empty uid in token")
	}
	tokenInfo.UserUid, ok = userUid.(uuid.UUID)
	if !ok {
		return TokenInfo{}, errors.New("invalid token")
	}

	userEmail, ok := claims[tokenFieldEmail]
	if !ok {
		return TokenInfo{}, errors.New("empty userEmail in token")
	}
	tokenInfo.Email, ok = userEmail.(string)
	if !ok {
		return TokenInfo{}, errors.New("invalid token")
	}

	perms, ok := claims[tokenFieldPermissions]
	if !ok {
		return TokenInfo{}, errors.New("empty permissions in token")
	}
	permsStr, ok := perms.(string)
	if !ok {
		return TokenInfo{}, errors.New("invalid permissions")
	}
	tokenInfo.Permissions = strings.Split(permsStr, permissionsSeparator)

	role, ok := claims[tokenFieldRole]
	if !ok {
		return TokenInfo{}, errors.New("invalid role")
	}
	tokenInfo.Role, ok = role.(string)
	if !ok {
		return TokenInfo{}, errors.New("invalid role")
	}

	expiredAt, ok := claims[tokenFieldExpiredAt]
	if !ok {
		return TokenInfo{}, errors.New("empty expired_at in token")
	}
	tokenInfo.ExpiredAt, ok = expiredAt.(int64)
	if !ok {
		return TokenInfo{}, errors.New("invalid expired_at")
	}

	return tokenInfo, nil
}
