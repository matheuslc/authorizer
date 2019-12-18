package transaction

import es "github.com/matheuslc/authorizer/internal/eventstore"

import "encoding/json"

// JSONPresenter
func JSONPresenter(event es.Event) string {
	formatted, _ := json.Marshal(event)
	return string(formatted)
}
