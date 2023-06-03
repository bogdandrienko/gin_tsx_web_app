package backend

//

// IMPORTS ////////////////////////////////////////////////////////////////////////////////////////////////////////////

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// IMPORTS ////////////////////////////////////////////////////////////////////////////////////////////////////////////

//

// STRUCTS ////////////////////////////////////////////////////////////////////////////////////////////////////////////
// STRUCTS ////////////////////////////////////////////////////////////////////////////////////////////////////////////

//

// METHODS ////////////////////////////////////////////////////////////////////////////////////////////////////////////
// METHODS ////////////////////////////////////////////////////////////////////////////////////////////////////////////

//

// GLOBALS ////////////////////////////////////////////////////////////////////////////////////////////////////////////
// GLOBALS ////////////////////////////////////////////////////////////////////////////////////////////////////////////

//

// ACTIONS ////////////////////////////////////////////////////////////////////////////////////////////////////////////

func IndexHandler(context *gin.Context) {
	context.HTML(http.StatusOK, "index.html", gin.H{
		"title": "SPA",
	})
}

func CreateTaskHandler(context *gin.Context) {
	// get and check param
	title := context.PostForm("title")
	if title == "" {
		ErrHandlerWithContext(errors.New("title incorrect"), context)
		return
	}

	// insert to db
	err := ExecuteInsertOrDeleteDb("insert into tasks (title) values ($1);", title)
	if err != nil {
		ErrHandlerWithContext(err, context)
		return
	}

	context.JSON(http.StatusCreated, map[string]string{"response": "successfully created"})
}

func ReadTaskHandler(context *gin.Context) {
	// get id from url
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		ErrHandlerWithContext(err, context)
		return
	}

	// create object
	task := Task{}

	// select from db
	err = ExecuteSelectOneDb([]any{&task.Id, &task.Title}, "select id, title from tasks where id = $1", id)
	if err != nil {
		ErrHandlerWithContext(err, context)
		return
	}

	context.JSON(http.StatusOK, map[string]any{"response": task})
}

func ReadTasksHandler(context *gin.Context) {
	// select from db
	rows, err := ExecuteRowsDb("select id, title from tasks order by id asc;")
	if err != nil {
		ErrHandlerWithContext(err, context)
		return
	}

	// create objects
	tasks := make([]Task, 0)

	// fulling objects
	for rows.Next() {
		task := Task{}
		err = rows.Scan(&task.Id, &task.Title)
		if err != nil {
			ErrHandlerWithContext(err, context)
			return
		}
		tasks = append(tasks, task)
	}

	context.JSON(http.StatusOK, map[string]map[string]any{"response": {"list": tasks}})
}

func UpdateTaskHandler(context *gin.Context) {
	// get id from url
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		ErrHandlerWithContext(err, context)
		return
	}

	// get and check param
	title := context.PostForm("title")
	if title == "" {
		ErrHandlerWithContext(errors.New("title incorrect"), context)
		return
	}

	// update into db
	err = ExecuteInsertOrDeleteDb("update tasks set title=$1 where id = $2;", title, id)
	if err != nil {
		ErrHandlerWithContext(err, context)
		return
	}

	context.JSON(http.StatusOK, map[string]string{"response": "successfully updated"})
}

func DeleteTaskHandler(context *gin.Context) {
	// get id from url
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		ErrHandlerWithContext(err, context)
		return
	}

	// delete from db
	err = ExecuteInsertOrDeleteDb("delete from tasks where id = $1;", id)
	if err != nil {
		ErrHandlerWithContext(err, context)
		return
	}

	context.JSON(http.StatusOK, map[string]string{"response": "successfully deleted"})
}

// ACTIONS ////////////////////////////////////////////////////////////////////////////////////////////////////////////

//

// EXTRAS /////////////////////////////////////////////////////////////////////////////////////////////////////////////
// EXTRAS /////////////////////////////////////////////////////////////////////////////////////////////////////////////

//
