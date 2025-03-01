package repository

import (
	"Go-Gin-Basic-Template/types/requestTypes"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *gorm.DB, error) {
	// SQL 모의 객체 생성
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	// GORM 설정
	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 mockDB,
		PreferSimpleProtocol: true,
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)

	return mockDB, mock, db, err
}

func TestProductRepository_Insert(t *testing.T) {
	// 테스트 설정
	mockDB, mock, db, err := setupMockDB(t)
	require.NoError(t, err)
	defer mockDB.Close()

	repo := &ProductRepository{DB: db}

	// 테스트 데이터
	productReq := &requestTypes.ProductRequest{
		Name:  "테스트 상품",
		Price: 10000.0,
	}

	// SQL 쿼리 모의 설정
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "products"`)).
		WithArgs(
			sqlmock.AnyArg(), // ID
			sqlmock.AnyArg(), // CreateAt
			sqlmock.AnyArg(), // UpdateAt
			sqlmock.AnyArg(), // DeleteAt
			productReq.Name,
			productReq.Price,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// 테스트 실행
	err = repo.Insert(productReq)

	// 검증
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProductRepository_Update(t *testing.T) {
	// 테스트 설정
	mockDB, mock, db, err := setupMockDB(t)
	require.NoError(t, err)
	defer mockDB.Close()

	repo := &ProductRepository{DB: db}

	// 테스트 데이터
	testUUID := uuid.New()
	testIDStr := testUUID.String()
	productReq := &requestTypes.ProductRequest{
		Name:  "업데이트된 상품",
		Price: 15000.0,
	}

	// SQL 쿼리 모의 설정
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "products" WHERE id = $1 AND "products"."delete_at" IS NULL ORDER BY "products"."id" LIMIT $2`)).
		WithArgs(testIDStr, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "create_at", "update_at", "delete_at", "name", "price"}).
			AddRow(testUUID, time.Now(), time.Now(), nil, "원래 상품", 10000.0))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "products" SET`)).
		WithArgs(
			sqlmock.AnyArg(), // CreateAt
			sqlmock.AnyArg(), // UpdateAt
			sqlmock.AnyArg(), // DeleteAt
			productReq.Name,  // Name
			productReq.Price, // Price
			sqlmock.AnyArg(), // WHERE 조건의 ID
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// 테스트 실행
	err = repo.Update(testIDStr, productReq)

	// 검증
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProductRepository_Delete(t *testing.T) {
	// 테스트 설정
	mockDB, mock, db, err := setupMockDB(t)
	require.NoError(t, err)
	defer mockDB.Close()

	repo := &ProductRepository{DB: db}

	// 테스트 데이터
	testIDStr := uuid.New().String()

	// SQL 쿼리 모의 설정
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "products" SET "delete_at"=$1 WHERE id = $2 AND "products"."delete_at" IS NULL`)).
		WithArgs(sqlmock.AnyArg(), testIDStr).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// 테스트 실행
	err = repo.Delete(testIDStr)

	// 검증
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProductRepository_GetAll(t *testing.T) {
	// 테스트 설정
	mockDB, mock, db, err := setupMockDB(t)
	require.NoError(t, err)
	defer mockDB.Close()

	repo := &ProductRepository{DB: db}

	// 테스트 데이터
	testTime := time.Now()
	testUUID1 := uuid.New()
	testUUID2 := uuid.New()

	// SQL 쿼리 모의 설정
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "products" WHERE "products"."delete_at" IS NULL`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "create_at", "update_at", "delete_at", "name", "price"}).
			AddRow(testUUID1, testTime, testTime, nil, "상품1", 10000.0).
			AddRow(testUUID2, testTime, testTime, nil, "상품2", 20000.0))

	// 테스트 실행
	products, err := repo.GetAll()

	// 검증
	assert.NoError(t, err)
	assert.NotNil(t, products)
	assert.Len(t, *products, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProductRepository_GetByID(t *testing.T) {
	// 테스트 설정
	mockDB, mock, db, err := setupMockDB(t)
	require.NoError(t, err)
	defer mockDB.Close()

	repo := &ProductRepository{DB: db}

	// 테스트 데이터
	testUUID := uuid.New()
	testIDStr := testUUID.String()
	testTime := time.Now()
	deletedAt := sql.NullTime{Time: time.Now(), Valid: true}

	// SQL 쿼리 모의 설정
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "products" WHERE id =$1 AND "products"."delete_at" IS NULL`)).
		WithArgs(testIDStr).
		WillReturnRows(sqlmock.NewRows([]string{"id", "create_at", "update_at", "delete_at", "name", "price"}).
			AddRow(testUUID, testTime, testTime, deletedAt, "테스트 상품", 10000.0))

	// 테스트 실행
	product, err := repo.GetByID(testIDStr)

	// 검증
	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, testUUID, product.ID)
	assert.Equal(t, "테스트 상품", product.Name)
	assert.Equal(t, 10000.0, product.Price)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProductRepository_GetByID_NotFound(t *testing.T) {
	// 테스트 설정
	mockDB, mock, db, err := setupMockDB(t)
	require.NoError(t, err)
	defer mockDB.Close()

	repo := &ProductRepository{DB: db}

	// 테스트 데이터
	testIDStr := uuid.New().String()

	// SQL 쿼리 모의 설정
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "products" WHERE id =$1 AND "products"."delete_at" IS NULL`)).
		WithArgs(testIDStr).
		WillReturnRows(sqlmock.NewRows([]string{"id", "create_at", "update_at", "delete_at", "name", "price"}))

	// 테스트 실행
	product, err := repo.GetByID(testIDStr)

	// 검증
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.Nil(t, product)
	assert.NoError(t, mock.ExpectationsWereMet())
} 