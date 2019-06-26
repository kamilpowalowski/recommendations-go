package users

import (
	"net/http"

	utils "github.com/kamilpowalowski/recommendations-go/api/utils"
)

// Like - add new like item to database
func Like(w http.ResponseWriter, r *http.Request) {
	client := utils.DBClient()
	entry := CreateEntry(r, 1)

	_, err := entry.DBCreateOrUpdate(client)

	if err != nil {
		utils.SendInternalServerError(w, err)
		return
	}

	utils.SendSuccess(w)
}

// DeleteLike - remove like item from database
func DeleteLike(w http.ResponseWriter, r *http.Request) {
	client := utils.DBClient()
	entry := CreateEntry(r, 0)

	_, err := entry.DBDelete(client)

	if err != nil {
		utils.SendInternalServerError(w, err)
		return
	}

	utils.SendSuccess(w)
}
