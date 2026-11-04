## print() and auto-print of vectors and lists. A line holding only #--- separates cells.
#---
1:3
#---
c(1.5, 2, 3.25)
#---
c(TRUE, NA, FALSE)
#---
c("a", "bb", NA)
#---
NULL
#---
character(0); numeric(0); integer(0); logical(0); list()
#---
1:30
#---
seq_len(100)
#---
c(0.1, 0.25, 100)
#---
c(1e-10, 1, 1e10)
#---
c(123456789, 0.5)
#---
1e5; 1e-5; 100000L; 123456.7; 0.1 + 0.2; 1/3; 2/3; 1e15; 1e16; 22/7
#---
c(1.123456789, 22.1, 333)
#---
c(NA, 1.5, NaN, Inf, -Inf)
#---
c(-1.5, 2, NA)
#---
100 * (1:12)
#---
c(0.00001, 123456)
#---
c(2^100, 2^-100)
#---
x <- c(1.5, NA, 3); x
#---
c(a = 1, b = 2, c = 3)
#---
c(first = 1.5, second = NA, third = 3)
#---
c(alpha = "x", beta = "yy")
#---
c(a = TRUE, b = NA)
#---
x <- 1:26; names(x) <- letters_ <- c("a","b","c","d","e","f","g","h","i","j","k","l","m","n","o","p","q","r","s","t","u","v","w","x","y","z"); x
#---
c(a = 1L)
#---
x <- c(1, 2); names(x) <- c("a", NA); x
#---
c("a", "b")[0]; c(a = 1)[0]; c(a = "x")[0]
#---
list(1, "a", TRUE)
#---
list(a = 1, b = "x", c = NULL)
#---
list(1, list(2, 3), "z")
#---
list(a = list(b = list(c = 1)))
#---
list(list())
#---
list(x = list())
#---
list(NULL)
#---
l <- list(1, 2); names(l) <- c("a", ""); l
#---
l <- list(1); names(l) <- "a b"; l
#---
x <- 1:3; attr(x, "foo") <- "bar"; x
#---
x <- structure(1:3, class = "foo"); x
#---
x <- structure(list(a = 1), class = "rec"); x
#---
x <- structure(1:3, a = 1, b = c("x", "y")); x
#---
list(a = structure(1, u = "v"))
#---
structure(c(a = 1, b = 2), k = 3)
#---
quote(x + y)
#---
quote(x)
#---
quote(f(a, b = 2))
#---
quote(if (a) b else c)
#---
quote(function(x, y = 2) x + y)
#---
x <- 1:6; dim(x) <- c(2, 3); x
#---
x <- c(1.5, 2.5, 3.5, 4.5); dim(x) <- c(2, 2); x
#---
x <- c("a", "bb", "ccc", "d"); dim(x) <- c(2, 2); x
#---
x <- c(TRUE, FALSE, NA, TRUE); dim(x) <- c(2, 2); x
#---
x <- 1:6; dim(x) <- c(2, 3); dimnames(x) <- list(c("r1", "r2"), c("a", "b", "c")); x
#---
x <- 1:6; dim(x) <- c(2, 3); dimnames(x) <- list(NULL, c("a", "b", "c")); x
#---
x <- 1:6; dim(x) <- c(2, 3); dimnames(x) <- list(c("r1", "r2"), NULL); x
#---
x <- c(0.5, 100, 2, 3.25, NA, 7); dim(x) <- c(3, 2); x
#---
x <- 1:30; dim(x) <- c(15, 2); x
#---
x <- 1:24; dim(x) <- c(2, 3, 4); x
#---
x <- 1:8; dim(x) <- c(2, 2, 2); dimnames(x) <- list(c("a", "b"), c("c", "d"), c("e", "f")); x
#---
x <- numeric(0); dim(x) <- c(0, 3); x
#---
x <- 1:40; dim(x) <- c(2, 20); x
#---
x <- (1:40) * 1000; dim(x) <- c(2, 20); x
#---
factor_like <- structure(c(1L, 2L, 1L), levels = c("lo", "hi"), class = "factor"); factor_like
#---
structure(c(1L, NA, 2L), levels = c("a", "b"), class = "factor")
#---
structure(c(2L, 1L), levels = c("x", "y"), class = "factor", names = c("p", "q"))
#---
structure(1:2, levels = c("alpha", "beta", "gamma", "delta", "epsilon", "zeta", "eta", "theta", "iota", "kappa", "lambda", "mu", "nu", "xi", "omicron"), class = "factor")
#---
x <- 1:3; names(x) <- c("a", "b", "c"); x[2]
#---
invisible(5)
#---
(invisible(5))
#---
x <- c(a = 1)
#---
if (FALSE) 1
#---
for (i in 1:2) i
#---
f <- function() invisible(7); f(); (f()); y <- f(); y
#---
"a\tb"; "q\"uote"; "back\\slash"; "nl\nx"; "é"; "日本語"
#---
x <- c("short", "a much longer string here", "mid length"); x
#---
paste("word", 1:30)
#---
c(1.5, 2.25, 3)[c(TRUE, FALSE, TRUE)]
#---
TRUE; NA; NA_integer_; NA_real_; NA_character_; -0.5; -1L
#---
100000; 1e5 + 1; 0.0001; 0.00001234; 123456789; 1234567890123; 0.1 + 0.7
#---
c(1.5, 2, 3) * 1e6; c(1.5, 2, 3) * 1e-6
#---
c(10, 200, 3000) / 7
#---
-(1:3) / 3
#---
x <- c(0.1, 0.12, 0.123); x
#---
exp(0); sqrt(2); sqrt(16)

#---
array(0, c(2, 0), list(c("")))
#---
matrix(1:4, 2, dimnames = list("a"))
#---
matrix(1:4, 2, dimnames = list(c("a", "b", "c"), NULL))
#---
matrix(1:4, 2, dimnames = 5)
#---
c(-2.5, 0, 2^1000, NaN); c(-1, 2^1000, NA_real_, 7); -c(7, -2^1000, 0)
