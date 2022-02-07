package cartIt

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	main "github.com/vincentandr/shopping-microservice/cmd/cart"
	db "github.com/vincentandr/shopping-microservice/cmd/cart/internal/db"
	catalogGrpc "github.com/vincentandr/shopping-microservice/internal/grpc/catalog"
	paymentGrpc "github.com/vincentandr/shopping-microservice/internal/grpc/payment"
	catalogpb "github.com/vincentandr/shopping-microservice/internal/proto/catalog"
	paymentpb "github.com/vincentandr/shopping-microservice/internal/proto/payment"
	rds "github.com/vincentandr/shopping-microservice/internal/redis"
)

// Run integration test using docker-compose test

type TestServer struct {
	catalogClient catalogpb.CatalogServiceClient
	paymentClient paymentpb.PaymentServiceClient
	repo *db.Repository
}

type ITCartSuite struct {
	suite.Suite
	srv *main.Server
	rdb *rds.Redis
	catalogConn *catalogGrpc.CatalogGrpc
	paymentConn *paymentGrpc.PaymentGrpc
}

func (s *ITCartSuite) SetupSuite() {
	idx, _ := strconv.Atoi(os.Getenv("REDIS_DB_INDEX"))
	// New db connection
	rdb := rds.NewDb(idx)

	// Database repo
	repo := db.NewRepository(rdb.Conn)

	// GRPC clients
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	
	catalogRpc, err := catalogGrpc.NewGrpcClient(ctx)
	if err != nil {
		fmt.Println(err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	paymentRpc, err := paymentGrpc.NewGrpcClient(ctx)
	if err != nil {
		fmt.Println(err)
	}

	s.srv = &main.Server{CatalogClient: catalogRpc.Client, PaymentClient: paymentRpc.Client, Repo: repo}
	s.rdb = rdb
	s.catalogConn = catalogRpc
	s.paymentConn = paymentRpc
}

func (s *ITCartSuite) TearDownSuite() {
	s.rdb.Close()
	s.catalogConn.Close()
	s.paymentConn.Close()
}

func (s *ITCartSuite) TestGepc_GetCartItems(t *testing.T) {
	// Add initial items
	for i := 1; i<3; i++ {
		_, err := s.srv.Repo.AddOrUpdateCart(context.Background(), "user1", fmt.Sprintf("productid%d",i), i)
		if err != nil {
			t.Fatalf("failed to add initial items: %v", err)
		}
	}
}