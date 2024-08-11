package model

import (
    "testing"
)

func TestGraph(t *testing.T) {
    graph := NewGraph()

    warehouse1 := GraphNode{ID: 1, Name: "Warehouse 1", Type: "WH"}
    warehouse2 := GraphNode{ID: 2, Name: "Warehouse 2", Type: "WH"}
    graph.AddNode(warehouse1)
    graph.AddNode(warehouse2)

    truck1 := GraphNode{ID: 3, Name: "Delivery Truck 1"}
    truck2 := GraphNode{ID: 4, Name: "Delivery Truck 2"}
    graph.AddNode(truck1)
    graph.AddNode(truck2)

    graph.AddEdge(GraphEdge{Source: 1, Target: 3})
    graph.AddEdge(GraphEdge{Source: 2, Target: 3})
    graph.AddEdge(GraphEdge{Source: 1, Target: 4})

    expectedNumNodes := 4
    if len(graph.Nodes) != expectedNumNodes {
        t.Errorf("Expected %d nodes, but got %d", expectedNumNodes, len(graph.Nodes))
    }

    expectedNumEdges := 3
    if len(graph.Edges) != expectedNumEdges {
        t.Errorf("Expected %d edges, but got %d", expectedNumEdges, len(graph.Edges))
    }

    // Check the Connected flag for each warehouse
    for _, node := range graph.Nodes {
        if node.Type == "WH" && !node.Connected {
            t.Errorf("%s should be connected, but it is not", node.Name)
        }
    }
}

func TestGetConnectedNodes(t *testing.T) {
    // Create a new graph
    graph := NewGraph()

    // Add some nodes to the graph
    nodeA := GraphNode{ID: 1, Type: "TypeA"}
    nodeB := GraphNode{ID: 2, Type: "TypeB"}
    nodeC := GraphNode{ID: 3, Type: "TypeA"}
    nodeD := GraphNode{ID: 4, Type: "TypeB"}
    graph.AddNode(nodeA)
    graph.AddNode(nodeB)
    graph.AddNode(nodeC)
    graph.AddNode(nodeD)

    // Add some edges to the graph
    edge1 := GraphEdge{Source: 1, Target: 2}
    edge2 := GraphEdge{Source: 1, Target: 3}
    edge3 := GraphEdge{Source: 4, Target: 2}
    graph.AddEdge(edge1)
    graph.AddEdge(edge2)
    graph.AddEdge(edge3)

    // Get connected nodes for node with ID 1 and Type "TypeA"
    connectedNodes := graph.GetConnectedNodes(1, "TypeA")

    // Check the number of connected nodes
    expectedNumNodes := 1
    if len(connectedNodes) != expectedNumNodes {
        t.Errorf("Expected %d connected nodes, but got %d", expectedNumNodes, len(connectedNodes))
    }

    // Check the connected node's ID
    expectedNodeID := uint(3)
    if connectedNodes[0].ID != expectedNodeID {
        t.Errorf("Expected connected node ID %d, but got %d", expectedNodeID, connectedNodes[0].ID)
    }
}
