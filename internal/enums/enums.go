package enums

// Role represents the member's role in the GitHub organization.
// One can be either admin or member.
type Role string

func (r Role) String() string {
	return string(r)
}

const (
	RoleMember Role = "member"
	RoleAdmin  Role = "admin"
)
