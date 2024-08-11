package operator

import (
	"errors"
	"github.com/ivanbulyk/clients_logistics_engine_api/internal/logistics/services/generator"
	"math"
	"math/rand"

	"github.com/ivanbulyk/clients_logistics_engine_api/internal/logistics/model"
)

// GlobalOperator that handles world and units movements
type GlobalOperator struct {
	world *model.Graph
}

// New GlobalOperator instance
func New() *GlobalOperator {
	return &GlobalOperator{
		world: model.NewGraph(),
	}
}

func (g *GlobalOperator) Populate(maxWarehouses, maxCargoUnits uint32) error {
	if uint64(maxWarehouses)+uint64(maxCargoUnits) >= 4294967295 {
		return errors.New("world actor count overflow")
	}

	generator.AddNewActors(model.Warehouses, g.world, uint(maxWarehouses), 0)
	generator.AddNewActors(model.CargoUnits, g.world, uint(maxCargoUnits), uint(maxWarehouses))

	var warehouseIDs []uint
	var deliveryUnitIDs []uint
	for _, node := range g.world.Nodes {
		if node.Type == model.Warehouses {
			warehouseIDs = append(warehouseIDs, node.ID)
		} else if node.Type == model.CargoUnits {
			deliveryUnitIDs = append(deliveryUnitIDs, node.ID)
		}
	}

	for _, warehouseID := range warehouseIDs {
		if len(deliveryUnitIDs) == 0 {
			break
		}

		rand.Shuffle(len(deliveryUnitIDs), func(i, j int) {
			deliveryUnitIDs[i], deliveryUnitIDs[j] = deliveryUnitIDs[j], deliveryUnitIDs[i]
		})

		numDeliveryUnits := rand.Intn(len(deliveryUnitIDs)) + 1 // Random number of units to connect (at least 1)
		for i := 0; i < numDeliveryUnits; i++ {
			unitID := deliveryUnitIDs[i]

			g.world.AddEdge(model.GraphEdge{Source: unitID, Target: warehouseID})
		}

		deliveryUnitIDs = deliveryUnitIDs[numDeliveryUnits:]
	}

	return nil
}

// GetDeliveryUnit from the world
func (g *GlobalOperator) GetDeliveryUnit() []*model.GraphNode {
	return g.world.GetNodesByType(model.CargoUnits)
}

// FindEntityByCoordinate in the world
func (g *GlobalOperator) FindEntityByCoordinate(coordinate model.Coordinate, entityType model.ActorType) *model.GraphNode {
	return g.world.FindNodesByLocation(coordinate, entityType)
}

// MoveDeliveryUnitToNearestWarehouse moves the given unit to the nearest connected warehouse based on their X and Y locations
func (g *GlobalOperator) MoveDeliveryUnitToNearestWarehouse(unitID uint) model.Coordinate {
	deliveryUnitNode := g.world.GetNodeByID(unitID)
	unitX := deliveryUnitNode.X
	unitY := deliveryUnitNode.Y

	connectedWarehouses := g.world.GetConnectedNodes(unitID, model.Warehouses)

	// Initialize variables for tracking the nearest warehouse
	minDistance := math.MaxFloat64
	nearestWarehouseID := uint(0)

	for _, warehouseNode := range connectedWarehouses {
		warehouseX := warehouseNode.X
		warehouseY := warehouseNode.Y

		distance := math.Sqrt(math.Pow(float64(unitX-warehouseX), 2) + math.Pow(float64(unitY-warehouseY), 2))

		// Update nearest warehouse if distance is smaller
		if distance < minDistance {
			minDistance = distance
			nearestWarehouseID = warehouseNode.ID
		}
	}

	// Move unit to goal
	if unitX < g.world.GetNodeByID(nearestWarehouseID).X {
		deliveryUnitNode.X++
	} else if unitX > g.world.GetNodeByID(nearestWarehouseID).X {
		deliveryUnitNode.X--
	}
	if unitY < g.world.GetNodeByID(nearestWarehouseID).Y {
		deliveryUnitNode.Y++
	} else if unitY > g.world.GetNodeByID(nearestWarehouseID).Y {
		deliveryUnitNode.Y--
	}

	return model.Coordinate{X: deliveryUnitNode.X, Y: deliveryUnitNode.Y}
}
