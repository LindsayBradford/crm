// Copyright (c) 2019 Australian Rivers Institute.

package archive

import (
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
	"testing"

	"github.com/LindsayBradford/crem/internal/pkg/model/models/modumb"
	. "github.com/onsi/gomega"
)

const equalTo = "=="

func TestArchivist_Store(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	archivistUnderTest := new(Archivist)
	testModel := modumb.NewModel()

	expectedVariableSize := len(*testModel.DecisionVariables())
	expectedActionSize := len(testModel.ManagementActions())

	// when
	actualArchive := archivistUnderTest.Store(testModel)

	// then
	g.Expect(expectedVariableSize).To(BeNumerically(equalTo, len(*actualArchive.Variables())))
	g.Expect(expectedActionSize).To(BeNumerically(equalTo, actualArchive.Actions().Len()))
}

func TestArchivist_Retrieve_InitialModel(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	archivistUnderTest := new(Archivist)

	modelToStore := buildMultiObjectiveDumbModel()
	modelToRetrieve := modelToStore.DeepClone()

	g.Expect(modelToStore.DecisionVariables()).To(Equal(modelToRetrieve.DecisionVariables()))
	g.Expect(modelToStore.ManagementActions()).To(Equal(modelToRetrieve.ManagementActions()))

	// when
	storedArchive := archivistUnderTest.Store(modelToStore)
	archivistUnderTest.Retrieve(storedArchive, modelToRetrieve)

	// then
	g.Expect(modelToStore.DecisionVariables()).To(Equal(modelToRetrieve.DecisionVariables()))
	g.Expect(modelToStore.ManagementActions()).To(Equal(modelToRetrieve.ManagementActions()))
}

func TestArchivist_Retrieve_AlteredModel(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	archivistUnderTest := new(Archivist)
	modelToStore := buildMultiObjectiveDumbModel()

	numberOfRandomChanges := 7
	for change := 0; change < numberOfRandomChanges; change++ {
		modelToStore.DoRandomChange()
	}

	modelToRetrieve := modelToStore.DeepClone()

	// when
	storedArchive := archivistUnderTest.Store(modelToStore)
	archivistUnderTest.Retrieve(storedArchive, modelToRetrieve)

	// then
	g.Expect(modelToStore.DecisionVariables()).To(Equal(modelToRetrieve.DecisionVariables()))
	g.Expect(modelToStore.ManagementActions()).To(Equal(modelToRetrieve.ManagementActions()))
}

func buildMultiObjectiveDumbModel() *modumb.Model {
	model := modumb.NewModel().WithId("Test Mo Dumb Model")
	model.SetEventNotifier(loggers.DefaultTestingEventNotifier)
	model.Initialise()
	return model
}
