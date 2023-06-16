package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyGetTasksHandlerShowCompleteTrue(t *testing.T) {
	var taskList TaskList
	req := httptest.NewRequest(http.MethodGet, "/tasks?showCompleted=true", nil)
	w := httptest.NewRecorder()
	taskList.TasksHandler(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)

	assert.Nil(t, err) //check that error is nil since thats what we expect
	expectedRespBody := "Getting all tasks...\nThere are no tasks!"
	assert.Equal(t, string(data), expectedRespBody)
	assert.Equal(t, http.StatusOK, res.StatusCode)

}

func TestEmptyGetTasksHandlerShowCompleteFalse(t *testing.T) {
	var taskList TaskList
	req := httptest.NewRequest(http.MethodGet, "/tasks?showCompleted=false", nil)
	w := httptest.NewRecorder()
	taskList.TasksHandler(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)

	assert.Nil(t, err) //check that error is nil since thats what we expect
	expectedRespBody := "Getting all tasks...\nThere are no tasks!"
	assert.Equal(t, string(data), expectedRespBody)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestEmptyGetTasksHandler(t *testing.T) {
	var taskList TaskList
	cases := []struct {
		endP, params, want string
		respCode           int
	}{
		{"/tasks?", "showCompleted=false", "Getting all tasks...\nThere are no tasks!", http.StatusOK},
		{"/tasks?", "showCompleted=true", "Getting all tasks...\nThere are no tasks!", http.StatusOK},
	}

	for _, c := range cases {
		requestStr := c.endP + c.params
		req := httptest.NewRequest(http.MethodGet, requestStr, nil)
		w := httptest.NewRecorder()
		taskList.TasksHandler(w, req)
		res := w.Result()
		defer res.Body.Close()
		data, err := ioutil.ReadAll(res.Body)

		assert.Nil(t, err)
		assert.Equal(t, string(data), c.want)
		assert.Equal(t, c.respCode, res.StatusCode)
	}

}

func TestAddCompleteGetTaskCorrect(t *testing.T) {
	var taskList TaskList
	cases := []struct {
		url, op, id, title, description, completed, want string
		respCode                                         int
	}{
		{"/tasks", "/add", "2", "secondTask", "boo2", "false", "Adding the following task to your task list\nTask:\n\tId = 2\n\tTitle = secondTask\n\tDescription = boo2\n\tCompleted = false", http.StatusCreated},
		{"/tasks", "/complete", "2", "", "", "", "Completed task with id 2\n", http.StatusOK},
		{"/tasks", "", "", "", "", "true", "Getting all tasks...\nTask:\n\tId = 2\n\tTitle = secondTask\n\tDescription = boo2\n\tCompleted = true\n", http.StatusOK},
	}

	for _, c := range cases {

		switch c.op {
		case "": //get tasks (GET)
			url := c.url + c.op + "?showCompleted=" + c.completed
			req := httptest.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()
			taskList.TasksHandler(w, req)
			res := w.Result()
			defer res.Body.Close()
			data, err := ioutil.ReadAll(res.Body)
			assert.Nil(t, err)
			assert.Equal(t, c.want, string(data))
			assert.Equal(t, c.respCode, res.StatusCode)
		case "/add": //add a task (POST)
			url := c.url + c.op
			completedBool, _ := strconv.ParseBool(c.completed)
			idInt, _ := strconv.ParseInt(c.id, 10, 64)
			values := map[string]interface{}{"id": idInt, "title": c.title, "description": c.description, "completed": completedBool}
			jsonValue, _ := json.Marshal(values)
			req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonValue))
			w := httptest.NewRecorder()
			taskList.AddTaskHandler(w, req)
			res := w.Result()
			defer res.Body.Close()
			data, err := ioutil.ReadAll(res.Body)
			assert.Nil(t, err)
			assert.Equal(t, c.want, string(data))
			assert.Equal(t, c.respCode, res.StatusCode)

		case "/complete": //complete a task (PATCH)
			url := c.url + c.op
			idInt, _ := strconv.ParseInt(c.id, 10, 64)
			values := map[string]interface{}{"id": idInt}
			jsonValue, _ := json.Marshal(values)
			req := httptest.NewRequest(http.MethodPatch, url, bytes.NewBuffer(jsonValue))
			w := httptest.NewRecorder()
			taskList.CompleteTaskHandler(w, req)
			res := w.Result()
			defer res.Body.Close()
			data, err := ioutil.ReadAll(res.Body)
			assert.Nil(t, err)
			assert.Equal(t, c.want, string(data))
			assert.Equal(t, c.respCode, res.StatusCode)

		case "/clear": //clear tasks (DELETE)
			url := c.url
			req := httptest.NewRequest(http.MethodDelete, url, nil)
			w := httptest.NewRecorder()
			taskList.TasksHandler(w, req)
			res := w.Result()
			defer res.Body.Close()
			data, err := ioutil.ReadAll(res.Body)
			assert.Nil(t, err)
			assert.Equal(t, c.want, string(data))
			assert.Equal(t, c.respCode, res.StatusCode)

		}
	}

}

func TestAddManyCompleteGetTaskCorrect(t *testing.T) {
	var taskList TaskList
	cases := []struct {
		url, op, id, title, description, completed, want string
		respCode                                         int
	}{
		{"/tasks", "/add", "1", "task1", "boo1", "false", "Adding the following task to your task list\nTask:\n\tId = 1\n\tTitle = task1\n\tDescription = boo1\n\tCompleted = false", http.StatusCreated},
		{"/tasks", "/add", "2", "task2", "boo2", "false", "Adding the following task to your task list\nTask:\n\tId = 2\n\tTitle = task2\n\tDescription = boo2\n\tCompleted = false", http.StatusCreated},
		{"/tasks", "/add", "3", "task3", "boo3", "false", "Adding the following task to your task list\nTask:\n\tId = 3\n\tTitle = task3\n\tDescription = boo3\n\tCompleted = false", http.StatusCreated},
		{"/tasks", "/add", "4", "task4", "boo4", "false", "Adding the following task to your task list\nTask:\n\tId = 4\n\tTitle = task4\n\tDescription = boo4\n\tCompleted = false", http.StatusCreated},
		{"/tasks", "/add", "5", "task5", "boo5", "false", "Adding the following task to your task list\nTask:\n\tId = 5\n\tTitle = task5\n\tDescription = boo5\n\tCompleted = false", http.StatusCreated},
		{"/tasks", "/complete", "1", "", "", "", "Completed task with id 1\n", http.StatusOK},
		{"/tasks", "/complete", "2", "", "", "", "Completed task with id 2\n", http.StatusOK},
		{"/tasks", "/complete", "6", "", "", "", "No task with ID = 6 to complete\n", http.StatusOK},
		{"/tasks", "", "", "", "", "false",
			"Getting all tasks...\n" +
				"Task:\n\tId = 3\n\tTitle = task3\n\tDescription = boo3\n\tCompleted = false\n" +
				"Task:\n\tId = 4\n\tTitle = task4\n\tDescription = boo4\n\tCompleted = false\n" +
				"Task:\n\tId = 5\n\tTitle = task5\n\tDescription = boo5\n\tCompleted = false\n", http.StatusOK},
	}

	for _, c := range cases {

		switch c.op {
		case "": //get tasks (GET)
			url := c.url + c.op + "?showCompleted=" + c.completed
			req := httptest.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()
			taskList.TasksHandler(w, req)
			res := w.Result()
			defer res.Body.Close()
			data, err := ioutil.ReadAll(res.Body)
			assert.Nil(t, err)
			assert.Equal(t, c.want, string(data))
			assert.Equal(t, c.respCode, res.StatusCode)
		case "/add": //add a task (POST)
			url := c.url + c.op
			completedBool, _ := strconv.ParseBool(c.completed)
			idInt, _ := strconv.ParseInt(c.id, 10, 64)
			values := map[string]interface{}{"id": idInt, "title": c.title, "description": c.description, "completed": completedBool}
			jsonValue, _ := json.Marshal(values)
			req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonValue))
			w := httptest.NewRecorder()
			taskList.AddTaskHandler(w, req)
			res := w.Result()
			defer res.Body.Close()
			data, err := ioutil.ReadAll(res.Body)
			assert.Nil(t, err)
			assert.Equal(t, c.want, string(data))
			assert.Equal(t, c.respCode, res.StatusCode)

		case "/complete": //complete a task (PATCH)
			url := c.url + c.op
			idInt, _ := strconv.ParseInt(c.id, 10, 64)
			values := map[string]interface{}{"id": idInt}
			jsonValue, _ := json.Marshal(values)
			req := httptest.NewRequest(http.MethodPatch, url, bytes.NewBuffer(jsonValue))
			w := httptest.NewRecorder()
			taskList.CompleteTaskHandler(w, req)
			res := w.Result()
			defer res.Body.Close()
			data, err := ioutil.ReadAll(res.Body)
			assert.Nil(t, err)
			assert.Equal(t, c.want, string(data))
			assert.Equal(t, c.respCode, res.StatusCode)

		case "/clear": //clear tasks (DELETE)
			url := c.url
			req := httptest.NewRequest(http.MethodDelete, url, nil)
			w := httptest.NewRecorder()
			taskList.TasksHandler(w, req)
			res := w.Result()
			defer res.Body.Close()
			data, err := ioutil.ReadAll(res.Body)
			assert.Nil(t, err)
			assert.Equal(t, c.want, string(data))
			assert.Equal(t, c.respCode, res.StatusCode)
		}
	}

}

func TestAddGetClearTasksCorrect(t *testing.T) {
	var taskList TaskList
	cases := []struct {
		url, op, id, title, description, completed, want string
		respCode                                         int
	}{
		{"/tasks", "/add", "1", "task1", "boo1", "false", "Adding the following task to your task list\nTask:\n\tId = 1\n\tTitle = task1\n\tDescription = boo1\n\tCompleted = false", http.StatusCreated},
		{"/tasks", "/add", "2", "task2", "boo2", "false", "Adding the following task to your task list\nTask:\n\tId = 2\n\tTitle = task2\n\tDescription = boo2\n\tCompleted = false", http.StatusCreated},
		{"/tasks", "/add", "3", "task3", "boo3", "false", "Adding the following task to your task list\nTask:\n\tId = 3\n\tTitle = task3\n\tDescription = boo3\n\tCompleted = false", http.StatusCreated},
		{"/tasks", "/complete", "6", "", "", "", "No task with ID = 6 to complete\n", http.StatusOK},
		{"/tasks", "/clear", "", "", "", "", "", http.StatusNoContent},
		{"/tasks", "", "", "", "", "true", "Getting all tasks...\nThere are no tasks!", http.StatusOK},
	}

	for _, c := range cases {

		switch c.op {
		case "": //get tasks (GET)
			url := c.url + c.op + "?showCompleted=" + c.completed
			req := httptest.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()
			taskList.TasksHandler(w, req)
			res := w.Result()
			defer res.Body.Close()
			data, err := ioutil.ReadAll(res.Body)
			assert.Nil(t, err)
			assert.Equal(t, c.want, string(data))
			assert.Equal(t, c.respCode, res.StatusCode)
		case "/add": //add a task (POST)
			url := c.url + c.op
			completedBool, _ := strconv.ParseBool(c.completed)
			idInt, _ := strconv.ParseInt(c.id, 10, 64)
			values := map[string]interface{}{"id": idInt, "title": c.title, "description": c.description, "completed": completedBool}
			jsonValue, _ := json.Marshal(values)
			req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonValue))
			w := httptest.NewRecorder()
			taskList.AddTaskHandler(w, req)
			res := w.Result()
			defer res.Body.Close()
			data, err := ioutil.ReadAll(res.Body)
			assert.Nil(t, err)
			assert.Equal(t, c.want, string(data))
			assert.Equal(t, c.respCode, res.StatusCode)

		case "/complete": //complete a task (PATCH)
			url := c.url + c.op
			idInt, _ := strconv.ParseInt(c.id, 10, 64)
			values := map[string]interface{}{"id": idInt}
			jsonValue, _ := json.Marshal(values)
			req := httptest.NewRequest(http.MethodPatch, url, bytes.NewBuffer(jsonValue))
			w := httptest.NewRecorder()
			taskList.CompleteTaskHandler(w, req)
			res := w.Result()
			defer res.Body.Close()
			data, err := ioutil.ReadAll(res.Body)
			assert.Nil(t, err)
			assert.Equal(t, c.want, string(data))
			assert.Equal(t, c.respCode, res.StatusCode)

		case "/clear": //clear tasks (DELETE) at endpoint /tasks just using /clear for testing
			url := c.url
			req := httptest.NewRequest(http.MethodDelete, url, nil)
			w := httptest.NewRecorder()
			taskList.TasksHandler(w, req)
			res := w.Result()
			defer res.Body.Close()
			data, err := ioutil.ReadAll(res.Body)
			assert.Nil(t, err)
			assert.Equal(t, c.want, string(data))
			assert.Equal(t, c.respCode, res.StatusCode)
		}
	}

}
