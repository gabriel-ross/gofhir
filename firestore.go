package gofhir

import (
	"context"

	"cloud.google.com/go/firestore"
)

func NewFirestoreClient(ctx context.Context, projectID string) (_ *firestore.Client, err error) {
	return firestore.NewClient(ctx, projectID)
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

func BuildListQuery(query firestore.Query, opts ListOptions) firestore.Query {

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

func NewListOptions(opts ...func(*ListOptions)) *ListOptions {
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
func WithQuery(key string, op operator, val interface{}) func(*ListOptions) {
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
