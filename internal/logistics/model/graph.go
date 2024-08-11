package model

import "sync"

// Graph model
type Graph struct {
    Nodes []GraphNode
    Edges []GraphEdge
    sync.Mutex
}

// GraphNode ...
type GraphNode struct {
    ID        uint
    Name      string
    Connected bool
    Type      any
    Metadata  any
    *Coordinate
}

// GraphEdge ...
type GraphEdge struct {
    Source uint
    Target uint
}

// NewGraph instance
func NewGraph() *Graph {
    return &Graph{}
}

// AddNode to the graph
func (g *Graph) AddNode(node GraphNode) {
    g.Lock()
    defer g.Unlock()

    g.Nodes = append(g.Nodes, node)
}

// AddEdge to the graph
func (g *Graph) AddEdge(edge GraphEdge) {
    g.Lock()
    defer g.Unlock()

    g.Edges = append(g.Edges, edge)

    for i := range g.Nodes {
        if g.Nodes[i].ID == edge.Source {
            g.Nodes[i].Connected = true
        }
    }
}

// GetNodeByID returns the node with the specified ID, or nil if it is not found
func (g *Graph) GetNodeByID(nodeID uint) *GraphNode {
    for _, node := range g.Nodes {
        if node.ID == nodeID {
            return &node
        }
    }
    return nil
}

// GetNodesByType returns a slice of nodes with the specified type
func (g *Graph) GetNodesByType(nodeType any) []*GraphNode {
    var nodesByType []*GraphNode

    for _, node := range g.Nodes {
        if node.Type == nodeType {
            copyNode := node

            if !containsNode(nodesByType, &copyNode) {
                nodesByType = append(nodesByType, &copyNode)
            }
        }
    }

    return nodesByType
}

// GetConnectedNodes returns a slice of connected nodes of the given type to the node with the specified ID
func (g *Graph) GetConnectedNodes(nodeID uint, nodeType any) []*GraphNode {
    var connectedNodes []*GraphNode

    for _, edge := range g.Edges {
        if edge.Source == nodeID {
            targetNode := g.GetNodeByID(edge.Target)
            if targetNode != nil && targetNode.Type == nodeType {
                if !containsNode(connectedNodes, targetNode) {
                    connectedNodes = append(connectedNodes, targetNode)
                }
            }
        } else if edge.Target == nodeID {
            sourceNode := g.GetNodeByID(edge.Source)
            if sourceNode != nil && sourceNode.Type == nodeType {
                if !containsNode(connectedNodes, sourceNode) {
                    connectedNodes = append(connectedNodes, sourceNode)
                }
            }
        }
    }

    return connectedNodes
}

// FindNodesByLocation node in given coordinate
func (g *Graph) FindNodesByLocation(coordinate Coordinate, nodeType any) *GraphNode {
    for _, node := range g.Nodes {
        if node.Type != nodeType {
            continue
        }

        if *node.Coordinate == coordinate {
            return &node
        }
    }

    return nil
}

func containsNode(nodes []*GraphNode, node *GraphNode) bool {
    for _, n := range nodes {
        if n.ID == node.ID {
            return true
        }
    }
    return false
}
