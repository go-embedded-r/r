x <- c(-3, -0.5, 0, 0.7, 2, 8); signif(dnorm(x), 12); signif(pnorm(x), 12); signif(dnorm(x, 1, 2, log = TRUE), 12)
#---
signif(pnorm(c(-40, -10, 5, 10), lower.tail = FALSE), 12); signif(pnorm(c(-40, -10, 5, 10), log.p = TRUE), 12); signif(pnorm(c(-40, -10, 5, 10), lower.tail = FALSE, log.p = TRUE), 12)
#---
p <- c(0, 1e-300, 1e-10, 0.001, 0.1, 0.5, 0.9, 0.999, 1); signif(qnorm(p), 12); signif(qnorm(p, 3, 0.5, lower.tail = FALSE), 12)
#---
lp <- c(-Inf, -1e5, -700, -30, -2, -0.7, -1e-5, 0); signif(qnorm(lp, log.p = TRUE), 12); signif(qnorm(lp, log.p = TRUE, lower.tail = FALSE), 12)
#---
x <- c(-1, 0, 0.25, 0.5, 1, 2.6, 4); signif(dunif(x, 0, 2.75), 12); signif(punif(x, 0, 2.75), 12); signif(punif(x, 0, 2.75, lower.tail = FALSE, log.p = TRUE), 12); signif(qunif(c(0, 0.1, 0.5, 1), -1, 3), 12); dunif(x, 0, 2.75, log = TRUE)
#---
x <- c(-1, 0, 0.1, 1, 5, 30); signif(dexp(x, 0.5), 12); signif(pexp(x, 3), 12); signif(pexp(x, 3, lower.tail = FALSE, log.p = TRUE), 12); signif(qexp(c(0, 0.001, 0.5, 0.999, 1), 2), 12); signif(qexp(c(-1e-5, -2, -30), 0.7, log.p = TRUE), 12)
#---
x <- c(0, 0.1, 1, 3, 10); signif(dgamma(x, 2, 3), 12); signif(pgamma(x, 7.5, scale = 2), 12); signif(pgamma(x, 0.5, lower.tail = FALSE), 12); signif(qgamma(c(0.001, 0.5, 0.999), 2.5, rate = 4), 12); signif(dgamma(x, 0.5, log = TRUE), 12)
#---
x <- c(0, 0.001, 0.1, 0.5, 0.9, 0.999, 1); signif(dbeta(x, 2, 3), 12); signif(pbeta(x, 0.5, 2.5), 12); signif(pbeta(x, 7.5, 100, lower.tail = FALSE), 12); signif(qbeta(c(0.001, 0.5, 0.999), 2, 5), 12); signif(dbeta(x, 2, 3, log = TRUE), 12)
#---
x <- c(-30, -3, -0.5, 0, 0.7, 2, 10); signif(dt(x, 5), 12); signif(pt(x, 2.5), 12); signif(pt(x, 30, lower.tail = FALSE), 12); signif(qt(c(0.001, 0.1, 0.5, 0.9, 0.999), 7), 12); signif(dt(x, 1000, log = TRUE), 12)
#---
x <- c(0, 0.01, 0.5, 1, 3, 20); signif(df(x, 3, 10), 12); signif(pf(x, 5, 30), 12); signif(pf(x, 2, 7, lower.tail = FALSE), 12); signif(qf(c(0.001, 0.9, 0.999), 4, 9), 12); signif(df(x, 2.5, 100, log = TRUE), 12)
#---
x <- c(0, 0.01, 0.5, 1, 3, 20); signif(dchisq(x, 3), 12); signif(pchisq(x, 5), 12); signif(pchisq(x, 2, lower.tail = FALSE), 12); signif(qchisq(c(0.001, 0.5, 0.999), 4), 12); signif(dchisq(x, 0.5, log = TRUE), 12)
#---
x <- c(0, 0.01, 0.5, 1, 3, 20); signif(dlnorm(x, 0, 1), 12); signif(plnorm(x, -2, 0.5), 12); signif(plnorm(x, 1, 2, lower.tail = FALSE, log.p = TRUE), 12); signif(qlnorm(c(0, 0.1, 0.5, 0.9, 1), 1, 0.1), 12); signif(qlnorm(c(-30, -2), 0, 1, log.p = TRUE), 12)
#---
x <- c(-800, -30, -3, 0, 0.7, 2, 30); signif(dlogis(x, 0, 1), 12); signif(plogis(x, 1, 0.5), 12); signif(plogis(x, 0, 1, log.p = TRUE), 12); signif(qlogis(c(0, 0.1, 0.5, 0.9, 1), 2, 7), 12); signif(qlogis(c(-30, -1e-5), log.p = TRUE), 12)
#---
x <- c(-1e10, -30, -3, 0, 0.7, 2, 30, 1e10); signif(dcauchy(x), 12); signif(pcauchy(x, 1, 0.5), 12); signif(pcauchy(x, 0, 1, lower.tail = FALSE, log.p = TRUE), 12); signif(qcauchy(c(0, 0.1, 0.5, 0.9, 1), 2, 7), 12); signif(qcauchy(c(-30, -1e-5), log.p = TRUE), 12)
#---
x <- c(0, 0.01, 0.5, 1, 3, 20); signif(dweibull(x, 2, 3), 12); signif(pweibull(x, 0.5, 3), 12); signif(pweibull(x, 7.5, 2, lower.tail = FALSE, log.p = TRUE), 12); signif(qweibull(c(0, 0.1, 0.5, 0.9, 1), 2, 3), 12); dweibull(x, 1, 1, log = TRUE)
#---
r <- pnorm(matrix(c(-1, 0, 1, 2), 2, dimnames = list(c("a", "b"), c("u", "v")))); r; attributes(pnorm(structure(c(a = 1, b = 2), foo = "bar"))); dnorm(numeric(0)); pnorm(numeric(0), 1:3); length(qnorm(c(0.5, 0.9), 1:3, 1:4))
#---
pnorm(c(1, NA, NaN, 2)); dnorm(NA); qnorm(c(NA, NaN, 0.5)); pnorm(1, NA); pnorm(NaN, NA); pnorm(NA, NaN); pexp(NaN, NA)
#---
signif(dnorm(c(1, 2), mean = c(0, 1, 2, 3), sd = c(1, 2)), 12); signif(pgamma(1:6, shape = c(1, 2), scale = c(1, 2, 3)), 12)
#---
signif(pt(1, 5, ncp = 0), 12); signif(pchisq(3, 4, ncp = 0), 12); signif(pchisq(3, 4, ncp = 1.5), 10)
#---
signif(dgamma(1, 2, rate = 3, scale = 1 / 3), 12); tryCatch(dgamma(1, 2, rate = 3, scale = 1), error = function(e) conditionMessage(e)); signif(dgamma(1, 2, 3), 12)
#---
tryCatch(dnorm(), error = function(e) conditionMessage(e)); tryCatch(pnorm(1, 2, 3, 4, 5, 6), error = function(e) conditionMessage(e)); tryCatch(pnorm("a"), error = function(e) conditionMessage(e)); tryCatch(pt(1), error = function(e) conditionMessage(e))
#---
signif(pnorm(1:3, lower.tail = NA), 12); signif(pnorm(1:3, lower.tail = 0), 12); signif(dnorm(1, log = 1), 12)
