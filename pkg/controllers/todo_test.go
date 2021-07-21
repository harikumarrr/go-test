package controllers

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetTodos(t *testing.T) {
	//todoStorage = nil

	ctx, err := gin.CreateTestContext(httptest.NewRecorder())
	fmt.Println("err:", err)
	GetTodos(ctx)

	fmt.Println("ctx:", ctx)

}
