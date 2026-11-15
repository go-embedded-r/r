package main

import (
	"fmt"
	"math/rand"
	"strings"
)

// The slice20 mode: random programs over the linear algebra of slice 20 (solve,
// det, determinant, chol, chol2inv, backsolve, forwardsolve, qr and its
// family, svd, eigen, norm, rcond, kappa) and the contrasts and model.matrix.
// The first line is "# family". The matrices are small, with integer entries
// and a dominant diagonal, so that no program sits on the edge of singularity,
// where two correct algorithms may disagree about which message to give. A
// value ends in "# libm": it is compared within --ulp units in the last place,
// after being rounded to seven digits, because the kernels are accurate to a
// stated bound and not bit-identical to LAPACK; the singular vectors and the
// eigenvectors are compared through what does not depend on their sign.

func fuzzMatrix20(r *rand.Rand, n int, dominant bool) string {
	v := make([]string, n*n)
	for i := range v {
		x := float64(r.Intn(7)-3) + float64(r.Intn(97))/97 // no exact ties between pivots, none between rounded digits
		if dominant && i%(n+1) == 0 {
			x += 9
		}
		v[i] = fmt.Sprint(x)
	}
	return "matrix(c(" + strings.Join(v, ", ") + "), " + fmt.Sprint(n) + ")"
}

func fuzzVec20(r *rand.Rand, n int) string {
	v := make([]string, n)
	for i := range v {
		v[i] = fmt.Sprint(r.Intn(11) - 5)
	}
	return "c(" + strings.Join(v, ", ") + ")"
}

func fuzzSlice20(r *rand.Rand) string {
	pick := func(p []string) string { return p[r.Intn(len(p))] }
	n := 2 + r.Intn(4)
	m := n + r.Intn(3)
	rd := func(e string) string { return "r7 <- function(z) round(z, 7) + 0; r7(" + e + ")" }
	var family, body string
	switch r.Intn(10) {
	case 0:
		family = "solve"
		a := fuzzMatrix20(r, n, true)
		body = "a <- " + a + "; " + pick([]string{
			rd("solve(a)"), rd("solve(a, " + fuzzVec20(r, n) + ")"), rd("solve(a, matrix(" + fuzzVec20(r, 2*n) + ", " + fmt.Sprint(n) + "))"),
			rd("solve(qr(a))"), rd("qr.solve(a, " + fuzzVec20(r, n) + ")"), "tryCatch(solve(a, " + fuzzVec20(r, n+1) + "), error = function(e) conditionMessage(e))",
			"tryCatch(solve(a[, -1]), error = function(e) conditionMessage(e))",
		}) + " # libm"
	case 1:
		family = "det"
		a := fuzzMatrix20(r, n, r.Intn(2) == 0)
		body = "a <- " + a + "; " + pick([]string{
			"x <- determinant(a); c(round(as.vector(x$modulus), 6), x$sign)", "round(det(a), 4)", "determinant(a, logarithm = FALSE)$sign",
		}) + " # libm"
	case 2:
		family = "chol"
		a := fuzzMatrix20(r, m, true)
		body = "a <- crossprod(" + a + "); " + pick([]string{
			rd("chol(a)"), rd("chol2inv(chol(a))"), rd("backsolve(chol(a), " + fuzzVec20(r, m) + ")"),
			rd("forwardsolve(t(chol(a)), " + fuzzVec20(r, m) + ")"), rd("backsolve(chol(a), " + fuzzVec20(r, m) + ", transpose = TRUE)"),
			rd("c(chol(a, pivot = TRUE))"), "attr(chol(a, pivot = TRUE), 'rank')",
			rd("backsolve(chol(a), matrix(" + fuzzVec20(r, 2*m) + ", " + fmt.Sprint(m) + "), k = " + fmt.Sprint(1+r.Intn(m)) + ")"),
		}) + " # libm"
	case 3:
		family = "qr"
		x := "matrix(c(" + strings.TrimSuffix(strings.TrimPrefix(fuzzVec20(r, m*n), "c("), ")") + "), " + fmt.Sprint(m) + ")"
		body = "x <- " + x + "; q <- qr(x" + pick([]string{"", ", LAPACK = TRUE", ", tol = 1e-3"}) + "); " + pick([]string{
			"q$rank", "q$pivot", rd("abs(qr.R(q))"), rd("abs(qr.Q(q))"), rd("qr.X(q)"), rd("qr.coef(q, " + fuzzVec20(r, m) + ")"),
			rd("qr.qty(q, " + fuzzVec20(r, m) + ")"), rd("qr.qy(q, " + fuzzVec20(r, m) + ")"),
			rd("qr.Q(q, complete = TRUE)[, 1:" + fmt.Sprint(m) + "] %*% t(qr.Q(q, complete = TRUE)[, 1:" + fmt.Sprint(m) + "])"),
		}) + " # libm"
	case 4:
		family = "svd"
		x := "matrix(c(" + strings.TrimSuffix(strings.TrimPrefix(fuzzVec20(r, m*n), "c("), ")") + "), " + fmt.Sprint(m) + ")"
		body = "x <- " + x + "; s <- svd(x" + pick([]string{"", ", nu = 0, nv = 0", ", nu = " + fmt.Sprint(m)}) + "); " + pick([]string{
			rd("s$d"), "dim(s$u)", "dim(s$v)",
			rd("La.svd(x, nu = 0, nv = 0)$d"), rd("crossprod(svd(x)$u)"),
		}) + " # libm"
	case 5:
		family = "eigen"
		a := fuzzMatrix20(r, n, false)
		body = "a <- crossprod(" + a + "); g <- " + fuzzMatrix20(r, n, true) + "; " + pick([]string{
			rd("eigen(a)$values"), rd("eigen(a, only.values = TRUE)$values"), rd("abs(eigen(a)$vectors %*% diag(eigen(a)$values, " + fmt.Sprint(n) + ") %*% t(eigen(a)$vectors))"),
			rd("Re(eigen(g, only.values = TRUE)$values)"), "class(eigen(a))", "names(eigen(g))",
			"e <- eigen(g); " + rd("Mod(g %*% e$vectors - e$vectors %*% diag(e$values, "+fmt.Sprint(n)+")) < 1e-9"),
		}) + " # libm"
	case 6:
		family = "norm"
		a := fuzzMatrix20(r, n, r.Intn(2) == 0)
		body = "a <- " + a + "; " + pick([]string{
			rd(`norm(a, "O")`), rd(`norm(a, "I")`), rd(`norm(a, "F")`), rd(`norm(a, "M")`), rd(`norm(a, "2")`),
			"tryCatch(norm(a, 'x'), error = function(e) conditionMessage(e))",
			rd("rcond(a)"), rd(`rcond(a, "I")`), rd("1 / kappa(a, exact = TRUE)"),
		}) + " # libm"
	case 7:
		family = "contr"
		k := 2 + r.Intn(5)
		body = pick([]string{
			"contr.treatment(" + fmt.Sprint(k) + ", base = " + fmt.Sprint(1+r.Intn(k)) + ")", "contr.sum(" + fmt.Sprint(k) + ")", "contr.helmert(" + fmt.Sprint(k) + ")",
			"contr.SAS(" + fmt.Sprint(k) + ")", "round(contr.poly(" + fmt.Sprint(k) + "), 7) + 0", "contr.treatment(" + fmt.Sprint(k) + ", contrasts = FALSE)",
			"contr.sum(letters[1:" + fmt.Sprint(k) + "])", "round(crossprod(contr.poly(" + fmt.Sprint(k) + ")), 7) + 0",
		})
	case 8:
		family = "contrasts"
		k := 2 + r.Intn(4)
		body = "f <- factor(letters[1:" + fmt.Sprint(k) + "]); " + pick([]string{
			"contrasts(f)", "contrasts(f, contrasts = FALSE)", "contrasts(f) <- 'contr.sum'; contrasts(f)", "contrasts(f) <- contr.helmert; attributes(f)",
			"contrasts(f, 1) <- contr.helmert(" + fmt.Sprint(k) + ")[, 1, drop = FALSE]; contrasts(f)", "o <- ordered(f); round(contrasts(o), 7) + 0",
		})
	default:
		family = "modelmatrix"
		nr := 4 + r.Intn(5)
		fv, gv, xv := make([]string, nr), make([]string, nr), make([]string, nr)
		for i := range fv {
			fv[i] = fmt.Sprintf("%q", []string{"a", "b", "c"}[r.Intn(3)])
			gv[i] = fmt.Sprintf("%q", []string{"u", "v"}[r.Intn(2)])
			xv[i] = fmt.Sprint(r.Intn(20)) + ".5"
		}
		body = "d <- data.frame(f = c(" + strings.Join(fv, ", ") + "), g = c(" + strings.Join(gv, ", ") + "), x = c(" + strings.Join(xv, ", ") + "), y = " + fuzzVec20(r, nr) + "); " + pick([]string{
			"model.matrix(~ f, d)", "model.matrix(y ~ f + g, d)", "model.matrix(~ f * g, d)", "model.matrix(~ f:g, d)", "model.matrix(~ 0 + f, d)",
			"model.matrix(~ x + f, d)", "model.matrix(~ x:f, d)", "model.matrix(~ f, d, contrasts.arg = list(f = 'contr.sum'))",
			"model.matrix(~ f / x, d)", "model.matrix(~ f + I(x^2), d)", "model.matrix(~ g - 1 + f, d)", "attr(model.matrix(~ f + g, d), 'assign')",
			"dim(model.matrix(~ f * g * x, d))",
		})
	}
	return "# " + family + "\n" + body
}
