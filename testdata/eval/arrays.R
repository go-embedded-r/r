as.matrix(1:3)
#---
as.matrix(c(a = 1, b = 2))
#---
as.matrix(matrix(1:4, 2))
#---
as.matrix(letters[1:2])
#---
mat.or.vec(3, 1)
#---
mat.or.vec(2, 3)
#---
lower.tri(matrix(1:9, 3))
#---
lower.tri(matrix(1:9, 3), diag = TRUE)
#---
upper.tri(matrix(1:9, 3))
#---
upper.tri(matrix(1:12, 3), diag = TRUE)
#---
lower.tri(1:3)
#---
m <- matrix(1:9, 3); m[upper.tri(m)] <- 0L; m
#---
m <- matrix(0, 3, 3); diag(m) <- 1:3; m
#---
m <- matrix(0, 2, 3); diag(m) <- 5; m
#---
m <- matrix(0, 3, 2); diag(m) <- c(7, 8); m
#---
m <- matrix(0, 2, 2); diag(m) <- 1:3
#---
x <- 1:4; diag(x) <- 1
#---
m <- matrix(1:6, 2, dimnames = list(c("a", "b"), NULL)); diag(m) <- 0L; m
#---
rownames(matrix(1:4, 2))
#---
rownames(matrix(1:4, 2, dimnames = list(c("a", "b"), c("x", "y"))))
#---
colnames(matrix(1:4, 2, dimnames = list(c("a", "b"), c("x", "y"))))
#---
rownames(matrix(1:4, 2), do.NULL = FALSE)
#---
colnames(matrix(1:4, 2), do.NULL = FALSE, prefix = "V")
#---
rownames(1:3)
#---
colnames(1:3)
#---
rownames(1:3, do.NULL = FALSE)
#---
colnames(1:3, do.NULL = FALSE)
#---
rownames(array(1:3, 3, list(letters[1:3])))
#---
m <- matrix(1:4, 2); rownames(m) <- c("a", "b"); m
#---
m <- matrix(1:4, 2); colnames(m) <- c("x", "y"); m
#---
m <- matrix(1:4, 2, dimnames = list(c("a", "b"), c("x", "y"))); rownames(m) <- NULL; m
#---
m <- matrix(1:4, 2, dimnames = list(c("a", "b"), c("x", "y"))); colnames(m) <- NULL; m
#---
m <- matrix(1:4, 2); rownames(m) <- NULL; attributes(m)
#---
x <- 1:3; rownames(x) <- c("a", "b", "c")
#---
x <- 1:3; colnames(x) <- "a"
#---
m <- matrix(1:4, 2); rownames(m) <- c("a", "b", "c")
#---
m <- matrix(1:4, 2); colnames(m) <- c("a", "b", "c")
#---
m <- matrix(1:4, 2); rownames(m) <- 1:2; rownames(m)
#---
aperm(array(1:6, c(1, 2, 3)))
#---
aperm(array(1:24, 2:4), c(2, 1, 3))[, , 1]
#---
aperm(array(1:24, 2:4), c(3, 1, 2))[2, , ]
#---
aperm(matrix(1:6, 2))
#---
aperm(matrix(1:6, 2, dimnames = list(c("a", "b"), c("x", "y", "z"))))
#---
a <- array(1:24, 2:4, list(A = c("a1", "a2"), B = c("b1", "b2", "b3"), C = NULL)); dimnames(aperm(a, c("C", "A", "B")))
#---
aperm(array(1:6, c(1, 2, 3)), resize = FALSE)
#---
aperm(1:3)
#---
aperm(matrix(1:6, 2), c(1, 1))
#---
aperm(matrix(1:6, 2), 1)
#---
aperm(array(1:8, c(2, 2, 2)), c(3, 2, 1))
#---
apply(matrix(1:6, 2), 1, sum)
#---
apply(matrix(1:6, 2), 2, sum)
#---
apply(matrix(1:6, 2), c(1, 2), function(v) v * 2)
#---
apply(matrix(1:6, 2), 1, function(v) v * 2)
#---
apply(matrix(1:6, 2), 2, range)
#---
apply(matrix(1:6, 2, dimnames = list(c("a", "b"), c("x", "y", "z"))), 1, sum)
#---
apply(matrix(1:6, 2, dimnames = list(c("a", "b"), c("x", "y", "z"))), 2, function(v) v)
#---
apply(matrix(1:6, 2, dimnames = list(c("a", "b"), c("x", "y", "z"))), 1, function(v) v * 2)
#---
apply(matrix(1:6, 2, dimnames = list(c("a", "b"), c("x", "y", "z"))), c(1, 2), identity)
#---
apply(matrix(1:6, 2), 1, function(v) list(v))
#---
apply(matrix(1:6, 2), 1, function(v) v[v > 2])
#---
apply(matrix(1:6, 2), 2, function(v) if (v[1] > 2) 1 else c(1, 2))
#---
apply(matrix(1:6, 2), 1, paste, collapse = "-")
#---
apply(matrix(1:6, 2), 1, sum, simplify = FALSE)
#---
apply(matrix(1:6, 2), 2, function(v, k) v * k, k = 10)
#---
apply(array(1:24, 2:4), 3, sum)
#---
apply(array(1:24, 2:4), c(1, 3), sum)
#---
apply(array(1:24, 2:4), 1, function(m) dim(m))
#---
apply(array(1:24, 2:4), c(2, 3), function(v) v)
#---
apply(array(1:24, 2:4, list(A = c("a", "b"), B = c("p", "q", "r"), C = NULL)), "A", sum)
#---
apply(array(1:24, 2:4, list(A = c("a", "b"), B = c("p", "q", "r"), C = NULL)), c("C", "A"), sum)
#---
apply(matrix(numeric(0), 0, 3), 2, sum)
#---
apply(matrix(numeric(0), 3, 0), 1, sum)
#---
apply(1:3, 1, sum)
#---
apply(matrix(1:4, 2), 3, sum)
#---
apply(matrix(1:4, 2), "x", sum)
#---
apply(matrix(1:4, 2, dimnames = list(a = c("p", "q"), NULL)), "a", sum)
#---
apply(matrix(c(1, NA, 3, 4), 2), 1, sum, na.rm = TRUE)
#---
apply(matrix(1:6, 2), 1, function(v) NULL)
#---
apply(matrix(1:6, 2), 1, mean)
#---
apply(matrix(c("a", "b", "c", "d"), 2), 2, toupper)
#---
apply(matrix(c(TRUE, FALSE, TRUE, TRUE), 2), 1, all)
#---
apply(matrix(1:6, 3), 1, function(v) c(a = v[1], b = v[2]))
#---
apply(matrix(1:6, 2), 2, sort, decreasing = TRUE)
#---
sweep(matrix(1:6, 2), 2, c(1, 2, 3))
#---
sweep(matrix(1:6, 2), 1, c(10, 20), "+")
#---
sweep(matrix(1:6, 2), 2, colMeans(matrix(1:6, 2)))
#---
sweep(matrix(1:6, 2), 1, 1:2, "*")
#---
sweep(matrix(1:6, 2), 2, 1:3, FUN = function(a, b) a - b)
#---
sweep(matrix(1:6, 2, dimnames = list(c("a", "b"), NULL)), 1, 1:2)
#---
sweep(matrix(1:6, 2), 2, 1:2)
#---
sweep(matrix(1:6, 2), 2, 1:4)
#---
sweep(matrix(1:6, 2), 2, 1:2, check.margin = FALSE)
#---
sweep(array(1:24, 2:4), 3, 1:4)[, , 4]
#---
sweep(array(1:24, 2:4), c(1, 3), matrix(1:8, 2))[2, , ]
#---
sweep(matrix(1:6, 2), 1, 1:2, "/")
#---
scale(matrix(c(1, 2, 3, 4, 5, 6), 3))
#---
scale(matrix(c(1, 2, 3, 4, 5, 6), 3), center = FALSE)
#---
scale(matrix(c(1, 2, 3, 4, 5, 6), 3), scale = FALSE)
#---
scale(matrix(c(1, 2, 3, 4, 5, 6), 3), center = c(1, 1), scale = c(2, 2))
#---
scale(matrix(c(1, 2, 3, 4, 5, 6), 3), center = FALSE, scale = FALSE)
#---
scale(1:4)
#---
scale(matrix(c(1, NA, 3, 4, 5, 6), 3))
#---
scale(matrix(1:6, 3, dimnames = list(NULL, c("a", "b"))))
#---
scale(matrix(1:6, 3), center = 1)
#---
scale(matrix(1:6, 3), scale = 1:3)
#---
prop.table(matrix(1:4, 2))
#---
prop.table(matrix(1:4, 2), 1)
#---
prop.table(matrix(1:4, 2), 2)
#---
proportions(c(a = 1, b = 3))
#---
proportions(matrix(1:4, 2), 1)
#---
prop.table(1:4)
#---
margin.table(matrix(1:4, 2))
#---
margin.table(matrix(1:4, 2), 1)
#---
margin.table(matrix(1:4, 2), 2)
#---
marginSums(matrix(1:4, 2, dimnames = list(a = c("p", "q"), b = c("x", "y"))), 1)
#---
marginSums(array(1:24, 2:4), c(1, 3))
#---
marginSums(array(1:24, 2:4), 3)
#---
marginSums(1:3)
#---
1:3 %o% 1:2
#---
c(a = 1, b = 2) %o% c(x = 1, y = 2)
#---
force(5)
#---
f <- function(x) { force(x); function() x }; f(3)()
#---
NULL %||% 4
#---
2 %||% 4
#---
c(1, 5) %notin% c(5, 6)
#---
"a" %notin% character(0)
#---
grepv("a", c("abc", "xyz", "bca"))
#---
grepv("a", c("abc", "xyz", "bca"), invert = TRUE)
#---
grepv("A", c("abc", "xyz"), ignore.case = TRUE)
#---
grepv("a", c(x = "abc", y = "q"))
#---
f <- function() { x <- 1; g <- function() eval.parent(quote(x)); g() }; f()
#---
f <- function() { x <- 2; g <- function() eval.parent(quote(x), 1); g() }; f()
#---
replace(1:5, 2, 10L)
#---
replace(c(a = 1, b = 2), "b", 9)
#---
replace(1:5, c(1, 5), c(0L, 0L))
#---
replace(1:3, 5, 1L)
#---
replace(1:5, 1:5 > 3, NA)
#---
replace(list(1, 2), 2, list("a"))
#---
set.seed <- NULL; length(replicate(5, 1))
#---
replicate(3, 1 + 1)
#---
replicate(2, c(1, 2))
#---
replicate(2, "a")
#---
replicate(0, 1)
#---
replicate(2, list(1), simplify = FALSE)
#---
replicate(3, c(a = 1, b = 2))
#---
x <- 0; replicate(3, x <- x + 1); x
#---
replicate(2, NULL)
#---
f <- function(k) replicate(2, k * 2); f(4)
#---
replicate(2, matrix(1:4, 2))
#---
replicate(2, matrix(1:4, 2), simplify = FALSE)
#---
replicate(2, matrix(1:4, 2), simplify = TRUE)
#---
replicate(-1, 1)
#---
setequal(1:3, c(3, 2, 1))
#---
setequal(1:3, 1:4)
#---
setequal(c(1, 1, 2), c(2, 1))
#---
setequal(character(0), character(0))
#---
setequal(NULL, 1)
#---
setequal(c(1, NA), c(NA, 1))
#---
setequal(list(1, 2), list(2, 1))
#---
getElement(list(a = 1, b = 2), "b")
#---
getElement(list(abc = 1), "ab")
#---
getElement(c(a = 1), "a")
#---
getElement(list(a = 1), "z")
#---
factorial(5)
#---
factorial(0:5)
#---
signif(factorial(c(0.5, 10)), 10)
#---
factorial(-1)
#---
lfactorial(10)
#---
signif(lfactorial(c(0, 1, 100)), 10)
#---
zapsmall(c(1, 1e-20))
#---
zapsmall(c(1.00000001, 2e-9), digits = 4)
#---
zapsmall(c(0, 0))
#---
zapsmall(c(NA, 1e-20, 1))
#---
zapsmall(c(NA_real_, NA_real_))
#---
zapsmall(c(1e-300, 1))
#---
zapsmall(c(123456.7, 0.001), digits = 3)
#---
zapsmall(numeric(0))
#---
zapsmall(c(1, 1e-12), mFUN = function(x, ina) 100)
#---
zapsmall(c(1e-3, 1), min.d = 5, digits = 3)
#---
zapsmall(c(1, -1e-20, Inf))
#---
sequence(c(3, 2))
#---
sequence(c(3, 2), from = 2L)
#---
sequence(c(3, 2), by = 2L)
#---
sequence(c(3, 2), from = c(10L, 20L), by = c(1L, -1L))
#---
sequence(integer(0))
#---
sequence(0)
#---
sequence(3:1)
#---
sequence(c(2, 2), from = 1:4)
#---
sequence(2.0)
#---
unclass(rle(c(1, 1, 2, 2, 2, 3)))
#---
unclass(rle(c("a", "a", "b")))
#---
unclass(rle(c(TRUE, TRUE, FALSE)))
#---
unclass(rle(integer(0)))
#---
unclass(rle(c(NA, NA, 1)))
#---
unclass(rle(c(1, NA, NA)))
#---
unclass(rle(1))
#---
unclass(rle(c(a = 1, b = 1, c = 2)))
#---
class(rle(1:3))
#---
rle(matrix(1:4, 2))
#---
unclass(rle(factor(c("a", "a"))))
#---
inverse.rle(rle(c(1, 1, 2)))
#---
inverse.rle(rle(c("a", "a", "b", "b", "b")))
#---
inverse.rle(list(lengths = 1:2, values = c("x", "y")))
#---
inverse.rle(list(lengths = 1:2))
#---
inverse.rle(list(lengths = 1:2, values = 1))
#---
inverse.rle(rle(integer(0)))
#---
x <- c(4, 4, 4, 1, 1, 9); identical(inverse.rle(rle(x)), x)
#---
split(1:6, c(1, 2, 1, 2, 1, 2))
#---
split(c("a", "b", "c"), c("u", "v", "u"))
#---
split(1:4, factor(c("a", "b", "a", "b"), levels = c("b", "a", "z")))
#---
split(1:4, factor(c("a", "b", "a", "b"), levels = c("b", "a", "z")), drop = TRUE)
#---
split(1:4, c(1, NA, 1, 2))
#---
split(c(a = 1, b = 2, c = 3), c("x", "y", "x"))
#---
split(1:4, 1:2)
#---
split(1:5, 1:2)
#---
split(1:3, 1:2)
#---
split(integer(0), character(0))
#---
split(list(1, "a", TRUE), c(1, 2, 1))
#---
split(as.Date(c("2024-01-01", "2024-02-01", "2024-01-15")), c(1, 2, 1))
#---
split(1:3, c("b", "a", "b"))
#---
split(1:3, factor(c("b", "a", "b"), levels = c("b", "a")))
#---
split(1:6, list(c(1, 1, 2, 2, 3, 3)))
#---
split(1:4, c(TRUE, FALSE, TRUE, FALSE))
#---
split(1:3, c(0.5, 1.5, 0.5))
#---
split(1:3, list(1:3, 1:3))
#---
split(factor(c("u", "v", "u")), c(1, 2, 1))
#---
x <- 1:6; split(x, c(1, 2)) <- list(c(10L, 11L, 12L), c(20L, 21L, 22L)); x
#---
x <- 1:4; split(x, c(1, 2, 1, 2)) <- 0L; x
#---
x <- c("a", "b", "c"); split(x, c(1, 1, 2)) <- list("p", "q"); x
#---
x <- 1:4; split(x, c("a", "b", "a", "b")) <- list(c(1L, 2L), 9L); x
#---
unsplit(split(1:6, c(1, 2, 1, 2, 1, 2)), c(1, 2, 1, 2, 1, 2))
#---
unsplit(split(c("a", "b", "c"), c("x", "y", "x")), c("x", "y", "x"))
#---
unsplit(list(`1` = 1:2, `2` = 3:4), c(1, 1, 2, 2))
#---
unsplit(list(a = c(5, 6), b = 7), c("a", "b", "a"))
#---
tapply(1:6, c(1, 1, 2, 2, 3, 3), sum)
#---
tapply(c(1, 2, 3, 4), c("a", "b", "a", "b"), mean)
#---
tapply(c(1, 2, 3, 4), c("a", "b", "a", "b"), function(v) v)
#---
tapply(c(1, 2, 3, 4), c("a", "b", "a", "b"), range)
#---
tapply(c(1, 2, 3, 4), factor(c("a", "b", "a", "b"), levels = c("a", "b", "c")), sum)
#---
tapply(c(1, 2, 3, 4), factor(c("a", "b", "a", "b"), levels = c("a", "b", "c")), sum, default = 0)
#---
tapply(c(1, 2, 3, 4), factor(c("a", "b", "a", "b"), levels = c("a", "b", "c")), function(v) v)
#---
tapply(1:4, c("a", "b", "a", "b"))
#---
tapply(1:4, list(c("a", "a", "b", "b"), c("x", "y", "x", "y")), sum)
#---
tapply(1:6, list(g = c(1, 1, 2, 2, 3, 3), h = c("a", "b", "a", "b", "a", "b")), sum)
#---
tapply(1:4, list(c("a", "a", "b", "b"), c("x", "x", "x", "y")), sum)
#---
tapply(1:4, list(c("a", "a", "b", "b"), c("x", "x", "x", "y")), sum, default = -1)
#---
tapply(1:4, c("a", "b", "a", "b"), sum, simplify = FALSE)
#---
tapply(c(1.5, 2.5, 3.5), c("u", "v", "u"), max)
#---
tapply(1:4, c(NA, "a", "a", "b"), sum)
#---
tapply(c(a = 1, b = 2, c = 3), c("x", "y", "x"), sum)
#---
tapply(1:4, c("a", "b", "a", "b"), paste, collapse = "")
#---
tapply(1:4, c("a", "b", "a", "b"), function(v) NULL)
#---
tapply(1:4, c("a", "b", "a"), sum)
#---
tapply(1:4, list(c("a", "b", "a", "b"), 1:3), sum)
#---
tapply(character(0), character(0), length)
#---
tapply(c(TRUE, FALSE, TRUE), c("a", "a", "b"), any)
#---
tapply(c("p", "q", "r"), c("a", "a", "b"), toupper)
#---
tapply(1:6, c(2, 2, 1, 1, 3, 3), length)
#---
bquote(x + .(1 + 2))
#---
y <- 5; bquote(f(.(y), z))
#---
bquote(.(quote(a)) * 2)
#---
bquote(x)
#---
bquote(a + b)
#---
bquote(g(.(c(1, 2))))
#---
bquote(function(x) x + .(3))
#---
e <- new.env(); assign("w", 9, e); bquote(h(.(w)), e)
#---
bquote(.(1))
#---
bquote(f(..(c(1, 2))), splice = TRUE)
#---
bquote(f(.(1), .(2 + 2), k = .(3)))
#---
matrix(1:4, 2, dimnames = list(a = c("x", "yy"), b = c("p", "q")))
#---
matrix(1:4, 2, dimnames = list(rowname = c("x", "y"), b = c("p", "q")))
#---
matrix(1:4, 2, dimnames = list(a = c("x", "y"), colname = c("p", "q")))
#---
matrix(1:4, 2, dimnames = list(a = c("x", "y"), NULL))
#---
matrix(1:4, 2, dimnames = list(NULL, b = c("p", "q")))
#---
matrix(c("a", "bb", "c", "d"), 2, dimnames = list(r = c("x", "y"), cc = c("p", "q")))
#---
matrix(c(1.5, 2, 3, 4), 2, dimnames = list(r = c("x", "y"), cc = c("p", "q")))
#---
array(1:8, c(2, 2, 2), dimnames = list(a = c("x", "y"), b = c("p", "q"), c = c("u", "v")))
#---
matrix(1:4, 2, dimnames = list(a = NULL, b = c("p", "q")))
#---
matrix(1:4, 2, dimnames = list(a = c("x", "y"), b = NULL))
#---
matrix(1:4, 2, dimnames = list(c("", ""), c("p", "q")))
#---
matrix(1:4, 2, dimnames = list(a = c("", ""), b = c("p", "q")))
#---
matrix(c(TRUE, FALSE, NA, TRUE), 2, dimnames = list(r = c("x", "y"), cc = c("p", "q")))
#---
matrix(numeric(0), 0, 2, dimnames = list(a = NULL, b = c("p", "q")))
#---
array(1:3, 3, dimnames = list(k = c("a", "b", "c")))
#---
array(1:3, 3, dimnames = list(c("a", "b", "c")))
#---
array(list(1, "a"), 2, dimnames = list(c("a", "b")))
#---
array(list(1, "a"), 2, dimnames = list(k = c("a", "b")))
#---
array(1:3, 3)
#---
x <- array(1:3, 3, dimnames = list(c("a", "b", "c"))); list(x[2], names(x), x[["b"]], x["c"])
#---
array(1:2, 2, dimnames = list(k = c("a", NA)))
#---
matrix(1:2, 1, dimnames = list(a = "x", b = c("p", "q")))
#---
is.vector(1:3)
#---
is.vector(c(a = 1))
#---
is.vector(matrix(1:4, 2))
#---
is.vector(list(1))
#---
is.vector(NULL)
#---
is.vector(factor("a"))
#---
is.vector(1:3, "numeric")
#---
is.vector(1:3, "integer")
#---
is.vector(1:3, "double")
#---
is.vector(1.5, "double")
#---
is.vector("a", "character")
#---
is.vector(list(1), "list")
#---
is.vector(list(1), "numeric")
#---
is.vector(quote(a))
#---
is.vector(quote(a), "symbol")
#---
is.vector(quote(f(x)))
#---
is.vector(function(x) x)
#---
is.vector(structure(1:2, foo = "x"))
#---
is.vector(TRUE, "logical")
#---
is.vector(1:2, "any")
#---
q <- quote(f(a, b = 2)); q[[1]]
#---
q <- quote(f(a, b = 2)); q[[3]]
#---
q <- quote(f(a, b = 2)); q[-1]
#---
q <- quote(f(a, b = 2)); q[1]
#---
q <- quote(f(a, b = 2)); q[2:3]
#---
q <- quote(f(a, b = 2)); q[[4]]
#---
q <- quote(f(a, b = 2)); length(q)
#---
q <- quote(f(a, b = 2)); as.list(q)
#---
sapply(1:2, function(i) matrix(i, 2, 2), simplify = "array")
#---
sapply(1:2, function(i) matrix(i, 2, 2), simplify = TRUE)
#---
sapply(c(a = 1, b = 2), function(i) matrix(i, 2, 2), simplify = "array")
#---
sapply(1:3, function(i) i, simplify = "array")
#---
sapply(1:2, function(i) c(i, i), simplify = "array")
#---
array(NULL)
#---
array(NULL, c(2, 2))
#---
sapply(1:3, function(i) x <- i)
#---
lapply(1:2, function(i) y <- i)
#---
Map(function(a) z <- a, 1:2)
#---
Reduce(function(a, b) w <- a + b, 1:3)
#---
Filter(function(a) { v <- a; a > 1 }, 1:3)
#---
aperm.default(matrix(1:6, 2))
#---
as.matrix.default(1:2)
#---
row(matrix(1:6, 2))
#---
col(matrix(1:6, 2))
#---
row(1:3)
#---
col(array(1:8, c(2, 2, 2)))
#---
scale.default(matrix(1:6, 3))
#---
sequence.default(c(2, 3))
#---
split.default(1:4, c(1, 2, 1, 2))
#---
x <- 1:4; `split<-.default`(x, c(1, 2, 1, 2), value = list(10L, 20L))
#---
unclass(row(matrix(1:4, 2, dimnames = list(c("a", "b"), NULL)), as.factor = TRUE))
#---
unclass(col(matrix(1:4, 2, dimnames = list(NULL, c("p", "q"))), as.factor = TRUE))
#---
m <- matrix(1:6, 2, dimnames = list(r = c("a", "b"), c = c("x", "y", "z"))); names(dimnames(apply(m, 1, function(v) c(lo = 1, hi = 2, mid = 3))))
#---
m <- matrix(1:6, 2, dimnames = list(r = c("a", "b"), c = c("x", "y", "z"))); names(dimnames(apply(m, 2, function(v) c(lo = 1, hi = 2))))
#---
m <- matrix(1:6, 2, dimnames = list(r = c("a", "b"), c = c("x", "y", "z"))); names(dimnames(apply(m, 2, function(v) c(v, 9))))
#---
m <- matrix(1:6, 2, dimnames = list(r = c("a", "b"), c = c("x", "y", "z"))); dimnames(apply(m, 2, function(v) unname(v * 2)))
#---
m <- matrix(1:6, 2, dimnames = list(r = c("a", "b"), c = c("x", "y", "z"))); apply(m, 1, rev)
#---
m <- matrix(1:6, 2, dimnames = list(r = c("a", "b"), c = c("x", "y", "z"))); apply(m, 1, function(v) if (v[1] > 1) list(v) else 0)
#---
apply(matrix(c(-3, 1, 4, -4), 2), c(1, 2), function(v) v[v > 0])
#---
apply(matrix(1:4, 2), c(1, 2), function(v) c(1, 2))
#---
identical(quote(f(x)), quote(f))
#---
identical(quote(f), quote(f(x)))
#---
bquote(0()())
