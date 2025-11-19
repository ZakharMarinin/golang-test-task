package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(ctx context.Context, storagePath string) (*Storage, error) {
	poolConfig, err := pgxpool.ParseConfig(storagePath)
	if err != nil {
		return nil, err
	}

	db, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() {
	s.db.Close()
}

func (s *Storage) PutNumber(ctx context.Context, num int) error {
	const op = "storage.PutNumber"

	query := "INSERT INTO nums (num) VALUES ($1)"

	_, err := s.db.Exec(ctx, query, num)
	if err != nil {
		return fmt.Errorf("%s: could not store num: %w", op, err)
	}
	return nil
}

func (s *Storage) GetSlice(ctx context.Context) (numbers []int, err error) {
	const op = "storage.GetSlice"

	query := "SELECT num FROM nums"

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: could not fetch nums: %w", op, err)
	}

	for rows.Next() {
		var num int
		err = rows.Scan(&num)
		if err != nil {
			return nil, fmt.Errorf("%s: could not fetch nums: %w", op, err)
		}
		numbers = append(numbers, num)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: could not fetch nums: %w", op, err)
	}

	return numbers, nil
}
