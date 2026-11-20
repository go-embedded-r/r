lt <- function(v) as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + v); x <- lt(c(0, 10, 20)); class(x); mean(x)
#---
lt <- function(v) as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + v); x <- lt(c(0, NA, 20)); mean(x)
#---
lt <- function(v) as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + v); x <- lt(c(0, NA, 20)); mean(x, na.rm = TRUE)
#---
lt <- function(v) as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + v); x <- lt(c(30, 10, 20)); sort(x)
#---
lt <- function(v) as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + v); x <- lt(c(30, 10, 20)); sort(x, decreasing = TRUE)
#---
lt <- function(v) as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + v); x <- lt(c(30, NA, 20)); sort(x, na.last = TRUE)
#---
lt <- function(v) as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + v); x <- lt(c(30, NA, 20)); sort(x)
#---
lt <- function(v) as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + v); x <- lt(c(10, 10, 20)); unique(x)
#---
lt <- function(v) as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + v); x <- lt(c(10, 10, 20)); duplicated(x)
#---
lt <- function(v) as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + v); x <- lt(c(10, 10, 20)); anyDuplicated(x)
#---
lt <- function(v) as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + v); x <- lt(0); rep(x, 3)
#---
lt <- function(v) as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + v); x <- lt(c(0, 10)); rep(x, each = 2)
#---
lt <- function(v) as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + v); x <- lt(c(0, 10)); rep(x, times = c(2, 1))
#---
lt <- function(v) as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + v); x <- lt(c(0, 10)); rep(x, length.out = 3)
#---
lt <- function(v) as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + v); x <- lt(c(0, 10)); class(rep(x, 2))
#---
lt <- function(v) as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + v); x <- lt(c(0, NA)); anyNA(x)
#---
lt <- function(v) as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + v); x <- lt(c(0, 10)); anyNA(x)
#---
lt <- function(v) as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + v); x <- lt(c(0, NA)); is.finite(x)
#---
lt <- function(v) as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + v); x <- lt(c(0, NA)); is.infinite(x)
#---
lt <- function(v) as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + v); x <- lt(c(0, NA)); is.nan(x)
#---
lt <- function(v) as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + v); x <- lt(c(0, 10, 20)); length(x) <- 2; x
#---
lt <- function(v) as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + v); x <- lt(c(0, 10)); length(x) <- 4; x
#---
lt <- function(v) as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + v); x <- lt(c(0, 10)); xtfrm(x)
#---
lt <- function(v) as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + v); x <- lt(c(0, 10)); as.double(x)
#---
lt <- function(v) as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + v); x <- lt(c(0, 10)); as.vector(x)
#---
lt <- function(v) as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + v); x <- lt(c(0, 10)); as.vector(x, "list")
#---
lt <- function(v) as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + v); x <- lt(c(0, 10)); as.vector(x, "character")
#---
lt <- function(v) as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + v); x <- lt(c(0, 10)); length(as.list(x))
#---
lt <- function(v) as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + v); x <- lt(c(0, 10)); class(as.list(x)[[2]])
#---
lt <- function(v) as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + v); x <- lt(c(0, 10)); as.matrix(x)
#---
lt <- function(v) as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + v); x <- lt(c(0, 10)); range(x)
#---
lt <- function(v) as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + v); x <- lt(c(0, 10)); min(x)
#---
lt <- function(v) as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + v); x <- lt(c(0, 10)); max(x)
#---
lt <- function(v) as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + v); x <- lt(c(0, 10)); sum(x)
#---
lt <- function(v) as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + v); x <- lt(c(0, 10)); Summary.POSIXlt(x)
#---
lt <- function(v) as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + v); x <- lt(0); julian(x)
#---
lt <- function(v) as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + v); x <- lt(0); julian(x, origin = as.POSIXct("2024-01-01", tz = "UTC"))
#---
lt <- function(v) as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + v); x <- lt(0); julian(x, origin = as.POSIXct(c("2024-01-01", "2024-01-02"), tz = "UTC"))
#---
lt <- function(v) as.POSIXlt(as.POSIXct("2024-01-02 03:04:05", tz = "UTC") + v); x <- lt(c(0, 10)); mtfrm(x)
#---
f <- function(x, y = 2) { x + y }; environment(f) <- globalenv(); print.function(f)
#---
f <- function(x, y = 2) { x + y }; environment(f) <- globalenv(); print(f, useSource = FALSE)
#---
f <- function(a) a; environment(f) <- globalenv(); r <- print.function(f); identical(r, f)
#---
tc <- textConnection("out", "w"); write.dcf(list(Package = "foo", Version = "1.0"), tc); close(tc); out
#---
write.dcf(data.frame(A = c("x", "y"), B = c("1", NA)))
#---
write.dcf(list(Package = "foo", Title = "A long description line that needs to be wrapped because it exceeds the width of the output"), width = 40)
#---
write.dcf(list(A = "multi\nline"))
#---
write.dcf(list(A = "x", B = "y"), keep.white = "B")
#---
write.dcf(list(A = "x", B = "y"), fields = "B")
#---
write.dcf(rbind(c(Name = "v", Ver = "1"), c(Name = "w", Ver = "2")))
#---
names(sys.status())
#---
f <- function() length(sys.status()$sys.calls) - length(sys.calls()); f()
#---
f <- function() sys.status(); g <- function() f(); length(g()$sys.parents) - length(f()$sys.parents)
#---
f <- function() sys.parents(); g <- function() f(); x <- g(); x[length(x)] - length(x) + 1L
#---
truncate(1)
#---
chooseOpsMethod(1, 2, `+`, `+`, NULL, FALSE)
#---
chooseOpsMethod(1, 2, `+`, `+`, NULL, TRUE)
#---
chooseOpsMethod.default(1, 2, `+`, `+`, NULL, FALSE)
#---
nameOfClass(1)
#---
nameOfClass.default("a")
#---
nameOfClass.foo <- function(x) "bar"; nameOfClass(structure(1, class = "foo"))
#---
srcfile("x.R")$filename
#---
sf <- srcfile("x.R"); class(sf)
#---
sf <- srcfilecopy("f.R", c("a <- 1", "b <- 2")); class(sf); sf$lines
#---
sf <- srcfilecopy("f.R", c("a <- 1", "b <- 2")); getSrcLines(sf, 1, 2)
#---
sf <- srcfilecopy("f.R", c("a <- 1", "b <- 2")); getSrcLines(sf, 2, 5)
#---
sf <- srcfilecopy("f.R", c("a <- 1", "b <- 2")); getSrcLines(sf, 3, 2)
#---
sf <- srcfilecopy("f.R", c("a <- 1", "b <- 2")); sr <- srcref(sf, c(1, 1, 1, 6)); as.character(sr)
#---
sf <- srcfilecopy("f.R", c("a <- 1", "b <- 2")); sr <- srcref(sf, c(1, 1, 2, 6)); print(sr); NULL
#---
sf <- srcfilecopy("f.R", c("a <- 1", "b <- 2")); sr <- srcref(sf, c(2, 1, 2, 6)); class(summary(sr)); NULL
#---
sf <- srcfilecopy("f.R", c("a <- 1", "b <- 2")); sr <- srcref(sf, c(2, 1, 2, 6)); class(sr); length(unclass(sr)); unclass(sr)[1:6]
#---
sf <- srcfilecopy("f.R", c("a <- 1", "b <- 2")); print(sf); NULL
#---
sf <- srcfilecopy("f.R", c("a <- 1", "b <- 2")); sf$filename; sf$isFile; sf$Enc
#---
srcref(1, 2)
#---
sf <- srcfilecopy("f.R", "a <- 1"); sa <- srcfilealias("g.R", sf); class(sa); sa$filename
#---
sf <- srcfilecopy("f.R", c("a <- 1", "b <- 2")); sa <- srcfilealias("g.R", sf); getSrcLines(sa, 1, 1)
