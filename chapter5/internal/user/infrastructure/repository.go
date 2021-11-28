package infrastructure

import (
	"context"
	"log"

	"gorm.io/gorm"

	"github.com/PacktPublishing/Domain-Driven-Design-with-Go/chapter5/internal/user/domain"
)

// userGorm is a DTO used for communication with SQLite database
type userGorm struct {
	ID        uint   `gorm:"primaryKey;column:id"`
	Username  string `gorm:"column:username"`
	Password  string `gorm:"column:password"`
	FirstName string `gorm:"column:firstname"`
	LastName  string `gorm:"column:lastname"`
}

// newUserGorm creates new DTO from domain.User
func newUserGorm(user domain.User, password string) userGorm {
	return userGorm{
		Username:  user.Username,
		Password:  password,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
}

// TableName provides a table name for domain.User
func (userGorm) TableName() string {
	return "users"
}

// toEntity transforms DTO into domain.User
func (u userGorm) toEntity() domain.User {
	return domain.User{
		ID:        u.ID,
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}
}

// sqliteUserRepository is an actual implementation of domain.UserRepository
type sqliteUserRepository struct {
	db *gorm.DB
}

// NewUserRepository initiates new sqliteUserRepository
func NewUserRepository(db *gorm.DB) domain.UserRepository {
	err := db.AutoMigrate(userGorm{})
	if err != nil {
		log.Fatalln(err)
	}
	return &sqliteUserRepository{
		db: db,
	}
}

// Create inserts domain.User data into SQLite database
func (r *sqliteUserRepository) Create(ctx context.Context, user domain.User, password string) (*domain.User, error) {
	row := newUserGorm(user, password)
	err := r.db.WithContext(ctx).Create(&row).Error
	if err != nil {
		log.Println(err)
		return nil, err
	}

	result := row.toEntity()
	return &result, nil
}

// SearchByUsername searches all domain.User from SQLite database
func (r *sqliteUserRepository) SearchByUsername(ctx context.Context, username string) ([]domain.User, error) {
	var rows []userGorm
	err := r.db.WithContext(ctx).Where("username = ?", username).Find(&rows).Error
	if err != nil {
		log.Println(err)
		return nil, err
	}

	result := make([]domain.User, 0, len(rows))
	for _, row := range rows {
		result = append(result, row.toEntity())
	}
	return result, nil
}
