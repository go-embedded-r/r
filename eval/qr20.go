package eval

import (
	"math"
)

// The QR decomposition and its family (slice 20). go-ndarray's linalg.QR is
// numpy's: unpivoted, one convention for the compact form. R's qr() is LINPACK's
// dqrdc2 (Householder with LIMITED pivoting: a column whose norm falls below tol
// times its original norm is moved to the end, which is what rank means), and
// qr(LAPACK = TRUE) is LAPACK's dgeqp3 (column pivoting on the largest norm). Both
// are written here from the netlib routines (public domain / BSD) and the help
// page of qr, and listed in the README as what go-ndarray lacks.

type qrObj struct {
	n, p   int
	a      []float64 // n x p, column-major: R above, Householder vectors below
	qraux  []float64
	pivot  []int // 1-based
	rank   int
	lapack bool
}

func nrm2(x []float64) float64 {
	scale, ssq := 0.0, 1.0
	for _, v := range x {
		if v != 0 {
			av := math.Abs(v)
			if scale < av {
				ssq = 1 + ssq*(scale/av)*(scale/av)
				scale = av
			} else {
				ssq += (av / scale) * (av / scale)
			}
		}
	}
	return scale * math.Sqrt(ssq)
}

func dot(x, y []float64) float64 {
	s := 0.0
	for i := range x {
		s += x[i] * y[i]
	}
	return s
}

// qrLinpack is dqrdc2.
func qrLinpack(n, p int, x []float64, tol float64) *qrObj {
	a := append([]float64(nil), x...)
	qraux := make([]float64, p)
	w1 := make([]float64, p)
	w2 := make([]float64, p)
	piv := make([]int, p)
	for j := 0; j < p; j++ {
		piv[j] = j + 1
		qraux[j] = nrm2plain(a[j*n : (j+1)*n])
		w1[j], w2[j] = qraux[j], qraux[j]
		if w2[j] == 0 {
			w2[j] = 1
		}
	}
	lup := min(n, p)
	k := p + 1
	for l := 1; l <= lup; l++ {
		for l < k && qraux[l-1] < w2[l-1]*tol {
			for i := 0; i < n; i++ {
				t := a[i+(l-1)*n]
				for j := l + 1; j <= p; j++ {
					a[i+(j-2)*n] = a[i+(j-1)*n]
				}
				a[i+(p-1)*n] = t
			}
			ip, t, tt, ttt := piv[l-1], qraux[l-1], w1[l-1], w2[l-1]
			for j := l + 1; j <= p; j++ {
				piv[j-2], qraux[j-2], w1[j-2], w2[j-2] = piv[j-1], qraux[j-1], w1[j-1], w2[j-1]
			}
			piv[p-1], qraux[p-1], w1[p-1], w2[p-1] = ip, t, tt, ttt
			k--
		}
		if l == n {
			break
		}
		ll := (l - 1) + (l-1)*n
		nrmxl := nrm2plain(a[ll : l*n])
		if nrmxl != 0 {
			if a[ll] != 0 {
				nrmxl = math.Copysign(nrmxl, a[ll])
			}
			for i := l - 1; i < n; i++ {
				a[i+(l-1)*n] /= nrmxl
			}
			a[ll] = 1 + a[ll]
			for j := l + 1; j <= p; j++ {
				col := a[(l-1)+(l-1)*n : l*n]
				cj := a[(l-1)+(j-1)*n : j*n]
				t := -dot(col, cj) / a[ll]
				for i := range col {
					cj[i] += t * col[i]
				}
				if qraux[j-1] != 0 {
					tt := 1 - (math.Abs(a[(l-1)+(j-1)*n])/qraux[j-1])*(math.Abs(a[(l-1)+(j-1)*n])/qraux[j-1])
					tt = math.Max(tt, 0)
					t = tt
					if math.Abs(t) < 1e-6 {
						qraux[j-1] = nrm2plain(a[l+(j-1)*n : j*n])
						w1[j-1] = qraux[j-1]
					} else {
						qraux[j-1] *= math.Sqrt(t)
					}
				}
			}
			qraux[l-1] = a[ll]
			a[ll] = -nrmxl
		}
	}
	return &qrObj{n: n, p: p, a: a, qraux: qraux, pivot: piv, rank: min(k-1, n)}
}

// qrLapack is dgeqp3 (unblocked, as dlaqp2): every column free to move.
func qrLapack(n, p int, x []float64) *qrObj {
	a := append([]float64(nil), x...)
	k := min(n, p)
	tau := make([]float64, k)
	piv := make([]int, p)
	vn1 := make([]float64, p)
	vn2 := make([]float64, p)
	for j := 0; j < p; j++ {
		piv[j] = j + 1
		vn1[j] = nrm2(a[j*n : (j+1)*n])
		vn2[j] = vn1[j]
	}
	tol3z := math.Sqrt(math.Pow(2, -53))
	for i := 0; i < k; i++ {
		pvt := i
		for j := i + 1; j < p; j++ {
			if vn1[j] > vn1[pvt] {
				pvt = j
			}
		}
		if pvt != i {
			for r := 0; r < n; r++ {
				a[r+pvt*n], a[r+i*n] = a[r+i*n], a[r+pvt*n]
			}
			piv[pvt], piv[i] = piv[i], piv[pvt]
			vn1[pvt], vn2[pvt] = vn1[i], vn2[i]
		}
		if i < n-1 {
			alpha := a[i+i*n]
			xn := nrm2(a[i+1+i*n : (i+1)*n])
			if xn == 0 {
				tau[i] = 0
			} else {
				beta := -math.Copysign(math.Hypot(alpha, xn), alpha)
				tau[i] = (beta - alpha) / beta
				sc := 1 / (alpha - beta)
				for r := i + 1; r < n; r++ {
					a[r+i*n] *= sc
				}
				a[i+i*n] = beta
			}
		} else {
			tau[i] = 0
		}
		if i < p-1 && tau[i] != 0 {
			for j := i + 1; j < p; j++ {
				t := a[i+j*n]
				for r := i + 1; r < n; r++ {
					t += a[r+i*n] * a[r+j*n]
				}
				t *= tau[i]
				a[i+j*n] -= t
				for r := i + 1; r < n; r++ {
					a[r+j*n] -= t * a[r+i*n]
				}
			}
		}
		for j := i + 1; j < p; j++ {
			if vn1[j] != 0 {
				r := math.Abs(a[i+j*n]) / vn1[j]
				temp := math.Max(1-r*r, 0)
				temp2 := temp * (vn1[j] / vn2[j]) * (vn1[j] / vn2[j])
				if temp2 <= tol3z {
					if i < n-1 {
						vn1[j] = nrm2(a[i+1+j*n : (j+1)*n])
						vn2[j] = vn1[j]
					} else {
						vn1[j], vn2[j] = 0, 0
					}
				} else {
					vn1[j] *= math.Sqrt(temp)
				}
			}
		}
	}
	return &qrObj{n: n, p: p, a: a, qraux: tau, pivot: piv, rank: k, lapack: true}
}

// reflect applies reflector j to y (length n).
func (q *qrObj) reflect(j int, y []float64) {
	n := q.n
	if q.lapack {
		tau := q.qraux[j]
		if tau == 0 {
			return
		}
		t := y[j]
		for r := j + 1; r < n; r++ {
			t += q.a[r+j*n] * y[r]
		}
		t *= tau
		y[j] -= t
		for r := j + 1; r < n; r++ {
			y[r] -= t * q.a[r+j*n]
		}
		return
	}
	u0 := q.qraux[j]
	if u0 == 0 {
		return
	}
	t := u0 * y[j]
	for r := j + 1; r < n; r++ {
		t += q.a[r+j*n] * y[r]
	}
	t = -t / u0
	y[j] += t * u0
	for r := j + 1; r < n; r++ {
		y[r] += t * q.a[r+j*n]
	}
}

func (q *qrObj) steps() int {
	if q.lapack {
		return min(q.n, q.p)
	}
	return min(q.rank, q.n-1)
}

func (q *qrObj) qty(y []float64) {
	for j := 0; j < q.steps(); j++ {
		q.reflect(j, y)
	}
}

func (q *qrObj) qy(y []float64) {
	for j := q.steps() - 1; j >= 0; j-- {
		q.reflect(j, y)
	}
}

// coef solves R b = (Q'y)[1:rank] for one column; the unsolved are NA.
func (q *qrObj) coef(y []float64) []float64 {
	qt := append([]float64(nil), y...)
	q.qty(qt)
	k := q.rank
	b := append([]float64(nil), qt[:k]...)
	for i := k - 1; i >= 0; i-- {
		r := b[i]
		for j := i + 1; j < k; j++ {
			r -= q.a[i+j*q.n] * b[j]
		}
		b[i] = r / q.a[i+i*q.n]
	}
	out := make([]float64, q.p)
	for i := range out {
		out[i] = naReal
	}
	for i := 0; i < k; i++ {
		out[q.pivot[i]-1] = b[i]
	}
	return out
}

// qrRFactor is the rows x cols upper trapezoid of the compact factor.
func qrRFactor(q *qrObj, rows, cols int) []float64 {
	r := make([]float64, rows*cols)
	for j := 0; j < cols; j++ {
		for i := 0; i < rows && i <= j; i++ {
			if i < q.n && j < q.p {
				r[i+j*rows] = q.a[i+j*q.n]
			}
		}
	}
	return r
}

// ---- R's layer -------------------------------------------------------------

func (s *Session) qrOf(c *callInfo, v Value) (*qrObj, *lmat, bool) {
	l, ok := v.(*List)
	if !ok || !inheritsFrom(v, "qr") {
		return nil, nil, false
	}
	get := func(name string) Value {
		if nm := namesOf(l); nm != nil {
			for i, e := range nm.V {
				if e.S == name {
					return l.V[i]
				}
			}
		}
		return nil
	}
	qm, rk, qa, pv := get("qr"), get("rank"), get("qraux"), get("pivot")
	if qm == nil || rk == nil || qa == nil || pv == nil || len(dimsOf(qm)) != 2 {
		return nil, nil, false
	}
	m := s.lmatrix(c, qm, "invalid 'qr'")
	q := &qrObj{n: m.nr, p: m.nc, a: m.v, qraux: asDoubles(qa), rank: int(asDoubles(rk)[0])}
	for _, e := range asDoubles(pv) {
		q.pivot = append(q.pivot, int(e))
	}
	if len(q.qraux) != q.p || len(q.pivot) != q.p || m.cplx {
		return nil, nil, false
	}
	if ul, ok := getAttr(v, "useLAPACK").(*Logical); ok && len(ul.V) == 1 && ul.V[0] == 1 {
		q.lapack = true
	}
	return q, m, true
}

func (s *Session) needQR(c *callInfo, v Value) (*qrObj, *lmat) {
	q, m, ok := s.qrOf(c, v)
	if !ok {
		s.errorf(c.call, "first argument must be a QR decomposition")
	}
	return q, m
}

// origColnames are the column names of the decomposed matrix, in its order.
func origColnames(q *qrObj, m *lmat) Value {
	if m.cn == nil {
		return Null
	}
	out := make([]Str, len(m.cn.V))
	for i, p := range q.pivot {
		if p >= 1 && p <= len(out) {
			out[p-1] = m.cn.V[i]
		}
	}
	return &String{V: out}
}

// yColumns reads y as the columns of an n-row matrix, vector or matrix.
func (s *Session) yColumns(c *callInfo, y Value, n int) (cols [][]float64, isMat bool, ny int) {
	if k := kindOf(y); k != kDouble && k != kInteger && k != kLogical {
		s.errorf(c.call, "'qr' and 'y' must have the same number of rows")
	}
	d := dimsOf(y)
	v := asDoubles(y)
	nr, nc := len(v), 1
	if len(d) == 2 {
		nr, nc, isMat = d[0], d[1], true
	}
	if nr != n {
		s.errorf(c.call, "'qr' and 'y' must have the same number of rows")
	}
	for j := 0; j < nc; j++ {
		cols = append(cols, append([]float64(nil), v[j*nr:(j+1)*nr]...))
	}
	return cols, isMat, nc
}

func (s *Session) qrApply(c *callInfo, which string) Value {
	a := s.matchFormals(c, "qr", "y")
	if a[0] == nil || a[1] == nil {
		s.errorf(c.call, "argument \"%s\" is missing, with no default", map[bool]string{true: "qr", false: "y"}[a[0] == nil])
	}
	q, _ := s.needQR(c, a[0])
	cols, isMat, _ := s.yColumns(c, a[1], q.n)
	out := make([]float64, 0, len(cols)*q.n)
	for _, col := range cols {
		switch which {
		case "qy":
			q.qy(col)
		case "qty":
			q.qty(col)
		}
		out = append(out, col...)
	}
	if q.lapack && !isMat {
		// LAPACK's dormqr works on matrices: a vector comes back as a column
		return dmat(q.n, 1, out, Null, Null)
	}
	return &Double{V: out, Attr: attrsOf(a[1])}
}

func (s *Session) qrCoefFn(c *callInfo, q *qrObj, m *lmat, y Value) Value {
	cols, isMat, ny := s.yColumns(c, y, q.n)
	s.checkMat(c, q.p, ny)
	out := make([]float64, 0, q.p*ny)
	for _, col := range cols {
		out = append(out, q.coef(col)...)
	}
	cn := origColnames(q, m)
	if !isMat {
		r := &Double{V: out}
		if _, isNull := cn.(nullT); !isNull {
			return setAttr(r, "names", cn)
		}
		return r
	}
	var ycn Value = Null
	if dn, ok := getAttr(y, "dimnames").(*List); ok && len(dn.V) == 2 {
		ycn = dn.V[1]
	}
	return dmat(q.p, ny, out, cn, ycn)
}

// firstDouble is the first element of a numeric argument; an empty one (a
// missing argument as in qr(,), or numeric()) is an error, not a panic.
func (s *Session) firstDouble(c *callInfo, v Value, name string) float64 {
	d := asDoubles(v)
	if len(d) == 0 {
		s.errorf(c.call, "invalid '%s' argument", name)
	}
	return d[0]
}

func (s *Session) qrAs(c *callInfo, x Value, tol float64, lapack bool) Value {
	m := s.lmatrix(c, x, "")
	_ = m
	return nil
}

func init() {
	regBuiltin("qr.default", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "tol", "LAPACK", "...")
		if a[0] == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		tol := 1e-7
		if a[1] != nil {
			tol = s.firstDouble(c, a[1], "tol")
		}
		lapack := s.optBool(c, a[2], false, "LAPACK")
		var m *lmat
		switch a[0].(type) {
		case *String:
			s.warningf(c.call, "NAs introduced by coercion")
			s.errorf(c.call, "NA/NaN/Inf in foreign function call (arg 1)")
		case *Complex:
			s.errorf(c.call, "qr of a complex matrix is not implemented (README, deferrals)")
		}
		m = s.lmatrix(c, a[0], "'x' must be a numeric matrix")
		if !finiteAll(m.v) {
			s.errorf(c.call, "NA/NaN/Inf in foreign function call (arg 1)")
		}
		s.checkMat(c, m.nr, m.nc)
		var q *qrObj
		if lapack {
			q = qrLapack(m.nr, m.nc, m.v)
		} else {
			q = qrLinpack(m.nr, m.nc, m.v, tol)
		}
		var cn Value = nullIfNilStr(m.cn)
		if m.cn != nil {
			o := make([]Str, len(m.cn.V))
			for i, p := range q.pivot {
				o[i] = m.cn.V[p-1]
			}
			cn = &String{V: o}
		}
		qm := dmat(m.nr, m.nc, q.a, nullIfNilStr(m.rn), cn)
		out := mkNamedList([]string{"qr", "rank", "qraux", "pivot"}, qm, mkInt(int32(q.rank)), &Double{V: q.qraux}, &Integer{V: toInt32(q.pivot)})
		var res Value = out
		if lapack {
			res = setAttr(res, "useLAPACK", mkLogical(true))
		}
		return setAttr(res, "class", mkString("qr"))
	})
	regBuiltin("is.qr", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x")
		if a[0] == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		return mkLogical(inheritsFrom(a[0], "qr"))
	})
	regBuiltin("as.qr", func(s *Session, c *callInfo) Value {
		s.errorf(c.call, "you cannot be serious")
		return Null
	})
	regBuiltin("qr.qy", func(s *Session, c *callInfo) Value { return s.qrApply(c, "qy") })
	regBuiltin("qr.qty", func(s *Session, c *callInfo) Value { return s.qrApply(c, "qty") })
	regBuiltin("qr.coef", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "qr", "y")
		if a[0] == nil || a[1] == nil {
			s.errorf(c.call, "argument \"%s\" is missing, with no default", map[bool]string{true: "qr", false: "y"}[a[0] == nil])
		}
		q, m := s.needQR(c, a[0])
		return s.qrCoefFn(c, q, m, a[1])
	})
	fitres := func(name string, resid bool) {
		regBuiltin(name, func(s *Session, c *callInfo) Value {
			a := s.matchFormals(c, "qr", "y", "k")
			if a[0] == nil || a[1] == nil {
				s.errorf(c.call, "argument \"%s\" is missing, with no default", map[bool]string{true: "qr", false: "y"}[a[0] == nil])
			}
			q, _ := s.needQR(c, a[0])
			if q.lapack {
				s.errorf(c.call, "not supported for LAPACK QR")
			}
			k := q.rank
			if !resid && a[2] != nil {
				k = int(s.firstDouble(c, a[2], "k"))
			}
			if k < 0 || k > q.rank {
				s.errorf(c.call, "invalid 'k' argument")
			}
			cols, _, _ := s.yColumns(c, a[1], q.n)
			out := make([]float64, 0, len(cols)*q.n)
			for _, col := range cols {
				q.qty(col)
				for i := range col {
					if (!resid && i >= k) || (resid && i < k) {
						col[i] = 0
					}
				}
				q.qy(col)
				out = append(out, col...)
			}
			return &Double{V: out, Attr: attrsOf(a[1])}
		})
	}
	fitres("qr.fitted", false)
	fitres("qr.resid", true)
	regBuiltin("qr.R", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "qr", "complete", "...")
		if a[0] == nil {
			s.errorf(c.call, "argument \"qr\" is missing, with no default")
		}
		q, m := s.needQR(c, a[0])
		complete := s.optBool(c, a[1], false, "complete")
		rows := min(q.n, q.p)
		if complete {
			rows = q.n
		}
		s.checkMat(c, rows, q.p)
		r := make([]float64, rows*q.p)
		for j := 0; j < q.p; j++ {
			for i := 0; i < rows && i <= j; i++ {
				r[i+j*rows] = q.a[i+j*q.n]
			}
		}
		var rn Value = Null
		if m.rn != nil {
			rn = &String{V: m.rn.V[:rows]}
		}
		return dmat(rows, q.p, r, rn, nullIfNilStr(m.cn))
	})
	regBuiltin("qr.Q", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "qr", "complete", "Dvec")
		if a[0] == nil {
			s.errorf(c.call, "argument \"qr\" is missing, with no default")
		}
		q, _ := s.needQR(c, a[0])
		complete := s.optBool(c, a[1], false, "complete")
		ncol := min(q.n, q.p)
		if complete {
			ncol = q.n
		}
		s.checkMat(c, q.n, ncol)
		var dvec []float64
		if a[2] != nil {
			dvec = asDoubles(a[2])
		} else {
			dvec = make([]float64, ncol)
			for i := range dvec {
				dvec[i] = 1
			}
		}
		out := make([]float64, q.n*ncol)
		for j := 0; j < ncol; j++ {
			col := out[j*q.n : (j+1)*q.n]
			if j < len(dvec) && j < q.n {
				col[j] = dvec[j]
			}
			q.qy(col)
		}
		return dmat(q.n, ncol, out, Null, Null)
	})
	regBuiltin("qr.X", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "qr", "complete", "ncol")
		if a[0] == nil {
			s.errorf(c.call, "argument \"qr\" is missing, with no default")
		}
		q, m := s.needQR(c, a[0])
		complete := s.optBool(c, a[1], false, "complete")
		ncol := q.p
		if a[2] != nil {
			nc := asDoubles(a[2])
			if len(nc) != 1 || math.IsNaN(nc[0]) || nc[0] < 0 {
				s.errorf(c.call, "invalid 'ncol' argument")
			}
			if nc[0] > float64(s.opts.MaxVector) {
				s.checkLen(c.call, s.opts.MaxVector+1) // not int(nc[0]): on 32 bits it wraps
			}
			ncol = int(nc[0])
		} else if complete && q.n > q.p {
			ncol = q.n
		}
		rows := q.n
		s.checkMat(c, rows, max(ncol, q.p))
		r := make([]float64, q.n*q.p)
		for j := 0; j < q.p; j++ {
			for i := 0; i < q.n && i <= j; i++ {
				r[i+j*q.n] = q.a[i+j*q.n]
			}
		}
		for j := 0; j < q.p; j++ {
			q.qy(r[j*q.n : (j+1)*q.n])
		}
		un := make([]float64, q.n*q.p)
		for j := 0; j < q.p; j++ {
			copy(un[(q.pivot[j]-1)*q.n:q.pivot[j]*q.n], r[j*q.n:(j+1)*q.n])
		}
		if ncol < q.p {
			un = un[:q.n*ncol]
		}
		cn := origColnames(q, m)
		if ncol != q.p {
			cn = Null
		}
		return dmat(q.n, min(ncol, q.p), un, nullIfNilStr(m.rn), cn)
	})
	solveQR := func(name string, fromMatrix bool) {
		regBuiltin(name, func(s *Session, c *callInfo) Value {
			var a []Value
			if fromMatrix {
				a = s.matchFormals(c, "a", "b", "tol")
			} else {
				a = s.matchFormals(c, "a", "b", "...")
			}
			if a[0] == nil {
				s.errorf(c.call, "argument \"a\" is missing, with no default")
			}
			x := a[0]
			if fromMatrix && !inheritsFrom(x, "qr") {
				tol := 1e-7
				if a[2] != nil {
					tol = s.firstDouble(c, a[2], "tol")
				}
				x = s.callNamed(c, "qr", []Value{x, mkDouble(tol)}, []string{"", "tol"})
			} else if !fromMatrix && !inheritsFrom(x, "qr") {
				s.errorf(c.call, "this is the \"qr\" method for the generic function solve()")
			}
			q, m := s.needQR(c, x)
			if q.rank != min(q.n, q.p) {
				s.errorf(c.call, "singular matrix 'a' in solve")
			}
			var b Value
			if a[1] == nil {
				if q.p != q.n {
					s.errorf(c.call, "only square matrices can be inverted")
				}
				s.checkMat(c, q.p, q.p)
				id := make([]float64, q.p*q.p)
				for i := 0; i < q.p; i++ {
					id[i+i*q.p] = 1
				}
				b = dmat(q.p, q.p, id, Null, Null)
			} else {
				b = a[1]
			}
			res := s.qrCoefFn(c, q, m, b)
			d := res.(*Double)
			for i, e := range d.V {
				if math.IsNaN(e) && isNAReal(e) {
					d.V[i] = 0
				}
			}
			return res
		})
	}
	solveQR("qr.solve", true)
	solveQR("solve.qr", false)
}

// nrm2plain is the Euclidean norm as the reference BLAS of LAPACK 3.10 and
// later computes it for values in the middle of the range: the square root of
// the plain sum of squares.
func nrm2plain(x []float64) float64 {
	s := 0.0
	for _, v := range x {
		av := math.Abs(v)
		if av != 0 && (av < 1e-140 || av > 1e140) || math.IsNaN(av) || math.IsInf(av, 0) {
			return nrm2(x)
		}
		s += v * v
	}
	return math.Sqrt(s)
}
