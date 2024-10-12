package article

import (
	"encoding/json"
	"net/http"
	"rwa/internal/db"
	"rwa/internal/handlers"
	"time"
)

func (h *ArticleHandler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Article db.Article `json:"article"`
	}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	session, ok := r.Context().Value("session").(*db.Session)
	if !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	user := session.User

	partialUser := db.User{
		Bio:       user.Bio,
		Username:  user.Username,
		CreatedAt: nil,
		UpdatedAt: nil,
	}
	article := reqBody.Article
	article.Author = partialUser
	article.CreatedAt = time.Now()
	article.UpdatedAt = time.Now()

	h.mu.Lock()
	h.articles = append(h.articles, article)
	h.mu.Unlock()

	response := map[string]interface{}{"article": article}
	handlers.JsonWrite(w, http.StatusCreated, response)
}

func (h *ArticleHandler) GetArticles(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	articlesSlice := h.FilterByAuthor(&h.articles, queryParams.Get("author"))
	articlesSlice = h.FilterByTag(&articlesSlice, queryParams.Get("tag"))

	articles := struct {
		Articles      []db.Article `json:"articles"`
		ArticlesCount uint         `json:"articlesCount"`
	}{
		Articles:      articlesSlice,
		ArticlesCount: uint(len(articlesSlice)),
	}
	handlers.JsonWrite(w, http.StatusOK, articles)
}
