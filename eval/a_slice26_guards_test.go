package eval

import (
	"math"
	"testing"

	"github.com/go-statistics/stats"
)

// These tests sit in a file that sorts before the others: they turn a crash in
// slice26.go into an assertion, so a mutation of one of these guards is
// reported by what it breaks and not by the panic that stops the run.

func noPanic(t *testing.T, what string, f func()) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("%s panicked: %v", what, r)
		}
	}()
	f()
}

func TestSlice26FormalsShape(t *testing.T) {
	for i := range families {
		f := &families[i]
		d, p, q := f.formals(kindD), f.formals(kindP), f.formals(kindQ)
		if len(d) == 0 || d[len(d)-1] != "log" {
			t.Errorf("%s: d formals %v do not end with log", f.name, d)
		}
		for _, n := range d {
			if n == "lower.tail" || n == "log.p" {
				t.Errorf("%s: d formals %v hold %s", f.name, d, n)
			}
		}
		for _, l := range [][]string{p, q} {
			if len(l) < 3 || l[len(l)-2] != "lower.tail" || l[len(l)-1] != "log.p" {
				t.Errorf("%s: p/q formals %v do not end with lower.tail, log.p", f.name, l)
			}
		}
		if len(p) != len(d)+1 || len(q) != len(p) {
			t.Errorf("%s: formals d %v p %v q %v", f.name, d, p, q)
		}
	}
}

func TestSlice26EmptyInputs(t *testing.T) {
	for _, c := range []struct{ src, want string }{
		{`length(pnorm(numeric(0)))`, "[1] 0\n"},
		{`length(dgamma(1, numeric(0)))`, "[1] 0\n"},
		{`length(qgamma(numeric(0), 2))`, "[1] 0\n"},
		{`length(pbeta(0.5, 2, numeric(0)))`, "[1] 0\n"},
	} {
		noPanic(t, c.src, func() {
			res, err := NewSession(Options{}).Eval(c.src)
			if err != nil {
				t.Errorf("%s: %v", c.src, err)
				return
			}
			if res.Output != c.want {
				t.Errorf("%s printed %q, want %q", c.src, res.Output, c.want)
			}
		})
	}
}

func TestSlice26ValueRefusedBuild(t *testing.T) {
	norm := fam(t, "norm")
	noPanic(t, "value with a refused build", func() {
		for _, k := range []distKind{kindD, kindP, kindQ} {
			if v := norm.value(k, 0.5, []float64{0, -1}, true, false); !math.IsNaN(v) {
				t.Errorf("norm kind %v with sd -1 = %v, want NaN", k, v)
			}
		}
	})
}

func TestSlice26QuantilePlain(t *testing.T) {
	d, ok := mustNew(stats.NewNormal(0, 1))
	if !ok {
		t.Fatal("normal")
	}
	noPanic(t, "plain quantile", func() {
		if got, want := quantile(d, 0.3, true, false), d.PPF(0.3); got != want {
			t.Errorf("lower = %v, want %v", got, want)
		}
		if got, want := quantile(d, 0.3, false, false), d.ISF(0.3); got != want {
			t.Errorf("upper = %v, want %v", got, want)
		}
	})
}

func TestSlice26QuantileInvalidProb(t *testing.T) {
	pm := pointMass{2}
	for _, c := range []struct {
		p           float64
		lower, logp bool
		wantNaN     bool
	}{
		{1.5, true, false, true},
		{-0.5, true, false, true},
		{0.5, true, false, false},
		{0.5, true, true, true},
		{-1, true, true, false},
		{math.NaN(), true, false, true},
	} {
		got := quantile(pm, c.p, c.lower, c.logp)
		if math.IsNaN(got) != c.wantNaN {
			t.Errorf("quantile(point mass, %v, lower=%v, log=%v) = %v, NaN wanted: %v", c.p, c.lower, c.logp, got, c.wantNaN)
		}
	}
}

// An infinite x equal to an infinite location is NaN in the density: z is NaN
// and every standard law carries it through.
func TestSlice26InfEqualsLoc(t *testing.T) {
	for _, name := range []string{"norm", "logis", "cauchy"} {
		for _, loc := range []float64{math.Inf(1), math.Inf(-1)} {
			d, ok := fam(t, name).build([]float64{loc, 1})
			if !ok {
				continue
			}
			if v := d.PDF(loc); !math.IsNaN(v) {
				t.Errorf("%s loc %v: PDF(loc) = %v, want NaN", name, loc, v)
			}
			if v := d.LogPDF(loc); !math.IsNaN(v) {
				t.Errorf("%s loc %v: LogPDF(loc) = %v, want NaN", name, loc, v)
			}
		}
	}
}
