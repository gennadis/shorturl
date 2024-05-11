package repository

import (
	"context"
	"encoding/json"
	"os"
)

type FileRepository struct {
	filename string
	urls     []URL
}

func NewFileRepository(filename string) (*FileRepository, error) {
	fs := &FileRepository{
		filename: filename,
		urls:     []URL{},
	}

	if err := fs.loadData(); err != nil {
		return nil, err
	}
	return fs, nil
}

func (fr *FileRepository) loadData() error {
	file, err := os.OpenFile(fr.filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	if fileInfo.Size() == 0 {
		fr.urls = []URL{}
		return nil
	}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&fr.urls); err != nil {
		return err
	}
	return nil
}

func (fr *FileRepository) saveData() error {
	file, err := os.OpenFile(fr.filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(fr.urls); err != nil {
		return err
	}
	return nil
}

func (fr *FileRepository) AddURL(ctx context.Context, url URL) error {
	// check if the original URL already exists for any user
	for _, entry := range fr.urls {
		if entry.OriginalURL == url.OriginalURL {
			return ErrURLAlreadyExists
		}
	}

	fr.urls = append(fr.urls, url)

	if err := fr.saveData(); err != nil {
		return err
	}
	return nil
}

func (fr *FileRepository) AddURLs(ctx context.Context, urls []URL) error {
	for _, url := range urls {
		if err := fr.AddURL(ctx, url); err != nil {
			return err
		}
	}
	return nil
}

func (fr *FileRepository) GetURL(ctx context.Context, slug string) (URL, error) {
	for _, url := range fr.urls {
		if url.Slug == slug {
			return url, nil
		}
	}
	return URL{}, ErrURLNotFound
}

func (fr *FileRepository) GetURLsByUser(ctx context.Context, userID string) ([]URL, error) {
	var userURLs []URL

	for _, url := range fr.urls {
		if url.UserID == userID {
			userURLs = append(userURLs, url)
		}
	}

	if len(userURLs) == 0 {
		return nil, ErrURLNotFound
	}
	return userURLs, nil
}

func (fr *FileRepository) GetURLByOriginalURL(ctx context.Context, originalURL string) (URL, error) {
	for _, url := range fr.urls {
		if url.OriginalURL == originalURL {
			return url, nil
		}
	}
	return URL{}, ErrURLNotFound
}

func (fr *FileRepository) Ping(ctx context.Context) error {
	return nil
}
