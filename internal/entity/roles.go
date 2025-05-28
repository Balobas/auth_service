package entity

type Role struct {
	Role        string
	Description string
}

type RoleWithPermissions struct {
	Role        string
	Permissions []string
}

const (
	UserRoleAdmin = "admin"
	UserRoleUser  = "user"
)

func HasRole(roles []string, role string) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

func RolesToStrings(roles []Role) []string {
	res := make([]string, len(roles))
	for i, r := range roles {
		res[i] = r.Role
	}
	return res
}