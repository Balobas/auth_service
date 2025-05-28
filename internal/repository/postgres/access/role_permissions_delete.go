package repositoryAccess

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/pkg/errors"
)

func (r *AccessRepository) DeletePermissionsFromRole(ctx context.Context, role string, permissions []string) error {
	log.Printf("accessRepository.DeletePermissionsFromRole: role %s permissions %v", role, permissions)

	if len(role) == 0 {
		return errors.New("empty role")
	}
	if len(permissions) == 0 {
		return errors.New("empty permissions")
	}

	stmt := strings.Builder{}
	stmt.WriteString("DELETE FROM roles_permissions WHERE role = $1 AND permission IN (")

	args := make([]interface{}, 0, len(permissions)+1)
	args = append(args, role)

	stmt.WriteString("$2")
	args = append(args, permissions[0])
	idx := 3

	for i := 1; i < len(permissions); i++ {
		stmt.WriteString(fmt.Sprintf(", $%d", idx))
		idx++
		args = append(args, permissions[i])
	}
	stmt.WriteByte(')')

	if _, err := r.DB().Exec(ctx, stmt.String(), args...); err != nil {
		log.Printf("accessRepository.DeletePermissionsFromRole: failed :%v", err)
		return errors.WithStack(err)
	}
	return nil
}
