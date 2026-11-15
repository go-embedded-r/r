x <- as.POSIXct("2024-01-31 10:30:00", tz = "UTC"); seq(x, by = "day", length.out = 3)
#---
x <- as.POSIXct("2024-01-31 10:30:00", tz = "UTC"); seq(x, by = "month", length.out = 4)
#---
x <- as.POSIXct("2024-01-31 10:30:00", tz = "UTC"); seq(x, by = "2 months", length.out = 3)
#---
x <- as.POSIXct("2024-01-31 10:30:00", tz = "UTC"); seq(x, by = "year", length.out = 3)
#---
x <- as.POSIXct("2024-01-31 10:30:00", tz = "UTC"); seq(x, by = "week", length.out = 3)
#---
x <- as.POSIXct("2024-01-31 10:30:00", tz = "UTC"); seq(x, by = "hour", length.out = 3)
#---
x <- as.POSIXct("2024-01-31 10:30:00", tz = "UTC"); seq(x, by = "30 min", length.out = 3)
#---
x <- as.POSIXct("2024-01-31 10:30:00", tz = "UTC"); seq(x, by = 3600, length.out = 3)
#---
x <- as.POSIXct("2024-01-31 10:30:00", tz = "UTC"); seq(x, as.POSIXct("2024-01-31 12:30:00", tz = "UTC"), by = "hour")
#---
x <- as.POSIXct("2024-01-31 10:30:00", tz = "UTC"); seq(x, as.POSIXct("2024-01-31 12:30:00", tz = "UTC"), length.out = 3)
#---
x <- as.POSIXct("2024-01-31 10:30:00", tz = "UTC"); seq(x, by = "-1 day", length.out = 3)
#---
x <- as.POSIXct("2024-01-31 10:30:00", tz = "UTC"); seq(x, as.POSIXct("2024-05-01", tz = "UTC"), by = "month")
#---
x <- as.POSIXct("2024-01-31 10:30:00", tz = "UTC"); seq(x, by = "DSTday", length.out = 3)
#---
x <- as.POSIXct("2024-01-31 10:30:00", tz = "UTC"); seq(x, by = "quarter", length.out = 3)
#---
x <- as.POSIXct("2024-01-31 10:30:00", tz = "UTC"); seq(x, by = "sec", length.out = 3)
#---
x <- as.POSIXct("2024-01-31 10:30:00", tz = "UTC"); class(seq(x, by = "day", length.out = 3))
#---
x <- as.POSIXct("2024-01-31 10:30:00", tz = "UTC"); attributes(seq(x, by = "day", length.out = 2))
#---
x <- as.POSIXct("2024-01-31 10:30:00", tz = "UTC"); tryCatch(seq(x, by = "fortnight", length.out = 3), error = function(e) conditionMessage(e))
#---
x <- as.POSIXct("2024-01-31 10:30:00", tz = "UTC"); tryCatch(seq(x), error = function(e) conditionMessage(e))
#---
x <- as.POSIXct("2024-01-31 10:30:00", tz = "UTC"); seq(x, by = as.difftime(2, units = "hours"), length.out = 3)
#---
seq(as.POSIXct("2024-03-30 12:00", tz = "Europe/Paris"), by = "DSTday", length.out = 3)
#---
seq(as.POSIXct("2024-03-30 12:00", tz = "Europe/Paris"), by = "day", length.out = 3)
#---
seq(as.POSIXct("2024-01-31", tz = "UTC"), by = "month", length.out = 3)
#---
seq(as.POSIXct("2024-01-31", tz = "UTC"), by = "month", length.out = 3)
#---
seq(as.POSIXct("2024-01-31", tz = "UTC"), as.POSIXct("2024-02-03", tz = "UTC"), by = "day")
#---
seq(as.POSIXct("2024-01-31", tz = "UTC"), as.POSIXct("2024-01-28", tz = "UTC"), by = "-1 day")
#---
tryCatch(seq.POSIXt(as.POSIXct("2024-01-31", tz = "UTC"), by = "x y z", length.out = 2), error = function(e) conditionMessage(e))
#---
x <- as.POSIXct("2024-01-31 10:30:00", tz = "UTC"); seq.POSIXt(x, by = "min", length.out = 2)
#---
d <- as.Date("2024-01-31"); round(d)
#---
trunc(as.Date("2024-01-31"), "month")
#---
trunc(as.Date("2024-05-17"), "month")
#---
trunc(as.Date("2024-05-17"), "year")
#---
trunc(as.Date("2024-05-17"), "day")
#---
trunc(as.Date("2024-05-17"))
#---
class(trunc(as.Date("2024-05-17"), "month"))
#---
dput(trunc(as.Date("2024-05-17"), "months"))
#---
dput(trunc(as.Date(c("2024-05-17", NA)), "months"))
#---
dput(trunc(structure(19860.7, class = "Date")))
#---
dput(trunc(structure(19860L, class = "Date")))
#---
dput(trunc(structure(-1.2, class = "Date")))
#---
dput(round(structure(19860.7, class = "Date")))
#---
dput(round(structure(-0.5, class = "Date")))
#---
dput(round(structure(c(a = 19860.5), class = "Date")))
#---
dput(trunc(as.Date(NA), "months"))
#---
x <- as.POSIXct("2024-01-31 10:30:00", tz = "UTC"); trunc(x, "hours")
#---
x <- as.POSIXct("2024-01-31 10:30:00", tz = "UTC"); trunc(x, "days")
#---
x <- as.POSIXct("2024-01-31 10:30:00", tz = "UTC"); trunc(x, "mins")
#---
x <- as.POSIXct("2024-01-31 10:30:00", tz = "UTC"); trunc(x, "secs")
#---
x <- as.POSIXct("2024-01-31 10:30:00", tz = "UTC"); trunc(x, "months")
#---
x <- as.POSIXct("2024-01-31 10:30:00", tz = "UTC"); trunc(x, "years")
#---
x <- as.POSIXct("2024-01-31 10:30:00", tz = "UTC"); class(trunc(x, "days"))
#---
x <- as.POSIXct("2024-01-31 10:30:00", tz = "UTC"); round(x, "hours")
#---
x <- as.POSIXct("2024-01-31 10:30:00", tz = "UTC"); round(x, "days")
#---
x <- as.POSIXct("2024-01-31 10:30:00", tz = "UTC"); round(x, "mins")
#---
round(as.POSIXct("2024-01-31 10:30:31", tz = "UTC"), "mins")
#---
round(as.POSIXct("2024-01-31 10:30:31.7", tz = "UTC"), "secs")
#---
x <- as.POSIXct("2024-01-31 10:30:00", tz = "UTC"); round(x, "months")
#---
x <- as.POSIXct("2024-01-31 10:30:00", tz = "UTC"); round(x, "years")
#---
x <- as.POSIXct("2024-01-31 10:30:00", tz = "UTC"); round(x)
#---
x <- as.POSIXct("2024-01-31 10:30:00", tz = "UTC"); trunc(x)
#---
x <- as.POSIXct("2024-01-31 10:30:00", tz = "UTC"); format(round(x, "days"))
#---
trunc(as.POSIXct(c("2024-01-31 10:30:00", "2024-06-15 23:59:59"), tz = "UTC"), "days")
#---
trunc(as.POSIXct(NA), "days")
#---
x <- as.POSIXct("2024-07-31 10:30:45.5", tz = "Europe/Paris"); dput(unclass(trunc(x, "secs")))
#---
x <- as.POSIXct("2024-07-31 10:30:45.5", tz = "Europe/Paris"); dput(unclass(trunc(x, "mins")))
#---
x <- as.POSIXct("2024-07-31 10:30:45.5", tz = "Europe/Paris"); dput(unclass(trunc(x, "hours")))
#---
x <- as.POSIXct("2024-07-31 10:30:45.5", tz = "Europe/Paris"); dput(unclass(trunc(x, "days")))
#---
x <- as.POSIXct("2024-07-31 10:30:45.5", tz = "Europe/Paris"); dput(unclass(trunc(x, "months")))
#---
x <- as.POSIXct("2024-07-31 10:30:45.5", tz = "Europe/Paris"); dput(unclass(trunc(x, "years")))
#---
x <- as.POSIXct("2024-07-31 10:30:45.5", tz = "Europe/Paris"); dput(unclass(round(x, "secs")))
#---
x <- as.POSIXct("2024-07-31 10:30:45.5", tz = "Europe/Paris"); dput(unclass(round(x, "mins")))
#---
x <- as.POSIXct("2024-07-31 10:30:45.5", tz = "Europe/Paris"); dput(unclass(round(x, "hours")))
#---
x <- as.POSIXct("2024-07-31 10:30:45.5", tz = "Europe/Paris"); dput(unclass(round(x, "days")))
#---
x <- as.POSIXct("2024-07-31 10:30:45.5", tz = "Europe/Paris"); dput(unclass(round(x, "months")))
#---
x <- as.POSIXct("2024-07-31 10:30:45.5", tz = "Europe/Paris"); dput(unclass(round(x, "years")))
#---
y <- as.POSIXct("2024-07-31 10:30:45.5", tz = "UTC"); dput(unclass(trunc(as.POSIXlt(y), "days")))
#---
y <- as.POSIXct("2024-07-31 10:30:45.5", tz = "UTC"); dput(unclass(trunc(as.POSIXct(c(y, NA)), "days")))
#---
y <- as.POSIXct("2024-07-31 10:30:45.5", tz = "UTC"); dput(unclass(trunc(as.POSIXct(c(a = y)), "days")))
#---
dput(unclass(round(as.POSIXct("2024-07-31 23:30:45", tz = "UTC"), "days")))
#---
dput(unclass(round(as.POSIXct("2024-12-31 23:30:45", tz = "UTC"), "years")))
#---
dput(unclass(round(as.POSIXct("2024-12-31 23:30:45", tz = "UTC"), "months")))
#---
dput(round(as.POSIXct("2024-07-31 23:30:45", tz = "UTC")))
#---
dput(trunc(as.POSIXct("2024-07-31 10:30:45", tz = "UTC"), units = "hours"))
#---
f <- function(s, u) format(round(as.POSIXct(s, tz = "UTC"), u), "%Y-%m-%d %H:%M:%S"); c(f("2024-07-16 00:00:00", "months"), f("2024-07-16 11:59:59", "months"), f("2024-07-16 12:00:00", "months"), f("2024-07-17 00:00:00", "months"))
#---
f <- function(s, u) format(round(as.POSIXct(s, tz = "UTC"), u), "%Y-%m-%d %H:%M:%S"); c(f("2024-02-14 00:00:00", "months"), f("2024-02-15 00:00:00", "months"), f("2024-02-15 12:00:00", "months"), f("2024-02-16 00:00:00", "months"))
#---
f <- function(s, u) format(round(as.POSIXct(s, tz = "UTC"), u), "%Y-%m-%d %H:%M:%S"); c(f("2023-02-14 12:00:00", "months"), f("2023-02-15 12:00:00", "months"), f("2024-06-30 23:00:00", "months"), f("2024-06-15 23:00:00", "months"), f("2024-06-16 00:00:00", "months"))
#---
f <- function(s, u) format(round(as.POSIXct(s, tz = "UTC"), u), "%Y-%m-%d %H:%M:%S"); c(f("2024-07-01 12:00:00", "years"), f("2024-06-30 12:00:00", "years"), f("2024-07-02 00:00:00", "years"), f("2024-07-01 00:00:00", "years"))
#---
f <- function(s, u) format(round(as.POSIXct(s, tz = "UTC"), u), "%Y-%m-%d %H:%M:%S"); c(f("2023-07-02 00:00:00", "years"), f("2023-07-01 00:00:00", "years"), f("2023-07-02 12:00:00", "years"), f("2023-07-03 00:00:00", "years"))
#---
f <- function(s, u) format(round(as.POSIXct(s, tz = "UTC"), u), "%Y-%m-%d %H:%M:%S"); c(f("2024-01-01 12:30:30", "days"), f("2024-01-01 12:00:00", "days"), f("2024-01-01 11:59:59", "days"), f("2024-01-01 00:30:00", "hours"), f("2024-01-01 00:29:59", "hours"))
#---
f <- function(s, u) format(round(as.POSIXct(s, tz = "UTC"), u), "%Y-%m-%d %H:%M:%S"); c(f("2024-01-01 00:00:29.9", "mins"), f("2024-01-01 00:00:30", "mins"), f("2024-01-01 00:00:00.5", "secs"), f("2024-01-01 00:00:00.4999", "secs"))
#---
d <- as.Date(c("2024-01-15", "2024-02-20", "2024-02-25", "2024-04-01")); cut(d, "month")
#---
d <- as.Date(c("2024-01-15", "2024-02-20", "2024-02-25", "2024-04-01")); cut(d, "week")
#---
d <- as.Date(c("2024-01-15", "2024-02-20", "2024-02-25", "2024-04-01")); cut(d, "year")
#---
d <- as.Date(c("2024-01-15", "2024-02-20", "2024-02-25", "2024-04-01")); cut(d, "2 months")
#---
d <- as.Date(c("2024-01-15", "2024-02-20", "2024-02-25", "2024-04-01")); cut(d, "day")
#---
d <- as.Date(c("2024-01-15", "2024-02-20", "2024-02-25", "2024-04-01")); cut(d, 2)
#---
d <- as.Date(c("2024-01-15", "2024-02-20", "2024-02-25", "2024-04-01")); cut(d, breaks = as.Date(c("2024-01-01", "2024-03-01", "2024-05-01")))
#---
d <- as.Date(c("2024-01-15", "2024-02-20", "2024-02-25", "2024-04-01")); cut(d, breaks = as.Date(c("2024-01-01", "2024-03-01", "2024-05-01")), labels = c("a", "b"))
#---
d <- as.Date(c("2024-01-15", "2024-02-20", "2024-02-25", "2024-04-01")); cut(d, "month", labels = FALSE)
#---
d <- as.Date(c("2024-01-15", "2024-02-20", "2024-02-25", "2024-04-01")); levels(cut(d, "month"))
#---
d <- as.Date(c("2024-01-15", "2024-02-20", "2024-02-25", "2024-04-01")); levels(cut(d, "quarter"))
#---
d <- as.Date(c("2024-01-15", "2024-02-20", "2024-02-25", "2024-04-01")); table(cut(d, "month"))
#---
d <- as.Date(c("2024-01-15", "2024-02-20", "2024-02-25", "2024-04-01")); cut(d, "week", start.on.monday = FALSE)
#---
d <- as.Date(c("2024-01-15", "2024-02-20", "2024-02-25", "2024-04-01")); cut(d, "month", right = TRUE)
#---
d <- as.Date(c("2024-01-15", "2024-02-20", "2024-02-25", "2024-04-01")); cut(d, "quarter")
#---
d <- as.Date(c("2024-01-15", "2024-02-20", "2024-02-25", "2024-04-01")); tryCatch(cut(d, "fortnight"), error = function(e) conditionMessage(e))
#---
d <- as.Date(c("2024-01-15", "2024-02-20", "2024-02-25", "2024-04-01")); cut(d, breaks = as.Date(c("2024-01-01", "2024-03-01", "2024-05-01")), right = TRUE)
#---
p <- as.POSIXct(c("2024-01-15 10:00:00", "2024-01-15 11:30:00", "2024-01-16 09:00:00"), tz = "UTC"); cut(p, "day")
#---
p <- as.POSIXct(c("2024-01-15 10:00:00", "2024-01-15 11:30:00", "2024-01-16 09:00:00"), tz = "UTC"); cut(p, "hour")
#---
p <- as.POSIXct(c("2024-01-15 10:00:00", "2024-01-15 11:30:00", "2024-01-16 09:00:00"), tz = "UTC"); cut(p, "2 hours")
#---
p <- as.POSIXct(c("2024-01-15 10:00:00", "2024-01-15 11:30:00", "2024-01-16 09:00:00"), tz = "UTC"); cut(p, "month")
#---
p <- as.POSIXct(c("2024-01-15 10:00:00", "2024-01-15 11:30:00", "2024-01-16 09:00:00"), tz = "UTC"); cut(p, 2)
#---
p <- as.POSIXct(c("2024-01-15 10:00:00", "2024-01-15 11:30:00", "2024-01-16 09:00:00"), tz = "UTC"); cut(p, "30 min")
#---
p <- as.POSIXct(c("2024-01-15 10:00:00", "2024-01-15 11:30:00", "2024-01-16 09:00:00"), tz = "UTC"); cut(p, breaks = as.POSIXct(c("2024-01-15", "2024-01-16", "2024-01-17"), tz = "UTC"))
#---
p <- as.POSIXct(c("2024-01-15 10:00:00", "2024-01-15 11:30:00", "2024-01-16 09:00:00"), tz = "UTC"); cut(p, "week")
#---
d <- as.Date(c("2024-01-15", "2024-02-20", "2024-02-25", "2024-04-01")); summary(d)
#---
d <- as.Date(c("2024-01-15", "2024-02-20", "2024-02-25", "2024-04-01")); class(summary(d))
#---
d <- as.Date(c("2024-01-15", "2024-02-20", "2024-02-25", "2024-04-01")); dput(summary(d))
#---
d <- as.Date(c("2024-01-15", "2024-02-20", "2024-02-25", "2024-04-01")); dput(summary(c(d, NA)))
#---
p <- as.POSIXct(c("2024-01-15 10:00:00", "2024-01-15 11:30:00", NA), tz = "UTC"); dput(summary(p))
#---
p <- as.POSIXct(c("2024-01-15 10:00:00", "2024-01-15 11:30:00", NA), tz = "UTC"); dput(summary(as.POSIXlt(p)))
#---
p <- as.POSIXct(c("2024-01-15 10:00:00", "2024-01-15 11:30:00", "2024-01-16 09:00:00"), tz = "UTC"); summary(p)
#---
d <- as.Date(c("2024-01-15", "2024-02-20", "2024-02-25", "2024-04-01")); summary(d, digits = 3)
#---
d <- as.Date(c("2024-01-15", "2024-02-20", "2024-02-25", "2024-04-01")); summary(c(d, NA))
#---
summary(as.POSIXct(c("2024-01-15 10:00:00", NA), tz = "UTC"))
#---
p <- as.POSIXct(c("2024-01-15 10:00:00", "2024-01-15 11:30:00", "2024-01-16 09:00:00"), tz = "UTC"); summary(as.POSIXlt(p))
#---
summary(as.difftime(c(1, 5, 10), units = "mins"))
#---
dput(summary(as.difftime(c(1, 5, 10), units = "mins")))
#---
dput(summary(as.difftime(c(1, 5, NA), units = "mins")))
#---
summary(as.difftime(c(1, 5, NA), units = "mins"))
#---
d <- as.Date(c("2024-01-15", "2024-02-20", "2024-02-25", "2024-04-01")); print(summary(d), quote = FALSE)
#---
print(summary(as.difftime(c(1, 5, 10), units = "mins")), digits = 2)
#---
d <- as.Date(c("2024-01-15", "2024-02-20", "2024-02-25", "2024-04-01")); format(summary(d))
#---
d <- as.Date(c("2024-01-15", "2024-02-20", "2024-02-25", "2024-04-01")); as.character(summary(d))
#---
summary(as.Date(c("2024-01-15", NA, NA)))
#---
x <- as.POSIXct("2024-01-31 10:30:00", tz = "UTC"); trunc.POSIXt(x, "hours")
#---
x <- as.POSIXct("2024-01-31 10:30:00", tz = "UTC"); round.POSIXt(x, "hours")
#---
trunc.Date(as.Date("2024-05-17"), "month")
#---
d <- as.Date(c("2024-01-15", "2024-02-20")); cut.Date(d, "month")
#---
p <- as.POSIXct(c("2024-01-15 10:00:00", "2024-01-16 09:00:00"), tz = "UTC"); cut.POSIXt(p, "day")
#---
summary.Date(as.Date(c("2024-01-15", "2024-02-20")))
#---
summary.POSIXct(as.POSIXct(c("2024-01-15 10:00:00", "2024-01-16 09:00:00"), tz = "UTC"))
#---
summary.POSIXlt(as.POSIXlt(as.POSIXct(c("2024-01-15 10:00:00", "2024-01-16 09:00:00"), tz = "UTC")))
#---
summary.difftime(as.difftime(c(1, 5), units = "mins"))
#---
as.character(as.POSIXct(c("2024-01-15 10:00:00", "2024-01-16 00:00:00", "2024-01-16 00:30:00"), tz = "UTC"))
#---
as.character(as.POSIXct("2024-01-15 10:30:00", tz = "UTC"))
#---
tryCatch(round.Date(as.Date("2024-05-17")), error = function(e) "failed")
