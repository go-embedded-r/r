package eval

import (
	"math"
	"slices"
	"testing"

	"github.com/go-statistics/stats"
)

// Direct tests of the laws and the guards of slice26.go: what the R-level
// goldens cannot tell apart because a later layer gives the same answer for
// the points they reach.

var (
	inf  = math.Inf(1)
	ninf = math.Inf(-1)
	nan  = math.NaN()
	nz   = math.Copysign(0, -1)
)

// same is equality of float64 with NaN equal to NaN and the sign of zero held.
func same(a, b float64) bool {
	if math.IsNaN(a) || math.IsNaN(b) {
		return math.IsNaN(a) && math.IsNaN(b)
	}
	return a == b && math.Signbit(a) == math.Signbit(b)
}

func fam(t *testing.T, name string) *family {
	t.Helper()
	for i := range families {
		if families[i].name == name {
			return &families[i]
		}
	}
	t.Fatalf("no family %s", name)
	return nil
}

func TestSlice26BuildGuards(t *testing.T) {
	cases := []struct {
		fam string
		p   []float64
		ok  bool
	}{
		{"norm", []float64{0, -1}, false},
		{"norm", []float64{0, ninf}, false},
		{"norm", []float64{0, 0}, true},
		{"norm", []float64{0, inf}, true},
		{"unif", []float64{1, 0}, false},
		{"unif", []float64{0, ninf}, false},
		{"unif", []float64{inf, 5}, false},
		{"unif", []float64{inf, inf}, false},
		{"unif", []float64{ninf, ninf}, false},
		{"unif", []float64{ninf, 5}, true},
		{"unif", []float64{2, 2}, true},
		{"exp", []float64{-1}, false},
		{"exp", []float64{ninf}, false},
		{"exp", []float64{0}, true},
		{"exp", []float64{inf}, true},
		{"gamma", []float64{-1, 1}, false},
		{"gamma", []float64{1, 0}, false},
		{"gamma", []float64{1, -1}, false},
		{"gamma", []float64{0, -1}, false},
		{"gamma", []float64{0, 1}, true},
		{"beta", []float64{-1, 1}, false},
		{"beta", []float64{1, -1}, false},
		{"beta", []float64{0, -1}, false},
		{"beta", []float64{-1, 0}, false},
		{"beta", []float64{0, 0}, true},
		{"t", []float64{0}, false},
		{"t", []float64{-1}, false},
		{"t", []float64{nan}, false},
		{"t", []float64{inf}, true},
		{"f", []float64{0, 1}, false},
		{"f", []float64{1, 0}, false},
		{"f", []float64{-1, 2}, false},
		{"f", []float64{2, -1}, false},
		{"f", []float64{inf, 0}, false},
		{"f", []float64{0, inf}, false},
		{"f", []float64{inf, nan}, false},
		{"f", []float64{inf, inf}, true},
		{"chisq", []float64{-1, 0}, false},
		{"chisq", []float64{1, -1}, false},
		{"chisq", []float64{inf, 0}, true},
		{"chisq", []float64{inf, 1}, false},
		{"chisq", []float64{1, inf}, false},
		{"chisq", []float64{0, 0}, true},
		{"lnorm", []float64{0, -1}, false},
		{"lnorm", []float64{0, 0}, true},
		{"logis", []float64{0, -1}, false},
		{"logis", []float64{0, 0}, true},
		{"cauchy", []float64{0, -1}, false},
		{"cauchy", []float64{0, 0}, true},
		{"weibull", []float64{0, 1}, false},
		{"weibull", []float64{1, 0}, false},
		{"weibull", []float64{-1, 1}, false},
		{"weibull", []float64{1, -1}, false},
		{"weibull", []float64{inf, 0}, false},
		{"weibull", []float64{0, inf}, false},
		{"weibull", []float64{inf, 2}, true},
	}
	for _, c := range cases {
		d, ok := fam(t, c.fam).build(c.p)
		if ok != c.ok {
			t.Errorf("build %s %v: ok = %v, want %v", c.fam, c.p, ok, c.ok)
		}
		if ok && d == nil {
			t.Errorf("build %s %v: ok with a nil law", c.fam, c.p)
		}
	}
	// the limit laws chosen
	for _, c := range []struct {
		fam  string
		p    []float64
		want any
	}{
		{"t", []float64{inf}, locScale{}},
		{"f", []float64{inf, inf}, fMass{}},
		{"f", []float64{inf, 3}, invChi{}},
	} {
		d, _ := fam(t, c.fam).build(c.p)
		switch c.want.(type) {
		case locScale:
			if _, isN := d.(stats.Normal); !isN {
				t.Errorf("build %s %v: %T, want the standard normal", c.fam, c.p, d)
			}
		case fMass:
			if _, is := d.(fMass); !is {
				t.Errorf("build %s %v: %T", c.fam, c.p, d)
			}
		case invChi:
			if _, is := d.(invChi); !is {
				t.Errorf("build %s %v: %T", c.fam, c.p, d)
			}
		}
	}
}

func TestSlice26Constructors(t *testing.T) {
	if _, ok := lsOf(stats.NewNormal(0, -1))(0, 1); ok {
		t.Error("lsOf accepted a standard law that failed to build")
	}
	if _, ok := lsOf(stats.NewNormal(0, 1))(0, -1); ok {
		t.Error("lsOf accepted a negative scale")
	}
	if _, ok := lsOf(stats.NewNormal(0, 1))(0, 0); !ok {
		t.Error("lsOf refused a zero scale")
	}
	if _, ok := mustNew(stats.NewNormal(0, -1)); ok {
		t.Error("mustNew passed an error")
	}
	if d, ok := mustNew(stats.NewNormal(0, 1)); !ok || d == nil {
		t.Error("mustNew refused a good law")
	}
}

func TestSlice26ValueEnds(t *testing.T) {
	norm, unif, gam, tt := fam(t, "norm"), fam(t, "unif"), fam(t, "gamma"), fam(t, "t")
	for _, c := range []struct {
		name string
		got  float64
		want float64
	}{
		// invalid parameters: NaN, but the ends of the support at p = 0 and 1
		{"pnorm sd<0", norm.value(kindP, 0, []float64{0, -1}, true, false), nan},
		{"dnorm sd<0", norm.value(kindD, 0, []float64{0, -1}, true, false), nan},
		{"qnorm(0) sd<0", norm.value(kindQ, 0, []float64{0, -1}, true, false), ninf},
		{"qnorm(1) sd<0", norm.value(kindQ, 1, []float64{0, -1}, true, false), inf},
		{"qnorm(.5) sd<0", norm.value(kindQ, 0.5, []float64{0, -1}, true, false), nan},
		{"qunif(0) invalid", unif.value(kindQ, 0, []float64{1, 0}, true, false), nan},
		{"qgamma(0) invalid", gam.value(kindQ, 0, []float64{-1, 1}, true, false), 0},
		{"qt(1) invalid", tt.value(kindQ, 1, []float64{-1}, true, false), inf},
		// the log density and log tails outside the support are -Inf, not NaN
		{"dunif(-1, log)", unif.value(kindD, -1, []float64{0, 1}, true, true), ninf},
		{"dunif(2, log)", unif.value(kindD, 2, []float64{0, 1}, true, true), ninf},
		{"dunif(.5, log) is -0", unif.value(kindD, 0.5, []float64{0, 1}, true, true), nz},
		{"dgamma(0, 1, 1, log) is -0", gam.value(kindD, 0, []float64{1, 1}, true, true), nz},
		{"pgamma(0, log) at the end", gam.value(kindP, 0, []float64{2, 1}, true, true), ninf},
		{"pgamma(Inf, upper, log) at the end", gam.value(kindP, inf, []float64{2, 1}, false, true), ninf},
		// go-statistics underflows these log tails inside the support: NaN
		{"pgamma(3, 1000, 1/2, log.p)", gam.value(kindP, 3, []float64{1000, 0.5}, true, true), nan},
		{"pt(-1e3, 1000, log.p)", tt.value(kindP, -1e3, []float64{1000}, true, true), nan},
	} {
		if !same(c.got, c.want) {
			t.Errorf("%s = %v, want %v", c.name, c.got, c.want)
		}
	}
}

func TestSlice26Quantile(t *testing.T) {
	e, _ := fam(t, "exp").build([]float64{1})
	p := 0.3
	for _, c := range []struct {
		name         string
		p            float64
		lower, logp  bool
		want         float64
		tolUlps      float64
		wantExactly  bool
		wantExactVal float64
	}{
		{"lower", p, true, false, -math.Log1p(-p), 2, false, 0},
		{"upper", p, false, false, -math.Log(p), 2, false, 0},
		{"lower log", math.Log(p), true, true, -math.Log1p(-p), 4, false, 0},
		{"upper log", math.Log(p), false, true, -math.Log(p), 4, false, 0},
		{"lower log, large p", math.Log(0.8), true, true, -math.Log1p(-0.8), 4, false, 0},
		{"upper log, large p", math.Log(0.8), false, true, -math.Log(0.8), 4, false, 0},
	} {
		got := quantile(e, c.p, c.lower, c.logp)
		if ulpErr(got, c.want) > c.tolUlps {
			t.Errorf("exp quantile %s = %v, want %v", c.name, got, c.want)
		}
	}
	for _, c := range []struct {
		name        string
		p           float64
		lower, logp bool
	}{
		{"p > 1", 2, true, false},
		{"p < 0", -0.5, true, false},
		{"log p > 0", 0.5, true, true},
		{"log p NaN? no: log p > 0", 1e-9, false, true},
	} {
		if v := quantile(e, c.p, c.lower, c.logp); !math.IsNaN(v) {
			t.Errorf("quantile %s = %v, want NaN", c.name, v)
		}
	}
	// the ends
	lo, hi := e.Support()
	for _, c := range []struct {
		name        string
		p           float64
		lower, logp bool
		want        float64
	}{
		{"0 lower", 0, true, false, lo},
		{"1 lower", 1, true, false, hi},
		{"0 upper", 0, false, false, hi},
		{"1 upper", 1, false, false, lo},
		{"-Inf lower log", ninf, true, true, lo},
		{"0 lower log", 0, true, true, hi},
		{"-Inf upper log", ninf, false, true, hi},
		{"0 upper log", 0, false, true, lo},
	} {
		if v := quantile(e, c.p, c.lower, c.logp); !same(v, c.want) {
			t.Errorf("quantile %s = %v, want %v", c.name, v, c.want)
		}
	}
}

// deep log-probability quantiles: the answers R gives where exp(p) underflows.
func TestSlice26DeepQuantile(t *testing.T) {
	cauchy, norm := fam(t, "cauchy"), fam(t, "norm")
	for _, c := range []struct {
		name string
		got  float64
		want float64
		ulps float64
	}{
		// the root is beyond the largest finite value: infinite
		{"qcauchy(-711, log) lower", cauchy.value(kindQ, -711, []float64{0, 1}, true, true), ninf, 0},
		{"qcauchy(-711, log) upper", cauchy.value(kindQ, -711, []float64{0, 1}, false, true), inf, 0},
		// an infinite scale: loc + Inf * z
		{"qnorm(-800, 0, Inf, log) lower", norm.value(kindQ, -800, []float64{0, inf}, true, true), ninf, 0},
		{"qnorm(-800, 0, Inf, log) upper", norm.value(kindQ, -800, []float64{0, inf}, false, true), inf, 0},
		// a finite scale: the standard quantile scaled
		{"qnorm(-1000, 0, 2, log)", norm.value(kindQ, -1000, []float64{0, 2}, true, true), -2 * 44.7, 0.01},
	} {
		if c.ulps == 0 {
			if !same(c.got, c.want) {
				t.Errorf("%s = %v, want %v", c.name, c.got, c.want)
			}
			continue
		}
		if math.Abs(c.got/c.want-1) > c.ulps {
			t.Errorf("%s = %v, want about %v", c.name, c.got, c.want)
		}
	}
}

func TestSlice26DeepQuantileEnds(t *testing.T) {
	// a finite end of the support: the root rounds onto it
	wb, bt, gm := fam(t, "weibull"), fam(t, "beta"), fam(t, "gamma")
	for _, c := range []struct {
		name string
		got  float64
		want float64
	}{
		{"qweibull(-1000, 1, 1, log.p)", wb.value(kindQ, -1000, []float64{1, 1}, true, true), 0},
		{"qbeta(-1000, 1, 1, upper, log.p)", bt.value(kindQ, -1000, []float64{1, 1}, false, true), 1},
		{"qgamma(-1000, 1, 1, log.p)", gm.value(kindQ, -1000, []float64{1, 1}, true, true), 0},
		// the log probability go-statistics cannot follow: NaN, never a number
		{"qgamma(-1000, 5, 1, upper, log.p)", gm.value(kindQ, -1000, []float64{5, 1}, false, true), nan},
	} {
		if !same(c.got, c.want) {
			t.Errorf("%s = %v, want %v", c.name, c.got, c.want)
		}
	}
}

func TestSlice26NegZero(t *testing.T) {
	for _, c := range []struct {
		v      float64
		inside bool
		want   float64
	}{
		{0, true, nz},
		{0, false, 0},
		{nz, false, 0},
		{nz, true, nz},
		{-3, true, -3},
		{-3, false, -3},
	} {
		if got := negZero(c.v, c.inside); !same(got, c.want) {
			t.Errorf("negZero(%v, %v) = %v, want %v", c.v, c.inside, got, c.want)
		}
	}
}

func TestSlice26ExpOfAndInvChi(t *testing.T) {
	n, _ := stats.NewNormal(0, 1)
	e := expOf{locScale{n, 0, 1}}
	chi, _ := stats.NewChiSquared(3, 0, 1)
	v := invChi{chi, 3}
	for _, c := range []struct {
		name string
		d    dist
	}{{"expOf", e}, {"invChi", v}} {
		for _, x := range []float64{-1, 0} {
			for _, k := range []struct {
				name string
				got  float64
				want float64
			}{
				{"PDF", c.d.PDF(x), 0},
				{"LogPDF", c.d.LogPDF(x), ninf},
				{"CDF", c.d.CDF(x), 0},
				{"SF", c.d.SF(x), 1},
				{"LogCDF", c.d.LogCDF(x), ninf},
				{"LogSF", c.d.LogSF(x), 0},
			} {
				if !same(k.got+0, k.want+0) && !(k.got == k.want) {
					t.Errorf("%s.%s(%v) = %v, want %v", c.name, k.name, x, k.got, k.want)
				}
			}
		}
	}
	// the density at +Inf of a log-normal with an infinite sdlog: NaN, on both scales
	w := expOf{locScale{n, 0, inf}}
	if !math.IsNaN(w.PDF(inf)) || !math.IsNaN(w.LogPDF(inf)) {
		t.Errorf("dlnorm(Inf, 0, Inf) = %v, log %v, want NaN", w.PDF(inf), w.LogPDF(inf))
	}
	if math.IsNaN(e.PDF(1)) || math.IsNaN(e.LogPDF(1)) {
		t.Error("dlnorm(1) is NaN")
	}
}

func TestSlice26LocScaleEdges(t *testing.T) {
	n, _ := stats.NewNormal(0, 1)
	// an infinite location met by x at the same infinity: NaN
	l := locScale{n, inf, 1}
	if !math.IsNaN(l.PDF(inf)) || !math.IsNaN(l.LogPDF(inf)) {
		t.Errorf("dnorm(Inf, Inf) = %v, log %v", l.PDF(inf), l.LogPDF(inf))
	}
	// with an infinite scale as well the density is 0, not NaN
	w := locScale{n, inf, inf}
	if w.PDF(inf) != 0 || !math.IsInf(w.LogPDF(inf), -1) {
		t.Errorf("dnorm(Inf, Inf, Inf) = %v, log %v", w.PDF(inf), w.LogPDF(inf))
	}
	// a point at infinity of a finite law: density 0
	f := locScale{n, 0, 1}
	if f.PDF(inf) != 0 || !math.IsInf(f.LogPDF(inf), -1) || f.PDF(ninf) != 0 {
		t.Errorf("dnorm(Inf) = %v, log %v", f.PDF(inf), f.LogPDF(inf))
	}
	if f.CDF(inf) != 1 || f.CDF(ninf) != 0 || f.SF(inf) != 0 || f.SF(ninf) != 1 {
		t.Error("tails of the normal at infinity")
	}
	if f.LogCDF(inf) != 0 || !math.IsInf(f.LogCDF(ninf), -1) || !math.IsInf(f.LogSF(inf), -1) || f.LogSF(ninf) != 0 {
		t.Error("log tails of the normal at infinity")
	}
	// an infinite scale saturates by the side of the location
	g := locScale{n, 0, inf}
	if g.CDF(inf) != 1 || g.CDF(ninf) != 0 {
		t.Errorf("pnorm(+-Inf, 0, Inf) = %v, %v", g.CDF(inf), g.CDF(ninf))
	}
	if g.PDF(3) != 0 || !math.IsInf(g.LogPDF(3), -1) {
		t.Errorf("dnorm(3, 0, Inf) = %v, log %v", g.PDF(3), g.LogPDF(3))
	}
	if !g.interior(0) || g.interior(inf) {
		t.Error("interior of an infinite scale")
	}
}

func TestSlice26InfScale(t *testing.T) {
	for _, c := range []struct {
		name string
		d    infScale
		x    float64
		pdf  float64
	}{
		{"exp at Inf", infScale{nanInf: true}, inf, nan},
		{"exp at 3", infScale{nanInf: true}, 3, 0},
		{"gamma at 0, shape < 1", infScale{atZero: true}, 0, inf},
		{"gamma at 0, shape >= 1", infScale{}, 0, 0},
		{"weibull at 2, shape < 1", infScale{atZero: true, nanPos: true}, 2, nan},
		{"weibull at Inf, shape < 1", infScale{atZero: true, nanPos: true}, inf, 0},
		{"weibull at 2, shape >= 1", infScale{}, 2, 0},
	} {
		if got := c.d.PDF(c.x); !same(got, c.pdf) {
			t.Errorf("%s: PDF = %v, want %v", c.name, got, c.pdf)
		}
		want := c.pdf
		switch {
		case c.pdf == 0:
			want = ninf
		case math.IsInf(c.pdf, 1):
			want = inf
		}
		if got := c.d.LogPDF(c.x); !same(got, want) {
			t.Errorf("%s: LogPDF = %v, want %v", c.name, got, want)
		}
	}
	d := infScale{}
	if !math.IsNaN(d.CDF(inf)) || d.CDF(3) != 0 || !math.IsNaN(d.SF(inf)) || d.SF(3) != 1 {
		t.Error("infScale tails")
	}
	if !math.IsNaN(d.LogCDF(inf)) || !math.IsInf(d.LogCDF(3), -1) || !math.IsNaN(d.LogSF(inf)) || d.LogSF(3) != 0 {
		t.Error("infScale log tails")
	}
}

func TestSlice26PointLaws(t *testing.T) {
	z := zeroPlus{pointMass{0}}
	if z.CDF(0) != 0 || z.CDF(1) != 1 || z.SF(0) != 1 || z.SF(1) != 0 {
		t.Error("zeroPlus tails")
	}
	if !math.IsInf(z.LogCDF(0), -1) || z.LogCDF(1) != 0 || z.LogSF(0) != 0 || !math.IsInf(z.LogSF(1), -1) {
		t.Error("zeroPlus log tails")
	}
	pn := pointNaN{zeroPlus{pointMass{0}}, true}
	if !pn.interior(1) || pn.interior(0) || pn.interior(-1) {
		t.Error("pointNaN interior")
	}
	if pu := (pointNaN{pointMass{2}, false}); pu.interior(1) {
		t.Error("pointNaN without neg has an interior")
	}
	if _, hi := pn.Support(); !math.IsNaN(hi) {
		t.Errorf("pointNaN(neg) upper end = %v, want NaN", hi)
	}
	if _, hi := (pointNaN{pointMass{2}, false}).Support(); math.IsNaN(hi) {
		t.Error("pointNaN upper end is NaN without neg")
	}
	nm := negMass{pointMass{inf}}
	if !nm.interior(1) || nm.interior(0) || nm.PDF(3) != 0 || !math.IsInf(nm.LogPDF(3), -1) {
		t.Error("negMass")
	}
	b := betaMass{both: true}
	for _, c := range []struct{ x, want float64 }{{-1, 0.5}, {0, 0.5}, {0.5, 0.5}, {1, 1}, {2, 1}} {
		if got := b.CDF(c.x); got != c.want {
			t.Errorf("beta(0,0) CDF(%v) = %v, want %v", c.x, got, c.want)
		}
	}
	if b.PPF(0.2) != 0 || b.PPF(0.5) != 0.5 || b.PPF(0.8) != 1 || b.ISF(0.2) != 0 || b.ISF(0.8) != 1 {
		t.Error("beta(0,0) quantiles")
	}
	if b.PDF(0) != inf || b.PDF(1) != inf || b.PDF(0.5) != 0 {
		t.Error("beta(0,0) density")
	}
	m0, m1, mh := betaMass{m: 0}, betaMass{m: 1}, betaMass{m: 0.5}
	if m0.CDF(0) != 1 || m1.CDF(0.5) != 0 || m1.CDF(1) != 1 || mh.CDF(0.4) != 0 || mh.CDF(0.5) != 1 {
		t.Error("beta point mass CDF")
	}
	if m0.PDF(0) != inf || m0.PDF(0.5) != 0 || m1.PDF(1) != inf || mh.PDF(0.5) != inf || m0.PPF(0.3) != 0 {
		t.Error("beta point mass PDF, PPF")
	}
	fm := fMass{pointMass{1}}
	if fm.CDF(0.5) != 0 || fm.CDF(1) != 0.5 || fm.CDF(2) != 1 || fm.SF(0.5) != 1 || fm.SF(1) != 0.5 || fm.SF(2) != 0 {
		t.Error("fMass tails")
	}
	w := weibMass{2}
	if w.PDF(-1) != 0 || w.PDF(inf) != 0 || w.PDF(2) != inf || !math.IsNaN(w.PDF(1)) || !math.IsNaN(w.PDF(3)) {
		t.Error("weibMass density")
	}
	if w.CDF(1) != 0 || w.CDF(3) != 1 || w.SF(1) != 1 || w.SF(3) != 0 {
		t.Error("weibMass tails")
	}
	if got := w.CDF(2); math.Abs(got-(-math.Expm1(-1))) > 1e-16 {
		t.Errorf("weibMass CDF at the scale = %v", got)
	}
	if w.PPF(0.3) != 2 || w.ISF(0.3) != 2 || !w.interior(1) || w.interior(0) {
		t.Error("weibMass quantile and interior")
	}
}

func TestSlice26PartialNames(t *testing.T) {
	formals := []string{"q", "mean", "sd", "lower.tail", "log.p"}
	for _, c := range []struct {
		in, want []string
	}{
		{[]string{"", "lower", "log"}, []string{"", "lower.tail", "log.p"}},
		{[]string{"", "l"}, []string{"", "l"}},
		{[]string{"", "mean", "me"}, []string{"", "mean", "me"}},
		{[]string{"", "m", "s"}, []string{"", "mean", "sd"}},
		{[]string{"", "zz"}, []string{"", "zz"}},
		{[]string{"", "lower.tail", "lower"}, []string{"", "lower.tail", "lower"}},
		{[]string{"", "log.p"}, []string{"", "log.p"}},
	} {
		got := partialNames(&callInfo{names: c.in}, formals).names
		if !slices.Equal(got, c.want) {
			t.Errorf("partialNames(%v) = %v, want %v", c.in, got, c.want)
		}
	}
	orig := []string{"", "lower"}
	partialNames(&callInfo{names: orig}, formals)
	if orig[1] != "lower" {
		t.Error("partialNames modified the caller's names")
	}
}
