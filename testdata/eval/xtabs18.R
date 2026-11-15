d <- data.frame(a = factor(c("x", "y", "x", "y", "x")), b = factor(c("u", "u", "v", "v", "u")), n = c(1, 2, 3, 4, 5)); xtabs(~ a + b, d)
#---
d <- data.frame(a = factor(c("x", "y", "x", "y", "x")), b = factor(c("u", "u", "v", "v", "u")), n = c(1, 2, 3, 4, 5)); x <- xtabs(~ a + b, d); attributes(x)
#---
d <- data.frame(a = factor(c("x", "y", "x", "y", "x")), b = factor(c("u", "u", "v", "v", "u")), n = c(1, 2, 3, 4, 5)); xtabs(n ~ a + b, d)
#---
d <- data.frame(a = factor(c("x", "y", "x", "y", "x")), b = factor(c("u", "u", "v", "v", "u")), n = c(1, 2, 3, 4, 5)); x <- xtabs(n ~ a, d); x; summary(x)
#---
d <- data.frame(a = factor(c("x", "y", "x", "y", "x")), b = factor(c("u", "u", "v", "v", "u")), n = c(1, 2, 3, 4, 5)); xtabs(cbind(n, n2 = n * 2) ~ a, d)
#---
d <- data.frame(a = factor(c("x", "y", "x", "y", "x")), b = factor(c("u", "u", "v", "v", "u")), n = c(1, 2, 3, 4, 5)); xtabs(~ a, d, subset = n > 1)
#---
d <- data.frame(a = c("x", "y", "x", "y", "x"), b = c("u", "u", "v", "v", "u"), n = c(1, 2, 3, 4, 5)); xtabs(~ a + b, d, subset = a == "x", drop.unused.levels = TRUE)
#---
d <- data.frame(a = factor(c("x", NA, "x", "y", "x")), b = factor(c("u", "u", "v", "v", "u"))); xtabs(~ a + b, d)
#---
d <- data.frame(a = factor(c("x", NA, "x", "y", "x")), b = factor(c("u", "u", "v", "v", "u"))); xtabs(~ a + b, d, addNA = TRUE)
#---
d <- data.frame(a = factor(c("x", NA, "x", "y", "x")), b = factor(c("u", "u", "v", "v", "u"))); x <- xtabs(~ a + b, d, na.action = na.pass, addNA = TRUE); x
#---
d <- data.frame(a = factor(c("x", "y", "x", "y", "x")), b = factor(c("u", "u", "v", "v", "u"))); summary(xtabs(~ a + b, d))$n.cases
#---
d <- data.frame(a = factor(c("x", "y", "x", "y", "x")), b = factor(c("u", "u", "v", "v", "u")), c = c("p", "q", "p", "p", "q")); xtabs(~ a + b + c, d)
#---
d <- data.frame(a = factor(c("x", "y", "x", "y", "x")), n = 1:5, m = 5:1); xtabs(n ~ ., d)
#---
xtabs(a ~ b, list(a = c(1, 2, 3), b = c("p", "q", "p")))
#---
xtabs(1:3 ~ c("a", "b", "a"))
#---
xtabs(~ a, data.frame(a = c(TRUE, FALSE, TRUE)))
#---
xtabs(y ~ g, data.frame(y = c(1.5, NA, 2), g = c("a", "a", "b")))
#---
xtabs(y ~ g, data.frame(y = c(1.5, NA, 2), g = c("a", "a", "b")), na.action = na.pass)
#---
xtabs(y ~ g, data.frame(y = c(1.5, NA, 2), g = c("a", "a", "b")), na.action = na.omit)
#---
xtabs(y ~ g, data.frame(y = c(1.5, 2, 2), g = c("a", NA, "b")), addNA = TRUE)
#---
print(table(c(1, 2, 2), c(1, 1, 2)), zero.print = ".")
#---
table(c(1, 2, 2)); summary(table(c(1, 2, 2)))
#---
summary(table(c("a", "b"), c("u", "v")))$n.vars
#---
xtabs(~ a, data.frame(a = integer(0)))
#---
x <- xtabs(~ a + b, data.frame(a = c(1, 1, 2), b = c("p", "q", "p"))); print(x, zero.print = "."); dim(x); dimnames(x)
#---
x <- summary.table(table(c(1, 2, 2))); x$n.cases; x$n.vars; class(x)
#---
print.summary.table(summary(table(c("a", "b", "b"))))
#---
summary.table(1:3)
#---
summary.table(table(c("a", "b"), c("u", "v")))$n.vars
