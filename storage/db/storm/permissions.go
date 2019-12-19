package storm

import (
	"github.com/asdine/storm/q"

	"github.com/ProxeusApp/proxeus-core/sys/model"
)

func IsReadGrantedFor(auth model.Auth, includeGrant bool) q.Matcher {
	if auth.AccessRights().IsGrantedFor(model.SUPERADMIN) {
		return q.True()
	}
	var m q.Matcher
	if includeGrant {
		m = q.Or(
			q.Eq("Owner", auth.UserID()),
			//use publicByID only for direct access by ID
			q.NewFieldMatcher("Grant", &GrantMatcher{Auth: auth, CheckWrite: false}),
			q.NewFieldMatcher("GroupAndOthers", &GroupAndOthersMatcher{Auth: auth, CheckWrite: false}),
		)
	} else {
		m = q.Or(
			q.Eq("Owner", auth.UserID()),
		)
	}
	return m
}

type PublicByIDMatcher struct {
	CheckWrite bool //false checks for read perm
}

func (me *PublicByIDMatcher) MatchField(v interface{}) (bool, error) {
	if prm, ok := v.(model.Permission); ok {
		if me.CheckWrite {
			return prm.IsWrite(), nil
		} else {
			return prm.IsRead(), nil
		}
	}
	return false, nil
}

type GrantMatcher struct {
	Auth       model.Auth
	CheckWrite bool //false checks for read perm
}

func (me *GrantMatcher) MatchField(v interface{}) (bool, error) {
	if me.Auth.UserID() != "" {
		if grant, ok := v.(map[string]model.Permission); ok {
			if prm, ok := grant[me.Auth.UserID()]; ok {
				if me.CheckWrite {
					return prm.IsWrite(), nil
				} else {
					return prm.IsRead(), nil
				}
			}
		}
	}
	return false, nil
}

type GroupAndOthersMatcher struct {
	Auth       model.Auth
	CheckWrite bool //false checks for read perm
}

func (me *GroupAndOthersMatcher) MatchField(v interface{}) (bool, error) {
	if grpAndOthrs, ok := v.(model.GroupAndOthers); ok {
		grpAndOthrs.IsOthersRead()
		if me.CheckWrite {
			if grpAndOthrs.IsGroupWrite(me.Auth) {
				return true, nil
			}
			if grpAndOthrs.IsOthersWrite() {
				return true, nil
			}
		} else {
			if grpAndOthrs.IsGroupRead(me.Auth) {
				return true, nil
			}
			if grpAndOthrs.IsOthersRead() {
				return true, nil
			}
		}
	}
	return false, nil
}
