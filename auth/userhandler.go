package auth

import (
	"context"
	"strconv"
	"strings"

	"github.com/crossworth/cartola-web-admin/database"
	"github.com/crossworth/cartola-web-admin/logger"
	"github.com/crossworth/cartola-web-admin/model"
)

type UserHandler struct {
	db          *database.PostgreSQL
	superAdmins []int
}

func NewUserHandler(db *database.PostgreSQL, superAdmins []int) *UserHandler {
	return &UserHandler{db: db, superAdmins: superAdmins}
}

func (u *UserHandler) GetUserType(id int) string {
	for _, sa := range u.superAdmins {
		if sa == id {
			return "super_admin"
		}
	}

	admins, err := u.db.GetAdministratorProfiles(context.TODO())
	if err != nil {
		logger.Log.Warn().Err(err).Msg("não foi possível conseguir a lista de administradores")
		return "membro"
	}

	for _, admin := range admins {
		if admin.ID == id {
			return "admin"
		}
	}

	return "membro"
}

func (u *UserHandler) IsUserAllowed(id int) bool {
	rules, err := u.db.SettingByName(context.TODO(), model.MembersRuleSettingName)
	if err != nil {
		logger.Log.Warn().Err(err).Int("id", id).Msg("não foi possível conseguir as regras de membro")
		return false
	}

	rules = strings.TrimSpace(rules)
	rules = strings.ReplaceAll(rules, "\n\r", "\n")
	rulesLines := strings.Split(rules, "\n")

	if len(rulesLines) == 0 {
		return false
	}

	allowed := false

	for _, line := range rulesLines {
		if strings.TrimSpace(line) == "*" {
			allowed = true
		}

		if line == strconv.Itoa(id) {
			allowed = true
		}

		if line == "!"+strconv.Itoa(id) {
			allowed = false
		}
	}

	return allowed
}
