package repository

import (
	"Go-Gin-Basic-Template/types"
	"Go-Gin-Basic-Template/types/requestTypes"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type ProductRepository struct {
	DB *gorm.DB
}

func (r *ProductRepository) Insert(input *requestTypes.ProductRequest) (err error) {
	dbRecord := &types.Product{
		BasicModel: types.BasicModel{
			ID:       uuid.New(),
			CreateAt: time.Now(),
		},
		Name:  input.Name,
		Price: input.Price,
	}

	if err = r.DB.Create(&dbRecord).Error; err != nil {
		return err
	}

	return nil
}

func (r *ProductRepository) Update(id string, input *requestTypes.ProductRequest) (err error) {
	dbRecord := &types.Product{}

	if err = r.DB.Where("id = ?", id).First(dbRecord).Error; err != nil {
		return err
	}

	dbRecord.Name = input.Name
	dbRecord.Price = input.Price
	dbRecord.UpdateAt = time.Now()

	if err = r.DB.Save(dbRecord).Error; err != nil {
		return err
	}

	return nil
}

func (r *ProductRepository) Delete(id string) error {
	dbRecord := &types.Product{}

	if err := r.DB.Where("id = ?", id).Delete(dbRecord).Error; err != nil {
		return err
	}

	return nil
}

func (r *ProductRepository) GetAll() (product *[]types.Product, err error) {
	if err = r.DB.Find(&product).Error; err != nil {
		return nil, err
	}

	return product, nil
}

func (r *ProductRepository) GetByID(id string) (product *[]types.Product, err error) {
	if err = r.DB.Where("id =?", id).Find(&product).Error; err != nil {
		return nil, err
	}

	return product, nil
}
