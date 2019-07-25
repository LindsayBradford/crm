// Copyright (c) 2019 Australian Rivers Institute.

package actions

import (
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
)

type RiverBankRestorationGroup struct {
	planningUnitTable tables.CsvTable
	parameters        parameters.Parameters

	actionMap map[planningunit.Id]*RiverBankRestoration
}

func (r *RiverBankRestorationGroup) Initialise(planningUnitTable tables.CsvTable, parameters parameters.Parameters) *RiverBankRestorationGroup {
	r.planningUnitTable = planningUnitTable
	r.parameters = parameters
	r.createManagementActions()

	return r
}

func (r *RiverBankRestorationGroup) ManagementActions() map[planningunit.Id]*RiverBankRestoration {
	return r.actionMap
}

func (r *RiverBankRestorationGroup) createManagementActions() {
	_, rowCount := r.planningUnitTable.ColumnAndRowSize()
	r.actionMap = make(map[planningunit.Id]*RiverBankRestoration, rowCount)

	for row := uint(0); row < rowCount; row++ {
		r.createManagementAction(row)
	}
}

func (r *RiverBankRestorationGroup) createManagementAction(rowNumber uint) {
	planningUnit := r.planningUnitTable.CellFloat64(planningUnitIndex, rowNumber)
	planningUnitAsId := planningunit.Float64ToId(planningUnit)

	vegetationTarget := r.parameters.GetFloat64(parameters.RiparianBufferVegetationProportionTarget)

	originalBufferVegetation := r.originalBufferVegetation(rowNumber)

	if originalBufferVegetation >= vegetationTarget {
		return
	}

	costInDollars := r.calculateImplementationCost(rowNumber)

	r.actionMap[planningUnitAsId] =
		NewRiverBankRestoration().
			WithPlanningUnit(planningUnitAsId).
			WithoRIGINALBufferVegetation(originalBufferVegetation).
			WithActionedBufferVegetation(vegetationTarget).
			WithImplementationCost(costInDollars)
}

func (r *RiverBankRestorationGroup) originalBufferVegetation(rowNumber uint) float64 {
	proportionOfRiparianVegetation := r.planningUnitTable.CellFloat64(riparianVegetationIndex, rowNumber)
	return proportionOfRiparianVegetation
}

func (r *RiverBankRestorationGroup) calculateChangeInBufferVegetation(rowNumber uint) float64 {
	proportionOfRiparianVegetation := r.originalBufferVegetation(rowNumber)
	vegetationTarget := r.parameters.GetFloat64(parameters.RiparianBufferVegetationProportionTarget)
	changeInRiparianVegetation := vegetationTarget - proportionOfRiparianVegetation
	return changeInRiparianVegetation
}

func (r *RiverBankRestorationGroup) calculateImplementationCost(rowNumber uint) float64 {
	riparianRevegetationCostPerKlmSquared := r.parameters.GetFloat64(parameters.RiparianRevegetationCostPerKilometer)
	riverLengthInMetres := r.planningUnitTable.CellFloat64(riverLengthIndex, rowNumber)
	riverLengthInKilometres := riverLengthInMetres / 1000

	vegetationChange := r.calculateChangeInBufferVegetation(rowNumber)

	vegetationChangeLengthInKilometresSquared := vegetationChange * riverLengthInKilometres

	implementationCost := vegetationChangeLengthInKilometresSquared * riparianRevegetationCostPerKlmSquared

	return implementationCost
}