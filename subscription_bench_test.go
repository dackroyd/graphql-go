package graphql_test

import (
	"context"
	"testing"

	"github.com/graph-gophers/graphql-go"
)

func BenchmarkSubscription_deepConditionalKeep(b *testing.B) {
	rr := &LibraryRootResolver{
		LibraryQueryResolver:        &LibraryQueryResolver{},
		LibrarySubscriptionResolver: &LibrarySubscriptionResolver{},
	}
	schema := graphql.MustParseSchema(complexSchema, rr, graphql.UseFieldResolvers())

	queryText := `
		subscription Books($excludeAuthors: Boolean!) {
			books {
				...BookFragment
				authors @skip(if: $excludeAuthors) {
					...AuthorFragment
					books {
						...BookFragment
						authors {
							...AuthorFragment
							books {
								...BookFragment
								authors {
									...AuthorFragment
									books {
										...BookFragment
										authors {
											...AuthorFragment
											books {
												...BookFragment
											}
										}
									}
								}
							}
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
	`

	variables := map[string]interface{}{
		"excludeAuthors": false,
	}

	b.Run("Subscribe", func(b *testing.B) {
		for b.Loop() {
			c, err := schema.Subscribe(context.Background(), queryText, "", variables)
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
		query := schema.MustPrepareQuery(context.Background(), queryText)

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

func BenchmarkSubscription_deepConditionalDrop(b *testing.B) {
	rr := &LibraryRootResolver{
		LibraryQueryResolver:        &LibraryQueryResolver{},
		LibrarySubscriptionResolver: &LibrarySubscriptionResolver{},
	}
	schema := graphql.MustParseSchema(complexSchema, rr, graphql.UseFieldResolvers())

	queryText := `
		subscription Books($excludeAuthors: Boolean!) {
			books {
				...BookFragment
				authors @skip(if: $excludeAuthors) {
					...AuthorFragment
					books {
						...BookFragment
						authors {
							...AuthorFragment
							books {
								...BookFragment
								authors {
									...AuthorFragment
									books {
										...BookFragment
										authors {
											...AuthorFragment
											books {
												...BookFragment
											}
										}
									}
								}
							}
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
	`

	variables := map[string]interface{}{
		"excludeAuthors": true,
	}

	b.Run("Subscribe", func(b *testing.B) {
		for b.Loop() {
			c, err := schema.Subscribe(context.Background(), queryText, "", variables)
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
		query := schema.MustPrepareQuery(context.Background(), queryText)

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
