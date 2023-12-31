// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0
// source: tag_subscriptions.sql

package db

import (
	"context"

	"github.com/lib/pq"
)

const bindTag = `-- name: BindTag :one
INSERT INTO tag_subscriptions (tag, channel_id)
VALUES ($1, $2)
RETURNING tag, channel_id
`

type BindTagParams struct {
	Tag       string
	ChannelID string
}

func (q *Queries) BindTag(ctx context.Context, arg BindTagParams) (TagSubscription, error) {
	row := q.db.QueryRowContext(ctx, bindTag, arg.Tag, arg.ChannelID)
	var i TagSubscription
	err := row.Scan(&i.Tag, &i.ChannelID)
	return i, err
}

const getSubscriberChannels = `-- name: GetSubscriberChannels :many
SELECT 
    channels.id, channels.channel_name, channels.created_at, channels.workspace_id, channels.bot_token
FROM 
    channels
JOIN 
    tag_subscriptions 
    ON channels.id = tag_subscriptions.channel_id
WHERE 
    tag_subscriptions.tag = $1
`

func (q *Queries) GetSubscriberChannels(ctx context.Context, tag string) ([]Channel, error) {
	rows, err := q.db.QueryContext(ctx, getSubscriberChannels, tag)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Channel
	for rows.Next() {
		var i Channel
		if err := rows.Scan(
			&i.ID,
			&i.ChannelName,
			&i.CreatedAt,
			&i.WorkspaceID,
			&i.BotToken,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTagSubscriptions = `-- name: GetTagSubscriptions :many
SELECT tag, channel_id FROM tag_subscriptions
`

func (q *Queries) GetTagSubscriptions(ctx context.Context) ([]TagSubscription, error) {
	rows, err := q.db.QueryContext(ctx, getTagSubscriptions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []TagSubscription
	for rows.Next() {
		var i TagSubscription
		if err := rows.Scan(&i.Tag, &i.ChannelID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTagSubscriptionsWithChannelId = `-- name: GetTagSubscriptionsWithChannelId :many
SELECT tag, channel_id FROM tag_subscriptions WHERE channel_id = $1
`

func (q *Queries) GetTagSubscriptionsWithChannelId(ctx context.Context, channelID string) ([]TagSubscription, error) {
	rows, err := q.db.QueryContext(ctx, getTagSubscriptionsWithChannelId, channelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []TagSubscription
	for rows.Next() {
		var i TagSubscription
		if err := rows.Scan(&i.Tag, &i.ChannelID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTagSubscriptionsWithName = `-- name: GetTagSubscriptionsWithName :many
SELECT tag, channel_id FROM tag_subscriptions WHERE tag = $1
`

func (q *Queries) GetTagSubscriptionsWithName(ctx context.Context, tag string) ([]TagSubscription, error) {
	rows, err := q.db.QueryContext(ctx, getTagSubscriptionsWithName, tag)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []TagSubscription
	for rows.Next() {
		var i TagSubscription
		if err := rows.Scan(&i.Tag, &i.ChannelID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTagsOfChannel = `-- name: GetTagsOfChannel :many
SELECT 
    tags.id, tags.has_synonyms, tags.synonyms, tags.is_moderator_only, tags.is_required, tags.count, tags.name, tags.status
FROM 
    tags
JOIN 
    tag_subscriptions 
    ON tags.name = tag_subscriptions.tag
WHERE 
    tag_subscriptions.channel_id = $1
`

func (q *Queries) GetTagsOfChannel(ctx context.Context, channelID string) ([]Tag, error) {
	rows, err := q.db.QueryContext(ctx, getTagsOfChannel, channelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Tag
	for rows.Next() {
		var i Tag
		if err := rows.Scan(
			&i.ID,
			&i.HasSynonyms,
			pq.Array(&i.Synonyms),
			&i.IsModeratorOnly,
			&i.IsRequired,
			&i.Count,
			&i.Name,
			&i.Status,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const unbindTag = `-- name: UnbindTag :one
DELETE FROM tag_subscriptions
WHERE tag = $1 AND channel_id = $2
RETURNING tag, channel_id
`

type UnbindTagParams struct {
	Tag       string
	ChannelID string
}

func (q *Queries) UnbindTag(ctx context.Context, arg UnbindTagParams) (TagSubscription, error) {
	row := q.db.QueryRowContext(ctx, unbindTag, arg.Tag, arg.ChannelID)
	var i TagSubscription
	err := row.Scan(&i.Tag, &i.ChannelID)
	return i, err
}
