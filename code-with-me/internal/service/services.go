package service

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/archive"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"log/slog"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type CodeService struct {
}

func New() CodeService {
	return CodeService{}
}

func (s CodeService) Echo(message string) string {
	return "echo " + message
}

type ErrorLine struct {
	Error       string      `json:"error"`
	ErrorDetail ErrorDetail `json:"errorDetail"`
}

type ErrorDetail struct {
	Message string `json:"message"`
}

func (s CodeService) CreateNewContainerFromImage() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		slog.Error("error creating docker client : " + err.Error())
		return
	}

	ctx := context.Background()
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:        "nginx:latest",
		ExposedPorts: nat.PortSet{"8082": struct{}{}},
	}, &container.HostConfig{
		PortBindings: map[nat.Port][]nat.PortBinding{"8082": {{HostIP: "127.0.0.1", HostPort: "8082"}}},
	}, nil, nil, "mongo-go-cli")
	if err != nil {
		slog.Error("error creating container : " + err.Error())
		return
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		slog.Error("error starting container : " + err.Error())
		return
	}
	return
}

func (s CodeService) CreateNewContainerFromFile() (string, error) {

	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()

	tar, err := archive.Tar("/home/user/GolandProjects/code-with-me/app/test", 2)
	if err != nil {
		return "", fmt.Errorf("error archiving test : %w", err)
	}
	var result []byte
	n, err := tar.Read(result)
	imageName := "test-app"
	fmt.Println("RESULT : ", string(result), "N : ", n, "ERR: ", err)
	opts := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{imageName},
	}
	res, err := dockerClient.ImageBuild(ctx, tar, opts)
	if err != nil {
		return "", fmt.Errorf("error bilding image : %w", err)
	}

	defer res.Body.Close()

	err = print(res.Body)
	if err != nil {
		return "", err
	}

	resp, err := dockerClient.ContainerCreate(ctx, &container.Config{
		Image:        imageName,
		OpenStdin:    true,
		AttachStdout: true,
		AttachStderr: true,
		AttachStdin:  true,
	}, nil, nil, nil, imageName)
	if err != nil {
		log.Printf("Failed to create container: %s\n", err)
		return "", err
	}

	return resp.ID, nil
}

func (s CodeService) StartContainerByID(containerID string, conn *websocket.Conn) error {
	isContainerRemoved := false
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()

	attachResp, err := dockerClient.ContainerAttach(ctx, containerID, container.AttachOptions{Stream: true, Stdin: true, Stdout: true, Stderr: true})
	if err != nil {
		return err
	}

	go func() {
		buf := make([]byte, 1024)
		for {
			if isContainerRemoved {
				return
			}
			n, err := attachResp.Reader.Read(buf)
			if err != nil {
				// Handle other errors
				log.Println("Error reading from container output:", err)
				return
			}

			err = conn.WriteMessage(websocket.TextMessage, buf[:n])
			if err != nil {
				isContainerRemoved = true
				_ = dockerClient.ContainerStop(ctx, containerID, container.StopOptions{})
				_ = dockerClient.ContainerRemove(ctx, containerID, container.RemoveOptions{})
				if err != nil {
					log.Println("error removing container : ", err.Error())
				}
				log.Println("Error writing to WebSocket connection:", err)
				return
			}
		}
		//conn.WriteMessage()
		//_, _ = io.Copy(os.Stdout, attachResp.Reader) // Print the container output to stdout
	}()

	// Read input from the user and write it to the container's stdin
	go func() {
		for {
			if isContainerRemoved {
				return
			}
			_, message, err := conn.ReadMessage()
			if err != nil {
				_ = dockerClient.ContainerStop(ctx, containerID, container.StopOptions{})
				_ = dockerClient.ContainerRemove(ctx, containerID, container.RemoveOptions{})
				log.Println("error reading from client : ", err.Error())
				conn.Close()
				return
			}
			_, err = attachResp.Conn.Write(message)
			if err != nil {
				log.Println("error writing data to container : ", err.Error())
				return
			}
		}
		//_, _ = io.Copy(attachResp.Conn, os.Stdin)
	}()

	if err := dockerClient.ContainerStart(ctx, containerID, container.StartOptions{}); err != nil {
		slog.Error("error starting container : " + err.Error())
		return err
	}
	fmt.Println("Docker container started!")

	waitChan, errChan := dockerClient.ContainerWait(ctx, containerID, container.WaitConditionNotRunning)
	select {
	case <-waitChan:
		err := dockerClient.ContainerStop(ctx, containerID, container.StopOptions{})
		if err != nil {
			log.Println("error stopping container : ", err.Error())
		}
		err = dockerClient.ContainerRemove(ctx, containerID, container.RemoveOptions{})
		if err != nil {
			log.Println("error removing container : ", err.Error())
		}
		return errors.New("not running")
	case <-errChan:
		return <-errChan
	}

}

func print(rd io.Reader) error {
	var lastLine string

	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		lastLine = scanner.Text()
		fmt.Println(scanner.Text())
	}

	errLine := &ErrorLine{}
	json.Unmarshal([]byte(lastLine), errLine)
	if errLine.Error != "" {
		return errors.New(errLine.Error)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
