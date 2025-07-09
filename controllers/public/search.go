package public

import (
	"BookMyTurf/config"
	"BookMyTurf/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SearchTurfs(c *gin.Context) {

	location := c.Query("location")
	date := c.Query("date")

	if location == "" || date == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "location and date are required",
		})
		return
	}

	var turfs []models.Turf

	// Finding turfs

	subQuery := config.DB.Table("bookings").Select("turf_id").Where("DATE(start_time) = ?", date)

	err := config.DB.Where("location LIKE ?", "%"+location+"%").Where("id NOT IN (?)", subQuery).Find(&turfs).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to search turfs",
		})
		return
	}

	c.JSON(http.StatusOK, turfs)

}
