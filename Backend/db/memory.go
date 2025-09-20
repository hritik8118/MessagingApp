package db

import model "MessasingApp/Backend/models"

var users = make(map[string]model.User)

func UserExists(username string) bool {
	_, ok := users[username]
	return ok
}

func SaveUser(u model.User) {
	users[u.Username] = u
}

func GetUser(username string) (model.User, bool) {
	u, ok := users[username]
	return u, ok
}
