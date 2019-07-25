package actions

import (
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	assert "github.com/LindsayBradford/crem/pkg/assert/debug"
	"math"
)

const (
	hillSlopeAreaIndex                   = 11
	proportionOfHillSlopeVegetationIndex = 12
	hillSlopeRKLSIndex                   = 13
	hillSlopeDistanceIndex               = 14
)

type hillslopeSedimentTracker struct {
	rslk                         float64
	originalProportionVegetation float64
	distanceToCatchment          float64
	area                         float64
}

type HillSlopeSedimentContribution struct {
	planningUnitTable tables.CsvTable
	parameters        parameters.Parameters

	contributionMap map[planningunit.Id]hillslopeSedimentTracker
}

func (h *HillSlopeSedimentContribution) Initialise(planningUnitTable tables.CsvTable, parameters parameters.Parameters) {
	h.planningUnitTable = planningUnitTable
	h.parameters = parameters
	h.populateContributionMap()
}

func (h *HillSlopeSedimentContribution) populateContributionMap() {
	_, rowCount := h.planningUnitTable.ColumnAndRowSize()
	h.contributionMap = make(map[planningunit.Id]hillslopeSedimentTracker, rowCount)

	for row := uint(0); row < rowCount; row++ {
		h.populateContributionMapEntry(row)
	}
}

func (h *HillSlopeSedimentContribution) populateContributionMapEntry(rowNumber uint) {
	planningUnit := h.planningUnitTable.CellFloat64(planningUnitIndex, rowNumber)
	mapKey := planningunit.Float64ToId(planningUnit)

	h.contributionMap[mapKey] = hillslopeSedimentTracker{
		rslk:                         h.hillslopeRkls(rowNumber),
		originalProportionVegetation: h.originalHillSlopeVegetation(rowNumber),
		distanceToCatchment:          h.distanceToCatchment(rowNumber),
		area:                         h.hillslopArea(rowNumber),
	}
}

func (h *HillSlopeSedimentContribution) hillslopeRkls(rowNumber uint) float64 {
	// rkls: Rainfall erosivity factor  (R) * Soil Erodibility Factor (K) * Slope length (L) * Slope Steepness (S)
	// See Catchment Rehabilitation Planner final report, section 3.2.3
	rkls := h.planningUnitTable.CellFloat64(hillSlopeRKLSIndex, rowNumber)
	return rkls
}

func (h *HillSlopeSedimentContribution) hillslopArea(rowNumber uint) float64 {
	return h.planningUnitTable.CellFloat64(hillSlopeAreaIndex, rowNumber)
}

func (h *HillSlopeSedimentContribution) originalHillSlopeVegetation(rowNumber uint) float64 {
	return h.planningUnitTable.CellFloat64(proportionOfHillSlopeVegetationIndex, rowNumber)
}

func (h *HillSlopeSedimentContribution) distanceToCatchment(rowNumber uint) float64 {
	return h.planningUnitTable.CellFloat64(hillSlopeDistanceIndex, rowNumber)
}

func (h *HillSlopeSedimentContribution) OriginalPlanningUnitVegetationProportion(id planningunit.Id) float64 {
	planningUnitSedimentTracker, planningUnitIsPresent := h.contributionMap[id]
	assert.That(planningUnitIsPresent).Holds()

	return planningUnitSedimentTracker.originalProportionVegetation
}

func (h *HillSlopeSedimentContribution) OriginalPlanningUnitSedimentContribution(id planningunit.Id) float64 {
	planningUnitSedimentTracker, planningUnitIsPresent := h.contributionMap[id]
	assert.That(planningUnitIsPresent).Holds()

	originalVegetationCover := h.calculateVegetationCover(id, planningUnitSedimentTracker.originalProportionVegetation)
	sedimentContribution := planningUnitSedimentTracker.rslk * originalVegetationCover
	return sedimentContribution
}

func (h *HillSlopeSedimentContribution) PlanningUnitSedimentContribution(planningUnit planningunit.Id, proportionOfHillSlopeVegetation float64) float64 {
	planningUnitSedimentTracker, planningUnitIsPresent := h.contributionMap[planningUnit]
	assert.That(planningUnitIsPresent).Holds()

	sedimentContribution := planningUnitSedimentTracker.rslk * h.calculateVegetationCover(planningUnit, proportionOfHillSlopeVegetation)

	return sedimentContribution
}

func (h *HillSlopeSedimentContribution) calculateVegetationCover(planningUnit planningunit.Id, proportionOfHillSlopeVegetation float64) float64 {
	// See: CRP final report, section 3.2.3, pg 24.
	distanceToRiparianBuffer := h.contributionMap[planningUnit].distanceToCatchment
	area := h.contributionMap[planningUnit].area
	groundCover := area * proportionOfHillSlopeVegetation

	unmodifiedVegetationCoverFactor := 0.5665 * math.Exp(-0.0487*groundCover)

	b := 0.001 * math.Exp(0.053*groundCover)
	hillSlopeSedimentDeliveryRatio := 0.1336 * math.Exp((-1*b)*distanceToRiparianBuffer)

	vegetationCoverFactor := unmodifiedVegetationCoverFactor * hillSlopeSedimentDeliveryRatio

	return vegetationCoverFactor
}