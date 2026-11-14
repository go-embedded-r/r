data.frame(a = 1:3, b = c("x", "y", "z"))
#---
attributes(data.frame(a = 1:3, b = c("x", "y", "z")))
#---
data.frame(a = 1:4, b = 1:2)
#---
data.frame(a = 1:3, b = 1:2)
#---
data.frame(x = c(a = 1, b = 2))
#---
attributes(data.frame(x = c(a = 1, b = 2)))
#---
data.frame(a = 1:2, row.names = c("p", "q"))
#---
attributes(data.frame(a = 1:2, b = 3:4, row.names = 2))
#---
data.frame(a = 1, a = 2)
#---
data.frame(a = 1, a = 2, check.names = FALSE)
#---
names(data.frame(`a b` = 1, `1x` = 2))
#---
x <- 1:2; names(data.frame(x, x * 2))
#---
names(data.frame(1:2, c("a", "b")))
#---
names(data.frame(a = data.frame(x = 1)))
#---
names(data.frame(a = data.frame(x = 1, y = 2)))
#---
names(data.frame(matrix(1:4, 2)))
#---
names(data.frame(m = matrix(1:4, 2)))
#---
data.frame(matrix(1:4, 2, dimnames = list(c("r1", "r2"), c("p", "q"))))
#---
data.frame(a = c("u", "v"), stringsAsFactors = TRUE)$a
#---
str1 <- data.frame(a = c("u", "v", "u"), stringsAsFactors = TRUE); levels(str1$a)
#---
attributes(data.frame())
#---
data.frame()
#---
data.frame(a = integer(0))
#---
data.frame(a = 1:2, b = letters[1:2])[0, ]
#---
data.frame(a = c(1.5, 10, NA), b = c("x", NA, "zz"), c = c(TRUE, FALSE, NA))
#---
data.frame(a = factor(c("u", NA)), b = c(0.001234, 123.4))
#---
is.data.frame(data.frame(a = 1)); is.data.frame(list(a = 1))
#---
as.data.frame(matrix(1:4, 2))
#---
attributes(as.data.frame(matrix(1:4, 2)))
#---
as.data.frame(matrix(1:4, 2, dimnames = list(NULL, c("p", "q"))), stringsAsFactors = FALSE)
#---
as.data.frame(1:2)
#---
attributes(as.data.frame(1:2))
#---
as.data.frame(c(a = 1, b = 2))
#---
as.data.frame(list(a = 1:2, b = c("q", "r")))
#---
as.data.frame(letters[1:2])
#---
df <- data.frame(a = 1:3, b = c("x", "y", "z")); df[2:3, ]
#---
df <- data.frame(a = 1:3, b = c("x", "y", "z")); attributes(df[2:3, ])
#---
df <- data.frame(a = 1:3, b = c("x", "y", "z")); df[df$a > 1, "b"]
#---
df <- data.frame(a = 1:3, b = c("x", "y", "z")); df[["a"]]
#---
df <- data.frame(a = 1:3, b = c("x", "y", "z")); df$b
#---
df <- data.frame(a = 1:3, b = c("x", "y", "z")); df[2]
#---
df <- data.frame(a = 1:3, b = c("x", "y", "z")); attributes(df[2])
#---
df <- data.frame(a = 1:3, b = c("x", "y", "z")); df[, 1]
#---
df <- data.frame(a = 1:3, b = c("x", "y", "z")); df[, 1, drop = FALSE]
#---
df <- data.frame(a = 1:3, b = c("x", "y", "z")); df[c(1, 1), ]
#---
df <- data.frame(a = 1:3, b = c("x", "y", "z")); df[c(1, NA), ]
#---
df <- data.frame(a = 1:3, b = c("x", "y", "z")); df[-1, ]
#---
df <- data.frame(a = 1:3, b = c("x", "y", "z")); df[c("a", "zz")]
#---
df <- data.frame(a = 1:3, b = c("x", "y", "z")); df[2, ]
#---
df <- data.frame(a = 1:3, b = c("x", "y", "z")); df[2, , drop = TRUE]
#---
df <- data.frame(a = 1:3, b = c("x", "y", "z")); df[[2, "b"]]
#---
df <- data.frame(a = 1:3, b = c("x", "y", "z")); dim(df); nrow(df); ncol(df); names(df); dimnames(df); row.names(df); rownames(df); colnames(df)
#---
df <- data.frame(a = 1:3, b = c("x", "y", "z")); df$c <- df$a * 2; df
#---
df <- data.frame(a = 1:3, b = c("x", "y", "z")); df$c <- df$a * 2; attributes(df)
#---
df <- data.frame(a = 1:3, b = c("x", "y", "z")); df$b <- NULL; df
#---
df <- data.frame(a = 1:3, b = c("x", "y", "z")); df$c <- 1:2
#---
df <- data.frame(a = 1:4); df$c <- 1:2; df
#---
df <- data.frame(a = 1:3, b = c("x", "y", "z")); df[["c"]] <- 3:1; df
#---
df <- data.frame(a = 1:3, b = c("x", "y", "z")); df[2, "a"] <- 10L; df
#---
df <- data.frame(a = 1:3, b = c("x", "y", "z")); df[df$a > 1, "a"] <- 0L; df
#---
df <- data.frame(a = 1:3, b = c("x", "y", "z")); df["a"] <- list(c(7, 8, 9)); df
#---
df <- data.frame(a = 1:3, b = c("x", "y", "z")); df[4, ] <- list(4L, "w"); df
#---
df <- data.frame(a = 1:3, b = c("x", "y", "z")); df[, "z"] <- 0; df
#---
df <- data.frame(a = 1:3, b = c("x", "y", "z")); names(df) <- c("p", "q"); df
#---
df <- data.frame(a = 1:3, b = c("x", "y", "z")); names(df)[2] <- "w"; df
#---
df <- data.frame(a = 1:3, b = c("x", "y", "z")); row.names(df) <- c("r1", "r2", "r3"); df
#---
df <- data.frame(a = 1:3, b = c("x", "y", "z")); rownames(df) <- c("r1", "r2", "r2")
#---
df <- data.frame(a = 1:3, b = c("x", "y", "z"), row.names = c("r1", "r2", "r3")); df["r2", ]
#---
df <- data.frame(a = 1:3, b = c("x", "y", "z")); rownames(df) <- NULL; attributes(df)
#---
rbind(data.frame(a = 1, b = "x"), data.frame(a = 2, b = "y"))
#---
attributes(rbind(data.frame(a = 1, b = "x"), data.frame(a = 2, b = "y")))
#---
rbind(data.frame(a = 1, b = "x"), data.frame(a = 2:3, b = c("y", "z")))
#---
rbind(data.frame(a = 1, b = "x"), list(5, "w"))
#---
rbind(data.frame(a = 1, b = factor("x")), data.frame(a = 2, b = factor("y")))$b
#---
rbind(data.frame(a = 1, row.names = "p"), data.frame(a = 2, row.names = "q"))
#---
rbind(u = data.frame(a = 1), v = data.frame(a = 2))
#---
rbind(data.frame(a = 1), data.frame(b = 2))
#---
cbind(data.frame(a = 1:2), b = c("x", "y"))
#---
attributes(cbind(data.frame(a = 1), b = 2))
#---
cbind(data.frame(a = 1:2), data.frame(b = 3:4))
#---
merge(data.frame(k = c(1, 2, 3), a = 4:6), data.frame(k = c(3, 1, 5), b = 7:9))
#---
merge(data.frame(k = c(1, 2, 3), a = 4:6), data.frame(k = c(3, 1, 5), b = 7:9), all = TRUE)
#---
merge(data.frame(k = c(1, 2, 3), a = 4:6), data.frame(k = c(3, 1, 5), b = 7:9), all.x = TRUE)
#---
merge(data.frame(k = c("a", "b", "a"), v = 1:3), data.frame(k = c("a", "a", "c"), w = 4:6))
#---
merge(data.frame(k1 = 1:3, v = 4:6), data.frame(k2 = c(3, 1), w = 7:8), by.x = "k1", by.y = "k2")
#---
merge(data.frame(k = 1:2, v = 1:2), data.frame(k = 1:2, v = 3:4), by = "k")
#---
attributes(merge(data.frame(k = 1:2, v = 1:2), data.frame(k = 1:2, w = 3:4)))
#---
merge(data.frame(a = 1:2), data.frame(b = 3:4))
#---
split(data.frame(a = 1:4, g = c("u", "v", "u", "v")), c("u", "v", "u", "v"))
#---
df <- data.frame(a = 1:4, g = c("u", "v", "u", "v")); split(df, df$g)
#---
df <- data.frame(a = 1:4, g = c("u", "v", "u", "v")); lapply(df, class)
#---
df <- data.frame(a = 1:4, b = c(2.5, 3, 1, 8)); sapply(df, mean)
#---
df <- data.frame(a = 1:4, b = c(2.5, 3, 1, 8)); with(df, a + b)
#---
df <- data.frame(a = 1:4, b = c(2.5, 3, 1, 8)); subset(df, a > 2)
#---
df <- data.frame(a = 1:4, b = c(2.5, 3, 1, 8)); subset(df, a > 2, select = b)
#---
df <- data.frame(a = 1:4, b = c(2.5, 3, 1, 8)); subset(df, b < 5 & a > 1, select = c(b, a))
#---
df <- data.frame(a = 1:4, b = c(2.5, 3, 1, 8)); transform(df, c = a * 2)
#---
df <- data.frame(a = 1:4, b = c(2.5, 3, 1, 8)); transform(df, a = a + 100L, d = b / 2)
#---
df <- data.frame(a = 1:4, b = c(2.5, 3, 1, 8)); df[order(df$b), ]
#---
df <- data.frame(a = c(2, 1, 2), b = c("x", "y", "x")); order(df$a, df$b); unique(df); duplicated(df)
#---
df <- data.frame(a = 1:2, b = c("x", "y")); as.matrix(df)
#---
df <- data.frame(a = c(1, 10), b = 3:4); as.matrix(df)
#---
df <- data.frame(a = c(1, NA), b = c("x", NA)); is.na(df)
#---
df <- data.frame(a = 1:2, b = c("x", "y")); t(df)
#---
df <- data.frame(a = 1:2, b = c("x", "y")); format(df)
#---
df <- data.frame(a = 1:2, b = c("x", "y")); as.list(df)
#---
df <- data.frame(a = 1:2, b = c("x", "y")); length(df); unclass(df)
#---
print(data.frame(a = 1:2, b = c("x", "yy")), row.names = FALSE)
#---
data.frame(a = factor(c("u", "v", "u")), b = 1:3)[2:3, ]
#---
df <- data.frame(g = factor(c("u", "v", "u"))); droplevels(df[df$g == "u", , drop = FALSE])$g
#---
f <- factor(c("a", "b", "a")); f[2] <- "a"; f
#---
f <- factor(c("a", "b", "a")); f[[1]] <- "q"; f
#---
c(factor("a"), factor(c("b", "a")))
#---
addNA(factor(c("a", NA)))
#---
levels(factor(c("b", "a", NA), exclude = NULL))
#---
factor(c(1, 2, NA), exclude = NULL)
#---
droplevels(factor(c("a", "b"))[1])
#---
make.names(c("a b", "1x", "if", "", NA, "ok", "_a", ".1", "a-b"))
#---
make.names(c("a", "a", "b"), unique = TRUE)
#---
.bincode(c(1, 5, 10, 11, NA), c(0, 5, 10))
#---
.bincode(c(0, 5, 10), c(0, 5, 10), TRUE, TRUE)
#---
.bincode(c(0, 5, 10), c(0, 5, 10), FALSE, TRUE)
#---
.row_names_info(data.frame(a = 1:3)); .row_names_info(data.frame(a = 1:3), 2L); .row_names_info(data.frame(a = 1:3, row.names = c("a", "b", "c")))
#---
.set_row_names(4)
#---
x <- structure(list(a = 1), row.names = c(NA, -3L), class = "data.frame"); attr(x, "row.names")
#---
cbind(1:2, c(a = 3, b = 4))
#---
a <- 1:2; cbind(a, b = 3:4, 5)
#---
rbind(1:3, c(x = 1, y = 2, z = 3))
#---
rbind(a = 1:2, b = 3:4)
#---
cbind(matrix(1:4, 2, dimnames = list(c("r1", "r2"), c("A", "B"))), z = 9:10)
#---
rbind(1:2, matrix(5:8, 2))
#---
cbind(1:3, 1:2)
#---
cbind(c("a", "b"), 1:2)
#---
do.call(rbind, list(1:2, 3:4))
#---
cbind(NULL, 1:2)
#---
rbind(list(1, 2), 3:4)
#---
cbind(matrix(1:4, 2), matrix(1:6, 3))
#---
print(matrix(c("a", "bb", "ccc", "d"), 2, dimnames = list(c("r1", "r2"), c("A", "B"))), quote = FALSE, right = TRUE)
#---
f <- function(x, i, j) nargs(); x <- structure(list(1), class = "foo"); "[.foo" <- function(x, i, j) nargs(); c(x[1], x[1, ], x[, 2], x[1, 2], x[])
