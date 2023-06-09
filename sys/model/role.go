package model

import (
	"sort"
	"strings"
)

const (
	PUBLIC Role = -1
	//Unset/Zero value Role = 0
	GUEST      Role = 1
	USER       Role = 5
	CREATOR    Role = 7
	ADMIN      Role = 10
	SUPERADMIN Role = 50
	ROOT       Role = 100

	Title NameFormat = 1
	Lower NameFormat = 2
	Upper NameFormat = 3
)

type NameFormat int
type Role int

type RoleSet struct {
	Name string `json:"name"`
	Role Role   `json:"role"`
}

var roles = map[string]Role{
	"public":     PUBLIC,
	"guest":      GUEST,
	"user":       USER,
	"creator":    CREATOR,
	"admin":      ADMIN,
	"superadmin": SUPERADMIN,
	"root":       ROOT,
}

var sortedRoles []Role

func init() {
	sortedRoles = make([]Role, len(roles))
	i := 0
	for _, v := range roles {
		sortedRoles[i] = v
		i++
	}
	sort.SliceStable(sortedRoles, func(i, j int) bool {
		return sortedRoles[i] < sortedRoles[j]
	})
}

// String returns every character in lower case like "admin"
func (r Role) String() string {
	for k, v := range roles {
		if v == r {
			return k
		}
	}
	return "unknown"
}

// StringToRole takes role string case insensitive and returns a GUEST role if no role was found with the provided string
func StringToRole(role string) Role {
	lowerCaseRole := strings.ToLower(role)
	if r, ok := roles[lowerCaseRole]; ok {
		return r
	}
	return GUEST
}

// Lower returns every character in lower case like "admin"
func (r Role) Lower() string {
	return r.String()
}

// Title returns the first character with upper case like "Admin"
func (r Role) Title() string {
	return strings.Title(r.String())
}

// Upper returns every character in upper case like "ADMIN"
func (r Role) Upper() string {
	return strings.ToUpper(r.String())
}

// Is checks if current role and provided one are equal
func (r Role) Is(pr Role) bool {
	return pr == r
}

// IsGrantedFor checks if current is higher or equal to the provided role
func (r Role) IsGrantedFor(pr Role) bool {
	return pr <= r
}

func (r Role) AllowedToCreateUserData() bool {
	return USER <= r
}

func (r Role) AllowedToCreateEntities() bool {
	return CREATOR <= r
}

// IsGrantedForUserModifications checks if current is higher or equal to SUPERADMIN
func (r Role) IsGrantedForUserModifications() bool {
	return SUPERADMIN <= r
}

// RolesInRange provides all roles hierarchically sorted in your range
func (r Role) RolesInRange() []RoleSet {
	return r.RolesInRangeWithNameFormat(Title)
}

func (r Role) nameAs(nf NameFormat) string {
	switch nf {
	case Lower:
		return r.Lower()
	case Upper:
		return r.Upper()
	default:
		return r.Title()
	}
}

func (r Role) RolesInRangeWithNameFormat(nf NameFormat) []RoleSet {
	rs := make([]RoleSet, 0)
	for _, ro := range sortedRoles {
		if r.IsGrantedFor(ro) {
			rs = append(rs, RoleSet{Name: ro.nameAs(nf), Role: ro})
		}
	}
	return rs
}
