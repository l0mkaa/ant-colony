package simulation

import (
	"testing"
)

func BenchmarkStep(b *testing.B) {
	s := NewSimulation(600, 600, defaultSimulationVars)
	for i := 0; i < b.N; i++ {
		s.step()
	}
}
