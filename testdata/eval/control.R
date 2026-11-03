## Closures, promises, argument matching, control flow. A line holding only #--- separates cells.
#---
f <- function(x) x + 1; f(2)
#---
f <- function(x, y = 2) x * y; f(3); f(3, 4); f(y = 5, x = 2)
#---
f <- function(x, y) x - y; f(y = 1, 10); f(10, y = 1); f(x = 1, 10)
#---
f <- function(alpha, beta) alpha - beta; f(al = 10, be = 1)
#---
f <- function(value, verbose) c(value, verbose); f(v = 1, 2)
#---
f <- function(value, verbose) 1; f(v = 1, v = 2)
#---
f <- function(value, vector) 1; f(v = 1)
#---
f <- function(x) 1; f(1, 2)
#---
f <- function(x) 1; f(1, y = 2)
#---
f <- function(x) 1; f(1, 2, z = 3)
#---
f <- function(x, x2) 1; f(x = 1, x = 2)
#---
f <- function(x) x; f()
#---
f <- function(x) 1; f()
#---
f <- function(x, y) { if (missing(y)) "no y" else y }; f(1); f(1, 2)
#---
f <- function(x, y = 2) missing(y); f(1); f(1, 3)
#---
f <- function(x) missing(x); f(); f(1)
#---
f <- function(x) g(x); g <- function(y) missing(y); f()
#---
f <- function(x) { x <- 5; missing(x) }; f()
#---
f <- function(...) length(list(...)); f(); f(1); f(1, 2, 3)
#---
f <- function(...) c(...); f(1, 2, 3)
#---
f <- function(...) list(...); f(a = 1, 2)
#---
f <- function(x, ...) list(x, ...); f(1, 2, a = 3)
#---
f <- function(..., last) c(..., last); f(1, 2, last = 3)
#---
f <- function(..., x) x; f(1, 2)
#---
f <- function(a, ...) a; f(b = 1, 2)
#---
f <- function(x) { x; 1 }; f(1 + "a")
#---
f <- function(x) 1; f(1 + "a")
#---
f <- function(x, y = x * 2) y; f(4)
#---
f <- function(x, y = x * 2) { x <- 10; y }; f(4)
#---
f <- function(x = y, y = x) x; f()
#---
f <- function(x) { x; function() x }; g <- f(2); g()
#---
counter <- 0; f <- function() counter <<- counter + 1; f(); f(); counter
#---
f <- function() { x <- 1; g <- function() x <<- x + 1; g(); x }; f()
#---
x <- 10; f <- function() x; g <- function() { x <- 20; f() }; g()
#---
f <- function() { x <- 1; g <- function() x; x <- 2; g() }; f()
#---
add <- function(x) function(y) x + y; add(1)(2)
#---
f <- function(n) if (n <= 1) 1 else n * f(n - 1); f(10)
#---
f <- function(n) if (n < 2) n else f(n - 1) + f(n - 2); f(15)
#---
f <- function() return(5); f()
#---
f <- function() { return(1); 2 }; f()
#---
f <- function() { for (i in 1:10) if (i == 3) return(i); 0 }; f()
#---
f <- function() invisible(7); f()
#---
f <- function() invisible(7); y <- f(); y
#---
f <- function(x) { if (x) return("yes"); "no" }; f(TRUE); f(FALSE)
#---
f <- function() return(); f()
#---
f <- function() { g <- function() return(1); g(); 2 }; f()
#---
f <- function() { lapply_like <- function(g) g(); lapply_like(function() return(9)); 10 }; f()
#---
(function(x, y) x + y)(1, 2)
#---
if (TRUE) 1 else 2
#---
if (FALSE) 1 else 2
#---
if (FALSE) 1
#---
x <- if (FALSE) 1; x
#---
if (1) "one"; if (0) "zero" else "else"; if ("TRUE") "str"; if ("T") "T"
#---
if (NA) 1
#---
if (NULL) 1
#---
if (as.logical(NULL)) 1
#---
if (c(TRUE, FALSE)) 1
#---
if ("maybe") 1
#---
if (NA_integer_) 1
#---
if (NaN) 1
#---
if (list(TRUE)) 1
#---
x <- 5; if (x > 3) "big" else "small"
#---
if (TRUE) if (FALSE) 1 else 2
#---
for (i in 1:3) i
#---
x <- 0; for (i in 1:5) x <- x + i; x
#---
x <- 0; for (i in 1:10) { if (i %% 2 == 0) next; x <- x + i }; x
#---
x <- 0; for (i in 1:10) { if (i > 4) break; x <- x + i }; x
#---
for (i in NULL) 1 + "a"; "done"
#---
for (i in as.integer(NULL)) 1 + "a"; i_unset <- 1
#---
for (i in 1:3) NULL; i
#---
for (x in c(1.5, 2.5)) y <- x * 2; y
#---
for (s in c("a", "b")) last <- s; last
#---
for (e in list(1, "a", TRUE)) last <- e; last
#---
for (i in c(a = 1, b = 2)) last <- i; last
#---
r <- for (i in 1:3) i; r
#---
i <- 100; for (i in 1:2) NULL; i
#---
for (i in 1:2) for (j in 1:2) last <- i * 10 + j; last
#---
for (i in 1:3) { for (j in 1:3) { if (j == 2) break }; last <- c(i, j) }; last
#---
for (i in function() 1) 1
#---
i <- 0; while (i < 5) i <- i + 1; i
#---
i <- 0; while (TRUE) { i <- i + 1; if (i >= 7) break }; i
#---
i <- 0; n <- 0; while (i < 10) { i <- i + 1; if (i %% 2 == 1) next; n <- n + 1 }; n
#---
while (FALSE) 1
#---
r <- while (FALSE) 1; r
#---
i <- 0; repeat { i <- i + 1; if (i == 4) break }; i
#---
i <- 0; r <- repeat { i <- i + 1; if (i == 2) break }; r
#---
f <- function() { i <- 0; repeat { i <- i + 1; if (i > 2) return(i) } }; f()
#---
while (NA) 1
#---
while (c(TRUE, TRUE)) 1
#---
break
#---
next
#---
f <- function() break; f()
#---
f <- function() { for (i in 1:3) g <- function() break; 1 }; f()
#---
{ 1; 2; 3 }
#---
{}
#---
{ x <- 1 }
#---
(x <- 1)
#---
(invisible(2))
#---
quote(x + y)
#---
quote(f(x, y = 2))
#---
quote(x)
#---
quote(1)
#---
quote(if (a) b else c)
#---
quote(function(x, y = 1) x + y)
#---
quote("a")
#---
quote(x[1])
#---
quote(-x)
#---
quote(a <- b)
#---
quote()
#---
quote(a, b)
#---
class_of_quote <- is.call(quote(f(x))); class_of_quote
#---
identity_fn <- function(x) x; identity_fn(quote(a + b))
#---
f <- function(x) quote(x); f(1)
#---
"x" <- 5; x
#---
`my var` <- 3; `my var` + 1
#---
x <- y <- 2; x + y
#---
1 <- 2
#---
NULL <- 1
#---
f <- function() { x <- 1; g <- function() { x <<- 2; x }; c(g(), x) }; f()
#---
x <- 1; f <- function() x <<- 5; f(); x
#---
f <- function() newvar <<- 7; f(); newvar
#---
a_variable_that_does_not_exist
#---
f <- function() undefined_var; f()
#---
undefined_function(1)
#---
f <- function() undefined_function(1); f()
#---
(1)(2)
#---
"a"(1)
#---
x <- 5; x(1)
#---
f <- 1; f(2)
#---
sum <- 3; sum(1, 2)
#---
c <- function(...) "mine"; c(1, 2)
#---
`+` <- function(a, b) "plus"; 1 + 2
#---
`if` <- function(...) "no"; if (TRUE) 1
#---
f <- function(x) UseMethod("f"); 1
#---
1 + "a"
#---
"a" + 1
#---
-"a"
#---
-NULL
#---
!"a"
#---
list(1) + 1
#---
sum("a")
#---
f <- function(x) x + "a"; f(1)
#---
f <- function(x) { y <- x + "a"; y }; f(1)
#---
f <- function(n) f(n + 1); f(1)
#---
sqrt("a")
#---
TRUE & "a"
#---
1:3 + 1:2
#---
identity <- function(x) x; identity(1:2) + 1:3
