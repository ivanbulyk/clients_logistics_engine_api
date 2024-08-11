package grpc_client

import (
	"context"
	logistics_v1 "github.com/ivanbulyk/clients_logistics_engine_api/internal/generated/logistics/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// APILogisticsClient to send requests about cargo unit movements
type APILogisticsClient struct {
	apiClientGRPC logistics_v1.LogisticsEngineAPIClient

	conn *grpc.ClientConn
}

// NewLogisticsClient instance
func NewLogisticsClient() *APILogisticsClient {
	return &APILogisticsClient{}
}

// Connect to gRPC API
func (lc *APILogisticsClient) Connect(serverAddr string, ctx context.Context) error {

	conn, dialErr := grpc.DialContext(
		ctx,
		serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if dialErr != nil {
		return dialErr
	}

	lc.conn = conn
	lc.apiClientGRPC = logistics_v1.NewLogisticsEngineAPIClient(lc.conn)

	return nil

}

// Disconnect from gRPC API
func (lc *APILogisticsClient) Disconnect() error {
	return lc.conn.Close()
}

// MoveUnit to new location
func (lc *APILogisticsClient) MoveUnit(ctx context.Context, req *logistics_v1.MoveUnitRequest) (responseErr error) {

	_, responseErr = lc.apiClientGRPC.MoveUnit(ctx, req)
	return
}

// UnitReachedWarehouse report that reach warehouse
func (lc *APILogisticsClient) UnitReachedWarehouse(ctx context.Context, req *logistics_v1.UnitReachedWarehouseRequest) (responseErr error) {

	_, responseErr = lc.apiClientGRPC.UnitReachedWarehouse(ctx, req)
	return

}
