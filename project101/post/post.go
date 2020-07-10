package post

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

// Client is instance of module, the attribute is DB instance
type Client struct {
	DB *sqlx.DB
}

// NewClient is function to create new instance of Post module.
func NewClient(db *sqlx.DB) (c Client) {
	return Client{
		DB: db,
	}
}

// Post is object structure of real-world post model
type Post struct {
	ID         int64     `json:"id"`
	Content    string    `json:"content"`
	CreateTime time.Time `json:"create_time"`
}

// CreatePost is function to create new content, require a content in string
func (c *Client) CreatePost(content string) (id int64, err error) {
	// run the insert script to db
	result, err := c.DB.Exec(`
		INSERT INTO post (content, create_time)
		VALUES (?, now())
	`,
		content)
	if err != nil {
		log.Println(err)
		return
	}

	// get the last insert id
	id, err = result.LastInsertId()
	if err != nil {
		log.Println(err)
		return
	}

	return
}

// ReadPost is function to read single post
func (c *Client) ReadPost(id int64) (post Post, err error) {
	// run the query with given id
	row := c.DB.QueryRow(`
			SELECT id, content, create_time
			FROM post
			WHERE id = ?`,
		id)

	// assign the query result to 'post' object
	err = row.Scan(&post.ID,
		&post.Content,
		&post.CreateTime)

	if err != nil {
		log.Println(err)
		return
	}

	return
}

// ListPost is function to read list of post with limit
func (c *Client) ListPost(limit int64) (posts []Post, err error) {
	// run the query with limit
	rows, err := c.DB.Query(`
			SELECT id, content, create_time
			FROM post
			ORDER BY id DESC
			LIMIT ?`,
		limit)

	if err != nil {
		log.Println(err)
		return
	}

	// loop the result and assign each row to single object
	defer rows.Close()
	for rows.Next() {
		post := Post{}

		err = rows.Scan(&post.ID,
			&post.Content,
			&post.CreateTime)

		if err != nil {
			log.Println(err)
			return
		}
		posts = append(posts, post)

	}

	return
}
