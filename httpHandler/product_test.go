package httpHandler

import (
	"Go-Gin-Basic-Template/types"
	"Go-Gin-Basic-Template/types/requestTypes"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ProductControllerInterface는 ProductController의 인터페이스입니다.
type ProductControllerInterface interface {
	Insert(product *requestTypes.ProductRequest) (int, string, error)
	Update(id string, product *requestTypes.ProductRequest) (int, string, error)
	Delete(id string) (int, string, error)
	GetAll() (int, *[]types.Product, error)
	Get(id string) (int, *types.Product, error)
}

// ProductControllerMock은 ProductControllerInterface의 모의 구현체입니다.
type ProductControllerMock struct {
	mock.Mock
}

// Insert는 ProductController.Insert의 모의 구현입니다.
func (m *ProductControllerMock) Insert(product *requestTypes.ProductRequest) (int, string, error) {
	args := m.Called(product)
	return args.Int(0), args.String(1), args.Error(2)
}

// Update는 ProductController.Update의 모의 구현입니다.
func (m *ProductControllerMock) Update(id string, product *requestTypes.ProductRequest) (int, string, error) {
	args := m.Called(id, product)
	return args.Int(0), args.String(1), args.Error(2)
}

// Delete는 ProductController.Delete의 모의 구현입니다.
func (m *ProductControllerMock) Delete(id string) (int, string, error) {
	args := m.Called(id)
	return args.Int(0), args.String(1), args.Error(2)
}

// GetAll은 ProductController.GetAll의 모의 구현입니다.
func (m *ProductControllerMock) GetAll() (int, *[]types.Product, error) {
	args := m.Called()
	if args.Get(1) == nil {
		return args.Int(0), nil, args.Error(2)
	}
	return args.Int(0), args.Get(1).(*[]types.Product), args.Error(2)
}

// Get은 ProductController.Get의 모의 구현입니다.
func (m *ProductControllerMock) Get(id string) (int, *types.Product, error) {
	args := m.Called(id)
	if args.Get(1) == nil {
		return args.Int(0), nil, args.Error(2)
	}
	return args.Int(0), args.Get(1).(*types.Product), args.Error(2)
}

// 테스트용 ProductHandler 구조체
type TestProductHandler struct {
	ProductController ProductControllerInterface
}

// Insert 메서드 구현
func (h *TestProductHandler) Insert(c *gin.Context) {
	var product requestTypes.ProductRequest
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request payload", "error": err.Error()})
		return
	}

	statusCode, message, err := h.ProductController.Insert(&product)
	if err != nil {
		c.JSON(statusCode, gin.H{"message": message, "error": err.Error()})
		return
	}

	c.JSON(statusCode, gin.H{"message": message})
}

// Update 메서드 구현
func (h *TestProductHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var product requestTypes.ProductRequest
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request payload", "error": err.Error()})
		return
	}

	statusCode, message, err := h.ProductController.Update(id, &product)
	if err != nil {
		c.JSON(statusCode, gin.H{"message": message, "error": err.Error()})
		return
	}

	c.JSON(statusCode, gin.H{"message": message})
}

// Delete 메서드 구현
func (h *TestProductHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	statusCode, message, err := h.ProductController.Delete(id)
	if err != nil {
		c.JSON(statusCode, gin.H{"message": message, "error": err.Error()})
		return
	}

	c.JSON(statusCode, gin.H{"message": message})
}

// GetAll 메서드 구현
func (h *TestProductHandler) GetAll(c *gin.Context) {
	statusCode, products, err := h.ProductController.GetAll()
	if err != nil {
		c.JSON(statusCode, gin.H{"message": "SELECT 오류", "error": err.Error()})
		return
	}

	c.JSON(statusCode, *products)
}

// GetByID 메서드 구현
func (h *TestProductHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	statusCode, product, err := h.ProductController.Get(id)
	if err != nil {
		c.JSON(statusCode, gin.H{"message": "SELECT 오류", "error": err.Error()})
		return
	}
	var response []types.Product
	response = append(response, *product)

	c.JSON(statusCode, response)
}

// 테스트 설정 함수
func setupTest() (*gin.Engine, *ProductControllerMock) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockController := new(ProductControllerMock)
	return r, mockController
}

func TestProductHandler_Insert_Success(t *testing.T) {
	// 테스트 설정
	r, mockController := setupTest()
	handler := &TestProductHandler{
		ProductController: mockController,
	}

	// 라우터 설정
	r.POST("/products", handler.Insert)

	// 테스트 데이터
	productReq := requestTypes.ProductRequest{
		Name:  "테스트 상품",
		Price: 10000.0,
	}
	jsonValue, _ := json.Marshal(productReq)

	// 모의 동작 설정
	mockController.On("Insert", mock.AnythingOfType("*requestTypes.ProductRequest")).Return(http.StatusCreated, "성공", nil)

	// 테스트 요청 생성
	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 검증
	assert.Equal(t, http.StatusCreated, w.Code)
	mockController.AssertExpectations(t)

	// 응답 검증
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "성공", response["message"])
}

func TestProductHandler_Insert_InvalidPayload(t *testing.T) {
	// 테스트 설정
	r, mockController := setupTest()
	handler := &TestProductHandler{
		ProductController: mockController,
	}

	// 라우터 설정
	r.POST("/products", handler.Insert)

	// 잘못된 JSON 데이터
	invalidJSON := []byte(`{"name": "테스트 상품", "price": "invalid"}`)

	// 테스트 요청 생성
	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 검증
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 응답 검증
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Invalid request payload", response["message"])
}

func TestProductHandler_Insert_ControllerError(t *testing.T) {
	// 테스트 설정
	r, mockController := setupTest()
	handler := &TestProductHandler{
		ProductController: mockController,
	}

	// 라우터 설정
	r.POST("/products", handler.Insert)

	// 테스트 데이터
	productReq := requestTypes.ProductRequest{
		Name:  "테스트 상품",
		Price: 10000.0,
	}
	jsonValue, _ := json.Marshal(productReq)

	// 모의 동작 설정 - 컨트롤러 오류 반환
	mockController.On("Insert", mock.AnythingOfType("*requestTypes.ProductRequest")).Return(
		http.StatusInternalServerError, "데이터베이스 저장 실패", errors.New("데이터베이스 오류"))

	// 테스트 요청 생성
	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 검증
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockController.AssertExpectations(t)

	// 응답 검증
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "데이터베이스 저장 실패", response["message"])
}

func TestProductHandler_Update_Success(t *testing.T) {
	// 테스트 설정
	r, mockController := setupTest()
	handler := &TestProductHandler{
		ProductController: mockController,
	}

	// 라우터 설정
	r.PUT("/products/:id", handler.Update)

	// 테스트 데이터
	testID := uuid.New().String()
	productReq := requestTypes.ProductRequest{
		Name:  "업데이트된 상품",
		Price: 15000.0,
	}
	jsonValue, _ := json.Marshal(productReq)

	// 모의 동작 설정
	mockController.On("Update", testID, mock.AnythingOfType("*requestTypes.ProductRequest")).Return(http.StatusOK, "성공", nil)

	// 테스트 요청 생성
	req, _ := http.NewRequest("PUT", "/products/"+testID, bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 검증
	assert.Equal(t, http.StatusOK, w.Code)
	mockController.AssertExpectations(t)

	// 응답 검증
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "성공", response["message"])
}

func TestProductHandler_Delete_Success(t *testing.T) {
	// 테스트 설정
	r, mockController := setupTest()
	handler := &TestProductHandler{
		ProductController: mockController,
	}

	// 라우터 설정
	r.DELETE("/products/:id", handler.Delete)

	// 테스트 데이터
	testID := uuid.New().String()

	// 모의 동작 설정
	mockController.On("Delete", testID).Return(http.StatusOK, testID, nil)

	// 테스트 요청 생성
	req, _ := http.NewRequest("DELETE", "/products/"+testID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 검증
	assert.Equal(t, http.StatusOK, w.Code)
	mockController.AssertExpectations(t)

	// 응답 검증
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, testID, response["message"])
}

func TestProductHandler_GetAll_Success(t *testing.T) {
	// 테스트 설정
	r, mockController := setupTest()
	handler := &TestProductHandler{
		ProductController: mockController,
	}

	// 라우터 설정
	r.GET("/products", handler.GetAll)

	// 테스트 데이터
	testTime := time.Now()
	testUUID1 := uuid.New()
	testUUID2 := uuid.New()
	testProducts := &[]types.Product{
		{
			BasicModel: types.BasicModel{
				ID:       testUUID1,
				CreateAt: testTime,
				UpdateAt: testTime,
			},
			Name:  "상품1",
			Price: 10000.0,
		},
		{
			BasicModel: types.BasicModel{
				ID:       testUUID2,
				CreateAt: testTime,
				UpdateAt: testTime,
			},
			Name:  "상품2",
			Price: 20000.0,
		},
	}

	// 모의 동작 설정
	mockController.On("GetAll").Return(http.StatusOK, testProducts, nil)

	// 테스트 요청 생성
	req, _ := http.NewRequest("GET", "/products", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 검증
	assert.Equal(t, http.StatusOK, w.Code)
	mockController.AssertExpectations(t)

	// 응답 검증
	var responseData interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseData)
	assert.NoError(t, err)
	
	// 응답 구조 확인을 위한 출력
	fmt.Printf("GetAll 응답 구조: %+v\n", responseData)
	
	// 응답이 배열인지 확인
	responseArray, ok := responseData.([]interface{})
	assert.True(t, ok, "응답이 배열이어야 합니다")
	assert.Len(t, responseArray, 2)
	
	product1 := responseArray[0].(map[string]interface{})
	product2 := responseArray[1].(map[string]interface{})
	
	assert.Equal(t, "상품1", product1["Name"])
	assert.Equal(t, float64(10000), product1["Price"])
	assert.Equal(t, "상품2", product2["Name"])
	assert.Equal(t, float64(20000), product2["Price"])
}

func TestProductHandler_GetAll_Error(t *testing.T) {
	// 테스트 설정
	r, mockController := setupTest()
	handler := &TestProductHandler{
		ProductController: mockController,
	}

	// 라우터 설정
	r.GET("/products", handler.GetAll)

	// 모의 동작 설정 - 오류 반환
	mockController.On("GetAll").Return(http.StatusInternalServerError, nil, errors.New("데이터베이스 오류"))

	// 테스트 요청 생성
	req, _ := http.NewRequest("GET", "/products", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 검증
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockController.AssertExpectations(t)

	// 응답 검증
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "SELECT 오류", response["message"])
}

func TestProductHandler_GetByID_Success(t *testing.T) {
	// 테스트 설정
	r, mockController := setupTest()
	handler := &TestProductHandler{
		ProductController: mockController,
	}

	// 라우터 설정
	r.GET("/products/:id", handler.GetByID)

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
	mockController.On("Get", testID).Return(http.StatusOK, testProduct, nil)

	// 테스트 요청 생성
	req, _ := http.NewRequest("GET", "/products/"+testID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 검증
	assert.Equal(t, http.StatusOK, w.Code)
	mockController.AssertExpectations(t)

	// 응답 검증
	var responseData interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseData)
	assert.NoError(t, err)
	
	// 응답 구조 확인을 위한 출력
	fmt.Printf("GetByID 응답 구조: %+v\n", responseData)
	
	// 응답이 배열인지 확인
	responseArray, ok := responseData.([]interface{})
	assert.True(t, ok, "응답이 배열이어야 합니다")
	assert.Len(t, responseArray, 1)
	
	product := responseArray[0].(map[string]interface{})
	assert.Equal(t, "테스트 상품", product["Name"])
	assert.Equal(t, float64(10000), product["Price"])
}

func TestProductHandler_GetByID_Error(t *testing.T) {
	// 테스트 설정
	r, mockController := setupTest()
	handler := &TestProductHandler{
		ProductController: mockController,
	}

	// 라우터 설정
	r.GET("/products/:id", handler.GetByID)

	// 테스트 데이터
	testID := uuid.New().String()

	// 모의 동작 설정 - 오류 반환
	mockController.On("Get", testID).Return(http.StatusInternalServerError, nil, errors.New("데이터베이스 오류"))

	// 테스트 요청 생성
	req, _ := http.NewRequest("GET", "/products/"+testID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 검증
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockController.AssertExpectations(t)

	// 응답 검증
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "SELECT 오류", response["message"])
} 