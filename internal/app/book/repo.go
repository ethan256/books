package book

import (
	"context"
	"encoding/json"
	"math/rand"
	"strconv"
	"time"

	"github.com/ethan256/books/pkg/log"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const (
	RedisDefaultExpiration     = 7200 // redis key 过期时间
	RedisDefaultRandExpiration = 600  // redis key 过期时间随机值上限
)

var (
	RedisUpdateError = errors.New("update redis cache failed")
)

type BookRepo interface {
	GetBookInfoByName(ctx context.Context, name string) (*Book, error)
	ListBooksByKind(ctx context.Context, kind int) ([]*Book, error)
	UpdateBookInfo(ctx context.Context, name string, book *Book) error
	SaveBook(ctx context.Context, book *Book) error
}

var _ BookRepo = (*bookRepo)(nil)

type bookRepo struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewBookRepo(db *gorm.DB, rdb *redis.Client) BookRepo {
	return &bookRepo{
		db:  db,
		rdb: rdb,
	}
}

// GetBookInfoByName implements BookRepo
func (b *bookRepo) GetBookInfoByName(ctx context.Context, name string) (*Book, error) {
	var book *Book
	value, err := b.rdb.Get(ctx, name).Result()
	if err == nil {
		_ = json.Unmarshal([]byte(value), book)
		return book, nil
	}

	if err := b.db.WithContext(ctx).Model(&Book{}).Where("name=?", name).First(book).Error; err != nil {
		return nil, errors.Wrapf(err, "not found this book, name: %s", name)
	}

	if err = b.updateRedisCahe(ctx, name, book); err != nil {
		log.Logger.Warn().Err(err).Str("name", name).Msg("update redis cache failed")
	}
	return book, nil
}

// ListBooksByKind implements BookRepo
func (b *bookRepo) ListBooksByKind(ctx context.Context, kind int) ([]*Book, error) {
	var books []*Book
	if err := b.db.WithContext(ctx).Model(&Book{}).Where("kind=?", kind).Find(books).Error; err != nil {
		return nil, errors.Wrap(err, "not found books")
	}
	return books, nil
}

// SaveBook implements BookRepo
func (b *bookRepo) SaveBook(ctx context.Context, book *Book) error {
	if err := b.db.WithContext(ctx).Model(&Book{}).Create(book).Error; err != nil {
		return err
	}
	return nil
}

// UpdateBookInfo implements BookRepo
func (b *bookRepo) UpdateBookInfo(ctx context.Context, name string, book *Book) error {
	var info Book
	tx := b.db.WithContext(ctx).Begin()
	if err := tx.Model(&Book{}).Where("name=?", name).Updates(&info).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	// 更新redis缓存
	return b.updateRedisCahe(ctx, name, book)
}

func (b *bookRepo) updateRedisCahe(ctx context.Context, name string, book *Book) error {
	bookBytes, err := json.Marshal(book)
	if err != nil {
		return errors.Wrap(RedisUpdateError, err.Error())
	}

	exp := rand.Intn(RedisDefaultRandExpiration) + RedisDefaultExpiration
	expiration, _ := time.ParseDuration(strconv.Itoa(exp) + "s")

	if err := b.rdb.Set(ctx, name, string(bookBytes), expiration).Err(); err != nil {
		return errors.Wrap(RedisUpdateError, err.Error())
	}
	return nil
}
