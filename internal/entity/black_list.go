package entity

import uuid "github.com/satori/go.uuid"

type BlackListElement struct {
	Uid    uuid.UUID
	Reason string
}

type BlackListUser struct {
	User   User
	Reason string
}

type BlackListDevice struct {
	Device UserDevice
	Reason string
}
