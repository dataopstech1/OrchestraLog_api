package middleware

import (
	"net/http"

	"github.com/orchestralog/api/pkg/apierror"
	"github.com/orchestralog/api/pkg/response"
)

type Permission struct {
	CanCreate         bool
	CanEdit           bool
	CanDelete         bool
	CanViewAll        bool
	CanManageClusters bool
	CanManageUsers    bool
}

var rolePermissions = map[string]Permission{
	"admin": {
		CanCreate: true, CanEdit: true, CanDelete: true,
		CanViewAll: true, CanManageClusters: true, CanManageUsers: true,
	},
	"operator": {
		CanCreate: true, CanEdit: true, CanDelete: false,
		CanViewAll: true, CanManageClusters: true, CanManageUsers: false,
	},
	"viewer": {
		CanCreate: false, CanEdit: false, CanDelete: false,
		CanViewAll: true, CanManageClusters: false, CanManageUsers: false,
	},
}

func GetPermissions(role string) Permission {
	return rolePermissions[role]
}

func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole := GetUserRole(r)
			for _, role := range roles {
				if userRole == role {
					next.ServeHTTP(w, r)
					return
				}
			}
			response.Error(w, apierror.ErrForbidden)
		})
	}
}

func RequireAdmin() func(http.Handler) http.Handler {
	return RequireRole("admin")
}

func RequireAdminOrOperator() func(http.Handler) http.Handler {
	return RequireRole("admin", "operator")
}
