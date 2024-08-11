package generator

import (
	"testing"

	"github.com/ivanbulyk/clients_logistics_engine_api/internal/logistics/model"
)

func TestNewCoordinates(t *testing.T) {
	numCoordinates := 10
	xRange := 100
	yRange := 100

	coordinates := NewCoordinates(numCoordinates, xRange, yRange)

	// Check if the number of coordinates is correct
	if len(coordinates) != numCoordinates {
		t.Errorf("Expected %d coordinates, but got %d", numCoordinates, len(coordinates))
	}

	// Check for duplicate coordinates
	visited := make(map[model.Coordinate]bool)
	for _, c := range coordinates {
		if visited[c] {
			t.Errorf("Duplicate coordinate found: (%d, %d)", c.X, c.Y)
		}
		visited[c] = true

		// Check if the coordinates are within the specified range
		if c.X < 0 || c.X >= xRange || c.Y < 0 || c.Y >= yRange {
			t.Errorf("Coordinate out of range: (%d, %d)", c.X, c.Y)
		}
	}
}
