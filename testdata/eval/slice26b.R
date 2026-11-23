x <- c(-1, 0, 1, 2); signif(dnorm(x, 0, 0), 6); signif(pnorm(x, 0, 0), 6); signif(pnorm(x, 0, 0, lower.tail = FALSE, log.p = TRUE), 6); signif(qnorm(c(0, 0.3, 1), 5, 0), 6); signif(dnorm(x, 0, 0, log = TRUE), 6)
#---
signif(dnorm(x, Inf), 6); signif(pnorm(x, -Inf), 6); signif(pnorm(x, 0, Inf), 6); signif(qnorm(0.3, Inf), 6); signif(dnorm(x, 0, -1), 6); signif(qnorm(0.3, 0, -1), 6); signif(qnorm(c(0, 1), 0, -1), 6)
#---
signif(qnorm(c(-1e5, -700, -0.7), mean = c(-2.5, 0), sd = c(0.01, 0.01, 0), log.p = TRUE), 6); signif(qnorm(-1e5, 0.3, 0, log.p = TRUE), 6)
#---
signif(qnorm(c(-0.5, 1.5, NA, NaN)), 6); signif(qnorm(c(-0.5, 1.5), log.p = TRUE), 6); signif(qnorm(0.5, NA), 6)
#---
x <- c(-1, 0, 1, 2, 3); signif(dunif(x, 1, 1), 6); signif(punif(x, 1, 1), 6); signif(qunif(c(0, 0.5, 1), 1, 1), 6); signif(dunif(x, 2, 1), 6); signif(punif(x, 2, 1), 6); signif(qunif(c(0, 0.5, 1), 2, 1), 6)
#---
x <- c(-1, 0, 1, 2, 3); signif(dunif(x, 0, Inf), 6); signif(punif(x, 0, Inf), 6); signif(punif(x, -Inf, 2), 6); signif(qunif(c(0, 0.5, 1), 0, Inf), 6); signif(qunif(0.5, NA, 1), 6); signif(punif(x, NA), 6)
#---
x <- c(-1, 0, 1, 5, Inf); signif(dexp(x, 0), 6); signif(pexp(x, 0), 6); signif(qexp(c(0, 0.5, 1), 0), 6); signif(dexp(x, Inf), 6); signif(pexp(x, Inf), 6); signif(qexp(c(0, 0.5, 1), Inf), 6); signif(dexp(x, -1), 6); signif(pexp(x, -1), 6); signif(qexp(c(0, 0.5, 1), -1), 6)
#---
x <- c(-1, 0, 1, 5, Inf); signif(dexp(x, Inf, log = TRUE), 6); signif(pexp(x, 0, log.p = TRUE), 6); signif(pexp(x, 0, lower.tail = FALSE), 6); signif(pexp(x, Inf, lower.tail = FALSE, log.p = TRUE), 6)
#---
x <- c(-1, 0, 1, 5, Inf); signif(dgamma(x, 0), 6); signif(pgamma(x, 0), 6); signif(qgamma(c(0, 0.5, 1), 0), 6); signif(dgamma(x, Inf), 6); signif(pgamma(x, Inf), 6); signif(qgamma(c(0, 0.5, 1), Inf), 6); signif(dgamma(x, -1), 6); signif(pgamma(x, -1), 6); signif(qgamma(c(0, 0.5, 1), -1), 6)
#---
x <- c(0, 1, 5, Inf); signif(dgamma(x, 2, scale = 0), 6); signif(pgamma(x, 2, scale = 0), 6); signif(dgamma(x, 2, scale = Inf), 6); signif(pgamma(x, 2, scale = Inf), 6); signif(qgamma(c(0, 0.5, 1), 2, scale = Inf), 6); signif(dgamma(x, 2, rate = 0), 6); signif(dgamma(x, 2, rate = Inf), 6)
#---
x <- c(0, 1, Inf); signif(dgamma(x, 0.01, 0, log = TRUE), 6); signif(dgamma(x, 1, scale = Inf, log = TRUE), 6); signif(dgamma(x, 0, log = TRUE), 6); signif(dgamma(x, Inf, log = TRUE), 6); signif(pgamma(x, 0, log.p = TRUE), 6); signif(pgamma(x, Inf, lower.tail = FALSE, log.p = TRUE), 6)
#---
x <- c(-0.5, 0, 0.5, 1, 1.5); signif(dbeta(x, 0, 2), 6); signif(pbeta(x, 0, 2), 6); signif(qbeta(c(0, 0.5, 1), 0, 2), 6); signif(dbeta(x, 2, 0), 6); signif(pbeta(x, 2, 0), 6); signif(qbeta(c(0, 0.5, 1), 2, 0), 6); signif(dbeta(x, 0, 0), 6); signif(pbeta(x, 0, 0), 6); signif(qbeta(c(0, 0.3, 0.5, 0.7, 1), 0, 0), 6)
#---
x <- c(-0.5, 0, 0.5, 1, 1.5); signif(dbeta(x, Inf, 2), 6); signif(pbeta(x, Inf, 2), 6); signif(qbeta(c(0, 0.5, 1), Inf, 2), 6); signif(dbeta(x, 2, Inf), 6); signif(pbeta(x, 2, Inf), 6); signif(qbeta(c(0, 0.5, 1), 2, Inf), 6); signif(dbeta(x, Inf, Inf), 6); signif(pbeta(x, Inf, Inf), 6); signif(qbeta(c(0, 0.5, 1), Inf, Inf), 6); signif(pbeta(x, -1, 2), 6); signif(qbeta(0.5, 2, -1), 6)
#---
x <- c(-0.5, 0, 0.5, 1, 1.5); signif(pbeta(x, 0, 2, lower.tail = FALSE, log.p = TRUE), 6); signif(pbeta(x, Inf, Inf, log.p = TRUE), 6); signif(dbeta(x, 0, 0, log = TRUE), 6); signif(dbeta(x, Inf, 2, log = TRUE), 6)
#---
x <- c(-1, 0, 1, 5); signif(dt(x, 0), 6); signif(pt(x, 0), 6); signif(qt(c(0, 0.3, 0.5, 0.7, 1), 0), 6); signif(dt(x, Inf), 6); signif(pt(x, Inf), 6); signif(qt(c(0, 0.3, 0.9, 1), Inf), 6); signif(dt(x, -1), 6); signif(pt(x, -1), 6); signif(qt(c(0, 0.5, 1), -1), 6)
#---
x <- c(-1, 0, 1, 5); signif(dt(x, Inf, log = TRUE), 6); signif(pt(x, Inf, lower.tail = FALSE, log.p = TRUE), 6); signif(pt(x, 0, log.p = TRUE), 6); signif(qt(c(-1e-5, -2), Inf, log.p = TRUE), 6)
#---
x <- c(-1, 0, 1, 5, Inf); signif(df(x, 0, 3), 6); signif(pf(x, 0, 3), 6); signif(qf(c(0, 0.5, 1), 0, 3), 6); signif(df(x, 3, 0), 6); signif(pf(x, 3, 0), 6); signif(qf(c(0, 0.5, 1), 3, 0), 6); signif(df(x, 0, 0), 6); signif(pf(x, 0, 0), 6); signif(qf(c(0, 0.5, 1), 0, 0), 6)
#---
x <- c(-1, 0, 1, 5, Inf); signif(df(x, Inf, 3), 6); signif(pf(x, Inf, 3), 6); signif(qf(c(0, 0.5, 1), Inf, 3), 6); signif(df(x, 3, Inf), 6); signif(pf(x, 3, Inf), 6); signif(qf(c(0, 0.5, 1), 3, Inf), 6); signif(df(x, Inf, Inf), 6); signif(pf(x, Inf, Inf), 6); signif(qf(c(0, 0.5, 1), Inf, Inf), 6); signif(pf(x, -1, 2), 6)
#---
x <- c(-1, 0, 0.5, 1, 5, Inf); signif(dchisq(x, 0), 6); signif(pchisq(x, 0), 6); signif(qchisq(c(0, 0.5, 1), 0), 6); signif(pchisq(x, 0, lower.tail = FALSE), 6); signif(dchisq(x, -1), 6); signif(pchisq(x, -1), 6); signif(qchisq(c(0, 0.5, 1), -1), 6); signif(pchisq(x, 0, log.p = TRUE), 6); signif(qchisq(c(0, 1), NA), 6)
#---
x <- c(-1, 0, 0.5, 1, 5, Inf); signif(dlnorm(x, 0, 0), 6); signif(plnorm(x, 0, 0), 6); signif(qlnorm(c(0, 0.5, 1), 0, 0), 6); signif(plnorm(x, 0, 0, lower.tail = FALSE, log.p = TRUE), 6); signif(dlnorm(x, 0, -1), 6); signif(plnorm(x, 0, -1), 6); signif(qlnorm(c(0, 0.5, 1), 0, -1), 6); signif(qlnorm(-1e5, 0, 0, log.p = TRUE), 6)
#---
x <- c(-1, 0, 0.5, 1, 5, Inf); signif(dlnorm(x, -Inf, 1), 6); signif(plnorm(x, Inf, 1), 6); signif(dlnorm(x, 0, Inf), 6); signif(plnorm(x, 0, Inf), 6); signif(qlnorm(c(0.2, 0.7), 0, Inf), 6)
#---
x <- c(-Inf, -1, 0, 1, Inf); signif(dlogis(x, 0, 0), 6); signif(plogis(x, 0, 0), 6); signif(qlogis(c(0, 0.5, 1), 2, 0), 6); signif(dlogis(x, 0, -1), 6); signif(plogis(x, 0, -1), 6); signif(qlogis(c(0, 0.5, 1), 0, -1), 6); signif(qlogis(c(-30, -1e-5), 1, 0, log.p = TRUE), 6)
#---
x <- c(-Inf, -1, 0, 1, Inf); signif(dcauchy(x, 0, 0), 6); signif(pcauchy(x, 0, 0), 6); signif(qcauchy(c(0, 0.5, 1), 2, 0), 6); signif(dcauchy(x, 0, -1), 6); signif(pcauchy(x, 0, -1), 6); signif(qcauchy(c(0, 0.5, 1), 0, -1), 6); signif(qcauchy(c(-30, -1e-5), 1, 0, log.p = TRUE), 6)
#---
x <- c(-1, 0, 1, 5, Inf); signif(dweibull(x, 0, 1), 6); signif(pweibull(x, 0, 1), 6); signif(qweibull(c(0, 0.5, 1), 0, 1), 6); signif(dweibull(x, 1, 0), 6); signif(pweibull(x, 1, 0), 6); signif(qweibull(c(0, 0.5, 1), 1, 0), 6); signif(dweibull(x, 1, Inf), 6); signif(pweibull(x, 1, Inf), 6); signif(qweibull(c(0, 0.5, 1), 1, Inf), 6); signif(dweibull(x, -1, 1), 6); signif(pweibull(x, 1, -1), 6)
#---
x <- c(0, 1, 5, Inf); signif(dweibull(x, 0.5, Inf), 6); signif(dweibull(x, 2, Inf), 6); signif(dweibull(x, 0.5, Inf, log = TRUE), 6); signif(dgamma(x, 0.5, scale = Inf), 6); signif(dgamma(x, 1, scale = Inf), 6)
#---
x <- c(-1, 0, 1, 5, Inf); signif(dweibull(x, Inf, 1), 6); signif(pweibull(x, Inf, 1), 6); signif(qweibull(c(0, 0.5, 1), Inf, 1), 6); signif(pweibull(x, 2, 3, log.p = TRUE), 6); signif(qweibull(c(-1e-5, -2, -30), 2, 3, log.p = TRUE), 6); signif(qweibull(c(-1e-5, -2), 2, 3, lower.tail = FALSE, log.p = TRUE), 6)
#---
pnorm(1, 0, 0) == 1; pnorm(NA_real_, 0, 0); qnorm(NA, 0, 0); pexp(NA, 0); pgamma(NA, 0); pbeta(NA, 0, 0); pt(NA, 0); pf(NA, 0, 0); pchisq(NA, 0); plnorm(NA, 0, 0); plogis(NA, 0, 0); pcauchy(NA, 0, 0); pweibull(NA, 0, 0)
#---
qnorm(c(0, 1), 0, -1); qgamma(c(0, 1), -1); qbeta(c(0, 1), -1, 2); qt(c(0, 1), -1); qf(c(0, 1), -1, 2); qchisq(c(0, 1), -1); qlnorm(c(0, 1), 0, -1); qlogis(c(0, 1), 0, -1); qcauchy(c(0, 1), 0, -1); qweibull(c(0, 1), -1, 1); qexp(c(0, 1), -1); qunif(c(0, 1), 2, 1)
#---
pnorm(1, 0, -1); pgamma(1, -1); pbeta(1, -1, 2); pt(1, -1); pf(1, -1, 2); pchisq(1, -1); plnorm(1, 0, -1); plogis(1, 0, -1); pcauchy(1, 0, -1); pweibull(1, -1, 1); pexp(1, -1); punif(1, 2, 1)
