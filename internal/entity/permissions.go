package entity

type Permission struct {
	Key         string
	Description string
}

type ResourcePermissions struct {
	URI        string
	Method     string
	Permission string
}
