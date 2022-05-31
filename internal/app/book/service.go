package book

import "context"

type BookService interface {
	i()

	FindBookByName(ctx context.Context, name string) (*Book, error)
	ListBooksByKind(ctx context.Context, kind int) ([]*Book, error)
	UpdateBook(ctx context.Context, name string, book *Book) error
	SaveBook(ctx context.Context, book *Book) error
}

var _ BookService = (*bookService)(nil)

type bookService struct {
	repo BookRepo
}

// FindBookByName implements BookService
func (b *bookService) FindBookByName(ctx context.Context, name string) (*Book, error) {
	return b.repo.GetBookInfoByName(ctx, name)
}

// ListBooksByKind implements BookService
func (b *bookService) ListBooksByKind(ctx context.Context, kind int) ([]*Book, error) {
	return b.repo.ListBooksByKind(ctx, kind)
}

// SaveBook implements BookService
func (b *bookService) SaveBook(ctx context.Context, book *Book) error {
	return b.repo.SaveBook(ctx, book)
}

// UpdateBook implements BookService
func (b *bookService) UpdateBook(ctx context.Context, name string, book *Book) error {
	return b.repo.UpdateBookInfo(ctx, name, book)
}

func NewBookService(repo BookRepo) BookService {
	return &bookService{
		repo: repo,
	}
}

func (b *bookService) i() {}
