package entity

import (
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Service struct {
	Uid       uuid.UUID
	Name      string
	Domain    string
	Roles     []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewService(uid uuid.UUID, name string, domain string, createdTime time.Time) (Service, error) {
	if len(name) == 0 {
		return Service{}, fmt.Errorf("empty name")
	}
	if len(domain) == 0 {
		return Service{}, fmt.Errorf("empty domain")
	}

	return Service{
		Uid:       uid,
		Name:      name,
		Domain:    domain,
		CreatedAt: createdTime.UTC(),
		UpdatedAt: createdTime.UTC(),
	}, nil
}

func (s Service) WithRoles(roles []string) Service {
	s.Roles = roles
	return s
}

type RegisterServiceRequest struct {
	Uid      uuid.UUID
	Name     string
	Domain   string
	Password string
}
