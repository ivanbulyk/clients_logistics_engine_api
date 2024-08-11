package app

import (
	"context"
	"errors"
	"fmt"
	logistics_v1 "github.com/ivanbulyk/clients_logistics_engine_api/internal/generated/logistics/api/v1"
	"github.com/ivanbulyk/clients_logistics_engine_api/internal/grpc_client"
	"github.com/ivanbulyk/clients_logistics_engine_api/internal/logistics/config"
	"github.com/ivanbulyk/clients_logistics_engine_api/internal/logistics/model"
	"github.com/ivanbulyk/clients_logistics_engine_api/internal/logistics/services/operator"
	"github.com/ivanbulyk/clients_logistics_engine_api/internal/pkg/printer"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

const (
	appName = "Logistics Engine Client"

	maxWarehouses = 255
	maxCargoUnits = 1024
)

// App is instance of application
type App struct {
	ctx       context.Context
	ctxCancel context.CancelFunc

	logisticsClient *grpc_client.APILogisticsClient
	globalOperator  *operator.GlobalOperator

	maxMoveWaitNumber int
	reportTable       *printer.ASCIITablePrinter
	statistics        *model.Statistics
}

// New returns a service instance
func New(lc *grpc_client.APILogisticsClient, g *operator.GlobalOperator, cfg *config.ClientAppConfig) (*App, error) {
	log.Printf("%s, initializing...\n", appName)

	serviceCtx, serviceCtxCancel := context.WithCancel(context.Background())
	connCtx, connCtxCancel := context.WithTimeout(serviceCtx, 30*time.Second)
	defer connCtxCancel()

	log.Printf("%s, trying to connect to API - %s...\n", appName, cfg.GetCombinedAddress())
	if connErr := lc.Connect(cfg.GetCombinedAddress(), connCtx); connErr != nil {
		serviceCtxCancel()
		err := errors.New(fmt.Sprintf(
			"%s, failed to connect to API (%s), error: %v",
			appName,
			cfg.GetCombinedAddress(),
			connErr,
		))

		return nil, err
	}

	app := &App{
		ctx:       serviceCtx,
		ctxCancel: serviceCtxCancel,

		logisticsClient: lc,
		globalOperator:  g,

		maxMoveWaitNumber: 100,
		reportTable:       printer.NewASCIITablePrinter(),
		statistics: &model.Statistics{
			ExecTime: time.Now(),
			Operation: []*model.Operation{
				{Name: "MoveUnit"},
				{Name: "UnitReachedWarehouse"},
			},
		},
	}

	app.reportTable.AddHeader([]string{"Operation", "Count", "Errors"})
	worldPopulationErr := g.Populate(
		uint32(rand.Intn(maxWarehouses-10+1)+10),
		uint32(rand.Intn(maxCargoUnits-10+1)+10),
	)
	if worldPopulationErr != nil {
		return nil, worldPopulationErr
	}

	return app, nil
}

// MustRun is wrapper around run() and it panics if any error occurs.
func MustRun() {
	if err := run(); err != nil {
		panic(err)
	}
}

// Run app
func run() error {
	cfg := &config.ClientAppConfig{}
	cfg.LoadFromEnv()

	apiLogisticsClient := grpc_client.NewLogisticsClient()
	worldOperator := operator.New()
	app, err := New(apiLogisticsClient, worldOperator, cfg)
	if err != nil {
		panic(err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() { // Handle graceful shutdown
		<-signals // Wait for the signal

		log.Printf("%s, shutting down...\n", appName)

		app.ctxCancel()
		if app.logisticsClient != nil {
			_ = app.logisticsClient.Disconnect()
		}

		log.Printf("%s, stopped!\n", appName)

		os.Exit(0)
	}()

	deliveryUnits := app.globalOperator.GetDeliveryUnit()
	totalDeliveryUnits := len(deliveryUnits)

	for {
		var wg sync.WaitGroup
		unitsReachedObjective := 0

		// Check if all units reached goal
		for _, unit := range deliveryUnits {
			if unit.Metadata == true {
				unitsReachedObjective++
			}
		}

		if unitsReachedObjective == totalDeliveryUnits {
			log.Println("All delivery units reached warehouse...")
			break
		}

		for _, unit := range deliveryUnits {
			if unit.Metadata == true {
				continue
			}

			wg.Add(1)
			go app.processDelivery(unit, &wg)

		}

		wg.Wait()
	}

	for _, o := range app.statistics.Operation {
		app.reportTable.AddRow([]string{
			o.Name,
			strconv.FormatUint(o.A, 10),
			strconv.FormatUint(o.B, 10),
		})
	}

	fmt.Println("\nExecution time:", time.Since(app.statistics.ExecTime))
	fmt.Println(app.reportTable)

	return nil
}

func (a *App) processDelivery(unit *model.GraphNode, wg *sync.WaitGroup) {
	defer wg.Done()

	time.Sleep(time.Duration(a.maxMoveWaitNumber) * time.Microsecond)
	a.maxMoveWaitNumber = rand.Intn(a.maxMoveWaitNumber+1) + 1
	if a.maxMoveWaitNumber >= 1 {
		a.maxMoveWaitNumber = a.maxMoveWaitNumber >> 1
	}

	oldCoordinate := *unit.Coordinate
	newCoordinate := a.globalOperator.MoveDeliveryUnitToNearestWarehouse(unit.ID)
	unitMessage := fmt.Sprintf("%s moving to - Latitude:%d, Longitude:%d", unit.Name, newCoordinate.X, newCoordinate.Y)

	log.Println(unitMessage)

	a.statistics.Operation[0].AddA()
	moveErr := a.logisticsClient.MoveUnit(
		a.ctx,
		&logistics_v1.MoveUnitRequest{
			CargoUnitId: int64(unit.ID),
			Location: &logistics_v1.Location{
				Latitude:  uint32(newCoordinate.X),
				Longitude: uint32(newCoordinate.Y),
			},
		},
	)
	if moveErr != nil {
		log.Printf("filed to send MoveUnit %s, API error: %v\n", unitMessage, moveErr)
		a.statistics.Operation[0].AddB()

		return
	} else if newCoordinate != oldCoordinate {
		return
	}

	announcement := fmt.Sprintf("%s - Reached Objective.", unitMessage)
	warehouse := a.globalOperator.FindEntityByCoordinate(newCoordinate, model.Warehouses)
	if warehouse == nil {
		log.Printf("Warehouses not found in coordinates Latitude:%d Longitude:%d", newCoordinate.X, newCoordinate.Y)
		return
	}

	a.statistics.Operation[1].AddA()
	reachErr := a.logisticsClient.UnitReachedWarehouse(
		a.ctx,
		&logistics_v1.UnitReachedWarehouseRequest{
			Location: &logistics_v1.Location{Latitude: uint32(newCoordinate.X), Longitude: uint32(newCoordinate.Y)},
			Announcement: &logistics_v1.WarehouseAnnouncement{
				CargoUnitId: int64(unit.ID),
				WarehouseId: int64(warehouse.ID),
				Message:     announcement,
			},
		},
	)
	if reachErr != nil {
		log.Printf("filed to send UnitReachedWarehouse %s, API error: %v\n", unitMessage, moveErr)
		a.statistics.Operation[1].AddB()
		return
	}

	log.Println(announcement)
	unit.Metadata = true // Unit reached Warehouse

	return
}
