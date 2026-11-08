e <- new.env(); assign("a", 1, e); lockBinding("a", e); bindingIsLocked("a", e); bindingIsLocked(quote(a), e)
#---
e <- new.env(); assign("a", 1, e); lockBinding("a", e); assign("a", 2, e)
#---
e <- new.env(); assign("a", 1, e); lockBinding("a", e); e$a <- 2
#---
e <- new.env(); assign("a", 1, e); lockBinding("a", e); e[["a"]] <- 2
#---
e <- new.env(); assign("a", 1, e); lockBinding("a", e); unlockBinding("a", e); e$a <- 3; e$a; bindingIsLocked("a", e)
#---
e <- new.env(); lockBinding("nope", e)
#---
e <- new.env(); unlockBinding("nope", e)
#---
e <- new.env(); bindingIsLocked("nope", e)
#---
e <- new.env(); assign("a", 1, e); lockBinding(1, e)
#---
lockBinding("a", 5)
#---
f <- function() { a <- 1; lockBinding("a", environment()); a <- 2 }; f()
#---
f <- function() { a <- 1; lockBinding("a", environment()); a[1] <- 2 }; f()
#---
f <- function() { a <- 1; lockBinding("a", environment()); g <- function() a <<- 5; g() }; f()
#---
f <- function() { a <- 1; b <- 2; lockBinding("a", environment()); b <- 3; b }; f()
#---
e <- new.env(); assign("a", 1, e); lockBinding("a", e); list2env(list(a = 5), e)
#---
e <- new.env(); lockEnvironment(e); environmentIsLocked(e); environmentIsLocked(new.env())
#---
e <- new.env(); lockEnvironment(e); assign("b", 1, e)
#---
e <- new.env(); lockEnvironment(e); e$b <- 1
#---
e <- new.env(); lockEnvironment(e); e[["b"]] <- 1
#---
e <- new.env(); lockEnvironment(e); list2env(list(b = 5), e)
#---
e <- new.env(); assign("a", 1, e); lockEnvironment(e); e$a <- 2; e$a
#---
e <- new.env(); assign("a", 1, e); lockEnvironment(e, bindings = TRUE); e$a <- 2
#---
e <- new.env(); assign("a", 1, e); lockEnvironment(e, TRUE); bindingIsLocked("a", e)
#---
e <- new.env(); assign("a", 1, e); lockEnvironment(e, FALSE); bindingIsLocked("a", e)
#---
e <- new.env(); assign("a", 1, e); lockEnvironment(e); rm("a", envir = e)
#---
e <- new.env(); lockEnvironment(e); rm(list = "zz", envir = e)
#---
e <- new.env(); lockEnvironment(e); delayedAssign("k", 1, assign.env = e)
#---
e <- new.env(); lockEnvironment(e); makeActiveBinding("q", function() 1, e)
#---
lockEnvironment(5)
#---
environmentIsLocked(5)
#---
x <- lockEnvironment(new.env()); is.null(x)
#---
withVisible(lockBinding("a", { e <- new.env(); assign("a", 1, e); e }))$visible
#---
withVisible(lockEnvironment(new.env()))$visible
#---
e <- new.env(); makeActiveBinding("z", function() 42, e); e$z; get("z", e); bindingIsActive("z", e)
#---
e <- new.env(); makeActiveBinding("z", function() 42, e); bindingIsActive("nope", e)
#---
e <- new.env(); assign("r", 1, e); bindingIsActive("r", e)
#---
e <- new.env(); assign("r", 1, e); makeActiveBinding("r", function() 1, e)
#---
makeActiveBinding("z", 5, new.env())
#---
e <- new.env(); makeActiveBinding("z", function() 42, e); f <- activeBindingFunction("z", e); f()
#---
e <- new.env(); assign("r", 1, e); activeBindingFunction("r", e)
#---
e <- new.env(); makeActiveBinding("z", function(v) if (missing(v)) "read" else cat("write", v, "\n"), e); e$z; e$z <- 5; e[["z"]]; get0("z", e)
#---
f <- local({ v <- 0; function(x) { if (missing(x)) v else { v <<- x; invisible() } } }); e <- new.env(); makeActiveBinding("ab", f, e); e$ab <- 5; e$ab; assign("ab", 6, e); e$ab
#---
n <- 0; makeActiveBinding("tick", function() { n <<- n + 1; n }, environment()); tick; tick; tick
#---
makeActiveBinding("two", function() 2, environment()); two * 21
#---
e <- new.env(); makeActiveBinding("z", function() 42, e); exists("z", e); exists("z", envir = e, inherits = FALSE); sort(ls(e))
#---
e <- new.env(); makeActiveBinding("z", function() 42, e); mget("z", envir = e)
#---
e <- new.env(); makeActiveBinding("z", function() 42, e); as.list(e)
#---
e <- new.env(); makeActiveBinding("z", function() 42, e); rm("z", envir = e); exists("z", e, inherits = FALSE)
#---
e <- new.env(); makeActiveBinding("z", function() 42, e); lockBinding("z", e); makeActiveBinding("z", function() 1, e)
#---
e <- new.env(); makeActiveBinding("z", function() 42, e); local(z + 1, e)
#---
delayedAssign("lazy", {cat("forced\n"); 7}); cat("before\n"); lazy; lazy
#---
y0 <- 10; delayedAssign("lz", y0 * 2); y0 <- 20; lz
#---
e <- new.env(); delayedAssign("x1", 1 + 1, assign.env = e); get("x1", e)
#---
e <- new.env(); delayedAssign("x1", 1 + 1, assign.env = e); exists("x1", e, inherits = FALSE)
#---
e <- new.env(); delayedAssign("x2", stop("boom"), assign.env = e); exists("x2", e, inherits = FALSE)
#---
e <- new.env(); delayedAssign("x2", {cat("evaluated\n"); 1}, assign.env = e); ls(e); e$x2; e$x2
#---
ev <- new.env(); assign("p", 5, ev); delayedAssign("viaenv", p + 1, eval.env = ev); viaenv
#---
delayedAssign(1, 1)
#---
withVisible(delayedAssign("q", 1))$visible
#---
e <- new.env(); delayedAssign("d", 3, assign.env = e); mget("d", envir = e)
#---
e <- new.env(); delayedAssign("d", 3, assign.env = e); as.list(e)
#---
delayedAssign("w", 1); w <- 2; w
#---
f <- function() { delayedAssign("loc", {cat("in f\n"); 9}); loc + 1 }; f()
#---
e <- new.env(); assign("b", 2, e); assign("a", 1, e); assign(".h", 3, e); sort(names(as.list(e))); sort(names(as.list(e, all.names = TRUE)))
#---
e <- new.env(); assign("b", 2, e); assign("a", 1, e); assign(".h", 3, e); as.list(e, sorted = TRUE)
#---
e <- new.env(); assign("b", 2, e); assign("a", 1, e); assign(".h", 3, e); as.list(e, all.names = TRUE, sorted = TRUE)
#---
e <- new.env(); assign("b", 2, e); assign("a", 1, e); as.list.environment(e, sorted = TRUE)
#---
e <- new.env(); assign("B", 2, e); assign("a", 1, e); assign("c", 3, e); assign("A", 0, e); names(as.list(e, sorted = TRUE))
#---
as.list(new.env())
#---
as.list(emptyenv())
#---
e <- new.env(); assign("a", 1, e); as.list(e, all.names = NA)
#---
e <- new.env(); assign("a", 1, e); as.list(e, sorted = 5)
#---
as.list(1:3)
#---
as.list(c(a = 1, b = 2))
#---
as.list(list(1, 2))
#---
as.list("x")
#---
as.list(NULL)
#---
as.list(quote(f(x, y)))
#---
as.list(quote(f(a = x, y)))
#---
as.list(factor(c("u", "v")))
#---
as.list(factor(c(a = "u", b = "v")))
#---
as.list(function(x, y = 2) x + y)
#---
as.list(5, 6)
#---
as.list(1:2, 3)
#---
as.list.default(1:2)
#---
as.list.default(new.env())
#---
as.list(TRUE)
#---
as.list(c(1.5, NA))
#---
as.list(matrix(1:4, 2))
#---
x <- structure(list(1, 2), class = "foo"); class(as.list(x))
#---
as.list.foo <- function(x, ...) "method"; as.list(structure(1, class = "foo"))
#---
e <- new.env(); assign("a", 1, e); assign("b", 2, e); r <- eapply(e, function(v) v * 10); r[order(names(r))]
#---
e <- new.env(); assign("a", 1, e); assign(".h", 2, e); names(eapply(e, identity)); names(eapply(e, identity, all.names = TRUE)[1:0])
#---
e <- new.env(); assign("a", 1, e); eapply(e, function(v, k) v + k, 100)
#---
e <- new.env(); assign("a", 1, e); eapply(e, identity, USE.NAMES = FALSE)
#---
eapply(new.env(), identity)
#---
with(list(a = 1, b = 2), a + b)
#---
with(list(a = 1), a + zz)
#---
zz <- 100; with(list(a = 1), a + zz)
#---
e <- new.env(); assign("q", 7, e); with(e, q * 2)
#---
with(NULL, 1 + 1)
#---
f <- function() { loc <- 5; with(list(a = 1), a + loc) }; f()
#---
with(c(a = 1, b = 2), a + b)
#---
e <- new.env(); assign("a", 1, e); assign("b", 2, e); mget(c("a", "b"), envir = e)
#---
e <- new.env(); assign("a", 1, e); mget(c("a", "zz"), envir = e)
#---
e <- new.env(); assign("a", 1, e); mget(c("a", "zz"), envir = e, ifnotfound = list(NA))
#---
e <- new.env(); assign("a", 1, e); mget(c("a", "zz"), envir = e, ifnotfound = NA)
#---
e <- new.env(); assign("a", 1, e); mget(c("a", "zz"), envir = e, ifnotfound = list(function(x) paste("none", x)))
#---
e <- new.env(); assign("a", 1, e); mget(c("zz", "yy"), envir = e, ifnotfound = list("p", "q"))
#---
e <- new.env(); assign("a", 1, e); mget("a", envir = e, mode = "function")
#---
e <- new.env(); names(mget("sum", envir = e, mode = "function", inherits = TRUE))
#---
e <- new.env(); mget("sum", envir = e, inherits = FALSE)
#---
e <- new.env(); mget("pi", envir = e, inherits = TRUE)
#---
e <- new.env(); mget(1, envir = e)
#---
e <- new.env(); assign("a", 1, e); assign("b", 2, e); mget(c("a", "b"), envir = e, ifnotfound = list(1, 2, 3))
#---
e <- new.env(); mget(character(0), envir = e)
#---
e <- new.env(); assign("a", 1, e); mget(c("a", "a"), envir = e)
#---
mget("not_there_at_all")
#---
e <- new.env(); assign("a", 1, e); assign("b", "x", e); mget(c("a", "b"), envir = e, mode = c("numeric", "character"))
#---
e <- new.env(); assign("a", 1, e); assign("b", 2, e); mget(c("a", "b"), envir = e, mode = c("numeric", "character"))
#---
x <- 1; y <- "s"; mget(c("x", "y"))
#---
f <- function() { loc <- 3; mget("loc") }; f()
#---
e <- new.env(); assign("a", 1, e); get("a", envir = e, mode = "numeric")
#---
e <- new.env(); assign("a", 1, e); get("a", envir = e, mode = "character")
#---
e <- new.env(); assign("a", 1, e); get("a", e, inherits = FALSE); get("sum", e, inherits = FALSE)
#---
e <- new.env(); get("sum", e)(1:3)
#---
get(c("a", "b"))
#---
get(1)
#---
get(character(0))
#---
get("sum", mode = "function")(2, 3)
#---
e <- new.env(); assign("a", 1, e); assign("a", "chr", new.env(parent = e)); get("a", mode = "numeric", envir = new.env(parent = e))
#---
e <- new.env(); get0("zz", envir = e, inherits = FALSE); get0("sum", envir = e, inherits = FALSE, ifnotfound = "none")
#---
e <- new.env(); assign("a", 1, e); get0("a", envir = e, mode = "character", ifnotfound = "nomode")
#---
e <- new.env(); assign("a", 1, e); get0("a", envir = e, mode = "numeric")
#---
get0(c("a", "b"))
#---
e <- new.env(); assign("a", 1, e); exists("a", where = e); exists("a", e, inherits = FALSE); exists("sum", e, inherits = FALSE); exists("sum", e)
#---
e <- new.env(); assign("a", 1, e); exists("a", envir = e, mode = "numeric"); exists("a", envir = e, mode = "character"); exists("a", envir = e, mode = "any")
#---
exists(c("a", "b"))
#---
exists(1)
#---
exists("a", frame = 1)
#---
exists("a", where = 1)
#---
exists("sum", mode = "function"); exists("pi", mode = "function"); exists("pi", mode = "numeric")
#---
e <- new.env(); assign("a", 1, e); rm("a", envir = e); exists("a", envir = e, inherits = FALSE)
#---
e <- new.env(); assign("a", 1, e); assign("b", 2, e); rm(list = c("a", "b"), envir = e); ls(e)
#---
e1 <- new.env(); e2 <- new.env(parent = e1); assign("v", 1, e1); rm("v", envir = e2); exists("v", envir = e1, inherits = FALSE)
#---
e1 <- new.env(); e2 <- new.env(parent = e1); assign("v", 1, e1); rm("v", envir = e2, inherits = TRUE); exists("v", envir = e1, inherits = FALSE)
#---
e1 <- new.env(); e2 <- new.env(parent = e1); rm("v", envir = e2, inherits = TRUE)
#---
a <- 1; b <- 2; rm(a, b); exists("a"); exists("b")
#---
x1 <- 1; x2 <- 2; rm(list = ls(pattern = "^x[0-9]$")); ls()
#---
rm("never_defined")
#---
f <- function() { v <- 1; rm(v); exists("v", inherits = FALSE) }; f()
#---
rm(1)
#---
e <- new.env(); assign("a", 1, e); rm(a, envir = e); ls(e)
#---
assign("sum", 1, envir = baseenv())
#---
rm("sum", envir = baseenv())
#---
environmentIsLocked(baseenv()); environmentIsLocked(globalenv()); environmentIsLocked(emptyenv())
#---
bindingIsLocked("sum", baseenv()); bindingIsLocked("pi", baseenv()); bindingIsLocked("T", baseenv())
#---
a <- 1; lockBinding("a", environment()); rm(a); exists("a")
#---
assign("x", 1, envir = emptyenv())
#---
T <- 0; T; rm(T); T
#---
assign("brandnew", 1, envir = baseenv())
#---
asNamespace("base")
#---
getNamespace("base")
#---
.BaseNamespaceEnv
#---
environmentName(asNamespace("base")); environmentName(.BaseNamespaceEnv); environmentName(baseenv()); environmentName(globalenv()); environmentName(emptyenv())
#---
identical(asNamespace("base"), baseenv()); identical(asNamespace("base"), .BaseNamespaceEnv); identical(getNamespace("base"), asNamespace("base")); identical(asNamespace(asNamespace("base")), asNamespace("base"))
#---
isNamespace(baseenv()); isNamespace(asNamespace("base")); isNamespace(globalenv()); isNamespace(1); isNamespace("base"); isNamespace(NULL)
#---
isBaseNamespace(asNamespace("base")); isBaseNamespace(baseenv()); isBaseNamespace(1)
#---
getNamespaceName(asNamespace("base")); getNamespaceName("base")
#---
getNamespaceName(baseenv())
#---
asNamespace(baseenv())
#---
asNamespace(1)
#---
asNamespace(NULL)
#---
asNamespace("foo")
#---
getNamespace("foo")
#---
getNamespace(1)
#---
asNamespace(c("base", "base"))
#---
parent.env(asNamespace("base"))
#---
identical(parent.env(asNamespace("base")), globalenv())
#---
parent.env(baseenv())
#---
topenv(asNamespace("base")); topenv(baseenv()); topenv(emptyenv())
#---
topenv(); topenv(new.env()); topenv(new.env(parent = baseenv())); topenv(new.env(parent = emptyenv()))
#---
f <- function() topenv(); f()
#---
topenv(new.env(), matchThisEnv = NULL)
#---
e <- new.env(); identical(topenv(new.env(parent = e), matchThisEnv = e), e)
#---
isNamespaceLoaded("base"); isNamespaceLoaded("foo")
#---
requireNamespace("base"); requireNamespace("base", quietly = TRUE)
#---
requireNamespace("foo", quietly = TRUE)
#---
withVisible(requireNamespace("base"))$visible
#---
is.character(getNamespaceExports("base")); "sum" %in% getNamespaceExports("base")
#---
exists("sum", envir = asNamespace("base"), inherits = FALSE); is.function(get("sum", envir = asNamespace("base")))
#---
identical(asNamespace("base")$sum, sum); identical(asNamespace("base")[["sum"]], sum); asNamespace("base")$nonexistent
#---
class(asNamespace("base")); is.environment(asNamespace("base"))
#---
environment(sum)
#---
base::sum(1:3); base::`+`(1, 2); base:::paste("a", "b"); base::pi
#---
(base::paste)("a", "b")
#---
"base"::sum(1)
#---
base::"sum"(1)
#---
base::letters[1:3]; base::T
#---
base::nope
#---
base:::nope
#---
x <- "base"; base::x
#---
foo::bar
#---
foo:::bar
#---
base::get0("nada", ifnotfound = 3)
#---
f <- function() base::sum(1, 2); f()
#---
sapply(1:3, base::sqrt)
#---
do.call(base::paste, list("a", "b", sep = "-"))
#---
identical(base::sum, sum); identical(base::paste, paste)
#---
is.function(base::sum)
#---
.GlobalEnv
#---
identical(.GlobalEnv, globalenv())
#---
as.environment("package:base")
#---
identical(as.environment("package:base"), baseenv())
#---
as.environment(".GlobalEnv"); as.environment(1)
#---
as.environment("foo")
#---
as.environment(NULL)
#---
identical(as.environment(asNamespace("base")), asNamespace("base"))
#---
baseenv(); emptyenv(); globalenv()
#---
e <- new.env(); attr(e, "name") <- "myenv"; environmentName(e)
#---
e <- new.env(); attr(e, "name") <- "myenv"; attributes(e)
#---
e <- new.env(); class(e) <- "foo"; class(e); inherits(e, "foo"); is.environment(e)
#---
e <- new.env(); class(e) <- "foo"; print.foo <- function(x, ...) cat("<a foo>\n"); print(e); TRUE
#---
e <- new.env(); class(e) <- c("child", "parent"); format.parent <- function(x, ...) "p"; hello <- function(o) UseMethod("hello"); hello.parent <- function(o) "hello parent"; hello(e)
#---
e <- new.env(); attr(e, "k") <- 1; f <- function(x) { attr(x, "k") <- 2; x }; g <- f(e); attr(e, "k")
#---
e <- new.env(); class(e) <- "foo"; e$x <- 1; e$x; unclass(e)$x
#---
environmentName(5); environmentName(NULL); environmentName("base"); environmentName(new.env())
#---
f <- function() 1; environment(f) <- asNamespace("base"); environmentName(environment(f))
#---
e <- new.env(); parent.env(e) <- asNamespace("base"); environmentName(parent.env(e)); identical(parent.env(e), asNamespace("base"))
#---
e <- new.env(); parent.env(e) <- emptyenv(); environmentName(parent.env(e)); exists("sum", e)
#---
local({ a <- 1; b <- 2; a + b })
#---
x <- 10; local({ x <- 20; x }); x
#---
e <- new.env(); local({ y <- 5 }, envir = e); get("y", e)
#---
f <- function() { local({ z <- 1; environment() }) }; class(f())
#---
local({ v <- 1; g <- function() v <<- v + 1; g(); g(); v })
#---
e <- new.env(); local(a <- 1, e); local(a <- a + 1, e); get("a", e)
#---
local()
#---
evalq(environment(), new.env()) -> ee; class(ee)
#---
e <- new.env(); e$one <- 1; evalq(one + 1, e); eval(quote(one + 2), e); eval(quote(one), list(one = 9), e)
#---
f <- function() { x <- 1; g <- function() parent.frame(); identical(g(), environment()) }; f()
#---
f <- function() sys.function(); identical(f(), f)
#---
e <- new.env(); assign("a", 1, e); environmentName(parent.env(e)); identical(parent.env(e), globalenv())
#---
f <- function() { e <- environment(); lockEnvironment(e); tryCatch(assign("b", 1, e), error = function(err) conditionMessage(err)) }; f()
#---
e <- new.env(); lockBinding("zzz", e)
#---
res <- tryCatch({ e <- new.env(); assign("a", 1, e); lockBinding("a", e); e$a <- 2 }, error = function(err) conditionMessage(err)); res
#---
e <- new.env(); assign("a", 1, e); lockBinding("a", e); withCallingHandlers(tryCatch(assign("a", 9, e), error = function(err) class(err)), warning = function(w) 1)
#---
e <- new.env(); lockEnvironment(e); tryCatch(e$new <- 1, error = function(err) conditionMessage(err)); tryCatch(rm("new", envir = e), error = function(err) conditionMessage(err))
#---
f <- function() { delayedAssign("lazy", stop("never forced")); "fine" }; f()
#---
e <- new.env(); assign("x", 1, e); lockBinding("x", e); isTRUE(tryCatch({ e$x <- 2; FALSE }, error = function(err) TRUE))
#---
counter <- local({ n <- 0; function() { n <<- n + 1; n } }); lockBinding("counter", environment()); counter(); counter()
#---
cfg <- new.env(); cfg$debug <- FALSE; lockEnvironment(cfg, bindings = TRUE); tryCatch(cfg$debug <- TRUE, error = function(e) cat("locked:", conditionMessage(e), "\n")); cfg$debug
#---
e <- new.env(); makeActiveBinding("now", function() "tick", e); with(e, now)
#---
e <- list2env(list(a = 1, b = 2)); lockEnvironment(e, TRUE); sort(ls(e)); mget(c("a", "b"), envir = e)
#---
as.list.factor(factor(c("p", "q", "p")))
#---
as.list.function(function(a, b = 3) a * b)
#---
with.default(list(a = 4), a * a)
