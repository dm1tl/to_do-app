package repository

import (
	"context"
	"fmt"
	"strings"
	todo "todo-app/app-models"

	"github.com/jmoiron/sqlx"
)

type ToDoListPostgres struct {
	db *sqlx.DB
}

func NewToDoListPostgres(db *sqlx.DB) *ToDoListPostgres {
	return &ToDoListPostgres{
		db: db,
	}
}

func (r *ToDoListPostgres) Create(ctx context.Context, userId int64, list todo.ToDoList) (int64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	var id int64
	createListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES($1, $2) RETURNING id", todoListTable)
	row := tx.QueryRowContext(ctx, createListQuery, list.Title, list.Description)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}
	createUsersListQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES($1, $2)", usersListsTable)

	if _, err := tx.ExecContext(ctx, createUsersListQuery, userId, id); err != nil {
		tx.Rollback()
		return 0, err
	}
	return id, tx.Commit()
}

func (r *ToDoListPostgres) GetAll(ctx context.Context, userId int64) ([]todo.ToDoList, error) {
	var lists []todo.ToDoList
	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul ON tl.id = ul.list_id WHERE ul.user_id = $1", todoListTable, usersListsTable)
	err := r.db.SelectContext(ctx, &lists, query, userId)
	return lists, err
}

func (r *ToDoListPostgres) GetById(ctx context.Context, userId, listId int64) (todo.ToDoList, error) {
	var list todo.ToDoList
	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul ON tl.id = ul.list_id WHERE ul.user_id = $1 AND ul.list_id = $2",
		todoListTable, usersListsTable)
	err := r.db.GetContext(ctx, &list, query, userId, listId)
	return list, err
}

func (r *ToDoListPostgres) Delete(ctx context.Context, userId, listId int64) error {
	query := fmt.Sprintf("DELETE FROM %s tl using %s ul WHERE tl.id = ul.list_id AND ul.user_id = $1 AND ul.list_id = $2", todoListTable, usersListsTable)
	_, err := r.db.ExecContext(ctx, query, userId, listId)
	return err
}

func (r *ToDoListPostgres) Update(ctx context.Context, userId, listId int64, updateData todo.UpdateListInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1
	if updateData.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title = $%d", argId))
		args = append(args, updateData.Title)
		argId++
	}
	if updateData.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description = $%d", argId))
		args = append(args, updateData.Description)
		argId++
	}

	queryValues := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s tl SET %s FROM %s ul WHERE tl.id = ul.list_id AND list_id = %d AND user_id = %d", todoListTable,
		queryValues, usersListsTable, argId, argId+1)
	args = append(args, listId, userId)
	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}
