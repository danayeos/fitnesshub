package handlers

import (
	"bytes"
	"encoding/json"
	"fitnesshub/server/dal"
	"fitnesshub/server/models"
	"fitnesshub/server/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Настройка тестовой базы данных
func setupTestDB() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=123456 dbname=mydb port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Миграция таблицы users
	db.AutoMigrate(&models.User{})
	return db, nil
}

// Очистка базы данных после тестов
func teardownTestDB(db *gorm.DB) {
	db.Exec("DELETE FROM users")
}

// Тест регистрации нового пользователя
func TestUserRegistrationHandler(t *testing.T) {
	// Настроить тестовую БД
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("[Ошибка] Не удалось подключиться к тестовой базе: %v", err)
	}
	defer teardownTestDB(db)

	// Создать DAL и сервис пользователя
	userDal := dal.NewUserDal(db)
	emailService := service.NewEmailService()
	userService := service.NewUserService(userDal, emailService)
	handler := UserRegistrationHandler(userService, emailService)

	// Создать тестовые данные пользователя
	userData := map[string]string{
		"name":            "Test User",
		"email":           "testuser@example.com",
		"password":        "password123",
		"confirmPassword": "password123",
	}

	// Преобразовать в JSON
	requestBody, err := json.Marshal(userData)
	if err != nil {
		t.Fatalf("[Ошибка] Не удалось преобразовать данные пользователя в JSON: %v", err)
	}

	// Создать HTTP-запрос
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Вызвать обработчик регистрации
	handler(w, req)

	// Проверить статус ответа
	assert.Equal(t, http.StatusCreated, w.Code)

	// Проверить, что пользователь добавлен в БД
	var createdUser models.User
	if err := db.First(&createdUser, "email = ?", "testuser@example.com").Error; err != nil {
		t.Fatalf("[Ошибка] Пользователь не найден в БД: %v", err)
	}

	// Проверить, что данные соответствуют введённым
	assert.Equal(t, "Test User", createdUser.Name)
	assert.Equal(t, "testuser@example.com", createdUser.Email)
}
