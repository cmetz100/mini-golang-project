package main

import (
	"fmt"
    "io/ioutil"
    "net/http"
    "time"
	"bytes"
	"encoding/json"
	"log"
	"github.com/DataDog/datadog-go/statsd"
	"strconv"
)

func addTask (id int, t string, d string, cb bool){
	c := http.Client{Timeout: time.Duration(1) * time.Second}
	values := map[string]interface{}{"id": id, "title": t, "description" : d, "completed": cb}
	jsonValue, _ := json.Marshal(values)
	resp, err := c.Post("http://127.0.0.1:9000/tasks/add", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Printf("Body : %s", body)
}

func completeTask (id int){
	values := map[string]interface{}{"id": id}
	jsonValue, _ := json.Marshal(values)
	c := http.Client{Timeout: time.Duration(1) * time.Second}
	req,err := http.NewRequest(http.MethodPatch,"http://127.0.0.1:9000/tasks/complete",bytes.NewBuffer(jsonValue))
    if err != nil {
        log.Fatal(err)
    }
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.Do(req)
    if err != nil {
        log.Fatal(err)
    }

    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }
    log.Println(string(body))

}

func getTasks(showCompleted bool){
	url := "http://127.0.0.1:9000/tasks?showCompleted=" +strconv.FormatBool(showCompleted)
	resp, err := http.Get(url)
	if err != nil {
	   log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
      log.Fatalln(err)
    }
	log.Println(string(body))
}


func newClient() *statsd.Client {
	statsd, err := statsd.New("127.0.0.1:8125")
    if err != nil {
        log.Fatal(err)
    }
	return statsd
}

func main(){
	var client *statsd.Client =newClient()

	for {
		for i := 1; i < 10; i++ {
			title := "task #" + strconv.Itoa(i)
	
			//collect data on add tasks
			start := time.Now()
			addTask(i,title,"boo1",false)
			elasped := time.Since(start).Seconds()
			client.Histogram("add_task_exec_time.histogram",elasped, []string{"environment:dev"}, 1)
			
			//collect data on get tasks
			start = time.Now()
			getTasks(true)
			elasped = time.Since(start).Seconds()
			client.Histogram("get_tasks_exec_time.histogram",elasped, []string{"environment:dev"}, 1)
			
		}
		time.Sleep(120*time.Second)
	}
}