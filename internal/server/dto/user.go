package dto

import "github.com/topvennie/spotify_organizer/internal/database/model"

type User struct {
	ID          int    `json:"id"`
	UID         string `json:"uid"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
}

func UserDTO(user *model.User) User {
	name := user.Name
	if user.DisplayName != "" {
		name = user.DisplayName
	}

	return User{
		ID:          user.ID,
		UID:         user.UID,
		Name:        name,
		DisplayName: user.DisplayName,
		Email:       user.Email,
	}
}

func (u *User) ToModel() *model.User {
	user := model.User(*u)
	return &user
}
