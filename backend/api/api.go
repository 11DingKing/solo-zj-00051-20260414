package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type ListItem struct {
	Id   string `json:"id"`
	Item string `json:"item"`
	Done bool   `json:"done"`
}

var db *sql.DB
var err error

func SetupPostgres() {
	db, err = sql.Open("postgres", "postgres://postgres:password@postgres/todo?sslmode=disable")

	if err != nil {
		fmt.Println(err.Error())
	}

	if err = db.Ping(); err != nil {
		fmt.Println(err.Error())
	}

	log.Println("connected to postgres")
}

// CRUD: Create Read Update Delete API Format

// List all todo items with search, filter, and pagination
func TodoItems(c *gin.Context) {
	search := c.Query("search")
	done := c.Query("done")
	pageStr := c.Query("page")
	pageSizeStr := c.Query("page_size")

	page := 1
	pageSize := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			pageSize = ps
		}
	}

	offset := (page - 1) * pageSize

	query := "SELECT id, item, done FROM list WHERE 1=1"
	countQuery := "SELECT COUNT(*) FROM list WHERE 1=1"
	args := []interface{}{}
	countArgs := []interface{}{}
	argIndex := 1

	if search != "" {
		query += fmt.Sprintf(" AND item ILIKE $%d", argIndex)
		countQuery += fmt.Sprintf(" AND item ILIKE $%d", argIndex)
		args = append(args, "%"+search+"%")
		countArgs = append(countArgs, "%"+search+"%")
		argIndex++
	}

	if done != "" {
		doneBool, err := strconv.ParseBool(done)
		if err == nil {
			query += fmt.Sprintf(" AND done = $%d", argIndex)
			countQuery += fmt.Sprintf(" AND done = $%d", argIndex)
			args = append(args, doneBool)
			countArgs = append(countArgs, doneBool)
			argIndex++
		}
	}

	var total int
	err := db.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error with DB"})
		return
	}

	query += fmt.Sprintf(" ORDER BY id LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, pageSize, offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error with DB"})
		return
	}

	items := make([]ListItem, 0)

	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			item := ListItem{}
			if err := rows.Scan(&item.Id, &item.Item, &item.Done); err != nil {
				fmt.Println(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"message": "error with DB"})
				return
			}
			item.Item = strings.TrimSpace(item.Item)
			items = append(items, item)
		}
	}

	totalPages := (total + pageSize - 1) / pageSize

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
	c.JSON(http.StatusOK, gin.H{
		"items":       items,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	})
}

// Create todo item and add to DB
func CreateTodoItem(c *gin.Context) {
	item := c.Param("item")

	// Validate item
	if len(item) == 0 {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "please enter an item"})
	} else {
		// Create todo item
		var TodoItem ListItem

		TodoItem.Item = item
		TodoItem.Done = false

		// Insert item to DB
		_, err := db.Query("INSERT INTO list(item, done) VALUES($1, $2);", TodoItem.Item, TodoItem.Done)
		if err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"message": "error with DB"})

		}

		// Log message
		log.Println("created todo item", item)

		// Return success response
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
		c.JSON(http.StatusCreated, gin.H{"items": &TodoItem})
	}
}

// Update todo item
func UpdateTodoItem(c *gin.Context) {
	id := c.Param("id")
	done := c.Param("done")

	// Validate id and done
	if len(id) == 0 {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "please enter an id"})
	} else if len(done) == 0 {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "please enter a done state"})
	} else {
		// Find and update the todo item
		var exists bool
		err := db.QueryRow("SELECT * FROM list WHERE id=$1;", id).Scan(&exists)
		if err != nil && err == sql.ErrNoRows {
			fmt.Println(err.Error())
			c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
		} else {
			_, err := db.Query("UPDATE list SET done=$1 WHERE id=$2;", done, id)
			if err != nil {
				fmt.Println(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"message": "error with DB"})
			}

			// Log message
			log.Println("updated todo item", id, done)

			// Return success response
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
			c.JSON(http.StatusOK, gin.H{"message": "successfully updated todo item", "todo": id})
		}
	}
}

// Delete todo item
func DeleteTodoItem(c *gin.Context) {
	id := c.Param("id")

	// Validate id
	if len(id) == 0 {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "please enter an id"})
	} else {
		// Find and delete the todo item
		var exists bool
		err := db.QueryRow("SELECT * FROM list WHERE id=$1;", id).Scan(&exists)
		if err != nil && err == sql.ErrNoRows {
			fmt.Println(err.Error())
			c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
		} else {
			_, err = db.Query("DELETE FROM list WHERE id=$1;", id)
			if err != nil {
				fmt.Println(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"message": "error with DB"})
			}

			// Log message
			log.Println("deleted todo item", id)

			// Return success response
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
			c.JSON(http.StatusOK, gin.H{"message": "successfully deleted todo item", "todo": id})
		}
	}
}

// Add Filter API
