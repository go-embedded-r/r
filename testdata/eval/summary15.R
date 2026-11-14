summary(c(1.5, 2, 10, 100))
#---
unclass(summary(c(1.5, 2, 10, 100)))
#---
class(summary(c(1.5, 2, 10, 100)))
#---
summary(c(1, 2, NA, 4))
#---
unclass(summary(c(1, 2, NA, 4)))
#---
summary(1:10)
#---
summary(c(123456, 234567, 345678))
#---
summary(c(0.00012345, 0.5, 3))
#---
summary(c(TRUE, FALSE, NA))
#---
unclass(summary(c(TRUE, FALSE, NA)))
#---
summary(c(TRUE, TRUE))
#---
summary(letters)
#---
summary(c("a", NA, ""))
#---
unclass(summary(c(a = "x")))
#---
summary(factor(c("a", "b", "a", NA)))
#---
summary(factor(c("a", "b", "a")), maxsum = 1)
#---
summary(factor(c("a", "b", "c", "a", "b", "a")), maxsum = 2)
#---
summary(factor(character(0)))
#---
summary(list(1, "a", 1:3))
#---
summary(list(a = 1, b = list(2, 3), f = mean))
#---
summary(NULL)
#---
summary(quote(a + b))
#---
summary(c(1.5, 2), digits = 2)
#---
summary(1:3, digits = 2)
#---
format(summary(c(1.5, 2, 100)), digits = 3)
#---
format(summary(c(1, 2, 3)))
#---
x <- summary(c(1, 2, 3)); x[["Mean"]]
#---
summary(numeric(0))
#---
summary(c(NA_real_, NA))
#---
summary(c(a = 1, b = 2))
#---
summary(data.frame(a = 1:4, b = c("x", "y", "x", "z"), c = factor(c("u", "v", "u", "u"))))
#---
summary(data.frame(x = c(1.5, NA, 3), y = c(TRUE, FALSE, TRUE)))
#---
summary(data.frame(g = factor(c("a", "b", "c", "d", "e", "f", "g", "h", "a"))))
#---
summary(data.frame(a = 1:3, b = 4:6), digits = 2)
#---
summary(data.frame())
#---
summary(matrix(1:6, 3))
#---
summary(matrix(c(1.5, 2.5, 3.5, 10), 2, dimnames = list(NULL, c("p", "q"))))
#---
dim(summary(data.frame(a = 1:4, b = c("x", "y", "x", "z"))))
#---
dimnames(summary(data.frame(a = 1:4, bb = 5:8)))
#---
format.pval(c(0.1, 0.0001, 1e-27, NA), digits = 3)
#---
format.pval(c(0.05, 0.5), eps = 0.1)
#---
format.pval(1e-20)
#---
format.pval(c(1e-20, 0.5), digits = 4)
#---
format.pval(1e-20, digits = 1)
#---
format.pval(1e-20, digits = 7)
#---
format.pval(1e-20, eps = 0.001)
#---
format.pval(c(1e-20, 0.00123, 0.5), eps = 0.001)
#---
format.pval(c(1e-20, 0.00123, 0.5), eps = 0.0000001)
#---
format.pval(c(0, 0.5))
#---
format.pval(c(1e-5, 0.5))
#---
format.pval(c(1e-5, 0.123456789, 0.5), digits = 4)
#---
format.pval(c(0.0004, 0.05), digits = 3)
#---
format.pval(c(0.04512, 0.003), digits = 2)
#---
format.pval(NA)
#---
format.pval(numeric(0))
#---
summary.default(1:4)
#---
summary.factor(factor(c("x", "y", "x")))
#---
summary.matrix(matrix(1:4, 2))
#---
summary.data.frame(data.frame(a = 1:2))
#---
print.summaryDefault(summary(c(1, 2, 4)))
#---
# called directly, R fails with "incompatible types (from list to char)"; named here, not run
quote(format.summaryDefault(summary(c(1, 2, 4))))
