package postgres

import (
	"L0/config"
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	queryCreateTableItems = `CREATE TABLE IF NOT EXISTS items (
		id SERIAL PRIMARY KEY,
		chrt_id 		INTEGER,
    	track_number 	VARCHAR(50),
    	price 			INTEGER,
    	rid 			VARCHAR(50),
    	name 			VARCHAR(50),
    	sale 			INTEGER,
    	size 			VARCHAR(50),
    	total_price 	INTEGER,
    	nm_id 			INTEGER,
    	brand 			VARCHAR(50),
    	status 			INTEGER
	  );`
	queryCreateTablePaymant = `CREATE TABLE IF NOT EXISTS paymant (
		id SERIAL PRIMARY KEY,
		transact 		VARCHAR(50),
    	request_id 		VARCHAR(50),
    	currency 		VARCHAR(50),
    	provider 		VARCHAR(50),
    	amount 			INTEGER,
    	payment_dt 		INTEGER,
    	bank 			VARCHAR(50),
    	delivery_cost 	INTEGER,
    	goods_total 	INTEGER,
    	custom_fee 		INTEGER
	  );`
	queryCreateTableDelivery = `CREATE TABLE IF NOT EXISTS delivery (
		id SERIAL PRIMARY KEY,
		name 	VARCHAR(50),
		phone 	VARCHAR(50),
		zip 	VARCHAR(50),
		city 	VARCHAR(50),
		address VARCHAR(50),
		region 	VARCHAR(50),
		email 	VARCHAR(50)
	  );`
	queryCreateTableOrder = `CREATE TABLE IF NOT EXISTS order_us (
		order_uid 			VARCHAR(50),
		track_number 		VARCHAR(50),
		entry 				VARCHAR(50),
		deliveryId 			INTEGER,
		paymantId 			INTEGER,
		itemsId 			INTEGER[],
		locale 				VARCHAR(50),
	    internal_signature 	VARCHAR(50),
	    customer_id 		VARCHAR(50),
	    delivery_service 	VARCHAR(50),
	    shardkey 			VARCHAR(50),
	    sm_id 				INTEGER,
	    date_created 		VARCHAR(50),
	    oof_shard 			VARCHAR(50)
	  );`
	queryInsert = `INSERT INTO social_connections ( sender, recipient, number_of_communications)
					VALUES ($1, $2, 0);`
)

func InitPsqlDB(ctx context.Context, cfg *config.Config) (*sqlx.DB, error) {
	connectionURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DBName,
		cfg.Postgres.SSLMode,
	)
	database, err := sqlx.Open("postgres", connectionURL)
	if err != nil {
		return nil, err
	}
	
	if err = database.Ping(); err != nil {
		return nil, err
	}
	database.MustExec(queryCreateTablePaymant)
	database.MustExec(queryCreateTableItems)
	database.MustExec(queryCreateTableDelivery)
	database.MustExec(queryCreateTableOrder)

	return database, nil
}