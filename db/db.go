package db

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"database/sql"

	"github.com/assaidy/personal-blog-api/types"
	"github.com/mattn/go-sqlite3"
	// _ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// NOTE: this functions runs automatically
func init() {
	dbPath := "./db/data.db"

	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("failed to start db:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("failed to start db:", err)
	}

	if err := Migrate(db); err != nil {
		log.Fatal("failed apply migrations:", err)
	}
}

func CreatePost(post *types.Post) (int, error) {
	// begin a transaction
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	// insert into posts
	insertPostQuery := "INSERT INTO posts (title, content, category, created_at, updated_at) VALUES (?, ?, ?, ?, ?)"
	res, err := tx.Exec(insertPostQuery, post.Title, post.Content, post.Category, post.CreatedAt, post.UpdatedAt)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			tx.Rollback()
			return 0, types.AlreadyExistsError(fmt.Errorf("the title '%s' already exists", post.Title))
		}

		tx.Rollback()
		return 0, err
	}
	postId, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// insert into tags
	insertTagsQuery := "INSERT INTO tags (name, post_id) VALUES (?, ?)"
	for _, tag := range post.Tags {
		_, err := tx.Exec(insertTagsQuery, tag, postId)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return int(postId), nil
}

func GetPost(id int) (*types.Post, error) {
	query := `
        SELECT 
            p.title,
            p.content,
            p.category,
            p.created_at,
            p.updated_at,
            IFNULL(GROUP_CONCAT(t.name), '') AS tags
        FROM posts p
        LEFT JOIN tags t ON t.post_id = p.id
        WHERE id = ?
        GROUP BY p.id`

	post := &types.Post{Id: id}
	var tagsStr string

	row := db.QueryRow(query, id)
	err := row.Scan(&post.Title, &post.Content, &post.Category, &post.CreatedAt, &post.UpdatedAt, &tagsStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, types.NotFoundError(fmt.Errorf("no post with id %d found", id))
		}
		return nil, err
	}

	if tagsStr == "" {
		post.Tags = []string{}
	} else {
		post.Tags = strings.Split(tagsStr, ",")
	}

	return post, nil
}

func UpdatePost(post *types.Post) error {
	// Begin a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	updatePostQuery := `
        UPDATE posts
        SET title = ?, content = ?, category = ?, updated_at = ?
        WHERE id = ?`
	_, err = tx.Exec(updatePostQuery, post.Title, post.Content, post.Category, post.UpdatedAt, post.Id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete existing tags for the post
	deleteTagsQuery := `DELETE FROM tags WHERE post_id = ?`
	_, err = tx.Exec(deleteTagsQuery, post.Id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// insert into tags
	insertTagsQuery := "INSERT INTO tags (name, post_id) VALUES (?, ?)"
	for _, tag := range post.Tags {
		_, err := tx.Exec(insertTagsQuery, tag, post.Id)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func DeletePost(id int) error {
	query := "DELETE FROM posts WHERE id = ?"
	res, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affectedRows == 0 {
		return types.NotFoundError(fmt.Errorf("no post with id %d found", id))
	}

	return nil
}

func getTagsFromPost(id int) ([]string, error) {
	query := `
        SELECT t.name
        FROM tags t
        WHERE t.post_id = ?`

	rows, err := db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := make([]string, 0)

	for rows.Next() {
		var tag string
		err := rows.Scan(&tag)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func GetAllPostsByTerm(term string) ([]types.Post, error) {
	query := `
        SELECT DISTINCT 
            p.id,
            p.title,
            p.content,
            p.category,
            p.created_at,
            p.updated_at,
            IFNULL(GROUP_CONCAT(t.name), '') AS tags
        FROM posts p
        LEFT JOIN tags t ON t.post_id = p.id
        WHERE p.title LIKE '%' || ? || '%'
        OR p.content LIKE '%' || ? || '%'
        OR p.category LIKE '%' || ? || '%'
        OR (t.name IS NOT NULL AND t.name LIKE '%' || ? || '%')
        GROUP BY p.id`

	rows, err := db.Query(query, term, term, term, term)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]types.Post, 0)

	for rows.Next() {
		var post types.Post
		var tagsStr string
		err := rows.Scan(&post.Id, &post.Title, &post.Content, &post.Category, &post.CreatedAt, &post.UpdatedAt, &tagsStr)
		if err != nil {
			return nil, err
		}

		if tagsStr == "" {
			post.Tags = []string{}
		} else {
			post.Tags = strings.Split(tagsStr, ",")
		}

		posts = append(posts, post)
	}

	return posts, nil
}
