package database

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DBPool *pgxpool.Pool

// Configuration de la DB dans une optique de scaling
func InitDB(ctx context.Context, dbURL string) error {

	if dbURL == "" {
		return fmt.Errorf("db URL not set")
	}

	poolConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return fmt.Errorf("invalid db URL : %w", err)
	}
	numCPU := int32(runtime.NumCPU())
	poolConfig.MaxConnIdleTime = 5 * time.Minute
	poolConfig.MaxConnLifetime = 1 * time.Hour
	poolConfig.MaxConns = numCPU * 4
	poolConfig.MinConns = numCPU

	DBPool, err = pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return fmt.Errorf("unable to create connection pool : %w", err)
	}
	// Ping pour s'assurer du fonctionnement de la DB , éventuellement ajouter un context local avec timeout pour éviter un ping trop long
	err = DBPool.Ping(ctx)
	if err != nil {
		return fmt.Errorf("db cannot be reached : %w", err)
	}

	log.Println("DB Connexions pool initialized...")
	return nil
}
