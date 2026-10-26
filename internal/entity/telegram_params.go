package entity

// TelegramAuthParams — вход/регистрация через Telegram Mini App.
type TelegramAuthParams struct {
	InitData        string
	Device          LoginDeviceData
	CreateIfMissing bool // true — завести аккаунт, если telegram ещё не привязан
}

// TelegramLinkParams — привязка Telegram к существующему аккаунту (вход по email+паролю).
type TelegramLinkParams struct {
	InitData string
	Device   LoginDeviceData
	Email    string
	Password string
}
