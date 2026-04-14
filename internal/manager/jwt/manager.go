package jwtManager

import (
	"fmt"
	"strings"
	"time"

	"github.com/balobas/auth_service/internal/entity"
	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
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
	tokenFieldUserUid    = "user_uid"
	tokenFieldDeviceUid  = "device_uid"
	tokenFieldEmail      = "email"
	tokenFieldRoles      = "roles"
	tokenFieldSessionUid = "session_uid"
	tokenFieldExpiredAt  = "expired_at"
	tokenFieldIssuedAt   = "issued_at"
)

const rolesSeparator = ","

func (p *JwtManager) NewToken(info entity.TokenInfo, ttl time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims[tokenFieldUserUid] = info.UserUid.String()
	claims[tokenFieldDeviceUid] = info.DeviceUid.String()
	claims[tokenFieldEmail] = info.Email
	claims[tokenFieldRoles] = strings.Join(info.Roles, rolesSeparator)
	claims[tokenFieldSessionUid] = info.SessionUid.String()
	claims[tokenFieldExpiredAt] = time.Now().Add(ttl).Unix()
	claims[tokenFieldIssuedAt] = info.IssuedAt

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

func (p *JwtManager) ParseToken(tokenStr string) (t entity.TokenInfo, err error) {
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

	defer func() {
		if err != nil {
			err = errors.Wrap(serviceErrors.ErrInvalidToken, err.Error())
		}
	}()

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
		return entity.TokenInfo{}, errors.New("empty user uid in token")
	}

	tokenInfo.UserUid, err = uuid.FromString(userUid.(string))
	if err != nil {
		return entity.TokenInfo{}, errors.New("invalid user uid in token")
	}

	deviceUid, ok := claims[tokenFieldDeviceUid]
	if !ok {
		return entity.TokenInfo{}, errors.New("empty device uid in token")
	}

	tokenInfo.DeviceUid, err = uuid.FromString(deviceUid.(string))
	if err != nil {
		return entity.TokenInfo{}, errors.New("invalid device uid in token")
	}

	userEmail, ok := claims[tokenFieldEmail]
	if !ok {
		return entity.TokenInfo{}, errors.New("empty userEmail in token")
	}
	tokenInfo.Email, ok = userEmail.(string)
	if !ok {
		return entity.TokenInfo{}, errors.New("invalid token user email")
	}

	roles, ok := claims[tokenFieldRoles]
	if !ok {
		return entity.TokenInfo{}, errors.New("empty roles in token")
	}

	rolesStr, ok := roles.(string)
	if !ok {
		return entity.TokenInfo{}, errors.New("invalid roles")
	}
	if len(rolesStr) != 0 {
		tokenInfo.Roles = strings.Split(rolesStr, rolesSeparator)
	}

	sessionUid, ok := claims[tokenFieldSessionUid]
	if !ok {
		return entity.TokenInfo{}, errors.New("invalid session uid")
	}

	tokenInfo.SessionUid = uuid.FromStringOrNil(sessionUid.(string))

	expiredAt, ok := claims[tokenFieldExpiredAt].(float64)
	if !ok {
		return entity.TokenInfo{}, errors.New("empty expired_at in token")
	}

	tokenInfo.ExpiredAt = int64(expiredAt)

	issuedAt, ok := claims[tokenFieldIssuedAt].(float64)
	if !ok {
		return entity.TokenInfo{}, errors.New("empty issued_at in token")
	}

	tokenInfo.IssuedAt = int64(issuedAt)

	return tokenInfo, nil
}
