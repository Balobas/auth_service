package repositoryAccess

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/pkg/errors"
)

func (r *AccessRepository) AddPermissionsToRole(ctx context.Context, role string, permissions []string) error {
	log.Printf("accessRepository.AddPermissionsToRole: role %s permissions %v", role, permissions)

	if len(role) == 0 {
		return errors.New("empty role")
	}
	if len(permissions) == 0 {
		return errors.New("empty permissions")
	}

	stmt := strings.Builder{}
	stmt.WriteString("INSERT INTO roles_permissions (role, permission) VALUES ")

	args := make([]interface{}, 0, len(permissions)*2)

	stmt.WriteString("($1, $2)")
	args = append(args, role, permissions[0])
	idx := 3

	for i := 1; i < len(permissions); i++ {
		stmt.WriteString(fmt.Sprintf(", ($%d, $%d)", idx, idx+1))
		idx += 2
		args = append(args, role, permissions[i])
	}

	if _, err := r.Exec(ctx, stmt.String(), args...); err != nil {
		log.Printf("accessRepository.AddPermissionsToRole: failed :%v", err)
		return errors.WithStack(err)
	}
	return nil
}
