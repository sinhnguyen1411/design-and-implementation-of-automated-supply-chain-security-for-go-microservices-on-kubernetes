package fees

import "testing"

func TestCalculate(t *testing.T) {
	r := Calculate(10000, TierStandard)
	if r.Commission != 10.0 { t.Fatalf("want 10.0, got %f", r.Commission) }
	if r.NetAmount != 9990.0 { t.Fatalf("want 9990.0, got %f", r.NetAmount) }
}

func TestCalculatePro(t *testing.T) {
	r := Calculate(10000, TierPro)
	if r.Commission != 5.0 { t.Fatalf("want 5.0, got %f", r.Commission) }
}
