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
	tokenFieldType       = "type"
	tokenFieldEmail      = "email"
	tokenFieldRoles      = "roles"
	tokenFieldSessionUid = "session_uid"
	tokenFieldExpiredAt  = "expired_at"
	tokenFieldIssuedAt   = "issued_at"

	tokenFieldServiceUid  = "service_uid"
	tokenFieldServiceName = "service_name"
	tokenFieldDomain      = "domain"
)

const rolesSeparator = ","

func (p *JwtManager) NewUserTokens(info entity.TokenInfo, accessTtl, refreshTtl time.Duration) (access string, refresh string, err error) {
	return p.newTokens(info, fillUserTokenClaims, accessTtl, refreshTtl)
}

func (p *JwtManager) NewSystemTokens(info entity.TokenInfo, accessTtl, refreshTtl time.Duration) (access string, refresh string, err error) {
	return p.newTokens(info, fillSystemTokenClaims, accessTtl, refreshTtl)
}

func (p *JwtManager) newTokens(info entity.TokenInfo, fillClaims TokenClaimsFiller, accessTtl, refreshTtl time.Duration) (access string, refresh string, err error) {
	access, err = p.newToken(info, accessTtl, fillClaims)
	if err != nil {
		return
	}
	refresh, err = p.newToken(info, refreshTtl, fillClaims)
	if err != nil {
		return
	}
	return
}

func (p *JwtManager) newToken(info entity.TokenInfo, ttl time.Duration, fillClaims TokenClaimsFiller) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	fillClaims(claims, info, ttl)

	pk, err := p.keysProvider.GetPrivateKey()
	if err != nil {
		return "", err
	}

	return token.SignedString(pk)
}

type TokenClaimsFiller func(jwt.MapClaims, entity.TokenInfo, time.Duration)

func fillUserTokenClaims(claims jwt.MapClaims, info entity.TokenInfo, ttl time.Duration) {
	claims[tokenFieldUserUid] = info.UserUid.String()
	claims[tokenFieldDeviceUid] = info.DeviceUid.String()
	claims[tokenFieldEmail] = info.Email
	claims[tokenFieldRoles] = strings.Join(info.Roles, rolesSeparator)
	claims[tokenFieldSessionUid] = info.SessionUid.String()
	claims[tokenFieldExpiredAt] = time.Now().Add(ttl).Unix()
	claims[tokenFieldIssuedAt] = info.IssuedAt
	claims[tokenFieldType] = string(info.Type)
}

func fillSystemTokenClaims(claims jwt.MapClaims, info entity.TokenInfo, ttl time.Duration) {
	claims[tokenFieldServiceUid] = info.ServiceUid.String()
	claims[tokenFieldServiceName] = info.ServiceName
	claims[tokenFieldDeviceUid] = info.DeviceUid
	claims[tokenFieldDomain] = info.Domain
	claims[tokenFieldRoles] = strings.Join(info.Roles, rolesSeparator)
	claims[tokenFieldSessionUid] = info.SessionUid.String()
	claims[tokenFieldExpiredAt] = time.Now().Add(ttl).Unix()
	claims[tokenFieldIssuedAt] = info.IssuedAt
	claims[tokenFieldType] = string(info.Type)
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

	tokenTypeVal, ok := claims[tokenFieldType]
	if !ok {
		return entity.TokenInfo{}, errors.New("invalid token type")
	}
	tokenType := entity.TokenType(tokenTypeVal.(string))
	if err := tokenType.Validate(); err != nil {
		return entity.TokenInfo{}, err
	}

	switch tokenType {
	case entity.TokenTypeUser:
		return parseUserTokenClaims(claims)
	case entity.TokenTypeSystem:
		return parseSystemTokenClaims(claims)
	default:
		return entity.TokenInfo{}, errors.New("unknow token type")
	}
}

func parseUserTokenClaims(claims jwt.MapClaims) (entity.TokenInfo, error) {
	tokenInfo := entity.TokenInfo{
		Type: entity.TokenTypeUser,
	}

	var err error

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

func parseSystemTokenClaims(claims jwt.MapClaims) (entity.TokenInfo, error) {
	tokenInfo := entity.TokenInfo{
		Type: entity.TokenTypeSystem,
	}

	var err error

	serviceUid, ok := claims[tokenFieldServiceUid]
	if !ok {
		return entity.TokenInfo{}, errors.New("empty service uid in token")
	}
	tokenInfo.ServiceUid, err = uuid.FromString(serviceUid.(string))
	if err != nil {
		return entity.TokenInfo{}, errors.New("invalid service uid in token")
	}

	serviceName, ok := claims[tokenFieldServiceName]
	if !ok {
		return entity.TokenInfo{}, errors.New("empty service name in token")
	}
	tokenInfo.ServiceName, ok = serviceName.(string)
	if !ok {
		return entity.TokenInfo{}, errors.New("invalid token service name")
	}

	deviceUid, ok := claims[tokenFieldDeviceUid]
	if !ok {
		return entity.TokenInfo{}, errors.New("empty device uid in token")
	}
	tokenInfo.DeviceUid, err = uuid.FromString(deviceUid.(string))
	if err != nil {
		return entity.TokenInfo{}, errors.New("invalid device uid in token")
	}

	domain, ok := claims[tokenFieldDomain]
	if !ok {
		return entity.TokenInfo{}, errors.New("empty domain in token")
	}
	tokenInfo.Domain, ok = domain.(string)
	if !ok {
		return entity.TokenInfo{}, errors.New("invalid token domain")
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
