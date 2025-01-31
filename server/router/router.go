package router

import (
	"fitnesshub/server/dal"
	"fitnesshub/server/handlers"
	"fitnesshub/server/middleware"
	"fitnesshub/server/service"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// NewRouter создает новый маршрутизатор и настраивает маршруты
func NewRouter(db *gorm.DB) *mux.Router {
	// Создаем DAL (Data Access Layer)
	userDal := dal.NewUserDal(db)
	emailService := service.NewEmailService() // Используем конструктор без аргументов

	adminDal := dal.NewUserDal(db)
	adminService := service.NewAdminService(adminDal)

	// Создаем Service, передавая DAL
	userService := service.NewUserService(userDal, emailService)

	// Инициализируем роутер
	r := mux.NewRouter()

	// Регистрируем маршруты для статичных страниц
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "pages/index.html")
	}).Methods("GET")
	r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "pages/signup.html")
	}).Methods("GET")
	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "pages/login.html")
	}).Methods("GET")

	// Роут для профиля
	r.HandleFunc("/profile", middleware.AuthMiddleware(userService)(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "pages/profile.html")
	})).Methods("GET")

	// Роут для админ панели
	r.HandleFunc("/admin", middleware.AdminAuthMiddleware(userService)(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "pages/admin_panel.html")
	})).Methods("GET")

	// Регистрируем маршруты для пользователя
	r.HandleFunc("/register", handlers.UserRegistrationHandler(userService, emailService)).Methods("POST")
	r.HandleFunc("/verify", handlers.UserVerificationHandler(userService)).Methods("GET")
	r.HandleFunc("/login", handlers.UserLoginHandler(userService)).Methods("POST")
	r.HandleFunc("/logout", handlers.UserLogoutHandler(userService)).Methods("POST")

	// Регистрируем маршруты для повторной отправки письма с подтверждением
	r.HandleFunc("/resend-verification-email", handlers.ResendVerificationEmailHandler(emailService)).Methods("POST")

	// Инициализируем обработчик администратора
	adminHandler := handlers.NewAdminHandler(adminService)

	// Регистрируем маршруты для пользователей
	r.HandleFunc("/api/admin/users", adminHandler.GetAllUsersHandler).Methods("GET")          // Получить всех пользователей
	r.HandleFunc("/api/admin/users/{id}/block", adminHandler.BlockUserHandler).Methods("PUT") // Заблокировать пользователя
	r.HandleFunc("/api/admin/users/{id}", adminHandler.UpdateUserRoleHandler).Methods("PUT")

	return r
}
