package entity

import uuid "github.com/satori/go.uuid"

type LoginParams struct {
	Email    string
	Password string
	Device   LoginDeviceData
}

type LoginDeviceData struct {
	Uid      uuid.UUID
	Name     string
	Agent    string
	Language string
	Domain   string // для системных девайсов
}

type ServiceLoginParams struct {
	Uid      uuid.UUID
	Password string
	Device   LoginDeviceData
}
