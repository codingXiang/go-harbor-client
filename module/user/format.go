package user

import (
	"time"
)

// User holds the details of a user.
type User struct {
	UserID       int       `json:"user_id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	Realname     string    `json:"realname"`
	Comment      string    `json:"comment"`
	Deleted      bool      `json:"deleted"`
	Rolename     string    `json:"role_name"`
	Role         int       `json:"role_id"`
	HasAdminRole bool      `json:"has_admin_role"`
	ResetUUID    string    `json:"reset_uuid"`
	Salt         string    `json:"-"`
	CreationTime time.Time `json:"creation_time"`
	UpdateTime   time.Time `json:"update_time"`
}

type UpdateRole struct {
	UserID       int `short:"i" long:"user_id" description:"(REQUIRED) Registered user ID." required:"yes" json:"-"`
	HasAdminRole int `short:"r" long:"has_admin_role" description:"(REQUIRED) Toggle a user to admin or not." required:"yes" json:"has_admin_role"`
}

type UpdatePassword struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
