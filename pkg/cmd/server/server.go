package cmd

import (
	"context"
	"flag"
	"fmt"

	"github.com/mark/todo/services/pkg/mysql"
	"github.com/mark/todo/services/pkg/protocol/grpc"
	v1 "github.com/mark/todo/services/pkg/service/v1"
)

// Config - server configurations including database configurations
type Config struct {
	GRPCPort string

	DBHost string
	DBUser string
	DBPassword string
	DBSchema string
}

// RunServer run grpc server
func RunServer() error {
	ctx := context.Background()
	var dbManager mysql.Manager
	var cfgs Config
	flag.StringVar(&cfgs.GRPCPort, "grpc-port", "", "gRPC port to bind")
	flag.StringVar(&cfgs.DBHost, "db-host", "", "Database host")
	flag.StringVar(&cfgs.DBUser, "db-user", "", "Database user")
	flag.StringVar(&cfgs.DBPassword, "db-password", "", "Database password")
	flag.StringVar(&cfgs.DBSchema, "db-schema", "", "Database schema")
	flag.Parse()

	if len(cfgs.GRPCPort) == 0 {
		return fmt.Errorf("invalid TCP port for gRPC server: '%s'", cfgs.GRPCPort)
	}

	param := "parseTime=true"
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", cfgs.DBUser, cfgs.DBPassword, cfgs.DBHost, cfgs.DBSchema, param)
	err := dbManager.Open(dsn)
	if err != nil {
		fmt.Errorf("Failed to open database: %v", err)
	}
	defer dbManager.Close()
	V1API := v1.NewToDoServiceServer(dbManager.DB)
	return grpc.RunServer(ctx, V1API, cfgs.GRPCPort)
}