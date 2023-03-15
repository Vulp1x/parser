package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	datasetsservice "github.com/inst-api/parser/gen/datasets_service"
	"github.com/inst-api/parser/internal/config"
	"github.com/inst-api/parser/internal/dbmodel"
	"github.com/inst-api/parser/internal/mw"
	"github.com/inst-api/parser/internal/postgres"
	"github.com/inst-api/parser/internal/service"
	"github.com/inst-api/parser/internal/store/datasets"
	"github.com/inst-api/parser/internal/workers"
	"github.com/inst-api/parser/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	loc := time.FixedZone("Moscow", 3*60*60)
	time.Local = loc

	// Define command line flags, add any other flag required to configure the
	// service.
	var (
		debugFlag = flag.Bool("debug", false, "Log request and response bodies")
	)
	flag.Parse()

	configMode := os.Getenv("CONFIG_MODE")

	conf := &config.Config{}

	fmt.Println(configMode, *debugFlag)

	err := conf.ParseConfiguration(configMode)
	if err != nil {
		log.Fatal("Failed to parser configuration: ", err)
	}

	err = logger.InitLogger(conf.Logger)
	if err != nil {
		log.Fatal("Failed to create logger: ", err)

		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	dbTXFunc, err := postgres.NewDBTxFunc(ctx, conf.Postgres)
	if err != nil {
		logger.Fatalf(ctx, "Failed to connect to database: %v", err)
	}

	txFunc, err := postgres.NewTxFunc(ctx, conf.Postgres)
	if err != nil {
		logger.Fatalf(ctx, "Failed to connect to transaction's database: %v", err)
	}

	now, _ := dbmodel.New(dbTXFunc(ctx)).SelectNow(ctx)
	logger.Infof(ctx, "db now: %s", now)

	logger.Infof(ctx, "connecting to insta-proxy at '%s'", conf.Listen.InstaProxyURL)

	conn, err := grpc.DialContext(
		ctx,
		conf.Listen.InstaProxyURL,
		grpc.WithUnaryInterceptor(mw.UnaryClientLog()),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Fatalf(ctx, "failed to connect to parser: %v", err)
	}

	queue := workers.NewQueuue(ctx, dbTXFunc(ctx), dbTXFunc, conn)

	datasetsStore := datasets.NewStore(5*time.Second, dbTXFunc, txFunc, conf.Instagrapi.Hostname, queue)
	// botsStore := bots.NewStore(dbTXFunc, txFunc)
	//
	// Initialize the services.
	// adminServiceSvc := service.NewAdminService(botsStore)

	// Initialize the services.
	datasetsService := service.NewDatasetsService(conf.Security, datasetsStore, conn)

	// potentially running in different processes.
	datasetsEndpoints := datasetsservice.NewEndpoints(datasetsService)

	// Create channel used by both the signal handler and server goroutines
	// to notify the main goroutine when to stop the server.
	errc := make(chan error, 5)

	// Setup interrupt handler. This optional step configures the process so
	// that SIGINT and SIGTERM signals cause the services to stop gracefully.
	go func() {
		c := make(chan os.Signal, 2)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	var wg sync.WaitGroup

	wg.Add(1)

	// Start the servers and send errors (if any) to the error channel.
	go handleHTTPServer(
		ctx,
		conf.Listen.BindIP,
		conf.Listen.Port,
		datasetsEndpoints,
		&wg,
		errc,
		*debugFlag,
	)

	// go handleGRPCServer(
	// 	ctx,
	// 	&url.URL{Host: fmt.Sprintf("%s:%s", conf.Listen.BindIP, conf.Listen.GRPCPort)},
	// 	adminServiceSvc,
	// 	&wg,
	// 	errc,
	// 	*debugFlag,
	// )

	// Wait for signal.
	logger.Infof(ctx, "exiting from main: (%v)", <-errc)

	// Send cancellation signal to the goroutines.
	cancel()

	wg.Wait()
	logger.Info(ctx, "exited")
}
