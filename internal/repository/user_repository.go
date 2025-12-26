package repository

import (
	"fmt"
	"strings"

	"starter-kit-fullstack-gonethttp-template/internal/models"
	"starter-kit-fullstack-gonethttp-template/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *userRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	return &user, err
}

func (r *userRepository) FindAll(filters map[string]interface{}, search string, searchFields []string, pagination *utils.PaginationScope) ([]models.User, int64, error) {
	var users []models.User
	var totalRows int64

	query := r.db.Model(&models.User{})

	// 1. Apply Strict Filters (e.g. Role)
	for key, value := range filters {
		if value != "" {
			query = query.Where(fmt.Sprintf("%s = ?", key), value)
		}
	}

	// 2. Apply Fuzzy Search
	if search != "" && len(searchFields) > 0 {
		var searchConditions []string
		var searchParams []interface{}

		for _, field := range searchFields {
			if field == "name" || field == "email" || field == "role" || field == "id" {
				if r.db.Dialector.Name() == "postgres" && field == "id" {
					searchConditions = append(searchConditions, "CAST(id AS TEXT) LIKE ?")
				} else {
					searchConditions = append(searchConditions, fmt.Sprintf("%s LIKE ?", field))
				}
				searchParams = append(searchParams, "%"+search+"%")
			}
		}

		if len(searchConditions) > 0 {
			query = query.Where(strings.Join(searchConditions, " OR "), searchParams...)
		}
	}

	// 3. Count Total Rows (Before Pagination)
	query.Count(&totalRows)

	// 4. Sorting
	if pagination.Sort != "" {
		// Handle sort format "field:asc" or "field:desc"
		parts := strings.Split(pagination.Sort, ":")
		if len(parts) == 2 {
			query = query.Order(fmt.Sprintf("%s %s", parts[0], parts[1]))
		} else {
			query = query.Order("created_at desc")
		}
	} else {
		query = query.Order("created_at desc")
	}

	// 5. Pagination
	err := query.Scopes(pagination.Paginate()).Find(&users).Error
	return users, totalRows, err
}

func (r *userRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.User{}, id).Error
}