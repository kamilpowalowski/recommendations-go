package users

import (
	"net/http"

	utils "github.com/kamilpowalowski/recommendations-go/api/utils"
)

// Dislike - add new dislike item to database
func Dislike(w http.ResponseWriter, r *http.Request) {
	client := utils.DBClient()
	entry := CreateEntry(r, -1)

	_, err := entry.DBCreateOrUpdate(client)

	if err != nil {
		utils.SendInternalServerError(w, err)
		return
	}

	utils.SendSuccess(w)
}

// DeleteDislike - remove dislike item from database
func DeleteDislike(w http.ResponseWriter, r *http.Request) {
	client := utils.DBClient()
	entry := CreateEntry(r, 0)

	_, err := entry.DBDelete(client)

	if err != nil {
		utils.SendInternalServerError(w, err)
		return
	}

	utils.SendSuccess(w)
}
