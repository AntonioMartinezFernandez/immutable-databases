package main

import (
	"context"
	"fmt"

	config "github.com/AntonioMartinezFernandez/immutable-databases/cmd/simple-api/config"
	"github.com/AntonioMartinezFernandez/immutable-databases/cmd/simple-api/db"
	"github.com/AntonioMartinezFernandez/immutable-databases/cmd/simple-api/events"
	router "github.com/AntonioMartinezFernandez/immutable-databases/cmd/simple-api/router"
)

func main() {
	ctx := context.Background()
	cfg := config.LoadEnvConfig()

	dbClient, err := db.NewImmuDbSqlClient(ctx, cfg)
	if err != nil {
		panic(err)
	}
	eventRepository := events.NewImmudbEventRepository(dbClient, cfg.DbTable)

	r := router.SetupRouter(cfg, eventRepository)
	r.Run(fmt.Sprintf(":%s", cfg.HttpPort))
}
