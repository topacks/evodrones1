package services

import (
	"github.com/root/evodrones/backend/models"
	"github.com/root/evodrones/backend/repositories"
)

type DroneService struct {
	droneRepo *repositories.DroneRepository
	userRepo  *repositories.UserRepository
}

func NewDroneService(droneRepo *repositories.DroneRepository, userRepo *repositories.UserRepository) *DroneService {
	return &DroneService{
		droneRepo: droneRepo,
		userRepo:  userRepo,
	}
}

func (s *DroneService) Create(userID string, input models.CreateDroneInput) (*models.Drone, error) {
	// Verificar se o usuário existe
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// Criar novo drone
	drone := models.NewDrone(userID, input)

	// Salvar no banco de dados
	if err := s.droneRepo.Create(drone); err != nil {
		return nil, err
	}

	return drone, nil
}

func (s *DroneService) GetByID(userID string, droneID string) (*models.Drone, error) {
	drone, err := s.droneRepo.GetByID(droneID)
	if err != nil {
		return nil, err
	}
	if drone == nil {
		return nil, ErrDroneNotFound
	}

	// Verificar se o drone pertence ao usuário
	if drone.UserID != userID {
		return nil, ErrUserNotFound
	}

	return drone, nil
}

func (s *DroneService) GetByUserID(userID string) ([]*models.Drone, error) {
	return s.droneRepo.GetByUserID(userID)
}

func (s *DroneService) Update(userID string, droneID string, input models.UpdateDroneInput) (*models.Drone, error) {
	// Verificar se o drone existe e pertence ao usuário
	drone, err := s.GetByID(userID, droneID)
	if err != nil {
		return nil, err
	}

	// Atualizar campos
	drone.Model = input.Model
	drone.Serial = input.Serial
	drone.Status = input.Status

	// Salvar no banco de dados
	if err := s.droneRepo.Update(drone); err != nil {
		return nil, err
	}

	return drone, nil
}

func (s *DroneService) Delete(userID string, droneID string) error {
	// Verificar se o drone existe e pertence ao usuário
	_, err := s.GetByID(userID, droneID)
	if err != nil {
		return err
	}

	return s.droneRepo.Delete(droneID)
}

func (s *DroneService) List(offset, limit int) ([]*models.Drone, error) {
	return s.droneRepo.List(offset, limit)
}

func (s *DroneService) GetByStatus(status string) ([]*models.Drone, error) {
	return s.droneRepo.GetByStatus(status)
}

func (s *DroneService) GetDueForCheck() ([]*models.Drone, error) {
	return s.droneRepo.GetDueForCheck()
}
