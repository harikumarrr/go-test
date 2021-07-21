package main

import (
	"fmt"
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"go-test/pkg/routes"
)

func main() {
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	r := routes.SetupRouter()

	err := r.Run("localhost:8085")
	fmt.Println("Error: ", err)
}
