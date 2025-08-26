package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"start_golang_with_docker/internal/database"
	"start_golang_with_docker/internal/models"
	"start_golang_with_docker/internal/redis"
	"strconv"
	"strings"
	"time"
)

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func GetItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		writeError(w, http.StatusNotFound, "Invalid URL")
		return
	}

	idStr := pathParts[2]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusNotFound, "Invalid ID")
		return
	}

	cacheKey := fmt.Sprintf("item: %d", id)
	cacheData, err := redis.Get(cacheKey)
	if err == nil {
		var item models.Item
		if json.Unmarshal([]byte(cacheData), &item) == nil {
			writeJSON(w, http.StatusOK, map[string]interface{}{
				"source": "cache",
				"item":   item,
			})
			return
		}
	}

	db := database.GetDB()
	row := db.QueryRow(`SELECT id, name, price, created_at, updated_at
FROM items
WHERE id = $1`, id)

	item, err := models.ScanItem(row)
	if err != nil {
		if err == sql.ErrNoRows {
			writeError(w, http.StatusNotFound, "Item not found")
		} else {
			log.Println(err)
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	itemJSON, _ := json.Marshal(item)
	redis.Set(cacheKey, string(itemJSON), time.Hour)

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"source": "database",
		"item":   item,
	})
}

func CreateItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	var item struct {
		Name  string  `json:"name"`
		Price float64 `json:"price"`
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err = json.Unmarshal(body, &item); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if item.Name == "" || item.Price <= 0 {
		writeError(w, http.StatusBadRequest, "Invalid item name")
		return
	}

	db := database.GetDB()
	var newItem models.Item

	err = db.QueryRow(`INSERT INTO items(name, price) VALUES ($1, $2) RETURNING id, name, price, created_at, updated_at`, item.Name, item.Price).Scan(
		&newItem.ID,
		&newItem.Name,
		&newItem.Price,
		&newItem.CreateAt,
		&newItem.UpdateAt,
	)
	if err != nil {
		log.Printf("Database error: %s", err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{"item": newItem})
}

func GetItems(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	db := database.GetDB()
	rows, err := db.Query(`SELECT id, name, price, created_at, updated_at
FROM items ORDER BY created_at DESC`)
	if err != nil {
		log.Printf("Database error: %s", err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		err := rows.Scan(&item.ID,
			&item.Name,
			&item.Price,
			&item.CreateAt,
			&item.UpdateAt,
		)
		if err != nil {
			log.Printf("Database error: %s", err)
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		items = append(items, item)
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"items": items})
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"status": "ok"})
}
