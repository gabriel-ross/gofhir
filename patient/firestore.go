package patient

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/gabriel-ross/gofhir"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (svc *Service) create(ctx context.Context, data gofhir.Patient) (_ gofhir.Patient, err error) {
	data.ID = svc.db.Collection("patients").NewDoc().ID
	_, err = svc.db.Collection("patients").Doc(data.ID).Set(ctx, data)
	if err != nil {
		return gofhir.Patient{}, err
	}
	return data, nil
}

func (svc *Service) list(ctx context.Context, offset, limit int) (_ []gofhir.Patient, err error) {
	resp := []gofhir.Patient{}
	iter := svc.db.Collection("patients").OrderBy("id", firestore.Asc).StartAt(offset).Limit(limit).Documents(ctx)
	for {
		dsnap, err := iter.Next()
		if err == iterator.Done {
			break
		}

		var m gofhir.Patient
		dsnap.DataTo(&m)
		resp = append(resp, m)
	}
	return resp, nil
}

func (svc *Service) count(ctx context.Context) (_ int, err error) {
	docs, err := svc.db.Collection("patients").Documents(ctx).GetAll()
	if err != nil {
		return 0, err
	}
	return len(docs), nil
}

func (svc *Service) read(ctx context.Context, id string) (_ gofhir.Patient, err error) {
	dsnap, err := svc.db.Collection("patients").Doc(id).Get(ctx)
	if err != nil {
		return gofhir.Patient{}, err
	}

	var m gofhir.Patient
	dsnap.DataTo(&m)
	return m, nil
}

func (svc *Service) exists(ctx context.Context, id string) (_ bool, err error) {
	_, err = svc.db.Collection("patients").Doc(id).Get(ctx)
	if status.Code(err) == codes.NotFound {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (svc *Service) set(ctx context.Context, id string, data gofhir.Patient) (_ gofhir.Patient, err error) {
	data.ID = id
	_, err = svc.db.Collection("patients").Doc(data.ID).Set(ctx, data)
	if err != nil {
		return gofhir.Patient{}, err
	}
	return data, nil
}

func (svc *Service) updateNonZero(ctx context.Context, id string, data gofhir.Patient) (_ gofhir.Patient, err error) {

	// Build update slice
	updates := []firestore.Update{}

	if data.Name != "" {
		updates = append(updates, firestore.Update{
			Path:  "name",
			Value: data.Name,
		})
	}
	if data.ID != "" {
		updates = append(updates, firestore.Update{
			Path:  "id",
			Value: data.ID,
		})
	}

	_, err = svc.db.Collection("patients").Doc(id).Update(ctx, updates)
	if err != nil {
		return gofhir.Patient{}, err
	}

	return data, nil
}

func (svc *Service) delete(ctx context.Context, id string) (err error) {
	_, err = svc.db.Collection("patients").Doc(id).Delete(ctx)
	return err
}
