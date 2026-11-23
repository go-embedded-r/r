signif(pnorm(0.3, lower.tail = NA_integer_), 12); signif(pnorm(0.3, lower.tail = NaN), 12); signif(pnorm(0.3, log.p = NA_integer_), 12); signif(pnorm(0.3, log.p = NaN), 12); signif(dnorm(0.3, log = NA_integer_), 12); signif(dnorm(0.3, log = NaN), 12)
#---
signif(qnorm(0.3, lower.tail = NA_integer_), 12); signif(qnorm(-1, log.p = NA_integer_), 12); signif(qnorm(-1, log.p = NaN), 12)
#---
pnorm(1, mean = c(a = 0, b = 1)); pnorm(c(a = 1, b = 2), mean = c(c = 0, d = 1)); qnorm(0.5, sd = c(a = 1, b = 2))
#---
signif(pgamma(1, 2, scale = c(a = 1, b = 2)), 12); signif(pgamma(1, 2, rate = c(a = 1, b = 2)), 12); signif(pgamma(1, 2, rate = c(a = 1), scale = c(c = 1)), 12)
#---
pgamma(c(a = 1, b = 2), c(c = 2, d = 3)); dgamma(1, shape = c(a = 2, b = 3), rate = 1)
#---
pnorm(numeric(0)); pnorm(1, numeric(0)); pnorm(1, sd = numeric(0)); qgamma(numeric(0), 2); dgamma(1, numeric(0)); pgamma(1, 2, scale = numeric(0)); pbeta(0.5, numeric(0), 2)
#---
tryCatch(dnorm(0, log.p = TRUE), error = function(e) conditionMessage(e)); tryCatch(pnorm(0, log = TRUE), error = function(e) conditionMessage(e)); tryCatch(qnorm(0.5, lower = FALSE, log = TRUE), error = function(e) conditionMessage(e))
#---
tryCatch(pgamma(1), error = function(e) conditionMessage(e)); tryCatch(pbeta(0.5, 1), error = function(e) conditionMessage(e)); tryCatch(pweibull(1), error = function(e) conditionMessage(e))
#---
signif(qnorm(c(-1, -0.2), log.p = TRUE), 12); signif(qnorm(c(-1, -0.2), lower.tail = FALSE, log.p = TRUE), 12); signif(qexp(c(-1, -0.2), log.p = TRUE), 12); signif(qexp(c(-1, -0.2), lower.tail = FALSE, log.p = TRUE), 12)
#---
qnorm(c(-0.5, 1.5, 2)); qnorm(c(0.5, 0.1), 0, -1); qnorm(0.5, log.p = TRUE); qexp(1.0000001); qgamma(-0.1, 2)
#---
dnorm(Inf, Inf); dnorm(Inf, Inf, log = TRUE); dnorm(-Inf, -Inf); dnorm(Inf, Inf, Inf); dnorm(Inf, Inf, Inf, log = TRUE); dnorm(Inf); dnorm(Inf, log = TRUE); dnorm(0, 0, Inf); dnorm(0, 0, Inf, log = TRUE)
#---
dlnorm(Inf, 0, Inf); dlnorm(Inf, 0, Inf, log = TRUE); dlnorm(-1); dlnorm(-1, log = TRUE); plnorm(-1); plnorm(-1, lower.tail = FALSE); plnorm(-1, log.p = TRUE); plnorm(-1, lower.tail = FALSE, log.p = TRUE)
#---
dlnorm(0); dlnorm(0, log = TRUE); plnorm(0); plnorm(0, lower.tail = FALSE); plnorm(0, log.p = TRUE); plnorm(0, lower.tail = FALSE, log.p = TRUE)
#---
pf(-1, Inf, 3); pf(-1, Inf, 3, lower.tail = FALSE); pf(-1, Inf, 3, log.p = TRUE); pf(-1, Inf, 3, lower.tail = FALSE, log.p = TRUE); df(-1, Inf, 3); df(-1, Inf, 3, log = TRUE)
#---
pf(0, Inf, 3); pf(0, Inf, 3, lower.tail = FALSE); pf(0, Inf, 3, log.p = TRUE); pf(0, Inf, 3, lower.tail = FALSE, log.p = TRUE); signif(df(0.5, Inf, 3), 12); signif(df(0.5, Inf, 3, log = TRUE), 12)
#---
signif(pf(c(0.5, 2), Inf, 3), 12); signif(pf(c(0.5, 2), Inf, 3, lower.tail = FALSE), 12); signif(qf(c(0.1, 0.9), Inf, 3), 12); signif(qf(c(0.1, 0.9), Inf, 3, lower.tail = FALSE), 12)
#---
signif(pf(c(0.5, 2), Inf, 3, log.p = TRUE), 12); signif(pf(c(0.5, 2), Inf, 3, lower.tail = FALSE, log.p = TRUE), 12); signif(qf(c(-2, -0.1), Inf, 3, log.p = TRUE), 12)
#---
dgamma(1, 1, 1, log = TRUE); dgamma(0, 1, log = TRUE); dunif(0.5, log = TRUE); dunif(-1, log = TRUE); dunif(2, log = TRUE)
#---
pgamma(0, 2, log.p = TRUE); pgamma(Inf, 2, lower.tail = FALSE, log.p = TRUE); pgamma(3, 2, scale = Inf); pgamma(0, 0); pgamma(1, 0); pgamma(1, 0, log.p = TRUE); pgamma(1, 0, lower.tail = FALSE)
#---
qbeta(0.3, 0, 0); qbeta(0.5, 0, 0); qbeta(0.7, 0, 0); qbeta(0.3, 0, 0, lower.tail = FALSE); qbeta(0.5, 0, 0, lower.tail = FALSE); pbeta(c(-1, 0, 0.5, 1, 2), 0, 0); dbeta(c(0, 0.5, 1), 0, 0)
#---
pbeta(c(0, 0.4, 0.5, 1), 0, 2); pbeta(c(0, 0.4, 0.5, 1), 2, 0); pbeta(c(0, 0.4, 0.5, 1), Inf, Inf); dbeta(c(0, 0.5, 1), 0, 2); dbeta(c(0, 0.5, 1), 2, 0); dbeta(c(0, 0.5, 1), Inf, Inf); qbeta(0.3, 0, 2); qbeta(0.3, 2, 0)
#---
pweibull(c(1, 2, 3), Inf, 2); dweibull(c(-1, 1, 2, 3, Inf), Inf, 2); pweibull(c(1, 2, 3), Inf, 2, lower.tail = FALSE); qweibull(0.3, Inf, 2); signif(pweibull(2, Inf, 2, log.p = TRUE), 12); signif(pweibull(2, Inf, 2, lower.tail = FALSE, log.p = TRUE), 12)
#---
pexp(c(-1, 0, 1, Inf), Inf); pexp(c(-1, 0, 1, Inf), Inf, lower.tail = FALSE); dexp(c(0, 1, Inf), Inf); dexp(c(-1, 0, 1, Inf), 0); pexp(c(-1, 0, 1, Inf), 0); qexp(0.3, 0); qexp(0.3, Inf); qexp(1, Inf)
#---
dunif(c(1, 3), 2, 2); punif(c(1, 2, 3), 2, 2); qunif(c(0, 0.5, 1), 2, 2); dunif(1, 0, Inf); dunif(1, -Inf, Inf); punif(1, 0, Inf); dunif(1, 0, -Inf); dunif(1, Inf, Inf); dunif(1, 5, 2)
#---
pt(c(-2, 0, 2), Inf); signif(qt(c(0.1, 0.9), Inf), 12); pt(1, 0); pt(1, -1); signif(dt(1, Inf), 12)
#---
pchisq(1, 0); pchisq(0, 0); pchisq(1, 0, lower.tail = FALSE); pchisq(1, Inf); pchisq(1, -1); qchisq(0.5, 0); qchisq(0.5, Inf)
#---
pcauchy(c(1, 2), 1, 0); pcauchy(c(0, 1, 2), 1, 0); dcauchy(c(0, 1, 2), 1, 0); qcauchy(0.5, 1, 0); plogis(c(0, 1, 2), 1, 0); dlogis(c(0, 1, 2), 1, 0); qlogis(0.3, 1, 0); pcauchy(1, 0, -1); plogis(1, 0, -1)
#---
signif(qcauchy(c(-711, -750), log.p = TRUE), 12); signif(qcauchy(c(-711, -750), lower.tail = FALSE, log.p = TRUE), 12); signif(qcauchy(-1000, 3, 2, log.p = TRUE), 12); signif(qcauchy(-700, log.p = TRUE), 12)
#---
qnorm(c(-800, -1e5), 0, Inf, log.p = TRUE); qnorm(c(-800, -1e5), 0, Inf, lower.tail = FALSE, log.p = TRUE); signif(qnorm(-1000, 0, 2, log.p = TRUE), 12); signif(qnorm(-1000, 5, 0.5, lower.tail = FALSE, log.p = TRUE), 12)
#---
signif(qweibull(-1000, 1, 1, log.p = TRUE), 12); signif(qbeta(-1000, 1, 1, lower.tail = FALSE, log.p = TRUE), 12); signif(qgamma(-1000, 1, 1, log.p = TRUE), 12)#---
signif(qlnorm(c(-800, -1e5), 0, 1, log.p = TRUE), 12); signif(qlnorm(-800, 0, Inf, log.p = TRUE), 12); signif(qlnorm(-800, 0, 1, lower.tail = FALSE, log.p = TRUE), 12)
