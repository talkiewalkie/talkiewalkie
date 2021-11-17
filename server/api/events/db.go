package events

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/volatiletech/sqlboiler/v4/queries"

	"github.com/talkiewalkie/talkiewalkie/models"
)

func BatchInsert(ctx context.Context, db *sqlx.DB, slice models.EventSlice) (models.EventSlice, error) {
	q := sq.Insert(models.TableNames.Event).Columns(
		models.EventColumns.Type,
		models.EventColumns.RecipientID,
		models.EventColumns.ConversationID,
		models.EventColumns.DeletedMessageUUID,
		models.EventColumns.MessageID,
	)

	for _, p := range slice {
		q = q.Values(
			p.Type,
			p.RecipientID,
			p.ConversationID,
			p.DeletedMessageUUID,
			p.MessageID,
		)
	}

	query, args, _ := q.Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	var dbEvs models.EventSlice
	if err := queries.Raw(query, args...).Bind(ctx, db, &dbEvs); err != nil {
		return nil, err
	}

	return dbEvs, nil
}
