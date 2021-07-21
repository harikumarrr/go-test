package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Todo struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type storage struct {
	todo []Todo
}

var todoStorage *storage

func GetTodos(c *gin.Context) {

	if todoStorage == nil {
		todoStorage = &storage{}
	}

	if len(todoStorage.todo) == 0 {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, &todoStorage.todo)
	}
}

func CreateATodo(c *gin.Context) {
	if todoStorage == nil {
		todoStorage = &storage{}
	}
	var todo Todo
	c.BindJSON(&todo)
	fmt.Println(todo)
	todoStorage.todo = append(todoStorage.todo, todo)
	if len(todoStorage.todo) == 0 {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, &todoStorage.todo)
	}
}

func GetATodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}
	var todo Todo
	var isFound bool
	for _, td := range todoStorage.todo {
		if td.ID == uint(id) {
			todo = td
			isFound = true
		}
	}
	if !isFound {
		c.AbortWithStatus(http.StatusNotFound)
	}
	c.JSON(http.StatusOK, todo)
}

func DeleteATodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}
	for idx := range todoStorage.todo {
		if todoStorage.todo[idx].ID == uint(id) {
			todoStorage.todo = append(todoStorage.todo[:idx], todoStorage.todo[idx+1:]...)
		}
	}
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, gin.H{fmt.Sprintf("%d is ", id): "deleted"})
	}
}
