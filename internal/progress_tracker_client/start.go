package progress_tracker_client

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
)

func Start() {
	// Read input from the user
	fmt.Println("Select a method to query the server with (1-4), 0 for exit:", "GET", "POST", "PUT", "DELETE")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		data := scanner.Text()
		switch data {
		case "1":
			go getRequest()
		case "2":
			go postRequest()
		case "3":
			go putRequest()
		case "4":
			go deleteRequest()
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
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Response received")
	}

	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)

	if readErr != nil {
		fmt.Println("Error:", readErr)
	} else {
		fmt.Println("Body:", string(body))
	}
}

func postRequest() {
	resp, err := http.Post("http://localhost:8080/", "application/json", nil)

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Response received")
	}

	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)

	if readErr != nil {
		fmt.Println("Error:", readErr)
	} else {
		fmt.Println("Body:", string(body))
	}
}

func putRequest() {
	req, err := http.NewRequest("PUT", "http://localhost:8080/", nil)

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Request created")
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Response received")
	}

	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)

	if readErr != nil {
		fmt.Println("Error:", readErr)
	} else {
		fmt.Println("Body:", string(body))
	}
}

func deleteRequest() {
	req, err := http.NewRequest("DELETE", "http://localhost:8080/", nil)

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Request created")
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Response received")
	}

	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)

	if readErr != nil {
		fmt.Println("Error:", readErr)
	} else {
		fmt.Println("Body:", string(body))
	}
}
