package repository

import "gorm.io/gorm"

type Repository[T any] struct {
}

func (e Repository[T]) Save(tx *gorm.DB, entity *T) error {
	return tx.Save(entity).Error
}

func (e Repository[T]) Delete(tx *gorm.DB, entity *T) error {
	return tx.Delete(entity).Error
}

func (e Repository[T]) CountByID(tx *gorm.DB, id any) (int64, error) {
	var count int64
	err := tx.Model(new(T)).Where("id = ?", id).Count(&count).Error

	return count, err
}

func (e Repository[T]) FindByID(tx *gorm.DB, entity *T, id any) error {
	return tx.First(entity, "id = ?", id).Error
}
