package database

import (
	"context"
	"fmt"
	"kiln-projects/database"
	poller "kiln-projects/pollers"
	"log"
	"runtime"
	"time"

	"github.com/jackc/pgx/v5"
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

func BulkAddingDelegations(parentsContext context.Context, DelegationsList []poller.Delegations) error {
	ctx, cancel := context.WithTimeout(parentsContext, 10*time.Second)
	defer cancel()

	_, err := database.DBPool.CopyFrom(ctx, pgx.Identifier{"delegations"}, []string{"Timestamp", "SenderAddress", "Amount", "BlockHeight"}, pgx.CopyFromSlice(len(DelegationsList), func(i int) ([]any, error) {
		return []any{DelegationsList[i].Timestamp, DelegationsList[i].Sender.Address, DelegationsList[i].Amount, DelegationsList[i].BlockHeight}, nil
	}))

	if err != nil {
		return fmt.Errorf("ERR | Error inserting delegations : %v", err)
	}

	log.Printf("%d Delegations added successfully, last BlockHeight %v", len(DelegationsList), DelegationsList[len(DelegationsList)-1].BlockHeight)
	return nil

}

// Par défault on récupère les informations par 100 , les plus récents en premier
func DelegationsRetrieval(parentsContext context.Context, year int) ([]poller.Delegations, error) {
	ctx, cancel := context.WithTimeout(parentsContext, 10*time.Second)
	defer cancel()

	var DelegationsBulk []poller.Delegations

	query := `SELECT adress,timestamp,amout,blockhaight FROM delegations WHERE YEAR(timestamp) = $1 ORDER BY timestamp LIMIT 100`

	rows, err := DBPool.Query(ctx, query, year)
	if err != nil {
		return DelegationsBulk, err
	}
	defer rows.Close()

	for rows.Next() {
		var Delegation poller.Delegations
		err = rows.Scan(Delegation.Sender.Address, Delegation.Timestamp, Delegation.Amount, Delegation.BlockHeight)
		if err != nil {
			return DelegationsBulk, err
		}
		DelegationsBulk = append(DelegationsBulk, Delegation)
	}

	return DelegationsBulk, nil

}
