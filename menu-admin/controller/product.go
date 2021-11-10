package controller

import (
	"menu_admin/model"
	"menu_admin/storage"
)

// CreateProduct create a new product
func CreateProduct(m *model.Product) error {
	return storage.DB().Create(m).Error
}

// UpdateProduct Update an existing product
func UpdateProduct(m *model.Product) error {
	return storage.DB().Save(m).Error
}

// DeleteProduct use soft delete to remove a product
func DeleteProduct(id uint) error {
	return storage.DB().Delete(&model.Product{}, id).Error
}
