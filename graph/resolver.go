package graph

import "graphql-server/graph/model"

type Store struct {
	Projects []*model.Project
}

type Resolver struct {
	Store *Store
}

func SeedStore() *Store {
	return &Store{Projects: []*model.Project{}}
}
