package kyc_test

import (
	"testing"

	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/kyc-service/internal/kyc"
)

func TestCreate(t *testing.T) {
	r := kyc.NewRegistry()
	p := r.Create("user-1", "US")
	if p.Status != kyc.StatusPending {
		t.Fatalf("mong đợi pending, nhận được %s", p.Status)
	}
	if p.RiskScore != 0 {
		t.Fatalf("risk score ban đầu phải là 0, nhận được %d", p.RiskScore)
	}
}

func TestApprove(t *testing.T) {
	r := kyc.NewRegistry()
	r.Create("user-1", "US")
	r.AddDocument("user-1", kyc.DocPassport, nil)
	if err := r.Approve("user-1"); err != nil {
		t.Fatal(err)
	}
	p, _ := r.Get("user-1")
	if p.Status != kyc.StatusVerified {
		t.Fatalf("mong đợi verified, nhận được %s", p.Status)
	}
	for _, doc := range p.Documents {
		if !doc.Verified {
			t.Fatal("tất cả document phải được verified khi approve")
		}
	}
}

func TestReject(t *testing.T) {
	r := kyc.NewRegistry()
	r.Create("user-1", "US")
	if err := r.Reject("user-1", "document mismatch"); err != nil {
		t.Fatal(err)
	}
	p, _ := r.Get("user-1")
	if p.Status != kyc.StatusRejected {
		t.Fatalf("mong đợi rejected, nhận được %s", p.Status)
	}
	if p.Reason == "" {
		t.Fatal("lý do từ chối phải được lưu lại")
	}
}

func TestFlagAML_riskScore(t *testing.T) {
	r := kyc.NewRegistry()
	r.Create("user-1", "US")

	r.FlagAML("user-1", "UNUSUAL_ACTIVITY", "large transfer", kyc.RiskMedium)
	p, _ := r.Get("user-1")
	if p.RiskScore != 20 {
		t.Fatalf("mong đợi risk score 20, nhận được %d", p.RiskScore)
	}

	r.FlagAML("user-1", "PEP", "politically exposed person", kyc.RiskHigh)
	p, _ = r.Get("user-1")
	if p.RiskScore != 60 {
		t.Fatalf("mong đợi risk score 60, nhận được %d", p.RiskScore)
	}
	if p.RiskLevel != kyc.RiskHigh {
		t.Fatalf("mong đợi high risk, nhận được %s", p.RiskLevel)
	}
}

func TestFlagAML_scoreCap(t *testing.T) {
	r := kyc.NewRegistry()
	r.Create("user-1", "US")
	for i := 0; i < 5; i++ {
		r.FlagAML("user-1", "SANCTIONS", "", kyc.RiskHigh)
	}
	p, _ := r.Get("user-1")
	if p.RiskScore > 100 {
		t.Fatalf("risk score phải được giới hạn ở 100, nhận được %d", p.RiskScore)
	}
}

func TestFlagsFor(t *testing.T) {
	r := kyc.NewRegistry()
	r.Create("user-1", "US")
	r.Create("user-2", "VN")
	r.FlagAML("user-1", "PEP", "", kyc.RiskHigh)
	r.FlagAML("user-2", "UNUSUAL_ACTIVITY", "", kyc.RiskLow)
	r.FlagAML("user-1", "SANCTIONS", "", kyc.RiskHigh)

	flags := r.FlagsFor("user-1")
	if len(flags) != 2 {
		t.Fatalf("mong đợi 2 flags cho user-1, nhận được %d", len(flags))
	}
}

func TestGet_notFound(t *testing.T) {
	r := kyc.NewRegistry()
	_, err := r.Get("nonexistent")
	if err == nil {
		t.Fatal("phải trả về lỗi khi profile không tồn tại")
	}
}
