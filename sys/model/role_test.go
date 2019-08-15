package model

import (
	"testing"
)

func TestRoleIsAndIsGrantedFor(t *testing.T) {
	admin := ADMIN
	superadmin := SUPERADMIN
	if admin.Is(superadmin) {
		t.Error("cannot be true")
	}
	if admin.IsGrantedFor(SUPERADMIN) {
		t.Error("cannot be true")
	}
	if !admin.IsGrantedFor(ADMIN) {
		t.Error("cannot be false")
	}
	if !admin.IsGrantedFor(USER) {
		t.Error("cannot be false")
	}
}

func TestRole_Strings(t *testing.T) {
	if ADMIN.Lower() != "admin" {
		t.Error("expected admin", ADMIN.Lower())
	}
	if ADMIN.Title() != "Admin" {
		t.Error("expected Admin", ADMIN.Title())
	}
	if ADMIN.Upper() != "ADMIN" {
		t.Error("expected ADMIN", ADMIN.Upper())
	}
}

func TestRole_StringToRole(t *testing.T) {
	if StringToRole("admin") != ADMIN {
		t.Error("expected Admin", StringToRole("admin"))
	}
	if StringToRole("Admin") != ADMIN {
		t.Error("expected Admin", StringToRole("Admin"))
	}
	if StringToRole("adn") != GUEST {
		t.Error("expected Admin", StringToRole("adn"))
	}
}

func TestRolesByRole(t *testing.T) {
	cr := CREATOR.RolesInRange()
	if len(cr) == 0 || cr[len(cr)-1].Role != CREATOR {
		t.Error("expected Creator got", cr[len(cr)-1].Name)
	}
}
