package service

import (
	todo "todo-app/app-models"
	"todo-app/pkg/repository"
)

type ToDoListService struct {
	repo repository.ToDoList
}

func NewToDoListService(repo repository.ToDoList) *ToDoListService {
	return &ToDoListService{
		repo: repo,
	}
}

func (s *ToDoListService) Create(userId int, list todo.ToDoList) (int, error) {
	return s.repo.Create(userId, list)
}

func (s *ToDoListService) GetAll(userId int) ([]todo.ToDoList, error) {
	return s.repo.GetAll(userId)
}

func (s *ToDoListService) GetById(userId, listId int) (todo.ToDoList, error) {
	return s.repo.GetById(userId, listId)
}

func (s *ToDoListService) Delete(userId, listId int) error {
	return s.repo.Delete(userId, listId)
}

func (s *ToDoListService) Update(userId, listId int, updateData todo.UpdateListInput) error {
	if err := updateData.Validate(); err != nil {
		return err
	}
	return s.repo.Update(userId, listId, updateData)
}
