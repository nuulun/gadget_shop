package repository

import (
	"account-service/internal/model"
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

var ErrNotFound = errors.New("user not found")

type userRow struct {
	ID         uint64    `gorm:"column:id;primaryKey"`
	Login      string    `gorm:"column:login;uniqueIndex;not null"`
	Email      string    `gorm:"column:email;uniqueIndex;not null"`
	Phone      string    `gorm:"column:phone"`
	FirstName  string    `gorm:"column:first_name"`
	LastName   string    `gorm:"column:last_name"`
	MiddleName string    `gorm:"column:middle_name"`
	Age        uint32    `gorm:"column:age"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (userRow) TableName() string { return "users" }

type Repository struct{ db *gorm.DB }

func New(db *gorm.DB) *Repository { return &Repository{db: db} }

func (r *Repository) Create(ctx context.Context, in model.CreateUserInput) (model.User, error) {
	row := userRow{ID: in.ID, Login: in.Login, Email: in.Email, Phone: in.Phone,
		FirstName: in.FirstName, LastName: in.LastName, MiddleName: in.MiddleName, Age: in.Age}
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return model.User{}, fmt.Errorf("create user: %w", err)
	}
	return toModel(row), nil
}

func (r *Repository) GetByID(ctx context.Context, id uint64) (model.User, error) {
	var row userRow
	err := r.db.WithContext(ctx).First(&row, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.User{}, ErrNotFound
	}
	return toModel(row), err
}

func (r *Repository) List(ctx context.Context, limit, offset int) ([]model.User, error) {
	var rows []userRow
	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&rows).Error; err != nil {
		return nil, err
	}
	users := make([]model.User, len(rows))
	for i, row := range rows {
		users[i] = toModel(row)
	}
	return users, nil
}

func (r *Repository) Update(ctx context.Context, id uint64, in model.UpdateUserInput) (model.User, error) {
	updates := map[string]interface{}{"updated_at": time.Now()}
	if in.Email != nil { updates["email"] = *in.Email }
	if in.Phone != nil { updates["phone"] = *in.Phone }
	if in.FirstName != nil { updates["first_name"] = *in.FirstName }
	if in.LastName != nil { updates["last_name"] = *in.LastName }
	if in.MiddleName != nil { updates["middle_name"] = *in.MiddleName }
	if in.Age != nil { updates["age"] = *in.Age }
	res := r.db.WithContext(ctx).Model(&userRow{}).Where("id = ?", id).Updates(updates)
	if res.Error != nil {
		return model.User{}, res.Error
	}
	if res.RowsAffected == 0 {
		return model.User{}, ErrNotFound
	}
	return r.GetByID(ctx, id)
}

func (r *Repository) Delete(ctx context.Context, id uint64) error {
	res := r.db.WithContext(ctx).Where("id = ?", id).Delete(&userRow{})
	if res.Error != nil { return res.Error }
	if res.RowsAffected == 0 { return ErrNotFound }
	return nil
}

func toModel(r userRow) model.User {
	return model.User{ID: r.ID, Login: r.Login, Email: r.Email, Phone: r.Phone,
		FirstName: r.FirstName, LastName: r.LastName, MiddleName: r.MiddleName,
		Age: r.Age, CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt}
}
