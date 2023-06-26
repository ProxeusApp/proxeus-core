package model

import (
	"fmt"
	"time"
)

/*
*
The actual version of the User struct.
If the structure changes, increase this value, to upgrade persisted data and implement the upgrade!
*/
const userVersion = 2

type (
	User struct {
		ID      string    `json:"id" storm:"id"`
		Name    string    `json:"name" storm:"index"`
		Detail  string    `json:"detail"`
		Updated time.Time `json:"updated" storm:"index"`
		Created time.Time `json:"created" storm:"index"`

		//Item
		Email        string      `json:"email" storm:"unique"`
		Role         Role        `json:"role,omitempty"`
		Data         interface{} `json:"data,omitempty"`
		Active       bool        `json:"active"`
		EthereumAddr string      `json:"etherPK" storm:"index"`
		PhotoPath    string      `json:"-"`
		Photo        string      `json:"photo"`
		//the user is able to tell if he lets other people see his photo, name or email instead of just the blockchain address
		WantToBeFound bool `json:"wantToBeFound"`

		ApiKeys []*ApiKey `json:"apiKeys"`
	}

	UserData struct {
		ID          string      `json:"id" storm:"id"`
		Name        string      `json:"name" storm:"index"`
		Detail      string      `json:"detail"`
		Permissions Permissions `json:"permissions"`
		Updated     time.Time   `json:"updated" storm:"index"`
		Created     time.Time   `json:"created" storm:"index"`

		WorkflowID string      `json:"workflowId"`
		Lang       string      `json:"lang"`
		LangForm   string      `json:"langForm"`
		LangTmpl   string      `json:"langTmpl"`
		Finished   bool        `json:"finished"`
		Data       interface{} `json:"data"`
	}
)

// IsGrantedFor check if user has enough permissions
func (me *User) IsGrantedFor(role Role) bool {
	return role <= me.Role
}

// ----Auth interface----------------
func (me *User) UserID() string {
	return me.ID
}

func (me *User) AccessRights() Role {
	return me.Role
}

//-------------------------------------------

func (me *User) CheckIfAuthIsAllowedToReadPersonalData(auth Auth) bool {
	if auth.AccessRights().IsGrantedForUserModifications() {
		return true
	}
	if auth.UserID() != me.ID && !me.WantToBeFound {
		//clear all personal data except blockchain address
		me.Email = ""
		me.Detail = ""
		me.Photo = ""
		me.PhotoPath = ""
		me.Name = ""
		return false
	}
	//allowed
	return true
}

func (me *User) String() string {
	return fmt.Sprintf(`{"id":"%v", "name":"%s", "detail":"%s", "updated":"%v", "created": "%v", "role":"%v", "active":%v, "etherAddr":"%s"}`, me.ID, me.Name, me.Detail, me.Updated, me.Created, me.Role, me.Active, me.EthereumAddr)
}

func (me *User) GetVersion() int {
	return userVersion
}

func (me *User) SetApiKey(name string) (*ApiKey, error) {
	apiKey, err := NewApiKey(name, me.ID)
	if err != nil {
		return nil, err
	}
	me.ApiKeys = append(me.ApiKeys, apiKey)
	return apiKey, nil
}

func (me *User) Close() {}
