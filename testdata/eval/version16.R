x <- numeric_version(c("1.2.3", "10.0", "2")); x
#---
x <- numeric_version(c("1.2.3", "10.0", "2")); unclass(x)
#---
x <- numeric_version(c("1.2.3", "10.0", "2")); class(x)
#---
x <- numeric_version(c("1.2.3", "10.0", "2")); format(x)
#---
x <- numeric_version(c("1.2.3", "10.0", "2")); as.character(x)
#---
x <- numeric_version(c("1.2.3", "10.0", "2")); x[2]
#---
x <- numeric_version(c("1.2.3", "10.0", "2")); x[[1]]
#---
x <- numeric_version(c("1.2.3", "10.0", "2")); x[[1, 2]]
#---
x <- numeric_version(c("1.2.3", "10.0", "2")); x[[1]][2]
#---
x <- numeric_version(c("1.2.3", "10.0", "2")); x > "2.0"
#---
x <- numeric_version(c("1.2.3", "10.0", "2")); x == "2"
#---
x <- numeric_version(c("1.2.3", "10.0", "2")); max(x)
#---
x <- numeric_version(c("1.2.3", "10.0", "2")); range(x)
#---
x <- numeric_version(c("1.2.3", "10.0", "2")); sort(x)
#---
x <- numeric_version(c("1.2.3", "10.0", "2")); c(x, numeric_version("3.1"))
#---
x <- numeric_version(c("1.2.3", "10.0", "2")); rev(x)
#---
x <- numeric_version(c("1.2.3", "10.0", "2")); rep(x, 2)
#---
x <- numeric_version(c("1.2.3", "10.0", "2")); unique(c(x, x))
#---
x <- numeric_version(c("1.2.3", "10.0", "2")); duplicated(c(x, x))
#---
x <- numeric_version(c("1.2.3", "10.0", "2")); is.na(x)
#---
x <- numeric_version(c("1.2.3", "10.0", "2")); is.numeric_version(x)
#---
is.numeric_version(1)
#---
p <- package_version("3.4.5"); p$major
#---
p <- package_version("3.4.5"); p$minor
#---
p <- package_version("3.4.5"); p$patchlevel
#---
p <- package_version("3.4.5"); p[["major"]]
#---
p <- package_version("3.4.5"); class(p)
#---
x <- numeric_version(c("1.2.3", "10.0", "2")); unlist(x)
#---
x <- numeric_version(c("1.2.3", "10.0", "2")); as.list(x)
#---
x <- numeric_version(c("1.2.3", "10.0", "2")); length(x)
#---
numeric_version("1.a", strict = FALSE)
#---
cat(tryCatch(numeric_version("1.a"), error = function(e) conditionMessage(e)), "\n")
#---
cat(tryCatch(package_version("1"), error = function(e) conditionMessage(e)), "\n")
#---
numeric_version("1.a")
#---
numeric_version(1.5)
#---
package_version("1")
#---
as.numeric_version("1.2")
#---
as.package_version("1.2-3")
#---
numeric_version("1.2-3")
#---
numeric_version(character())
#---
numeric_version(c("1.2", "3"))[5]
#---
numeric_version(c("1.2", "3"))[-1]
#---
numeric_version(NA_character_)
#---
cat(tryCatch(numeric_version(c("1.a", "b")), error = function(e) conditionMessage(e)), "\n")
#---
cat(tryCatch(numeric_version(""), error = function(e) conditionMessage(e)), "\n")
#---
numeric_version("", strict = FALSE)
#---
cat(tryCatch(package_version(c("1.2", "3")), error = function(e) conditionMessage(e)), "\n")
#---
numeric_version(" 1.2 ", strict = FALSE)
#---
numeric_version("1..2", strict = FALSE)
#---
tryCatch(numeric_version(list(1)), error = function(e) conditionMessage(e))
#---
numeric_version("01.002")
#---
unclass(numeric_version("01.002"))
#---
numeric_version("1.2") == "1.2.0"
#---
numeric_version("1.2") < "1.2.0"
#---
numeric_version("1") == "1.0.0"
#---
numeric_version("1") < "1.0.0.1"
#---
numeric_version("1.2.0.1") > "1.2"
#---
numeric_version(c("1.2", NA), strict = FALSE) < "1.0"
#---
numeric_version("1.2") + 1
#---
tryCatch(numeric_version("1.2") + 1, error = function(e) conditionMessage(e))
#---
tryCatch(numeric_version("1.2") & 1, error = function(e) conditionMessage(e))
#---
numeric_version("1.2") > 1
#---
tryCatch(sum(numeric_version("1.2")), error = function(e) conditionMessage(e))
#---
names(numeric_version(c(a = "1.2", b = "3")))
#---
x <- numeric_version(c("1.2", "3.4.5")); x[[2]] <- "9.9"; x
#---
x <- numeric_version(c("1.2", "3.4.5")); x[1] <- "0.1.2"; x
#---
x <- numeric_version(c("1.2", "3.4.5")); x[3] <- "7"; x
#---
x <- numeric_version(c("1.2", "3.4.5")); is.na(x) <- 2; x
#---
x <- numeric_version(c("1.2", "3.4.5")); cat(tryCatch(x[[1]] <- "a", error = function(e) conditionMessage(e)), "\n")
#---
x <- numeric_version(c("1.2", "3.4.5")); cat(tryCatch(x[1] <- "a", error = function(e) conditionMessage(e)), "\n")
#---
x <- numeric_version(c("1.2", "3.4.5")); x[[c(1, 2)]]
#---
x <- numeric_version(c("1.2", "3.4.5")); x[[2, 1]]
#---
x <- numeric_version(c("1.2", "3.4.5")); x[[1, 5]]
#---
x <- numeric_version(c("1.2", "3.4.5")); tryCatch(x[[4]], error = function(e) conditionMessage(e))
#---
p <- package_version(c("3.4.5", "1.2")); p$major
#---
p <- package_version(c("3.4.5", "1.2")); p$minor
#---
p <- package_version(c("3.4.5", "1.2")); p$patchlevel
#---
p <- package_version(c("3.4.5", "1.2")); tryCatch(p$foo, error = function(e) conditionMessage(e))
#---
numeric_version("1.2")$major
#---
c(numeric_version("1.2"), package_version("1.2.3"))
#---
class(c(package_version("1.2.3"), package_version("1.2")))
#---
class(c(package_version("1.2.3"), numeric_version("1.2")))
#---
c(numeric_version("1.2"), "3")
#---
format(numeric_version(c("1.2", NA), strict = FALSE))
#---
as.character(numeric_version(c("1.2", NA), strict = FALSE))
#---
print(numeric_version(c("1.2", NA), strict = FALSE))
#---
print(package_version("1.2"), quote = FALSE)
#---
paste(numeric_version("1.2"))
#---
as.numeric_version(numeric_version("2.1"))
#---
as.numeric_version(3)
#---
cat(tryCatch(as.package_version("1"), error = function(e) conditionMessage(e)), "\n")
#---
as.package_version(numeric_version("1.2"))
#---
as.package_version(R.version)
#---
as.numeric_version(R.version)
#---
as.package_version(list(major = "1", minor = "2.3"))
#---
R_system_version("4.1.2")
#---
class(R_system_version("4.1.2"))
#---
is.package_version(R_system_version("4.1.2"))
#---
is.package_version(numeric_version("1.2"))
#---
anyNA(numeric_version(c("1.2", NA), strict = FALSE))
#---
unique(numeric_version(c("1.2", "1.2.0", "1.2")))
#---
as.list(package_version(c("1.2", "1.3")))
#---
sort(numeric_version(c("1.10", "1.9", "1.2")), decreasing = TRUE)
#---
order(numeric_version(c("1.10", "1.9", "1.2")))
#---
rev(numeric_version(c("1.10", "1.9", "1.2")))
#---
xtfrm(numeric_version(c("1.2.3", "10.0", "2", NA), strict = FALSE))
#---
sort(numeric_version(c("1.2.3", "10.0", "2", NA), strict = FALSE))
#---
sort(numeric_version(c("1.2.3", "10.0", "2", NA), strict = FALSE), na.last = TRUE)
#---
order(numeric_version(c("1.2.3", "10.0", "2", NA), strict = FALSE))
#---
numeric_version(c(a = "1.2", b = "3"))[2]
#---
x <- numeric_version(c("1.2", "3.4.5")); x[[1]][[1]]
#---
x <- numeric_version(c("1.2", "3.4.5")); x[3]
#---
x <- numeric_version(c("1.2", "3.4.5")); x[c(TRUE, FALSE)]
#---
x <- numeric_version(c("1.2", "3.4.5")); x[0]
#---
x <- numeric_version(c("1.2", "3.4.5")); x[NA]
#---
x <- numeric_version(c("1.2", "3.4.5")); x[["a"]]
#---
x <- numeric_version(c("1.2", "3.4.5")); x[[-1]]
#---
x <- numeric_version(c("1.2", "3.4.5")); rep(x, times = c(2, 1))
#---
x <- numeric_version(c("1.2", "3.4.5")); rep(x, each = 2)
#---
x <- numeric_version(c("1.2", "3.4.5")); rep(x, length.out = 3)
#---
rep(numeric_version("1.2"), 0)
#---
c(numeric_version("1"), NULL)
#---
c(numeric_version("1"), 2)
#---
tryCatch(c(numeric_version("1"), NA), error = function(e) conditionMessage(e))
#---
unique(numeric_version(c("1", "1.0")))
#---
duplicated(numeric_version(c("1", "1.0", "2")))
#---
unique(numeric_version(c("1", "1.0", "2")), fromLast = TRUE)
#---
is.na(numeric_version(c("1", NA), strict = FALSE))
#---
sapply(numeric_version(c("1.2", "3")), format)
#---
numeric_version(c("1.2", "3"))[[2, 1]]
#---
ifelse(numeric_version(c("1.2", "3")) > "2", "new", "old")
#---
identical(numeric_version("1.2"), numeric_version("1.2"))
#---
is.list(numeric_version("1.2"))
#---
inherits(package_version("1.2"), "numeric_version")
#---
as.list(numeric_version(c(a = "1.2")))
#---
as.character(package_version("1.2.3"))
#---
toString(numeric_version(c("1.2", "3")))
#---
paste("v", numeric_version(c("1.2", "3")))
#---
tryCatch(as.numeric(numeric_version("1.2")), error = function(e) conditionMessage(e))
#---
as.numeric_version(NULL)
#---
as.numeric_version(factor("1.2"))
#---
names(.standard_regexps())
#---
getRversion()
#---
class(getRversion())
#---
getRversion() >= "4.0.0"
#---
R.version$major
#---
R.version$language
#---
R.version.string
#---
class(R.version)
#---
names(R.version)
#---
identical(version, R.version)
#---
compareVersion("1.2.3", "1.2.4")
#---
compareVersion("1.10", "1.9")
#---
compareVersion("1.2", "1.2.0")
#---
compareVersion("2.0", "2.0")
#---
compareVersion("1.2-3", "1.2.3")
#---
compareVersion("1.2.0", "1.2")
#---
.make_numeric_version(c("1.2", "3"), TRUE, "[0-9.]+")
#---
.make_numeric_version("1.2", TRUE, "[0-9.]+", "package_version")
#---
`$.package_version`(package_version("3.4.5"), "minor")
#---
`[.numeric_version`(numeric_version(c("1.2", "3")), 2)
#---
`[[.numeric_version`(numeric_version(c("1.2", "3")), 1)
#---
`[<-.numeric_version`(numeric_version(c("1.2", "3")), 1, value = "9.9")
#---
`[[<-.numeric_version`(numeric_version(c("1.2", "3")), 2, value = "7.7")
#---
`is.na<-.numeric_version`(numeric_version(c("1.2", "3")), 1)
#---
anyNA.numeric_version(numeric_version(c("1.2", "3")))
#---
as.character.numeric_version(numeric_version(c("1.2", "3")))
#---
as.list.numeric_version(numeric_version(c("1.2", "3")))
#---
c.numeric_version(numeric_version("1.2"), numeric_version("3"))
#---
duplicated.numeric_version(numeric_version(c("1.2", "1.2")))
#---
format.numeric_version(numeric_version(c("1.2", "3")))
#---
is.na.numeric_version(numeric_version(c("1.2", "3")))
#---
print.numeric_version(numeric_version(c("1.2", "3")))
#---
rep.numeric_version(numeric_version("1.2"), 2)
#---
unique.numeric_version(numeric_version(c("1.2", "1.2")))
#---
xtfrm.numeric_version(numeric_version(c("1.2", "3")))
#---
tryCatch(Ops.numeric_version(numeric_version("1.2"), "1.3"), error = function(e) "failed")
#---
tryCatch(Summary.numeric_version(numeric_version("1.2")), error = function(e) "failed")
