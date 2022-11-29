package gofhir

import (
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Filter operators
var (
	Eq  operator = "=="
	Ne           = "!="
	Lt           = "<"
	Leq          = "<="
	Gt           = ">"
	Geq          = ">="
)

func NewFirestoreClient(ctx context.Context, projectID string) (_ *firestore.Client, err error) {
	return firestore.NewClient(ctx, projectID)
}

func Create[T any](ctx context.Context, client *firestore.Client, collectionPath string, data T, options ...CreateOption) (_ string, _ T, err error) {
	opts := &createOptions{}
	for _, option := range options {
		option(opts)
	}

	var id string
	if opts.idAutoGenerated == true {
		id = client.Collection(collectionPath).NewDoc().ID
	} else {
		id = opts.id
	}
	_, err = client.Collection(collectionPath).Doc(id).Set(ctx, data)
	if err != nil {
		return "", data, err
	}
	return id, data, nil
}

// Exists returns whether a document with the given id exists in the given
// collectionPath.
func Exists(ctx context.Context, client *firestore.Client, collectionPath string, id string) (_ bool, err error) {
	_, err = client.Collection(collectionPath).Doc(id).Get(ctx)
	if status.Code(err) == codes.NotFound {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func Read[T any](ctx context.Context, client *firestore.Client, collectionPath string, id string) (_ T, err error) {
	var data T
	dsnap, err := client.Collection(collectionPath).Doc(id).Get(ctx)
	if err != nil {
		return data, err
	}
	dsnap.DataTo(&data)
	return data, nil
}

func List[T any](ctx context.Context, client *firestore.Client, collectionPath string, options ...ListOption) (_ []T, err error) {
	resp := []T{}
	query := BuildListQueryFromListOptions(client.Collection(collectionPath).Query, options...)
	iter := query.Documents(ctx)
	for {
		dsnap, done := iter.Next()
		if done == iterator.Done {
			break
		}
		var data T
		dsnap.DataTo(&data)
		resp = append(resp, data)
	}
	return resp, nil
}

// Set sets the document at a given id to the given data. If no document
// exists at this id a new one will be created.
func Set[T any](ctx context.Context, client *firestore.Client, collectionPath string, id string, data T, options ...firestore.SetOption) (_ T, err error) {
	_, err = client.Collection(collectionPath).Doc(id).Set(ctx, data, options...)
	return data, err
}

func Delete(ctx context.Context, client *firestore.Client, collectionPath, id string) (err error) {
	_, err = client.Collection(collectionPath).Doc(id).Delete(ctx)
	return err
}

type CreateOption func(*createOptions)

type createOptions struct {
	idAutoGenerated bool
	id              string
}

func WithAutoGeneratedID() CreateOption {
	return func(opts *createOptions) {
		opts.idAutoGenerated = true
		opts.id = ""
	}
}

func WithID(id string) CreateOption {
	return func(opts *createOptions) {
		opts.idAutoGenerated = false
		opts.id = id
	}
}

type operator string

// BuildListQueryFromListOptions builds a firestore.Query from a variadic
// number of ListOption.
// TODO: Enable and get working adding sort keys, offset, and limiting.
// TODO: Error on limitations that a sort key can only follow a filter on the same key
// or something like that
func BuildListQueryFromListOptions(query firestore.Query, options ...ListOption) firestore.Query {

	opts := newListOptions(options...)

	// Add filter params
	for _, filterParam := range opts.filterParameters {
		query = query.Where(filterParam.key, string(filterParam.op), filterParam.val)
	}

	// Add sort keys
	// for _, sortParam := range opts.sortParameters {
	// 	query = query.OrderBy(sortParam.key, sortParam.direction)
	// }

	// if opts.offset > -1 {
	// 	query = query.StartAt(opts.offset)
	// }
	// if opts.limit > -1 {
	// 	query = query.Limit(opts.limit)
	// }

	return query
}

// See https://cloud.google.com/firestore/docs/query-data/order-limit-data
// for limitations on filtering and ordering data
type ListOption func(*ListOptions)

type ListOptions struct {
	limit            int
	offset           int
	sortParameters   []sortParameter
	filterParameters []filterParameter
}

func newListOptions(opts ...ListOption) *ListOptions {
	lo := &ListOptions{
		limit:            100,
		offset:           0,
		sortParameters:   []sortParameter{},
		filterParameters: []filterParameter{},
	}
	for _, opt := range opts {
		opt(lo)
	}
	return lo
}

func WithLimit(limit int) ListOption {
	return func(l *ListOptions) {
		l.limit = limit
	}
}

func WithOffset(offset int) ListOption {
	return func(l *ListOptions) {
		l.offset = offset
	}
}

// WithSortKey adds a sort key to the sort chain. Sort keys will be applied in
// the order they are added. Max of two sort keys - any more will be ignored.
func WithSortKey(key string, direction firestore.Direction) ListOption {
	return func(l *ListOptions) {
		if len(l.sortParameters) < 2 {
			l.sortParameters = append(l.sortParameters, sortParameter{
				key:       key,
				direction: direction,
			})
		}
	}
}

// WithQuery adds a filter query to the query chain. Calling this multiple times
// adds additional filter queries to the chain rather than overwriting previous
// filters.
func WithFilterQuery(key string, op operator, val interface{}) func(*ListOptions) {
	return func(l *ListOptions) {
		l.filterParameters = append(l.filterParameters, filterParameter{
			key: key,
			op:  op,
			val: val,
		})
	}
}

type filterParameter struct {
	key string
	op  operator
	val interface{}
}

type sortParameter struct {
	key       string
	direction firestore.Direction
}