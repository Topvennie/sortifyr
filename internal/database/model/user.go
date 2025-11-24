// Package model contains all databank models
package model

import "github.com/topvennie/spotify_organizer/pkg/sqlc"

type User struct {
	ID          int
	UID         string
	Name        string
	DisplayName string
	Email       string
}

func UserModel(user sqlc.User) *User {
	displayName := ""
	if user.DisplayName.Valid {
		displayName = user.DisplayName.String
	}

	return &User{
		ID:          int(user.ID),
		UID:         user.Uid,
		Name:        user.Name,
		DisplayName: displayName,
		Email:       user.Email,
	}
}

// Equal returns true if all non unique values are equal
func (u *User) Equal(u2 User) bool {
	return u.Name == u2.Name && u.DisplayName == u2.DisplayName && u.Email == u2.Email
}
