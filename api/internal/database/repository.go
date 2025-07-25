package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"errors"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	UpsertUser(ctx context.Context, user *User) error
	GetUserByGoogleID(ctx context.Context, googleID string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)

	GetAllTags(ctx context.Context) ([]Tag, error)

	GetNewsByID(ctx context.Context, id int) (*News, error)
	GetAllNews(ctx context.Context) ([]News, error)
	GetNewsByTag(ctx context.Context, tagID int) ([]News, error)
	GetTagsForNews(ctx context.Context, newsID int) ([]Tag, error)

	AddFavoriteTag(ctx context.Context, userID string, tagID int) error
	RemoveFavoriteTag(ctx context.Context, userID string, tagID int) error
	GetFavoriteTags(ctx context.Context, userID string) ([]Tag, error)
	GetFavoriteNews(ctx context.Context, userID string) ([]News, error)
}

type SQLRepository struct {
	db *sqlx.DB
}

func NewSQLRepository(db *sqlx.DB) Repository {
	return &SQLRepository{db: db}
}

func (r *SQLRepository) UpsertUser(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (google_id, email)
		VALUES ($1, $2)
		ON CONFLICT (google_id) DO UPDATE
		SET email = EXCLUDED.email
		RETURNING id, created_at
	`

	return r.db.QueryRowxContext(
		ctx,
		query,
		user.GoogleID,
		user.Email,
	).Scan(&user.ID, &user.CreatedAt)
}

func (r *SQLRepository) GetUserByGoogleID(ctx context.Context, googleID string) (*User, error) {
	user := &User{}
	query := `SELECT * FROM users WHERE google_id = $1`
	err := r.db.GetContext(ctx, user, query, googleID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return user, err
}

func (r *SQLRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	user := &User{}
	query := `SELECT * FROM users WHERE email = $1`
	err := r.db.GetContext(ctx, user, query, email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return user, err
}

func (r *SQLRepository) GetAllTags(ctx context.Context) ([]Tag, error) {
	var tags []Tag
	query := `SELECT * FROM tags`
	err := r.db.SelectContext(ctx, &tags, query)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return tags, err
}

func (r *SQLRepository) GetNewsByID(ctx context.Context, id int) (*News, error) {
	news := &News{}
	query := `SELECT * FROM news WHERE id = $1`
	err := r.db.GetContext(ctx, news, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	tags, err := r.GetTagsForNews(ctx, id)
	if err != nil {
		return nil, err
	}
	news.Tags = tags

	return news, nil
}

func (r *SQLRepository) GetNewsByTag(ctx context.Context, tagID int) ([]News, error) {
	query := `
			WITH filtered_news AS (
				SELECT DISTINCT n.id
				FROM news n
				JOIN news_tags nt ON n.id = nt.news_id
				WHERE nt.tag_id = $1
			)
			SELECT n.*, t.id AS tag_id, t.name AS tag_name
			FROM news n
			JOIN news_tags nt ON n.id = nt.news_id
			JOIN tags t ON nt.tag_id = t.id
			WHERE n.id IN (SELECT id FROM filtered_news)
			ORDER BY n.created_at DESC
    `

	var newsWithTags []NewsWithTags
	if err := r.db.SelectContext(ctx, &newsWithTags, query, tagID); err != nil {
		return nil, fmt.Errorf("failed to get news by tag: %w", err)
	}

	return combineNewsWithTags(newsWithTags), nil
}

func (r *SQLRepository) GetAllNews(ctx context.Context) ([]News, error) {
	query := `
		SELECT n.*, t.id AS tag_id, t.name AS tag_name
		FROM news n
		LEFT JOIN news_tags nt ON n.id = nt.news_id
		LEFT JOIN tags t ON nt.tag_id = t.id
		ORDER BY n.created_at DESC
	`

	var newsWithTags []NewsWithTags
	if err := r.db.SelectContext(ctx, &newsWithTags, query); err != nil {
		return nil, err
	}

	return combineNewsWithTags(newsWithTags), nil
}

func (r *SQLRepository) AddTagsToNews(ctx context.Context, newsID int, tagIDs []int) error {
	if len(tagIDs) == 0 {
		return nil
	}

	query := `INSERT INTO news_tags (news_id, tag_id) VALUES `
	valueStrings := make([]string, 0, len(tagIDs))
	valueArgs := make([]any, 0, len(tagIDs)*2)

	for i, tagID := range tagIDs {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
		valueArgs = append(valueArgs, newsID, tagID)
	}

	query += strings.Join(valueStrings, ",")
	query += " ON CONFLICT DO NOTHING"

	_, err := r.db.ExecContext(ctx, query, valueArgs...)
	return err
}

func (r *SQLRepository) GetTagsForNews(ctx context.Context, newsID int) ([]Tag, error) {
	query := `
		SELECT t.* 
		FROM tags t
		JOIN news_tags nt ON t.id = nt.tag_id
		WHERE nt.news_id = $1
	`
	var tags []Tag
	err := r.db.SelectContext(ctx, &tags, query, newsID)
	return tags, err
}

func (r *SQLRepository) AddFavoriteTag(ctx context.Context, userID string, tagID int) error {
	query := `
		INSERT INTO user_favorite_tags (user_id, tag_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`
	_, err := r.db.ExecContext(ctx, query, userID, tagID)
	return err
}

func (r *SQLRepository) RemoveFavoriteTag(ctx context.Context, userID string, tagID int) error {
	query := `DELETE FROM user_favorite_tags WHERE user_id = $1 AND tag_id = $2`
	_, err := r.db.ExecContext(ctx, query, userID, tagID)
	return err
}

func (r *SQLRepository) GetFavoriteTags(ctx context.Context, userID string) ([]Tag, error) {
	query := `
		SELECT t.* 
		FROM tags t
		JOIN user_favorite_tags uft ON t.id = uft.tag_id
		WHERE uft.user_id = $1
	`
	var tags []Tag
	err := r.db.SelectContext(ctx, &tags, query, userID)
	return tags, err
}

func (r *SQLRepository) GetFavoriteNews(ctx context.Context, userID string) ([]News, error) {
	query := `
		SELECT DISTINCT n.*, t.id AS tag_id, t.name AS tag_name
		FROM news n
		JOIN news_tags nt ON n.id = nt.news_id
		JOIN user_favorite_tags uft ON nt.tag_id = uft.tag_id
		JOIN tags t ON nt.tag_id = t.id
		WHERE uft.user_id = $1
		ORDER BY n.created_at DESC
	`

	var newsWithTags []NewsWithTags
	if err := r.db.SelectContext(ctx, &newsWithTags, query, userID); err != nil {
		return nil, err
	}

	return combineNewsWithTags(newsWithTags), nil
}

func combineNewsWithTags(newsWithTags []NewsWithTags) []News {
	newsMap := make(map[int]*News)
	for _, nwt := range newsWithTags {
		if _, exists := newsMap[nwt.ID]; !exists {
			newsMap[nwt.ID] = &News{
				ID:        nwt.ID,
				Title:     nwt.Title,
				Content:   nwt.Content,
				Author:    nwt.Author,
				CreatedAt: nwt.CreatedAt,
				Tags:      []Tag{},
			}
		}

		if nwt.TagID != nil {
			newsMap[nwt.ID].Tags = append(newsMap[nwt.ID].Tags, Tag{
				ID:   *nwt.TagID,
				Name: *nwt.TagName,
			})
		}
	}

	result := make([]News, 0, len(newsMap))
	for _, n := range newsMap {
		result = append(result, *n)
	}
	return result
}
