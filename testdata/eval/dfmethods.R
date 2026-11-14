df <- data.frame(a = 1:3, b = c("x", "y", "z")); `[.data.frame`(df, 2:3, )
#---
df <- data.frame(a = 1:3, b = c("x", "y", "z")); `[[.data.frame`(df, "b")
#---
df <- data.frame(a = 1:3, b = c("x", "y", "z")); `$<-.data.frame`(df, "c", 3:1)
#---
df <- data.frame(a = 1:3, b = c("x", "y", "z")); `[[<-.data.frame`(df, "c", value = 3:1)
#---
df <- data.frame(a = 1:3, b = c("x", "y", "z")); `[<-.data.frame`(df, 2, "a", value = 9L)
#---
f <- factor(c("a", "b")); `[<-.factor`(f, 1, value = "b")
#---
f <- factor(c("a", "b")); `[[<-.factor`(f, 2, value = "a")
#---
df <- data.frame(a = 1:3); `.rowNamesDF<-`(df, value = c("p", "q", "r"))
#---
df <- data.frame(a = 1:3); `row.names<-.data.frame`(df, value = c("p", "q", "r"))
#---
m <- matrix(1:4, 2); `row.names<-.default`(m, value = c("p", "q"))
#---
df <- data.frame(a = 1:3, row.names = c("p", "q", "r")); row.names.data.frame(df)
#---
row.names.default(matrix(1:4, 2, dimnames = list(c("p", "q"), NULL)))
#---
df <- data.frame(a = 1:3); dim.data.frame(df); dimnames.data.frame(df)
#---
df <- data.frame(a = 1:2); `dimnames<-.data.frame`(df, list(c("p", "q"), "z"))
#---
as.data.frame.Date(as.Date("2020-01-02"))
#---
as.data.frame.POSIXct(as.POSIXct("2020-01-02", tz = "UTC"))
#---
as.data.frame.character(c("a", "b"))
#---
as.data.frame.complex(1i)
#---
df <- data.frame(a = 1:2); as.data.frame.data.frame(df, row.names = c("p", "q"))
#---
as.data.frame.default(sum)
#---
as.data.frame.factor(factor(c("a", "b")))
#---
as.data.frame.integer(1:2)
#---
as.data.frame.list(list(a = 1:2, b = 3:4))
#---
as.data.frame.logical(c(TRUE, NA))
#---
as.data.frame.matrix(matrix(1:4, 2))
#---
as.data.frame.numeric(c(1.5, 2))
#---
as.data.frame.vector(1:3)
#---
as.data.frame.table(table(c("a", "b", "a")))
#---
as.list.data.frame(data.frame(a = 1:2))
#---
as.matrix.data.frame(data.frame(a = 1:2, b = 3:4))
#---
c.factor(factor("a"), factor("b"))
#---
cbind.data.frame(data.frame(a = 1:2), b = 3:4)
#---
rbind.data.frame(data.frame(a = 1), data.frame(a = 2))
#---
cut.default(c(1, 5), c(0, 2, 6))
#---
droplevels.factor(factor(c("a", "b"))[1])
#---
droplevels.data.frame(data.frame(g = factor(c("a", "b")))[1, , drop = FALSE])$g
#---
duplicated.data.frame(data.frame(a = c(1, 1)))
#---
unique.data.frame(data.frame(a = c(1, 1)))
#---
format.data.frame(data.frame(a = c(1.5, 10)))
#---
is.na.data.frame(data.frame(a = c(1, NA)))
#---
merge.data.frame(data.frame(k = 1:2, a = 3:4), data.frame(k = 2:3, b = 5:6))
#---
print.data.frame(data.frame(a = 1:2))
#---
print.table(table(c("a", "b")))
#---
split.data.frame(data.frame(a = 1:2), c("u", "v"))
#---
subset.data.frame(data.frame(a = 1:3), a > 1)
#---
subset.default(1:5, c(TRUE, FALSE, TRUE, NA, FALSE))
#---
t.data.frame(data.frame(a = 1:2, b = 3:4))
#---
transform.data.frame(data.frame(a = 1:2), b = a * 2)
