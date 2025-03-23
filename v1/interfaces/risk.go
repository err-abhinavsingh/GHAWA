package interfaces

import "GHAWA/entities"

type RiskService interface {
	GetRisk(string) *entities.Risk
	GetRisks() []*entities.Risk
	CreateRisk(risk entities.Risk) (*entities.Risk, error)
}
