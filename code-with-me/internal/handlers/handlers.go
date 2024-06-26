package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

type CodeHandler struct {
	service CodeService
}

type CodeService interface {
	Echo(message string) string
	CreateNewContainerFromImage()
	CreateNewContainerFromFile() (string, error)
	StartContainerByID(containerID string, conn *websocket.Conn) error
}

func New(service CodeService) *CodeHandler {
	return &CodeHandler{service: service}
}

func InitRouter(router *gin.Engine, handlers *CodeHandler) *gin.Engine {
	codeHandlers := router.Group("/code")
	codeHandlers.POST("/create", handlers.CreateContainerFromFile)
	codeHandlers.GET("/start", handlers.StartContainer)

	return router
}

func (h CodeHandler) CreateContainerFromFile(ctx *gin.Context) {
	fmt.Println(ctx.RemoteIP())
	file, err := ctx.FormFile("file")
	if err != nil {
		fmt.Println(err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(file.Filename)
	if err != nil {
		fmt.Println(err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "can't create directory"})
		return
	}
	err = ctx.SaveUploadedFile(file, "/home/user/GolandProjects/code-with-me/app/test/test.go")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	id, err := h.service.CreateNewContainerFromFile()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	ctx.AbortWithStatusJSON(http.StatusCreated, gin.H{"container_id": id})

	return
}

func (h CodeHandler) StartContainer(ctx *gin.Context) {
	//_ = ctx.Param("container_id")

	upgrader := websocket.Upgrader{}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println("error upgrading", err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "can't update connection to web socket"})
		return
	}
	if conn == nil {
		fmt.Println("connection is nil")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "connection is nil"})
		return
	}
	_, containerID, err := conn.ReadMessage()
	fmt.Println(string(containerID))
	err = conn.WriteMessage(websocket.TextMessage, []byte("hello from WebSocket!"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusTeapot, gin.H{"error": err.Error()})
		return
	}

	err = h.service.StartContainerByID(string(containerID), conn)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	return
}
