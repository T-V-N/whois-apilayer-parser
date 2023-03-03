package db

import (
	"context"
	"log"

	"github.com/T-V-N/whois-api-parser/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBStorage struct {
	conn *pgxpool.Pool // connection pool for performing db requests
	cfg  config.Config // config containing dsn link
}

type UpdateEntry struct {
	Domain      string
	isAvailable bool
}

func Init(cfg *config.Config) (*DBStorage, error) {
	conn, err := pgxpool.New(context.Background(), cfg.DatabaseDSN)
	if err != nil {
		log.Panic("Unable to connect to database: %v\n", err.Error())
	}

	_, err = conn.Exec(context.Background(), `
	CREATE TABLE IF NOT EXISTS 
	DOMAINS 
	(domain varchar, available varchar default null);
	`)

	if err != nil {
		log.Printf("Unable to create db: %v\n", err.Error())
		return nil, err
	}

	return &DBStorage{conn, *cfg}, nil
}

func (db *DBStorage) UpdateDomainAvailability(ctx context.Context, domain, available string) error {
	_, err := db.conn.Exec(ctx, "UPDATE domains set available = $1 WHERE domain = $2", available, domain)

	if err != nil {
		return err
	}

	return nil
}

func (db *DBStorage) GetUnproccessedDomains(ctx context.Context) ([]string, error) {
	rows, err := db.conn.Query(ctx, "Select * from DOMAINS where available is null limit 100")

	if err != nil {
		return nil, err
	}

	var domains []string

	defer rows.Close()

	for rows.Next() {
		var d string
		var b interface{}
		err = rows.Scan(&d, &b)

		if err != nil {
			return nil, err
		}

		domains = append(domains, d)
	}

	return domains, nil
}
