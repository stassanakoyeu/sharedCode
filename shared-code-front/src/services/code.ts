import axios from "axios";
import * as React from 'react';
import { useState } from 'react';

export function runCode(file: File) {
    let a
    //const [messages, setMessages] = useState([]);
    const formData = new FormData();
    formData.append('file', file);
    console.log(file.name);

    axios.post('http://localhost:8003/code/create', formData)
        .then(function (response) {
            console.log(response.data.container_id);
            localStorage.setItem('container_id', response.data.container_id);

            // Create the WebSocket instance
            const socket = new WebSocket("ws://localhost:8003/code/start");

            // Handle WebSocket events
            socket.addEventListener("open", () => {
                console.log("WebSocket connection established");
                socket.send(localStorage.getItem('container_id'));
            });

            socket.addEventListener("message", (event) => {
                console.log("Message from server:", event.data);
                a += event.data
                //setMessages((prevMessages) => [...prevMessages, event.data]);
            });

            socket.addEventListener("error", (error) => {
                console.error("WebSocket error:", error);
            });

            socket.addEventListener("close", () => {
                console.log("WebSocket connection closed");
            });
        })
        .catch(function (error) {
            console.error("HTTP request error:", error);
        });
    console.log('a :',a)
    return a;
}

export function createTextFile(value: string) {
    const fileContent = value;
    const fileName = 'test.go';

    // Create a Blob object with the file content and specify the MIME type
    const file = new File([fileContent], fileName, { type: 'text/plain' });

    return file;
}