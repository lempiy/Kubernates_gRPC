package main

import (
	"github.com/lempiy/kubegrpc/pb"
	"google.golang.org/grpc"
	"github.com/gin-gonic/gin"
	"strconv"
	"net/http"
	"fmt"
)

func main() {
	conn, err := grpc.Dial("gcd-service:3000", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("Dial failed: %v\n", err)
	}
	gcdClient := pb.NewGCDServiceClient(conn)
	r := gin.Default()
	r.GET("/gcd/:a/:b", func(c *gin.Context) {
		// Parse parameters
		a, err := strconv.ParseUint(c.Param("a"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameter A"})
			return
		}
		b, err := strconv.ParseUint(c.Param("b"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameter B"})
			return
		}
		// Call GCD service
		req := &pb.GCDRequest{A: a, B: b}
		if res, err := gcdClient.Compute(c, req); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"result": fmt.Sprint(res.Result),
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	})
	if err := r.Run(":3000"); err != nil {
		fmt.Printf("Failed to run server: %v\n", err)
	}
}