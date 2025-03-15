package graphql_test

import (
	"context"
	"slices"
	"testing"

	"github.com/graph-gophers/graphql-go"
)

const (
	complexSchema = `
		schema {
			query: Query
			subscription: Subscription
		}

		type Query {
			authors: [Author!]!
			books: [Book!]!
			series: [Series!]!
		}

		type Subscription {
			books: Book!
		}

		type Book {
			id: ID!
			title: String!
			authors: [Author!]!
			series: Series
		}

		type Series {
			id: ID!
			title: String!
			authors: [Author!]!
			books: [Book!]!
		}

		type Author {
			id: ID!
			name: String!
			books: [Book!]!
			series: [Series!]!
		}
	`
	complexQuery = `
		query Library(
			$excludeSeriesAuthors: Boolean = true,
			$excludeBookAuthors: Boolean = true
		) {
			authors {
				...AuthorFragment
				books {
					...BookFragment
					authors @skip(if: $excludeBookAuthors) {
						...AuthorFragment
						books {
							...BookFragment
						}
					}
				}
				series {
					...SeriesFragment
					authors @skip(if: $excludeSeriesAuthors) {
						...AuthorFragment
					}
					books {
						...BookFragment
						authors @skip(if: $excludeBookAuthors) {
							...AuthorFragment
						}
					}
				}
			}
			books {
				...BookFragment
				authors @skip(if: $excludeBookAuthors) {
					...AuthorFragment
					books {
						...BookFragment
					}
					series {
						...SeriesFragment
						authors @skip(if: $excludeSeriesAuthors) {
							...AuthorFragment
							books {
								...BookFragment
							}
						}
						books {
							...BookFragment
						}
					}
				}
			}
			series {
				...SeriesFragment
				authors @skip(if: $excludeSeriesAuthors) {
					...AuthorFragment
					books {
						...BookFragment
					}
				}
				books {
					...BookFragment
					authors @skip(if: $excludeBookAuthors) {
						...AuthorFragment
						books {
							...BookFragment
						}
					}
				}
			}
		}

		fragment AuthorFragment on Author {
			id
			name
		}

		fragment BookFragment on Book {
			id
			title
		}

		fragment SeriesFragment on Series {
			id
			title
		}
	`
	complexSubscriptionQuery = `
		subscription Library(
			$excludeSeriesAuthors: Boolean!
		){
			books {
				...BookFragment
				authors {
					...AuthorFragment
					books {
						...BookFragment
					}
					series {
						...SeriesFragment
						authors @skip(if: $excludeSeriesAuthors) {
							...AuthorFragment
						}
						books {
							...BookFragment
						}
					}
				}
				series {
					...SeriesFragment
					// authors @skip(if: $excludeSeriesAuthors) {
					// 	...AuthorFragment
					// }
					books {
						...BookFragment
					}
				}
			}
		}

		fragment AuthorFragment on Author {
			id
			name
		}

		fragment BookFragment on Book {
			id
			title
		}

		fragment SeriesFragment on Series {
			id
			title
		}
	`
)

func BenchmarkComplexSchemaWithFragments(b *testing.B) {
	rr := &LibraryRootResolver{
		LibraryQueryResolver:        &LibraryQueryResolver{},
		LibrarySubscriptionResolver: &LibrarySubscriptionResolver{},
	}
	schema := graphql.MustParseSchema(complexSchema, rr, graphql.UseFieldResolvers())

	b.Run("Exec", func(b *testing.B) {
		for b.Loop() {
			r := schema.Exec(context.Background(), complexQuery, "", nil)
			if len(r.Errors) > 0 {
				b.Fatalf("failed to execute query: %v", r.Errors)
			}
		}
	})

	b.Run("PrepareQuery", func(b *testing.B) {
		query := schema.MustPrepareQuery(context.Background(), complexQuery)

		for b.Loop() {
			r := query.Exec(context.Background(), "", nil)
			if len(r.Errors) > 0 {
				b.Fatalf("failed to execute query: %v", r.Errors)
			}
		}
	})
}

func BenchmarkSubscriptionComplexSchemaWithFragments(b *testing.B) {
	rr := &LibraryRootResolver{
		LibraryQueryResolver:        &LibraryQueryResolver{},
		LibrarySubscriptionResolver: &LibrarySubscriptionResolver{},
	}
	schema := graphql.MustParseSchema(complexSchema, rr, graphql.UseFieldResolvers())

	variables := map[string]interface{}{
		"excludeSeriesAuthors": true,
	}

	b.Run("Subscribe", func(b *testing.B) {
		for b.Loop() {
			c, err := schema.Subscribe(context.Background(), complexSubscriptionQuery, "", variables)
			if err != nil {
				b.Fatalf("failed to execute query: %v", err)
			}

			for cc := range c {
				r := cc.(*graphql.Response)

				if len(r.Errors) > 0 {
					b.Fatalf("failed to execute query: %v", r.Errors)
				}
			}
		}
	})

	b.Run("PrepareQuery", func(b *testing.B) {
		query := schema.MustPrepareQuery(context.Background(), complexSubscriptionQuery)

		for b.Loop() {
			c, err := query.Subscribe(context.Background(), "", variables)
			if err != nil {
				b.Fatalf("failed to execute query: %v", err)
			}

			for cc := range c {
				r := cc.(*graphql.Response)

				if len(r.Errors) > 0 {
					b.Fatalf("failed to execute query: %v", r.Errors)
				}
			}
		}
	})
}

type Book struct {
	Title   string
	id      graphql.ID
	authors []graphql.ID
	series  *graphql.ID
}

func (b *Book) ID() graphql.ID {
	return b.id
}

func (b *Book) Authors() []*Author {
	authors := make([]*Author, len(b.authors))
	for i, id := range b.authors {
		authors[i] = authorsData[id]
	}

	return authors
}

func (b *Book) Series() *Series {
	if b.series == nil {
		return nil
	}

	return seriesData[*b.series]
}

type Author struct {
	Name   string
	id     graphql.ID
	books  []graphql.ID
	series []graphql.ID
}

func (a *Author) ID() graphql.ID {
	return a.id
}

func (a *Author) Books() []*Book {
	books := make([]*Book, len(a.books))
	for i, id := range a.books {
		books[i] = booksData[id]
	}

	return books
}

func (a *Author) Series() []*Series {
	series := make([]*Series, len(a.series))
	for i, id := range a.series {
		series[i] = seriesData[id]
	}

	return series
}

type Series struct {
	Title string
	id    graphql.ID
	books []graphql.ID
}

func (s *Series) ID() graphql.ID {
	return s.id
}

func (s *Series) Books() []*Book {
	books := make([]*Book, len(s.books))
	for i, id := range s.books {
		books[i] = booksData[id]
	}

	return books
}

func (s *Series) Authors() []*Author {
	authorSet := make(map[graphql.ID]struct{})
	for _, bID := range s.books {
		b := booksData[bID]
		for _, aID := range b.authors {
			authorSet[aID] = struct{}{}
		}
	}

	authors := make([]*Author, 0, len(authorSet))
	for aID := range authorSet {
		authors = append(authors, authorsData[aID])
	}

	return authors
}

type LibraryRootResolver struct {
	*LibraryQueryResolver
	*LibrarySubscriptionResolver
}

func (r *LibraryRootResolver) Query() *LibraryQueryResolver {
	return r.LibraryQueryResolver
}

func (r *LibraryRootResolver) Subscription() *LibrarySubscriptionResolver {
	return r.LibrarySubscriptionResolver
}

var authorsData = make(map[graphql.ID]*Author)
var booksData = make(map[graphql.ID]*Book)
var seriesData = make(map[graphql.ID]*Series)

func init() {
	authors := []*Author{
		{
			id:   "2000",
			Name: "Author 1",
		},
		{
			id:   "2001",
			Name: "Author 2",
		},
		{
			id:   "2002",
			Name: "Author 3",
		},
	}
	books := []*Book{
		{
			id:      "1000",
			Title:   "Book 1",
			authors: []graphql.ID{"2000", "2001"},
		},
		{
			id:      "1001",
			Title:   "Book 2",
			authors: []graphql.ID{"2001"},
		},
		{
			id:      "1002",
			Title:   "Book 3",
			authors: []graphql.ID{"2000", "2002"},
		},
		{
			id:      "1003",
			Title:   "Book 4",
			authors: []graphql.ID{"2002"},
		},
	}
	series := []*Series{
		{
			id:    "3000",
			Title: "Series 1",
			books: []graphql.ID{"1000", "1001"},
		},
		{
			id:    "3001",
			Title: "Series 2",
			books: []graphql.ID{"1002", "1003"},
		},
	}

	for _, a := range authors {
		authorsData[a.id] = a
	}

	for _, b := range books {
		booksData[b.id] = b

		for _, aID := range b.authors {
			a := authorsData[aID]
			a.books = append(a.books, b.id)
		}
	}

	for _, s := range series {
		seriesData[s.id] = s

		for _, bID := range s.books {
			b := booksData[bID]
			b.series = &s.id

			for _, aID := range b.authors {
				a := authorsData[aID]

				if !slices.Contains(a.series, s.id) {
					a.series = append(a.series, s.id)
				}
			}
		}
	}
}

type LibraryQueryResolver struct{}

func (r *LibraryQueryResolver) Books() []*Book {
	books := make([]*Book, 0, len(booksData))

	for _, b := range booksData {
		books = append(books, b)
	}

	return books
}

func (r *LibraryQueryResolver) Authors() []*Author {
	authors := make([]*Author, 0, len(authorsData))

	for _, a := range authorsData {
		authors = append(authors, a)
	}

	return authors
}

func (r *LibraryQueryResolver) Series() []*Series {
	series := make([]*Series, 0, len(seriesData))

	for _, s := range seriesData {
		series = append(series, s)
	}

	return series
}

type LibrarySubscriptionResolver struct{}

func (r *LibrarySubscriptionResolver) Books() <-chan *Book {
	c := make(chan *Book)

	go func() {
		for _, b := range booksData {
			c <- b
		}

		close(c)
	}()

	return c
}
