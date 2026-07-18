package graph

import (
	"context"
	"errors"
	"graphql-server/graph/model"

	"github.com/google/uuid"
)

func (r *queryResolver) Projects(ctx context.Context) ([]*model.Project, error) {
	return r.Store.Projects, nil
}

func (r *queryResolver) Project(ctx context.Context, id string) (*model.Project, error) {
	for _, project := range r.Store.Projects {
		if project.ID == id {
			return project, nil
		}
	}
	return nil, nil
}

func (r *mutationResolver) CreateProject(ctx context.Context, input model.CreateProjectInput) (*model.Project, error) {
	project := &model.Project{
		ID:          uuid.NewString(),
		Name:        input.Name,
		Description: input.Description,
		Tasks:       []*model.Task{},
	}
	r.Store.Projects = append(r.Store.Projects, project)
	return project, nil
}

func (r *mutationResolver) AddTask(ctx context.Context, input model.AddTaskInput) (*model.Task, error) {
	for _, project := range r.Store.Projects {
		if project.ID != input.ProjectID {
			continue
		}

		task := &model.Task{
			ID:          uuid.NewString(),
			Title:       input.Title,
			Description: input.Description,
			Status:      input.Status,
		}
		project.Tasks = append(project.Tasks, task)
		return task, nil
	}

	return nil, errors.New("project not found")
}

func (r *mutationResolver) UpdateTaskStatus(ctx context.Context, taskID string, status model.TaskStatus) (*model.Task, error) {
	for _, project := range r.Store.Projects {
		for _, task := range project.Tasks {
			if task.ID == taskID {
				task.Status = status
				return task, nil
			}
		}
	}

	return nil, errors.New("task not found")
}

func (r *projectResolver) Tasks(ctx context.Context, obj *model.Project, status *model.TaskStatus) ([]*model.Task, error) {
	if status == nil {
		return obj.Tasks, nil
	}

	filtered := make([]*model.Task, 0, len(obj.Tasks))
	for _, task := range obj.Tasks {
		if task.Status == *status {
			filtered = append(filtered, task)
		}
	}

	return filtered, nil
}

func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

func (r *Resolver) Project() ProjectResolver { return &projectResolver{r} }

func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type (
	mutationResolver struct{ *Resolver }
	projectResolver  struct{ *Resolver }
	queryResolver    struct{ *Resolver }
)
