package events

import (
	"context"
	"database/sql"
	"fmt"
)

type EventRepository interface {
	Save(ctx context.Context, events []EventDto) error
	GetByStreamId(ctx context.Context, streamId string) ([]EventDto, error)
}

type ImmudbEventRepository struct {
	table  string
	client *sql.DB
}

func NewImmudbEventRepository(client *sql.DB, table string) *ImmudbEventRepository {
	return &ImmudbEventRepository{
		table:  table,
		client: client,
	}
}

func (i *ImmudbEventRepository) Save(ctx context.Context, events []EventDto) error {
	sqlQuery := fmt.Sprintf(`
		INSERT INTO %s (
			id,
			streamId,
			content
		) VALUES ($1, $2, $3);`, i.table)

	tx, err := i.client.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, event := range events {
		tx.ExecContext(
			ctx,
			sqlQuery,
			event.Id,
			event.StreamId,
			event.Content,
		)
	}

	commitErr := tx.Commit()
	if commitErr != nil {
		return commitErr
	}

	return nil
}

func (i *ImmudbEventRepository) GetByStreamId(ctx context.Context, streamId string) ([]EventDto, error) {
	sqlQuery := fmt.Sprintf(`
		SELECT id, streamId, content
		FROM %s
		WHERE streamId = $1;`, i.table)

	rows, err := i.client.QueryContext(ctx, sqlQuery, streamId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []EventDto

	for rows.Next() {
		var event EventDto
		err := rows.Scan(&event.Id, &event.StreamId, &event.Content)
		if err != nil {
			return nil, err
		}
		result = append(result, event)
	}

	return result, nil
}
