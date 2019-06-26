package users

import (
	"net/http"

	f "github.com/fauna/faunadb-go/faunadb"
	"github.com/gorilla/mux"

	utils "github.com/kamilpowalowski/recommendations-go/api/utils"
)

// Entry - db entry model
type Entry struct {
	Token  string `fauna:"token"`
	UserID string `fauna:"user_id"`
	ItemID string `fauna:"item_id"`
	Value  int    `fauna:"value"`
}

// CreateEntry - returns entry object from data send in request
func CreateEntry(r *http.Request, value int) Entry {
	token := utils.ExtractToken(r.Header)
	vars := mux.Vars(r)
	userID := vars["user_id"]
	itemID := vars["item_id"]

	return Entry{
		Token:  token,
		UserID: userID,
		ItemID: itemID,
		Value:  value,
	}
}

// DBGetAllRefs - get all elements
func DBGetAllRefs(client *f.FaunaClient, token string) (refs []f.RefV, err error) {
	value, err := client.Query(
		f.Paginate(
			f.MatchTerm(
				f.Index("entries_with_token"),
				token,
			),
		),
	)
	if err != nil {
		return nil, err
	}

	value.At(f.ObjKey("data")).Get(&refs)
	return refs, nil
}

// DBGetFromRefs - get all elements
func DBGetFromRefs(client *f.FaunaClient, refs []f.RefV) (entries []Entry, err error) {
	request := mapRefV(refs, func(ref f.RefV) interface{} {
		return f.Get(ref)
	})
	value, err := client.Query(f.Arr(request))

	if err != nil {
		return nil, err
	}

	var elements f.ArrayV
	value.Get(&elements)

	results := make([]Entry, len(elements))
	for index, element := range elements {
		var object f.ObjectV
		element.At(f.ObjKey("data")).Get(&object)
		var entry Entry
		object.Get(&entry)
		results[index] = entry
	}

	return results, nil
}

// DBGet - get existing element from database
func (entry Entry) DBGet(client *f.FaunaClient) (value f.Value, err error) {
	return client.Query(
		f.Get(
			f.MatchTerm(
				f.Index("entry_with_token_user_item"),
				f.Arr{entry.Token, entry.UserID, entry.ItemID},
			),
		),
	)
}

// DBCreate - create new Entry object
func (entry Entry) DBCreate(client *f.FaunaClient) (value f.Value, err error) {
	return client.Query(
		f.Create(
			f.Class("entries"),
			f.Obj{"data": entry},
		),
	)
}

// DBUpdate - update existing object provided in result parameter
func (entry Entry) DBUpdate(client *f.FaunaClient, result f.Value) (value f.Value, err error) {
	var ref f.RefV
	result.At(f.ObjKey("ref")).Get(&ref)
	return client.Query(
		f.Update(
			ref,
			f.Obj{"data": entry},
		),
	)
}

// DBCreateOrUpdate - combine DBGet, DBCreate and DBUpdate to make uperation easier
func (entry Entry) DBCreateOrUpdate(client *f.FaunaClient) (value f.Value, err error) {
	value, _ = entry.DBGet(client)

	if value == nil {
		value, err = entry.DBCreate(client)
	} else {
		value, err = entry.DBUpdate(client, value)
	}
	return value, err
}

// DBDelete - remove Entry object from database
func (entry Entry) DBDelete(client *f.FaunaClient) (value f.Value, err error) {
	result, err := entry.DBGet(client)
	if result != nil {
		var ref f.RefV
		result.At(f.ObjKey("ref")).Get(&ref)
		return client.Query(f.Delete(ref))
	}
	return result, err
}

func mapRefV(vs []f.RefV, f func(f.RefV) interface{}) []interface{} {
	vsm := make([]interface{}, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}
