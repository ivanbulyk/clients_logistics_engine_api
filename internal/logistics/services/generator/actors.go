package generator

import (
	"fmt"
	"sync"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/ivanbulyk/clients_logistics_engine_api/internal/logistics/model"
)

// AddNewActors by type to the model.Graph with actorNumber and from what ID it must be added (idPrefix)
func AddNewActors(t model.ActorType, g *model.Graph, actorNumber uint, idPrefix uint) {
	locsAndRange := int(actorNumber)
	locations := NewCoordinates(locsAndRange, 255, 255)

	var wg sync.WaitGroup
	wg.Add(int(actorNumber))

	for i := uint(0); i < actorNumber; i++ {
		go func(i uint) {
			defer wg.Done()

			actorNode := model.GraphNode{ID: idPrefix + i}

			switch t {
			case model.Warehouses:
				actorNode.Name = fmt.Sprintf("Warehouse: %s - %s", gofakeit.City(), gofakeit.Company())
				actorNode.Type = model.Warehouses
			case model.CargoUnits:
				actorNode.Name = fmt.Sprintf("CargoUnit: %s - %s", gofakeit.CarMaker(), gofakeit.CarModel())
				actorNode.Type = model.CargoUnits
				actorNode.Metadata = false // Used to indicate if unit reached objective
			}

			actorNode.Coordinate = &locations[i]

			g.AddNode(actorNode)
		}(i)
	}

	wg.Wait()
}
