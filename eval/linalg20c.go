package eval

import (
	"math"
	"math/cmplx"

	"github.com/go-ndarray/ndarray"
	"github.com/go-ndarray/ndarray/linalg"
)

// svd, La.svd and eigen (slice 20) over go-ndarray's linalg.SVD, Eigh and Eig.

func ndData(a *ndarray.Array) []float64 {
	d, _ := ndarray.Data[float64](a.Ravel())
	return d
}

func ndDataZ(a *ndarray.Array) []complex128 {
	d, _ := ndarray.Data[complex128](a.Ravel())
	return d
}

// svdParts is La.svd's core: d, and u (nr x nu), vt (nv x nc) when asked.
func (s *Session) svdFn(c *callInfo, vt bool) Value {
	a := s.matchFormals(c, "x", "nu", "nv", "LINPACK")
	if a[0] == nil {
		s.errorf(c.call, "argument \"x\" is missing, with no default")
	}
	x := s.realMatrix(c, a[0], "a numeric or complex matrix is required")
	n, p := x.nr, x.nc
	k := min(n, p)
	count := func(v Value, name string) int {
		if v == nil {
			return k
		}
		d := asDoubles(v)
		if len(d) != 1 || math.IsNaN(d[0]) || d[0] < 0 {
			s.errorf(s.parseCallText("seq_len(min(n, "+name+"))"), "argument must be coercible to non-negative integer")
		}
		return int(math.Min(d[0], float64(math.MaxInt32)))
	}
	nu, nv := count(a[1], "nu"), count(a[2], "nv")
	if n == 0 || p == 0 {
		s.errorf(c.call, "a dimension is zero")
	}
	if !finiteAll(x.v) {
		s.errorf(c.call, "infinite or missing values in 'x'")
	}
	full := nu > k || nv > k
	if nu > 0 || nv > 0 {
		if full {
			s.checkMat(c, n, n)
			s.checkMat(c, p, p)
		} else {
			s.checkMat(c, n, k)
			s.checkMat(c, k, p)
		}
	}
	an := toND(x, n, p)
	names := []string{"d"}
	if nu == 0 && nv == 0 {
		sv, err := linalg.SVDVals(an)
		if err != nil {
			s.errorf(c.call, "error code 1 from Lapack routine 'dgesdd'")
		}
		return mkNamedList(names, &Double{V: ndData(sv)})
	}
	u, sv, vh, err := linalg.SVD(an, full)
	if err != nil {
		s.errorf(c.call, "error code 1 from Lapack routine 'dgesdd'")
	}
	vals := []Value{&Double{V: ndData(sv)}}
	ucols, vrows := k, k
	if full {
		ucols, vrows = n, p
	}
	if nu > 0 {
		uv := fromND(u, n, ucols)
		keep := min(nu, ucols)
		vals = append(vals, dmat(n, keep, uv[:n*keep], Null, Null))
		names = append(names, "u")
	}
	if nv > 0 {
		vv := fromND(vh, vrows, p)
		keep := min(nv, vrows)
		if vt {
			out := make([]float64, keep*p)
			for i := 0; i < keep; i++ {
				for j := 0; j < p; j++ {
					out[i+j*keep] = vv[i+j*vrows]
				}
			}
			vals = append(vals, dmat(keep, p, out, Null, Null))
			names = append(names, "vt")
		} else {
			out := make([]float64, p*keep)
			for i := 0; i < keep; i++ {
				for j := 0; j < p; j++ {
					out[j+i*p] = vv[i+j*vrows]
				}
			}
			vals = append(vals, dmat(p, keep, out, Null, Null))
			names = append(names, "v")
		}
	}
	return mkNamedList(names, vals...)
}

func init() {
	regBuiltin("svd", func(s *Session, c *callInfo) Value { return s.svdFn(c, false) })
	regBuiltin("La.svd", func(s *Session, c *callInfo) Value { return s.svdFn(c, true) })

	regBuiltin("eigen", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "symmetric", "only.values", "EISPACK")
		if a[0] == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		x := s.lmatrix(c, a[0], "a numeric or complex matrix is required")
		if x.nr != x.nc {
			s.errorf(c.call, "non-square matrix in 'eigen'")
		}
		n := x.nr
		if n == 0 {
			s.errorf(c.call, "0 x 0 matrix")
		}
		var sym bool
		if a[1] != nil {
			sym = s.optBool(c, a[1], false, "symmetric")
		} else if x.cplx {
			sym = hermitian(n, x.z)
		} else {
			un := a[0]
			if len(attrsOf(a[0])) > 0 {
				un = &Double{V: x.v}
				if x.cplx {
					un = nil
				}
				if un != nil {
					un = setAttr(un, "dim", &Integer{V: []int32{int32(n), int32(n)}})
				} else {
					un = setAttr(&Complex{V: x.z}, "dim", &Integer{V: []int32{int32(n), int32(n)}})
				}
			}
			r := s.callNamed(c, "isSymmetric", []Value{un}, []string{""})
			l, ok := r.(*Logical)
			sym = ok && len(l.V) == 1 && l.V[0] == 1
		}
		only := s.optBool(c, a[2], false, "only.values")
		if x.cplx {
			if !finiteZ(x.z) {
				s.errorf(c.call, "infinite or missing values in 'x'")
			}
		} else if !finiteAll(x.v) {
			s.errorf(c.call, "infinite or missing values in 'x'")
		}
		s.checkMat(c, n, n)
		an := toND(x, n, n)
		var vals Value
		var vecs Value = Null
		if sym {
			var w, v *ndarray.Array
			var err error
			if only {
				w, err = linalg.Eigvalsh(an, linalg.Lower)
			} else {
				w, v, err = linalg.Eigh(an, linalg.Lower)
			}
			if err != nil {
				s.errorf(c.call, "error code 1 from Lapack routine 'dsyevr'")
			}
			wd := ndData(w)
			rev := make([]float64, n)
			for i := range rev {
				rev[i] = wd[n-1-i]
			}
			vals = &Double{V: rev}
			if !only {
				if x.cplx {
					vz := fromNDz(v, n, n)
					out := make([]complex128, n*n)
					for j := 0; j < n; j++ {
						copy(out[j*n:(j+1)*n], vz[(n-1-j)*n:(n-j)*n])
					}
					vecs = zmat(n, n, out, Null, Null)
				} else {
					vv := fromND(v, n, n)
					out := make([]float64, n*n)
					for j := 0; j < n; j++ {
						copy(out[j*n:(j+1)*n], vv[(n-1-j)*n:(n-j)*n])
					}
					vecs = dmat(n, n, out, Null, Null)
				}
			}
		} else {
			var w, v *ndarray.Array
			var err error
			if only {
				w, err = linalg.Eigvals(an)
			} else {
				w, v, err = linalg.Eig(an)
			}
			if err != nil {
				s.errorf(c.call, "error code 1 from Lapack routine 'dgeev'")
			}
			wd := ndDataZ(w)
			ord := sortedDesc(n, func(i int) float64 { return cmplx.Abs(wd[i]) })
			anyIm := false
			for _, e := range wd {
				if imag(e) != 0 {
					anyIm = true
				}
			}
			if anyIm {
				out := make([]complex128, n)
				for i, o := range ord {
					out[i] = wd[o]
				}
				vals = &Complex{V: out}
			} else {
				out := make([]float64, n)
				for i, o := range ord {
					out[i] = real(wd[o])
				}
				vals = &Double{V: out}
			}
			if !only {
				vz := fromNDz(v, n, n)
				if anyIm {
					out := make([]complex128, n*n)
					for j, o := range ord {
						copy(out[j*n:(j+1)*n], vz[o*n:(o+1)*n])
					}
					vecs = zmat(n, n, out, Null, Null)
				} else {
					out := make([]float64, n*n)
					for j, o := range ord {
						for i := 0; i < n; i++ {
							out[i+j*n] = real(vz[i+o*n])
						}
					}
					vecs = dmat(n, n, out, Null, Null)
				}
			}
		}
		res := mkNamedList([]string{"values", "vectors"}, vals, vecs)
		if only {
			return res
		}
		return setAttr(res, "class", mkString("eigen"))
	})
	regBuiltin("print.eigen", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "...")
		if a[0] == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		s.callNamed(c, "cat", []Value{mkString("eigen() decomposition\n")}, []string{""})
		u := s.callNamed(c, "unclass", []Value{a[0]}, []string{""})
		s.callNamed(c, "print", []Value{u}, []string{""})
		s.visible = false
		return a[0]
	})
}

// hermitian is isSymmetric for a complex matrix: the mean modulus of x - conj(t(x))
// against the mean modulus of x, within 100 machine epsilons, as for a real one.
func hermitian(n int, z []complex128) bool {
	var diff, size float64
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			diff += cmplx.Abs(z[i+j*n] - cmplx.Conj(z[j+i*n]))
			size += cmplx.Abs(z[i+j*n])
		}
	}
	if size == 0 {
		return true
	}
	return diff/size < 100*math.Pow(2, -52)
}
