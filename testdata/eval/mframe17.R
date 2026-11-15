`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
get_all_vars(y ~ u, d)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
get_all_vars(y ~ u, d, w = v)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
tryCatch(get_all_vars(y ~ u + zzz, d), error = function(e) conditionMessage(e))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
get_all_vars(y ~ u + v, d)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
x <- 1:3; w <- c(2, 4, 6); get_all_vars(w ~ x)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
get_all_vars(d)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
get_all_vars(~ u, d)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
model.frame(y ~ u + v, d)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
model.frame(y ~ u + v, d, na.action = na.omit)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
attributes(model.frame(y ~ u + v, d, na.action = na.omit))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
attr(model.frame(y ~ u + v, d, na.action = na.exclude), "na.action")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
model.frame(y ~ u + v, d, subset = y > 1, na.action = na.pass)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
model.frame(y ~ u + v, d, subset = c(TRUE, FALSE, TRUE, TRUE), na.action = NULL)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
model.frame(y ~ u + v, d, subset = 2:3, na.action = na.pass)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
model.frame(~ u, d, na.action = na.pass)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
model.frame(abs(y) ~ sqrt(u) + v, d, na.action = na.pass)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
model.frame(y ~ ., d, na.action = na.pass)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
model.frame(y ~ u, d, weights = v, na.action = na.pass)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
attr(model.frame(y ~ u + I(u^2), d, na.action = na.pass), "terms")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
mf <- model.frame(y ~ v, d[1:2, ], drop.unused.levels = TRUE); levels(mf$v)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
mf <- model.frame(y ~ v, d[1:2, ]); levels(mf$v)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
model.response(model.frame(y ~ u, d, na.action = na.pass))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
model.response(model.frame(~ u, d, na.action = na.pass))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
model.response(model.frame(y ~ u, d, na.action = na.pass), "numeric")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
model.response(model.frame(cbind(y, u) ~ v, d, na.action = na.pass))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
mf <- model.frame(y ~ u, d, na.action = na.pass); list(names(mf), class(mf), dim(mf))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
x <- 1:3; y2 <- c(2, 4, 6); model.frame(y2 ~ x)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
x <- 1:3; y2 <- c(2, 4, 6); model.frame(y2 ~ x, na.action = NULL)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
model.frame(y ~ u, data = list(y = 1:2, u = 3:4))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
tryCatch(model.frame(y ~ qqq, d), error = function(e) conditionMessage(e))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tryCatch(model.frame(y ~ u, 1:3), error = function(e) "error")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
tryCatch(model.frame(y ~ u, as.matrix(d[, 1:2])), error = function(e) conditionMessage(e))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
na.omit(c(1, NA, 3))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
na.omit(data.frame(a = c(1, NA), b = c("x", "y")))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
na.omit(c(a = 1, b = NA, c = 3))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
na.omit(matrix(c(1, NA, 3, 4), 2))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
na.omit(1:3)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
na.omit(list(1, NA))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
na.fail(c(1, 2))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tryCatch(na.fail(c(1, NA)), error = function(e) conditionMessage(e))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tryCatch(na.fail(data.frame(a = c(1, NA))), error = function(e) conditionMessage(e))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
na.exclude(c(a = 1, b = NA))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
na.exclude(data.frame(a = c(1, NA, 3)))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
attr(na.omit(data.frame(a = c(1, NA, 3))), "na.action")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
na.pass(c(1, NA))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
complete.cases(data.frame(a = c(1, NA, 3), b = c("x", "y", NA)))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
complete.cases(c(1, NA), c(1, 2))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
complete.cases(matrix(c(1, NA, 3, 4), 2))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tryCatch(complete.cases(1:3, 1:2), error = function(e) conditionMessage(e))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tryCatch(complete.cases(), error = function(e) conditionMessage(e))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
complete.cases(list(a = c(1, NA), b = c(2, 3)))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
mf <- model.frame(y ~ u, d); nrow(mf)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
mf <- model.frame(y ~ u, d, na.action = "na.exclude"); attr(mf, "na.action")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
mf <- model.frame(y ~ u + v, d, na.action = na.omit, xlev = list(v = c("a", "b", "c", "d", "e"))); levels(mf$v)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
tryCatch(model.frame(y ~ v, d, xlev = list(v = c("a", "b"))), error = function(e) conditionMessage(e))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
mf <- model.frame(y ~ u, d, na.action = na.pass); attr(attr(mf, "terms"), "dataClasses")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
mf <- model.frame(y ~ u + v + I(u > 1), d, na.action = na.pass); attr(attr(mf, "terms"), "dataClasses")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:2, s = c("p", "q"), l = c(TRUE, FALSE)); attr(attr(model.frame(y ~ s + l, d), "terms"), "dataClasses")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:2, o = factor(c("p", "q"), ordered = TRUE)); attr(attr(model.frame(y ~ o, d), "terms"), "dataClasses")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:3); d$m <- matrix(1:6, 3); attr(attr(model.frame(y ~ m, d), "terms"), "dataClasses")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
mf <- model.frame(y ~ u + v, d, na.action = na.pass); model.response(mf, "numeric")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
mf <- model.frame(v ~ u, d, na.action = na.pass); model.response(mf, "numeric")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
tryCatch(model.response(d), error = function(e) conditionMessage(e))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
mf <- model.frame(y ~ u, d, na.action = na.pass); tryCatch(model.response(mf, "foo"), error = function(e) conditionMessage(e))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
mf <- model.frame(y ~ u, d, na.action = na.pass); model.frame(mf)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
mf <- model.frame(y ~ u, d, na.action = na.pass); identical(model.frame(mf), mf)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
mf <- model.frame(y ~ u, d, na.action = na.pass); model.frame(y ~ u, mf, na.action = na.omit)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
mf <- model.frame(y ~ u, d, na.action = na.omit); model.frame(terms(y ~ u), d, na.action = na.omit)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
fm <- y ~ u; mf <- model.frame(fm, d, na.action = na.omit); model.frame(fm, mf)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
f <- y ~ u; fn <- function(dat) model.frame(f, dat, na.action = na.pass); fn(d)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
w <- c(1, 1, 2, 2); model.frame(y ~ u, d, weights = w, na.action = na.pass)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
model.frame(y ~ u, d, weights = NULL, na.action = na.pass)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
names(model.frame(y ~ u + v, d, offset = u, na.action = na.pass))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
dim(model.frame(y ~ 1, d))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
names(model.frame(y ~ 1, d))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
dim(model.frame(~ 0, d))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
model.frame(y ~ u:v, d, na.action = na.pass)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
model.frame(y ~ abs(y) + u, d, na.action = na.pass)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
dim(model.frame(cbind(y, u) ~ v, d, na.action = na.pass))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
class(model.frame(cbind(y, u) ~ v, d, na.action = na.pass)[[1]])
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:4, u = c(1, 2, NA, 4), v = factor(letters[1:4]))
row.names(model.frame(y ~ u, d, subset = y > 2, na.action = na.pass))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:3, u = 4:6, row.names = c("a", "b", "c")); model.frame(y ~ u, d)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = 1:3, u = 4:6, row.names = c("a", "b", "c")); model.response(model.frame(y ~ u, d))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
y <- c(a = 1, b = 2); x <- 3:4; model.frame(y ~ x)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
y <- c(a = 1, b = 2); x <- 3:4; model.response(model.frame(y ~ x))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = c(1, NA, 3), u = 1:3); model.frame(y ~ u, d, na.action = na.exclude)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = c(1, NA, 3), u = 1:3); mf <- model.frame(y ~ u, d); attr(mf, "na.action")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = c(1, NA, 3), u = 1:3); mf <- model.frame(y ~ u, d); names(attributes(mf))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = c(1, NA, 3), u = 1:3); op <- options(na.action = "na.fail"); r <- tryCatch(model.frame(y ~ u, d), error = function(e) conditionMessage(e)); options(op); r
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = c(1, NA, 3), u = 1:3); op <- options(na.action = "na.pass"); r <- model.frame(y ~ u, d); options(op); r
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = c(1, NA, 3), u = 1:3); getOption("na.action")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = c(1, NA, 3), u = 1:3); model.frame(y ~ u, d, na.action = "na.pass")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
d <- data.frame(y = c(1, NA, 3), u = 1:3); model.frame(y ~ u, d, na.action = function(x) x[1, ])
