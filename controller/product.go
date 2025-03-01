package controller

import (
	"Go-Gin-Basic-Template/repository"
	"Go-Gin-Basic-Template/types"
	"Go-Gin-Basic-Template/types/requestTypes"
	"net/http"
)

type ProductController struct {
	ProductRepository *repository.ProductRepository
}

func (c *ProductController) Insert(product *requestTypes.ProductRequest) (statusCode int, message string, err error) {
	err = c.ProductRepository.Insert(product)
	if err != nil {
		return http.StatusInternalServerError, "데이터베이스 저장 실패", err
	}

	return http.StatusCreated, "성공", nil
}

func (c *ProductController) Update(id string, product *requestTypes.ProductRequest) (statusCode int, message string, err error) {
	err = c.ProductRepository.Update(id, product)
	if err != nil {
		return http.StatusInternalServerError, "데이터베이스 저장 실패", err
	}

	return http.StatusOK, "성공", nil
}

func (c *ProductController) Delete(id string) (statusCode int, message string, err error) {
	err = c.ProductRepository.Delete(id)
	if err != nil {
		return http.StatusInternalServerError, "데이터베이스 삭제 실패", err
	}

	return http.StatusOK, id, nil
}

func (c *ProductController) GetAll() (statusCode int, product *[]types.Product, err error) {
	product, err = c.ProductRepository.GetAll()
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, product, nil
}

func (c *ProductController) Get(id string) (statusCode int, product *types.Product, err error) {
	product, err = c.ProductRepository.GetByID(id)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, product, nil
}
