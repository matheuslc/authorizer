package account

import es "github.com/matheuslc/authorizer/internal/eventstore"

import "encoding/json"

// JSONResponse
type JSONResponse struct {
	Account es.Event `json:"account"`
}

// JSONPresenter
func JSONPresenter(event es.Event) string {
	response := JSONResponse{Account: event}
	formatted, _ := json.Marshal(response)

	return string(formatted)
}
