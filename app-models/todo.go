package todo

import "errors"

type ToDoList struct {
	Id          int64  `json:"id" db:"id"`
	Title       string `json:"title" db:"title" binding:"required"`
	Description string `json:"description" db:"description"`
}

type UsersList struct {
	Id     int64
	UserId int64
	ListId int64
}

type ToDoItem struct {
	Id          int64  `json:"id"`
	Title       string `json:"title" db:"title" binding:"required"`
	Description string `json:"description" db:"description"`
	Done        bool   `json:"done" db:"done"`
}

type UpdateItemInput struct {
	Title       *string `json:"title" db:"title" binding:"required"`
	Description *string `json:"description" db:"description"`
	Done        *bool   `json:"done" db:"done"`
}

func (u *UpdateItemInput) Validate() error {
	if u.Title == nil && u.Description == nil && u.Done == nil {
		return errors.New("empty request for updating item input")
	}
	return nil
}

type ListsItem struct {
	Id     int64
	ListId int64
	ItemId int64
}

type UpdateListInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

func (u *UpdateListInput) Validate() error {
	if u.Title == nil && u.Description == nil {
		return errors.New("invalid body for update request, empty request")
	}
	return nil
}
