package controllers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/root/evodrones/backend/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockTechnicalReportService struct {
	mock.Mock
}

func (m *MockTechnicalReportService) Create(input *models.CreateTechnicalReportInput, userID uuid.UUID) (*models.TechnicalReport, error) {
	args := m.Called(input, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TechnicalReport), args.Error(1)
}

func (m *MockTechnicalReportService) GetByID(id uuid.UUID, userID uuid.UUID, userRole string) (*models.TechnicalReportWithDetails, error) {
	args := m.Called(id, userID, userRole)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TechnicalReportWithDetails), args.Error(1)
}

func (m *MockTechnicalReportService) GetByMaintenanceID(maintenanceID uuid.UUID, userID uuid.UUID, userRole string) ([]*models.TechnicalReportWithDetails, error) {
	args := m.Called(maintenanceID, userID, userRole)
	return args.Get(0).([]*models.TechnicalReportWithDetails), args.Error(1)
}

func (m *MockTechnicalReportService) GetByUserID(userID uuid.UUID, requestingUserID uuid.UUID, userRole string) ([]*models.TechnicalReportWithDetails, error) {
	args := m.Called(userID, requestingUserID, userRole)
	return args.Get(0).([]*models.TechnicalReportWithDetails), args.Error(1)
}

func (m *MockTechnicalReportService) Update(id uuid.UUID, input *models.UpdateTechnicalReportInput, userID uuid.UUID, userRole string) error {
	args := m.Called(id, input, userID, userRole)
	return args.Error(0)
}

func (m *MockTechnicalReportService) Delete(id uuid.UUID, userID uuid.UUID, userRole string) error {
	args := m.Called(id, userID, userRole)
	return args.Error(0)
}

func (m *MockTechnicalReportService) List(offset, limit int, userID uuid.UUID, userRole string) ([]*models.TechnicalReportWithDetails, error) {
	args := m.Called(offset, limit, userID, userRole)
	return args.Get(0).([]*models.TechnicalReportWithDetails), args.Error(1)
}

func TestTechnicalReportController_Create(t *testing.T) {
	mockService := new(MockTechnicalReportService)
	controller := NewTechnicalReportController(mockService)

	app := fiber.New()
	app.Post("/reports", controller.Create)

	maintenanceID := uuid.New()
	userID := uuid.New()

	// Configurar mock
	mockService.On("Create", mock.AnythingOfType("*models.CreateTechnicalReportInput"), userID).Return(&models.TechnicalReport{
		ID:              uuid.New(),
		MaintenanceID:   maintenanceID,
		UserID:          userID,
		Title:           "Test Report",
		Description:     "Test Description",
		Findings:        "Test Findings",
		Recommendations: "Test Recommendations",
	}, nil)

	// Criar requisição
	body := map[string]interface{}{
		"maintenance_id":  maintenanceID,
		"title":           "Test Report",
		"description":     "Test Description",
		"findings":        "Test Findings",
		"recommendations": "Test Recommendations",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/reports", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")

	// Executar requisição
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	// Verificar resposta
	var response models.TechnicalReport
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, maintenanceID, response.MaintenanceID)
	assert.Equal(t, userID, response.UserID)
	assert.Equal(t, "Test Report", response.Title)

	mockService.AssertExpectations(t)
}

func TestTechnicalReportController_GetByID(t *testing.T) {
	mockService := new(MockTechnicalReportService)
	controller := NewTechnicalReportController(mockService)

	app := fiber.New()
	app.Get("/reports/:id", controller.GetByID)

	reportID := uuid.New()
	userID := uuid.New()

	// Configurar mock
	mockService.On("GetByID", reportID, userID, "user").Return(&models.TechnicalReportWithDetails{
		TechnicalReport: models.TechnicalReport{
			ID:     reportID,
			UserID: userID,
		},
	}, nil)

	// Criar requisição
	req := httptest.NewRequest("GET", "/reports/"+reportID.String(), nil)
	req.Header.Set("Authorization", "Bearer test-token")

	// Executar requisição
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Verificar resposta
	var response models.TechnicalReportWithDetails
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, reportID, response.ID)
	assert.Equal(t, userID, response.UserID)

	mockService.AssertExpectations(t)
}

func TestTechnicalReportController_GetByMaintenanceID(t *testing.T) {
	mockService := new(MockTechnicalReportService)
	controller := NewTechnicalReportController(mockService)

	app := fiber.New()
	app.Get("/reports/maintenance/:maintenanceID", controller.GetByMaintenanceID)

	maintenanceID := uuid.New()
	userID := uuid.New()

	// Configurar mock
	mockService.On("GetByMaintenanceID", maintenanceID, userID, "user").Return([]*models.TechnicalReportWithDetails{
		{
			TechnicalReport: models.TechnicalReport{
				ID:            uuid.New(),
				MaintenanceID: maintenanceID,
				UserID:        userID,
			},
		},
	}, nil)

	// Criar requisição
	req := httptest.NewRequest("GET", "/reports/maintenance/"+maintenanceID.String(), nil)
	req.Header.Set("Authorization", "Bearer test-token")

	// Executar requisição
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Verificar resposta
	var response []*models.TechnicalReportWithDetails
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, maintenanceID, response[0].MaintenanceID)
	assert.Equal(t, userID, response[0].UserID)

	mockService.AssertExpectations(t)
}

func TestTechnicalReportController_GetByUserID(t *testing.T) {
	mockService := new(MockTechnicalReportService)
	controller := NewTechnicalReportController(mockService)

	app := fiber.New()
	app.Get("/reports/user/:userID", controller.GetByUserID)

	userID := uuid.New()

	// Configurar mock
	mockService.On("GetByUserID", userID, userID, "user").Return([]*models.TechnicalReportWithDetails{
		{
			TechnicalReport: models.TechnicalReport{
				ID:     uuid.New(),
				UserID: userID,
			},
		},
	}, nil)

	// Criar requisição
	req := httptest.NewRequest("GET", "/reports/user/"+userID.String(), nil)
	req.Header.Set("Authorization", "Bearer test-token")

	// Executar requisição
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Verificar resposta
	var response []*models.TechnicalReportWithDetails
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, userID, response[0].UserID)

	mockService.AssertExpectations(t)
}

func TestTechnicalReportController_Update(t *testing.T) {
	mockService := new(MockTechnicalReportService)
	controller := NewTechnicalReportController(mockService)

	app := fiber.New()
	app.Put("/reports/:id", controller.Update)

	reportID := uuid.New()
	userID := uuid.New()

	// Configurar mock
	mockService.On("Update", reportID, mock.AnythingOfType("*models.UpdateTechnicalReportInput"), userID, "user").Return(nil)

	// Criar requisição
	title := "Updated Title"
	description := "Updated Description"
	findings := "Updated Findings"
	recommendations := "Updated Recommendations"

	body := map[string]interface{}{
		"title":           title,
		"description":     description,
		"findings":        findings,
		"recommendations": recommendations,
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("PUT", "/reports/"+reportID.String(), bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")

	// Executar requisição
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestTechnicalReportController_Delete(t *testing.T) {
	mockService := new(MockTechnicalReportService)
	controller := NewTechnicalReportController(mockService)

	app := fiber.New()
	app.Delete("/reports/:id", controller.Delete)

	reportID := uuid.New()
	userID := uuid.New()

	// Configurar mock
	mockService.On("Delete", reportID, userID, "user").Return(nil)

	// Criar requisição
	req := httptest.NewRequest("DELETE", "/reports/"+reportID.String(), nil)
	req.Header.Set("Authorization", "Bearer test-token")

	// Executar requisição
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestTechnicalReportController_List(t *testing.T) {
	mockService := new(MockTechnicalReportService)
	controller := NewTechnicalReportController(mockService)

	app := fiber.New()
	app.Get("/reports", controller.List)

	userID := uuid.New()

	// Configurar mock
	mockService.On("List", 0, 10, userID, "user").Return([]*models.TechnicalReportWithDetails{
		{
			TechnicalReport: models.TechnicalReport{
				ID:     uuid.New(),
				UserID: userID,
			},
		},
	}, nil)

	// Criar requisição
	req := httptest.NewRequest("GET", "/reports", nil)
	req.Header.Set("Authorization", "Bearer test-token")

	// Executar requisição
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Verificar resposta
	var response []*models.TechnicalReportWithDetails
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, userID, response[0].UserID)

	mockService.AssertExpectations(t)
}
