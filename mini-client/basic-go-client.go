package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/DataDog/datadog-go/statsd"
)

const endpointURL string = "http://10.244.0.49:9000" //local ip of task manager container

func addTask(id int, t string, d string, cb bool) {
	c := http.Client{Timeout: time.Duration(1) * time.Second}
	values := map[string]interface{}{"id": id, "title": t, "description": d, "completed": cb}
	jsonValue, err := json.Marshal(values)
	if err != nil {
		log.Panic(err)
	}
	resp, err := c.Post(endpointURL+"/tasks/add", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Panicln(err)
	}
	if resp.StatusCode != http.StatusCreated {
		log.Panicln("invalid status code: ", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("Body : %s", body)
}

func completeTask(id int) {
	values := map[string]interface{}{"id": id}
	jsonValue, _ := json.Marshal(values)
	c := http.Client{Timeout: time.Duration(1) * time.Second}
	req, err := http.NewRequest(http.MethodPatch, endpointURL+"/tasks/complete", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.Do(req)
	if err != nil {
		log.Panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Panicln("invalid status code: ", resp.StatusCode)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}
	log.Println(string(body))

}

func getTasks(showCompleted bool) {
	url := endpointURL + "/tasks?showCompleted=" + strconv.FormatBool(showCompleted)
	resp, err := http.Get(url)
	if err != nil {
		log.Panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Panicln("invalid status code: ", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}
	log.Println(string(body))
}

func newClient() *statsd.Client {
	statsd, err := statsd.New("")
	if err != nil {
		log.Panic(err)
	}
	return statsd
}

func main() {
	var client *statsd.Client = newClient()

	for {
		for i := 1; i < 10; i++ {
			title := "task #" + strconv.Itoa(i)

			//collect data on add tasks
			start := time.Now()
			addTask(i, title, "boo1", false)
			elasped := time.Since(start).Seconds()
			client.Histogram("add_task_exec_time.histogram", elasped, []string{"environment:dev"}, 1)

			//collect data on get tasks
			start = time.Now()
			getTasks(true)
			elasped = time.Since(start).Seconds()
			client.Histogram("get_tasks_exec_time.histogram", elasped, []string{"environment:dev"}, 1)

			//collect data on time to complete tasks
			start = time.Now()
			completeTask(1) //completes all tasks with id of 1
			elasped = time.Since(start).Seconds()
			client.Histogram("get_tasks_exec_time.histogram", elasped, []string{"environment:dev"}, 1)

		}
		time.Sleep(120 * time.Second)
	}
}
