package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ONSdigital/log.go/log"
)

// json policy structs
type mockConditions struct {
	Attribute string
	Operator  string
	Value     string
}

type mockPolicy struct {
	Id         string
	Members    []string
	Roles      []string
	Conditions []mockConditions
}

//GetPoliciesHandler is a handler that returns a mocked policy to caller
func (api *API) GetPoliciesHandler(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	logdata := log.Data{"policy-id": "36A0C894-8902-4EAD-87B3-429C0C2EBBC7"}
	policy := mockPolicy{
		Id:      "36A0C894-8902-4EAD-87B3-429C0C2EBBC7",
		Members: []string{"GDPTeam"},
		Roles:   []string{"Readonly", "Publisher", "Admin"},
		Conditions: []mockConditions{
			{
				Attribute: "collection_id",
				Operator:  "equals",
				Value:     "some-collection-0207358238o57234802735925812739725",
			},
		},
	}

	var b []byte
	var err error
	b, err = json.Marshal(policy)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("json is:", string(b))

	//Set headers
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if _, err := w.Write(b); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Event(ctx, "getRole Handler: Successfully retrieved role", log.INFO, logdata)
}
