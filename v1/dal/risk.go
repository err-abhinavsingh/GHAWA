package dal

import (
	"GHAWA/entities"
	"sync"
)

type RiskDal struct {

	// Risks indexed at Uuid
	risks       map[string]*entities.Risk
	createMutex sync.Mutex
}

func NewRiskDal() *RiskDal {
	return &RiskDal{
		risks: make(map[string]*entities.Risk),
	}
}

func (d *RiskDal) GetRisks() []*entities.Risk {

	var risks []*entities.Risk
	for _, risk := range d.risks {
		risks = append(risks, risk)
	}
	return risks
}

func (d *RiskDal) GetRisk(riskId string) *entities.Risk {

	risk := d.risks[riskId]
	return risk
}

func (d *RiskDal) CreateRisk(risk entities.Risk) {

	d.createMutex.Lock()
	defer d.createMutex.Unlock()

	d.risks[risk.Uuid] = &risk
}
