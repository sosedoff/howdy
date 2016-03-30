package main

import (
	"github.com/gin-gonic/gin"
)

func handleGetServices(c *gin.Context) {
	names := []string{}

	DB.
		Table("check_records").
		Order("service asc").
		Pluck("distinct service", &names)

	c.JSON(200, names)
}

func handleGetService(c *gin.Context) {
	records := []CheckRecord{}

	DB.
		Table("check_records").
		Where("service = ?", c.Param("name")).
		Order("timestamp DESC").
		Limit(30).
		Find(&records)

	c.JSON(200, records)
}

func runHttpServer() {
	router := gin.Default()
	router.GET("/", handleGetServices)
	router.GET("/services", handleGetServices)
	router.GET("/services/:name", handleGetService)
	router.Run("0.0.0.0:5000")
}
