package mana

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func BenchmarkCalculator_GenerateMana(b *testing.B) {
	calculator := NewCalculator(10, 0.1)

	for i := 0; i < b.N; i++ {
		calculator.GenerateMana(1000000, 100000000000)
	}
}

func TestCalculator_GenerateMana(t *testing.T) {
	calculator := NewCalculator(500, 0.1, ManaScaleFactor(2))

	generatedMana, _ := calculator.GenerateMana(1000, 0)
	assert.Equal(t, generatedMana, uint64(0))

	generatedMana, _ = calculator.GenerateMana(1000, 500)
	assert.Equal(t, generatedMana, uint64(200))

	generatedMana, _ = calculator.GenerateMana(1000, 5000000)
	assert.Equal(t, generatedMana, uint64(2000))
}

func TestCalculator_ManaSymmetry(t *testing.T) {
	calculator := NewCalculator(500, 0.1, ManaScaleFactor(2))

	// 1st case: generate mana by spending two times
	generatedManaStep1, _ := calculator.GenerateMana(1000, 500)
	generatedManaStep2, _ := calculator.GenerateMana(1000, 500)
	generatedManaStep3, _ := calculator.GenerateMana(1000, 500)

	// the first "realized" mana part starts decaying while the coins of the 2nd spend are gaining weight again
	erodedMana1, _ := calculator.ErodeMana(generatedManaStep1, 1000)
	erodedMana2, _ := calculator.ErodeMana(generatedManaStep2, 500)

	// 2nd case: generate mana by spending only once
	generatedManaWithoutSpends, _ := calculator.GenerateMana(1000, 1500)

	// the two mana values should be equal
	assert.Equal(t, generatedManaWithoutSpends, erodedMana1+erodedMana2+generatedManaStep3)
}
