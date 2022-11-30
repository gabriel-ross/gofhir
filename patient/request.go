package patient

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gabriel-ross/gofhir"
)

type request struct {
	Name string `json:"name"`
}

// BindRequest binds the fields defined in request of a request to a User.
// This method also extracts the token from the header "Token".
func BindRequest(r *http.Request, p *gofhir.Patient) (err error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	var reqBody request
	err = json.Unmarshal(body, &reqBody)
	if err != nil {
		return err
	}

	p.Name = reqBody.Name
	return nil
}
