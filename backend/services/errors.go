package services

import "errors"

var (
	ErrUnauthorized             = errors.New("não autorizado")
	ErrUserNotFound             = errors.New("usuário não encontrado")
	ErrDroneNotFound            = errors.New("drone não encontrado")
	ErrMaintenanceOrderNotFound = errors.New("ordem de manutenção não encontrada")
	ErrTechnicalReportNotFound  = errors.New("relatório técnico não encontrado")
	ErrInvalidResourceType      = errors.New("tipo de recurso inválido")
	ErrResourceNotFound         = errors.New("recurso não encontrado")
)
