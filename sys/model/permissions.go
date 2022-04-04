package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

var ErrAuthorityMissing = fmt.Errorf("authority missing")
var ErrAuthorityInvalid = fmt.Errorf("invalid authority key")
var ErrAuthorityNotFound = fmt.Errorf("user not found for key")

//Permission holds an byte slice for the pattern --/rw or ----/r-r- and so on
// none     read    write
//0 == - | 1 == r | 2 == w
//to support []byte we implement the json marshaller interface otherwise use []int
type Permission []byte

const read = 1
const write = 2

type Auth interface {
	UserID() string
	AccessRights() Role
}

type MemoryAuth interface {
	Auth
	GetMemory(k string) (interface{}, bool)
	PutMemory(k string, val interface{})
	DeleteMemory(k string)
	GetSessionDir() string
}

type AccessibleItem interface {
	IsReadGrantedFor(auth Auth) bool
	IsWriteGrantedFor(auth Auth) bool
	OwnedBy(auth Auth) bool
}

type GroupAndOthers struct {
	//Allowed to modify: @Owner only!
	Group Role `json:"group,omitempty"`
	//Rights pattern:       	    					group others
	//		                     	 	 				 --     --
	//default value is:                	    				----
	//example for group and others with read perm:	    	r-r-
	//Allowed to modify: @Owner only!
	Rights Permission `json:"rights,omitempty"`
}

type Permissions struct {
	//Allowed to modify: @Owner only!
	Owner string `json:"owner,omitempty"`
	//Grant can be modified by the owner only and is an optional field to whitelist user's directly
	//Allowed to modify: everyone with write rights!
	Grant map[string]Permission `json:"grant,omitempty"`

	GroupAndOthers GroupAndOthers `json:"groupAndOthers,omitempty"`
	//Accessible by everyone who has the ID
	//Allowed to modify: everyone with write rights!
	PublicByID Permission `json:"publicByID,omitempty"`

	//Execute only! If read or write not set.
	Published bool `json:"published"`
}

func PermissionFrom(readablePattern string) (Permission, error) {
	if readablePattern == "" {
		return nil, os.ErrInvalid
	}
	if len(readablePattern)%2 != 0 {
		return nil, os.ErrInvalid
	}
	p := make(Permission, len(readablePattern)/2)
	for i, ii := 0, 0; ii < len(readablePattern); i, ii = i+1, ii+2 {
		if readablePattern[ii:ii+1] == "r" {
			p[i] = read
		}
		if readablePattern[ii+1:ii+2] == "w" {
			p[i] = write
		}
	}
	return p, nil
}

func (me Permission) IsGroupRead() bool {
	return me.posHasRead(0)
}

func (me Permission) IsGroupWrite() bool {
	return me.posHasWrite(0)
}

func (me Permission) IsRead() bool {
	return me.posHasRead(0)
}

func (me Permission) IsWrite() bool {
	return me.posHasWrite(0)
}

func (me Permission) IsOthersRead() bool {
	return me.posHasRead(1)
}

func (me Permission) IsOthersWrite() bool {
	return me.posHasWrite(1)
}

func (me Permission) MarshalJSON() ([]byte, error) {
	i := make([]int, len(me))
	for a, b := range me {
		i[a] = int(b)
	}
	return json.Marshal(i)
}

func (me Permission) posHasRead(pos int) bool {
	return hasRead(me, pos)
}

func hasRead(me []byte, pos int) bool {
	l := len(me)
	if l == 0 {
		return false
	}
	if l-1 < pos || pos < 0 {
		return false
	}
	return me[pos] >= read
}

func (me Permission) posHasWrite(pos int) bool {
	return hasWrite(me, pos)
}

func hasWrite(me []byte, pos int) bool {
	l := len(me)
	if l == 0 {
		return false
	}
	if l-1 < pos || pos < 0 {
		return false
	}
	return me[pos] == write
}

func (me Permission) ToReadablePattern() string {
	if len(me) == 0 {
		return "--"
	}
	b := &bytes.Buffer{}
	for _, v := range me {
		if v == 0 {
			b.WriteString("--")
		} else if v == read {
			b.WriteString("r-")
		} else if v == write {
			b.WriteString("rw")
		}
	}
	return b.String()
}

func (me Permission) String() string {
	return me.ToReadablePattern()
}

//TODO improve perm modification
func (me *Permissions) Change(auth Auth, changed *Permissions) *Permissions {
	if me.OwnedBy(auth) {
		if changed.Owner != "" {
			me.Owner = changed.Owner
		}
		me.GroupAndOthers = changed.GroupAndOthers
		me.Grant = changed.Grant
		me.PublicByID = changed.PublicByID
		me.Published = changed.Published
		return me
	}
	if me.IsWriteGrantedFor(auth) {
		me.Grant = changed.Grant
		me.PublicByID = changed.PublicByID
		me.Published = changed.Published
	}
	return me
}

func (me *Permissions) UserIdsMap(usrIds map[string]bool) {
	if me.Owner != "" {
		usrIds[me.Owner] = true
	}
	for usrID := range me.Grant {
		if usrID != "" {
			usrIds[usrID] = true
		}
	}
}

func (me *Permissions) UpdateUserID(usrOldIdNewID map[string]string) {
	checks := len(me.Grant)
	if me.Owner != "" {
		if newId, exists := usrOldIdNewID[me.Owner]; exists {
			me.Owner = newId
		}
	}
	if checks > 0 {
		for oldId, newId := range usrOldIdNewID {
			if checks > 0 {
				if perm, ok := me.Grant[oldId]; ok {
					me.Grant[newId] = perm
					delete(me.Grant, oldId)
					checks--
				}
			} else {
				break
			}
		}
	}
}

//IsReadGrantedFor is checking whether the User has read rights or not
//Rights pattern:       	    					group others
//		                     	 	 				 --     --
//default value is:                	    				----
//example for group and others with read perm:	    	r-r-
func (me *Permissions) IsReadGrantedFor(auth Auth) bool {
	if auth == nil {
		return false
	}

	if auth.AccessRights().IsGrantedFor(SUPERADMIN) {
		return true
	}

	//check public
	if me.PublicByID.IsRead() {
		return true
	}

	//check owner
	if auth.UserID() != "" {
		if auth.UserID() == me.Owner {
			return true
		}
		if len(me.Grant) > 0 {
			//check specific grants
			if perm, ok := me.Grant[auth.UserID()]; ok {
				if perm.IsGroupRead() {
					return true
				}
			}
		}
	}

	//check group
	if me.GroupAndOthers.IsGroupRead(auth) {
		return true
	}

	//check other
	if me.GroupAndOthers.IsOthersRead() {
		return true
	}

	return false
}

func (me *GroupAndOthers) IsGroupRead(auth Auth) bool {
	//check group
	if auth.AccessRights() != 0 {
		if me.Group != 0 {
			if me.Group <= auth.AccessRights() && me.Rights.IsGroupRead() {
				//same or higher group and has read rights
				return true
			}
		}
	}
	return false
}

func (me *GroupAndOthers) IsGroupWrite(auth Auth) bool {
	//check group
	if auth.AccessRights() != 0 {
		if me.Group != 0 {
			if me.Group <= auth.AccessRights() && me.Rights.IsGroupWrite() {
				//same or higher group has write rights
				return true
			}
		}
	}
	return false
}

func (me *GroupAndOthers) IsOthersRead() bool {
	//check other
	return me.Rights.IsOthersRead()
}

func (me *GroupAndOthers) IsOthersWrite() bool {
	//check other
	return me.Rights.IsOthersWrite()
}

//IsWriteGrantedFor is checking whether the User has write rights or not
//Rights pattern:       	    					group others
//		                     	 	 				 --     --
//default value is:                	    				----
//example for group and others with read perm:	    	r-r-
func (me *Permissions) IsWriteGrantedFor(auth Auth) bool {
	if auth == nil {
		return false
	}

	if auth.AccessRights().IsGrantedFor(ROOT) {
		return true
	}

	//check public
	if me.PublicByID.IsWrite() {
		return true
	}

	//check owner
	if auth.UserID() != "" {
		if auth.UserID() == me.Owner {
			return true
		}
		//check specific grants
		if len(me.Grant) > 0 {
			if rights, ok := me.Grant[auth.UserID()]; ok && rights.IsWrite() {
				return true
			}
		}
	}

	//check group
	if me.GroupAndOthers.IsGroupWrite(auth) {
		return true
	}

	//check other
	if me.GroupAndOthers.IsOthersWrite() {
		return true
	}

	return false
}

func (me *Permissions) IsPublishedOrReadGrantedFor(auth Auth) bool {
	return me.IsReadGrantedFor(auth) || me.IsPublishedFor(auth)
}

func (me *Permissions) IsPublishedFor(auth Auth) bool {
	return /*activate if needed auth.AccessRights().IsGrantedFor(CREATOR) && */ me.Published
}

func (me *Permissions) OwnedBy(auth Auth) bool {
	if auth != nil && auth.UserID() != "" {
		if auth.UserID() == me.Owner {
			return true
		}
	}
	return false
}
