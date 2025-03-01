package controller

import (
	"Go-Gin-Basic-Template/types"
	"Go-Gin-Basic-Template/types/requestTypes"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ProductRepositoryInterface는 ProductRepository의 인터페이스입니다.
type ProductRepositoryInterface interface {
	Insert(input *requestTypes.ProductRequest) error
	Update(id string, input *requestTypes.ProductRequest) error
	Delete(id string) error
	GetAll() (*[]types.Product, error)
	GetByID(id string) (*types.Product, error)
}

// ProductRepositoryMock은 ProductRepositoryInterface의 모의 구현체입니다.
type ProductRepositoryMock struct {
	mock.Mock
}

// Insert는 ProductRepository.Insert의 모의 구현입니다.
func (m *ProductRepositoryMock) Insert(input *requestTypes.ProductRequest) error {
	args := m.Called(input)
	return args.Error(0)
}

// Update는 ProductRepository.Update의 모의 구현입니다.
func (m *ProductRepositoryMock) Update(id string, input *requestTypes.ProductRequest) error {
	args := m.Called(id, input)
	return args.Error(0)
}

// Delete는 ProductRepository.Delete의 모의 구현입니다.
func (m *ProductRepositoryMock) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// GetAll은 ProductRepository.GetAll의 모의 구현입니다.
func (m *ProductRepositoryMock) GetAll() (*[]types.Product, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]types.Product), args.Error(1)
}

// GetByID는 ProductRepository.GetByID의 모의 구현입니다.
func (m *ProductRepositoryMock) GetByID(id string) (*types.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.Product), args.Error(1)
}

// 테스트용 ProductController 구조체
type TestProductController struct {
	ProductRepository ProductRepositoryInterface
}

// Insert 메서드 구현
func (c *TestProductController) Insert(product *requestTypes.ProductRequest) (statusCode int, message string, err error) {
	err = c.ProductRepository.Insert(product)
	if err != nil {
		return http.StatusInternalServerError, "데이터베이스 저장 실패", err
	}

	return http.StatusCreated, "성공", nil
}

// Update 메서드 구현
func (c *TestProductController) Update(id string, product *requestTypes.ProductRequest) (statusCode int, message string, err error) {
	err = c.ProductRepository.Update(id, product)
	if err != nil {
		return http.StatusInternalServerError, "데이터베이스 저장 실패", err
	}

	return http.StatusOK, "성공", nil
}

// Delete 메서드 구현
func (c *TestProductController) Delete(id string) (statusCode int, message string, err error) {
	err = c.ProductRepository.Delete(id)
	if err != nil {
		return http.StatusInternalServerError, "데이터베이스 삭제 실패", err
	}

	return http.StatusOK, id, nil
}

// GetAll 메서드 구현
func (c *TestProductController) GetAll() (statusCode int, product *[]types.Product, err error) {
	product, err = c.ProductRepository.GetAll()
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, product, nil
}

// Get 메서드 구현
func (c *TestProductController) Get(id string) (statusCode int, product *types.Product, err error) {
	product, err = c.ProductRepository.GetByID(id)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, product, nil
}

func TestProductController_Insert_Success(t *testing.T) {
	// 모의 객체 설정
	mockRepo := new(ProductRepositoryMock)
	controller := &TestProductController{
		ProductRepository: mockRepo,
	}

	// 테스트 데이터
	productReq := &requestTypes.ProductRequest{
		Name:  "테스트 상품",
		Price: 10000.0,
	}

	// 모의 동작 설정
	mockRepo.On("Insert", productReq).Return(nil)

	// 테스트 실행
	statusCode, message, err := controller.Insert(productReq)

	// 검증
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, statusCode)
	assert.Equal(t, "성공", message)
	mockRepo.AssertExpectations(t)
}

func TestProductController_Insert_Failure(t *testing.T) {
	// 모의 객체 설정
	mockRepo := new(ProductRepositoryMock)
	controller := &TestProductController{
		ProductRepository: mockRepo,
	}

	// 테스트 데이터
	productReq := &requestTypes.ProductRequest{
		Name:  "테스트 상품",
		Price: 10000.0,
	}
	expectedErr := errors.New("데이터베이스 오류")

	// 모의 동작 설정
	mockRepo.On("Insert", productReq).Return(expectedErr)

	// 테스트 실행
	statusCode, message, err := controller.Insert(productReq)

	// 검증
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, http.StatusInternalServerError, statusCode)
	assert.Equal(t, "데이터베이스 저장 실패", message)
	mockRepo.AssertExpectations(t)
}

func TestProductController_Update_Success(t *testing.T) {
	// 모의 객체 설정
	mockRepo := new(ProductRepositoryMock)
	controller := &TestProductController{
		ProductRepository: mockRepo,
	}

	// 테스트 데이터
	testID := uuid.New().String()
	productReq := &requestTypes.ProductRequest{
		Name:  "업데이트된 상품",
		Price: 15000.0,
	}

	// 모의 동작 설정
	mockRepo.On("Update", testID, productReq).Return(nil)

	// 테스트 실행
	statusCode, message, err := controller.Update(testID, productReq)

	// 검증
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, "성공", message)
	mockRepo.AssertExpectations(t)
}

func TestProductController_Update_Failure(t *testing.T) {
	// 모의 객체 설정
	mockRepo := new(ProductRepositoryMock)
	controller := &TestProductController{
		ProductRepository: mockRepo,
	}

	// 테스트 데이터
	testID := uuid.New().String()
	productReq := &requestTypes.ProductRequest{
		Name:  "업데이트된 상품",
		Price: 15000.0,
	}
	expectedErr := errors.New("데이터베이스 오류")

	// 모의 동작 설정
	mockRepo.On("Update", testID, productReq).Return(expectedErr)

	// 테스트 실행
	statusCode, message, err := controller.Update(testID, productReq)

	// 검증
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, http.StatusInternalServerError, statusCode)
	assert.Equal(t, "데이터베이스 저장 실패", message)
	mockRepo.AssertExpectations(t)
}

func TestProductController_Delete_Success(t *testing.T) {
	// 모의 객체 설정
	mockRepo := new(ProductRepositoryMock)
	controller := &TestProductController{
		ProductRepository: mockRepo,
	}

	// 테스트 데이터
	testID := uuid.New().String()

	// 모의 동작 설정
	mockRepo.On("Delete", testID).Return(nil)

	// 테스트 실행
	statusCode, message, err := controller.Delete(testID)

	// 검증
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, testID, message)
	mockRepo.AssertExpectations(t)
}

func TestProductController_Delete_Failure(t *testing.T) {
	// 모의 객체 설정
	mockRepo := new(ProductRepositoryMock)
	controller := &TestProductController{
		ProductRepository: mockRepo,
	}

	// 테스트 데이터
	testID := uuid.New().String()
	expectedErr := errors.New("데이터베이스 오류")

	// 모의 동작 설정
	mockRepo.On("Delete", testID).Return(expectedErr)

	// 테스트 실행
	statusCode, message, err := controller.Delete(testID)

	// 검증
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, http.StatusInternalServerError, statusCode)
	assert.Equal(t, "데이터베이스 삭제 실패", message)
	mockRepo.AssertExpectations(t)
}

func TestProductController_GetAll_Success(t *testing.T) {
	// 모의 객체 설정
	mockRepo := new(ProductRepositoryMock)
	controller := &TestProductController{
		ProductRepository: mockRepo,
	}

	// 테스트 데이터
	testTime := time.Now()
	testProducts := &[]types.Product{
		{
			BasicModel: types.BasicModel{
				ID:       uuid.New(),
				CreateAt: testTime,
				UpdateAt: testTime,
			},
			Name:  "상품1",
			Price: 10000.0,
		},
		{
			BasicModel: types.BasicModel{
				ID:       uuid.New(),
				CreateAt: testTime,
				UpdateAt: testTime,
			},
			Name:  "상품2",
			Price: 20000.0,
		},
	}

	// 모의 동작 설정
	mockRepo.On("GetAll").Return(testProducts, nil)

	// 테스트 실행
	statusCode, products, err := controller.GetAll()

	// 검증
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, testProducts, products)
	assert.Len(t, *products, 2)
	mockRepo.AssertExpectations(t)
}

func TestProductController_GetAll_Failure(t *testing.T) {
	// 모의 객체 설정
	mockRepo := new(ProductRepositoryMock)
	controller := &TestProductController{
		ProductRepository: mockRepo,
	}

	// 테스트 데이터
	expectedErr := errors.New("데이터베이스 오류")

	// 모의 동작 설정
	mockRepo.On("GetAll").Return(nil, expectedErr)

	// 테스트 실행
	statusCode, products, err := controller.GetAll()

	// 검증
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, http.StatusInternalServerError, statusCode)
	assert.Nil(t, products)
	mockRepo.AssertExpectations(t)
}

func TestProductController_Get_Success(t *testing.T) {
	// 모의 객체 설정
	mockRepo := new(ProductRepositoryMock)
	controller := &TestProductController{
		ProductRepository: mockRepo,
	}

	// 테스트 데이터
	testID := uuid.New().String()
	testTime := time.Now()
	testProduct := &types.Product{
		BasicModel: types.BasicModel{
			ID:       uuid.MustParse(testID),
			CreateAt: testTime,
			UpdateAt: testTime,
		},
		Name:  "테스트 상품",
		Price: 10000.0,
	}

	// 모의 동작 설정
	mockRepo.On("GetByID", testID).Return(testProduct, nil)

	// 테스트 실행
	statusCode, product, err := controller.Get(testID)

	// 검증
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, testProduct, product)
	assert.Equal(t, "테스트 상품", product.Name)
	assert.Equal(t, 10000.0, product.Price)
	mockRepo.AssertExpectations(t)
}

func TestProductController_Get_Failure(t *testing.T) {
	// 모의 객체 설정
	mockRepo := new(ProductRepositoryMock)
	controller := &TestProductController{
		ProductRepository: mockRepo,
	}

	// 테스트 데이터
	testID := uuid.New().String()
	expectedErr := errors.New("데이터베이스 오류")

	// 모의 동작 설정
	mockRepo.On("GetByID", testID).Return(nil, expectedErr)

	// 테스트 실행
	statusCode, product, err := controller.Get(testID)

	// 검증
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, http.StatusInternalServerError, statusCode)
	assert.Nil(t, product)
	mockRepo.AssertExpectations(t)
} 