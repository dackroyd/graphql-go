package graphql_test

import (
	"context"
	"testing"

	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/example/social"
	"github.com/graph-gophers/graphql-go/example/starwars"
)

func BenchmarkStarWarsSchema_heroDetails(b *testing.B) {
	schema := graphql.MustParseSchema(starwars.Schema, &starwars.Resolver{}, graphql.UseFieldResolvers())

	queryText := `
		query HeroDetails($episode: Episode!) {
			hero(episode: $episode) {
				name
				friends {
					name
					appearsIn
				}
			}
		}
	`

	variables := map[string]interface{}{
		"episode": "EMPIRE",
	}

	b.Run("Exec", func(b *testing.B) {
		for b.Loop() {
			r := schema.Exec(context.Background(), queryText, "", variables)
			if len(r.Errors) > 0 {
				b.Fatalf("failed to execute query: %v", r.Errors)
			}
		}
	})

	b.Run("PrepareQuery", func(b *testing.B) {
		query := schema.MustPrepareQuery(context.Background(), queryText)

		for b.Loop() {
			r := query.Exec(context.Background(), "", variables)
			if len(r.Errors) > 0 {
				b.Fatalf("failed to execute query: %v", r.Errors)
			}
		}
	})
}

func BenchmarkStarWarsSchema_deepConditionalKeep(b *testing.B) {
	schema := graphql.MustParseSchema(starwars.Schema, &starwars.Resolver{}, graphql.UseFieldResolvers())

	queryText := `
		query HeroDetails($episode: Episode!, $excludeFriends: Boolean!) {
			hero(episode: $episode) {
				name
				friends @skip(if: $excludeFriends) {
					name
					appearsIn
					friends {
						name
						appearsIn
						friends {
							name
							appearsIn
							friends {
								name
								appearsIn
								friends {
									name
									appearsIn
									friends {
										name
										appearsIn
										friends {
											name
											appearsIn
										}
									}
								}
							}
						}
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"episode":        "EMPIRE",
		"excludeFriends": false,
	}

	b.Run("Exec", func(b *testing.B) {
		for b.Loop() {
			r := schema.Exec(context.Background(), queryText, "", variables)
			if len(r.Errors) > 0 {
				b.Fatalf("failed to execute query: %v", r.Errors)
			}
		}
	})

	b.Run("PrepareQuery", func(b *testing.B) {
		query := schema.MustPrepareQuery(context.Background(), queryText)

		for b.Loop() {
			r := query.Exec(context.Background(), "", variables)
			if len(r.Errors) > 0 {
				b.Fatalf("failed to execute query: %v", r.Errors)
			}
		}
	})
}

func BenchmarkStarWarsSchema_heroesConnection(b *testing.B) {
	queryText := `
		query Heroes(
			$heroFriendsAfter: ID
			$moreFriendsAfter: ID
		) {
			hero {
				name
				friendsConnection(first: 1, after: $heroFriendsAfter) {
					totalCount
					pageInfo {
						startCursor
						endCursor
						hasNextPage
					}
					edges {
						cursor
						node {
							name
						}
					}
				}
			},
			moreFriends: hero {
				name
				friendsConnection(first: 1, after: $moreFriendsAfter) {
					totalCount
					pageInfo {
						startCursor
						endCursor
						hasNextPage
					}
					edges {
						cursor
						node {
							name
						}
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"heroFriendsAfter": "Y3Vyc29yMQ==",
		"moreFriendsAfter": "Y3Vyc29yMg==",
	}

	b.Run("Exec", func(b *testing.B) {
		for b.Loop() {
			r := starwarsSchema.Exec(context.Background(), queryText, "", variables)
			if len(r.Errors) > 0 {
				b.Fatalf("failed to execute query: %v", r.Errors)
			}
		}
	})

	b.Run("PrepareQuery", func(b *testing.B) {
		query := starwarsSchema.MustPrepareQuery(context.Background(), queryText)

		for b.Loop() {
			r := query.Exec(context.Background(), "", variables)
			if len(r.Errors) > 0 {
				b.Fatalf("failed to execute query: %v", r.Errors)
			}
		}
	})
}

func BenchmarkStarWarsSchema_deepConditionalDrop(b *testing.B) {
	schema := graphql.MustParseSchema(starwars.Schema, &starwars.Resolver{}, graphql.UseFieldResolvers())

	queryText := `
		query HeroDetails($episode: Episode!, $excludeFriends: Boolean!) {
			hero(episode: $episode) {
				name
				friends @skip(if: $excludeFriends) {
					name
					appearsIn
					friends {
						name
						appearsIn
						friends {
							name
							appearsIn
							friends {
								name
								appearsIn
								friends {
									name
									appearsIn
									friends {
										name
										appearsIn
										friends {
											name
											appearsIn
										}
									}
								}
							}
						}
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"episode":        "EMPIRE",
		"excludeFriends": true,
	}

	b.Run("Exec", func(b *testing.B) {
		for b.Loop() {
			r := schema.Exec(context.Background(), queryText, "", variables)
			if len(r.Errors) > 0 {
				b.Fatalf("failed to execute query: %v", r.Errors)
			}
		}
	})

	b.Run("PrepareQuery", func(b *testing.B) {
		query := schema.MustPrepareQuery(context.Background(), queryText)

		for b.Loop() {
			r := query.Exec(context.Background(), "", variables)
			if len(r.Errors) > 0 {
				b.Fatalf("failed to execute query: %v", r.Errors)
			}
		}
	})
}

// BenchmarkStarWarsSchema_deepAsync benchmarks a query where many fields expect to be executed async, necessetating the
// use of concurrent Go routines to resolve
func BenchmarkStarWarsSchema_deepAsync(b *testing.B) {
	schema := graphql.MustParseSchema(starwars.Schema, &starwars.Resolver{}, graphql.UseFieldResolvers())

	queryText := `
		query HeroDetails($episode: Episode!) {
			hero(episode: $episode) {
				name
				friendsConnection(first: 10) {
					friends {
						name
						appearsIn
						friendsConnection(first: 10) {
							friends {
								name
								appearsIn
								friendsConnection(first: 10) {
									friends {
										name
										appearsIn
										friendsConnection(first: 10) {
											friends {
												name
												appearsIn
												friendsConnection(first: 10) {
													friends {
														name
														appearsIn
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
			}
		}
	`

	variables := map[string]interface{}{
		"episode": "EMPIRE",
	}

	b.Run("Exec", func(b *testing.B) {
		for b.Loop() {
			r := schema.Exec(context.Background(), queryText, "", variables)
			if len(r.Errors) > 0 {
				b.Fatalf("failed to execute query: %v", r.Errors)
			}
		}
	})

	b.Run("PrepareQuery", func(b *testing.B) {
		query := schema.MustPrepareQuery(context.Background(), queryText)

		for b.Loop() {
			r := query.Exec(context.Background(), "", variables)
			if len(r.Errors) > 0 {
				b.Fatalf("failed to execute query: %v", r.Errors)
			}
		}
	})
}

func BenchmarkSocialSchema(b *testing.B) {
	schema := graphql.MustParseSchema(social.Schema, &social.Resolver{}, graphql.UseFieldResolvers())

	queryText := `
		query User($userID: ID!) {
			user(id: $userID) {
				id
				name
				friends {
					id
					name
				}
			}
		}
	`

	variables := map[string]interface{}{
		"userID": "0x02",
	}

	b.Run("Exec", func(b *testing.B) {
		for b.Loop() {
			r := schema.Exec(context.Background(), queryText, "", variables)
			if len(r.Errors) > 0 {
				b.Fatalf("failed to execute query: %v", r.Errors)
			}
		}
	})

	b.Run("PrepareQuery", func(b *testing.B) {
		query := schema.MustPrepareQuery(context.Background(), queryText)

		for b.Loop() {
			r := query.Exec(context.Background(), "", variables)
			if len(r.Errors) > 0 {
				b.Fatalf("failed to execute query: %v", r.Errors)
			}
		}
	})
}
