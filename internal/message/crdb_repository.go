package message

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/cockroachdb/cockroach-go/crdb"
	_ "github.com/lib/pq"
)

// CockroachDBRepository implements the Repository interface using CockroachDB.
type CockroachDBRepository struct {
	db *sql.DB
}

// NewCockroachDBRepository creates a new instance of CockroachDBRepository.
func NewCockroachDBRepository(db *sql.DB) *CockroachDBRepository {
	return &CockroachDBRepository{db: db}
}

// Add inserts a new message into the database.
func (repo *CockroachDBRepository) Add(ctx context.Context, message Message) (string, error) {
	const query = "INSERT INTO message (text) VALUES ($1) RETURNING id"

	var id string
	err := repo.db.QueryRowContext(ctx, query, message.Text).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("failed to insert message %w", err)
	}

	return id, nil
}

// GetAll retrieves all messages from the database.
func (repo *CockroachDBRepository) GetAll(ctx context.Context) ([]Message, error) {
	const query = "SELECT text FROM message"

	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query all messages %w", err)
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var message Message
		err := rows.Scan(&message.Text)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}
