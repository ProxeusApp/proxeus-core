package model

import "time"

type Session struct {
	ID         string    `json:"id"`
	Rights     Role      `json:"rights"`
	UsrID      string    `json:"usrId"`
	UserName   string    `json:"userName"`
	SessionDir string    `json:"sessionDir"`
	Updated    time.Time `json:"updated"`
}
