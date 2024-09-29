package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// const (
// 	host     = "127.0.0.1"
// 	port     = 5432
// 	user     = "postgres"
// 	password = "password"
// 	dbname   = "postgres"
// )

func Connect(connstring string) (*pgx.Conn, error) {
	var conn *pgx.Conn

	ctx := context.Background()
	connConfig, err := pgx.ParseConfig(connstring)
	if err != nil {
		return conn, err
	}

	conn, err = pgx.ConnectConfig(ctx, connConfig)
	if err != nil {
		return conn, err
	}

	//     defer conn.Close(ctx)

	var version string

	err = conn.QueryRow(ctx, "select version()").Scan(&version)
	if err != nil {
		return conn, err
	}

	fmt.Println(version)
	return conn, nil
}
