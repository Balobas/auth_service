package entity

import (
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Device struct {
	Uid           uuid.UUID
	MaintainerUid uuid.UUID
	Type          DeviceType
	Name          string
	Agent         string
	Language      string
	CreatedAt     time.Time
}

type DeviceType string

const (
	DeviceTypeUserDevice   = DeviceType("user_device")
	DeviceTypeSystemDevice = DeviceType("system_device")
)

func (dt DeviceType) Validate() error {
	switch dt {
	case DeviceTypeUserDevice, DeviceTypeSystemDevice:
		return nil
	default:
		return fmt.Errorf("invalid device type")
	}
}

func (ud Device) Validate() error {
	if uuid.Equal(ud.Uid, uuid.UUID{}) {
		return fmt.Errorf("empty device uid")
	}
	return nil
}

func NewUserDevice(userUid uuid.UUID, deviceData LoginDeviceData) (Device, error) {
	d := Device{
		Uid:           deviceData.Uid,
		MaintainerUid: userUid,
		Type:          DeviceTypeUserDevice,
		Name:          deviceData.Name,
		Agent:         deviceData.Agent,
		Language:      deviceData.Agent,
	}

	return d, d.Validate()
}

func NewSystemDevice(serviceUid uuid.UUID, deviceData LoginDeviceData) (Device, error) {
	d := Device{
		Uid:           deviceData.Uid,
		MaintainerUid: serviceUid,
		Type:          DeviceTypeSystemDevice,
		Name:          deviceData.Name,
		Agent:         deviceData.Agent,
		Language:      deviceData.Agent,
	}

	return d, d.Validate()
}

func (ud Device) WithUserUid(userUid uuid.UUID) Device {
	ud.MaintainerUid = userUid
	return ud
}

func (ud Device) Authorize(loginTime time.Time) AuthorizedDevice {
	return AuthorizedDevice{
		Device:       ud,
		AuthorizedAt: loginTime,
	}
}

func (ud Device) IsSystem() bool {
	return ud.Type == DeviceTypeSystemDevice
}

type AuthorizedDevice struct {
	Device
	AuthorizedAt   time.Time
	UnauthorizedAt *time.Time
}

func (ad AuthorizedDevice) Unauthorize(logoutTime time.Time) AuthorizedDevice {
	ad.UnauthorizedAt = &logoutTime
	return ad
}
