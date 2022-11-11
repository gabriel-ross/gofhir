package gofhir

import (
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type FirestoreClient struct {
	*firestore.Client
	ctx context.Context
}

func NewFirestoreClient(ctx context.Context, projectID string) (_ *FirestoreClient, err error) {
	fs, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return &FirestoreClient{
		Client: fs,
		ctx:    ctx,
	}, nil
}

// Set sets the document at a given id to the given data. If no document
// exists at this id a new one will be created.
func (fc *FirestoreClient) Set(ctx context.Context, collectionPath string, id string, data interface{}, options ...firestore.SetOption) (_ interface{}, err error) {
	_, err = fc.Client.Collection(collectionPath).Doc(id).Set(ctx, data, options)
	return data, err
}

func (fc *FirestoreClient) List(ctx context.Context, collectionPath string, options ...func(*ListOptions)) (_ []interface{}, err error) {
	resp := []interface{}{}
	query := buildListQueryFromListOptions(fc.Client.Collection(collectionPath).Query, *newListOptions(options...))
	iter := query.Documents(ctx)
	for {
		dsnap, done := iter.Next()
		if done == iterator.Done {
			break
		}
		var data interface{}
		dsnap.DataTo(&data)
		resp = append(resp, data)
	}
	return resp, nil
}

func (fc *FirestoreClient) Read(ctx context.Context, collectionPath string, id string) (_ interface{}, err error) {
	dsnap, err := fc.Client.Collection(collectionPath).Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}
	var data interface{}
	dsnap.DataTo(&data)
	return data, nil
}

func (fc *FirestoreClient) Delete(ctx context.Context, collectionPath, id string) (err error) {
	_, err = fc.Client.Collection(collectionPath).Doc(id).Delete(ctx)
	return err
}

// Filter operators
var (
	Eq  operator = "=="
	Ne           = "!="
	Lt           = "<"
	Leq          = "<="
	Gt           = ">"
	Geq          = ">="
)

type operator string

func buildListQueryFromListOptions(query firestore.Query, opts ListOptions) firestore.Query {

	// Add sort keys
	for _, sortParam := range opts.sortParameters {
		query = query.OrderBy(sortParam.key, sortParam.direction)
	}

	// Add filter params
	for _, filterParam := range opts.filterParameters {
		query = query.Where(filterParam.key, string(filterParam.op), filterParam.val)
	}

	query = query.StartAt(opts.offset).Limit(opts.limit)

	return query
}

type ListOptions struct {
	limit            int
	offset           int
	sortParameters   []sortParameter
	filterParameters []filterParameter
}

func newListOptions(opts ...func(*ListOptions)) *ListOptions {
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

func WithLimit(limit int) func(*ListOptions) {
	return func(l *ListOptions) {
		l.limit = limit
	}
}

func WithOffset(offset int) func(*ListOptions) {
	return func(l *ListOptions) {
		l.offset = offset
	}
}

// WithSortKey adds a sort key to the sort chain. Sort keys will be applied in
// the order they are added. Max of two sort keys - any more will be ignored.
func WithSortKey(key string, direction firestore.Direction) func(*ListOptions) {
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
