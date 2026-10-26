// Package telegram валидирует initData Telegram Mini App (WebApp) и достаёт данные
// пользователя. Проверка строго серверная по алгоритму Telegram.
package telegram

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

// TODO: временная диагностика hash mismatch — убрать после отладки.
const tgValidateDebug = true

var (
	ErrInvalidInitData = fmt.Errorf("invalid telegram init data")
	ErrInitDataExpired = fmt.Errorf("telegram init data expired")
)

// TelegramUser — распарсенные доверенные данные пользователя из initData.
type TelegramUser struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// ValidateInitData проверяет подпись initData ботовым токеном и свежесть auth_date,
// затем возвращает пользователя. maxAge — максимальный возраст auth_date (anti-replay).
//
// Алгоритм (Telegram WebApp):
//
//	data_check_string = отсортированные "key=value" (кроме hash), join "\n"
//	secret = HMAC_SHA256(key="WebAppData", data=bot_token)
//	hash   = hex(HMAC_SHA256(key=secret, data=data_check_string))
func ValidateInitData(botToken, initData string, maxAge time.Duration) (TelegramUser, error) {
	// Страховка от скрытых символов в конце токена (CRLF/пробел из .env).
	botToken = strings.TrimSpace(botToken)

	values, err := url.ParseQuery(initData)
	if err != nil {
		return TelegramUser{}, fmt.Errorf("%w: parse: %v", ErrInvalidInitData, err)
	}

	hash := values.Get("hash")
	if hash == "" {
		return TelegramUser{}, fmt.Errorf("%w: no hash", ErrInvalidInitData)
	}

	// data_check_string: все пары кроме hash, отсортированы по ключу.
	keys := make([]string, 0, len(values))
	for k := range values {
		if k == "hash" || k == "signature" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	pairs := make([]string, 0, len(keys))
	for _, k := range keys {
		pairs = append(pairs, k+"="+values.Get(k))
	}
	dataCheckString := strings.Join(pairs, "\n")

	secret := hmacSHA256([]byte("WebAppData"), []byte(botToken))
	computed := hex.EncodeToString(hmacSHA256(secret, []byte(dataCheckString)))

	if tgValidateDebug {
		// Фингерпринт токена (без утечки самого токена) — чтобы отличить скрытые
		// символы; а также альтернативный хэш С включением signature, чтобы
		// проверить гипотезу об алгоритме.
		tokenFp := sha256.Sum256([]byte(botToken))

		altKeys := make([]string, 0, len(values))
		for k := range values {
			if k == "hash" {
				continue
			}
			altKeys = append(altKeys, k)
		}
		sort.Strings(altKeys)
		altPairs := make([]string, 0, len(altKeys))
		for _, k := range altKeys {
			altPairs = append(altPairs, k+"="+values.Get(k))
		}
		altComputed := hex.EncodeToString(
			hmacSHA256(secret, []byte(strings.Join(altPairs, "\n"))))

		_, hasSig := values["signature"]
		log.Printf("[tg-debug] tokenLen=%d tokenFp=%x initDataLen=%d keys=%v hasSig=%v recv=%s calc(noSig)=%s calc(withSig)=%s",
			len(botToken), tokenFp[:6], len(initData), keys, hasSig, hash, computed, altComputed)
	}

	if !hmac.Equal([]byte(computed), []byte(hash)) {
		return TelegramUser{}, fmt.Errorf("%w: hash mismatch", ErrInvalidInitData)
	}

	// Свежесть auth_date.
	authDateUnix, err := strconv.ParseInt(values.Get("auth_date"), 10, 64)
	if err != nil {
		return TelegramUser{}, fmt.Errorf("%w: bad auth_date", ErrInvalidInitData)
	}
	if maxAge > 0 && time.Since(time.Unix(authDateUnix, 0)) > maxAge {
		return TelegramUser{}, ErrInitDataExpired
	}

	var user TelegramUser
	if err := json.Unmarshal([]byte(values.Get("user")), &user); err != nil {
		return TelegramUser{}, fmt.Errorf("%w: parse user: %v", ErrInvalidInitData, err)
	}
	if user.ID == 0 {
		return TelegramUser{}, fmt.Errorf("%w: empty user id", ErrInvalidInitData)
	}

	return user, nil
}

func hmacSHA256(key, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}
