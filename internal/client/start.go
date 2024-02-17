package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/offlaneDefender/progress-tracker-go/internal/common"
)

func Start() {
	log.SetPrefix("greetings: ")
	log.SetFlags(0)
	// Read input from the user
	fmt.Println("Select a method to query the server with (1-4), 0 for exit:", "GET", "POST", "PUT", "DELETE")
	scanner := bufio.NewScanner(os.Stdin)
	goal := common.GoalPutBody{Name: "GoalTest", MaxTicks: 10}

	for scanner.Scan() {
		data := scanner.Text()
		switch data {
		case "1":
			go getRequest()
		case "2":
			go postRequest(goal)
		case "3":
			go putRequest(goal)
		case "4":
			go deleteRequest(goal)
		case "0":
			fmt.Println("Exiting")
			return
		default:
			fmt.Println("Invalid input")
		}
		fmt.Println("Select a method to query the server with (1-4):", "GET", "POST", "PUT", "DELETE")
	}
}

func getRequest() {
	resp, err := http.Get("http://localhost:8080/")

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)

	if readErr != nil {
		log.Fatalf("Error: %v", readErr)
	} else {
		log.Println("Body:", string(body))
	}
}

func postRequest(g common.Goal) {
	jsonGoal, err := json.Marshal(g)
	payload := bytes.NewBuffer(jsonGoal)

	if err != nil {
		log.Fatalf("JSON marshall error: %v", err)
	}

	resp, err := http.Post("http://localhost:8080/", "application/json", payload)

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)

	if readErr != nil {
		log.Fatalf("Error: %v", readErr)
	} else {
		log.Println("POST Body:", string(body))
	}
}

func putRequest(g common.Goal) {
	jsonGoal, err := json.Marshal(g)

	if err != nil {
		log.Fatalf("JSON marshall error: %v", err)
	}

	payload := bytes.NewBuffer(jsonGoal)

	req, err := http.NewRequest(http.MethodPut, "http://localhost:8080/", payload)

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)

	if readErr != nil {
		log.Fatalf("Error: %v", err)
	} else {
		log.Println("PUT Body:", string(body))
	}
}

func deleteRequest(g common.Goal) {
	jsonGoal, err := json.Marshal(g)

	if err != nil {
		log.Fatalf("JSON marshall error: %v", err)
	}

	payload := bytes.NewBuffer(jsonGoal)

	req, err := http.NewRequest("DELETE", "http://localhost:8080/", payload)

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)

	if readErr != nil {
		log.Fatalf("Error: %v", err)
	} else {
		log.Println("Delete Body:", string(body))
	}
}
