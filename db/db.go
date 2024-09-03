package db

import (
	"errors"
	"fmt"
	"log"

	"database/sql"

	"github.com/assaidy/personal-blog-api/types"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// NOTE: this functions runs automatically
func init() {
	dbPath := "./data.db"

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
	insertPostQuery := "INSERT INTO posts (title, content, category, created_at, updated_at) VALUES (?, ?, ?, ?, ?)"
	insertTagsQuery := "INSERT INTO tags (name) VALUES (?)"
	insertPostTagsQuery := "INSERT INTO post_tags (post_id, tag_id) VALUES (?, ?)"

	// begin a transaction
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	// insert into posts
	res, err := tx.Exec(insertPostQuery, post.Title, post.Content, post.Category, post.CreatedAt, post.UpdatedAt)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	postId, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// insert into tags
	tagIds := []int64{}
	for _, tag := range post.Tags {
		res, err = tx.Exec(insertTagsQuery, tag)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		tagId, err := res.LastInsertId()
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		tagIds = append(tagIds, tagId)
	}

	// insert into post_tags
	for _, tagId := range tagIds {
		res, err = tx.Exec(insertPostTagsQuery, postId, tagId)
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
            title,
            content,
            category,
            created_at,
            updated_at
        FROM posts 
        WHERE id = ?`

	post := &types.Post{Id: id}

	row := db.QueryRow(query, id)
	err := row.Scan(&post.Title, &post.Content, &post.Category, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, types.NotFoundError(fmt.Errorf("no post with id %d found", id))
		}
		return nil, err
	}

	tags, err := getTagsFromPost(id)
	if err != nil {
		return nil, err
	}
	post.Tags = tags

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
	_, err = tx.Exec(updatePostQuery, post.Title, post.Content, post.Category, post.UpdatedAt)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete existing tags for the post
	deleteTagsQuery := `DELETE FROM post_tags WHERE post_id = ?`
	_, err = tx.Exec(deleteTagsQuery, post.Id)
	if err != nil {
		tx.Rollback()
		return err
	}

	tagIds := []int64{}
	insertTagsQuery := "INSERT INTO tags (name) VALUES (?)"
	for _, tag := range post.Tags {
		result, err := tx.Exec(insertTagsQuery, tag)
		if err != nil {
			tx.Rollback()
			return err
		}
		tagId, err := result.LastInsertId()
		if err != nil {
			tx.Rollback()
			return err
		}
		tagIds = append(tagIds, tagId)
	}

	// Insert new post_tags entries
	insertPostTagsQuery := "INSERT INTO post_tags (post_id, tag_id) VALUES (?, ?)"
	for _, tagId := range tagIds {
		_, err = tx.Exec(insertPostTagsQuery, post.Id, tagId)
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

func GetAllPosts() ([]types.Post, error) {
	query := `
        SELECT 
            id,
            title,
            content,
            category,
            created_at,
            updated_at
        FROM posts`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []types.Post

	for rows.Next() {
		var post types.Post
		err := rows.Scan(&post.Id, &post.Title, &post.Content, &post.Category, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tags, err := getTagsFromPost(post.Id)
		if err != nil {
			return nil, err
		}
		post.Tags = tags
		posts = append(posts, post)
	}

	return posts, nil
}

func getTagsFromPost(id int) ([]string, error) {
	query := `
        SELECT t.name
        FROM tags t
        INNER JOIN post_tags pt ON pt.tag_id = t.id
        WHERE pt.post_id = ?`

	rows, err := db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []string

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
