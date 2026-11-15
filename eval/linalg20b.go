package eval

import (
	"math"
	"math/cmplx"
	"sort"
	"strings"

	"github.com/go-embedded-r/r/syntax"
	"github.com/go-ndarray/ndarray"
	"github.com/go-ndarray/ndarray/linalg"
)

func (s *Session) parseCallText(src string) *Lang {
	f, err := syntax.Parse(src, syntax.Options{})
	if err != nil || len(f.Exprs) != 1 {
		return nil
	}
	l, _ := s.fromSyntax(f.Exprs[0]).(*Lang)
	return l
}

func (s *Session) optBool(c *callInfo, v Value, def bool, what string) bool {
	return s.boolArg(c, v, def, what)
}

func detList(mod float64, sign float64, log bool) Value {
	m := &Double{V: []float64{mod}, Attr: []Attr{{Name: "logarithm", Value: mkLogical(log)}}}
	l := mkNamedList([]string{"modulus", "sign"}, m, mkInt(int32(sign)))
	return setAttr(l, "class", mkString("det"))
}

func sortedDesc(n int, key func(i int) float64) []int {
	ord := make([]int, n)
	for i := range ord {
		ord[i] = i
	}
	sort.SliceStable(ord, func(a, b int) bool { return key(ord[a]) > key(ord[b]) })
	return ord
}

// triSolve solves T x = b in place for the k x k triangle of t (leading
// dimension ld), upper or lower, transposed or not, for each column of b
// (k rows, ncb columns). It returns the 1-based index of the first zero on
// the diagonal, or 0.
func triSolve(t []float64, ld, k int, upper, trans bool, b []float64, ncb int) int {
	for i := 0; i < k; i++ {
		if t[i+i*ld] == 0 {
			return i + 1
		}
	}
	// the matrix actually solved is T or T', so: effective upper-ness
	up := upper != trans
	for col := 0; col < ncb; col++ {
		x := b[col*k : (col+1)*k]
		at := func(i, j int) float64 {
			if trans {
				return t[j+i*ld]
			}
			return t[i+j*ld]
		}
		if up {
			for i := k - 1; i >= 0; i-- {
				r := x[i]
				for j := i + 1; j < k; j++ {
					r -= at(i, j) * x[j]
				}
				x[i] = r / at(i, i)
			}
		} else {
			for i := 0; i < k; i++ {
				r := x[i]
				for j := 0; j < i; j++ {
					r -= at(i, j) * x[j]
				}
				x[i] = r / at(i, i)
			}
		}
	}
	return 0
}

func (s *Session) backsolveFn(c *callInfo, defUpper bool) Value {
	a := s.matchFormals(c, map[bool][]string{true: {"r", "x", "k", "upper.tri", "transpose"}, false: {"l", "x", "k", "upper.tri", "transpose"}}[defUpper]...)
	if a[0] == nil || a[1] == nil {
		nm := "x"
		if a[0] == nil {
			nm = map[bool]string{true: "r", false: "l"}[defUpper]
		}
		s.errorf(c.call, "argument \"%s\" is missing, with no default", nm)
	}
	r := s.lmatrix(c, a[0], "'r' must be a numeric matrix")
	if r.cplx {
		s.errorf(c.call, "invalid argument values in 'backsolve'")
	}
	xm := s.lmatrix(c, a[1], "'x' must be a numeric matrix")
	xIsMat := len(dimsOf(a[1])) == 2
	k := r.nc
	if a[2] != nil {
		d := asDoubles(a[2])
		if len(d) != 1 || math.IsNaN(d[0]) || d[0] < 1 || d[0] > float64(math.MaxInt32) {
			s.errorf(c.call, "invalid 'k' argument")
		}
		k = int(d[0])
	}
	upper := s.optBool(c, a[3], defUpper, "upper.tri")
	trans := s.optBool(c, a[4], false, "transpose")
	if k > r.nr || k > r.nc || k > xm.nr {
		s.errorf(c.call, "invalid 'k' argument")
	}
	if xm.cplx {
		s.errorf(c.call, "invalid argument values in 'backsolve'")
	}
	s.checkMat(c, k, xm.nc)
	b := make([]float64, k*xm.nc)
	for j := 0; j < xm.nc; j++ {
		copy(b[j*k:(j+1)*k], xm.v[j*xm.nr:j*xm.nr+k])
	}
	if z := triSolve(r.v, r.nr, k, upper, trans, b, xm.nc); z > 0 {
		s.errorf(c.call, "singular matrix in 'backsolve'. First zero in diagonal [%d]", z)
	}
	if xIsMat {
		return dmat(k, xm.nc, b, Null, Null)
	}
	return &Double{V: b}
}

// cholFail is the order of the first leading minor that is not positive
// (LAPACK's dpotrf INFO), reading the upper triangle of the n x n a.
func cholFail(n int, a []float64) int {
	u := make([]float64, n*n)
	for j := 0; j < n; j++ {
		for i := 0; i <= j; i++ {
			u[i+j*n] = a[i+j*n]
		}
	}
	for j := 0; j < n; j++ {
		d := u[j+j*n]
		for k := 0; k < j; k++ {
			d -= u[k+j*n] * u[k+j*n]
		}
		if d <= 0 || math.IsNaN(d) {
			return j + 1
		}
		d = math.Sqrt(d)
		u[j+j*n] = d
		for c := j + 1; c < n; c++ {
			v := u[j+c*n]
			for k := 0; k < j; k++ {
				v -= u[k+j*n] * u[k+c*n]
			}
			u[j+c*n] = v / d
		}
	}
	return 0
}

// cholPivot is LAPACK's pivoted Cholesky (dpstf2, upper): the upper factor in
// place, the permutation (1-based) and the rank. tol < 0 is n eps max(diag).
func cholPivot(n int, a []float64, tol float64) (piv []int, rank int) {
	piv = make([]int, n)
	for i := range piv {
		piv[i] = i + 1
	}
	pvt, ajj := 0, a[0]
	for i := 1; i < n; i++ {
		if a[i+i*n] > ajj {
			pvt, ajj = i, a[i+i*n]
		}
	}
	if ajj <= 0 || math.IsNaN(ajj) {
		return piv, 0
	}
	dstop := tol
	if tol < 0 {
		dstop = float64(n) * math.Pow(2, -53) * 2 * ajj
	}
	work := make([]float64, 2*n)
	for j := 0; j < n; j++ {
		for i := j; i < n; i++ {
			if j > 0 {
				work[i] += a[(j-1)+i*n] * a[(j-1)+i*n]
			}
			work[n+i] = a[i+i*n] - work[i]
		}
		if j > 0 {
			pvt = j
			ajj = work[n+j]
			for i := j + 1; i < n; i++ {
				if work[n+i] > ajj {
					pvt, ajj = i, work[n+i]
				}
			}
			if ajj <= dstop || math.IsNaN(ajj) {
				a[j+j*n] = ajj
				return piv, j
			}
		}
		if j != pvt {
			a[pvt+pvt*n] = a[j+j*n]
			for k := 0; k < j; k++ {
				a[k+j*n], a[k+pvt*n] = a[k+pvt*n], a[k+j*n]
			}
			for k := pvt + 1; k < n; k++ {
				a[j+k*n], a[pvt+k*n] = a[pvt+k*n], a[j+k*n]
			}
			for k := j + 1; k < pvt; k++ {
				a[j+k*n], a[k+pvt*n] = a[k+pvt*n], a[j+k*n]
			}
			work[j], work[pvt] = work[pvt], work[j]
			piv[j], piv[pvt] = piv[pvt], piv[j]
		}
		ajj = math.Sqrt(ajj)
		a[j+j*n] = ajj
		for k := j + 1; k < n; k++ {
			v := a[j+k*n]
			for i := 0; i < j; i++ {
				v -= a[i+j*n] * a[i+k*n]
			}
			a[j+k*n] = v / ajj
		}
	}
	return piv, n
}

func (s *Session) realMatrix(c *callInfo, v Value, msg string) *lmat {
	m := s.lmatrix(c, v, msg)
	if m.cplx {
		s.errorf(c.call, "complex matrices are not supported here")
	}
	return m
}

// kappaTri is .kappa_tri: the condition number of the triangular z.
func (s *Session) kappaTri(c *callInfo, z *lmat, exact bool, norm string, uplo string) float64 {
	if z.nr != z.nc {
		s.errorf(c.call, "triangular matrix should be square")
	}
	n := z.nr
	upper := uplo != "L"
	if exact {
		t := make([]float64, n*n)
		for j := 0; j < n; j++ {
			for i := 0; i < n; i++ {
				if (upper && i <= j) || (!upper && i >= j) {
					t[i+j*n] = z.v[i+j*n]
				}
			}
		}
		return svdCond(s, c, n, n, t)
	}
	inf := false
	switch norm {
	case "", "1", "O":
	case "I":
		inf = true
	case "2":
		s.warningf(c.call, "using 1-norm approximation for approximate 2-norm")
	case "F":
		s.warningf(c.call, "norm=\"F\" not available here; using norm=\"1\"")
	default:
		s.errorf(c.call, "invalid 'norm' argument")
	}
	return 1 / rcondTriangular(n, z.v, upper, inf)
}

func svdCond(s *Session, c *callInfo, nr, nc int, v []float64) float64 {
	a := &lmat{nr: nr, nc: nc, v: v}
	sv, err := linalg.SVDVals(toND(a, nr, nc))
	if err != nil {
		s.errorf(c.call, "error code 1 from Lapack routine 'dgesdd'")
	}
	d, _ := ndarray.Data[float64](sv.Ravel())
	if len(d) == 0 {
		return math.NaN()
	}
	mx, mn := d[0], math.Inf(1)
	for _, e := range d {
		if e > mx {
			mx = e
		}
		if e < mn {
			mn = e
		}
	}
	return mx / mn
}

func init() {
	regBuiltin("determinant.matrix", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "logarithm", "...")
		if a[0] == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		d := dimsOf(a[0])
		if len(d) != 2 || d[0] != d[1] {
			s.errorf(c.call, "'x' must be a square matrix")
		}
		if _, ok := a[0].(*Complex); ok {
			s.errorf(c.call, "'determinant' not currently defined for complex matrices")
		}
		logarithm := s.optBool(c, a[1], true, "logarithm")
		x := s.lmatrix(c, a[0], "'x' must be a square matrix")
		n := x.nr
		if n == 0 {
			if logarithm {
				return detList(0, 1, true)
			}
			return detList(1, 1, false)
		}
		if !finiteAll(x.v) {
			return detList(naReal, 1, logarithm)
		}
		an := toND(x, n, n)
		if logarithm {
			sg, lg, err := linalg.Slogdet(an)
			if err != nil {
				s.errorf(c.call, "Lapack routine dgetrf: matrix is singular")
			}
			sv, _ := ndarray.Data[float64](sg.Ravel())
			lv, _ := ndarray.Data[float64](lg.Ravel())
			if sv[0] == 0 {
				return detList(math.Inf(-1), 1, true)
			}
			return detList(lv[0], sv[0], true)
		}
		dt, err := linalg.Det(an)
		if err != nil {
			s.errorf(c.call, "Lapack routine dgetrf: matrix is singular")
		}
		dv, _ := ndarray.Data[float64](dt.Ravel())
		sign := 1.0
		if dv[0] < 0 {
			sign = -1
		}
		return detList(math.Abs(dv[0]), sign, false)
	})
	regBuiltin("det", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "...")
		if a[0] == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		d := dimsOf(a[0])
		if len(d) != 2 || d[0] != d[1] {
			s.errorf(s.parseCallText("determinant.matrix(x, logarithm = TRUE, ...)"), "'x' must be a square matrix")
		}
		z := s.callNamed(c, "determinant", []Value{a[0], mkLogical(true)}, []string{"", "logarithm"})
		l := z.(*List)
		mod := asDoubles(l.V[0])[0]
		sign := asDoubles(l.V[1])[0]
		return mkDouble(sign * math.Exp(mod))
	})

	regBuiltin("chol.default", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "pivot", "LINPACK", "tol", "...")
		if a[0] == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		x := s.realMatrix(c, a[0], "'a' must be a numeric matrix")
		pivot := s.optBool(c, a[1], false, "pivot")
		tol := -1.0
		if a[3] != nil {
			tol = asDoubles(a[3])[0]
		}
		if x.nr != x.nc {
			s.errorf(c.call, "'a' must be a square matrix")
		}
		n := x.nr
		if n == 0 {
			s.errorf(c.call, "'a' must have dims > 0")
		}
		s.checkMat(c, n, n)
		var dn Value
		if d := getAttr(a[0], "dimnames"); d != nil && len(dimsOf(a[0])) == 2 {
			dn = d
		}
		if pivot {
			u := append([]float64(nil), x.v...)
			piv, rank := cholPivot(n, u, tol)
			for j := 0; j < n; j++ {
				for i := j + 1; i < n; i++ {
					u[i+j*n] = 0
				}
			}
			if rank < n {
				s.warningf(c.call, "the matrix is either rank-deficient or not positive definite")
			}
			out := Value(dmat(n, n, u, Null, Null))
			if dn != nil {
				l := dn.(*List)
				if len(l.V) == 2 {
					permute := func(v Value) Value {
						if sv, ok := v.(*String); ok && len(sv.V) == n {
							o := make([]Str, n)
							for i, p := range piv {
								o[i] = sv.V[p-1]
							}
							return &String{V: o}
						}
						return v
					}
					out = withDimnames(out, permute(l.V[0]), permute(l.V[1]))
				}
			}
			out = setAttr(out, "pivot", &Integer{V: toInt32(piv)})
			out = setAttr(out, "rank", mkInt(int32(rank)))
			return out
		}
		if !finiteAll(x.v) {
			s.errorf(c.call, "the leading minor of order 1 is not positive")
		}
		u, err := linalg.Cholesky(toND(x, n, n), linalg.Upper)
		if err != nil {
			k := cholFail(n, x.v)
			if k == 0 {
				k = n
			}
			s.errorf(c.call, "the leading minor of order %d is not positive", k)
		}
		out := Value(dmat(n, n, fromND(u, n, n), Null, Null))
		if dn != nil {
			out = setAttr(out, "dimnames", dn)
		}
		return out
	})

	regBuiltin("chol2inv", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "size", "LINPACK")
		if a[0] == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		x := s.realMatrix(c, a[0], "'a' must be a numeric matrix")
		size := x.nc
		if a[1] != nil {
			d := asDoubles(a[1])
			if len(d) != 1 || math.IsNaN(d[0]) || d[0] < 0 || d[0] > float64(math.MaxInt32) {
				s.errorf(c.call, "'size' argument must be a positive integer")
			}
			size = int(d[0])
		}
		if size > x.nc {
			s.errorf(c.call, "'size' cannot exceed ncol(x) = %d", x.nc)
		}
		if size > x.nr {
			s.errorf(c.call, "'size' cannot exceed nrow(x) = %d", x.nr)
		}
		s.checkMat(c, size, size)
		r := make([]float64, size*size)
		for j := 0; j < size; j++ {
			for i := 0; i <= j; i++ {
				r[i+j*size] = x.v[i+j*x.nr]
			}
		}
		for i := 0; i < size; i++ {
			if r[i+i*size] == 0 {
				s.errorf(c.call, "element (%d, %d) is zero, so the inverse cannot be computed", i+1, i+1)
			}
		}
		inv, _ := triInverse(size, r, true)
		out := make([]float64, size*size)
		for i := 0; i < size; i++ {
			for j := i; j < size; j++ {
				sum := 0.0
				for k := j; k < size; k++ {
					sum += inv[i+k*size] * inv[j+k*size]
				}
				out[i+j*size], out[j+i*size] = sum, sum
			}
		}
		return dmat(size, size, out, Null, Null)
	})

	regBuiltin("backsolve", func(s *Session, c *callInfo) Value { return s.backsolveFn(c, true) })
	regBuiltin("forwardsolve", func(s *Session, c *callInfo) Value { return s.backsolveFn(c, false) })

	regBuiltin("norm", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "type", "...")
		if a[0] == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		typ := "O"
		if a[1] != nil {
			typ = s.optString(c, a[1], "O", "type")
		}
		switch k := kindOf(a[0]); k {
		case kDouble, kInteger, kLogical, kComplex:
		default:
			s.errorf(c.call, "invalid 'x': type \"%s\"", typeName(a[0]))
		}
		if len(dimsOf(a[0])) != 2 {
			s.errorf(c.call, "'A' must be a numeric matrix")
		}
		t := strings.ToUpper(typ)
		switch t {
		case "O", "1", "I", "F", "E", "M", "2":
		default:
			s.errorf(c.call, "argument type[1]='%s' must be one of 'M','1','O','I','F' or 'E'", typ)
		}
		x := s.lmatrix(c, a[0], "")
		if x.nr == 0 || x.nc == 0 {
			return mkDouble(0)
		}
		mag := make([]float64, x.nr*x.nc)
		if x.cplx {
			for i, e := range x.z {
				mag[i] = cmplx.Abs(e)
			}
		} else {
			copy(mag, x.v)
		}
		for _, e := range mag {
			if math.IsNaN(e) {
				if isNAReal(e) {
					return mkDouble(naReal)
				}
				return mkDouble(math.NaN())
			}
		}
		switch t {
		case "O", "1":
			return mkDouble(norm1(x.nr, x.nc, abs64(mag)))
		case "I":
			return mkDouble(normInf(x.nr, x.nc, abs64(mag)))
		case "M":
			best := 0.0
			for _, e := range mag {
				best = math.Max(best, math.Abs(e))
			}
			return mkDouble(best)
		case "F", "E":
			nrm, err := linalg.Norm(toND(&lmat{nr: x.nr, nc: x.nc, v: mag}, x.nr, x.nc), linalg.Fro, nil, false)
			if err != nil {
				s.errorf(c.call, "error in 'norm'")
			}
			d, _ := ndarray.Data[float64](nrm.Ravel())
			return mkDouble(d[0])
		}
		sv, err := linalg.SVDVals(toND(x, x.nr, x.nc))
		if err != nil {
			s.errorf(c.call, "error code 1 from Lapack routine 'dgesdd'")
		}
		d, _ := ndarray.Data[float64](sv.Ravel())
		return mkDouble(d[0])
	})

	regBuiltin("rcond", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "norm", "triangular", "uplo", "...")
		if a[0] == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		if len(dimsOf(a[0])) != 2 {
			s.errorf(c.call, "length(d <- dim(x)) == 2L is not TRUE")
		}
		nm := "O"
		if a[1] != nil {
			nm = s.optString(c, a[1], "O", "norm")
			if nm != "O" && nm != "I" && nm != "1" {
				s.errorf(callLang("match.arg", "norm"), "'arg' should be one of “O”, “I”, “1”")
			}
		}
		tri := s.optBool(c, a[2], false, "triangular")
		uplo := s.optString(c, a[3], "U", "uplo")
		x := s.realMatrix(c, a[0], "'x' must be a numeric matrix")
		inf := nm == "I"
		if x.nr != x.nc {
			src := x
			if x.nr < x.nc {
				src = transposeL(x)
			}
			q := qrLinpack(src.nr, src.nc, src.v, 1e-7)
			r := qrRFactor(q, src.nc, src.nc)
			return mkDouble(rcondTriangular(src.nc, r, true, inf))
		}
		if !finiteAll(x.v) {
			s.errorf(c.call, "NA/NaN/Inf in foreign function call (arg 1)")
		}
		if tri {
			return mkDouble(rcondTriangular(x.nr, x.v, uplo != "L", inf))
		}
		return mkDouble(rcondGeneral(x.nr, x.v, inf))
	})

	regBuiltin(".kappa_tri", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "z", "exact", "norm", "uplo", "...")
		z := s.realMatrix(c, a[0], "'z' must be a numeric matrix")
		return mkDouble(s.kappaTri(c, z, s.optBool(c, a[1], false, "exact"), s.normArg(c, a[2]), s.optString(c, a[3], "U", "uplo")))
	})
	regBuiltin("kappa.default", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "z", "exact", "norm", "method", "inv_z", "triangular", "uplo", "...")
		if a[0] == nil {
			s.errorf(c.call, "argument \"z\" is missing, with no default")
		}
		z := s.realMatrix(c, a[0], "'z' must be a numeric matrix")
		exact := s.optBool(c, a[1], false, "exact")
		norm := s.normArg(c, a[2])
		method := "qr"
		if a[3] != nil {
			m := s.optString(c, a[3], "qr", "method")
			switch {
			case m != "" && strings.HasPrefix("qr", m):
				method = "qr"
			case m != "" && strings.HasPrefix("direct", m):
				method = "direct"
			default:
				s.errorf(callLang("match.arg", "method"), "'arg' should be one of “qr”, “direct”")
			}
		}
		tri := s.optBool(c, a[5], false, "triangular")
		uplo := s.optString(c, a[6], "U", "uplo")
		if exact {
			return mkDouble(svdCond(s, c, z.nr, z.nc, triPart(z, tri, uplo)))
		}
		if method == "direct" {
			if norm != "" && norm != "1" && norm != "O" && norm != "I" {
				s.errorf(c.call, "argument type[1]='%s' must be one of '1','O', or 'I'", norm)
			}
			if z.nr != z.nc {
				s.errorf(callLang("solve.default", "z"), "'a' (%d x %d) must be square", z.nr, z.nc)
			}
			return mkDouble(1 / rcondGeneral(z.nr, triPart(z, tri, uplo), norm == "I"))
		}
		if tri {
			return mkDouble(s.kappaTri(c, z, false, norm, uplo))
		}
		if z.nr < z.nc {
			z = transposeL(z)
		}
		q := qrLinpack(z.nr, z.nc, z.v, 1e-7)
		k := min(z.nr, z.nc)
		r := &lmat{nr: k, nc: z.nc, v: qrRFactor(q, k, z.nc)}
		return mkDouble(s.kappaTri(c, r, false, norm, "U"))
	})
}

func (s *Session) normArg(c *callInfo, v Value) string {
	if v == nil {
		return ""
	}
	if _, isNull := v.(nullT); isNull {
		return ""
	}
	return strings.ToUpper(s.optString(c, v, "", "norm"))
}

func triPart(z *lmat, tri bool, uplo string) []float64 {
	if !tri {
		return z.v
	}
	out := make([]float64, len(z.v))
	for j := 0; j < z.nc; j++ {
		for i := 0; i < z.nr; i++ {
			if (uplo != "L" && i <= j) || (uplo == "L" && i >= j) {
				out[i+j*z.nr] = z.v[i+j*z.nr]
			}
		}
	}
	return out
}

func abs64(v []float64) []float64 {
	out := make([]float64, len(v))
	for i, e := range v {
		out[i] = math.Abs(e)
	}
	return out
}

func transposeL(m *lmat) *lmat {
	out := &lmat{nr: m.nc, nc: m.nr, v: make([]float64, len(m.v))}
	for i := 0; i < m.nr; i++ {
		for j := 0; j < m.nc; j++ {
			out.v[j+i*m.nc] = m.v[i+j*m.nr]
		}
	}
	return out
}

func toInt32(v []int) []int32 {
	out := make([]int32, len(v))
	for i, e := range v {
		out[i] = int32(e)
	}
	return out
}
