package article

import (
	"rwa/internal/db"
	"rwa/internal/handlers"
	"rwa/internal/utils"
	"sync"
)

type ArticleHandler struct {
	handlers.RwaHandler
	articles []db.Article
	mu       sync.Mutex
}

func (h *ArticleHandler) FilterByAuthor(articles *[]db.Article, author string) []db.Article {
	if author == "" {
		return *articles
	}
	var filterArticles []db.Article
	for _, a := range *articles {
		if a.Author.Username == author {
			filterArticles = append(filterArticles, a)
		}
	}
	return filterArticles
}

func (h *ArticleHandler) FilterByTag(articles *[]db.Article, tag string) []db.Article {
	if tag == "" {
		return *articles
	}
	var filterArticles []db.Article
	for _, a := range *articles {
		if utils.Contains(a.TagList, tag) {
			filterArticles = append(filterArticles, a)
		}
	}
	return filterArticles
}
