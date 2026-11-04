environmentName(globalenv()); environmentName(baseenv()); environmentName(emptyenv()); environmentName(new.env())
#---
f <- function() environment(); e <- f(); class(e); is.environment(e)
#---
f <- function(x) x; environmentName(environment(f)); environment(sum)
#---
e <- new.env(); assign("a", 10, envir = e); get("a", envir = e); exists("a", envir = e); exists("zz", envir = e)
#---
e <- new.env(); assign("a", 1, e); e$b <- 2; e[["c"]] <- 3; sort(ls(e)); mget(c("a", "b"), envir = e)
#---
e <- new.env(); local({ y <- 5; z <- y * 2 }, envir = e); get("z", e)
#---
x <- 1; f <- function() { x <- 2; g <- function() parent.frame(); identical(g(), environment()) }; f()
#---
f <- function() parent.frame(); identical(f(), globalenv())
#---
f <- function() { x <- "f"; g() }; g <- function() get("x", envir = parent.frame()); f()
#---
f <- function(...) sys.call(); f(1, b = 2)
#---
f <- function() sys.function(); class(f()); identical(f(), f)
#---
g <- function() sys.call(-1); f <- function() g(); f()
#---
get("nope")
#---
get("nope", mode = "function")
#---
exists("sum"); exists("sum", mode = "function"); exists("pi_not_defined_anywhere")
#---
get0("nope"); get0("nope", ifnotfound = "dflt"); a <- 3; get0("a")
#---
a <- 1; b <- 2; ls(); rm(a); ls(); rm(list = "b"); ls()
#---
rm(not_there)
#---
e <- new.env(); assign("v", 1, envir = e); rm("v", envir = e); exists("v", envir = e, inherits = FALSE)
#---
x <- 10; local({ x <- 20; x }); x
#---
local({ a <- 1; b <- 2; a + b })
#---
e <- new.env(); evalq(w <- 7, e); get("w", e)
#---
eval(quote(1 + 2)); eval(quote(x * 2), list(x = 21)); eval(parse_not_here <- 5)
#---
e <- new.env(); e$q <- 4; eval(quote(q + 1), e)
#---
f <- function() { v <- 3; eval(quote(v)) }; f()
#---
l <- list(a = 1, b = 2); e <- list2env(l); sort(ls(e)); get("b", e)
#---
e <- as.environment(list(k = 9)); get("k", e)
#---
as.environment(1)
#---
e1 <- new.env(); e2 <- new.env(parent = e1); assign("pp", 1, e1); identical(parent.env(e2), e1); exists("pp", envir = e2); exists("pp", envir = e2, inherits = FALSE)
#---
e <- new.env(); parent.env(e) <- emptyenv(); environmentName(parent.env(e))
#---
parent.env(emptyenv())
#---
f <- function(x) x + 1; environment(f) <- new.env(); class(environment(f))
#---
f <- function(x) { x + 1 }; body(f)
#---
counter <- function() { i <- 0; function() { i <<- i + 1; i } }; cn <- counter(); cn(); cn(); cn()
#---
make <- function() { v <- 1; list(get = function() v, set = function(n) v <<- n) }; o <- make(); o$set(42); o$get()
#---
e <- new.env(); assign("x", 1, e); assign("x", 2, e); get("x", e)
#---
assign("g1", 99); g1
#---
f <- function() { assign("created_by_f", 1, envir = globalenv()); invisible() }; f(); created_by_f
#---
assign(1, 2)
#---
topenv()
#---
identity(5); identity("a")
