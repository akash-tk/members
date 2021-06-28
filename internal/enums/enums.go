package enums

type Role string

func (r Role) String() string {
	return string(r)
}


const (
	RoleMember Role = "member"
	RoleAdmin  Role = "admin"
)
