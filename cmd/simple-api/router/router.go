package router

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	config "github.com/AntonioMartinezFernandez/immutable-databases/cmd/simple-api/config"
	"github.com/AntonioMartinezFernandez/immutable-databases/cmd/simple-api/events"
)

func SetupRouter(
	cfg config.Config,
	eventRepository events.EventRepository,
) *gin.Engine {
	var r *gin.Engine
	if cfg.Mode == "debug" {
		r = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
		r.Use(gin.Recovery())
	}

	r.GET("/holi", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	authorized := r.Group(
		"/",
		gin.BasicAuth(gin.Accounts{
			"manoli": "holi", // Header: Authorization: Basic bWFub2xpOmhvbGk=
		}),
	)

	authorized.POST("/api/tracking/events", func(c *gin.Context) {
		// Always return 202 Accepted
		c.JSON(http.StatusAccepted, gin.H{"status": "Accepted"})

		// Parse JSON request
		var jsonBody RequestDto
		if c.Bind(&jsonBody) != nil {
			fmt.Println("invalid request")
		}

		var dtoEvents []events.EventDto
		for _, event := range jsonBody.Events {
			e := events.EventDto{
				Id:       uuid.New().String(),
				StreamId: jsonBody.TransactionId,
				Content:  fmt.Sprintf("%v", event),
			}
			dtoEvents = append(dtoEvents, e)
		}
		err := eventRepository.Save(c.Request.Context(), dtoEvents)
		if err != nil {
			fmt.Println("ERROR SAVING EVENTS ==>", err)
			return
		}

		// fmt.Println("Processed event ==>", jsonBody.Events)
		fmt.Print(".")
	})

	return r
}
