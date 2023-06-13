package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/DataDog/datadog-go/statsd"
	log "github.com/sirupsen/logrus"
	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

type Task struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

type UpdateTask struct {
	Id        int64 `json:"id"`
	Completed bool  `json:"completed"`
}

type TaskList struct {
	tasks       []Task
	numTasks    int
	numComplete int
}

var client *statsd.Client
var standardFields log.Fields

func getTaskAsString(task Task, res http.ResponseWriter) {
	fmt.Fprintf(res, "Task:\n\tId = %d\n\tTitle = %s\n\tDescription = %s\n\tCompleted = %t", task.Id, task.Title, task.Description, task.Completed)
}

// handler to deal with all /tasks routes (so far: Get all tasks and clear all tasks)
func (t *TaskList) TasksHandler(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		showCompleted := req.URL.Query().Get("showCompleted")
		showCompletedBool, err := strconv.ParseBool(showCompleted)
		if err != nil {
			log.Panic(err)
		}
		if len(t.tasks) == 0 {
			res.WriteHeader(http.StatusOK)
			fmt.Fprint(res, "Getting all tasks...\n")
			fmt.Fprint(res, "There are no tasks!")

		} else {
			res.WriteHeader(http.StatusOK)
			fmt.Fprint(res, "Getting all tasks...\n")

			// log number of tasks requested by the user
			info := fmt.Sprintf("User requested %d tasks", len(t.tasks))
			log.WithFields(standardFields).Info(info)

			for _, task := range t.tasks {
				if showCompletedBool {
					getTaskAsString(task, res)
					fmt.Fprint(res, "\n")
				} else {
					if !task.Completed {
						getTaskAsString(task, res)
						fmt.Fprint(res, "\n")
					}
				}
			}

		}
	case "DELETE":
		t.tasks = t.tasks[:0]
		res.WriteHeader(http.StatusNoContent)
		//send metrics
		client.Gauge("num_total_tasks.gauge", 0.0, []string{"environment:dev"}, 1)
		client.Gauge("num_complete_tasks.gauge", 0.0, []string{"environment:dev"}, 1)
		client.Gauge("num_incomplete_tasks.gauge", 0.0, []string{"environment:dev"}, 1)
	}
}

func (t *TaskList) AddTaskHandler(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/tasks/add" {
		http.NotFound(res, req)
		return
	}

	var task Task
	err := json.NewDecoder(req.Body).Decode(&task)
	if err != nil {
		http.Error(res, err.Error(), 400)
		return
	}
	res.WriteHeader(http.StatusCreated)
	t.tasks = append(t.tasks, task)
	fmt.Fprint(res, "Adding the following task to your task list\n")
	getTaskAsString(task, res)

	//log user action
	info := fmt.Sprintf("Added task with Id = %d, Title = %s, Description = %s\n", task.Id, task.Title, task.Description)
	log.WithFields(standardFields).Info(info)

	//send metrics
	t.numTasks += 1
	numIncomplete := t.numTasks - t.numComplete
	client.Gauge("num_total_tasks.gauge", float64(t.numTasks), []string{"environment:dev"}, 1)
	client.Gauge("num_complete_tasks.gauge", float64(t.numComplete), []string{"environment:dev"}, 1)
	client.Gauge("num_incomplete_tasks.gauge", float64(numIncomplete), []string{"environment:dev"}, 1)

}

func (t *TaskList) CompleteTaskHandler(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/tasks/complete" {
		http.NotFound(res, req)
		return
	}
	var update UpdateTask
	err := json.NewDecoder(req.Body).Decode(&update)
	if err != nil {
		http.Error(res, err.Error(), 400)
		return
	}
	res.WriteHeader(http.StatusOK)
	id := update.Id
	completedSomething := false
	for i, task := range t.tasks {
		if task.Id == id {
			if task.Completed {
				fmt.Fprintf(res, "Task %d is already completed\n", id)

			} else {
				t.tasks[i].Completed = true
				completedSomething = true
				fmt.Fprintf(res, "Completed task with id %d\n", id)
				log.Printf("Completed task with id %d\n", id)
				//send metrics
				t.numComplete += 1
				numIncomplete := t.numTasks - t.numComplete
				client.Gauge("num_complete_tasks.gauge", float64(t.numComplete), []string{"environment:dev"}, 1)
				client.Gauge("num_incomplete_tasks.gauge", float64(numIncomplete), []string{"environment:dev"}, 1)
			}
		}
	}
	if !completedSomething {
		fmt.Fprintf(res, "No task with ID = %d to complete\n", id)
	}

}

func (t *TaskList) MainPageHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	fmt.Fprintf(res, "Welcome to your super simple task manager\n")
	log.WithFields(standardFields).Info("Main page accessed")
}

func main() {
	//ingest command-line args to get optionally specified log path
	logPath := os.Args[1]
	if logPath == "" {
		logPath = "mini-server-logs.log" //path is from our current directory not absolute
	}
	//configure port to send statsd metric over

	//for my server running on host
	//client, _ = statsd.New("127.0.0.1:8125")

	//for my server running in docker container
	client, _ = statsd.New("host.docker.internal:8125")
	//for this to work make sure to expose this port on the host by running basic http server with "python -m http.server 8125"

	//configure log location
	f, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Panicf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.SetFormatter(&log.JSONFormatter{})

	//configure standard log fields
	standardFields = log.Fields{
		"hostname": "Calebs Mac",
		"appname":  "mini-golang-http-server",
		"session":  "testing",
	}

	//example log with fields
	log.WithFields(standardFields).WithFields(log.Fields{"string": "foo", "int": 1, "float": 1.1}).Info("My first ssl event from Golang")
	log.WithFields(standardFields).Info("Server started")

	// task list
	var taskList TaskList

	//congifure and set up apm and http routing and multiplexer
	tracer.Start(
		tracer.WithService("task-manager"),
		tracer.WithEnv("dev"),
	)
	defer tracer.Stop()

	err = profiler.Start(
		profiler.WithService("task-manager"),
		profiler.WithEnv("dev"),
		profiler.WithProfileTypes(
			profiler.CPUProfile,
			profiler.HeapProfile,
		),
	)
	if err != nil {
		log.Panic(err)
	}
	defer profiler.Stop()

	// Create a traced mux router
	mux := httptrace.NewServeMux()
	mux.HandleFunc("/", taskList.MainPageHandler)
	mux.HandleFunc("/tasks", taskList.TasksHandler)
	mux.HandleFunc("/tasks/add", taskList.AddTaskHandler)
	mux.HandleFunc("/tasks/complete", taskList.CompleteTaskHandler)
	http.ListenAndServe(":9000", mux)

}
