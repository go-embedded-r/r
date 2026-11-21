x <- structure(list(a = 1, b = 2, c = 3), class = "listof"); y <- `[.listof`(x, 1:2); class(y); length(y)
#---
x <- structure(list(a = 1, b = 2, c = 3), class = "listof"); print.listof(x)
#---
x <- noquote(c("a", "b")); y <- `[.noquote`(x, 2); class(y); unclass(y)
#---
x <- noquote(c("a", "b")); y <- c.noquote(x, "c"); class(y); unclass(y)
#---
x <- noquote(matrix(c("a", "b", "c", "d"), 2)); m <- as.matrix.noquote(x); class(m); dim(m)
#---
x <- structure(list(1, 2), class = "simple.list"); y <- `[.simple.list`(x, 1); class(y); length(y)
#---
t1 <- table(c("a", "b", "a")); as.vector(t1[1])
#---
f <- factor(c("a", "b")); `[[.factor`(f, 2)
#---
f <- factor(c("TRUE", "FALSE")); as.logical.factor(f)
#---
f <- factor(c("a", "b")); as.vector.factor(f)
#---
as.vector.data.frame(data.frame(a = 1:2))
#---
anyDuplicated.data.frame(data.frame(a = c(1, 1)))
#---
anyNA.data.frame(data.frame(a = c(1, NA)))
#---
all.equal.formula(y ~ x, y ~ x)
#---
all.equal.formula(y ~ x, y ~ z)
#---
all.equal.POSIXt(as.POSIXct("2024-01-01", tz = "UTC"), as.POSIXct("2024-01-01", tz = "UTC"))
#---
f <- approxfun(1:3, c(2, 4, 8)); f(c(1.5, 2.5, 4))
#---
f <- approxfun(1:3, c(2, 4, 8), method = "constant"); f(c(1.5, 2.5))
#---
is.unsorted(c(1, 3, 2)); is.unsorted(1:3); is.unsorted(c(1, 1), strictly = TRUE)
#---
lt <- as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + c(0, 10, 10)); as.double.POSIXlt(lt)
#---
lt <- as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + c(0, 10, 10)); length(as.list.POSIXlt(lt))
#---
lt <- as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + c(0, 10, 10)); dim(as.matrix.POSIXlt(lt))
#---
lt <- as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + c(0, 10, 10)); length(as.vector.POSIXlt(lt, "list"))
#---
lt <- as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + c(0, 10, 10)); duplicated.POSIXlt(lt)
#---
lt <- as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + c(0, 10, 10)); unique.POSIXlt(lt)
#---
lt <- as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + c(0, 10, 10)); anyNA.POSIXlt(lt)
#---
lt <- as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + c(0, 10, 10)); is.finite.POSIXlt(lt); is.infinite.POSIXlt(lt); is.nan.POSIXlt(lt)
#---
lt <- as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + c(0, 10, 10)); julian.POSIXt(lt)
#---
lt <- as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + c(0, 10, 10)); length(`length<-.POSIXlt`(lt, 2))
#---
lt <- as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + c(0, 10, 10)); mean.POSIXlt(lt)
#---
lt <- as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + c(0, 10, 10)); length(mtfrm.POSIXlt(lt))
#---
lt <- as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + c(0, 10, 10)); length(rep.POSIXlt(lt, 2))
#---
lt <- as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + c(10, 0, 5)); sort.POSIXlt(lt)
#---
lt <- as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + c(10, 0, 5)); xtfrm.POSIXlt(lt)
#---
x <- ftable(table(a = c(1, 1, 2), b = c("x", "y", "x"))); as.data.frame(x)
#---
x <- ftable(table(a = c(1, 1, 2), b = c("x", "y", "x"))); as.matrix(x)
#---
x <- ftable(table(a = c(1, 1, 2), b = c("x", "y", "x"))); as.table(x)
#---
x <- ftable(table(a = c(1, 1, 2), b = c("x", "y", "x"))); format(x)
#---
x <- ftable(table(a = c(1, 1, 2), b = c("x", "y", "x"))); print(x)
#---
Fn <- ecdf(c(3, 1, 2, 2)); summary(Fn)
#---
Fn <- ecdf(c(3, 1, 2, 2)); quantile(Fn, 0.5)
#---
Fn <- ecdf(c(3, 1, 2, 2)); knots(Fn)
#---
x <- polym(1:5, 1:5 + c(0, 1, 0, 1, 0), degree = 1); dim(x)
#---
x <- poly(1:5, 2); round(predict(x, 6:7)[, 1], 6) + 0 # ulp
#---
regexec("b+", c("abbc", "xyz"))
#---
sf <- srcfilecopy("f.R", c("a <- 1", "b <- 2")); sr <- srcref(sf, c(1, 1, 1, 6)); as.character.srcref(sr)
#---
sf <- srcfilecopy("f.R", c("a <- 1", "b <- 2")); sr <- srcref(sf, c(1, 1, 1, 6)); print.srcref(sr); NULL
#---
sf <- srcfilecopy("f.R", c("a <- 1", "b <- 2")); sr <- srcref(sf, c(1, 1, 1, 6)); class(summary.srcref(sr)); NULL
#---
sf <- srcfilecopy("f.R", c("a <- 1", "b <- 2")); print.srcfile(sf); NULL
#---
sf <- srcfilecopy("f.R", c("a <- 1", "b <- 2")); sa <- srcfilealias("g.R", sf); open.srcfilealias(sa); close.srcfilealias(sa); NULL
#---
x <- ftable(table(a = c(1, 1, 2), b = c("x", "y", "x"))); stats:::as.data.frame.ftable(x)
#---
x <- ftable(table(a = c(1, 1, 2), b = c("x", "y", "x"))); stats:::as.matrix.ftable(x)
#---
x <- ftable(table(a = c(1, 1, 2), b = c("x", "y", "x"))); stats:::as.table.ftable(x)
#---
x <- ftable(table(a = c(1, 1, 2), b = c("x", "y", "x"))); stats:::format.ftable(x)
#---
x <- ftable(table(a = c(1, 1, 2), b = c("x", "y", "x"))); stats:::print.ftable(x)
#---
stats:::ftable.default(table(a = c(1, 1, 2), b = c("x", "y", "x")))
#---
stats:::ftable.formula(a ~ b, data.frame(a = c(1, 1, 2), b = c("x", "y", "x")))
#---
Fn <- ecdf(c(3, 1, 2, 2)); stats:::summary.ecdf(Fn)
#---
Fn <- ecdf(c(3, 1, 2, 2)); stats:::print.summary.ecdf(summary(Fn))
#---
Fn <- ecdf(c(3, 1, 2, 2)); stats:::quantile.ecdf(Fn, 0.5)
#---
Fn <- ecdf(c(3, 1, 2, 2)); stats:::knots.stepfun(Fn)
#---
Fn <- ecdf(c(3, 1, 2, 2)); x <- capture.output(stats:::print.ecdf(Fn)); x[-1]
#---
x <- poly(1:5, 2); round(stats:::predict.poly(x, 6:7)[, 1], 6) + 0 # ulp
#---
t1 <- table(c("a", "b", "a")); tryCatch(`[.table`(t1, 1), error = function(e) conditionMessage(e)); t1[2]
#---
sf <- srcfilecopy("f.R", c("a <- 1", "b <- 2")); r <- withVisible(close.srcfile(sf)); r$visible; r$value
#---
ts <- as.POSIXct("2024-01-01", tz = "UTC"); sf <- srcfilecopy("f.R", c("a <- 1", "b <- 2"), timestamp = ts); tryCatch(open.srcfile(sf, 1), error = function(e) "failed")
#---
sf <- srcfilecopy("f.R", c("a <- 1", "b <- 2")); tryCatch(open.srcfilecopy(sf), error = function(e) conditionMessage(e))
#---
sf <- srcfilecopy("f.R", c("a <- 1", "b <- 2")); cn <- open.srcfilecopy(sf, 2); class(cn); readLines(cn, 1); close.srcfile(sf)
#---
ts <- as.POSIXct("2024-01-01", tz = "UTC"); sf <- srcfilecopy("f.R", c("a <- 1", "b <- 2"), timestamp = ts); summary.srcfile(sf); NULL
#---
t1 <- table(c("a", "b", "a", "c")); t1[2]; t1[1:2]; t1[2, drop = FALSE]; t1[0]; t1["a"]
#---
t2 <- table(x = c("a", "b", "a"), y = c("u", "v", "v")); t2[1, ]; t2[1, , drop = FALSE]; t2[1:2, 1:2]
