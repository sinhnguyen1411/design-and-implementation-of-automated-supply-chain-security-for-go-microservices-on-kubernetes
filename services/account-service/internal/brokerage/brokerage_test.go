package brokerage

import "testing"

func TestDepositWithdraw(t *testing.T) {
	acc := New("owner-1", TypeCash)
	if acc.CashBalance != 0 { t.Fatal("initial balance should be 0") }

	tx, err := acc.Deposit(5000)
	if err != nil { t.Fatalf("deposit failed: %v", err) }
	if acc.CashBalance != 5000 { t.Fatalf("want 5000, got %f", acc.CashBalance) }
	if tx.Kind != "deposit" { t.Fatal("wrong tx kind") }

	tx, err = acc.Withdraw(2000)
	if err != nil { t.Fatalf("withdraw failed: %v", err) }
	if acc.CashBalance != 3000 { t.Fatalf("want 3000, got %f", acc.CashBalance) }
	if tx.Balance != 3000 { t.Fatalf("tx balance want 3000, got %f", tx.Balance) }
}

func TestWithdraw_insufficient(t *testing.T) {
	acc := New("owner-1", TypeCash)
	acc.Deposit(100)
	_, err := acc.Withdraw(500)
	if err == nil { t.Fatal("should fail: insufficient funds") }
}

func TestDeposit_suspended(t *testing.T) {
	acc := New("owner-1", TypeCash)
	acc.Suspend()
	_, err := acc.Deposit(100)
	if err == nil { t.Fatal("should fail: account suspended") }
}

func TestBuyingPower_margin(t *testing.T) {
	acc := New("owner-1", TypeMargin)
	acc.Deposit(10000)
	acc.MarginBalance = 5000
	if acc.BuyingPower() != 15000 { t.Fatalf("want 15000, got %f", acc.BuyingPower()) }
}

func TestBuyingPower_cash(t *testing.T) {
	acc := New("owner-1", TypeCash)
	acc.Deposit(8000)
	acc.MarginBalance = 4000
	if acc.BuyingPower() != 8000 { t.Fatalf("cash account buying power should equal cash only") }
}
