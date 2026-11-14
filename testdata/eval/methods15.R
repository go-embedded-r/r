e <- simpleError("boom", quote(f(x))); conditionMessage.condition(e)
#---
e <- simpleError("boom", quote(f(x))); conditionCall.condition(e)
#---
e <- simpleError("boom", quote(f(x))); as.character(e)
#---
e <- simpleError("boom"); as.character(e)
#---
w <- simpleWarning("careful", quote(g())); as.character(w)
#---
w <- simpleCondition("plain"); class(w) <- c("custom", "condition"); as.character.condition(w)
#---
e <- simpleError("boom", quote(f(x))); as.character.error(e)
#---
e <- simpleError("boom"); as.character.error(e)
#---
e <- simpleError("boom", quote(f(x))); print(e)
#---
e <- simpleError("boom"); print.condition(e)
#---
w <- simpleWarning("careful", quote(g())); w
#---
pmax.int(c(1, 5, 3), c(4, 2, 6))
#---
pmin.int(c(1, 5, NA), 3, na.rm = TRUE)
#---
pmax.int(1:3, 2)
#---
d <- as.Date(c("2024-03-05", "2024-01-01", NA, "2025-12-31")); mtfrm(d)
#---
mtfrm(factor(c("a", "b")))
#---
mtfrm(1:2)
#---
mtfrm(as.POSIXct("2024-01-01 10:00:00", tz = "UTC"))
#---
mtfrm(list(1, 2))
#---
mtfrm(list(1:2, 3))
#---
mtfrm.default(c(a = 1.5, b = 2))
#---
mtfrm.Date(as.Date("2024-01-02"))
#---
mtfrm.POSIXct(as.POSIXct("2024-01-01 10:00:00", tz = "UTC"))
#---
d <- as.Date(c("2024-03-05", "2024-01-01", NA, "2025-12-31")); xtfrm(d)
#---
xtfrm.Date(as.Date("2024-01-02"))
#---
xtfrm.POSIXct(as.POSIXct("2024-01-01 10:00:00", tz = "UTC"))
#---
xtfrm(factor(c("b", "a")))
#---
xtfrm.factor(factor(c("b", "a", "b")))
#---
xtfrm(as.difftime(c(3, 1), units = "mins"))
#---
xtfrm.difftime(as.difftime(c(3, 1), units = "mins"))
#---
xtfrm(I(c(3, 1)))
#---
xtfrm.AsIs(I(c(3, 1)))
#---
xtfrm(data.frame(a = c(2, 1)))
#---
xtfrm.data.frame(data.frame(a = 1))
#---
d <- as.Date(c("2024-03-05", "2024-01-01", NA, "2025-12-31")); range(d)
#---
d <- as.Date(c("2024-03-05", "2024-01-01", NA, "2025-12-31")); range(d, na.rm = TRUE)
#---
range(as.POSIXct(c("2024-01-01 10:00:00", "2024-01-02 11:00:00"), tz = "UTC"))
#---
range.Date(as.Date(c("2024-03-05", "2024-01-01")), as.Date("2026-01-01"))
#---
range.POSIXct(as.POSIXct(c("2024-01-01 10:00:00", "2024-01-02 11:00:00"), tz = "UTC"))
#---
d <- as.Date(c("2024-03-05", "2024-01-01", NA, "2025-12-31")); split(d, c(1, 2, 1, 2))
#---
split(as.POSIXct(c("2024-01-01", "2024-01-02"), tz = "UTC"), 1:2)[[1]]
#---
split.Date(as.Date(c("2024-03-05", "2024-01-01")), c("a", "b"))
#---
split.POSIXct(as.POSIXct(c("2024-01-01", "2024-01-02"), tz = "UTC"), c(1, 1))
#---
x <- as.Date(c("2024-03-05", "2024-01-01")); length(x) <- 4; x
#---
`length<-.Date`(as.Date("2024-03-05"), 2)
#---
f <- factor(c("a", "b")); length(f) <- 3; f
#---
`length<-.factor`(factor(c("a", "b")), 1)
#---
dt <- as.difftime(1:2, units = "mins"); length(dt) <- 3; dt
#---
`length<-.difftime`(as.difftime(1:2, units = "secs"), 1)
#---
x <- as.POSIXct(c("2024-01-01", "2024-01-02"), tz = "UTC"); length(x) <- 3; x
#---
`length<-.POSIXct`(as.POSIXct("2024-01-01", tz = "UTC"), 2)
#---
tryCatch(sqrt(as.Date("2024-01-01")), error = function(e) cat(conditionMessage(e), "\n"))
#---
tryCatch(cumsum(as.Date("2024-01-01")), error = function(e) cat(conditionMessage(e), "\n"))
#---
tryCatch(Math.Date(as.Date("2024-01-01")), error = function(e) cat(conditionMessage(e), "\n"))
#---
tryCatch(sqrt(as.POSIXct("2024-01-01", tz = "UTC")), error = function(e) cat(conditionMessage(e), "\n"))
#---
tryCatch(Math.POSIXt(as.POSIXct("2024-01-01", tz = "UTC")), error = function(e) cat(conditionMessage(e), "\n"))
#---
tryCatch(sqrt(factor("a")), error = function(e) cat(conditionMessage(e), "\n"))
#---
tryCatch(Math.factor(factor("a")), error = function(e) cat(conditionMessage(e), "\n"))
#---
tryCatch(max(factor(c("a", "b"))), error = function(e) cat(conditionMessage(e), "\n"))
#---
tryCatch(Summary.factor(factor(c("a", "b"))), error = function(e) cat(conditionMessage(e), "\n"))
#---
max(ordered(c("a", "b")))
#---
range(ordered(c("b", "a", "c")))
#---
tryCatch(sum(ordered(c("a", "b"))), error = function(e) cat(conditionMessage(e), "\n"))
#---
Summary.ordered(ordered(c("a", "b")))
#---
withCallingHandlers(factor(c("a", "b")) + 1, warning = function(w) { print(conditionMessage(w)); invokeRestart("muffleWarning") })
#---
tryCatch(Ops.factor(factor("a"), 1), warning = function(w) cat(conditionMessage(w), "\n"))
#---
factor(c("a", "b", "a")) == "a"
#---
tryCatch(ordered(c("a", "b")) < ordered("a", levels = c("b", "a")), error = function(e) cat(conditionMessage(e), "\n"))
#---
ordered(c("a", "b")) < "b"
#---
Ops.ordered(ordered(c("a", "b")), "b")
#---
df <- data.frame(a = c(1, 4, 9), b = c(16, 25, 36)); sqrt(df)
#---
df <- data.frame(a = c(1, 4, 9), b = c(16, 25, 36)); Math.data.frame(df)
#---
df <- data.frame(a = c(1.5, -2.5), b = c(0.5, 3)); round(df)
#---
df <- data.frame(a = c(1, 4), b = c("x", "y")); tryCatch(sqrt(df), error = function(e) cat(conditionMessage(e), "\n"))
#---
df <- data.frame(a = c(1, 4), b = c(2, 3)); sum(df)
#---
df <- data.frame(a = c(1, 4), b = c(2, 3)); max(df, 10)
#---
df <- data.frame(a = c(1, 4), b = c(2, 3)); range(df)
#---
df <- data.frame(a = c(1, NA), b = c(2, 3)); sum(df, na.rm = TRUE)
#---
tryCatch(sum(data.frame(a = "x")), error = function(e) cat(conditionMessage(e), "\n"))
#---
Summary.data.frame(data.frame(a = 1:3), na.rm = FALSE)
#---
df <- data.frame(a = 1:3, b = c(0.5, 1, 1.5)); df * 2
#---
df <- data.frame(a = 1:3, b = c(0.5, 1, 1.5)); 10 - df
#---
df <- data.frame(a = 1:3, b = c(0.5, 1, 1.5)); df + df
#---
df <- data.frame(a = 1:3, b = c(0.5, 1, 1.5)); df > 1
#---
df <- data.frame(a = 1:3, b = c(0.5, 1, 1.5)); df == df
#---
df <- data.frame(a = 1:3, b = c(0.5, 1, 1.5)); -df
#---
df <- data.frame(a = c(TRUE, FALSE), b = c(FALSE, FALSE)); !df
#---
df <- data.frame(a = 1:2, b = 3:4, row.names = c("r1", "r2")); df * c(1, 10)
#---
df <- data.frame(a = 1:2, b = 3:4); tryCatch(df + data.frame(z = 1), error = function(e) cat(conditionMessage(e), "\n"))
#---
df <- data.frame(a = 1:4, b = 5:8); df %% 2 == 0
#---
Ops.data.frame(data.frame(a = 1:2), 5)
#---
df <- data.frame(a = 1:3, b = 4:6); within(df, c <- a + b)
#---
df <- data.frame(a = 1:3, b = 4:6); within(df, { c <- a + b; d <- c * 2 })
#---
df <- data.frame(a = 1:3, b = 4:6); within(df, rm(b))
#---
df <- data.frame(a = 1:3, b = 4:6); within(df, a <- a * 10)
#---
df <- data.frame(a = 1:3, b = 4:6); within(df, b <- NULL)
#---
within.data.frame(data.frame(a = 1:2), z <- a - 1)
