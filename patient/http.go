package patient

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gabriel-ross/gofhir"
	"github.com/gabriel-ross/gofhir/hook"
	"github.com/go-chi/chi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (svc *Service) handleCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		interceptorCtx := hook.NewContext(strconv.Itoa(svc.NewRequestID()))
		ctx := context.TODO()
		var err error
		data := gofhir.Patient{}

		reqEvent := &hook.RequestEvent{
			Context:        interceptorCtx,
			Timestamp:      time.Now(),
			Request:        r,
			ResponseWriter: w,
		}
		for _, interceptor := range svc.RequestInterceptors {
			interceptor.OnRequestReceived(reqEvent)
			if interceptorCtx.ShouldAbort {
				return
			}
		}

		err = BindRequest(r, &data)
		if err != nil {
			gofhir.RenderError(w, r, http.StatusBadRequest, err)
			return
		}

		databaseQueryEvent := &hook.DatabaseQueryEvent{
			Context:    interceptorCtx,
			Timestamp:  time.Now(),
			Ctx:        ctx,
			Query:      "Collection(patients).Create()",
			Successful: true,
		}
		for _, interceptor := range svc.DatabaseInterceptors {
			interceptor.BeforeDatabaseQuery(databaseQueryEvent)
			if interceptorCtx.ShouldAbort {
				return
			}
		}
		resp, err := svc.create(ctx, data)

		if err != nil {
			databaseQueryEvent.Successful = false
			databaseQueryEvent.Error = err
		}
		for _, interceptor := range svc.DatabaseInterceptors {
			interceptor.AfterDatabaseQuery(databaseQueryEvent)
			if interceptorCtx.ShouldAbort {
				return
			}
		}

		if err != nil {
			gofhir.RenderError(w, r, http.StatusInternalServerError, err)
			return
		}

		responseEvent := &hook.ResponseEvent{
			Context:        interceptorCtx,
			Timestamp:      time.Now(),
			Request:        r,
			ResponseWriter: w,
			HTTPStatusCode: http.StatusOK,
		}

		// TODO: this is a good example as to why the interceptor calling func should be offloaded
		// to the service performing the action (in this case responding)
		for _, interceptor := range svc.ResponseInterceptors {
			interceptor.OnServerResponse(responseEvent)
			if interceptorCtx.ShouldAbort {
				return
			}
		}
		svc.RenderResponse(w, r, resp, http.StatusCreated)
	}
}

func (svc *Service) handleList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.TODO()
		var err error

		offset, limit, err := extractPaginate(r)
		if err != nil {
			gofhir.RenderError(w, r, http.StatusBadRequest, err)
			return
		}
		resp, err := svc.list(ctx, offset, limit)
		if err != nil {
			gofhir.RenderError(w, r, http.StatusInternalServerError, err)
			return
		}

		count, err := svc.count(ctx)
		if err != nil {
			gofhir.RenderError(w, r, http.StatusInternalServerError, err)
			return
		}

		svc.RenderListResponse(w, r, http.StatusOK, resp, offset, limit, count)
	}
}

func extractPaginate(r *http.Request) (_ int, _ int, err error) {
	offset := 0
	limit := 5

	if offsetParam := r.URL.Query().Get("offset"); offsetParam != "" {
		offset, err = strconv.Atoi(offsetParam)
		if err != nil {
			return 0, 0, err
		}
	}

	if limitParam := r.URL.Query().Get("limit"); limitParam != "" {
		limit, err = strconv.Atoi(limitParam)
		if err != nil {
			return 0, 0, err
		}
	}

	return offset, limit, nil
}

func (svc *Service) handleGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.TODO()
		var err error
		id := chi.URLParam(r, "id")

		resp, err := svc.read(ctx, id)
		if status.Code(err) == codes.NotFound {
			gofhir.RenderError(w, r, http.StatusNotFound, errors.New("resource not found"))
			return
		} else if err != nil {
			gofhir.RenderError(w, r, http.StatusInternalServerError, err)
			return
		}

		svc.RenderResponse(w, r, resp, http.StatusOK)
	}
}

func (svc *Service) handleUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.TODO()
		var err error
		id := chi.URLParam(r, "id")
		data := gofhir.Patient{}

		err = BindRequest(r, &data)
		if err != nil {
			w.Write([]byte("error binding: " + err.Error()))
			return
		}

		resp, err := svc.update(ctx, id, data)
		if err != nil {
			w.Write([]byte("error creating: " + err.Error()))
			return
		}

		svc.RenderResponse(w, r, resp, http.StatusNoContent)
	}
}

func (svc *Service) handleDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.TODO()
		var err error
		id := chi.URLParam(r, "id")

		exists, err := svc.exists(ctx, id)
		if err != nil {
			gofhir.RenderError(w, r, http.StatusInternalServerError, err)
			return
		}

		if !exists {
			gofhir.RenderError(w, r, http.StatusNotFound, errors.New("resource not found"))
			return
		}

		err = svc.delete(ctx, id)
		if err != nil {
			gofhir.RenderError(w, r, http.StatusInternalServerError, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
