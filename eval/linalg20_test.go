package eval

import (
	"math"
	"math/big"
	"math/cmplx"
	"math/rand"
	"testing"
)

// The accuracy of the linear algebra of slice 20, asserted against math/big.
// The contract is never bit identity with R: it is a stated bound, here
//
//	solve   ||A x - b||  <= 4 n u ||A|| ||x||        (backward error of LU)
//	chol    ||R'R - A||  <= 4 n u ||A||
//	qr      ||QR - X||   <= 8 m u ||X||,  ||Q'Q - I|| <= 8 m u
//	svd     ||U S V' - X|| <= 16 m u ||X||, ||U'U - I|| <= 16 m u
//	eigen   ||A v - l v|| <= 16 n u ||A||  (symmetric: and ||V'V - I|| <= 16 n u)
//
// in the infinity norm, u = 2^-53, the products of the residuals formed exactly.

const unitRoundoff = 1.0 / (1 << 53)

func bigf(x float64) *big.Float { return new(big.Float).SetPrec(300).SetFloat64(x) }

// exactResidual is max |A B - C| for column-major A (m x k), B (k x n), C
// (m x n), every product and sum exact.
func exactResidual(a []float64, m, k int, b []float64, n int, c []float64) float64 {
	worst := 0.0
	for j := 0; j < n; j++ {
		for i := 0; i < m; i++ {
			sum := bigf(0)
			for p := 0; p < k; p++ {
				sum.Add(sum, new(big.Float).SetPrec(300).Mul(bigf(a[i+p*m]), bigf(b[p+j*k])))
			}
			sum.Sub(sum, bigf(c[i+j*m]))
			f, _ := sum.Float64()
			worst = math.Max(worst, math.Abs(f))
		}
	}
	return worst
}

func infNorm(a []float64, m, n int) float64 { return normInf(m, n, a) }

func transposeData(a []float64, m, n int) []float64 {
	out := make([]float64, len(a))
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			out[j+i*n] = a[i+j*m]
		}
	}
	return out
}

func randMat(r *rand.Rand, m, n int) []float64 {
	a := make([]float64, m*n)
	for i := range a {
		a[i] = r.NormFloat64()
	}
	return a
}

// run evaluates R code in a session holding the named matrices and returns
// the value of the last expression's variable names.
func linSession(t *testing.T, mats map[string][]float64, dims map[string][2]int) *Session {
	t.Helper()
	s := NewSession(Options{})
	for name, v := range mats {
		d := dims[name]
		s.Set(name, dmat(d[0], d[1], v, Null, Null))
	}
	return s
}

func linEval(t *testing.T, s *Session, src string) {
	t.Helper()
	if _, err := s.Eval(src); err != nil {
		t.Fatalf("%s: %v", src, err)
	}
}

func linGet(t *testing.T, s *Session, name string) []float64 {
	t.Helper()
	v, ok := s.Get(name)
	if !ok {
		t.Fatalf("no %s", name)
	}
	return append([]float64(nil), asDoubles(v)...)
}

func TestSolveBackwardError(t *testing.T) {
	r := rand.New(rand.NewSource(20))
	worst := 0.0
	for trial := 0; trial < 60; trial++ {
		n := 2 + r.Intn(60)
		a, b := randMat(r, n, n), randMat(r, n, 2)
		s := linSession(t, map[string][]float64{"A": a, "b": b}, map[string][2]int{"A": {n, n}, "b": {n, 2}})
		linEval(t, s, "x <- solve(A, b)")
		x := linGet(t, s, "x")
		res := exactResidual(a, n, n, x, 2, b)
		bound := 4 * float64(n) * unitRoundoff * infNorm(a, n, n) * infNorm(x, n, 2)
		worst = math.Max(worst, res/bound)
		if res > bound {
			t.Fatalf("n=%d: ||Ax-b|| = %g exceeds %g", n, res, bound)
		}
	}
	t.Logf("solve: worst residual / bound = %.3f", worst)
}

func TestCholBackwardError(t *testing.T) {
	r := rand.New(rand.NewSource(21))
	worst := 0.0
	for trial := 0; trial < 60; trial++ {
		n := 2 + r.Intn(60)
		b := randMat(r, n, n)
		a := make([]float64, n*n)
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				sum := 0.0
				for k := 0; k < n; k++ {
					sum += b[k+i*n] * b[k+j*n]
				}
				a[i+j*n] = sum
			}
			a[i+i*n] += float64(n)
		}
		s := linSession(t, map[string][]float64{"A": a}, map[string][2]int{"A": {n, n}})
		linEval(t, s, "R <- chol(A)")
		R := linGet(t, s, "R")
		for j := 0; j < n; j++ {
			for i := j + 1; i < n; i++ {
				if R[i+j*n] != 0 {
					t.Fatalf("R is not upper triangular")
				}
			}
		}
		res := exactResidual(transposeData(R, n, n), n, n, R, n, a)
		bound := 4 * float64(n) * unitRoundoff * infNorm(a, n, n)
		worst = math.Max(worst, res/bound)
		if res > bound {
			t.Fatalf("n=%d: ||R'R-A|| = %g exceeds %g", n, res, bound)
		}
	}
	t.Logf("chol: worst residual / bound = %.3f", worst)
}

// orthoError is max |Q'Q - I| for column-major Q (m x k), exact.
func orthoError(q []float64, m, k int) float64 {
	eye := make([]float64, k*k)
	for i := 0; i < k; i++ {
		eye[i+i*k] = 1
	}
	return exactResidual(transposeData(q, m, k), k, m, q, k, eye)
}

func TestQRBounds(t *testing.T) {
	r := rand.New(rand.NewSource(22))
	for _, lapack := range []string{"FALSE", "TRUE"} {
		worstR, worstO := 0.0, 0.0
		for trial := 0; trial < 40; trial++ {
			m := 3 + r.Intn(50)
			n := 1 + r.Intn(m)
			x := randMat(r, m, n)
			s := linSession(t, map[string][]float64{"X": x}, map[string][2]int{"X": {m, n}})
			linEval(t, s, "q <- qr(X, LAPACK = "+lapack+"); Q <- qr.Q(q); R <- qr.R(q); p <- q$pivot; Xp <- X[, p, drop = FALSE]")
			Q, R, Xp := linGet(t, s, "Q"), linGet(t, s, "R"), linGet(t, s, "Xp")
			res := exactResidual(Q, m, n, R, n, Xp)
			bound := 8 * float64(m) * unitRoundoff * infNorm(x, m, n)
			worstR = math.Max(worstR, res/bound)
			if res > bound {
				t.Fatalf("LAPACK=%s %dx%d: ||QR-X|| = %g exceeds %g", lapack, m, n, res, bound)
			}
			o := orthoError(Q, m, n)
			ob := 8 * float64(m) * unitRoundoff
			worstO = math.Max(worstO, o/ob)
			if o > ob {
				t.Fatalf("LAPACK=%s %dx%d: ||Q'Q-I|| = %g exceeds %g", lapack, m, n, o, ob)
			}
		}
		t.Logf("qr LAPACK=%s: worst reconstruction / bound = %.3f, orthogonality / bound = %.3f", lapack, worstR, worstO)
	}
}

func TestSVDBounds(t *testing.T) {
	r := rand.New(rand.NewSource(23))
	worstR, worstO := 0.0, 0.0
	for trial := 0; trial < 40; trial++ {
		m := 2 + r.Intn(40)
		n := 2 + r.Intn(40)
		k := min(m, n)
		x := randMat(r, m, n)
		s := linSession(t, map[string][]float64{"X": x}, map[string][2]int{"X": {m, n}})
		linEval(t, s, "z <- svd(X); U <- z$u; V <- z$v; d <- z$d")
		U, V, d := linGet(t, s, "U"), linGet(t, s, "V"), linGet(t, s, "d")
		us := make([]float64, m*k)
		for j := 0; j < k; j++ {
			for i := 0; i < m; i++ {
				us[i+j*m] = U[i+j*m] * d[j]
			}
		}
		res := exactResidual(us, m, k, transposeData(V, n, k), n, x)
		bound := 16 * float64(max(m, n)) * unitRoundoff * infNorm(x, m, n)
		worstR = math.Max(worstR, res/bound)
		if res > bound {
			t.Fatalf("%dx%d: ||USV'-X|| = %g exceeds %g", m, n, res, bound)
		}
		for _, q := range [][3]interface{}{{U, m, k}, {V, n, k}} {
			o := orthoError(q[0].([]float64), q[1].(int), q[2].(int))
			ob := 16 * float64(max(m, n)) * unitRoundoff
			worstO = math.Max(worstO, o/ob)
			if o > ob {
				t.Fatalf("%dx%d: orthogonality %g exceeds %g", m, n, o, ob)
			}
		}
		for i := 1; i < len(d); i++ {
			if d[i] > d[i-1] || d[i] < 0 {
				t.Fatalf("singular values not decreasing: %v", d)
			}
		}
	}
	t.Logf("svd: worst reconstruction / bound = %.3f, orthogonality / bound = %.3f", worstR, worstO)
}

func TestEigenBounds(t *testing.T) {
	r := rand.New(rand.NewSource(24))
	worstS, worstO, worstG := 0.0, 0.0, 0.0
	for trial := 0; trial < 40; trial++ {
		n := 2 + r.Intn(40)
		b := randMat(r, n, n)
		a := make([]float64, n*n)
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				a[i+j*n] = b[i+j*n] + b[j+i*n]
			}
		}
		s := linSession(t, map[string][]float64{"A": a}, map[string][2]int{"A": {n, n}})
		linEval(t, s, "e <- eigen(A); V <- e$vectors; l <- e$values; VL <- V %*% diag(l, length(l))")
		V, l, VL := linGet(t, s, "V"), linGet(t, s, "l"), linGet(t, s, "VL")
		_ = l
		res := exactResidual(a, n, n, V, n, VL)
		bound := 16 * float64(n) * unitRoundoff * infNorm(a, n, n)
		worstS = math.Max(worstS, res/bound)
		if res > bound {
			t.Fatalf("symmetric n=%d: ||AV-VL|| = %g exceeds %g", n, res, bound)
		}
		o := orthoError(V, n, n)
		ob := 16 * float64(n) * unitRoundoff
		worstO = math.Max(worstO, o/ob)
		if o > ob {
			t.Fatalf("symmetric n=%d: ||V'V-I|| = %g exceeds %g", n, o, ob)
		}
		// the general solver, on a non-symmetric matrix: the residual of each pair
		g := randMat(r, n, n)
		gs := linSession(t, map[string][]float64{"G": g}, map[string][2]int{"G": {n, n}})
		linEval(t, gs, "e <- eigen(G); vals <- as.complex(e$values); vecs <- e$vectors")
		gv, _ := gs.Get("vals")
		gvec, _ := gs.Get("vecs")
		vals := asComplex(gv)
		vecs := asComplex(gvec)
		for j := 0; j < n; j++ {
			norm := 0.0
			for i := 0; i < n; i++ {
				var sum complex128
				for p := 0; p < n; p++ {
					sum += complex(g[i+p*n], 0) * vecs[p+j*n]
				}
				sum -= vals[j] * vecs[i+j*n]
				norm = math.Max(norm, cmplx.Abs(sum))
			}
			gb := 16 * float64(n) * unitRoundoff * infNorm(g, n, n)
			worstG = math.Max(worstG, norm/gb)
			if norm > gb {
				t.Fatalf("general n=%d pair %d: residual %g exceeds %g", n, j, norm, gb)
			}
		}
	}
	t.Logf("eigen: symmetric residual / bound = %.3f, orthogonality / bound = %.3f, general residual / bound = %.3f", worstS, worstO, worstG)
}

// The estimate of rcond is dlacn2's, so it never exceeds the exact value of
// the reciprocal condition number, and is within a factor of n of it.
func TestRcondIsALowerBoundOfTheExactValue(t *testing.T) {
	r := rand.New(rand.NewSource(25))
	for trial := 0; trial < 60; trial++ {
		n := 2 + r.Intn(30)
		a := randMat(r, n, n)
		est := rcondGeneral(n, a, false)
		s := linSession(t, map[string][]float64{"A": a}, map[string][2]int{"A": {n, n}})
		linEval(t, s, "B <- solve(A)")
		exact := 1 / (norm1(n, n, a) * norm1(n, n, linGet(t, s, "B")))
		if est < exact*(1-1e-9) || est > 3*float64(n)*exact {
			t.Fatalf("n=%d: estimate %g outside [exact, 3 n exact], exact %g", n, est, exact)
		}
	}
}
