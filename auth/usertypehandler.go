package auth

import (
	"github.com/crossworth/cartola-web-admin/database"
)

type UserTypeHandler struct {
	db          *database.PostgreSQL
	superAdmins []int
}

func NewUserTypeHandler(db *database.PostgreSQL, superAdmins []int) *UserTypeHandler {
	return &UserTypeHandler{db: db, superAdmins: superAdmins}
}

func (u *UserTypeHandler) GetUserType(id int) string {
	for _, sa := range u.superAdmins {
		if sa == id {
			return "super_admin"
		}
	}

	// fixme: handle database admin

	return "membro"
}
