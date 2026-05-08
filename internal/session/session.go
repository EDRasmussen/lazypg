package session

import (
	"context"
	"fmt"
	"os"

	"charm.land/bubbles/v2/table"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Session struct {
	conn *pgx.Conn
}

func New(ctx context.Context) (*Session, error) {
	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	if err := conn.Ping(ctx); err != nil {
		conn.Close(ctx)
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return &Session{conn: conn}, nil
}

func (s *Session) Close(ctx context.Context) error {
	return s.conn.Close(ctx)
}

func (s *Session) ExecuteQuery(ctx context.Context, query string, args ...any) ([]table.Column, []table.Row, error) {
	rows, err := s.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	resCols := make([]table.Column, len(rows.FieldDescriptions()))
	for i, header := range rows.FieldDescriptions() {
		resCols[i] = table.Column{
			Title: header.Name,
			Width: len(header.Name),
		}
	}

	resRows := make([]table.Row, 0)
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, nil, err
		}

		row := make([]string, len(values))
		for i, value := range values {
			row[i] = formatValue(value)
		}

		resRows = append(resRows, row)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	return resCols, resRows, nil
}

func formatValue(value any) string {
	switch v := value.(type) {
	case nil:
		return "null"
	case [16]byte:
		return formatUUIDBytes(v[:])
	case []byte:
		if len(v) == 16 {
			return formatUUIDBytes(v)
		}
		return fmt.Sprintf("%v", v)
	case pgtype.UUID:
		if !v.Valid {
			return "null"
		}
		return formatUUIDBytes(v.Bytes[:])
	default:
		return fmt.Sprintf("%v", value)
	}
}

func formatUUIDBytes(bytes []byte) string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", bytes[0:4], bytes[4:6], bytes[6:8], bytes[8:10], bytes[10:16])
}
