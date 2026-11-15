signif(pchisq(3.84, 1), 10)
#---
signif(pchisq(c(0.5, 1, 5, 10, 50, 100), 3), 10)
#---
signif(pchisq(5, c(1, 2, 10), lower.tail = FALSE), 10)
#---
signif(pchisq(2, 4, log.p = TRUE), 10)
#---
signif(pchisq(2, 4, lower.tail = FALSE, log.p = TRUE), 10)
#---
pchisq(c(-1, 0, Inf), 2)
#---
pchisq(c(0, 1), 0)
#---
pchisq(1, -1)
#---
pchisq(c(NA, NaN, 1), 2)
#---
pchisq(numeric(0), 2)
#---
signif(pchisq(c(a = 1, b = 2), 3), 10)
#---
signif(pchisq(1:6, c(1, 2)), 10)
#---
signif(pchisq(1, 2, ncp = 1.5), 9)
#---
signif(pchisq(c(1, 10), 5, ncp = 3, lower.tail = FALSE), 9)
#---
signif(pchisq(2L, 3L), 10)
#---
pchisq("a", 1)
#---
pchisq(1)
#---
signif(qchisq(0.95, 1:5), 10)
#---
signif(qchisq(c(0, 0.5, 1), 4), 10)
#---
signif(qchisq(0.05, 3, lower.tail = FALSE), 10)
#---
signif(qchisq(log(0.5), 3, log.p = TRUE), 10)
#---
qchisq(c(-0.1, 1.5), 2)
#---
qchisq(0.5, 0)
#---
signif(qchisq(0.5, 3, ncp = 2), 8)
#---
signif(qchisq(c(x = 0.1, y = 0.9), 2), 10)
#---
t2 <- table(c(1, 1, 2, 2, 1, 2, 1, 2, 2, 1, 1, 2), c("a", "b", "a", "b", "a", "a", "b", "b", "a", "a", "a", "b")); capture.output(print(summary(t2)))
#---
t2 <- table(c(1, 1, 2, 2, 1, 2, 1, 2, 2, 1, 1, 2), c("a", "b", "a", "b", "a", "a", "b", "b", "a", "a", "a", "b")); s <- summary(t2); signif(s$statistic, 10)
#---
t2 <- table(c(1, 1, 2, 2, 1, 2, 1, 2, 2, 1, 1, 2), c("a", "b", "a", "b", "a", "a", "b", "b", "a", "a", "a", "b")); s <- summary(t2); c(s$parameter, signif(s$p.value, 10))
#---
t2 <- table(c(1, 1, 2, 2, 1, 2, 1, 2, 2, 1, 1, 2), c("a", "b", "a", "b", "a", "a", "b", "b", "a", "a", "a", "b")); s <- summary(t2); c(names(s), class(s), s$approx.ok, is.null(s$call))
#---
i <- 0:59; t3 <- table(c("u", "v")[1 + (i * 7 + i %/% 3) %% 2], c("p", "q", "r")[1 + (i * 5 + i %/% 4) %% 3], c("x", "y")[1 + (i %/% 5) %% 2]); capture.output(print(summary(t3)))
#---
round(1.607143, 3.79)
#---
round(c(2.5, 0.125), 0.5 + 1e-9)
#---
d <- data.frame(a = factor(c("x", "y", "x", "y", "x", "x", "y", "y", "x")), b = factor(c("u", "u", "v", "v", "u", "u", "v", "u", "v")), n = c(1, 2, 3, 4, 5, 6, 7, 8, 9)); capture.output(print(summary(xtabs(n ~ a + b, d))))
#---
d <- data.frame(a = factor(c("x", "y", "x", "y", "x")), b = factor(c("u", "u", "v", "v", "u")), n = c(1, 2, 3, 4, 5)); capture.output(print(summary(xtabs(n ~ a + b, d))))
#---
summary(table(c(1, 1, 2, 3)))$n.vars
#---
t2 <- table(c(1, 1, 2, 2, 1, 2, 1, 2, 2, 1, 1, 2), c("a", "b", "a", "b", "a", "a", "b", "b", "a", "a", "a", "b")); capture.output(print(summary(t2), digits = 6))
