package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

const serverPort = 3333

func main() {

	go func() {
		// Multiplexing request handlers
		mux := http.NewServeMux()

		// Setup the Handler Functions
		mux.HandleFunc("/", getRoot)

		// Create the HTTP server
		server := http.Server{
			Addr:    fmt.Sprintf(":%d", serverPort),
			Handler: mux,
		}

		// Listen every request on port 3333 irrespective of the IP address
		err := server.ListenAndServe()

		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Server closed\n")
		} else if err != nil {
			fmt.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	// Create a json of type byte slice
	jsonBody := []byte(`{"client_message": "hello, server!"}`)
	bodyReader := bytes.NewReader(jsonBody)

	// Create a POST request
	requestURL := fmt.Sprintf("http://localhost:%d?id=1234", serverPort)
	req, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)

	if err != nil {
		fmt.Printf("Client: could not create request: %s\n", err)
		os.Exit(1)
	}

	req.Header.Set("Content-Type", "application/json")

	// create a specific client and sent request
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	// send the Client request
	res, err := client.Do(req)

	if err != nil {
		fmt.Printf("Error making http request: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Client: got response!\n")
	fmt.Printf("Client: status code: %d\n", res.StatusCode)

	// read the response
	resBody, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Printf("Client: could not read response body: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Client: response body: %s\n", &resBody)
}

// HTTP Handler Functions
// w -> write information, r -> get information
func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("server: %s /\n", r.Method)

	// print various information received
	fmt.Printf("server: query id: %s\n", r.URL.Query().Get("id"))
	fmt.Printf("server: content-type: %s\n", r.Header.Get("Content-Type"))
	fmt.Printf("server: headers:\n")
	for headerName, headerValue := range r.Header {
		fmt.Printf("\t%s = %s\n", headerName, strings.Join(headerValue, ", "))
	}

	// read the request body
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("server: could not read request body: %s\n", err)
	}

	fmt.Printf("server: request body: %s\n", &reqBody)

	// sending JSON response
	fmt.Fprintf(w, `{"Message": "hello!"}`)

	// io.WriteString(w, "This is my website!\n")
}

