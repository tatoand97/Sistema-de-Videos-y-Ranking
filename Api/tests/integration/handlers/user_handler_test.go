package handlers_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"api/internal/application/useCase"
	"api/internal/domain"
	"api/internal/domain/entities"
	hdl "api/internal/presentation/handlers"

	"github.com/gin-gonic/gin"
)

type fakeUserRepoForReg struct {
	emailExists bool
	createErr   error
}

func (f *fakeUserRepoForReg) Create(ctx context.Context, user *entities.User) error {
	return f.createErr
}

func (f *fakeUserRepoForReg) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	if f.emailExists {
		return &entities.User{Email: email}, nil
	}
	return nil, domain.ErrNotFound
}

func (f *fakeUserRepoForReg) EmailExists(ctx context.Context, email string) (bool, error) {
	return f.emailExists, nil
}
func (f *fakeUserRepoForReg) GetPermissions(ctx context.Context, userID uint) ([]string, error) {
	return []string{}, nil
}

type fakeLocationRepoForReg struct {
	cityID int
	err    error
}

func (f *fakeLocationRepoForReg) GetCityID(ctx context.Context, country, city string) (int, error) {
	if f.err != nil {
		return 0, f.err
	}
	return f.cityID, nil
}

func setupUserRouter(userRepo *fakeUserRepoForReg, locationRepo *fakeLocationRepoForReg) *gin.Engine {
	gin.SetMode(gin.TestMode)
	svc := useCase.NewUserService(userRepo, locationRepo)
	h := hdl.NewUserHandlers(svc)
	r := gin.New()
	r.POST("/api/auth/signup", h.Register)
	return r
}

func TestRegister_Success(t *testing.T) {
	userRepo := &fakeUserRepoForReg{emailExists: false}
	locationRepo := &fakeLocationRepoForReg{cityID: 1}
	r := setupUserRouter(userRepo, locationRepo)

	body := `{
		"first_name": "John",
		"last_name": "Doe",
		"email": "john@example.com",
		"password1": "password123",
		"password2": "password123",
		"country": "Colombia",
		"city": "Bogotá"
	}`

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/auth/signup", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestRegister_InvalidJSON(t *testing.T) {
	userRepo := &fakeUserRepoForReg{}
	locationRepo := &fakeLocationRepoForReg{}
	r := setupUserRouter(userRepo, locationRepo)

	body := `{"first_name": "John"`

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/auth/signup", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestRegister_EmailExists(t *testing.T) {
	userRepo := &fakeUserRepoForReg{emailExists: true}
	locationRepo := &fakeLocationRepoForReg{}
	r := setupUserRouter(userRepo, locationRepo)

	body := `{
		"first_name": "John",
		"last_name": "Doe",
		"email": "existing@example.com",
		"password1": "password123",
		"password2": "password123",
		"country": "Colombia",
		"city": "Bogotá"
	}`

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/auth/signup", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestRegister_PasswordMismatch(t *testing.T) {
	userRepo := &fakeUserRepoForReg{emailExists: false}
	locationRepo := &fakeLocationRepoForReg{}
	r := setupUserRouter(userRepo, locationRepo)

	body := `{
		"first_name": "John",
		"last_name": "Doe",
		"email": "john@example.com",
		"password1": "password123",
		"password2": "different",
		"country": "Colombia",
		"city": "Bogotá"
	}`

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/auth/signup", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestRegister_InvalidCity(t *testing.T) {
	userRepo := &fakeUserRepoForReg{emailExists: false}
	locationRepo := &fakeLocationRepoForReg{err: domain.ErrNotFound}
	r := setupUserRouter(userRepo, locationRepo)

	body := `{
		"first_name": "John",
		"last_name": "Doe",
		"email": "john@example.com",
		"password1": "password123",
		"password2": "password123",
		"country": "Colombia",
		"city": "InvalidCity"
	}`

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/auth/signup", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestRegister_InvalidCountry(t *testing.T) {
	userRepo := &fakeUserRepoForReg{emailExists: false}
	locationRepo := &fakeLocationRepoForReg{err: domain.ErrInvalid}
	r := setupUserRouter(userRepo, locationRepo)

	body := `{
		"first_name": "John",
		"last_name": "Doe",
		"email": "john@example.com",
		"password1": "password123",
		"password2": "password123",
		"country": "InvalidCountry",
		"city": "Bogotá"
	}`

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/auth/signup", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestRegister_CreateUserError(t *testing.T) {
	userRepo := &fakeUserRepoForReg{emailExists: false, createErr: errors.New("database error")}
	locationRepo := &fakeLocationRepoForReg{cityID: 1}
	r := setupUserRouter(userRepo, locationRepo)

	body := `{
		"first_name": "John",
		"last_name": "Doe",
		"email": "john@example.com",
		"password1": "password123",
		"password2": "password123",
		"country": "Colombia",
		"city": "Bogotá"
	}`

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/auth/signup", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d, body=%s", w.Code, w.Body.String())
	}
}
