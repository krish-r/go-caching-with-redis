package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresClient struct {
	ctx  context.Context
	pool *pgxpool.Pool
}

func NewPostgresClient(ctx context.Context) (*PostgresClient, error) {
	url := getPostgresURL()

	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, err
	}

	return &PostgresClient{
		ctx:  ctx,
		pool: pool,
	}, nil
}

func (p *PostgresClient) Close() {
	p.pool.Close()
}

func (p *PostgresClient) Get(key string) (*TitleBasics, error) {
	query := `SELECT tconst, title_type, primary_title, original_title, is_adult, start_year, end_year, runtime_minutes, genres
 FROM title_basics WHERE TRIM(tconst) = $1`

	title := &TitleBasics{}
	row := p.pool.QueryRow(p.ctx, query, key)
	err := row.Scan(&title.Tconst,
		&title.TitleType,
		&title.PrimaryTitle,
		&title.OriginalTitle,
		&title.IsAdult,
		&title.StartYear,
		&title.EndYear,
		&title.RuntimeMinutes,
		&title.Genres,
	)
	if err != nil {
		return nil, err
	}
	return title, nil
}

func getPostgresURL() string {
	var (
		username string
		password string
		database string
		host     string
		port     string
	)

	if username = os.Getenv("POSTGRES_USERNAME"); username == "" {
		username = "postgres"
	}
	if password = os.Getenv("POSTGRES_PASSWORD"); password == "" {
		password = "postgres"
	}
	if database = os.Getenv("POSTGRES_DATABASE"); database == "" {
		database = "postgres"
	}
	if host = os.Getenv("POSTGRES_HOST"); host == "" {
		host = "localhost"
	}
	if port = os.Getenv("POSTGRES_PORT"); port == "" {
		port = "5432"
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)
}
