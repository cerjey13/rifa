// utils/utils_test.go
package utils

import (
	"math"
	"strconv"
	"testing"
)

func assertEqualIntSlices(t *testing.T, got, want []int) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf(
			"len mismatch: got %d, want %d (got=%v, want=%v)",
			len(got),
			len(want),
			got,
			want,
		)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf(
				"index %d: got %d, want %d (got=%v, want=%v)",
				i,
				got[i],
				want[i],
				got,
				want,
			)
		}
	}
}

func assertEqualStringSlices(t *testing.T, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf(
			"len mismatch: got %d, want %d (got=%v, want=%v)",
			len(got),
			len(want),
			got,
			want,
		)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf(
				"index %d: got %q, want %q (got=%v, want=%v)",
				i,
				got[i],
				want[i],
				got,
				want,
			)
		}
	}
}

func TestParseFloatOrZero(t *testing.T) {
	tests := []struct {
		name       string
		in         string
		want       float64
		wantIsNaN  bool
		wantIsPosI bool
		wantIsNegI bool
	}{
		{"empty_is_zero", "", 0, false, false, false},
		{"garbage_is_zero", "abc", 0, false, false, false},
		{"int_string", "42", 42, false, false, false},
		{"negative", "-13", -13, false, false, false},
		{"float_string", "3.14159", 3.14159, false, false, false},
		{"sci_notation", "1.2e3", 1200, false, false, false},
		{"leading_zeros", "0007.50", 7.5, false, false, false},
		{"nan_passthrough", "NaN", 0, true, false, false},
		{"pos_inf_passthrough", "+Inf", 0, false, true, false},
		{"neg_inf_passthrough", "-Inf", 0, false, false, true},
	}

	const eps = 1e-9

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseFloatOrZero(tt.in)

			switch {
			case tt.wantIsNaN:
				if !math.IsNaN(got) {
					t.Fatalf("ParseFloatOrZero(%q) = %v, want NaN", tt.in, got)
				}
			case tt.wantIsPosI:
				if !math.IsInf(got, +1) {
					t.Fatalf("ParseFloatOrZero(%q) = %v, want +Inf", tt.in, got)
				}
			case tt.wantIsNegI:
				if !math.IsInf(got, -1) {
					t.Fatalf("ParseFloatOrZero(%q) = %v, want -Inf", tt.in, got)
				}
			default:
				if math.IsNaN(got) || math.IsInf(got, 0) {
					t.Fatalf(
						"ParseFloatOrZero(%q) = %v, want finite number",
						tt.in,
						got,
					)
				}
				if diff := math.Abs(got - tt.want); diff > eps {
					t.Fatalf(
						"ParseFloatOrZero(%q) = %v, want %v (|diff|=%g)",
						tt.in,
						got,
						tt.want,
						diff,
					)
				}
			}
		})
	}
}

func TestConvertToIntSlice(t *testing.T) {
	tests := []struct {
		name    string
		in      []string
		want    []int
		wantErr bool
	}{
		{
			name:    "success",
			in:      []string{"0", "007", "-12", "2147483647"},
			want:    []int{0, 7, -12, 2147483647},
			wantErr: false,
		},
		{
			name:    "error_on_invalid_element",
			in:      []string{"1", "two", "3"},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertToIntSlice(tt.in)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error for input %v, got nil", tt.in)
				}
				if got != nil {
					t.Fatalf("on error, expected nil slice; got %v", got)
				}
				return
			}

			if err != nil {
				t.Fatalf("ConvertToIntSlice(%v) error = %v", tt.in, err)
			}
			assertEqualIntSlices(t, got, tt.want)
		})
	}
}

func TestConvertToStrSlice(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		want []string
	}{
		{
			name: "typical_values",
			in:   []int{0, 7, -12, 2147483647},
			want: []string{"0", "7", "-12", "2147483647"},
		},
		{
			name: "empty_slice",
			in:   []int{},
			want: []string{},
		},
		{
			name: "nil_slice_behaves_as_empty",
			in:   nil,
			want: []string{},
		},
		{
			name: "single_element",
			in:   []int{-1},
			want: []string{"-1"},
		},
		{
			name: "mixed_values",
			in:   []int{-5, 0, 10, 999999},
			want: []string{"-5", "0", "10", "999999"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertToStrSlice(tt.in)

			if len(got) != len(tt.in) {
				t.Fatalf("len mismatch: got %d, want %d", len(got), len(tt.in))
			}
			for i, s := range got {
				if s != strconv.Itoa(func() int {
					if tt.in == nil {
						return 0 // won't run because len==0 when nil
					}
					return tt.in[i]
				}()) {
					t.Fatalf("itoa check failed at %d: got %q", i, s)
				}
			}

			assertEqualStringSlices(t, got, tt.want)
		})
	}
}

func archMaxInt() int { return int(^uint(0) >> 1) }
func archMinInt() int { return -archMaxInt() - 1 }

func TestRoundTrip_IntsToStringsAndBack(t *testing.T) {
	maxI, minI := archMaxInt(), archMinInt()

	tests := []struct {
		name string
		in   []int
	}{
		{"empty", []int{}},
		{"nil_behaves_empty", nil},
		{"single_zero", []int{0}},
		{"single_negative", []int{-1}},
		{"typical_mixed", []int{-5, 0, 1, 10, 999999}},
		{"edges", []int{minI, -123456789, 0, 123456789, maxI}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strs := ConvertToStrSlice(tt.in)
			back, err := ConvertToIntSlice(strs)
			if err != nil {
				t.Fatalf("unexpected error on round-trip: %v", err)
			}
			assertEqualIntSlices(t, back, tt.in)
		})
	}
}

func TestRoundTrip_FailureOnTamperedStrings(t *testing.T) {
	// Start from valid numbers, then tamper one element to make parsing fail.
	valid := []int{-5, 0, 42, 999}
	base := ConvertToStrSlice(valid)

	overflowStr := func() string {
		// craft a value guaranteed to overflow for current int size
		if strconv.IntSize == 64 {
			// MaxInt64+1
			return "9223372036854775808"
		}
		// MaxInt32+1
		return "2147483648"
	}()

	tests := []struct {
		name string
		edit func([]string) []string
	}{
		{"non_numeric", func(s []string) []string { s[2] = "two"; return s }},
		{"empty_string", func(s []string) []string { s[1] = ""; return s }},
		{"embedded_space", func(s []string) []string { s[0] = "4 2"; return s }},
		{"leading_space", func(s []string) []string { s[0] = " 1"; return s }},
		{"trailing_space", func(s []string) []string { s[3] = "999 "; return s }},
		{"decimal_value", func(s []string) []string { s[1] = "1.0"; return s }},
		{"hex_prefix", func(s []string) []string { s[2] = "0x10"; return s }},
		{"thousands_separator", func(s []string) []string {
			s[3] = "1,000"
			return s
		}},
		// Arabic-Indic
		{"unicode_digits", func(s []string) []string { s[2] = "١٢٣"; return s }},
		{"overflow_for_arch", func(s []string) []string {
			s[1] = overflowStr
			return s
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// copy to avoid mutating base between cases
			mut := make([]string, len(base))
			copy(mut, base)
			mut = tt.edit(mut)

			got, err := ConvertToIntSlice(mut)
			if err == nil {
				t.Fatalf("expected error for tampered input %v, got nil", mut)
			}
			if got != nil {
				t.Fatalf("on error, expected nil slice; got %v", got)
			}
		})
	}
}

// -------- benchmarks --------
// Sizes chosen to simulate small/medium/large payloads.

func benchIntData(n int) []int {
	out := make([]int, n)
	// deterministically fill with a mix of small/large, +/- values
	for i := range n {
		v := (i*i + 17*i + 23)
		if i%2 == 0 {
			v = -v
		}
		out[i] = v
	}
	return out
}

func benchStrData(n int) []string {
	ints := benchIntData(n)
	return ConvertToStrSlice(ints)
}

func BenchmarkConvertToStrSlice_10(b *testing.B) {
	benchmarkConvertToStrSlice(b, 10)
}
func BenchmarkConvertToStrSlice_100(b *testing.B) {
	benchmarkConvertToStrSlice(b, 100)
}
func BenchmarkConvertToStrSlice_1000(b *testing.B) {
	benchmarkConvertToStrSlice(b, 1000)
}
func BenchmarkConvertToStrSlice_10000(b *testing.B) {
	benchmarkConvertToStrSlice(b, 10000)
}
func BenchmarkConvertToIntSlice_10(b *testing.B) {
	benchmarkConvertToIntSlice(b, 10)
}
func BenchmarkConvertToIntSlice_100(b *testing.B) {
	benchmarkConvertToIntSlice(b, 100)
}
func BenchmarkConvertToIntSlice_1000(b *testing.B) {
	benchmarkConvertToIntSlice(b, 1000)
}
func BenchmarkConvertToIntSlice_10000(b *testing.B) {
	benchmarkConvertToIntSlice(b, 10000)
}
func BenchmarkRoundTrip_10(b *testing.B)    { benchmarkRoundTrip(b, 10) }
func BenchmarkRoundTrip_100(b *testing.B)   { benchmarkRoundTrip(b, 100) }
func BenchmarkRoundTrip_1000(b *testing.B)  { benchmarkRoundTrip(b, 1000) }
func BenchmarkRoundTrip_10000(b *testing.B) { benchmarkRoundTrip(b, 10000) }

func benchmarkConvertToStrSlice(b *testing.B, n int) {
	data := benchIntData(n)
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_ = ConvertToStrSlice(data)
	}
}

func benchmarkConvertToIntSlice(b *testing.B, n int) {
	data := benchStrData(n)
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		got, err := ConvertToIntSlice(data)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
		if len(got) != n {
			b.Fatalf("len(got)=%d, want=%d", len(got), n)
		}
	}
}

func benchmarkRoundTrip(b *testing.B, n int) {
	ints := benchIntData(n)
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		strs := ConvertToStrSlice(ints)
		back, err := ConvertToIntSlice(strs)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
		if len(back) != n {
			b.Fatalf("len(back)=%d, want=%d", len(back), n)
		}
	}
}
