package model

import (
	"encoding/json"
	"testing"
)

func TestItem_IsReadWriteGrantedFor(t *testing.T) {
	it := &FormItem{}
	it.Permissions.Owner = "123"
	usr := &User{}
	it.Permissions.Grant = make(map[string]Permission)
	usr.ID = "123"
	if !it.Permissions.IsReadGrantedFor(usr) || !it.Permissions.IsWriteGrantedFor(usr) {
		t.Error("Owner must be able to read and write without additional settings!")
	}

	usr.ID = "1234"

	if it.Permissions.IsReadGrantedFor(usr) || it.Permissions.IsWriteGrantedFor(usr) {
		t.Error("Not the owner!")
	}

	usr.ID = "1234"
	usr.Role = ADMIN
	it.Permissions.GroupAndOthers.Rights, _ = PermissionFrom("r---")
	it.Permissions.GroupAndOthers.Group = GUEST

	if !it.Permissions.IsReadGrantedFor(usr) {
		t.Error("group hierarchy not working!")
	}

	if it.Permissions.IsWriteGrantedFor(usr) {
		t.Error("write rights not set..!")
	}

	usr.Role = SUPERADMIN
	it.Permissions.GroupAndOthers.Rights, _ = PermissionFrom("rw--")
	it.Permissions.GroupAndOthers.Group = GUEST

	if !it.Permissions.IsReadGrantedFor(usr) {
		t.Error("group hierarchy not working!")
	}

	if !it.Permissions.IsWriteGrantedFor(usr) {
		t.Error("write rights set!")
	}

	usr.Role = PUBLIC
	it.Permissions.GroupAndOthers.Rights, _ = PermissionFrom("--r-")
	it.Permissions.GroupAndOthers.Group = GUEST

	if !it.Permissions.IsReadGrantedFor(usr) {
		t.Error("group hierarchy not working!")
	}

	if it.Permissions.IsWriteGrantedFor(usr) {
		t.Error("write rights not set!")
	}

	usr.Role = PUBLIC
	it.Permissions.GroupAndOthers.Rights, _ = PermissionFrom("--rw")
	it.Permissions.GroupAndOthers.Group = GUEST

	if !it.Permissions.IsReadGrantedFor(usr) {
		t.Error("group hierarchy not working!")
	}

	if !it.Permissions.IsWriteGrantedFor(usr) {
		t.Error("write rights set!")
	}

	usr.Role = PUBLIC
	it.Permissions.GroupAndOthers.Rights, _ = PermissionFrom("----")
	it.Permissions.Grant[usr.ID], _ = PermissionFrom("r-")
	it.Permissions.GroupAndOthers.Group = GUEST

	if !it.Permissions.IsReadGrantedFor(usr) {
		t.Error("read rights set!")
	}

	if it.Permissions.IsWriteGrantedFor(usr) {
		t.Error("write rights not set!")
	}

	usr.Role = PUBLIC
	it.Permissions.GroupAndOthers.Rights, _ = PermissionFrom("----")
	//"-w" should be the same as "rw"
	it.Permissions.Grant[usr.ID], _ = PermissionFrom("-w")
	it.Permissions.GroupAndOthers.Group = GUEST

	if !it.Permissions.IsReadGrantedFor(usr) {
		t.Error("read rights set!")
	}

	if !it.Permissions.IsWriteGrantedFor(usr) {
		t.Error("write rights not set!")
	}
}

func TestPermissions_UpdateUserID(t *testing.T) {
	it := &FormItem{}
	it.Permissions.Owner = "123"
	r, _ := PermissionFrom("r---")
	it.Permissions.Grant = map[string]Permission{"111": r, "12": r}
	it.Permissions.UpdateUserID(map[string]string{"111": "222", "123": "321", "12": "21"})
	if it.Permissions.Owner != "321" {
		t.Error("owner 123 -> 321 expected")
	}
	if _, ok := it.Permissions.Grant["222"]; !ok {
		t.Error("grant 111 -> 222 expected")
	}
	if _, ok := it.Permissions.Grant["21"]; !ok {
		t.Error("grant 12 -> 21 expected")
	}
}

func TestPermissionFrom(t *testing.T) {
	p, err := PermissionFrom("r--wr-")
	if err != nil || p.String() != "r-rwr-" {
		t.Error("wrong pattern!", err)
	}
	p, err = PermissionFrom("r-")
	if err != nil || p.String() != "r-" {
		t.Error("wrong pattern!", err)
	}
	p, err = PermissionFrom("rw")
	if err != nil || p.String() != "rw" {
		t.Error("wrong pattern!", err)
	}
	p, err = PermissionFrom("-w")
	if err != nil || p.String() != "rw" {
		t.Error("wrong pattern!", err)
	}
}

func TestPermissionSerializeDeserialize(t *testing.T) {
	p, err := PermissionFrom("r--w")
	if err != nil {
		t.Error("wrong pattern!", err)
	}
	bts, err := json.Marshal(p)
	if err != nil {
		t.Error("wrong pattern!", err)
	}

	if len(bts) != 5 {
		t.Error("expected 5!", string(bts))
	}
	if bts[1] != byte('1') {
		t.Error("expected 1!", bts[0])
	}
	if bts[3] != byte('2') {
		t.Error("expected 2!", bts[3])
	}
	pp := Permission{}
	err = json.Unmarshal(bts, &pp)
	if err != nil {
		t.Error(err)
	}
	if pp.String() != "r-rw" {
		t.Errorf("addr %p, expected r-rw, got %s", &pp, pp.String())
	}
}
