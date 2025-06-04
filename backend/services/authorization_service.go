package services

import (
	"github.com/root/evodrones/backend/models"
	"github.com/root/evodrones/backend/repositories"
)

// Tipos de recursos do sistema
const (
	ResourceTypeDrone            = "drone"
	ResourceTypeMaintenanceOrder = "maintenance_order"
	ResourceTypeTechnicalReport  = "technical_report"
	ResourceTypeUser             = "user"
)

// Ações possíveis sobre recursos
const (
	ActionCreate = "create"
	ActionRead   = "read"
	ActionUpdate = "update"
	ActionDelete = "delete"
	ActionList   = "list"
)

// AuthorizationService centraliza a lógica de autorização do sistema
type AuthorizationService struct {
	userRepo        *repositories.UserRepository
	droneRepo       *repositories.DroneRepository
	maintenanceRepo *repositories.MaintenanceOrderRepository
	reportRepo      *repositories.TechnicalReportRepository
}

// NewAuthorizationService cria uma nova instância do serviço de autorização
func NewAuthorizationService(
	userRepo *repositories.UserRepository,
	droneRepo *repositories.DroneRepository,
	maintenanceRepo *repositories.MaintenanceOrderRepository,
	reportRepo *repositories.TechnicalReportRepository,
) *AuthorizationService {
	return &AuthorizationService{
		userRepo:        userRepo,
		droneRepo:       droneRepo,
		maintenanceRepo: maintenanceRepo,
		reportRepo:      reportRepo,
	}
}

// CanAccessResource verifica se um usuário pode acessar um recurso específico
func (s *AuthorizationService) CanAccessResource(resourceType, resourceID, userID, userRole, action string) (bool, error) {
	// Administradores têm acesso total a todos os recursos
	if userRole == models.UserRoleAdmin {
		return true, nil
	}

	// Técnicos têm acesso a todas as ordens de manutenção e relatórios técnicos
	if userRole == models.UserRoleTechnician {
		if resourceType == ResourceTypeMaintenanceOrder || resourceType == ResourceTypeTechnicalReport {
			return true, nil
		}
	}

	// Verificações específicas por tipo de recurso
	switch resourceType {
	case ResourceTypeDrone:
		return s.canAccessDrone(resourceID, userID, action)
	case ResourceTypeMaintenanceOrder:
		return s.canAccessMaintenanceOrder(resourceID, userID, action)
	case ResourceTypeTechnicalReport:
		return s.canAccessTechnicalReport(resourceID, userID, action)
	case ResourceTypeUser:
		return s.canAccessUser(resourceID, userID, action)
	default:
		return false, ErrInvalidResourceType
	}
}

// canAccessDrone verifica permissões para drones
func (s *AuthorizationService) canAccessDrone(droneID, userID, action string) (bool, error) {
	// Se for criação, qualquer usuário autenticado pode criar
	if action == ActionCreate {
		return true, nil
	}

	// Para outras ações, verificar propriedade
	drone, err := s.droneRepo.GetByID(droneID)
	if err != nil {
		return false, err
	}
	if drone == nil {
		return false, ErrDroneNotFound
	}

	// Verificar se o usuário é o proprietário do drone
	return drone.UserID == userID, nil
}

// canAccessMaintenanceOrder verifica permissões para ordens de manutenção
func (s *AuthorizationService) canAccessMaintenanceOrder(orderID, userID, action string) (bool, error) {
	// Se for criação, qualquer usuário autenticado pode criar
	if action == ActionCreate {
		return true, nil
	}

	// Para outras ações, verificar propriedade
	order, err := s.maintenanceRepo.GetByID(orderID)
	if err != nil {
		return false, err
	}
	if order == nil {
		return false, ErrMaintenanceOrderNotFound
	}

	// Verificar se o usuário é o proprietário da ordem
	return order.UserID == userID, nil
}

// canAccessTechnicalReport verifica permissões para relatórios técnicos
func (s *AuthorizationService) canAccessTechnicalReport(reportID, userID, action string) (bool, error) {
	// Se for criação, qualquer usuário autenticado pode criar
	if action == ActionCreate {
		return true, nil
	}

	// Para outras ações, verificar propriedade
	report, err := s.reportRepo.GetByID(reportID)
	if err != nil {
		return false, err
	}
	if report == nil {
		return false, ErrTechnicalReportNotFound
	}

	// Verificar se o usuário é o proprietário do relatório
	return report.UserID == userID, nil
}

// canAccessUser verifica permissões para usuários
func (s *AuthorizationService) canAccessUser(targetUserID, requestingUserID, action string) (bool, error) {
	// Usuários só podem acessar seus próprios dados
	return targetUserID == requestingUserID, nil
}

// CanListResources verifica se um usuário pode listar recursos de um determinado tipo
func (s *AuthorizationService) CanListResources(resourceType, userID, userRole string) (bool, error) {
	// Administradores podem listar todos os recursos
	if userRole == models.UserRoleAdmin {
		return true, nil
	}

	// Técnicos podem listar todas as ordens de manutenção e relatórios técnicos
	if userRole == models.UserRoleTechnician {
		if resourceType == ResourceTypeMaintenanceOrder || resourceType == ResourceTypeTechnicalReport {
			return true, nil
		}
	}

	// Usuários comuns só podem listar seus próprios recursos
	return true, nil
}

// GetAccessibleResources retorna os IDs dos recursos que um usuário pode acessar
func (s *AuthorizationService) GetAccessibleResources(resourceType, userID, userRole, action string) ([]string, error) {
	// Administradores podem acessar todos os recursos
	if userRole == models.UserRoleAdmin {
		return nil, nil // nil significa todos os recursos
	}

	// Técnicos podem acessar todas as ordens de manutenção e relatórios técnicos
	if userRole == models.UserRoleTechnician {
		if resourceType == ResourceTypeMaintenanceOrder || resourceType == ResourceTypeTechnicalReport {
			return nil, nil // nil significa todos os recursos
		}
	}

	// Para usuários comuns, retornar apenas os recursos que eles possuem
	switch resourceType {
	case ResourceTypeDrone:
		drones, err := s.droneRepo.GetByUserID(userID)
		if err != nil {
			return nil, err
		}
		ids := make([]string, len(drones))
		for i, drone := range drones {
			ids[i] = drone.ID
		}
		return ids, nil

	case ResourceTypeMaintenanceOrder:
		orders, err := s.maintenanceRepo.GetByUserID(userID)
		if err != nil {
			return nil, err
		}
		ids := make([]string, len(orders))
		for i, order := range orders {
			ids[i] = order.ID
		}
		return ids, nil

	case ResourceTypeTechnicalReport:
		reports, err := s.reportRepo.GetByUserID(userID)
		if err != nil {
			return nil, err
		}
		ids := make([]string, len(reports))
		for i, report := range reports {
			ids[i] = report.ID
		}
		return ids, nil

	default:
		return nil, ErrInvalidResourceType
	}
}
