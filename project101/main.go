package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/purnaresa/golang/project101/post"
)

var postClient post.Client

func main() {

	// create the db instance using valid db connection string
	db, err := sqlx.Connect("mysql",
		"user-test:password-test@tcp(127.0.0.1:3306)/db-test?parseTime=true")
	if err != nil {
		log.Fatalln(err.Error())
	}

	// create post instance module by passing db pool
	postClient = post.NewClient(db)

	// creaete maping of routing to function
	r := gin.Default()
	r.GET("/listPost", listPost)
	r.GET("/readPost/:id", readPost)
	r.POST("/createPost", createPost)
	r.Run(":8080")
}

func listPost(c *gin.Context) {
	// read the input from url in parameter
	limit := c.Query("limit")
	limitInt, err := strconv.ParseInt(limit, 10, 64)
	if err != nil || limitInt < 1 {
		limitInt = 2
	}

	// call function to list the posts
	post, err := postClient.ListPost(limitInt)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// response with the data in list
	c.JSON(http.StatusOK, gin.H{
		"data": post,
	})

	return
}

func readPost(c *gin.Context) {
	// read the input prom url in path
	id := c.Param("id")
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil || idInt < 1 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	//  call the  function to read the post
	post, err := postClient.ReadPost(idInt)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// response with the data in object
	c.JSON(http.StatusOK, gin.H{
		"data": post,
	})

	return
}

func createPost(c *gin.Context) {
	// read the input from request body
	content := c.PostForm("content")
	if content == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "empty content",
		})
		return
	}

	// call tthe function to create the post
	postID, err := postClient.CreatePost(content)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// response with the post id
	c.JSON(http.StatusCreated, gin.H{
		"data": postID,
	})
	return
}
