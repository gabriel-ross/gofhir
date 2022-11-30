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
	data.ID = svc.db.Collection("users").NewDoc().ID
	_, err = svc.db.Collection("users").Doc(data.ID).Set(ctx, data)
	if err != nil {
		return gofhir.Patient{}, err
	}
	return data, nil
}

func (svc *Service) list(ctx context.Context, offset, limit int) (_ []gofhir.Patient, err error) {
	resp := []gofhir.Patient{}
	iter := svc.db.Collection("users").OrderBy("id", firestore.Asc).StartAt(offset).Limit(limit).Documents(ctx)
	for {
		dsnap, err := iter.Next()
		if err == iterator.Done {
			break
		}

		var user gofhir.Patient
		dsnap.DataTo(&user)
		resp = append(resp, user)
	}
	return resp, nil
}

func (svc *Service) count(ctx context.Context) (_ int, err error) {
	docs, err := svc.db.Collection("users").Documents(ctx).GetAll()
	if err != nil {
		return 0, err
	}
	return len(docs), nil
}

func (svc *Service) read(ctx context.Context, id string) (_ gofhir.Patient, err error) {
	dsnap, err := svc.db.Collection("users").Doc(id).Get(ctx)
	if err != nil {
		return gofhir.Patient{}, err
	}

	var user gofhir.Patient
	dsnap.DataTo(&user)
	return user, nil
}

func (svc *Service) exists(ctx context.Context, id string) (_ bool, err error) {
	_, err = svc.db.Collection("users").Doc(id).Get(ctx)
	if status.Code(err) == codes.NotFound {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (svc *Service) update(ctx context.Context, id string, data gofhir.Patient) (_ gofhir.Patient, err error) {
	data.ID = id
	_, err = svc.db.Collection("users").Doc(data.ID).Set(ctx, data)
	if err != nil {
		return gofhir.Patient{}, err
	}
	return data, nil
}

func (svc *Service) delete(ctx context.Context, id string) (err error) {
	_, err = svc.db.Collection("users").Doc(id).Delete(ctx)
	return err
}
