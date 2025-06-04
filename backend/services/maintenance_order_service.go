package services

import (
	"github.com/root/evodrones/backend/models"
	"github.com/root/evodrones/backend/repositories"
)

type MaintenanceOrderService struct {
	maintenanceRepo *repositories.MaintenanceOrderRepository
	droneRepo       *repositories.DroneRepository
	userRepo        *repositories.UserRepository
}

func NewMaintenanceOrderService(
	maintenanceRepo *repositories.MaintenanceOrderRepository,
	droneRepo *repositories.DroneRepository,
	userRepo *repositories.UserRepository,
) *MaintenanceOrderService {
	return &MaintenanceOrderService{
		maintenanceRepo: maintenanceRepo,
		droneRepo:       droneRepo,
		userRepo:        userRepo,
	}
}

func (s *MaintenanceOrderService) Create(userID string, input models.CreateMaintenanceOrderInput) (*models.MaintenanceOrderWithDetails, error) {
	// Verificar se o usuário existe
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUnauthorized
	}

	// Verificar se o drone existe
	drone, err := s.droneRepo.GetByID(input.DroneID)
	if err != nil {
		return nil, err
	}
	if drone == nil {
		return nil, ErrDroneNotFound
	}

	// Criar ordem de manutenção
	order := &models.MaintenanceOrder{
		DroneID:     input.DroneID,
		UserID:      userID,
		Description: input.Description,
		StartDate:   input.StartDate,
	}

	if err := s.maintenanceRepo.Create(order); err != nil {
		return nil, err
	}

	// Retornar ordem com detalhes
	return s.maintenanceRepo.GetByID(order.ID)
}

func (s *MaintenanceOrderService) GetByID(userID, orderID string) (*models.MaintenanceOrderWithDetails, error) {
	// Verificar se o usuário existe
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUnauthorized
	}

	// Buscar ordem de manutenção
	order, err := s.maintenanceRepo.GetByID(orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrMaintenanceOrderNotFound
	}

	// Verificar se o usuário tem permissão
	if order.UserID != userID && user.Role != models.UserRoleAdmin {
		return nil, ErrUnauthorized
	}

	return order, nil
}

func (s *MaintenanceOrderService) GetByUserID(userID string) ([]*models.MaintenanceOrderWithDetails, error) {
	// Verificar se o usuário existe
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUnauthorized
	}

	return s.maintenanceRepo.GetByUserID(userID)
}

func (s *MaintenanceOrderService) GetByDroneID(userID, droneID string) ([]*models.MaintenanceOrderWithDetails, error) {
	// Verificar se o usuário existe
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUnauthorized
	}

	// Verificar se o drone existe
	drone, err := s.droneRepo.GetByID(droneID)
	if err != nil {
		return nil, err
	}
	if drone == nil {
		return nil, ErrDroneNotFound
	}

	// Verificar se o usuário tem permissão
	if drone.UserID != userID && user.Role != models.UserRoleAdmin {
		return nil, ErrUnauthorized
	}

	return s.maintenanceRepo.GetByDroneID(droneID)
}

func (s *MaintenanceOrderService) Update(userID, orderID string, input models.UpdateMaintenanceOrderInput) (*models.MaintenanceOrderWithDetails, error) {
	// Verificar se o usuário existe
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUnauthorized
	}

	// Buscar ordem de manutenção
	order, err := s.maintenanceRepo.GetByID(orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrMaintenanceOrderNotFound
	}

	// Verificar se o usuário tem permissão
	if order.UserID != userID && user.Role != models.UserRoleAdmin {
		return nil, ErrUnauthorized
	}

	// Atualizar ordem de manutenção
	if err := s.maintenanceRepo.Update(orderID, &input); err != nil {
		return nil, err
	}

	// Retornar ordem atualizada
	return s.maintenanceRepo.GetByID(orderID)
}

func (s *MaintenanceOrderService) Delete(userID, orderID string) error {
	// Verificar se o usuário existe
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUnauthorized
	}

	// Buscar ordem de manutenção
	order, err := s.maintenanceRepo.GetByID(orderID)
	if err != nil {
		return err
	}
	if order == nil {
		return ErrMaintenanceOrderNotFound
	}

	// Verificar se o usuário tem permissão
	if order.UserID != userID && user.Role != models.UserRoleAdmin {
		return ErrUnauthorized
	}

	return s.maintenanceRepo.Delete(orderID)
}

func (s *MaintenanceOrderService) List(offset, limit int) ([]*models.MaintenanceOrderWithDetails, error) {
	return s.maintenanceRepo.List(offset, limit)
}

func (s *MaintenanceOrderService) GetByStatus(status models.MaintenanceOrderStatus) ([]*models.MaintenanceOrderWithDetails, error) {
	return s.maintenanceRepo.GetByStatus(status)
}
