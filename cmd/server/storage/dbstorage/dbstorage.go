package dbstorage

import (
	"context"
	"errors"
	"log"

	"evgen3000/go-musthave-metrics-tpl.git/internal/dto"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBStorage struct {
	Pool *pgxpool.Pool
}

type Gauge struct {
	ID    string   `json:"id"`
	Value *float64 `json:"value"`
}

type Counter struct {
	ID    string `json:"id"`
	Value *int64 `json:"value"`
}

func (db *DBStorage) StorageType() string {
	return "db"
}

func (db *DBStorage) SetMetrics(metrics []dto.MetricsDTO) {
	tx, err := db.Pool.Begin(context.Background())
	if err != nil {
		log.Printf("Error starting transaction: %s", err)
		return
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(context.Background()); rollbackErr != nil {
				log.Fatalf("Unable to rollback transaction: %v", rollbackErr)
			}
		}
	}()
	for _, metric := range metrics {
		if metric.MType == dto.MetricTypeGauge && &metric.Value != nil {
			q := `INSERT INTO public.gauge (id, value)
					VALUES ($1, $2)
					ON CONFLICT (id) DO UPDATE
					SET value = excluded.value;`
			_, err = tx.Exec(context.Background(), q, metric.ID, &metric.Value)
			if err != nil {
				log.Printf("Error inserting gauge metric: %v", err)
			}
		} else if metric.MType == dto.MetricTypeCounter && &metric.Delta != nil {
			q := `INSERT INTO public.counter (id, value)
    				VALUES ($1, $2)
					ON CONFLICT (id) DO UPDATE
					SET value = public.counter.value + excluded.value;`
			_, err = tx.Exec(context.Background(), q, metric.ID, &metric.Value)
			if err != nil {
				log.Printf("Error inserting counter metric: %v", err)
			}
		} else {
			log.Printf("Unknown metric type or metric value is nil: %s, %s", metric.MType, metric.ID)
		}
	}
	err = tx.Commit(context.Background())
	if err != nil {
		log.Fatalf("Unable to commit transaction: %v", err)
	}
}

func (db *DBStorage) SetGauge(metricName string, value float64) {
	var gauge Gauge
	q := `SELECT id, value FROM public.gauge WHERE id = $1;`
	err := db.Pool.QueryRow(context.Background(), q, metricName).Scan(&gauge.ID, &gauge.Value)
	if errors.Is(err, pgx.ErrNoRows) {
		q = "INSERT INTO public.gauge (id, value) VALUES ($1, $2);"

		_, err = db.Pool.Exec(context.Background(), q, metricName, value)
		if err != nil {
			log.Printf("Error to create gauge %s with %v: %v", metricName, value, err)
		}
		return
	} else if err == nil {
		q = `UPDATE gauge SET value = $2 WHERE id = $1;`
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
	err := db.Pool.QueryRow(context.Background(), q, metricName).Scan(&counter.ID, &counter.Value)
	if errors.Is(err, pgx.ErrNoRows) {
		q = `INSERT INTO public.counter (id, value) VALUES ($1, $2);`
		_, err = db.Pool.Exec(context.Background(), q, metricName, value)
		if err != nil {
			log.Printf("Error to create gauge %s with %v: %v", metricName, value, err)
		}
		return
	} else if err == nil {
		q = `UPDATE public.counter SET value = $2 WHERE id = $1;`
		_, err = db.Pool.Exec(context.Background(), q, metricName, value+*counter.Value)
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
	q := `SELECT id, value FROM public.gauge WHERE id = $1;`
	err := db.Pool.QueryRow(context.Background(), q, metricName).Scan(&gauge.ID, &gauge.Value)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, false
	} else if err != nil {
		log.Printf("Can't get gauge %s from %v: %v", metricName, gauge, err)
		return 0, false
	} else {
		return *gauge.Value, true
	}
}

func (db *DBStorage) GetCounter(metricName string) (int64, bool) {
	var counter Counter
	q := `SELECT id, value FROM public.counter WHERE id = $1;`
	err := db.Pool.QueryRow(context.Background(), q, metricName).Scan(&counter.ID, &counter.Value)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, false
	} else if err != nil {
		log.Printf("Can't get gauge %s from %v: %v", metricName, counter, err)
		return 0, false
	} else {
		return *counter.Value, true
	}
}

func (db *DBStorage) GetAllGauges() map[string]float64 {
	var gauges = make(map[string]float64)
	q := `SELECT id, value FROM public.gauge;`
	rows, err := db.Pool.Query(context.Background(), q)
	if err != nil {
		log.Printf("Can't get all gauges: %v", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var gauge Gauge
		err := rows.Scan(&gauge.ID, &gauge.Value)
		if err != nil {
			log.Printf("Can't get all gauges: %v", err)
			return nil
		}
		gauges[gauge.ID] = *gauge.Value
	}
	return gauges
}

func (db *DBStorage) GetAllCounters() map[string]int64 {
	var counters = make(map[string]int64)
	q := `SELECT id, value FROM public.counter;`
	rows, err := db.Pool.Query(context.Background(), q)
	if err != nil {
		log.Printf("Can't get all gauges: %v", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {

		var counter Counter
		err := rows.Scan(&counter.ID, &counter.Value)

		if err != nil {
			log.Printf("Can't get all gauges: %v", err)
			return nil
		}

		counters[counter.ID] = *counter.Value
	}
	return counters

}
