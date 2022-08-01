package graph

import "github.xiaoliang.graphql.users/tasks"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Tasks *tasks.Tasks
	Atts  *tasks.Attachments
}
