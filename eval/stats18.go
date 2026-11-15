package eval

import (
	"math"
	"sort"
	"strings"
)

// The descriptive statistics of slice 18: median, var, sd, cov, cor,
// quantile (types 1 to 9), fivenum, IQR, mad, weighted.mean. The kernels are
// written here, from the help pages (the definitions of the nine quantile
// types are the ones of Hyndman and Fan that ?quantile cites) and from the
// oracle; go-numerics has none of them yet.
//
// u = 2^-53 is the unit roundoff. The stated bounds, asserted in
// stats18_test.go against math/big:
//
//	varKernel      |v - var| <= 4u var
//	covKernel      |c - cov| <= 4u S/(n-1) with S = sum |dx_i dy_i| (the
//	               componentwise bound), and |r - cor| <= 8u in absolute value
//	quantile       types 1 to 3: the exact order statistic (or the mean of two,
//	               within 2u of it), where the index is decided with the
//	               tolerance 4*2^-52 on n p + m; types 4 to 9: <= 4u max(|x_j|,
//	               |x_{j+1}|) from the interpolation at the exact position
//	median         exact for an odd count, <= 2u for an even one
//	fivenum        <= 2u (the average of two order statistics)
//	weighted.mean  <= 4u S_abs/|sum w| with S_abs = sum |x w| (compensated sums)
const u53 = 0x1p-53

// ---- compensated sums -----------------------------------------------------

// neumaier adds the elements of x with a running compensation; s + c is the
// sum with an error of at most n u^2 sum|x| (second order).
func neumaier(x []float64) (s, c float64) {
	for _, v := range x {
		t := s + v
		if math.Abs(s) >= math.Abs(v) {
			c += (s - t) + v
		} else {
			c += (v - t) + s
		}
		s = t
	}
	return
}

func finite(x float64) bool { return !math.IsNaN(x) && !math.IsInf(x, 0) }

// varKernel is the sample variance of x (no missing values): the squares of the
// deviations from the mean are added exactly (fused products, compensated
// sum) and the rounding of the mean is corrected. NaN when a value is not
// finite; NA (as NaN) when there are fewer than two values.
func varKernel(x []float64) float64 {
	n := len(x)
	if n < 2 {
		return naReal
	}
	ss, ok := sumSquares(x)
	if !ok {
		return math.NaN()
	}
	return ss / float64(n-1)
}

func sumSquares(x []float64) (float64, bool) {
	n := float64(len(x))
	s, c := neumaier(x)
	if !finite(s) || !finite(c) {
		return 0, false
	}
	mean := (s + c) / n
	// the residual r = sum(x) - n mean, so that sum (x - mean) = r
	p := n * mean
	e := math.FMA(n, mean, -p)
	r := ((s - p) + c) - e
	var a, comp float64
	for _, v := range x {
		d := v - mean
		sq := d * d
		comp += math.FMA(d, d, -sq)
		t := a + sq
		if a >= sq {
			comp += (a - t) + sq
		} else {
			comp += (sq - t) + a
		}
		a = t
	}
	return (a + comp) - r*r/n, true
}

// pairKernel answers the sums of squares and of cross products of two columns
// without missing values: sxx, syy and sxy about their means.
func pairKernel(x, y []float64) (sxx, syy, sxy float64, ok bool) {
	n := float64(len(x))
	sx, cx := neumaier(x)
	sy, cy := neumaier(y)
	if !finite(sx) || !finite(cx) || !finite(sy) || !finite(cy) {
		return 0, 0, 0, false
	}
	mx, my := (sx+cx)/n, (sy+cy)/n
	px := n * mx
	ex := math.FMA(n, mx, -px)
	rx := ((sx - px) + cx) - ex
	py := n * my
	ey := math.FMA(n, my, -py)
	ry := ((sy - py) + cy) - ey
	var axx, cxx, ayy, cyy, axy, cxy float64
	add := func(a, comp *float64, v float64) {
		t := *a + v
		if math.Abs(*a) >= math.Abs(v) {
			*comp += (*a - t) + v
		} else {
			*comp += (v - t) + *a
		}
		*a = t
	}
	for i := range x {
		dx, dy := x[i]-mx, y[i]-my
		p := dx * dx
		cxx += math.FMA(dx, dx, -p)
		add(&axx, &cxx, p)
		q := dy * dy
		cyy += math.FMA(dy, dy, -q)
		add(&ayy, &cyy, q)
		m := dx * dy
		cxy += math.FMA(dx, dy, -m)
		add(&axy, &cxy, m)
	}
	return (axx + cxx) - rx*rx/n, (ayy + cyy) - ry*ry/n, (axy + cxy) - rx*ry/n, true
}

// ---- ranks and Kendall ----------------------------------------------------

// ranks are the ranks of x, ties getting their average.
func ranks(x []float64) []float64 {
	n := len(x)
	idx := make([]int, n)
	for i := range idx {
		idx[i] = i
	}
	sort.SliceStable(idx, func(a, b int) bool { return x[idx[a]] < x[idx[b]] })
	r := make([]float64, n)
	for i := 0; i < n; {
		j := i
		for j+1 < n && x[idx[j+1]] == x[idx[i]] {
			j++
		}
		avg := float64(i+j)/2 + 1
		for k := i; k <= j; k++ {
			r[idx[k]] = avg
		}
		i = j + 1
	}
	return r
}

func sign(v float64) float64 {
	switch {
	case v > 0:
		return 1
	case v < 0:
		return -1
	}
	return 0
}

// kendallSums answers sum over ordered pairs of sign(dx) sign(dy), and the
// same for x with itself and y with itself.
func kendallSums(x, y []float64) (sxy, sxx, syy float64) {
	n := len(x)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i == j {
				continue
			}
			a, b := sign(x[i]-x[j]), sign(y[i]-y[j])
			sxy += a * b
			sxx += a * a
			syy += b * b
		}
	}
	return
}

// ---- cov and cor ----------------------------------------------------------

// statColumns are the columns of a numeric vector, matrix or data frame, the
// missing values as NaN, and the column names.
func (s *Session) statColumns(c *callInfo, v Value, what string, nullMsg string) (cols [][]float64, names *String, matrix bool) {
	if _, isNull := v.(nullT); isNull {
		s.errorf(c.call, "%s", nullMsg)
	}
	conv := func(e Value) []float64 {
		switch x := e.(type) {
		case *Double:
			return x.V
		case *Integer, *Logical:
			if inheritsFrom(e, "factor") {
				s.errorf(c.call, "'%s' must be numeric", what)
			}
			iv := asInts(e)
			out := make([]float64, len(iv))
			for i, k := range iv {
				if k == naInt {
					out[i] = math.NaN()
				} else {
					out[i] = float64(k)
				}
			}
			return out
		}
		s.errorf(c.call, "'%s' must be numeric", what)
		return nil
	}
	if l, ok := v.(*List); ok {
		if !inheritsFrom(v, "data.frame") {
			s.errorf(c.call, "'%s' must be numeric", what)
		}
		for _, e := range l.V {
			cols = append(cols, conv(e))
		}
		if len(cols) > 1 {
			for _, cl := range cols {
				if len(cl) != len(cols[0]) {
					s.errorf(c.call, "incompatible dimensions")
				}
			}
		}
		return cols, namesOf(v), true
	}
	all := conv(v)
	if d, ok := getAttr(v, "dim").(*Integer); ok && len(d.V) == 2 {
		nr, nc := int(d.V[0]), int(d.V[1])
		for j := 0; j < nc; j++ {
			cols = append(cols, all[j*nr:(j+1)*nr])
		}
		if dn, ok := getAttr(v, "dimnames").(*List); ok && len(dn.V) == 2 {
			if cn, ok := dn.V[1].(*String); ok {
				names = cn
			}
		}
		return cols, names, true
	}
	return [][]float64{all}, nil, false
}

var useModes = []string{"all.obs", "complete.obs", "pairwise.complete.obs", "everything", "na.or.complete"}

func (s *Session) matchUse(c *callInfo, v Value) string {
	st, ok := v.(*String)
	if !ok || len(st.V) != 1 {
		s.errorf(c.call, "invalid 'use' argument")
	}
	found := ""
	for _, m := range useModes {
		if st.V[0].S != "" && strings.HasPrefix(m, st.V[0].S) {
			if found != "" {
				s.errorf(c.call, "invalid 'use' argument")
			}
			found = m
		}
	}
	if found == "" {
		s.errorf(c.call, "invalid 'use' argument")
	}
	return found
}

func (s *Session) matchMethod(c *callInfo, v Value) string {
	if v == nil {
		return "pearson"
	}
	st, ok := v.(*String)
	opts := []string{"pearson", "kendall", "spearman"}
	if ok && len(st.V) > 0 {
		for _, m := range opts {
			if st.V[0].S != "" && strings.HasPrefix(m, st.V[0].S) {
				return m
			}
		}
	}
	s.errorf(&Lang{Fun: Intern("match.arg"), Args: []Arg{{Value: Intern("method")}}},
		"'arg' should be one of “pearson”, “kendall”, “spearman”")
	return ""
}

// covCor is var, cov and cor. cor is true for the correlation.
func (s *Session) covCor(c *callInfo, x, y Value, use, method string, cor bool) Value {
	xn := "'x' is NULL"
	xc, xnames, xmat := s.statColumns(c, x, "x", xn)
	var yc [][]float64
	var ynames *String
	ymat := false
	hasY := false
	if y != nil {
		if _, isNull := y.(nullT); !isNull {
			hasY = true
			yc, ynames, ymat = s.statColumns(c, y, "y", "'y' is NULL")
		}
	}
	if cor && !hasY && !xmat {
		s.errorf(c.call, "supply both 'x' and 'y' or a matrix-like 'x'")
	}
	n := 0
	if len(xc) > 0 {
		n = len(xc[0])
	}
	if hasY {
		ny := 0
		if len(yc) > 0 {
			ny = len(yc[0])
		}
		if ny != n {
			s.errorf(c.call, "incompatible dimensions")
		}
	}
	if !hasY {
		yc, ynames = xc, xnames
	}
	// rows to keep
	hasNA := func(cols [][]float64, i int) bool {
		for _, cl := range cols {
			if cl[i] != cl[i] {
				return true
			}
		}
		return false
	}
	keep := make([]int, 0, n)
	switch use {
	case "all.obs":
		for i := 0; i < n; i++ {
			if hasNA(xc, i) || (hasY && hasNA(yc, i)) {
				s.errorf(c.call, "missing observations in cov/cor")
			}
		}
	case "complete.obs", "na.or.complete":
		for i := 0; i < n; i++ {
			if !hasNA(xc, i) && !(hasY && hasNA(yc, i)) {
				keep = append(keep, i)
			}
		}
		if len(keep) == 0 && use == "complete.obs" {
			s.errorf(c.call, "no complete element pairs")
		}
	}
	sub := func(cols [][]float64) [][]float64 {
		out := make([][]float64, len(cols))
		for j, cl := range cols {
			o := make([]float64, len(keep))
			for k, i := range keep {
				o[k] = cl[i]
			}
			out[j] = o
		}
		return out
	}
	complete := use == "complete.obs" || use == "na.or.complete"
	if complete {
		xc = sub(xc)
		if hasY {
			yc = sub(yc)
		} else {
			yc = xc
		}
	}
	if method == "spearman" && (complete || use == "all.obs") {
		rk := func(cols [][]float64) [][]float64 {
			out := make([][]float64, len(cols))
			for j, cl := range cols {
				out[j] = ranks(cl)
			}
			return out
		}
		xc = rk(xc)
		if hasY {
			yc = rk(yc)
		} else {
			yc = xc
		}
		method = "pearson"
	}
	warned := false
	cell := func(a, b []float64, same bool) float64 {
		// the missing values of this pair, by the mode
		if use == "pairwise.complete.obs" || use == "everything" || use == "all.obs" {
			var aa, bb []float64
			skip := false
			for i := range a {
				if a[i] != a[i] || b[i] != b[i] {
					if use == "everything" {
						return naReal
					}
					skip = true
					continue
				}
				aa, bb = append(aa, a[i]), append(bb, b[i])
			}
			_ = skip
			a, b = aa, bb
		}
		if len(a) < 2 {
			return naReal
		}
		if method == "spearman" {
			a, b = ranks(a), ranks(b)
		}
		switch method {
		case "kendall":
			sxy, sxx, syy := kendallSums(a, b)
			if !cor {
				return sxy
			}
			if same {
				return 1
			}
			if sxx == 0 || syy == 0 {
				warned = true
				return naReal
			}
			return clip1(sxy / math.Sqrt(sxx*syy))
		}
		sxx, syy, sxy, ok := pairKernel(a, b)
		if !ok {
			return math.NaN()
		}
		if !cor {
			return sxy / float64(len(a)-1)
		}
		if same {
			return 1
		}
		if sxx == 0 || syy == 0 {
			warned = true
			return naReal
		}
		return clip1(sxy / math.Sqrt(sxx) / math.Sqrt(syy))
	}
	nr, nc := len(xc), len(yc)
	if len(xc) > 0 && len(yc) > 0 && len(xc[0]) == 0 && use == "na.or.complete" {
		out := make([]float64, nr*nc)
		for i := range out {
			out[i] = naReal
		}
		return s.statResult(out, nr, nc, xnames, ynames, xmat || ymat, hasY)
	}
	out := make([]float64, nr*nc)
	if nr*nc > s.opts.MaxVector {
		s.errorf(c.call, "cannot allocate vector of length %d", nr*nc)
	}
	for j := 0; j < nc; j++ {
		for i := 0; i < nr; i++ {
			if !hasY && i > j {
				break
			}
			v := cell(xc[i], yc[j], !hasY && i == j)
			out[j*nr+i] = v
			if !hasY {
				out[i*nr+j] = v
			}
		}
	}
	if warned {
		s.warningf(c.call, "the standard deviation is zero")
	}
	return s.statResult(out, nr, nc, xnames, ynames, xmat || ymat, hasY)
}

func clip1(r float64) float64 {
	switch {
	case r > 1:
		return 1
	case r < -1:
		return -1
	}
	return r
}

func (s *Session) statResult(out []float64, nr, nc int, xn, yn *String, matrix, hasY bool) Value {
	if !matrix {
		return &Double{V: out[:1]}
	}
	res := &Double{V: out, Attr: []Attr{{Name: "dim", Value: &Integer{V: []int32{int32(nr), int32(nc)}}}}}
	if xn != nil || yn != nil {
		var rn, cn Value = Null, Null
		if xn != nil {
			rn = &String{V: xn.V}
		}
		if yn != nil {
			cn = &String{V: yn.V}
		}
		res.Attr = append(res.Attr, Attr{Name: "dimnames", Value: &List{V: []Value{rn, cn}}})
	}
	return res
}

// ---- order statistics -----------------------------------------------------

// selectK partially sorts a so that a[k] is the k-th smallest for every k in
// ks (sorted); a is modified. Introselect: quickselect with median-of-three
// pivots, a full sort when the partitions go bad.
func selectK(a []float64, ks []int) {
	var rec func(lo, hi int, ks []int, depth int)
	rec = func(lo, hi int, ks []int, depth int) {
		for len(ks) > 0 && hi-lo > 16 {
			if depth == 0 {
				sort.Float64s(a[lo : hi+1])
				return
			}
			depth--
			p := partition(a, lo, hi)
			// ks below, at and above the pivot
			i := sort.SearchInts(ks, p)
			left := ks[:i]
			right := ks[i:]
			if len(right) > 0 && right[0] == p {
				right = right[1:]
			}
			if len(left) > 0 && len(right) > 0 {
				rec(lo, p-1, left, depth)
				lo, ks = p+1, right
				continue
			}
			if len(left) > 0 {
				hi, ks = p-1, left
			} else {
				lo, ks = p+1, right
			}
		}
		if len(ks) > 0 {
			sort.Float64s(a[lo : hi+1])
		}
	}
	rec(0, len(a)-1, ks, 2*bitsLen(len(a))+4)
}

func bitsLen(n int) int {
	b := 0
	for ; n > 0; n >>= 1 {
		b++
	}
	return b
}

// partition places a median-of-three pivot at its final index and answers it.
func partition(a []float64, lo, hi int) int {
	mid := lo + (hi-lo)/2
	if a[mid] < a[lo] {
		a[mid], a[lo] = a[lo], a[mid]
	}
	if a[hi] < a[lo] {
		a[hi], a[lo] = a[lo], a[hi]
	}
	if a[hi] < a[mid] {
		a[hi], a[mid] = a[mid], a[hi]
	}
	pv := a[mid]
	a[mid], a[hi-1] = a[hi-1], a[mid]
	i, j := lo, hi-1
	for {
		for i++; a[i] < pv; i++ {
		}
		for j--; a[j] > pv; j-- {
		}
		if i >= j {
			break
		}
		a[i], a[j] = a[j], a[i]
	}
	a[i], a[hi-1] = a[hi-1], a[i]
	return i
}

// half answers the mean of two values without overflow.
func halfSum(a, b float64) float64 {
	if a == b {
		return a
	}
	if math.Abs(a) < 1e300 && math.Abs(b) < 1e300 {
		return (a + b) / 2
	}
	return a/2 + b/2
}

// medianOf is the median of x (no missing values; x is modified).
func medianOf(x []float64) float64 {
	n := len(x)
	if n == 0 {
		return naReal
	}
	h := (n + 1) / 2
	if n%2 == 1 {
		selectK(x, []int{h - 1})
		return x[h-1]
	}
	selectK(x, []int{h - 1, h})
	return halfSum(x[h-1], x[h])
}

// A double-double constant for the positions of the continuous quantile types.
type dd struct{ hi, lo float64 }

// quantileParams are, for a type, the position n p + m as a + p (n + 1 - a - b)
// with a and k = 1 - a - b as double-doubles.
func quantileParams(typ int) (a, k dd) {
	third := dd{0x1.5555555555555p-2, 0x1.5555555555555p-56}
	switch typ {
	case 4:
		return dd{0, 0}, dd{0, 0}
	case 5:
		return dd{0.5, 0}, dd{0, 0}
	case 6:
		return dd{0, 0}, dd{1, 0}
	case 7:
		return dd{1, 0}, dd{-1, 0}
	case 8:
		return third, third
	}
	return dd{0.375, 0}, dd{0.25, 0}
}

// quantileFuzz is the tolerance on the position that decides whether it is an
// integer.
const quantileFuzz = 4 * 0x1p-52

// quantileAt answers the quantile of sorted-by-selection data a (modified) at
// probability p, of the given type. The positions are computed from the exact
// product n p (fused multiply-add) so that the fractional part is exact to
// within u.
func quantileAt(a []float64, p float64, typ int) float64 {
	n := len(a)
	nf := float64(n)
	if typ <= 3 {
		var m float64
		if typ == 3 {
			m = -0.5
		}
		hi := nf * p
		lo := math.FMA(nf, p, -hi)
		np := hi + m
		j := math.Floor(np + quantileFuzz*math.Max(1, math.Abs(np)))
		g := (np - j) + lo
		if math.Abs(g) < quantileFuzz*math.Max(1, math.Abs(np)) {
			g = 0
		}
		var gamma float64
		switch typ {
		case 1:
			if g > 0 {
				gamma = 1
			}
		case 2:
			if g > 0 {
				gamma = 1
			} else {
				gamma = 0.5
			}
		case 3:
			if g != 0 || int64(j)%2 != 0 {
				gamma = 1
			}
		}
		return interp(a, int(j), gamma)
	}
	ak, kk := quantileParams(typ)
	// position = ak + p (n + k), with k a double-double
	kn := dd{nf + kk.hi, kk.lo}
	h1 := p * kn.hi
	l1 := math.FMA(p, kn.hi, -h1)
	l1 += p * kn.lo
	pos := ak.hi + h1
	// error-free addition of ak.hi and h1
	bb := pos - ak.hi
	l2 := (ak.hi - (pos - bb)) + (h1 - bb)
	rest := l1 + l2 + ak.lo
	j := math.Floor(pos + quantileFuzz*math.Max(1, math.Abs(pos)))
	g := (pos - j) + rest
	if math.Abs(g) < quantileFuzz*math.Max(1, math.Abs(pos)) {
		g = 0
	}
	if g < 0 {
		g = 0
	}
	return interp(a, int(j), g)
}

// interp is (1 - g) x[j] + g x[j + 1] with 1-based j, the data clamped at both
// ends; a holds the data with the needed order statistics in place.
func interp(a []float64, j int, g float64) float64 {
	n := len(a)
	lo, hi := j, j+1
	if lo < 1 {
		lo = 1
	}
	if lo > n {
		lo = n
	}
	if hi < 1 {
		hi = 1
	}
	if hi > n {
		hi = n
	}
	xl, xh := a[lo-1], a[hi-1]
	switch {
	case g <= 0 || xl == xh:
		return xl
	case g >= 1:
		return xh
	}
	return (1-g)*xl + g*xh
}

// needed lists the 0-based indices that quantileAt of these probabilities reads.
func quantileIndices(n int, probs []float64, typ int) []int {
	set := map[int]bool{}
	for _, p := range probs {
		if p != p {
			continue
		}
		// a position within 2 of the exact one, clamped: cheap and safe
		var pos float64
		switch typ {
		case 1, 2, 3:
			pos = float64(n) * p
		default:
			pos = float64(n+1)*p + 1
		}
		j := int(math.Floor(pos))
		for d := -2; d <= 3; d++ {
			if k := j + d; k >= 1 && k <= n {
				set[k-1] = true
			}
		}
	}
	set[0], set[n-1] = true, true
	ks := make([]int, 0, len(set))
	for k := range set {
		ks = append(ks, k)
	}
	sort.Ints(ks)
	return ks
}

// ---- the builtins ---------------------------------------------------------

// numericVector answers the elements of x as doubles (missing as NaN) when x
// is a plain numeric or logical vector.
func numericVector(x Value) ([]float64, bool) {
	switch v := x.(type) {
	case *Double:
		return v.V, true
	case *Integer, *Logical:
		iv := asInts(x)
		out := make([]float64, len(iv))
		for i, k := range iv {
			if k == naInt {
				out[i] = math.NaN()
			} else {
				out[i] = float64(k)
			}
		}
		return out, true
	}
	return nil, false
}

func hasNaN(x []float64) bool {
	for _, v := range x {
		if v != v {
			return true
		}
	}
	return false
}

func dropNaN(x []float64) []float64 {
	out := make([]float64, 0, len(x))
	for _, v := range x {
		if v == v {
			out = append(out, v)
		}
	}
	return out
}

func init() {
	regBuiltin("var", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "y", "na.rm", "use")
		if a[0] == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		naRm := s.boolArg(c, a[2], false, "na.rm")
		use := "everything"
		if naRm {
			use = "na.or.complete"
		}
		if a[3] != nil {
			use = s.matchUse(c, a[3])
		}
		if inheritsFrom(a[0], "factor") {
			s.errorf(c.call, "Calling var(x) on a factor x is defunct.\n  Use something like 'all(duplicated(x)[-1L])' to test for a constant vector.")
		}
		return s.covCor(c, a[0], a[1], use, "pearson", false)
	})
	regBuiltin("sd", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "na.rm")
		if a[0] == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		naRm := s.boolArg(c, a[1], false, "na.rm")
		x, ok := numericVector(a[0])
		if !ok || inheritsFrom(a[0], "factor") {
			s.errorf(c.call, "'x' must be numeric")
		}
		if naRm {
			x = dropNaN(x)
		} else if hasNaN(x) {
			return mkDouble(naReal)
		}
		return mkDouble(math.Sqrt(varKernel(x)))
	})
	covcor := func(name string, cor bool) {
		regBuiltin(name, func(s *Session, c *callInfo) Value {
			a := s.matchFormals(c, "x", "y", "use", "method")
			if a[0] == nil {
				s.errorf(c.call, "argument \"x\" is missing, with no default")
			}
			use := "everything"
			if a[2] != nil {
				use = s.matchUse(c, a[2])
			}
			method := s.matchMethod(c, a[3])
			return s.covCor(c, a[0], a[1], use, method, cor)
		})
	}
	covcor("cov", false)
	covcor("cor", true)

	regBuiltin("median.default", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "na.rm", "...")
		if a[0] == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		naRm := s.boolArg(c, a[1], false, "na.rm")
		x := a[0]
		if inheritsFrom(x, "factor") || inheritsFrom(x, "data.frame") {
			s.errorf(c.call, "need numeric data")
		}
		if _, isNull := x.(nullT); isNull {
			return Null
		}
		return s.medianValue(c, x, naRm)
	})
}

// medianValue is median.default of any vector.
func (s *Session) medianValue(c *callInfo, x Value, naRm bool) Value {
	plain := classAttr(x) == nil
	// the vector without names; NA handling by the type
	var kept Value = x
	if getAttr(x, "names") != nil {
		kept = s.stripNames(x)
	}
	isna := s.callNamed(c, "is.na", []Value{kept}, []string{""})
	nas := asInts(isna)
	anyNA := false
	for _, v := range nas {
		if v == 1 {
			anyNA = true
		}
	}
	if anyNA {
		if !naRm {
			return s.callNamed(c, "[", []Value{kept, &Integer{V: []int32{naInt}}}, []string{"", ""})
		}
		notna := make([]int32, len(nas))
		for i, v := range nas {
			if v == 0 {
				notna[i] = 1
			}
		}
		kept = s.callNamed(c, "[", []Value{kept, &Logical{V: notna}}, []string{"", ""})
	}
	n := Length(kept)
	if n == 0 {
		return s.callNamed(c, "[", []Value{kept, &Integer{V: []int32{naInt}}}, []string{"", ""})
	}
	if plain {
		switch v := kept.(type) {
		case *Double:
			cp := append([]float64(nil), v.V...)
			return mkDouble(medianOf(cp))
		case *Integer:
			if n%2 == 1 {
				cp := append([]int32(nil), v.V...)
				sort.Slice(cp, func(i, j int) bool { return cp[i] < cp[j] })
				return &Integer{V: []int32{cp[n/2]}}
			}
			cp := make([]float64, n)
			for i, e := range v.V {
				cp[i] = float64(e)
			}
			return mkDouble(medianOf(cp))
		}
	}
	// any other vector: sort, take the middle one or the mean of the two
	sorted := s.callNamed(c, "sort", []Value{kept}, []string{""})
	h := (n + 1) / 2
	if n%2 == 1 {
		return s.callNamed(c, "[", []Value{sorted, mkInt(int32(h))}, []string{"", ""})
	}
	mid := s.callNamed(c, "[", []Value{sorted, &Integer{V: []int32{int32(h), int32(h + 1)}}}, []string{"", ""})
	return s.callNamed(c, "mean", []Value{mid}, []string{""})
}

func (s *Session) stripNames(x Value) Value {
	var keep []Attr
	for _, at := range attrsOf(x) {
		if at.Name != "names" {
			keep = append(keep, at)
		}
	}
	return withAttr(x, keep)
}
