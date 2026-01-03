package spotifyapi

import "github.com/topvennie/sortifyr/internal/database/model"

func accessKey(user model.User) string {
	return user.UID + ":spotify:access_token"
}

func refreshKey(user model.User) string {
	return user.UID + ":spotify:refresh_token"
}
