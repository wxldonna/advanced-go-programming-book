package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.xiaoliang.graphql.users/graph/generated"
	"github.xiaoliang.graphql.users/graph/model"
)

// CreateTask is the resolver for the createTask field.
func (r *mutationResolver) CreateTask(ctx context.Context, input model.NewTask) (*model.Task, error) {
	at := r.Atts.Create(input.Attchment.Name, input.Attchment.Date, input.Attchment.Contents)
	task := r.Tasks.CreateTask(input.Text, nil, input.Due, r.Atts.GetAttchment(at.ID))
	return task, nil
}

// CreateAttachment is the resolver for the createAttachment field.
func (r *mutationResolver) CreateAttachment(ctx context.Context, input *model.NewAttachment) (*model.Attachment, error) {
	at := r.Atts.Create(input.Name, input.Date, input.Contents)
	return at, nil
}

// GetTasks is the resolver for the getTasks field.
func (r *queryResolver) GetTasks(ctx context.Context) ([]*model.Task, error) {
	return r.Tasks.GetAllTasks(), nil
}

// GetAttachments is the resolver for the getAttachments field.
func (r *queryResolver) GetAttachments(ctx context.Context) ([]*model.Attachment, error) {
	return r.Atts.GetAllAtts(), nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
