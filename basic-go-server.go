package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Task struct {
	id          int64
	title       string
	description string
	completed   bool
}

type TaskList struct {
	tasks []Task
}

// task list
var taskList TaskList

func GetTasksHandler(res http.ResponseWriter, req *http.Request) {
	showCompleted := req.URL.Query().Get("showCompleted")
	showCompletedBool, _ := strconv.ParseBool(showCompleted)
	fmt.Fprint(res, "Getting all tasks\n")

	for _, task := range taskList.tasks {
		if showCompletedBool {
			fmt.Fprint(res, task)
			fmt.Fprint(res, "\n")
		} else {
			if !task.completed {
				fmt.Fprint(res, task)
				fmt.Fprint(res, "\n")
			}
		}
	}

}

func AddTaskHandler(res http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	idInt, _ := strconv.ParseInt(id, 10, 64)

	title := req.URL.Query().Get("title")
	description := req.URL.Query().Get("description")

	completed := req.URL.Query().Get("completed")
	completedBool, _ := strconv.ParseBool(completed)

	task := Task{id: idInt, title: title, description: description, completed: completedBool}
	taskList.tasks = append(taskList.tasks, task)

	fmt.Fprint(res, "Adding the following task to your tak list")
	fmt.Fprint(res, task)
}

func CompleteTaskHandler(res http.ResponseWriter, req *http.Request) {

	id := req.URL.Query().Get("id")
	idInt, _ := strconv.ParseInt(id, 10, 64)
	for i, task := range taskList.tasks {
		if task.id == idInt {
			if task.completed {
				fmt.Fprintf(res, "Task %d is already completed\n", idInt)
				return
			} else {
				taskList.tasks[i].completed = true
				fmt.Fprintf(res, "Completed task with id %d\n", idInt)
				return
			}
		}
	}
	fmt.Fprintf(res, "No task with ID = %d to complete\n", idInt)
}

func main() {

	http.HandleFunc("/tasks", GetTasksHandler)
	http.HandleFunc("/tasks/add", AddTaskHandler)
	http.HandleFunc("/tasks/complete", CompleteTaskHandler)

	// start HTTP server with `http.DefaultServeMux` handler
	log.Fatal(http.ListenAndServe(":9000", nil))

}
