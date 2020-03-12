package cmd

import (
	"context"
	"database/sql"
	"flag"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mark/todo/services/pkg/protocol/grpc"
	"github.com/mark/todo/services/pkg/protocol/rest"
	v1 "github.com/mark/todo/services/pkg/service/v1"
)

// Config - server configurations including database configurations
type Config struct {
	GRPCPort string
	HTTPPort string

	DBHost string
	DBUser string
	DBPassword string
	DBSchema string
}

// RunServer run grpc server
func RunServer() error {
	ctx := context.Background()
	var cfgs Config
	flag.StringVar(&cfgs.GRPCPort, "grpc-port", "", "gRPC port to bind")
	flag.StringVar(&cfgs.HTTPPort, "http-port", "", "HTPP port to bind")
	flag.StringVar(&cfgs.DBHost, "db-host", "", "Database host")
	flag.StringVar(&cfgs.DBUser, "db-user", "", "Database user")
	flag.StringVar(&cfgs.DBPassword, "db-password", "", "Database password")
	flag.StringVar(&cfgs.DBSchema, "db-schema", "", "Database schema")
	flag.Parse()

	if len(cfgs.GRPCPort) == 0 {
		return fmt.Errorf("invalid TCP port for gRPC server: '%s'", cfgs.GRPCPort)
	}

	if len(cfgs.HTTPPort) == 0 {
		return fmt.Errorf("invalid TCP port for HTTP gateway: '%s'", cfgs.HTTPPort)
	}

	param := "parseTime=true"
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", cfgs.DBUser, cfgs.DBPassword, cfgs.DBHost, cfgs.DBSchema, param)
	dbManager, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("Failed to open database: %v", err)
	}
	V1API := v1.NewToDoServiceServer(dbManager)

	go func() {
		_ = rest.RunServer(ctx, cfgs.GRPCPort, cfgs.HTTPPort)
	}()

	return grpc.RunServer(ctx, V1API, cfgs.GRPCPort)
}