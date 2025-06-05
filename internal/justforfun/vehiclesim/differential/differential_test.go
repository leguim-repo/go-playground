package differential

import (
	"testing"
)

// TestDifferentialBehavior demonstrates the behavior of the differential under different slipRatio values
// For now, only print the results for manual observation.
func TestDifferentialBehavior(t *testing.T) {
	diff := NewBasicDifferential(TypeRDiffRatio)

	testSlipRatios := []float64{0.0, 0.05, -0.05, 0.25, 2.0}
	inputRPM := 2000.0
	inputTorque := 150.0

	// The gear ratio is taken from the differential itself.
	gearRatio := diff.gearRatio
	expectedBaseRPM := inputRPM / gearRatio

	t.Log("--- Testing differential behavior ---")
	t.Logf("Input RPM: %.2f, Gear Ratio: %.2f, Expected Base Wheel RPM: %.2f", inputRPM, gearRatio, expectedBaseRPM)
	t.Log("---------------------------------------------------------------------")

	for _, slip := range testSlipRatios {
		diff.Update(inputRPM, inputTorque, slip)
		data := diff.GetData()

		// Use Logf for print results for manual observation only with the test is launched with flag -v
		t.Logf("SlipRatio: %-5.2f -> Left wheel: %7.2f RPM, Right wheel: %7.2f RPM",
			slip, data.WheelSpeedL, data.WheelSpeedR)

		// Check that the wheel speeds are equal for slipRatio = 0.0
		if slip == 0.0 {
			if data.WheelSpeedL != data.WheelSpeedR {
				//
				t.Errorf("For slipRatio 0.0, expected WheelSpeedL (%.2f) == WheelSpeedR (%.2f), but it wasn't like that", data.WheelSpeedL, data.WheelSpeedR)
			}
		}
	}
}
