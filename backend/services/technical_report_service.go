package services

import (
	"github.com/root/evodrones/backend/models"
	"github.com/root/evodrones/backend/repositories"
)

type TechnicalReportService struct {
	reportRepo      *repositories.TechnicalReportRepository
	maintenanceRepo *repositories.MaintenanceOrderRepository
	userRepo        *repositories.UserRepository
}

func NewTechnicalReportService(
	reportRepo *repositories.TechnicalReportRepository,
	maintenanceRepo *repositories.MaintenanceOrderRepository,
	userRepo *repositories.UserRepository,
) *TechnicalReportService {
	return &TechnicalReportService{
		reportRepo:      reportRepo,
		maintenanceRepo: maintenanceRepo,
		userRepo:        userRepo,
	}
}

func (s *TechnicalReportService) Create(input *models.CreateTechnicalReportInput, userID string) (*models.TechnicalReport, error) {
	// Verificar se a ordem de manutenção existe
	maintenance, err := s.maintenanceRepo.GetByID(input.MaintenanceID)
	if err != nil {
		return nil, err
	}
	if maintenance == nil {
		return nil, ErrMaintenanceOrderNotFound
	}

	// Verificar se o usuário existe
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	report := &models.TechnicalReport{
		MaintenanceID:   input.MaintenanceID,
		UserID:          userID,
		Title:           input.Title,
		Description:     input.Description,
		Findings:        input.Findings,
		Recommendations: input.Recommendations,
	}

	err = s.reportRepo.Create(report)
	if err != nil {
		return nil, err
	}

	return report, nil
}

func (s *TechnicalReportService) GetByID(id string, userID string, userRole string) (*models.TechnicalReportWithDetails, error) {
	report, err := s.reportRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if report == nil {
		return nil, ErrTechnicalReportNotFound
	}

	// Verificar permissão
	if userRole != "admin" && report.UserID != userID {
		return nil, ErrUnauthorized
	}

	return report, nil
}

func (s *TechnicalReportService) GetByMaintenanceID(maintenanceID string, userID string, userRole string) ([]*models.TechnicalReportWithDetails, error) {
	// Verificar se a ordem de manutenção existe
	maintenance, err := s.maintenanceRepo.GetByID(maintenanceID)
	if err != nil {
		return nil, err
	}
	if maintenance == nil {
		return nil, ErrMaintenanceOrderNotFound
	}

	// Verificar permissão
	if userRole != "admin" && maintenance.UserID != userID {
		return nil, ErrUnauthorized
	}

	reports, err := s.reportRepo.GetByMaintenanceID(maintenanceID)
	if err != nil {
		return nil, err
	}

	return reports, nil
}

func (s *TechnicalReportService) GetByUserID(userID string, requestingUserID string, userRole string) ([]*models.TechnicalReportWithDetails, error) {
	// Verificar permissão
	if userRole != "admin" && userID != requestingUserID {
		return nil, ErrUnauthorized
	}

	reports, err := s.reportRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	return reports, nil
}

func (s *TechnicalReportService) Update(id string, input *models.UpdateTechnicalReportInput, userID string, userRole string) error {
	// Verificar se o relatório existe
	report, err := s.reportRepo.GetByID(id)
	if err != nil {
		return err
	}
	if report == nil {
		return ErrTechnicalReportNotFound
	}

	// Verificar permissão
	if userRole != "admin" && report.UserID != userID {
		return ErrUnauthorized
	}

	return s.reportRepo.Update(id, input)
}

func (s *TechnicalReportService) Delete(id string, userID string, userRole string) error {
	// Verificar se o relatório existe
	report, err := s.reportRepo.GetByID(id)
	if err != nil {
		return err
	}
	if report == nil {
		return ErrTechnicalReportNotFound
	}

	// Verificar permissão
	if userRole != "admin" && report.UserID != userID {
		return ErrUnauthorized
	}

	return s.reportRepo.Delete(id)
}

func (s *TechnicalReportService) List(offset, limit int, userID string, userRole string) ([]*models.TechnicalReportWithDetails, error) {
	// Se não for admin, retornar apenas os relatórios do usuário
	if userRole != "admin" {
		return s.reportRepo.GetByUserID(userID)
	}

	return s.reportRepo.List(offset, limit)
}
