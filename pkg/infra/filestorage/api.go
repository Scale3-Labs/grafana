package filestorage

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

type StorageName string

const (
	StorageNameGrafanaDS StorageName = "grafanads"
)

var (
	ErrRelativePath          = errors.New("path cant be relative")
	ErrNonCanonicalPath      = errors.New("path must be canonical")
	ErrPathTooLong           = errors.New("path is too long")
	ErrPathInvalid           = errors.New("path is invalid")
	ErrPathEndsWithDelimiter = errors.New("path can not end with delimiter")
	Delimiter                = "/"
)

func Path(path string, storageName StorageName) string {
	if strings.HasPrefix(path, Delimiter) {
		return fmt.Sprintf("%s%s", string(storageName), path)
	}

	return fmt.Sprintf("%s%s%s", string(storageName), Delimiter, path)
}

func belongsToStorage(path string, storageName StorageName) bool {
	return strings.HasPrefix(path, string(storageName))
}

type File struct {
	Contents []byte
	FileMetadata
}

type FileMetadata struct {
	Name       string
	FullPath   string
	MimeType   string
	Modified   time.Time
	Created    time.Time
	Size       int64
	Properties map[string]string
}

type ListFilesResponse struct {
	Files    []FileMetadata
	HasMore  bool
	LastPath string
}

type Paging struct {
	After string
	First int
}

type UpsertFileCommand struct {
	Path       string
	MimeType   string
	Contents   *[]byte
	Properties map[string]string
}

type PathFilters struct {
	allowedPrefixes []string
}

func (f *PathFilters) isAllowed(path string) bool {
	if f == nil || f.allowedPrefixes == nil {
		return true
	}

	for i := range f.allowedPrefixes {
		if strings.HasPrefix(path, f.allowedPrefixes[i]) {
			return true
		}
	}

	return false
}

type ListOptions struct {
	Recursive bool
	PathFilters
}

type FileStorage interface {
	Get(ctx context.Context, path string) (*File, error)
	Delete(ctx context.Context, path string) error
	Upsert(ctx context.Context, command *UpsertFileCommand) error

	ListFiles(ctx context.Context, path string, paging *Paging, options *ListOptions) (*ListFilesResponse, error)
	ListFolders(ctx context.Context, path string, options *ListOptions) ([]FileMetadata, error)

	CreateFolder(ctx context.Context, path string, name string) error
	DeleteFolder(ctx context.Context, path string) error
}

// Get(ctx, "/myGit/dashboards/xyz123")
// Get(ctx, "/ryansGit/dashboards/xyz124")
// Get(ctx, "/general/dashboards/xyz124")    ?? general
//   VS
// s3Storage.Get(ctx, "/dashboards/xyz123")
