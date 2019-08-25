package grimoire

type Querier interface {
	Build(*Query)
}

func BuildQuery(collection string, queriers ...Querier) Query {
	q := Query{
		empty: true,
	}

	for _, querier := range queriers {
		querier.Build(&q)
		q.empty = false
	}

	if q.Collection == "" {
		q.Collection = collection
		q.empty = false
	}

	for i := range q.JoinQuery {
		q.JoinQuery[i].buildJoin(q)
	}

	return q
}

// Query defines information about query generated by query builder.
type Query struct {
	empty       bool // todo: use bit to mark what is updated and use it when building
	Collection  string
	SelectQuery SelectQuery
	JoinQuery   []JoinQuery
	WhereQuery FilterQuery
	GroupQuery  GroupQuery
	SortQuery   []SortQuery
	OffsetQuery Offset
	LimitQuery  Limit
	LockQuery  Lock
}

func (q Query) Build(query *Query) {
	if query.empty {
		*query = q
	} else {
		// manual merge
		if q.Collection != "" {
			query.Collection = q.Collection
		}

		if q.SelectQuery.Fields != nil {
			query.SelectQuery = q.SelectQuery
		}

		query.JoinQuery = append(query.JoinQuery, q.JoinQuery...)

		query.WhereQuery = query.WhereQuery.And(q.WhereQuery)

		if q.GroupQuery.Fields != nil {
			query.GroupQuery = q.GroupQuery
		}

		q.SortQuery = append(q.SortQuery, query.SortQuery...)

		if q.OffsetQuery != 0 {
			query.OffsetQuery = q.OffsetQuery
		}

		if q.LimitQuery != 0 {
			query.LimitQuery = q.LimitQuery
		}

		if q.LockQuery != "" {
			query.LockQuery = q.LockQuery
		}
	}
}

// Select filter fields to be selected from database.
func (q Query) Select(fields ...string) Query {
	q.SelectQuery = NewSelect(fields...)
	return q
}

func (q Query) From(collection string) Query {
	q.Collection = collection

	// TODO: no select default to select all

	// if len(q.SelectQuery.Fields) == 0 {
	// 	q.SelectQuery = NewSelect(collection + ".*")
	// }

	// if len(q.JoinQuery) > 0 {
	// 	for i := range q.JoinQuery {

	// 	}
	// }

	return q
}

func (q Query) Distinct() Query {
	q.SelectQuery.OnlyDistinct = true
	return q
}

// Join current collection with other collection.
func (q Query) Join(collection string) Query {
	return q.JoinOn(collection, "", "")
}

// Join current collection with other collection.
func (q Query) JoinOn(collection string, from string, to string) Query {
	return q.JoinWith("JOIN", collection, from, to)
}

// JoinWith current collection with other collection with custom join mode.
func (q Query) JoinWith(mode string, collection string, from string, to string) Query {
	NewJoinWith(mode, collection, from, to).Build(&q) // TODO: ensure this always called last

	return q
}

func (q Query) JoinFragment(expr string, args ...interface{}) Query {
	NewJoinFragment(expr, args...).Build(&q) // TODO: ensure this always called last

	return q
}

func (q Query) Where(filters ...FilterQuery) Query {
	q.WhereQuery = q.WhereQuery.And(filters...)
	return q
}

func (q Query) OrWhere(filters ...FilterQuery) Query {
	q.WhereQuery = q.WhereQuery.Or(And(filters...))
	return q
}

func (q Query) Group(fields ...string) Query {
	q.GroupQuery.Fields = fields
	return q
}

func (q Query) Having(filters ...FilterQuery) Query {
	q.GroupQuery.Filter = q.GroupQuery.Filter.And(filters...)
	return q
}

func (q Query) OrHaving(filters ...FilterQuery) Query {
	q.GroupQuery.Filter = q.GroupQuery.Filter.Or(And(filters...))
	return q
}

func (q Query) Sort(fields ...string) Query {
	return q.SortAsc(fields...)
}

func (q Query) SortAsc(fields ...string) Query {
	sorts := make([]SortQuery, len(fields))
	for i := range fields {
		sorts[i] = NewSortAsc(fields[i])
	}

	q.SortQuery = append(q.SortQuery, sorts...)
	return q
}

func (q Query) SortDesc(fields ...string) Query {
	sorts := make([]SortQuery, len(fields))
	for i := range fields {
		sorts[i] = NewSortDesc(fields[i])
	}

	q.SortQuery = append(q.SortQuery, sorts...)
	return q
}

// Offset the result returned by database.
func (q Query) Offset(offset Offset) Query {
	q.OffsetQuery = offset
	return q
}

// Limit result returned by database.
func (q Query) Limit(limit Limit) Query {
	q.LimitQuery = limit
	return q
}

// Lock query expression.
func (q Query) Lock(lock Lock) Query {
	q.LockQuery = lock
	return q
}

// From create query for collection.
func From(collection string) Query {
	return Query{
		Collection: collection,
	}
}

// Join current collection with other collection.
func Join(collection string) Query {
	return JoinOn(collection, "", "")
}

// JoinOn current collection with other collection.
func JoinOn(collection string, from string, to string) Query {
	return JoinWith("JOIN", collection, from, to)
}

// JoinWith current collection with other collection with custom join mode.
func JoinWith(mode string, collection string, from string, to string) Query {
	return Query{
		JoinQuery: []JoinQuery{
			NewJoinWith(mode, collection, from, to),
		},
	}
	// var q Query
	// NewJoinWith(mode, collection, from, to).Build(&q) // TODO: ensure this always called last

	// return q
}

func JoinFragment(expr string, args ...interface{}) Query {
	return Query{
		JoinQuery: []JoinQuery{
			NewJoinFragment(expr, args...),
		},
	}

	// var q Query
	// NewJoinFragment(expr, args...).Build(&q) // TODO: ensure this always called last

	// return q
}

func Where(filters ...FilterQuery) Query {
	return Query{
		WhereQuery: And(filters...),
	}
}

func Group(fields ...string) Query {
	return Query{
		GroupQuery: NewGroup(fields...),
	}
}
