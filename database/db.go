package database

import (
	"context"
	"fmt"
	poller "kiln-projects/pollers"
	"log"
	"runtime"
	"strings"
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

	_, err := DBPool.CopyFrom(ctx, pgx.Identifier{"delegations"}, []string{"timestamp", "delegator", "amount", "level"}, pgx.CopyFromSlice(len(DelegationsList), func(i int) ([]any, error) {
		return []any{DelegationsList[i].Timestamp, DelegationsList[i].Delegator, DelegationsList[i].Amount, DelegationsList[i].Level}, nil
	}))

	if err != nil {
		return fmt.Errorf("ERR | Error inserting delegations : %v", err)
	}

	log.Printf("%d Delegations added successfully, last level %v", len(DelegationsList), DelegationsList[len(DelegationsList)-1].Level)
	return nil

}

// Par défault on récupère les informations par 10000 , les plus récents en premier
func DelegationsRetrieval(parentsContext context.Context, year int, level int64) ([]poller.Delegations, error) {
	ctx, cancel := context.WithTimeout(parentsContext, 10*time.Second)
	defer cancel()
	query := `SELECT delegator,timestamp,amount,level FROM delegations`

	queryParameters := []string{}
	args := []interface{}{}
	argID := 1

	if year != 0 {
		queryParameters = append(queryParameters, fmt.Sprintf("DATE_PART('year', timestamp) = $%d", argID))
		args = append(args, year)
		argID++
	}

	if level != 0 {
		queryParameters = append(queryParameters, fmt.Sprintf("level = $%d", argID))
		args = append(args, level)
		argID++
	}

	if len(queryParameters) > 0 {
		query += " WHERE " + strings.Join(queryParameters, " AND ")
	}

	query += " ORDER BY timestamp DESC LIMIT 1000"

	rows, err := DBPool.Query(ctx, query, args...)
	var DelegationsBulk []poller.Delegations
	if err != nil {
		return DelegationsBulk, err
	}
	defer rows.Close()

	for rows.Next() {
		var Delegation poller.Delegations
		err = rows.Scan(&Delegation.Delegator, &Delegation.Timestamp, &Delegation.Amount, &Delegation.Level)
		if err != nil {
			return DelegationsBulk, err
		}
		DelegationsBulk = append(DelegationsBulk, Delegation)
	}

	return DelegationsBulk, nil

}
