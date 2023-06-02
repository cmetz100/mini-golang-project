package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"reflect"
	"strings"
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

func getTaskAsString(input interface{})(string){
	values := reflect.ValueOf(input)
	numFields := values.NumField()  
	types := values.Type()

	var sb strings.Builder
	sb.WriteString("Task:\n")
	for i := 0; i < numFields; i++ {
	  field := types.Field(i)
	  fieldValue := values.Field(i)
  
	  fmt.Fprintf(&sb,"\t%s = %v\n",field.Name ,fieldValue)
	}

	return sb.String()
}

func GetTasksHandler(res http.ResponseWriter, req *http.Request) {
	showCompleted := req.URL.Query().Get("showCompleted")
	showCompletedBool, _ := strconv.ParseBool(showCompleted)
	fmt.Fprint(res, "Getting all tasks\n")

	for _, task := range taskList.tasks {
		if showCompletedBool {
			fmt.Fprint(res, getTaskAsString(task))
			fmt.Fprint(res, "\n")
		} else {
			if !task.completed {
				fmt.Fprint(res, getTaskAsString(task))
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

	fmt.Fprint(res, "Adding the following task to your task list\n")
	fmt.Fprint(res, getTaskAsString(task) )
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

func MainPageHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "Welcome to your super simple task manager\n")

}

func main() {

	http.HandleFunc("/tasks", GetTasksHandler)
	http.HandleFunc("/tasks/add", AddTaskHandler)
	http.HandleFunc("/tasks/complete", CompleteTaskHandler)
	http.HandleFunc("/",MainPageHandler)

	// start HTTP server with `http.DefaultServeMux` handler
	log.Fatal(http.ListenAndServe(":9000", nil))

}
