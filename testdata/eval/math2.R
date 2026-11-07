log(1); log(0); log(-1); log(8, 2); log(100, 10); log(c(1, NA)); log(1L); log(8, base = 2)
#---
log("a")
#---
log()
#---
log10(1000); log10(c(1, 10, 100)); log2(8); log2(c(1, 2, 4)); log1p(0); expm1(0)
#---
sin(0); cos(0); tan(0); asin(0); acos(1); atan(0); sinh(0); cosh(0); tanh(0); asinh(0); acosh(1); atanh(0)
#---
asin(2); acos(2); acosh(0.5); atanh(2)
#---
sin(c(a = 0, b = 0)); cos(matrix(0, 2, 2))
#---
sinpi(0); sinpi(0.5); sinpi(1); sinpi(1.5); cospi(0); cospi(0.5); cospi(1); tanpi(0); tanpi(0.25)
#---
atan2(0, 1); atan2(1, 0) > 1.5; atan2(c(0, 0), c(1, -1)) > 1
#---
gamma(1); gamma(5); gamma(0.5) > 1.77; gamma(0); lgamma(1); lgamma(0); factorial_not_here <- 1
#---
trunc(2.7); trunc(-2.7); trunc(c(1.5, -1.5, NA)); trunc(5L); typeof(trunc(5L))
#---
sign(c(-2, 0, 3, NA)); sign(-3L); typeof(sign(3L)); sign(c(a = -1))
#---
round(2.5); round(3.5); round(-2.5); round(0.5); round(1.005, 2); round(2.567, 1); round(0.125, 2); round(1234.567, -2); round(c(1.234, 5.678), 1); round(5L); typeof(round(5L)); round(NA)
#---
round(2.5, digits = 0); round(x = 7.77, digits = 1); round(123.456, 10); round(-0.55, 1); round(0.15, 1); round(0.25, 1); round(0.35, 1)
#---
round("a")
#---
round()
#---
signif(123456, 2); signif(0.000123456, 3); signif(1234.5678); signif(c(1.5, 22.5, 333.5), 2); signif(5L, 1); signif(0); signif(-123.456, 4)
#---
cumsum(1:5); cumsum(c(1.5, 2.5, NA, 1)); cumsum(c(a = 1, b = 2)); cumsum(c(TRUE, TRUE)); cumsum(integer(0)); cumsum(NULL)
#---
cumsum(c(2147483647L, 1L))
#---
cumprod(1:5); cumprod(c(2, 3, 4)); cumprod(c(1, NA, 2)); typeof(cumprod(1:3))
#---
cummax(c(1, 3, 2, 5, 4)); cummin(c(5, 3, 4, 1, 2)); cummax(c(1L, 3L, 2L)); cummin(c(1, NA, 0)); cummax(c(1, NaN, 5)); typeof(cummax(1:3))
#---
cumsum("a")
#---
range(c(3, 1, 2)); range(1:5, 10); range(c(1, NA, 3)); range(c(1, NA, 3), na.rm = TRUE); range(c(1, Inf, -Inf, NA), finite = TRUE); range(integer(0))
#---
range("b", "a", "c")
#---
is.finite(c(1, Inf, NA, NaN)); is.infinite(c(1, Inf, -Inf, NA)); is.nan(c(1, NaN, NA)); is.finite(1:3); is.finite("a"); is.infinite(1L); is.nan(NA_real_); is.finite(c(a = 1)); is.finite(NULL)
#---
is.nan("a")
#---
is.finite(list(1))
#---
anyNA(c(1, NA)); anyNA(1:3); anyNA(list(1, NA)); anyNA(list(1, list(NA))); anyNA(list(1, list(NA)), recursive = TRUE); anyNA("a"); anyNA(NULL); anyNA(NaN)
#---
nzchar(c("a", "", NA)); nzchar(c("a", ""), keepNA = TRUE); nzchar(NA, keepNA = TRUE); nzchar(123); nzchar(c(x = "a")); nzchar(character(0))
#---
is.matrix(matrix(1:4, 2)); is.matrix(1:3); is.array(array(1:8, c(2, 2, 2))); is.array(1); is.matrix(array(1:8, c(2, 2, 2)))
#---
is.object(1); is.object(structure(1, class = "a")); is.object(factor("a")); is.recursive(list()); is.recursive(1); is.recursive(sum); is.recursive(quote(x + y)); is.language(quote(x)); is.language(1); is.pairlist(NULL); is.pairlist(1)
#---
x <- 1:3; length(x) <- 5; x
#---
x <- c(a = 1, b = 2); length(x) <- 3; x; length(x) <- 1; x
#---
l <- list(1, 2); length(l) <- 3; l
#---
x <- 1:3; length(x) <- -1
#---
x <- 1:3; storage.mode(x) <- "double"; typeof(x); storage.mode(x) <- "character"; x
#---
x <- 1:3; storage.mode(x) <- "bogus"
#---
xtfrm(c(3, 1, 2)); xtfrm(1:3)

#---
is.expression(1); is.complex(1); is.raw(1); isS4(1); interactive()

#---
round(digamma(1), 6); round(digamma(c(0.5, 2, 10)), 6); digamma(0); round(trigamma(1), 6); round(trigamma(c(0.5, 3)), 6); round(digamma(-0.5), 6); round(trigamma(-0.5), 5)
#---
seq.int(5); seq.int(2, 8); seq.int(1, 10, by = 3); seq.int(0, 1, length.out = 5); seq.int(5, 1)
#---
seq.int(1, 10, by = -1)
#---
.class2(1); .class2(1L); .class2(matrix(1:4, 2)); .class2("a"); .class2(sum); .class2(structure(1, class = "k"))
