package patient

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gabriel-ross/gofhir"
)

type response struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type listResponse struct {
	Data   []response `json:"data"`
	Count  int        `json:"count"`
	Offset int        `json:"offset"`
	Limit  int        `json:"limit"`
	Prev   string     `json:"prev,omitempty"`
	Next   string     `json:"next,omitempty"`
	Self   string     `json:"self"`
}

func (svc *Service) newResponse(p gofhir.Patient) response {
	return response{
		ID:   p.ID,
		Name: p.Name,
	}
}

func (svc *Service) RenderResponse(w http.ResponseWriter, r *http.Request, p gofhir.Patient, code int) {
	resp := svc.newResponse(p)
	body, err := json.Marshal(resp)
	if err != nil {
		gofhir.RenderError(w, r, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(code)
	w.Write(body)
}

func (svc *Service) RenderListResponse(w http.ResponseWriter, r *http.Request, code int, users []gofhir.Patient, offset, limit, count int) {
	data := []response{}
	for _, user := range users {
		data = append(data, svc.newResponse(user))
	}

	resp := listResponse{
		Data:   data,
		Count:  count,
		Offset: offset,
		Limit:  limit,
		Self:   svc.absolutePath,
	}

	if offset > 0 {
		resp.Prev = fmt.Sprintf("%s?offset=%d&limit=%d", svc.absolutePath, max(0, offset-limit), limit)
	}
	if offset+limit < count {
		resp.Next = fmt.Sprintf("%s?offset=%d&limit=%d", svc.absolutePath, offset+limit, limit)
	}

	body, err := json.Marshal(resp)
	if err != nil {
		gofhir.RenderError(w, r, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(code)
	w.Write(body)
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}
