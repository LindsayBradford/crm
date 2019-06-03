// Copyright (c) 2019 Australian Rivers Institute.

package archive

import (
	"github.com/LindsayBradford/crem/internal/pkg/rand"
	"testing"

	. "github.com/onsi/gomega"
)

const equalTo = "=="

func TestBankSedimentContribution_New(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	expectedSize := 130
	expectedArchiveSize := 3

	// when
	archiveUnderTest := New(expectedSize)

	// then
	g.Expect(archiveUnderTest.Len()).To(BeNumerically(equalTo, expectedSize))
	g.Expect(archiveUnderTest.ArchiveLen()).To(BeNumerically(equalTo, expectedArchiveSize))

	for index := 0; index < expectedSize; index++ {
		g.Expect(archiveUnderTest.Value(index)).To(BeFalse())
	}
}

func TestBankSedimentContribution_SetValue_OverArchiveRage(t *testing.T) {
	g := NewGomegaWithT(t)
	random := rand.NewTimeSeeded()

	// given
	expectedSize := 130
	archiveUnderTest := New(expectedSize)

	numberToSetTrue := 7
	expectedTrueIndexes := make([]int, numberToSetTrue)
	for current := 0; current < numberToSetTrue; current++ {
		indexToSetTrue := random.Intn(expectedSize)
		for previous := 0; previous < current; previous++ {
			duplicateIndexFound := true
			for duplicateIndexFound {
				if expectedTrueIndexes[previous] == indexToSetTrue {
					indexToSetTrue = random.Intn(expectedSize)
				} else {
					duplicateIndexFound = false
				}
			}
		}
		expectedTrueIndexes[current] = indexToSetTrue
	}
	t.Logf("Archive indexes that should be set to true: %v", expectedTrueIndexes)

	// when
	for current := 0; current < numberToSetTrue; current++ {
		archiveUnderTest.SetValue(expectedTrueIndexes[current], true)
	}

	actualTrueIndexes := make([]int, 0)
	falseCount := 0
	for index := 0; index < expectedSize; index++ {
		if archiveUnderTest.Value(index) == false {
			falseCount++
		} else {
			actualTrueIndexes = append(actualTrueIndexes, index)
		}
	}
	t.Logf("Archive indexes are set to true: %v", actualTrueIndexes)

	// then
	g.Expect(expectedTrueIndexes).To(ConsistOf(actualTrueIndexes))
	g.Expect(falseCount).To(BeNumerically(equalTo, expectedSize-numberToSetTrue))
}

func TestBankSedimentContribution_SetValue_Toggling(t *testing.T) {
	g := NewGomegaWithT(t)

	expectedSize := 5
	archiveUnderTest := New(expectedSize)
	indexToTest := 2

	// given
	valuesToAssign := []bool{false, false, true, true, false, false, true, true}

	for _, testValue := range valuesToAssign {
		// when
		archiveUnderTest.SetValue(indexToTest, testValue)

		// then
		if testValue {
			g.Expect(archiveUnderTest.Value(indexToTest)).To(BeTrue())
		} else {
			g.Expect(archiveUnderTest.Value(indexToTest)).To(BeFalse())
		}
	}
}

func TestBankSedimentContribution_SetValue_OutsideValidRange(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	expectedSize := 5
	archiveUnderTest := New(expectedSize)

	// when
	indexToTest := 4
	archiveUnderTest.SetValue(indexToTest, true)

	// then
	g.Expect(archiveUnderTest.Value(indexToTest)).To(BeTrue())

	// when
	indexToTest = 5
	outOfBoundsSet := func() {
		archiveUnderTest.SetValue(indexToTest, true)
	}

	// then
	g.Expect(outOfBoundsSet).To(Panic())
}

func TestBankSedimentContribution_Value_OutsideValidRange(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	expectedSize := 5
	archiveUnderTest := New(expectedSize)

	// when
	indexToTest := 4
	expectedValue := archiveUnderTest.Value(indexToTest)

	// then
	g.Expect(expectedValue).To(BeFalse())

	// when
	indexToTest = 5
	outOfBoundsValue := func() {
		archiveUnderTest.Value(indexToTest)
	}

	// then
	g.Expect(outOfBoundsValue).To(Panic())
}