package repository

import (
	"context"
	"errors"
	"github.com/agpprastyo/career-link/internal/common/pagination"
	"github.com/agpprastyo/career-link/internal/user/entity"
	"github.com/agpprastyo/career-link/pkg/database"
	"github.com/agpprastyo/career-link/pkg/logger"
	"github.com/agpprastyo/career-link/pkg/mail"
	"github.com/agpprastyo/career-link/pkg/minio"
	"github.com/agpprastyo/career-link/pkg/redis"
	"github.com/google/uuid"
	"io"
	"time"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrDuplicate          = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("users already exists")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")
	ErrUserNotActive      = errors.New("users not active")
	ErrInvalidInput       = errors.New("invalid input")
)

// UserRepository implements user repository using PostgreSQL
type UserRepository struct {
	mail          *mail.Client
	db            *database.PostgresDB
	redis         *redis.Client
	log           *logger.Logger
	minio         *minio.Client
	verifyBaseURL string
}

// Repository defines the interface for users repository operations
type Repository interface {
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	CreateUser(ctx context.Context, users *entity.User) error
	CreateToken(ctx context.Context, userID uuid.UUID, tokenID uuid.UUID, token string, tokenType entity.TokenType, expiry time.Time) error
	GetToken(ctx context.Context, token string) (*entity.VerificationToken, error)
	UpdateUser(ctx context.Context, users *entity.User) error
	TokenExists(ctx context.Context, token string) (bool, error)
	UpdateToken(ctx context.Context, tokenID uuid.UUID, usedAt time.Time) error

	SendVerificationEmail(ctx context.Context, username, email, token string) error
	SendPostVerificationEmail(ctx context.Context, username, email string) error

	UploadAvatarFile(ctx context.Context, fileName string, fileContent io.Reader, fileSize int64, contentType string) (avatarURL string, err error)

	DeleteAvatarFile(ctx context.Context, fileName string) error

	GetUsersPaginated(ctx context.Context, paging pagination.Pagination) ([]entity.User, int, error)

	CreateAdminUser(ctx context.Context, users *entity.Admin, usr *entity.User) error

	GetAdminByUserID(ctx context.Context, userID uuid.UUID) (*entity.Admin, error)

	UpdateAdminStatus(ctx context.Context, userID uuid.UUID) (bool, error)
	SoftDeleteAdmin(ctx context.Context, userID uuid.UUID) error

	GetCompanyByID(ctx context.Context, id uuid.UUID) (*entity.Company, error)
	//StoreUserSession(ctx context.Context, userID uuid.UUID, sessionID string) error
}

// NewUserRepository creates new UserRepository
func NewUserRepository(db *database.PostgresDB, log *logger.Logger, client *mail.Client, minio *minio.Client, verifyBaseURL string) *UserRepository {
	return &UserRepository{
		mail:          client,
		verifyBaseURL: verifyBaseURL,
		db:            db,
		minio:         minio,
		log:           log,
	}
}
