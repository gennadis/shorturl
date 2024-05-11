package repository

import "context"

type MemoryRepository struct {
	urls []URL
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		urls: []URL{},
	}
}

func (mr *MemoryRepository) AddURL(ctx context.Context, url URL) error {
	// check if the original URL already exists for any user
	for _, entry := range mr.urls {
		if entry.OriginalURL == url.OriginalURL {
			return ErrURLAlreadyExists
		}
	}

	mr.urls = append(mr.urls, url)
	return nil
}

func (mr *MemoryRepository) AddURLs(ctx context.Context, urls []URL) error {
	for _, url := range urls {
		if err := mr.AddURL(ctx, url); err != nil {
			return err
		}
	}
	return nil
}

func (mr *MemoryRepository) GetURL(ctx context.Context, slug string) (URL, error) {
	for _, url := range mr.urls {
		if url.Slug == slug {
			return url, nil
		}
	}
	return URL{}, ErrURLNotFound
}

func (mr *MemoryRepository) GetURLsByUser(ctx context.Context, userID string) ([]URL, error) {
	var userURLs []URL

	for _, url := range mr.urls {
		if url.UserID == userID {
			userURLs = append(userURLs, url)
		}
	}

	if len(userURLs) == 0 {
		return nil, ErrURLNotFound
	}
	return userURLs, nil
}

func (mr *MemoryRepository) GetURLByOriginalURL(ctx context.Context, originalURL string) (URL, error) {
	for _, url := range mr.urls {
		if url.OriginalURL == originalURL {
			return url, nil
		}
	}
	return URL{}, ErrURLNotFound
}

func (mr *MemoryRepository) Ping(ctx context.Context) error {
	return nil
}
