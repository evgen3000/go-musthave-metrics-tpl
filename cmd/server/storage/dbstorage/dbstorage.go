package dbstorage

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBStorage struct {
	Pool *pgxpool.Pool
}

type Gauge struct {
	Id    string  `json:"id"`
	Value float64 `json:"value"`
}

type Counter struct {
	Id    string `json:"id"`
	Value int64  `json:"value"`
}

func (db *DBStorage) StorageType() string {
	return "db"
}

func (db *DBStorage) SetGauge(metricName string, value float64) {
	var gauge Gauge
	q := `SELECT id, value FROM public.gauge WHERE id = $1;`
	err := db.Pool.QueryRow(context.Background(), q, metricName).Scan(&gauge.Id, &gauge.Value)
	if errors.Is(err, pgx.ErrNoRows) {
		q = "INSERT INTO public.gauge (id, value) VALUES ($1, $2);"
		_, err = db.Pool.Exec(context.Background(), q, metricName, value)
		if err != nil {
			log.Printf("Error to create gauge %s with %v: %v", metricName, value, err)
		}
		return
	} else if err == nil {
		q = `UPDATE gauge SET gauge.value = $2 WHERE id = $1;`
		_, err = db.Pool.Exec(context.Background(), q, metricName, value)
		if err != nil {
			log.Printf("Error to update gauge %s with %v: %v", metricName, value, err)
		}
		return
	} else {
		log.Printf("Can't update gauge %s to %v: %v", metricName, value, err)
	}
}

func (db *DBStorage) IncrementCounter(metricName string, value int64) {
	var counter Counter
	q := `SELECT id, value FROM public.counter WHERE id = $1;`
	err := db.Pool.QueryRow(context.Background(), q, metricName).Scan(&counter.Id, &counter.Value)
	if errors.Is(err, pgx.ErrNoRows) {
		q = `INSERT INTO public.counter (id, value) VALUES ($1, $2);`
		_, err = db.Pool.Exec(context.Background(), q, metricName, value)
		if err != nil {
			log.Printf("Error to create gauge %s with %v: %v", metricName, value, err)
		}
		return
	} else if err == nil {
		q = `UPDATE public.counter SET value = $2 WHERE id = $1;`
		_, err = db.Pool.Exec(context.Background(), q, metricName, value+counter.Value)
		if err != nil {
			log.Printf("Error to update gauge %s with %v: %v", metricName, value, err)
		}
		return
	} else {
		log.Printf("Can't update gauge %s to %v: %v", metricName, value, err)
	}
}

func (db *DBStorage) GetGauge(metricName string) (float64, bool) {
	var gauge Gauge
	q := `SELECT id, value FROM gauge WHERE id = $1;`
	err := db.Pool.QueryRow(context.Background(), q, metricName).Scan(&gauge.Id, &gauge.Value)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, false
	} else if err != nil {
		log.Printf("Can't get gauge %s from %v: %v", metricName, gauge, err)
		return 0, false
	} else {
		return gauge.Value, true
	}
}

func (db *DBStorage) GetCounter(metricName string) (int64, bool) {
	var counter Counter
	q := `SELECT id, value FROM gauge WHERE id = $1;`
	err := db.Pool.QueryRow(context.Background(), q, metricName).Scan(&counter.Id, &counter.Value)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, false
	} else if err != nil {
		log.Printf("Can't get gauge %s from %v: %v", metricName, counter, err)
		return 0, false
	} else {
		return counter.Value, true
	}
}

func (db *DBStorage) GetAllGauges() map[string]float64 {
	var gauges = make(map[string]float64)
	q := `SELECT id, value FROM gauge;`
	rows, err := db.Pool.Query(context.Background(), q)
	if err != nil {
		log.Printf("Can't get all gauges: %v", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var gauge Gauge
		err := rows.Scan(&gauge.Id, &gauge.Value)
		if err != nil {
			log.Printf("Can't get all gauges: %v", err)
			return nil
		}
		gauges[gauge.Id] = gauge.Value
	}
	return gauges
}

func (db *DBStorage) GetAllCounters() map[string]int64 {
	var counters = make(map[string]int64)
	q := `SELECT id, value FROM counter;`
	rows, err := db.Pool.Query(context.Background(), q)
	if err != nil {
		log.Printf("Can't get all gauges: %v", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var counter Counter
		err := rows.Scan(&counter.Id, &counter.Value)
		if err != nil {
			log.Printf("Can't get all gauges: %v", err)
			return nil
		}
		counters[counter.Id] = counter.Value
	}
	return counters
}
