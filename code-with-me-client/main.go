package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
)

var addr = flag.String("addr", "localhost:8003", "http service address")

type Container struct {
	ID string `json:"container_id"`
}

func main() {
	log.Println("Starting code-with-me client")

	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	file, err := os.Open("/home/user/GolandProjects/code-with-me-client/test.go")
	if err != nil {
		log.Println("error opening file : ", err)
		return
	}
	fmt.Println(file.Name(), file.Fd())
	part, err := writer.CreateFormFile("file", file.Name())
	if err != nil {
		log.Println("error creating form file")
		return
	}

	_, err = io.Copy(part, file)
	if err != nil {
		log.Println("error copying file data:", err)
		return
	}
	err = writer.Close()
	if err != nil {
		log.Println("error closing writer : ", err)
	}
	resp, err := http.Post("http://localhost:8003/code/echo", writer.FormDataContentType(), buf)
	if err != nil {
		log.Println(err)
		return
	}
	var respBuf []byte
	respBuf, _ = io.ReadAll(resp.Body)
	fmt.Println("Response status : ", resp.Status)
	fmt.Println(string(respBuf))
	var container Container
	err = json.Unmarshal(respBuf, &container)
	if err != nil {
		fmt.Println("error unmarshalling : ", err.Error())
	}

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/code/start"}
	log.Printf("connecting to %s", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Println("ERROR DIALING WITH WS SERVER : ", err.Error())
		return
	}
	defer conn.Close()

	err = conn.WriteMessage(websocket.TextMessage, []byte(container.ID))

	var message []byte
	_, message, err = conn.ReadMessage()
	if err != nil {
		log.Println("read:", err)
		return
	}
	log.Printf("recv: %s", string(message))

	go func() {
		buf := make([]byte, 1024)
		for {
			_, err := os.Stdin.Read(buf)
			if err != nil {
				log.Println("Error reading from stdin:", err)
				return
			}

			err = conn.WriteMessage(websocket.TextMessage, buf)
			if err != nil {
				log.Println("Error writing to WebSocket:", err)
				return
			}
		}
	}()

	// Start the goroutine to read from the WebSocket and write to stdout
	go func() {

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading from WebSocket:", err)
				return
			}

			_, err = os.Stdout.Write(message)
			if err != nil {
				log.Println("Error writing to stdout:", err)
				return
			}
		}
	}()
	select {}
}

//func main() {
//	resp, err := http.Get("http://localhost:8003/hello")
//	if err != nil {
//		fmt.Println("Error:", err)
//		return
//	}
//	defer resp.Body.Close()
//
//	// Read the response body
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		fmt.Println("Error:", err)
//		return
//	}
//
//	// Print the response
//	fmt.Println(string(body))
//}

//func main() {
//	// Open the file
//	file, err := os.Open("/home/user/GolandProjects/code-with-me-client/test.go")
//	if err != nil {
//		fmt.Println("Error:", err)
//		return
//	}
//	defer file.Close()
//
//	// Create a new HTTP request
//	req, err := http.NewRequest("POST", "http://localhost:8003/code/echo", file)
//	if err != nil {
//		fmt.Println("Error:", err)
//		return
//	}
//
//	// Send the request
//	client := &http.Client{}
//	resp, err := client.Do(req)
//	if err != nil {
//		fmt.Println("Error:", err)
//		return
//	}
//	defer resp.Body.Close()
//
//	// Read the response
//	body, err := io.ReadAll(resp.Body)
//	if err != nil {
//		fmt.Println("Error:", err)
//		return
//	}
//
//	// Print the response
//	fmt.Println(string(body))
//}
