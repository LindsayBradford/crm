// Copyright (c) 2019 Australian Rivers Institute.

package archive

import (
	model2 "github.com/LindsayBradford/crem/internal/pkg/model"
	"testing"

	"github.com/LindsayBradford/crem/internal/pkg/model/models/modumb"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
	. "github.com/onsi/gomega"
)

func TestModelCompressor_Compress(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	compressorUnderTest := new(ModelCompressor)
	testModel := modumb.NewModel()

	// when
	compressedModelState := compressorUnderTest.Compress(testModel)

	// then
	g.Expect(compressedModelState.MatchesStateOf(testModel)).To(BeTrue())
}

func TestModelCompressor_Decompress_InitialModel(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	compressorUnderTest := new(ModelCompressor)
	originalModel := buildMultiObjectiveDumbModel()
	decompressedModel := originalModel.DeepClone()

	// when
	compressedModelState := compressorUnderTest.Compress(originalModel)
	compressorUnderTest.Decompress(compressedModelState, decompressedModel)

	// then
	g.Expect(decompressedModel.IsEquivalentTo(originalModel)).To(BeTrue())
	g.Expect(compressedModelState.MatchesStateOf(originalModel)).To(BeTrue())
}

func TestModelCompressor_Decompress_AlteredModel(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	compressorUnderTest := new(ModelCompressor)

	originalModel := buildMultiObjectiveDumbModel()

	modifiedModel := buildMultiObjectiveDumbModel()
	numberOfRandomChanges := 7
	for change := 0; change < numberOfRandomChanges; change++ {
		modifiedModel.DoRandomChange()
	}

	decompressedModel := buildMultiObjectiveDumbModel()

	// when
	compressedModifiedModelState := compressorUnderTest.Compress(modifiedModel)
	compressorUnderTest.Decompress(compressedModifiedModelState, decompressedModel)

	// then
	g.Expect(originalModel.IsEquivalentTo(modifiedModel)).To(BeFalse())
	g.Expect(compressedModifiedModelState.MatchesStateOf(originalModel)).To(BeFalse())

	g.Expect(compressedModifiedModelState.MatchesStateOf(modifiedModel)).To(BeTrue())
	g.Expect(decompressedModel.IsEquivalentTo(modifiedModel)).To(BeTrue())
}

func buildMultiObjectiveDumbModel() *modumb.Model {
	model := modumb.NewModel().WithId("Test Mo Dumb Model")
	model.AddObserver(loggers.DefaultTestingAnnealingObserver)
	model.Initialise(model2.Random)
	return model
}
