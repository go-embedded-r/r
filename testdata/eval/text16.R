x <- "The quick brown fox jumps over the lazy dog. It was a sunny day  and everything was fine!  Really? Yes."; strwrap(x, width = 30)
#---
x <- "The quick brown fox jumps over the lazy dog. It was a sunny day  and everything was fine!  Really? Yes."; strwrap(x, width = 20)
#---
x <- "The quick brown fox jumps over the lazy dog. It was a sunny day  and everything was fine!  Really? Yes."; strwrap(x, width = 40, indent = 2, exdent = 4)
#---
x <- "The quick brown fox jumps over the lazy dog. It was a sunny day  and everything was fine!  Really? Yes."; strwrap(x, width = 40, prefix = "> ")
#---
x <- "The quick brown fox jumps over the lazy dog. It was a sunny day  and everything was fine!  Really? Yes."; strwrap(x, width = 40, prefix = "> ", initial = "* ")
#---
x <- "The quick brown fox jumps over the lazy dog. It was a sunny day  and everything was fine!  Really? Yes."; strwrap(x)
#---
x <- "The quick brown fox jumps over the lazy dog. It was a sunny day  and everything was fine!  Really? Yes."; strwrap(x, width = 30, indent = 1, exdent = 1, prefix = "#")
#---
strwrap(c("a b c", "d e f"), width = 5)
#---
strwrap(c("a b c", "d e f"), width = 5, simplify = FALSE)
#---
strwrap("", width = 10)
#---
strwrap(character(0), width = 10)
#---
strwrap(NA, width = 10)
#---
strwrap(c("x", NA), width = 10)
#---
strwrap("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa bbb", width = 10)
#---
strwrap("para one is here\n\npara two is here", width = 12)
#---
strwrap("line one\nline two", width = 40)
#---
strwrap("  leading and   multiple   spaces  ", width = 40)
#---
strwrap("tab\tseparated\twords", width = 40)
#---
strwrap("Hello there.  Next sentence here.", width = 50)
#---
strwrap("Hello there. Next sentence here.", width = 50)
#---
strwrap("end of sentence.  ", width = 50)
#---
strwrap("a b", width = 1)
#---
strwrap("a b", width = 3)
#---
strwrap("abc def", width = 7)
#---
strwrap("abc def", width = 8)
#---
strwrap("abc def ghi", width = 8)
#---
strwrap("héllo wörld ünï", width = 8)
#---
strwrap("日本語 テキスト", width = 8)
#---
strwrap(factor("a b c"), width = 3)
#---
strwrap(1:3)
#---
strwrap(c(a = "x y z"), width = 3)
#---
strwrap("x y z", width = 3, indent = 5)
#---
strwrap("The end. 'Quoted.' (Paren.) Next", width = 80)
#---
strwrap("The end.  'Quoted.'  (Paren.)  Next", width = 80)
#---
strwrap("Mr. Smith  went", width = 80)
#---
strwrap(c("a b c", "d e f"), width = 3, initial = "*", prefix = ">")
#---
strwrap(c("a b c", "d e f"), width = 5, initial = "*", prefix = ">")
#---
strwrap("aaa bbb.  ccc ddd", width = 13)
#---
strwrap("aaa bbb.  ccc ddd", width = 14)
#---
strwrap("aaa bbb. ccc ddd", width = 12)
#---
strwrap("aaa bbb. ccc ddd", width = 13)
#---
strwrap("p1 a b\n\np2 c d", width = 4, simplify = FALSE)
#---
strwrap("p1 a b\n  \np2 c d", width = 40)
#---
strwrap("p1\n\n\n\np2", width = 40)
#---
strwrap("p1\n\np2", width = 40, indent = 2, exdent = 3, prefix = "#", initial = "!")
#---
strwrap("p1 aaa bbb\n\np2 ccc ddd", width = 8, indent = 2, exdent = 3, prefix = "#", initial = "!")
#---
strwrap("\n\nleading blank", width = 40)
#---
strwrap("trailing\n\n", width = 40)
#---
strwrap("a.b.  c", width = 40)
#---
strwrap("e.g.  foo?  bar!  baz.\"  qux", width = 80)
#---
strwrap("What?!  ok", width = 80)
#---
strwrap("x.)  y", width = 80)
#---
strwrap("x.]  y", width = 80)
#---
strwrap("x:  y", width = 80)
#---
strwrap("a b c d e f g h", width = 4)
#---
strwrap("aa bb cc dd", width = 6)
#---
strwrap("aa bb cc dd", width = 5)
#---
strwrap("short", width = 80, indent = 90)
#---
A <- matrix(1:4, 2, dimnames = list(c("a", "b"), c("c", "d"))); B <- matrix(1:6, 2, dimnames = list(c("x", "y"), c("p", "q", "r"))); .kronecker(A, B)
#---
A <- matrix(1:4, 2, dimnames = list(c("a", "b"), c("c", "d"))); B <- matrix(1:6, 2, dimnames = list(c("x", "y"), c("p", "q", "r"))); .kronecker(A, B, make.dimnames = TRUE)
#---
A <- matrix(1:4, 2); B <- matrix(1:6, 2); A %x% B
#---
.kronecker(1:2, 1:3)
#---
.kronecker(diag(2), matrix(1, 2, 2))
#---
A <- matrix(1:4, 2); B <- matrix(1:6, 2); .kronecker(A, B, FUN = "+")
#---
dim(.kronecker(array(1:8, c(2, 2, 2)), matrix(1:4, 2)))
#---
1:2 %x% 1:3
#---
.kronecker(matrix(1:4, 2), 1:2)
#---
.kronecker(1:3, matrix(1:4, 2))
#---
dput(.kronecker(1:2, 1:3))
#---
.kronecker(matrix(1:4, 2), matrix(1, 1, 2))
#---
df <- data.frame(g = c("a", "b", "a", "b", "c"), x = 1:5, y = c(2.5, 3.5, 1, 8, 9)); by(df, df$g, function(d) sum(d$x))
#---
df <- data.frame(g = c("a", "b", "a", "b", "c"), x = 1:5, y = c(2.5, 3.5, 1, 8, 9)); b <- by(df, df$g, function(d) sum(d$x)); class(b)
#---
df <- data.frame(g = c("a", "b", "a", "b", "c"), x = 1:5, y = c(2.5, 3.5, 1, 8, 9)); b <- by(df, df$g, function(d) sum(d$x)); dput(b)
#---
df <- data.frame(g = c("a", "b", "a", "b", "c"), x = 1:5, y = c(2.5, 3.5, 1, 8, 9)); by(df, df$g, function(d) colMeans(d[, c("x", "y")]))
#---
df <- data.frame(g = c("a", "b", "a", "b", "c"), x = 1:5, y = c(2.5, 3.5, 1, 8, 9)); by(df$x, df$g, sum)
#---
df <- data.frame(g = c("a", "b", "a", "b", "c"), x = 1:5, y = c(2.5, 3.5, 1, 8, 9)); by(df, list(df$g), nrow)
#---
df <- data.frame(g = c("a", "b", "a", "b", "c"), x = 1:5, y = c(2.5, 3.5, 1, 8, 9)); by(df, list(grp = df$g), nrow)
#---
df <- data.frame(g = c("a", "b", "a", "b", "c"), x = 1:5, y = c(2.5, 3.5, 1, 8, 9)); by(df, df[, "g", drop = FALSE], nrow)
#---
df <- data.frame(g = c("a", "b", "a", "b", "c"), x = 1:5, y = c(2.5, 3.5, 1, 8, 9)); by(df$x, list(df$g, df$x > 2), sum)
#---
df <- data.frame(g = c("a", "b", "a", "b", "c"), x = 1:5, y = c(2.5, 3.5, 1, 8, 9)); names(dimnames(by(df, df$g, nrow)))
#---
df <- data.frame(g = c("a", "b", "a", "b", "c"), x = 1:5, y = c(2.5, 3.5, 1, 8, 9)); by(df, df$g, function(d) d)
#---
df <- data.frame(g = c("a", "b", "a", "b", "c"), x = 1:5, y = c(2.5, 3.5, 1, 8, 9)); by(df, df$g, function(d) NULL)
#---
df <- data.frame(g = c("a", "b", "a", "b", "c"), x = 1:5, y = c(2.5, 3.5, 1, 8, 9)); print(by(df$y, df$g, range))
#---
df <- data.frame(g = c("a", "b", "a", "b", "c"), x = 1:5, y = c(2.5, 3.5, 1, 8, 9)); b2 <- by(df, df$g, function(d) d$x); b2$a
#---
df <- data.frame(g = c("a", "b", "a", "b", "c"), x = 1:5, y = c(2.5, 3.5, 1, 8, 9)); b2 <- by(df, df$g, function(d) d$x); b2[["b"]]
#---
df <- data.frame(g = c("a", "b", "a", "b", "c"), x = 1:5, y = c(2.5, 3.5, 1, 8, 9)); b2 <- by(df, df$g, function(d) d$x); b2[1]
#---
df <- data.frame(g = c("a", "b", "a", "b", "c"), x = 1:5, y = c(2.5, 3.5, 1, 8, 9)); b2 <- by(df, df$g, function(d) d$x); c(b2)
#---
df <- data.frame(g = c("a", "b", "a", "b", "c"), x = 1:5, y = c(2.5, 3.5, 1, 8, 9)); tryCatch(by(df, df$g[-1], nrow), error = function(e) conditionMessage(e))
#---
df <- data.frame(g = c("a", "b", "a", "b", "c"), x = 1:5, y = c(2.5, 3.5, 1, 8, 9)); by(df$x, df$g, function(v) v * 2, simplify = FALSE)
#---
df <- data.frame(g = c("a", "b", "a", "b", "c"), x = 1:5, y = c(2.5, 3.5, 1, 8, 9)); as.list(by(df$x, df$g, sum))
#---
by(1:6, rep(c("u", "v"), 3), sum)
#---
by(matrix(1:6, 3), c(1, 1, 2), colSums)
#---
df <- data.frame(g = c("a", "b", "a", "b", "c"), x = 1:5); dput(by(df$x, df$g, sum))
#---
dput(by(matrix(1:6, 3), c(1, 1, 2), colSums))
#---
by.default(1:4, c("a", "b", "a", "b"), sum)
#---
by.data.frame(data.frame(x = 1:4), c("a", "b", "a", "b"), function(d) sum(d$x))
#---
print.by(by(1:4, c("a", "b", "a", "b"), sum))
#---
r <- rle(c(1, 1, 2, 2, 2, 3)); r
#---
r <- rle(c(1, 1, 2, 2, 2, 3)); print(r)
#---
rle(c("a", "a", "b"))
#---
rle(c(TRUE, TRUE, FALSE))
#---
rle(integer())
#---
rle(c(NA, NA, 1))
#---
print.rle(rle(c(1, 1, 2)))
#---
rle(rep(1:30, each = 2))
#---
rle(c(1.5, 1.5, 2.25))
#---
inverse.rle(rle(c("a", "a", "b")))
#---
file.info("nosuch.txt")
#---
dput(file.info("nosuch"))
#---
dput(file.info(c("nosuch", "other"), extra_cols = FALSE))
#---
dput(file.info(character()))
#---
f <- tempfile(); writeLines(c("hello", "world"), f); fi <- file.info(f); names(fi)
#---
f <- tempfile(); writeLines(c("hello", "world"), f); fi <- file.info(f); class(fi)
#---
f <- tempfile(); writeLines(c("hello", "world"), f); fi <- file.info(f); dim(fi)
#---
f <- tempfile(); writeLines(c("hello", "world"), f); fi <- file.info(f); fi$size
#---
f <- tempfile(); writeLines(c("hello", "world"), f); fi <- file.info(f); fi$isdir
#---
f <- tempfile(); writeLines(c("hello", "world"), f); fi <- file.info(f); fi$mode
#---
f <- tempfile(); writeLines(c("hello", "world"), f); fi <- file.info(f); class(fi$mode)
#---
f <- tempfile(); writeLines(c("hello", "world"), f); fi <- file.info(f); class(fi$mtime)
#---
f <- tempfile(); writeLines(c("hello", "world"), f); fi <- file.info(f); identical(rownames(fi), f)
#---
file.info("nosuch.txt")$size
#---
f <- tempfile(); writeLines(c("hello", "world"), f); file.size(f)
#---
f <- tempfile(); writeLines(c("hello", "world"), f); file.size(c(f, "nosuch"))
#---
f <- tempfile(); writeLines(c("hello", "world"), f); class(file.mtime(f))
#---
file.mtime("nosuch")
#---
f <- tempfile(); writeLines(c("hello", "world"), f); file.mode(f)
#---
f <- tempfile(); writeLines(c("hello", "world"), f); class(file.mode(f))
#---
f <- tempfile(); writeLines(c("hello", "world"), f); file.info(c(f, "nosuch"))$isdir
#---
file.info(".")$isdir
#---
f <- tempfile(); writeLines(c("hello", "world"), f); names(file.info(f, extra_cols = FALSE))
#---
f <- tempfile(); writeLines(c("hello", "world"), f); sapply(file.info(f), class)
#---
f <- tempfile(); writeLines(c("hello", "world"), f); file.info(f)$size == file.size(f)
#---
x <- c(1234567.891, 0.5, 0, 12); format(x, big.mark = ",")
#---
x <- c(1234567.891, 0.5, 0, 12); format(x, big.mark = ",", big.interval = 2)
#---
x <- c(1234567.891, 0.5, 0, 12); format(x, decimal.mark = ",")
#---
x <- c(1234567.891, 0.5, 0, 12); format(x, drop0trailing = TRUE)
#---
format(c(1.5, 2.25, 3), drop0trailing = TRUE)
#---
format(c(1, 2.5), nsmall = 2, drop0trailing = TRUE)
#---
x <- c(1234567.891, 0.5, 0, 12); format(x, zero.print = "-")
#---
format(c(0, 1, 2), zero.print = ".")
#---
format(1234.5678, small.mark = " ", digits = 8)
#---
format(0.123456789, small.mark = " ", small.interval = 3, digits = 9)
#---
format(1e6, big.mark = ",", scientific = FALSE)
#---
format(c(1e5, 2), drop0trailing = TRUE)
#---
x <- c(1234567.891, 0.5, 0, 12); formatC(x, format = "f", digits = 2, big.mark = ",")
#---
x <- c(1234567.891, 0.5, 0, 12); formatC(x, format = "f", digits = 2, drop0trailing = TRUE)
#---
formatC(c(1.5, 2, 3.25), format = "f", digits = 3, drop0trailing = TRUE)
#---
x <- c(1234567.891, 0.5, 0, 12); formatC(x, format = "f", digits = 2, decimal.mark = ",")
#---
x <- c(1234567.891, 0.5, 0, 12); formatC(x, format = "f", digits = 2, big.mark = ".", decimal.mark = ",")
#---
formatC(c(0, 1), zero.print = "-")
#---
formatC(c(0, 1), format = "d", zero.print = "zero")
#---
formatC(1234567L, big.mark = " ", big.interval = 3)
#---
formatC(1234567L, big.mark = ",", big.interval = 2, format = "d")
#---
formatC(1234.5678, format = "f", digits = 6, small.mark = "'", small.interval = 3)
#---
prettyNum(1234567, big.mark = ",")
#---
prettyNum(c(123456, 1234.5), big.mark = ",")
#---
prettyNum(1234567, big.mark = ",", big.interval = 2)
#---
prettyNum(0.000123456, small.mark = " ")
#---
prettyNum(c(0, 1), zero.print = "-")
#---
prettyNum(1234.5, decimal.mark = ",")
#---
prettyNum(c(1.5, 2), drop0trailing = TRUE)
#---
prettyNum(c("1.50", "2.00"), drop0trailing = TRUE)
#---
formatC(c(1, 2.5), format = "e", digits = 2, drop0trailing = TRUE)
#---
format(-1234567, big.mark = ",")
#---
format(c(-1, 1) * 1234.5, big.mark = " ")
#---
format(123456789L, big.mark = "'")
#---
format(c(1234.5, NA), big.mark = ",", nsmall = 1)
