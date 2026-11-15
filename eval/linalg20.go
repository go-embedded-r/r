package eval

import (
	"errors"
	"fmt"
	"math"
	"sort"

	"github.com/go-ndarray/ndarray"
	"github.com/go-ndarray/ndarray/linalg"
)

// Linear algebra (slice 20), written from the help pages of solve, qr, chol,
// svd, eigen, det, norm, rcond, kappa and backsolve and from the oracle. The
// decompositions that go-ndarray's linalg holds (solve, inverse, determinant,
// Cholesky, SVD, eigenvalues) are bound to it: this file keeps R's layer only,
// that is the call shapes, the defaults, the dimnames, the result lists and
// the messages. go-ndarray is row-major and R is column-major: a matrix is
// copied across once on each side, which costs O(n^2) against the O(n^3) of
// the kernel. What go-ndarray does not hold is written here (qr.go: LINPACK's
// dqrdc2 with limited pivoting, LAPACK's column-pivoted QR, the triangular
// solves), and listed in the README as gaps for its owner.

// lmat is an R matrix read for linear algebra: column-major, real or complex.
type lmat struct {
	nr, nc int
	v      []float64
	z      []complex128
	cplx   bool
	rn, cn *String // dimnames, nil when absent
}

func (m *lmat) at(i, j int) float64 { return m.v[i+j*m.nr] }

// lmatrix reads v as a matrix: a numeric (or complex) vector is a column.
// msg is the error of any other type.
func (s *Session) lmatrix(c *callInfo, v Value, msg string) *lmat {
	m := &lmat{}
	switch x := v.(type) {
	case *Double, *Integer, *Logical:
		m.v = asDoubles(x)
	case *Complex:
		m.z = x.V
		m.cplx = true
	default:
		s.errorf(c.call, "%s", msg)
	}
	d := dimsOf(v)
	n := Length(v)
	switch {
	case len(d) == 2:
		m.nr, m.nc = d[0], d[1]
	default:
		m.nr, m.nc = n, 1
	}
	if dn, ok := getAttr(v, "dimnames").(*List); ok && len(d) == 2 && len(dn.V) == 2 {
		m.rn, _ = dn.V[0].(*String)
		m.cn, _ = dn.V[1].(*String)
	}
	return m
}

// checkMat refuses an nr x nc result before anything is allocated.
func (s *Session) checkMat(c *callInfo, nr, nc int) {
	if nr > 0 && nc > s.opts.MaxVector/nr {
		s.errorf(c.call, "vector size too large")
	}
	s.checkLen(c.call, nr*nc)
}

func finiteAll(v []float64) bool {
	for _, e := range v {
		if math.IsNaN(e) || math.IsInf(e, 0) {
			return false
		}
	}
	return true
}

func finiteZ(v []complex128) bool {
	for _, e := range v {
		if !finiteAll([]float64{real(e), imag(e)}) {
			return false
		}
	}
	return true
}

// dmat builds an nr x nc double matrix with optional dimnames.
func dmat(nr, nc int, v []float64, rn, cn Value) *Double {
	out := &Double{V: v, Attr: []Attr{{Name: "dim", Value: &Integer{V: []int32{int32(nr), int32(nc)}}}}}
	return withDimnames(out, rn, cn).(*Double)
}

func zmat(nr, nc int, v []complex128, rn, cn Value) *Complex {
	out := &Complex{V: v, Attr: []Attr{{Name: "dim", Value: &Integer{V: []int32{int32(nr), int32(nc)}}}}}
	return withDimnames(out, rn, cn).(*Complex)
}

func nullIfNilStr(x *String) Value {
	if x == nil {
		return Null
	}
	return x
}

// withDimnames sets dimnames when either is not NULL.
func withDimnames(v Value, rn, cn Value) Value {
	if rn == nil {
		rn = Null
	}
	if cn == nil {
		cn = Null
	}
	_, a := rn.(nullT)
	_, b := cn.(nullT)
	if a && b {
		return v
	}
	return setAttr(v, "dimnames", &List{V: []Value{rn, cn}})
}

func mkNamedList(names []string, vals ...Value) *List {
	return &List{V: vals, Attr: namesAttr(strsOf(names))}
}

// toND copies a column-major matrix into a row-major go-ndarray array.
func toND(m *lmat, nr, nc int) *ndarray.Array {
	if m.cplx {
		r := make([]complex128, nr*nc)
		for i := 0; i < nr; i++ {
			for j := 0; j < nc; j++ {
				r[i*nc+j] = m.z[i+j*nr]
			}
		}
		a, _ := ndarray.FromSlice(r, nr, nc)
		return a
	}
	r := make([]float64, nr*nc)
	for i := 0; i < nr; i++ {
		for j := 0; j < nc; j++ {
			r[i*nc+j] = m.v[i+j*nr]
		}
	}
	a, _ := ndarray.FromData(r, nr, nc)
	return a
}

// fromND reads a row-major nr x nc float64 array back as column-major.
func fromND(a *ndarray.Array, nr, nc int) []float64 {
	d, ok := ndarray.Data[float64](a.Ravel())
	if !ok || len(d) != nr*nc {
		d = make([]float64, nr*nc)
	}
	out := make([]float64, nr*nc)
	for i := 0; i < nr; i++ {
		for j := 0; j < nc; j++ {
			out[i+j*nr] = d[i*nc+j]
		}
	}
	return out
}

func fromNDz(a *ndarray.Array, nr, nc int) []complex128 {
	d, ok := ndarray.Data[complex128](a.Ravel())
	if !ok || len(d) != nr*nc {
		d = make([]complex128, nr*nc)
	}
	out := make([]complex128, nr*nc)
	for i := 0; i < nr; i++ {
		for j := 0; j < nc; j++ {
			out[i+j*nr] = d[i*nc+j]
		}
	}
	return out
}

// luFirstZero is LAPACK's dgetrf diagnosis: the first k (1-based) with
// U[k,k] = 0 in the LU factorization with partial pivoting of the n x n
// column-major a, or 0 when there is none. It runs only to name the failure.
func luFirstZero(n int, a0 []float64) int {
	a := append([]float64(nil), a0...)
	first := 0
	for j := 0; j < n; j++ {
		p := j
		for i := j + 1; i < n; i++ {
			if math.Abs(a[i+j*n]) > math.Abs(a[p+j*n]) {
				p = i
			}
		}
		if a[p+j*n] == 0 {
			if first == 0 {
				first = j + 1
			}
			continue
		}
		if p != j {
			for k := 0; k < n; k++ {
				a[j+k*n], a[p+k*n] = a[p+k*n], a[j+k*n]
			}
		}
		for i := j + 1; i < n; i++ {
			a[i+j*n] /= a[j+j*n]
		}
		for k := j + 1; k < n; k++ {
			t := a[j+k*n]
			if t != 0 {
				for i := j + 1; i < n; i++ {
					a[i+k*n] -= a[i+j*n] * t
				}
			}
		}
	}
	return first
}

func norm1(nr, nc int, v []float64) float64 {
	best := 0.0
	for j := 0; j < nc; j++ {
		sum := 0.0
		for i := 0; i < nr; i++ {
			sum += math.Abs(v[i+j*nr])
		}
		if sum > best || math.IsNaN(sum) {
			best = sum
		}
	}
	return best
}

func normInf(nr, nc int, v []float64) float64 {
	best := 0.0
	for i := 0; i < nr; i++ {
		sum := 0.0
		for j := 0; j < nc; j++ {
			sum += math.Abs(v[i+j*nr])
		}
		if sum > best || math.IsNaN(sum) {
			best = sum
		}
	}
	return best
}

// rcondGeneral is the reciprocal condition number of the square n x n matrix
// in the 1-norm (inf = false) or the infinity norm, from the inverse that
// go-ndarray computes. LAPACK's dgecon estimates the norm of the inverse from
// the LU factors; this is the exact value, which the estimate bounds from
// below (README, divergences).
func rcondGeneral(n int, v []float64, inf bool) float64 {
	if n == 0 {
		return math.Inf(1)
	}
	a := &lmat{nr: n, nc: n, v: v}
	inv, err := linalg.Inv(toND(a, n, n))
	if err != nil {
		return 0
	}
	iv := fromND(inv, n, n)
	an := norm1(n, n, v)
	if inf {
		an = normInf(n, n, v)
	}
	est := lacn2(n, func(x []float64, trans bool) { matvec(n, iv, x, trans != inf) })
	if an == 0 || est == 0 || math.IsNaN(an*est) {
		return 0
	}
	return 1 / est / an
}

// matvec overwrites x with M x (or M' x) for the n x n column-major m.
func matvec(n int, m, x []float64, trans bool) {
	y := make([]float64, n)
	for j := 0; j < n; j++ {
		for i := 0; i < n; i++ {
			if trans {
				y[j] += m[i+j*n] * x[i]
			} else {
				y[i] += m[i+j*n] * x[j]
			}
		}
	}
	copy(x, y)
}

// rcondTriangular is rcondGeneral for a triangular matrix (upper or lower),
// the inverse by back substitution.
func rcondTriangular(n int, v []float64, upper, inf bool) float64 {
	if n == 0 {
		return math.Inf(1)
	}
	for i := 0; i < n; i++ {
		if v[i+i*n] == 0 {
			return 0
		}
	}
	t := make([]float64, n*n)
	for j := 0; j < n; j++ {
		for i := 0; i < n; i++ {
			if (upper && i <= j) || (!upper && i >= j) {
				t[i+j*n] = v[i+j*n]
			}
		}
	}
	an := norm1(n, n, t)
	if inf {
		an = normInf(n, n, t)
	}
	est := lacn2(n, func(x []float64, trans bool) { triSolve(t, n, n, upper, trans != inf, x, 1) })
	if an == 0 || est == 0 || math.IsNaN(an*est) {
		return 0
	}
	return 1 / est / an
}

// triInverse inverts a triangular matrix by substitution; ok is false when a
// diagonal element is zero.
func triInverse(n int, v []float64, upper bool) ([]float64, bool) {
	out := make([]float64, n*n)
	for i := 0; i < n; i++ {
		if v[i+i*n] == 0 {
			return nil, false
		}
	}
	for col := 0; col < n; col++ {
		x := out[col*n : (col+1)*n]
		if upper {
			for i := n - 1; i >= 0; i-- {
				r := 0.0
				if i == col {
					r = 1
				}
				for k := i + 1; k < n; k++ {
					r -= v[i+k*n] * x[k]
				}
				x[i] = r / v[i+i*n]
			}
		} else {
			for i := 0; i < n; i++ {
				r := 0.0
				if i == col {
					r = 1
				}
				for k := 0; k < i; k++ {
					r -= v[i+k*n] * x[k]
				}
				x[i] = r / v[i+i*n]
			}
		}
	}
	return out, true
}

func callLang(fun string, args ...string) *Lang {
	l := &Lang{Fun: Intern(fun)}
	for _, a := range args {
		l.Args = append(l.Args, Arg{Value: Intern(a)})
	}
	return l
}

func missingArg(c *callInfo, name string) {}

func init() {
	regBuiltin("solve.default", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "a", "b", "tol", "LINPACK", "...")
		if a[0] == nil {
			s.errorf(c.call, "argument \"a\" is missing, with no default")
		}
		am := s.lmatrix(c, a[0], "'a' must be a numeric matrix")
		tol := math.Pow(2, -52)
		if a[2] != nil {
			if k := kindOf(a[2]); k != kDouble && k != kInteger && k != kLogical || Length(a[2]) < 1 {
				s.errorf(c.call, "invalid 'tol' argument")
			}
			tol = asDoubles(a[2])[0]
		}
		if am.nr != am.nc {
			s.errorf(c.call, "'a' (%d x %d) must be square", am.nr, am.nc)
		}
		n := am.nr
		if n == 0 {
			s.errorf(c.call, "'a' is 0-diml")
		}
		var bm *lmat
		bVector := false
		if a[1] == nil {
			s.checkMat(c, n, n)
			bm = &lmat{nr: n, nc: n, v: make([]float64, n*n)}
			for i := 0; i < n; i++ {
				bm.v[i+i*n] = 1
			}
		} else {
			bm = s.lmatrix(c, a[1], "'b' must be numeric")
			bVector = len(dimsOf(a[1])) != 2
			if bm.nr != n {
				s.errorf(c.call, "'b' (%d x %d) must be compatible with 'a' (%d x %d)", bm.nr, bm.nc, n, n)
			}
		}
		cplx := am.cplx || bm.cplx
		if cplx {
			promoteZ(am)
			promoteZ(bm)
			if !finiteZ(am.z) || !finiteZ(bm.z) {
				return nanSolve(n, bm.nc, bVector, am, bm, true)
			}
		} else if !finiteAll(am.v) || !finiteAll(bm.v) {
			return nanSolve(n, bm.nc, bVector, am, bm, false)
		}
		an, bn := toND(am, n, n), toND(bm, n, bm.nc)
		x, err := linalg.Solve(an, bn)
		if err != nil {
			if !cplx {
				if k := luFirstZero(n, am.v); k > 0 {
					s.errorf(c.call, "Lapack routine dgesv: system is exactly singular: U[%d,%d] = 0", k, k)
				}
			} else {
				s.errorf(c.call, "Lapack routine zgesv: system is exactly singular")
			}
			s.errorf(c.call, "Lapack routine dgesv: system is exactly singular: U[%d,%d] = 0", n, n)
		}
		if !cplx {
			if rc := rcondGeneral(n, am.v, false); rc < tol {
				s.errorf(c.call, "system is computationally singular: reciprocal condition number = %.6g", rc)
			}
		}
		var rn, cn Value = nullIfNilStr(am.cn), Null
		if a[1] == nil {
			cn = nullIfNilStr(am.rn)
		} else if bm.cn != nil {
			cn = bm.cn
		}
		if bVector {
			var out Value
			if cplx {
				out = &Complex{V: fromNDz(x, n, 1)}
			} else {
				out = &Double{V: fromND(x, n, 1)}
			}
			if am.cn != nil {
				out = setAttr(out, "names", am.cn)
			}
			return out
		}
		if cplx {
			return zmat(n, bm.nc, fromNDz(x, n, bm.nc), rn, cn)
		}
		return dmat(n, bm.nc, fromND(x, n, bm.nc), rn, cn)
	})
}

func promoteZ(m *lmat) {
	if m.cplx {
		return
	}
	m.z = make([]complex128, len(m.v))
	for i, e := range m.v {
		m.z[i] = complex(e, 0)
	}
	m.cplx = true
}

// nanSolve answers a system with a non-finite entry: NA throughout, as LAPACK
// propagates it (README, divergences: an infinity is not reproduced).
func nanSolve(n, k int, vec bool, am, bm *lmat, cplx bool) Value {
	if cplx {
		z := make([]complex128, n*k)
		for i := range z {
			z[i] = complex(naReal, 0)
		}
		return zmat(n, k, z, nullIfNilStr(am.cn), nullIfNilStr(bm.cn))
	}
	v := make([]float64, n*k)
	for i := range v {
		v[i] = naReal
	}
	if vec {
		return &Double{V: v}
	}
	return dmat(n, k, v, nullIfNilStr(am.cn), nullIfNilStr(bm.cn))
}

// determinant.matrix, det, chol, chol2inv, backsolve, forwardsolve, norm,
// rcond: see linalg20b.go.

var _ = errors.Is
var _ = fmt.Sprintf
var _ = sort.Ints

// lacn2 estimates the 1-norm of the operator M (n x n) that apply(x, trans)
// multiplies a vector by (x := M x, or M' x when trans): Hager's method as
// improved by Higham, the algorithm of LAPACK's dlacn2 (netlib, public
// domain) that dgecon and dtrcon use. R's rcond and kappa report this
// estimate, so the estimate is what is reproduced.
func lacn2(n int, apply func(x []float64, trans bool)) float64 {
	const itmax = 5
	x := make([]float64, n)
	v := make([]float64, n)
	isgn := make([]int, n)
	sign1 := func(a float64) float64 {
		if a >= 0 {
			return 1
		}
		return -1
	}
	asum := func(a []float64) float64 {
		t := 0.0
		for _, e := range a {
			t += math.Abs(e)
		}
		return t
	}
	imax := func(a []float64) int {
		k := 0
		for i := 1; i < len(a); i++ {
			if math.Abs(a[i]) > math.Abs(a[k]) {
				k = i
			}
		}
		return k
	}
	for i := range x {
		x[i] = 1 / float64(n)
	}
	apply(x, false)
	if n == 1 {
		return math.Abs(x[0])
	}
	est := asum(x)
	for i := range x {
		x[i] = sign1(x[i])
		isgn[i] = int(x[i])
	}
	apply(x, true)
	j := imax(x)
	iter := 2
	for {
		for i := range x {
			x[i] = 0
		}
		x[j] = 1
		apply(x, false)
		copy(v, x)
		estold := est
		est = asum(v)
		same := true
		for i := range x {
			if int(sign1(x[i])) != isgn[i] {
				same = false
				break
			}
		}
		if same || est <= estold {
			break
		}
		for i := range x {
			x[i] = sign1(x[i])
			isgn[i] = int(x[i])
		}
		apply(x, true)
		jlast := j
		j = imax(x)
		if x[jlast] != math.Abs(x[j]) && iter < itmax {
			iter++
			continue
		}
		break
	}
	alt := 1.0
	for i := range x {
		x[i] = alt * (1 + float64(i)/float64(n-1))
		alt = -alt
	}
	apply(x, false)
	if t := 2 * (asum(x) / float64(3*n)); t > est {
		est = t
	}
	return est
}
