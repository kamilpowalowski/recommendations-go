package users

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/muesli/regommend"

	utils "github.com/kamilpowalowski/recommendations-go/api/utils"
)

// distancePair - struct that represents result from Recommendation function
type distancePair struct {
	ItemID   interface{} `json:"item_id"`
	Distance float64     `json:"distance"`
}

// Recommendations - get recommendations for given user
func Recommendations(w http.ResponseWriter, r *http.Request) {
	client := utils.DBClient()
	token := utils.ExtractToken(r.Header)
	vars := mux.Vars(r)
	userID := vars["user_id"]

	var err error

	refs, err := DBGetAllRefs(client, token)
	if err != nil {
		utils.SendInternalServerError(w, err)
		return
	}

	entries, err := DBGetFromRefs(client, refs)
	if err != nil {
		utils.SendInternalServerError(w, err)
		return
	}

	data := makeRegommendData(entries)

	recommendations, err := getRecommendations(token, data, userID)

	if err != nil {
		utils.SendNotFound(w, err)
		return
	}

	utils.SendJSON(w, recommendations)
}

func makeRegommendData(entries []Entry) map[string]map[interface{}]float64 {
	data := make(map[string]map[interface{}]float64)

	for _, entry := range entries {
		var values map[interface{}]float64
		values, ok := data[entry.UserID]
		if ok {
			values[entry.ItemID] = float64(entry.Value)
		} else {
			values = map[interface{}]float64{
				entry.ItemID: float64(entry.Value),
			}
		}
		data[entry.UserID] = values
	}

	return data
}

func getRecommendations(token string, data map[string]map[interface{}]float64, userID string) ([]distancePair, error) {
	table := regommend.Table(token)

	for key, value := range data {
		table.Add(key, value)
	}

	recs, err := table.Recommend(userID)
	log.Println(recs)

	if err != nil {
		return nil, err
	}

	return mapDistancePairs(recs, func(item regommend.DistancePair) distancePair {
		return distancePair{
			ItemID:   item.Key,
			Distance: item.Distance,
		}
	}), nil
}

func mapDistancePairs(vs []regommend.DistancePair, f func(regommend.DistancePair) distancePair) []distancePair {
	vsm := make([]distancePair, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}
