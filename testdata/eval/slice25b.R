print(signif(gamma(c(0, -1, -2.5, 171, -0.5, 0.5, 1.5, 5, 10, 20.5, 100.5, 1e-5, 3.25, -7.75)), 12)); gamma(c(Inf, -Inf, NA, NaN)); gamma(1:3); gamma(TRUE)
#---
print(signif(gamma(c(-1e-300, 1e-300, 2^-30, 30.25, 150.5)), 12)); gamma(c(a = 1, b = 5)); gamma(structure(c(1, 2), dim = 2L)); gamma(numeric(0))
#---
print(signif(lgamma(c(-0.5, 0.5, 1.5, 5, 10, 20.5, 100.5, 1e10, -170.5, 1e-5, 3.25, -7.75)), 12)); lgamma(c(0, -1, Inf, -Inf, NA)); print(signif(lgamma(2.5e305), 12)); lgamma(3e305)
#---
print(signif(digamma(c(0.5, 1, 2, 10.5, 100, -0.5, -2.25, 1e-5, 1e5, 30.5)), 12)); digamma(c(0, -1, Inf, -Inf, NA)); print(signif(digamma(1e300), 12))
#---
print(signif(trigamma(c(0.5, 1, 2, 10.5, 100, -0.5, -2.25, 1e-5, 1e5, 30.5)), 12)); trigamma(c(0, -1, Inf, -Inf, NA))
#---
print(signif(psigamma(c(0.5, 2.5, 10, 0.1, 30.5), 2), 12)); print(signif(psigamma(c(0.5, 2.5, 10), 3:5), 12)); print(signif(psigamma(c(-0.5, -2.5, -7.25), 2), 12)); print(signif(psigamma(-0.5, 2:5), 12))
#---
psigamma(1, c(-1, 0.4, 0.6, 100, 101, 1e9, NA, NaN)); psigamma(c(0, -1, Inf, -Inf), 3); psigamma(-1, 0:5); signif(psigamma(1), 12); signif(psigamma(1, deriv = 2L), 12)
#---
signif(psigamma(1:3, c(0, 1)), 12); signif(psigamma(c(a = 1, b = 2), 1), 12); signif(psigamma(matrix(1:4, 2), 1), 12); psigamma(numeric(0), 1); tryCatch(psigamma("a"), error = function(e) conditionMessage(e)); tryCatch(psigamma(), error = function(e) conditionMessage(e))
#---
print(signif(beta(c(0.5, 1, 2, 3.5, 100, 90, 1e-200, 10.5), c(0.5, 1, 3, 2.25, 80, 90, 1e-200, 1e-3)), 12)); beta(c(-1, 0, 1, Inf, NA), c(2, 2, 0, 2, 1)); beta(1e300, 1e300); signif(beta(c(a = 1, b = 2), 3), 12)
#---
print(signif(lbeta(c(0.5, 1, 2, 3.5, 100, 90, 1e-200, 10.5, 1e300), c(0.5, 1, 3, 2.25, 80, 90, 1e-200, 1e-3, 1e300)), 12)); lbeta(c(-1, 0, 1, Inf, NA), c(2, 2, 0, 2, 1)); lbeta(numeric(0), 1)
#---
print(signif(besselJ(c(0.5, 1, 2, 10, 50, 100.5, 1e4), 0), 12)); print(signif(besselJ(1, c(0.5, 1, 2, 7.25, 30)), 12)); print(signif(besselJ(c(0.1, 3, 40), c(-0.5, -1.5, -2.25)), 12)); signif(besselJ(c(0, 1, 10, -1, Inf, NA), 0), 12)
#---
signif(besselJ(1, c(0.5, -0.5, 1e5, NA)), 12); signif(besselJ(2, c(1000, 1e5)), 12); signif(besselJ(c(1e7, 5e7, 1e30), 0), 12)
#---
print(signif(besselY(c(0.5, 1, 2, 10, 50, 100.5, 1e4), 0), 12)); print(signif(besselY(1, c(0.5, 1, 2, 7.25, 30)), 12)); print(signif(besselY(c(0.1, 3, 40), c(-0.5, -1.5, -2.25)), 12)); signif(besselY(c(0, 1, 10, -1, Inf, NA), 0), 12)
#---
signif(besselY(1, c(0.5, -0.5, 1e8, -1e8)), 12); signif(besselY(2, c(1000, 1e5)), 12); signif(besselY(c(1e7, 3e8), 0), 12); besselY(0, c(0, 1, 2.5))
#---
print(signif(besselI(c(0, 0.5, 1, 2, 10, 50, 100.5), 0), 12)); print(signif(besselI(1, c(0.5, 1, 2, 7.25, 30)), 12)); print(signif(besselI(c(0.1, 3, 40), c(-0.5, -1.5, -2.25)), 12)); signif(besselI(c(0, 1, 10, -1, Inf, NA, 800), 0), 12)
#---
print(signif(besselI(c(0, 1, 10, 800, 1e4), 0, TRUE), 12)); print(signif(besselI(c(0.1, 3, 40), c(-0.5, -1.5, -2.25), expon.scaled = TRUE), 12)); signif(besselI(c(0, 1, 10, -1, Inf, NA, 800), 0, TRUE), 12); signif(besselI(1, 1, expon.scaled = 2), 12); besselI(1, 1, expon.scaled = NA); signif(besselI(1, 1, expon.scaled = c(TRUE, FALSE)), 12); besselI(1, c(1e5, -1e5)); besselI(1e6, 0, TRUE)
#---
print(signif(besselK(c(0.5, 1, 2, 10, 50, 100.5), 0), 12)); print(signif(besselK(1, c(0.5, 1, 2, 7.25, 30)), 12)); print(signif(besselK(c(0.1, 3, 40), c(-0.5, -1.5, -2.25)), 12)); signif(besselK(c(0, 1, 10, -1, Inf, NA, 800), 0), 12)
#---
print(signif(besselK(c(0, 1, 10, 800, 1e-300), 0, TRUE), 12)); signif(besselK(1, c(0.5, -0.5, 1e8, -1e8)), 12); besselK(0, 0:2); besselK(2, c(1000, 1e6))
#---
signif(besselI(1:3, c(0, 1)), 12); signif(besselJ(1:3, 0:1), 12); besselJ(numeric(0), 1); besselJ(1, numeric(0)); signif(besselJ(nu = 1, x = 2), 12); signif(besselJ(matrix(1:4, 2), 1), 12); signif(besselJ(c(a = 1, b = 2), 1), 12); signif(besselJ(1, c(a = 1, b = 2)), 12)
#---
tryCatch(besselJ("a", 1), error = function(e) conditionMessage(e)); tryCatch(besselJ(1), error = function(e) conditionMessage(e)); besselK(1, 1, "a"); tryCatch(beta(1), error = function(e) conditionMessage(e)); tryCatch(gamma("a"), error = function(e) conditionMessage(e)); tryCatch(gamma(1i), error = function(e) conditionMessage(e))
#---
signif(exp(c(1, -1, 709.78, 0.5, 20, -20, 1e-10, 100, 36.7, 0.1, 2.5)), 15); exp(c(709.79, -745.2, -745.14, -1e5, 1e-320, Inf, -Inf, NA)); sprintf("%.15g", exp(c(-708.5, -740, -744.4)))
#---
signif(2^c(0.5, -0.5, 10, 1023.9, 0.1, 3.7, -10.25), 15); 2^c(1024, -1074, -1075, -1074.5); signif(as.vector(sapply(c(1.1, 10, 0.3, 7), function(x) x^seq(-2, 2, by = 0.5))), 15)
#---
(-8)^(1/3); (-8)^3; (-8)^c(1e15, 1e16, 1e300, 0.5, -3); 10^(-5:5); signif(c(3, 5, 7)^c(40, 30, 20), 15); signif(1.0000001^1e7, 15); signif(0.9999999^1e9, 15)
#---
w <- function(e) { m <- character(); v <- withCallingHandlers(e, warning = function(c) { m <<- c(m, conditionMessage(c)); invokeRestart("muffleWarning") }); list(signif(v, 6), m) }; w(besselY(1, 2e7)); w(besselJ(-1, 1)); w(besselY(-1, 1)); w(besselI(-1, 1)); w(besselK(-1, 1)); w(psigamma(1, 101)); w(besselJ(1e25, 0)); w(gamma(-1)); w(gamma(200)); w(lgamma(c(-1, 1e306))); w(digamma(-1)); w(trigamma(0)); w(psigamma(0, 2)); w(beta(-1, 2)); w(lbeta(-1, 2)); w(exp(1000)); w(2^2000); w(psigamma(1e-300, 2))
