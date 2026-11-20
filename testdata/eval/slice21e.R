Fn <- ecdf(c(3, 1, 2, 2)); Fn(c(0, 1, 2, 2.5, 3, 4))
#---
Fn <- ecdf(c(3, 1, 2, 2)); summary(Fn)
#---
Fn <- ecdf(c(3, 1, 2, 2)); quantile(Fn)
#---
Fn <- ecdf(c(3, 1, 2, 2)); quantile(Fn, c(0.1, 0.5, 0.9))
#---
Fn <- ecdf(c(3, 1, 2, 2)); quantile(Fn, 0.5, type = 1)
#---
Fn <- ecdf(c(3, 1, 2, 2)); knots(Fn)
#---
Fn <- ecdf(c(3, 1, 2, 2)); class(Fn); is.function(Fn)
#---
Fn <- ecdf(c(3, 1, 2, 2)); Fn(NA)
#---
Fn <- ecdf(c(3, 1, 2, 2)); Fn(numeric(0))
#---
ecdf(numeric(0))
#---
Fn <- ecdf(c(1, NA, 2)); Fn(1.5)
#---
Fn <- ecdf(c(0.5, 0.25, 0.75, 0.25, 1)); Fn(c(0.25, 0.3, 0.9))
#---
Fn <- ecdf(c(0.5, 0.25, 0.75, 0.25, 1)); summary(Fn)
#---
Fn <- ecdf(1:10); Fn(c(-Inf, 5.5, Inf))
#---
Fn <- ecdf(1:10); quantile(Fn, seq(0, 1, 0.25))
#---
Fn <- ecdf(1:10); environment(Fn)$nobs
#---
Fn <- ecdf(c(1, 1, 1)); Fn(0:2)
#---
Fn <- ecdf(c(5, 3, 8)); knots(Fn)
#---
Fn <- ecdf(c(5, 3, 8)); Fn(c(3, 5, 8))
#---
Fn <- ecdf(c(5, 3, 8)); Fn(c(2.999, 3.001))
#---
ecdf("a")
#---
ecdf(c(Inf, -Inf, 0))(c(-Inf, 0, Inf))
#---
Fn <- ecdf(c(3, 1, 2)); Fn(2:3 + 0.5)
#---
Fn <- ecdf(c(3, 1, 2)); summary(Fn)
#---
Fn <- ecdf(c(3, 1, 2)); quantile(Fn, 0.3)
#---
Fn <- ecdf(c(3, 1, 2, 2)); s <- summary(Fn); class(s); names(attributes(s)); unclass(s)
#---
Fn <- ecdf(c(3, 1, 2, 2)); sort(ls(environment(Fn))); environment(Fn)$x; environment(Fn)$y; environment(Fn)$method; environment(Fn)$yleft; environment(Fn)$yright; environment(Fn)$f
#---
Fn <- ecdf(c(3, 1, 2, 2)); names(attributes(Fn)); attr(Fn, "call"); class(attr(Fn, "call"))
#---
Fn <- ecdf(c(3, 1, 2, 2)); Fn(c(a = 1, b = 2)); Fn(matrix(c(1, 2, 3, 4), 2)); Fn(factor("a")); Fn(c(NaN, 1))
#---
Fn <- ecdf(c(3, 1, 2, 2)); quantile(Fn, 0.5, names = FALSE); class(knots(Fn)); summary(knots(Fn))
#---
Fn <- ecdf(c(3, 1, 2, 2)); is.null(environment(Fn)$n); environment(Fn)$nobs; typeof(environment(Fn)$nobs)
#---
x <- poly(1:5, 2); round(x[, 1:2], 7) + 0 # ulp
#---
x <- poly(1:5, 2); attr(x, "coefs") # ulp
#---
x <- poly(1:5, 2); attr(x, "degree"); class(x); dim(x)
#---
poly(1:5, 2, raw = TRUE)
#---
poly(1:5, 3, raw = TRUE)[, 3]
#---
x <- poly(1:5, 2); round(crossprod(unclass(x)), 12) + 0
#---
x <- poly(c(1, 2, 3, 5, 8), 2); round(x[, 1:2], 9) + 0; attr(x, "degree"); class(x) # ulp
#---
x <- poly(c(1, 2, 3, 5, 8), 3); attr(x, "coefs") # ulp
#---
x <- poly(1:5, 2); predict(x, 6:7) # ulp
#---
x <- poly(1:5, 2); predict(x, c(1, 2.5)) # ulp
#---
x <- poly(1:5, 2); round(predict(x)[, 1:2], 7) + 0 # ulp
#---
x <- poly(1:5, 2, raw = TRUE); predict(x, 6)
#---
poly(1:3, 3)
#---
poly(1:5, 0)
#---
round(poly(1:5, 2, simple = TRUE), 7) + 0 # ulp
#---
poly(1:5, 2, coefs = list(alpha = c(3, 3), norm2 = c(1, 5, 10, 14))) # ulp
#---
x <- poly(1:5, 1); round(x[, 1], 7) + 0 # ulp
#---
poly(c(1, 2, NA), 1)
#---
poly("a", 1)
#---
x <- poly(cbind(1:5, c(2, 1, 4, 3, 5)), degree = 2); dimnames(x)[[2]]; dim(x)
#---
x <- poly(cbind(a = 1:5, b = c(2, 1, 4, 3, 5)), degree = 2); round(x[, 1:3], 6) + 0 # ulp
#---
x <- poly(1:5, 2); x[1:2, ] # ulp
#---
x <- poly(1:5, 2); round(x[, 2], 6) # ulp
#---
d <- data.frame(x = 1:5); round(model.matrix(~ poly(x, 2), d)[, 1:3], 7) + 0 # ulp
#---
d <- data.frame(x = c(1, 2, 3, 5, 8)); round(model.matrix(~ poly(x, 2), d)[, 1:3], 6) + 0 # ulp
#---
d <- data.frame(x = 1:5); colnames(model.matrix(~ poly(x, 2, raw = TRUE), d))
#---
d <- data.frame(x = 1:5, g = c("a", "b", "a", "b", "a")); dim(model.matrix(~ poly(x, 2) + g, d))
#---
x <- poly(seq(0, 1, length.out = 7), 3); round(attr(x, "coefs")$alpha, 12); round(attr(x, "coefs")$norm2, 12)
#---
x <- poly(1:20, 5); round(max(abs(crossprod(unclass(x)) - diag(5))), 10)
#---
x <- poly(cbind(1:5, c(2, 1, 4, 3, 5)), degree = 2); names(attributes(x)); attr(x, "degree"); names(attr(x, "coefs")); class(x); dimnames(x)
#---
x <- poly(cbind(1:5, c(2, 1, 4, 3, 5)), degree = 2, raw = TRUE); names(attributes(x)); attr(x, "degree"); round(x[, 1:3], 6)
#---
poly(1:2, 3, raw = TRUE)
#---
poly(c(1, NA), 1, raw = TRUE)
#---
x <- poly(1:3, 2, raw = TRUE, simple = TRUE); attributes(x)
#---
poly(1, 1, raw = TRUE)
#---
poly(c(2, 2, 2), 1)
#---
poly(1:5, 1:2)
#---
x <- poly(c(a = 1, b = 2, c = 3), 2); rownames(x); names(attributes(x))
#---
x <- poly(c(a = 1, b = 2, c = 3), 2, raw = TRUE); rownames(x); names(attributes(x))
#---
x <- poly(1:3, 1); predict(x, c(a = 1)); rownames(predict(x, c(a = 1, b = 2)))
#---
x <- poly(1:5, 2); predict(x, numeric(0))
#---
x <- poly(1:5, 2); predict(x, NA)
#---
x <- poly(1:5, 2); names(attributes(predict(x, 6)))
#---
x <- poly(1:5, 2, raw = TRUE); names(attributes(predict(x, 6))); predict(x, c(2, NA))
#---
x <- poly(1:5, 2); predict(x, 6:7, degree = 1) # ulp
#---
poly(1:5, degree = 2, coefs = list(alpha = 3, norm2 = c(1, 5, 10, 14))) # ulp
#---
poly(c(1, 2, 3, 4), 2, coefs = list(alpha = c(2.5, 2.5), norm2 = c(1, 4, 5, 4))) # ulp
#---
x <- c(1.5, 2.5, 4, 7, 11); p <- poly(x, 3); round(p[, 1:3], 8) + 0; round(crossprod(unclass(p)), 10) + 0 # ulp
#---
tf <- tempfile(); writeLines(c("123456789", "abcdefghi"), tf); read.fwf(tf, c(3, 2, 4))
#---
tf <- tempfile(); writeLines(c("123456789", "abcdefghi"), tf); read.fwf(tf, c(3, -2, 4))
#---
tf <- tempfile(); writeLines(c("123456789", "abcdefghi"), tf); read.fwf(tf, c(3, 2, 4), col.names = c("a", "b", "c"))
#---
tf <- tempfile(); writeLines(c("A  B", "1 2 ", "3 4 "), tf); read.fwf(tf, c(3, 1), header = TRUE)
#---
tf <- tempfile(); writeLines(c("12345", "67890", "11111"), tf); read.fwf(tf, c(2, 3), n = 2)
#---
tf <- tempfile(); writeLines(c("12345", "67890", "11111"), tf); read.fwf(tf, c(2, 3), skip = 1)
#---
tf <- tempfile(); writeLines(c("12345", "67890"), tf); read.fwf(tf, c(2, 3), colClasses = "character")
#---
tf <- tempfile(); writeLines(c("ab12", "cd34"), tf); read.fwf(tf, c(2, 2), stringsAsFactors = TRUE)$V1
#---
tf <- tempfile(); writeLines(c("ab12", "cd34"), tf); read.fwf(tf, c(2, 2), row.names = c("r1", "r2"))
#---
tf <- tempfile(); writeLines(c("ab12", "cd34"), tf); read.fwf(tf, c(2, 2), as.is = TRUE)
#---
tf <- tempfile(); writeLines(c("ab 2", "cd  "), tf); read.fwf(tf, c(2, 2))
#---
tf <- tempfile(); writeLines(c("ab12", "cd3"), tf); read.fwf(tf, c(2, 2))
#---
tf <- tempfile(); writeLines(c("ab12", "cd34", "ef56"), tf); read.fwf(tf, c(2, 2), buffersize = 1)
#---
tf <- tempfile(); writeLines(c("ab12"), tf); read.fwf(tf, c(2, 2), na.strings = "ab")
#---
tf <- tempfile(); writeLines(c("ab12"), tf); read.fwf(tf, 2:3, sep = ",")
#---
tf <- tempfile(); writeLines(c("ab12"), tf); read.fwf(tf, "a")
#---
tf <- tempfile(); writeLines(c("ab12"), tf); read.fwf(tf, c(-2, -2))
#---
tf <- tempfile(); writeLines(c("héllo", "wörld"), tf); x <- read.fwf(tf, c(2, 3)); nchar(x$V1); x$V2 == c("llo", "rld"); dim(x)
#---
tf <- tempfile(); writeLines(c("a b", "cd"), tf); read.fwf(tf, c(1, 1, 1))
#---
tf <- tempfile(); writeLines(c("12 34", "56 78"), tf); read.fwf(tf, c(2, -1, 2), col.names = c("p", "q"))
#---
tf <- tempfile(); writeLines(c("1234", "5678"), tf); x <- read.fwf(tf, c(1, 3)); sapply(x, class); x$V2
#---
tf <- tempfile(); writeLines(c("1.5 x", "2.5 y"), tf); x <- read.fwf(tf, c(3, 2), strip.white = TRUE); x; sapply(x, class)
#---
tf <- tempfile(); writeLines(c("# c", "12", "34"), tf); read.fwf(tf, c(1, 1), comment.char = "#")
#---
tf <- tempfile(); writeLines(c("12", "", "34"), tf); read.fwf(tf, c(1, 1))
#---
tf <- tempfile(); writeLines(c("12", "34"), tf); read.fwf(tf, c(1, 1), header = FALSE, check.names = FALSE, col.names = c("a b", "c d"))
#---
read.fwf(textConnection(c("12", "34")), c(1, 1))
#---
tf <- tempfile(); writeLines(c("12", "34"), tf); read.fwf(tf, widths = numeric(0))
#---
tf <- tempfile(); writeLines(c("1,5 2", "2,5 3"), tf); read.fwf(tf, widths = c(3, 2), dec = ","); sapply(read.fwf(tf, widths = c(3, 2), dec = ","), class)
#---
tf <- tempfile(); writeLines(c("TRUE F", "FALSE T"), tf); x <- read.fwf(tf, c(5, 2)); sapply(x, class); x
#---
tf <- tempfile(); writeLines(c("ab", "cd"), tf); x <- read.fwf(tf, c(1, 1), stringsAsFactors = TRUE); sapply(x, class)
