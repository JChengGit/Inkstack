package repository

import (
	"fmt"
	"inkstack-auth/internal/models"
	"strings"

	"gorm.io/gorm"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(user *models.User) error
	FindByID(id uint) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	FindByEmailOrUsername(identifier string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
	List(limit, offset int) ([]models.User, error)
	Count() (int64, error)
	ExistsByEmail(email string) (bool, error)
	ExistsByUsername(username string) (bool, error)
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create creates a new user
func (r *userRepository) Create(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// FindByID finds a user by ID
func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	return &user, nil
}

// FindByEmail finds a user by email (case-insensitive)
func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("LOWER(email) = LOWER(?)", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	return &user, nil
}

// FindByUsername finds a user by username (case-insensitive)
func (r *userRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("LOWER(username) = LOWER(?)", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	return &user, nil
}

// FindByEmailOrUsername finds a user by email or username
func (r *userRepository) FindByEmailOrUsername(identifier string) (*models.User, error) {
	var user models.User
	identifier = strings.ToLower(identifier)

	if err := r.db.Where("LOWER(email) = ? OR LOWER(username) = ?", identifier, identifier).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	return &user, nil
}

// Update updates a user
func (r *userRepository) Update(user *models.User) error {
	if err := r.db.Save(user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// Delete soft deletes a user
func (r *userRepository) Delete(id uint) error {
	if err := r.db.Delete(&models.User{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// List lists users with pagination
func (r *userRepository) List(limit, offset int) ([]models.User, error) {
	var users []models.User
	if err := r.db.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	return users, nil
}

// Count counts total users
func (r *userRepository) Count() (int64, error) {
	var count int64
	if err := r.db.Model(&models.User{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return count, nil
}

// ExistsByEmail checks if a user with given email exists
func (r *userRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.User{}).Where("LOWER(email) = LOWER(?)", email).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}
	return count > 0, nil
}

// ExistsByUsername checks if a user with given username exists
func (r *userRepository) ExistsByUsername(username string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.User{}).Where("LOWER(username) = LOWER(?)", username).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check username existence: %w", err)
	}
	return count > 0, nil
}
