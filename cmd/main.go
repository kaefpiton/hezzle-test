package main

import (
	"context"
	"fmt"
	"hezzle/cmd/providers"
	"hezzle/configs"
	"hezzle/internal/infrastructure/interfaces/httpControllers"
	queue2 "hezzle/internal/infrastructure/queue/nats"
	repository2 "hezzle/internal/infrastructure/repository"
	interactors "hezzle/internal/infrastructure/usecase/interractors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// todo убрать
type logDb struct {
	Id          int
	CampaignId  int
	Name        string
	Description string
	Priority    int
	Removed     bool
	EventTime   time.Time
}

const configPath = "configs/config.json"

func main() {
	cnf, err := configs.LoadConfig(configPath)
	if err != nil {
		log.Panic(err)
	}
	//logger
	logger, err := providers.ProvideConsoleLogger(cnf)
	if err != nil {
		log.Panic(err)
	}

	//postgres
	db, closeDB, err := providers.ProvidePostgres(cnf, logger)
	if err != nil {
		log.Panic(err)
	}

	//redis
	redisClient, err := providers.ProvideRedis(cnf)
	if err != nil {
		log.Panic(err)
	}

	//nats
	queue, err := providers.ProvideQueue(cnf)
	if err != nil {
		log.Panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	itemsRepository := repository2.NewItemsRepository(db, redisClient, logger)
	itemsInteractor := interactors.NewItemsInteractor(itemsRepository, redisClient, queue, logger)
	itemController := httpControllers.NewItemsController(ctx, itemsInteractor, logger)

	clickHouseConn := providers.ProvideClickhouse(cnf)
	logRepo := repository2.NewLogsRepository(clickHouseConn)
	eventListener := queue2.NewEventListener(ctx, queue, logRepo, logger)
	go eventListener.ListenTopic()

	server := providers.ProvideHTTPServer(cnf, itemController, logger)

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
		fmt.Println("Terminating the app")
		fmt.Println("Shutdown workers")
		cancel()

		fmt.Println("Close DB")
		closeDB()

		fmt.Println("Stop Server")
		server.Stop(ctx)
	}()

	server.Start()
}
