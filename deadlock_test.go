package puddlereproduce

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestDeadlock(t *testing.T) {
	poolConfig, err := pgxpool.ParseConfig("postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		t.Fatalf("error parsing config: %v", err)
	}

	ctx := context.Background()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		t.Fatalf("error creating pool: %v", err)
	}
	defer pool.Close()
	res := pool.SendBatch(ctx, &pgx.Batch{}) // we forget to close the batch
	_ = res
	// deadlock on defer pool.Close()

	/* more realistic scenario that results in deadlock:
	batch := &pgx.Batch{}
	batch.Queue("SELECT 1;")
	// mucle memory makes you think this returns an error, but in reality it is pgx.BatchResults, which needs to be executed and closed
	err := pool.SendBatch(ctx, batch)
	if err != nil {
		t.Fatalf("error sending batch: %v", err)
	}
	*/
}
