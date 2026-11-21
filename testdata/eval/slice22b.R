print(callCC(function(k) {k(3); 4}))
#---
print(callCC(function(k) 5))
#---
f <- function(x) callCC(function(k) { if (x < 0) k("neg"); "ok" })
print(c(f(-1), f(1)))
#---
saved <- NULL
callCC(function(k) { saved <<- k; 1 })
r <- try(saved(2), silent = TRUE)
cat(conditionMessage(attr(r, "condition")), "\n")
#---
print(callCC(function(k1) callCC(function(k2) { k1("outer"); "inner" })))
#---
p <- as.POSIXct("2020-01-01 00:00:00", tz = "UTC") + 0:2
p[2] <- as.POSIXct("2021-01-01", tz = "UTC"); print(p)
#---
p <- as.POSIXct("2020-01-01 00:00:00", tz = "UTC") + 0:2
p[3] <- "2022-02-02 10:00:00"; print(p)
#---
p <- as.POSIXct("2020-01-01 00:00:00", tz = "UTC") + 0:2
p[2] <- NA; print(p)
print(attributes(p))
#---
p <- as.POSIXct("2020-01-01 00:00:00", tz = "UTC") + 0:2
p[2] <- p[0]; print(p)
#---
p <- as.POSIXct("2020-01-01 00:00:00", tz = "UTC") + 0:2
p[5] <- p[1]; print(length(p)); print(class(p))
#---
dd <- as.difftime(c(1, 2, 3), units = "mins")
dd[2] <- as.difftime(1, units = "hours"); print(dd)
#---
dd <- as.difftime(c(1, 2, 3), units = "mins")
dd[3] <- 5; print(dd); print(attributes(dd))
#---
tt <- structure(1:4, tsp = c(2000, 2000.75, 4), class = "ts")
print(as.data.frame(tt))
#---
tt <- structure(1:4, tsp = c(2000, 2000.75, 4), class = "ts")
print(as.data.frame.ts(tt, row.names = c("a", "b", "c", "d")))
#---
tt <- structure(matrix(1:6, 3), tsp = c(2000, 2002, 1), class = c("mts", "ts"))
print(as.data.frame(tt))
#---
mm <- matrix(1:4, 2, dimnames = list(NULL, c("a", "b")))
d <- as.data.frame.model.matrix(mm)
print(d); print(names(d)); print(dim(d)); print(attributes(d)$row.names)
#---
mm <- matrix(1:4, 2, dimnames = list(c("r1", "r2"), c("a", "b")))
d <- as.data.frame.model.matrix(mm)
print(d); print(row.names(d))
#---
print(sort.int(c(3, 1, NA, 2)))
#---
print(sort.int(c(3, 1, NA, 2), na.last = TRUE, decreasing = TRUE))
#---
print(sort.int(c(b = 3, a = 1, c = 2)))
#---
print(sort.int(c(3, 1, 2), index.return = TRUE))
#---
print(sort.int(c(3, NA, 1), index.return = TRUE))
#---
print(sort.int(c("b", "A", "a", "B")))
#---
print(sort.list(c(3, 1, NA, 2)))
#---
print(sort.list(c(3, 1, NA, 2), na.last = NA, decreasing = TRUE))
#---
print(sort.list(c("b", "a", "c")))
#---
print(rapply(list(1, 2, list(3, "a")), function(x) x * 2, classes = "numeric"))
#---
print(rapply(list(a = 1, b = list(c = 2, d = "z")), function(x) x + 1, classes = "numeric", how = "replace"))
#---
print(rapply(list(a = 1, b = list(c = 2, d = "z")), function(x) x + 1, classes = "numeric", how = "list"))
#---
print(rapply(list(a = 1, b = list(c = 2, d = "z")), function(x) x + 1, classes = "numeric", deflt = NA, how = "unlist"))
#---
print(rapply(list(1:3, 4:6), sum))
#---
print(try(rapply(1:3, sum), silent = TRUE))
#---
x <- array(1:8, c(2, 2, 2))
print(asplit(x, 3))
#---
m <- matrix(1:6, 2, dimnames = list(c("a", "b"), c("x", "y", "z")))
print(asplit(m, 1))
print(asplit(m, 2))
#---
m <- matrix(1:6, 2, dimnames = list(r = c("a", "b"), c = c("x", "y", "z")))
print(asplit(m, "c"))
#---
print(try(asplit(1:3, 1), silent = TRUE))
#---
print(drop(matrix(1:3, 1)))
#---
print(drop(array(1:6, c(1, 2, 3))))
#---
print(drop(matrix(5, 1, 1, dimnames = list("a", "b"))))
#---
print(drop(matrix(1:2, 1, dimnames = list("a", c("x", "y")))))
#---
print(drop(1:3))
#---
print(drop(array(1:4, c(2, 1, 2), dimnames = list(c("a", "b"), "k", c("u", "v")))))
#---
print(args(function(a, b = 2) a))
#---
print(max.col(matrix(c(1, 5, 3, 9, 2, 2), 2), "first"))
#---
print(max.col(matrix(c(1, 5, 3, 9, 2, 2), 2, byrow = TRUE), "last"))
#---
print(max.col(matrix(c(1, 3, 2, 4, 5, 0), 2), "first"))
#---
print(max.col(matrix(c(1, 3, 2, 4, 5, 0), 2)))
#---
print(t.default(matrix(1:6, 2)))
#---
print(t.default(1:3))
