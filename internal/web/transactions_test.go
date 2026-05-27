package web

import (
	"net/http/httptest"
	"testing"
	"time"

	"cashout/internal/model"
)

func TestIsIncomeCategory(t *testing.T) {
	cases := []struct {
		cat  string
		want bool
	}{
		{string(model.CategorySalary), true},
		{string(model.CategoryOtherIncomes), true},
		{string(model.CategoryGrocery), false},
		{string(model.CategoryTravel), false},
		{"NotACategory", false},
		{"", false},
	}
	for _, c := range cases {
		if got := isIncomeCategory(c.cat); got != c.want {
			t.Errorf("isIncomeCategory(%q) = %v; want %v", c.cat, got, c.want)
		}
	}
}

func TestToTransactionDTO(t *testing.T) {
	d := time.Date(2026, 5, 21, 0, 0, 0, 0, time.UTC)
	tx := model.Transaction{
		ID:          7,
		TgID:        42,
		Date:        d,
		Type:        model.TypeExpense,
		Category:    model.CategoryGrocery,
		Amount:      19.90,
		Description: "weekly shop",
	}
	dto := toTransactionDTO(tx)
	if dto.ID != 7 || dto.Date != d || dto.Category != string(model.CategoryGrocery) ||
		dto.Description != "weekly shop" || dto.Amount != 19.90 || dto.Type != string(model.TypeExpense) {
		t.Fatalf("unexpected DTO: %+v", dto)
	}
}

func TestParseFloatQuery(t *testing.T) {
	t.Run("absent", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/x", nil)
		v, err := parseFloatQuery(r, "amountMin")
		if err != nil || v != nil {
			t.Fatalf("want nil/nil, got %v, %v", v, err)
		}
	})
	t.Run("valid", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/x?amountMin=12.5", nil)
		v, err := parseFloatQuery(r, "amountMin")
		if err != nil || v == nil || *v != 12.5 {
			t.Fatalf("want 12.5, got %v, %v", v, err)
		}
	})
	t.Run("invalid", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/x?amountMin=oops", nil)
		_, err := parseFloatQuery(r, "amountMin")
		if err == nil {
			t.Fatalf("want error")
		}
	})
}

func ptrFloat(v float64) *float64 { return &v }

func TestBuildSearchFilter(t *testing.T) {
	type tc struct {
		name      string
		query     string
		category  string
		txType    string
		dateFrom  string
		dateTo    string
		amountMin *float64
		amountMax *float64
		wantCode  int
	}

	cases := []tc{
		{name: "empty is fine"},
		{name: "category all is no-op", category: "all"},
		{name: "invalid category", category: "Bogus", wantCode: 400},
		{name: "valid category", category: string(model.CategoryGrocery)},
		{name: "invalid type", txType: "Maybe", wantCode: 400},
		{name: "valid income type", txType: string(model.TypeIncome)},
		{name: "invalid dateFrom", dateFrom: "21/05/2026", wantCode: 400},
		{name: "valid dateFrom", dateFrom: "2026-05-21"},
		{name: "dateFrom after dateTo", dateFrom: "2026-05-21", dateTo: "2026-05-01", wantCode: 400},
		{name: "amountMin negative", amountMin: ptrFloat(-1), wantCode: 400},
		{name: "amountMin > amountMax", amountMin: ptrFloat(50), amountMax: ptrFloat(10), wantCode: 400},
		{name: "valid amount range", amountMin: ptrFloat(5), amountMax: ptrFloat(100)},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, code, msg := buildSearchFilter(c.query, c.category, c.txType, c.dateFrom, c.dateTo, c.amountMin, c.amountMax)
			if code != c.wantCode {
				t.Fatalf("code=%d msg=%q; want code=%d", code, msg, c.wantCode)
			}
		})
	}
}

func TestBuildSearchFilter_DateToIsEndOfDay(t *testing.T) {
	f, code, _ := buildSearchFilter("", "", "", "", "2026-05-21", nil, nil)
	if code != 0 {
		t.Fatalf("unexpected error code %d", code)
	}
	if f.DateTo == nil {
		t.Fatalf("DateTo nil")
	}
	// Should be the last nanosecond of 2026-05-21 UTC.
	want := time.Date(2026, 5, 22, 0, 0, 0, 0, time.UTC).Add(-time.Nanosecond)
	if !f.DateTo.Equal(want) {
		t.Fatalf("DateTo = %v; want %v", *f.DateTo, want)
	}
}
