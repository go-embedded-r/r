`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
f <- y ~ x + z; f
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
f <- y ~ x + z; class(f)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
f <- y ~ x + z; attributes(f)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
f <- y ~ x + z; length(f)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
f <- y ~ x + z; f[[2]]
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
f <- y ~ x + z; f[[3]]
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
f <- y ~ x + z; deparse(f)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
~ a
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
g <- ~ a | b; length(g)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
f <- y ~ x; environment(f)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
f <- y ~ x; identical(environment(f), globalenv())
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
f <- y ~ x + z; as.character(f)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
as.character(~ x)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
is.call(y ~ x)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
is.function(y ~ x)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
inherits(y ~ x, "formula")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
identical(y ~ x, y ~ x)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
f <- y ~ x; identical(f, f)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
format(y ~ x)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
format(y ~ x + z, width = 5)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
class(format(y ~ x))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
x <- y ~ x + z; x[-3]
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
x <- y ~ x + z; class(x[-3])
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
x <- y ~ x + z; class(x[2])
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
x <- y ~ x + z; x[[3]] <- quote(a * b); x
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
x <- y ~ x + z; x[[3]] <- quote(a * b); class(x)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
x <- y ~ x + z; x[[2]] <- NULL; x
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
x <- y ~ x + z; attributes(x[-2])
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
x <- y ~ x + z; x[[3]][[3]] <- as.name("w"); x
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
rev(as.list(y ~ x))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
as.list(y ~ x)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
deparse(y ~ x | z)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
y ~ x | z
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
log(y) ~ sqrt(x) + I(z^2)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
y ~ .
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
~ 1
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
y ~ 0 + x
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
y ~ x %in% z
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
y ~ a * b + (c + d)^2 - 1
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
cat(deparse(y ~ x * z), "\n")
#---
f <- function() z ~ w; class(f())
#---
f <- function() { v <- 1; z ~ w }; is.environment(environment(f()))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
f <- y ~ x; environment(f) <- baseenv(); environment(f)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
f <- y ~ x; environment(f) <- baseenv(); f
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tryCatch(`environment<-`(y ~ x, 1), error = function(e) conditionMessage(e))
#---
q <- quote(y ~ x); class(q)
#---
q <- quote(y ~ x); class(eval(q))
#---
q <- quote(y ~ x); e <- eval(q); e2 <- eval(e); identical(e, e2)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
f <- y ~ x; g <- f; g[[2]] <- as.name("w"); f
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
sapply(list(y ~ x, ~ z), length)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
l <- list(y ~ x, ~ z); l
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
str2lang("y ~ x")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
class(str2lang("y ~ x"))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
all.vars(y ~ x + log(z))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
all.names(y ~ x + log(z))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
all.names(y ~ x + log(z), functions = FALSE)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
all.vars(quote(f(x, y) + x), unique = FALSE)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
all.names(quote(a + b + a), unique = TRUE)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
all.names(quote(a + b + c), max.names = 2)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
all.vars(expression(x + y, z))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
all.vars(quote(x$y))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
all.vars(1)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
all.vars(NULL)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
all.vars(quote(function(a) a + b))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
all.names(quote(if (a) b else c))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
all.vars("a")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
all.vars(list(1))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
all.vars(y ~ x + z, functions = TRUE)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
all.names(quote(f(g(h(x)))))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
all.vars(quote(a[b][[c]]))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
F <- function(x) { environment(x) <- globalenv(); x }
F(as.formula("y ~ x"))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
class(as.formula("y ~ x"))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
F <- function(x) { environment(x) <- globalenv(); x }
F(as.formula(quote(y ~ x)))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
as.formula(y ~ x)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
f <- y ~ x; identical(as.formula(f), f)
#---
e <- new.env(); f <- as.formula("y ~ x", env = e); identical(environment(f), e)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
F <- function(x) { environment(x) <- globalenv(); x }
F(formula("y ~ x + z"))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
formula(terms(y ~ a + b))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
attributes(formula(terms(y ~ a + b)))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
F <- function(x) { environment(x) <- globalenv(); x }
F(formula(data.frame(y = 1, a = 2, b = 3)))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
F <- function(x) { environment(x) <- globalenv(); x }
F(formula(data.frame(`a b` = 1, c = 2, check.names = FALSE)))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
formula(list(formula = y ~ x))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
formula(list(terms = terms(y ~ x)))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tryCatch(formula(structure(list(call = quote(lm(y ~ x))), class = "foo")), error = function(e) conditionMessage(e))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tryCatch(formula(1), error = function(e) conditionMessage(e))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
F <- function(x) { environment(x) <- globalenv(); x }
F(formula(NULL))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tryCatch(as.formula("1 +"), error = function(e) "parse error")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
formula(structure(list(), formula = y ~ x))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
F <- function(x) { environment(x) <- globalenv(); x }
F(reformulate(c("a", "b+c")))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
F <- function(x) { environment(x) <- globalenv(); x }
F(reformulate(c("a", "log(b)"), response = "y"))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
F <- function(x) { environment(x) <- globalenv(); x }
F(reformulate("a", response = quote(log(y))))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
F <- function(x) { environment(x) <- globalenv(); x }
F(reformulate("a", response = "log(y)"))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
F <- function(x) { environment(x) <- globalenv(); x }
F(reformulate("a", intercept = FALSE))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
F <- function(x) { environment(x) <- globalenv(); x }
F(reformulate(c("a", "b"), response = as.name("y"), intercept = FALSE))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
F <- function(x) { environment(x) <- globalenv(); x }
F(reformulate("a:b", "y"))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tryCatch(reformulate(1), error = function(e) conditionMessage(e))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tryCatch(reformulate("a", c("y", "z")), error = function(e) conditionMessage(e))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
F <- function(x) { environment(x) <- globalenv(); x }
attributes(F(reformulate("a")))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
update(y ~ x, . ~ . + z)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
update(y ~ x, log(.) ~ .)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
update(y ~ x, ~ . + z)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
update(y ~ x, z ~ .)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
update(~ x, y ~ .)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
update(~ x, . ~ . + z)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
update(y ~ x * z, . ~ . - x:z)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
update(y ~ x + z, . ~ . * w)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
update(y ~ x, . ~ . - 1)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
update(y ~ x - 1, . ~ . + 1)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
update(y ~ x, . ~ .)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
update(y ~ x, ~ .^2)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
update(y ~ a + b, . ~ .^2)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
update(y ~ x, "z ~ .")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
update.formula(y ~ x, . ~ . + w)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
class(update(y ~ x, . ~ . + w))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
update(y ~ a + b + c, . ~ . - b)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
update(y ~ a, sqrt(.) ~ . + I(a^2))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
update(y ~ a * b, . ~ . - a)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
update(cbind(y1, y2) ~ a, . ~ . + b)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
update(y ~ a + offset(o), . ~ . + b)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
f <- y ~ x; g <- update(f, . ~ . + w); f
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
x <- y ~ x; print(x)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
x <- y ~ x; r <- print(x); identical(r, x)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
x <- y ~ x; print(x, showEnv = TRUE)
#---
f <- function() z ~ w; x <- f(); capture.output(print(x))[1]
#---
f <- function() z ~ w; x <- f(); length(capture.output(print(x)))
#---
f <- function() z ~ w; x <- f(); capture.output(print(x, showEnv = FALSE))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
format(~ a + b)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
paste(y ~ x)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
toString(y ~ x)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
nchar(deparse(y ~ x + z))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
unclass(y ~ x)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
oldClass(y ~ x)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
inherits(y ~ x, c("a", "formula"), which = TRUE)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
is.language(y ~ x)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
typeof(y ~ x)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
mode(y ~ x)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
storage.mode(y ~ x)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
class(unclass(y ~ x))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
f <- y ~ x; attr(f, "foo") <- "bar"; attributes(f)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
f <- y ~ x; attr(f, "foo") <- "bar"; f
#---
f <- y ~ x; class(f) <- c("myform", "formula"); class(f)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
f <- y ~ x; class(f) <- c("myform", "formula"); inherits(f, "formula")
#---
f <- y ~ x; class(f) <- NULL; class(f)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
lapply(list(y ~ x), class)
#---
is.null(environment(y ~ x))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
f <- y ~ x; saveRDS(f, tf <- tempfile()); g <- readRDS(tf); identical(f, g)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
f <- y ~ x; saveRDS(f, tf <- tempfile()); g <- readRDS(tf); class(g)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
f <- y ~ x; dput(f)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
f <- y ~ x + z; f[[3]][[2]]
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
f <- y ~ x + z; length(f[[3]])
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
f <- y ~ x + z; f[[1]]
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
f <- y ~ x + z; as.character(f[[1]])
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
f <- y ~ x + z; f[2:3]
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
f <- y ~ x + z; class(f[2:3])
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
f <- y ~ x + z; f[c(1, 3)]
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
f <- y ~ x + z; f[0]
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
f <- y ~ x + z; sapply(as.list(f), class)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
f <- y ~ x + z; do.call("~", list(quote(a), quote(b)))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
f <- y ~ x + z; class(do.call("~", list(quote(a), quote(b))))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
f <- y ~ x + z; class(as.call(list(as.name("~"), quote(a))))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
f <- y ~ x + z; class(call("~", quote(a)))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
f <- y ~ x + z; class(eval(call("~", quote(a))))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
f <- y ~ x + z; Reduce(function(a, b) call("+", a, b), list(quote(a), quote(b), quote(c)))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
f <- y ~ x + z; substitute(a ~ b, list(a = quote(u), b = quote(v)))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
f <- y ~ x + z; class(substitute(a ~ b, list(a = quote(u), b = quote(v))))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
f <- y ~ x + z; bquote(.(f[[2]]) ~ w)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
f <- y ~ x + z; class(bquote(.(f[[2]]) ~ w))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
f <- y ~ x + z; eval(bquote(.(f[[2]]) ~ w))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
f <- y ~ x + z; match.call(sum, quote(sum(1, 2)))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
`~`(a, b)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
`~`(a)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
sapply(c(y ~ x, z ~ w), class)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
with(list(a = 1), class(y ~ a))
#---
terms.formula(y ~ a * b)$term.labels
#---
attr(terms.formula(y ~ a * b), "term.labels")
#---
.MFclass(1:3)
#---
sapply(list(TRUE, "a", factor("a"), factor("a", ordered = TRUE), 1.5, matrix(1:4, 2), list(1)), .MFclass)
#---
print.default(1:3)
#---
print.default("a", quote = FALSE)
#---
registerS3method("print", "zzfoo", function(x, ...) cat("zzfoo\n")); x <- structure(1, class = "zzfoo"); print(x); invisible(NULL)
#---
registerS3method("format", "zzbar", function(x, ...) "bar!"); format(structure(1, class = "zzbar"))
#---
r <- registerS3method("print", "zzfoo", function(x, ...) cat("zzfoo\n")); r
#---
d <- data.frame(y = 1:3, u = c(1, 2, 3)); mf <- model.frame.default(y ~ u, data = d); dim(mf)
#---
d <- data.frame(y = 1:3, u = c(1, 2, 3)); names(model.frame.default(y ~ u, data = d, subset = y > 1))
