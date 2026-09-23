package constants

// 用户角色枚举。
const (
	RoleAdmin  = "admin"  // 管理员
	RoleStaff  = "staff"  // 店员
	RoleMember = "member" // 会员
)

// AllRoles 所有角色列表。
var AllRoles = []string{RoleAdmin, RoleStaff, RoleMember}

// IsValidRole 判断角色是否合法。
func IsValidRole(role string) bool {
	switch role {
	case RoleAdmin, RoleStaff, RoleMember:
		return true
	}
	return false
}
