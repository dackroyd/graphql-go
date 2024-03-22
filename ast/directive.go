package ast

import "github.com/graph-gophers/graphql-go/errors"

// Directive is a representation of the GraphQL Directive.
//
// http://spec.graphql.org/draft/#sec-Language.Directives
type Directive struct {
	Name      Ident
	Arguments ArgumentList
}

// DirectiveDefinition is a representation of the GraphQL DirectiveDefinition.
//
// http://spec.graphql.org/draft/#sec-Type-System.Directives
type DirectiveDefinition struct {
	Name       string
	Desc       string
	Repeatable bool
	Locations  []string
	Arguments  ArgumentsDefinition
	Loc        errors.Location
}

type DirectiveList []*Directive

// Returns the Directive in the DirectiveList by name or nil if not found.
func (l DirectiveList) Get(name string) *Directive {
	for _, d := range l {
		if d.Name.Name == name {
			return d
		}
	}
	return nil
}

// DirectiveLocation represents the locations within a GraphQL schema that a Directive can be declared upon.
//
// https://spec.graphql.org/draft/#sec-Type-System.Directives
type DirectiveLocation interface {
	directiveLocation()
}

// ExecutableDirectiveLocation represents the locations within an executable request that a GraphQL Directive can be declared upon.
//
// https://spec.graphql.org/draft/#sec-Type-System.Directives
type ExecutableDirectiveLocation string

func (ExecutableDirectiveLocation) directiveLocation() {}

const (
	// DirectiveLocationQuery adjacent to a query operation.
	DirectiveLocationQuery ExecutableDirectiveLocation = "QUERY"
	// DirectiveLocationMutation adjacent to a mutation operation.
	DirectiveLocationMutation ExecutableDirectiveLocation = "MUTATION"
	// DirectiveLocationSubscription adjacent to a subscription operation.
	DirectiveLocationSubscription ExecutableDirectiveLocation = "SUBSCRIPTION"
	// DirectiveLocationField adjacent to a field.
	DirectiveLocationField ExecutableDirectiveLocation = "FIELD"
	// DirectiveLocationFragmentDefinition adjacent to a fragment definition.
	DirectiveLocationFragmentDefinition ExecutableDirectiveLocation = "FRAGMENT_DEFINITION"
	// DirectiveLocationFragmentSpread adjacent to a fragment spread.
	DirectiveLocationFragmentSpread ExecutableDirectiveLocation = "FRAGMENT_SPREAD"
	// DirectiveLocationInlineFragment adjacent to an inline fragment.
	DirectiveLocationInlineFragment ExecutableDirectiveLocation = "INLINE_FRAGMENT"

	// FIXME: new type? Need to add elsewhere in the library... on the meta schema, in __DirectiveLocation ... ?
	// DirectiveLocationVariableDefinition adjacent to a variable definition.
	DirectiveLocationVariableDefinition ExecutableDirectiveLocation = "VARIABLE_DEFINITION"
)

// TypeSystemDirectiveLocation represents the locations a GraphQL Directive can be declared upon.
//
// https://spec.graphql.org/draft/#sec-Type-System.Directives
type TypeSystemDirectiveLocation string

func (TypeSystemDirectiveLocation) directiveLocation() {}

const (
	// DirectiveLocationSchema adjacent to a schema definition.
	DirectiveLocationSchema TypeSystemDirectiveLocation = "SCHEMA"
	// DirectiveLocationScalar adjacent to a scalar definition.
	DirectiveLocationScalar TypeSystemDirectiveLocation = "SCALAR"
	// DirectiveLocationObject adjacent to an object type definition.
	DirectiveLocationObject TypeSystemDirectiveLocation = "OBJECT"
	// DirectiveLocationFieldDefinition adjacent to a field definition.
	DirectiveLocationFieldDefinition TypeSystemDirectiveLocation = "FIELD_DEFINITION"
	// DirectiveLocationArgumentDefinition adjacent to an argument definition.
	DirectiveLocationArgumentDefinition TypeSystemDirectiveLocation = "ARGUMENT_DEFINITION"
	// DirectiveLocationInterface adjacent to an interface definition.
	DirectiveLocationInterface TypeSystemDirectiveLocation = "INTERFACE"
	// DirectiveLocationUnion adjacent to a union definition.
	DirectiveLocationUnion TypeSystemDirectiveLocation = "UNION"
	// DirectiveLocationEnum adjacent to an enum definition.
	DirectiveLocationEnum TypeSystemDirectiveLocation = "ENUM"
	// DirectiveLocationEnumValue adjacent to an enum value definition.
	DirectiveLocationEnumValue TypeSystemDirectiveLocation = "ENUM_VALUE"
	// DirectiveLocationInputObject adjacent to an input object type definition.
	DirectiveLocationInputObject TypeSystemDirectiveLocation = "INPUT_OBJECT"
	// DirectiveLocationInputFieldDefiniton adjacent to an input object field definition.
	DirectiveLocationInputFieldDefiniton TypeSystemDirectiveLocation = "INPUT_FIELD_DEFINITION"
)
