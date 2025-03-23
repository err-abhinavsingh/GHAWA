package services

import (
	"GHAWA/dal"
	"GHAWA/entities"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type RiskService struct {
	riskDal *dal.RiskDal
}

func NewRiskService() *RiskService {
	return &RiskService{
		riskDal: dal.NewRiskDal(),
	}
}

func (s *RiskService) GetRisks() []*entities.Risk {
	return s.riskDal.GetRisks()
}

func (s *RiskService) GetRisk(riskId string) *entities.Risk {
	return s.riskDal.GetRisk(riskId)
}

func (s *RiskService) CreateRisk(risk entities.Risk) (*entities.Risk, error) {

	if !entities.ValidRisk(risk.State) {
		return nil, fmt.Errorf("invalid risk state: %s", risk.State)
	}
	id := uuid.New().String()
	risk.Uuid = id
	s.riskDal.CreateRisk(risk)
	log.Print("INFO: ", "Created risk with uuid ", risk.Uuid)
	return &risk, nil
}
