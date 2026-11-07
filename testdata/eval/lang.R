switch("b", a = 1, b = 2, c = 3)
#---
switch("z", a = 1, b = 2, 99)
#---
switch("z", a = 1, b = 2)
#---
switch("a", a = , b = "ab", c = "c")
#---
switch(2, "x", "y", "z"); switch(4, "x", "y", "z"); switch(1.9, "x", "y"); switch(0, "x")
#---
switch(NA_character_, a = 1, 2)
#---
switch("a", a = 1, 2, 3)
#---
switch(c("a", "b"), a = 1)
#---
switch(NULL, a = 1)
#---
switch("a")
#---
switch(2, a = "x", b = "y")
#---
f <- function(type) switch(type, num = 1, chr = "a", stop("unknown type: ", type)); f("num"); f("chr"); f("zzz")
#---
switch("b", a = 1, b = { x <- 5; x * 2 })
#---
x <- switch("q", a = 1); is.null(x)
#---
switch(1, )
#---
f <- function(x) substitute(x); f(a + b); f(1); f("s")
#---
f <- function(x, y) substitute(x + y); f(p, q * 2)
#---
f <- function(...) substitute(list(...)); f(a, b = 2)
#---
f <- function(x) { y <- 5; substitute(x + y) }; f(z)
#---
substitute(a + b, list(a = 1, b = quote(z)))
#---
e <- new.env(); assign("v", 10, e); substitute(v * 2, e)
#---
x <- 5; substitute(x)
#---
f <- function(x) deparse(substitute(x)); f(foo(bar)); f(1:10)
#---
substitute(x, 5)
#---
f <- function(a, b, c) nargs(); f(); f(1); f(1, 2, 3); f(b = 1)
#---
f <- function(...) ...length(); f(); f(1, 2, 3)
#---
f <- function(...) ...elt(2); f(10, 20, 30)
#---
f <- function(...) ...elt(5); f(1)
#---
f <- function(...) ...names(); f(a = 1, 2, b = 3); f(1, 2)
#---
...length()
#---
call("sum", 1, 2); eval(call("sum", 1, 2)); call("f", quote(x), a = 1)
#---
call(1)
#---
as.call(list(as.name("max"), 1, 5)); eval(as.call(list(as.name("max"), 1, 5))); as.call(list(quote(f), a = 1))
#---
as.call(list())
#---
as.call(5)
#---
f <- function(x) x * 2; forceAndCall(1, f, 21)
#---
.Primitive("sum")(1, 2, 3); .Primitive("+")(1, 2)
#---
.Primitive("nope")
#---
.subset(c(a = 1, b = 2), 2); .subset2(list(a = 1, b = 2), "b"); .subset(1:5, 2:3)
#---
x <- structure(list(1, 2), class = "foo"); `[.foo` <- function(x, i) "dispatched"; x[1]; .subset(x, 1)
#---
base::paste("a", "b"); base::sum(1:3)
#---
base::nope_not_there
#---
m <- matrix(c(1, 2, 3, 4, 5, 6), 2); n <- matrix(1:6, 3); m %*% n; dim(m %*% n)
#---
matrix(1:6, 2) %*% matrix(1:6, 2)
#---
m <- matrix(1:4, 2, dimnames = list(c("a", "b"), c("x", "y"))); m %*% m
#---
1:3 %*% matrix(1:6, 3); matrix(1:6, 2) %*% 1:3
#---
"a" %*% 1
#---
crossprod(matrix(1:6, 3)); tcrossprod(matrix(1:6, 3)); crossprod(1:3); tcrossprod(1:2); crossprod(matrix(1:4, 2), matrix(5:8, 2))
#---
crossprod(matrix(1:6, 3), matrix(1:4, 2))
#---
class(proc.time()); length(proc.time()); names(proc.time()); class(gc.time())

#---
m <- matrix(1:4, 2); m %*% m; m %*% c(1, 1); c(1, 1) %*% m; 1:3 %*% 1:3
#---
x <- c("a", "B", "b", "A", "_", "-", "1", "10", "2", "Z", "z", "\u00e9", "e", "E", "a b", "ab", "a-b", "a_b", ".", "~", "$", "aB", "Ab", "ab "); sort(x)
#---
"a" < "B"; "B" < "a"; "A" < "a"; "Z" < "a"; "e" < "\u00e9"; "10" < "2"; "ab" < "a b"; "a_b" < "ab"; "a" == "A"
#---
sort(c("x1", "x10", "x2", "X1")); max(c("a", "B")); min(c("a", "B")); ls_not <- 1
#---
a_B <- 1; Ab <- 2; aB <- 3; ab <- 4; ls()

#---
unique(c(1, 2, 2, 3, 1)); unique(c("b", "a", "b")); unique(c(1, NA, NA)); unique(list(1, 1, "a")); unique(c(a = 1, b = 1)); unique(NULL); unique(c(1L, 1L)); unique(c(0, -0)); unique(c(NaN, NA, NaN, NA))
#---
duplicated(c(1, 2, 1, 3, 2)); duplicated(c("a", "a")); duplicated(c(1, 2, 1), fromLast = TRUE); anyDuplicated(c(1, 2, 1)); anyDuplicated(1:3); anyDuplicated(NULL); duplicated(NULL)
#---
match(c("b", "z"), c("a", "b")); match(2, c(1, 2, 2)); match(NA, c(1, NA)); match(1L, 1.0); match("1", 1); match(c(1, 2), c(2, 1), nomatch = 0); match(NULL, 1); match(1, NULL); match(factor("a"), "a"); match(list(1, "a"), list("a"))
#---
c(1, 5) %in% c(1, 2); "a" %in% character(0); NA %in% c(1, NA); 1:3 %in% 2
#---
f <- factor(c("b", "a", "b", "c")); f; levels(f); nlevels(f); as.integer(f); class(f); unclass(f); table_not <- 1
#---
factor(c(3, 1, 2, 1)); factor(c("x", "y"), levels = c("y", "x")); factor(c("a", "b", NA)); factor(c("a", "z"), levels = c("a", "b")); factor(1:3, labels = c("one", "two", "three")); factor(c("a", "b"), ordered = TRUE); factor(character(0)); factor(c("a", "b", "a"), labels = "L")
#---
f <- factor(c("a", "b")); factor(f); factor(f, levels = "b"); as.factor(c(2, 1)); as.factor(f); factor(c(x = "a", y = "b")); levels(1); nlevels(1)
#---
format(factor(c("a", "bb", NA))); format(factor(character(0)))
#---
mode(1); mode(1L); mode("a"); mode(TRUE); mode(NULL); mode(list()); mode(sum); mode(quote(x)); mode(quote(x + 1)); mode(quote((x))); storage.mode(1L); storage.mode(sum); storage.mode(quote(x))
#---
as.symbol("abc"); as.name("x y"); as.name(quote(z)); typeof(as.name("q"))
#---
as.name(character(0))
#---
as.name("")
#---
base:::paste("a", "b")
