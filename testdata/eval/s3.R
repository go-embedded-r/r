## S3 dispatch: UseMethod, NextMethod, group generics, internal generics. A line holding only #--- separates cells.
#---
area <- function(s) UseMethod("area"); area.square <- function(s) s$side^2; area.default <- function(s) NA; area(structure(list(side = 3), class = "square")); area(1)
#---
f <- function(x) UseMethod("f"); f.foo <- function(x) "foo"; f.default <- function(x) "default"; f(structure(1, class = "foo")); f(2); f("a")
#---
f <- function(x) UseMethod("f"); f.numeric <- function(x) "numeric"; f.integer <- function(x) "integer"; f.default <- function(x) "default"; f(1); f(1L); f("a"); f(TRUE)
#---
f <- function(x) UseMethod("f"); f.double <- function(x) "double"; f.numeric <- function(x) "numeric"; f(1)
#---
f <- function(x) UseMethod("f"); f.function <- function(x) "fn"; f.matrix <- function(x) "matrix"; f.array <- function(x) "array"; f(sum); f(function(x) x); m <- 1:4; dim(m) <- c(2, 2); f(m); a <- 1:8; dim(a) <- c(2, 2, 2); f(a)
#---
f <- function(x) UseMethod("f"); f.character <- function(x) "chr"; f.list <- function(x) "lst"; f.logical <- function(x) "lgl"; f.NULL <- function(x) "null"; f("a"); f(list()); f(NA); f(NULL)
#---
f <- function(x) UseMethod("f"); f(1)
#---
f <- function(x) UseMethod("f"); f(structure(1, class = c("a", "b")))
#---
f <- function(x) UseMethod("f"); f(1:2)
#---
f <- function(x) UseMethod("f"); f(structure(1:3, dim = 3L))
#---
f <- function(x, ...) UseMethod("f"); f.foo <- function(x, extra = 0, ...) x + extra; f(structure(1, class = "foo")); f(structure(1, class = "foo"), extra = 5)
#---
f <- function(x, y) UseMethod("f"); f.foo <- function(x, y) y; f(structure(1, class = "foo"), 42)
#---
f <- function(x) { x <- 99; UseMethod("f") }; f.default <- function(x) x; f(1)
#---
f <- function(x) { UseMethod("f"); "not reached" }; f.default <- function(x) "m"; f(1)
#---
f <- function(x) UseMethod("f", 5); f.numeric <- function(x) "dispatched on 5"; f("a")
#---
f <- function(x) UseMethod("f"); f.a <- function(x) c("a", NextMethod()); f.b <- function(x) c("b", NextMethod()); f.default <- function(x) "default"; f(structure(1, class = c("a", "b")))
#---
f <- function(x) UseMethod("f"); f.a <- function(x) NextMethod(); f.default <- function(x) "dflt"; f(structure(1, class = "a"))
#---
f <- function(x, ...) UseMethod("f"); f.a <- function(x, ...) NextMethod(); f.default <- function(x, ...) list(...); f(structure(1, class = "a"), 1, k = 2)
#---
f <- function(x, n) UseMethod("f"); f.a <- function(x, n) { n <- n * 10; NextMethod() }; f.default <- function(x, n) n; f(structure(1, class = "a"), 3)
#---
f <- function(x) UseMethod("f"); f.a <- function(x) { x <- "changed"; NextMethod() }; f.default <- function(x) x; f(structure(1, class = "a"))
#---
f <- function(x) UseMethod("f"); f.a <- function(x) .Generic; f(structure(1, class = "a"))
#---
print.money <- function(x, ...) cat("money\n"); l <- list(a = structure(1, class = "money"), b = 2); l
#---
toString <- function(x) UseMethod("toString"); toString.default <- function(x) "d"; toString(1)
#---
length.stack <- function(x) 99L; s <- structure(list(1, 2), class = "stack"); length(s); length(unclass(s))
#---
"[.myvec" <- function(x, i) { structure(unclass(x)[i], class = "myvec") }; v <- structure(1:5, class = "myvec"); class(v[2:3]); unclass(v[2:3])
#---
"[[.rec" <- function(x, i) "custom"; r <- structure(list(a = 1), class = "rec"); r[["a"]]; r[[1]]
#---
"$.rec" <- function(x, name) paste("field", name); r <- structure(list(a = 1), class = "rec"); r$a; r$zzz
#---
"[<-.myvec" <- function(x, i, value) { cat("assigning\n"); y <- unclass(x); y[i] <- value; structure(y, class = "myvec") }; v <- structure(1:3, class = "myvec"); v[2] <- 10L; unclass(v)
#---
"$<-.rec" <- function(x, name, value) { cat("set", name, "\n"); x }; r <- structure(list(a = 1), class = "rec"); r$b <- 5; names(r)
#---
"[[<-.rec" <- function(x, i, value) { cat("set2\n"); unclass(x) }; r <- structure(list(a = 1), class = "rec"); r[["a"]] <- 2; class(r)
#---
names.rec <- function(x) "fake"; r <- structure(list(a = 1), class = "rec"); names(r); names(unclass(r))
#---
"names<-.rec" <- function(x, value) { cat("names<-\n"); x }; r <- structure(list(a = 1), class = "rec"); names(r) <- "z"; names(unclass(r))
#---
as.character.rec <- function(x, ...) "REC"; as.character(structure(list(1), class = "rec"))
#---
c.rec <- function(...) "combined"; c(structure(1, class = "rec"), 2)
#---
Ops.money <- function(e1, e2) .Generic; a <- structure(5, class = "money"); a + 1; 1 - a; a * a; a < 2; a & TRUE; !a
#---
Ops.money <- function(e1, e2) { if (missing(e2)) "unary" else "binary" }; a <- structure(5, class = "money"); -a; a - 1
#---
"+.money" <- function(e1, e2) "plus"; Ops.money <- function(e1, e2) "ops"; a <- structure(5, class = "money"); a + 1; a - 1
#---
Ops.a <- function(e1, e2) "a"; Ops.b <- function(e1, e2) "b"; x <- structure(1, class = "a"); y <- structure(2, class = "b"); x + y
#---
Ops.a <- function(e1, e2) { r <- NextMethod(); unclass(r) * 10 }; x <- structure(1, class = "a"); x + 2
#---
Math.money <- function(x, ...) paste("Math", .Generic); a <- structure(-4, class = "money"); abs(a); sqrt(a); floor(a); exp(a); ceiling(a)
#---
Math.money <- function(x, ...) { structure(NextMethod(), class = "money") }; a <- structure(-4.5, class = "money"); unclass(abs(a)); class(abs(a))
#---
Summary.money <- function(..., na.rm = FALSE) paste("Summary", .Generic); a <- structure(1, class = "money"); sum(a); max(a, 2); min(a); prod(a); any(a); all(a)
#---
f <- function(x) UseMethod("f"); f.foo <- function(x) "foo"; inherits(structure(1, class = "foo"), "foo"); x <- structure(1, class = "foo"); f(x)
#---
f <- function(x) UseMethod("f"); g <- function(y) f(y); f.default <- function(x) "default"; f.k <- function(x) "k"; g(structure(1, class = "k"))
#---
local_f <- function() { h <- function(x) UseMethod("h"); h.default <- function(x) "local default"; h(1) }; local_f()
#---
f <- function(x) UseMethod("f"); f.default <- function(x) UseMethod("g"); f(1)
#---
f <- function() UseMethod(); f()
#---
f <- function() UseMethod("f"); f.default <- function() "no args"; f()
#---
f <- function(x) UseMethod("f"); f.foo <- function(x) stop("boom"); f(structure(1, class = "foo"))
#---
f <- function(x) UseMethod("f"); f.foo <- function(x) undefined_thing; f(structure(1, class = "foo"))
#---
format.rec <- function(x, ...) "formatted"; format.rec(1)
#---
x <- structure(list(), class = c("child", "parent")); class(x); inherits(x, "parent"); inherits(x, c("a", "parent"), which = TRUE)
