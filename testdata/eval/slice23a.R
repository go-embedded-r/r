print(signif(lchoose(5, 2), 12))
#---
print(signif(lchoose(50, 25), 12))
#---
print(signif(lchoose(c(10, 100, 1000), 3), 12))
#---
print(lchoose(5, 0)); print(lchoose(5, -1)); print(lchoose(5, 6)); print(lchoose(7, 7))
#---
print(signif(lchoose(7, 1), 12)); print(signif(lchoose(-7, 1), 12))
#---
print(signif(lchoose(-3, 4), 12))
#---
print(signif(lchoose(5.5, 3), 12)); print(signif(lchoose(0.5, 4), 12))
#---
print(signif(lchoose(1000, 500), 12)); print(signif(lchoose(100000, 5), 12))
#---
print(lchoose(NA, 2)); print(lchoose(3, NA)); print(lchoose(NaN, 1))
#---
print(lchoose(Inf, 3)); print(lchoose(numeric(0), 2))
#---
x <- lchoose(c(a = 5, b = 6), 2); print(signif(x, 12))
#---
print(tryCatch(lchoose(5, 2.5), warning = conditionMessage))
#---
print(tryCatch(lchoose("a", 2), error = conditionMessage))
#---
print(abbreviate("International", 5))
#---
st <- c("Alabama", "Alaska", "Arizona", "Arkansas", "California", "Colorado", "Connecticut", "Delaware", "Florida", "Georgia", "Hawaii", "Idaho", "Illinois", "Indiana", "Iowa", "Kansas", "Kentucky", "Louisiana", "Maine", "Maryland", "Massachusetts", "Michigan", "Minnesota", "Mississippi", "Missouri", "Montana", "Nebraska", "Nevada", "New Hampshire", "New Jersey", "New Mexico", "New York", "North Carolina", "North Dakota", "West Virginia")
print(abbreviate(st[1:12], 3))
#---
st <- c("Illinois", "Indiana", "Iowa", "Kansas", "Kentucky", "Louisiana", "Maine", "Maryland", "Massachusetts", "Michigan", "Minnesota", "Mississippi", "Missouri", "Montana", "Nebraska", "Nevada", "New Hampshire", "New Jersey", "New Mexico", "New York", "North Carolina", "North Dakota", "West Virginia")
print(abbreviate(st, 5))
print(abbreviate(st, 4))
#---
print(abbreviate("International Business Machines", 8, dot = TRUE))
#---
print(abbreviate(c("Sepal.Length", "Sepal.Width", "Petal.Length"), 6))
#---
print(abbreviate(c("abcd", "abce"), 3, strict = TRUE))
#---
print(abbreviate(c("abcd", "abce"), 3))
#---
print(abbreviate(c("abcdef", "abcdef"), 3))
#---
print(abbreviate(c("abc", NA), 2))
#---
print(abbreviate(c("abc", "d"), 0))
#---
print(abbreviate(c("abc", "d"), 0, named = FALSE))
#---
print(abbreviate(character(0)))
#---
print(abbreviate(c(a = "Alabama", b = "Alaska"), 5))
#---
print(abbreviate(factor(c("hello", "world")), 3))
#---
print(abbreviate("a b c d e f", 3))
#---
print(abbreviate("Hello World Again", 6, use.classes = FALSE))
#---
print(abbreviate(c("  padded  ", "x y"), 4))
#---
print(abbreviate(12345678, 4))
#---
print(tryCatch(abbreviate("abc", NA), error = conditionMessage))
#---
print(tryCatch(abbreviate("abc", 1:2), error = conditionMessage))
#---
print(grepl("should be one of", tryCatch(abbreviate("abc", 3, method = "zz"), error = conditionMessage)))
#---
print(tryCatch(abbreviate("abc", 3, dot = NA), error = conditionMessage))
#---
print(abbreviate("abcdefghij", 2.7))
#---
print(abbreviate("café au lait", 5))
#---
print(agrep("lasy", c("1 lazy", "1", "1 LAZY", "lasy", "laser", "la", "", NA)))
#---
print(agrep("lasy", c("1 lazy", "1", "1 LAZY", "lasy", "laser", "la", "", NA), max.distance = 2))
#---
print(agrep("lasy", c("1 lazy", "1", "1 LAZY"), max.distance = 2, value = TRUE))
#---
print(agrep("lasy", c("1 lazy", "1", "1 LAZY"), max.distance = 2, ignore.case = TRUE))
#---
print(agrep("lasy", c(a = "lazy", b = "laser"), value = TRUE))
#---
print(agrepl("LASY", c("lazy", "laser", NA, "abc"), ignore.case = TRUE))
#---
print(agrepl("lasy", c("lazy", "laser", NA, "abc"), max.distance = list(sub = 0, all = 1)))
#---
print(agrepl("lasy", c("lazy", "laser", NA, "abc"), max.distance = list(ins = 0, del = 0, sub = 1)))
#---
y <- c("abde", "abXde", "abcXde", "xabcde", "abcdeXX", "bcde")
print(agrepl("abcde", y, max.distance = list(all = 1, ins = 0)))
print(agrepl("abcde", y, max.distance = list(all = 1, del = 0)))
print(agrepl("abcde", y, max.distance = list(all = 1, sub = 0)))
#---
print(agrepl("abcde", c("ab", "abc", "abXXe"), max.distance = 2))
print(agrepl("abcde", c("ab", "abc", "abXXe"), max.distance = 0))
#---
print(agrepl("ab", c("a", "b", "xx"), max.distance = 0.5))
#---
print(agrepl("abc", c("xbc", "abc"), max.distance = list(cost = 1)))
#---
print(agrep(NA, c("a", "b")))
#---
print(agrep(NA, c("a", "b"), value = TRUE))
#---
print(tryCatch(agrep("", "a"), error = conditionMessage))
#---
print(tryCatch(agrep(character(0), "a"), error = conditionMessage))
#---
print(tryCatch(agrep(c("a", "b"), "a"), warning = conditionMessage))
#---
print(agrep("lasy", character(0)))
#---
print(agrep("lasy", factor(c("lazy", "x"))))
#---
print(agrepl("lasy", list("lazy")))
#---
print(tryCatch(agrepl("lasy", "lazy", value = TRUE), error = conditionMessage))
#---
print(tryCatch(agrep("a", "b", max.distance = -1), error = conditionMessage))
#---
print(tryCatch(agrep("a", "b", max.distance = list(zz = 1)), warning = conditionMessage))
#---
print(tryCatch(agrep("a", "b", costs = -1), error = conditionMessage))
#---
print(agrep("lasy", c("lazy", "laser", "abc"), max.distance = list(ins = 1)))
#---
print(OlsonNames()[1:3]); print(length(OlsonNames()) > 400)
#---
z <- OlsonNames(); print(c("Europe/Paris", "UTC", "Etc/GMT+5", "US/Pacific", "GMT") %in% z)
print(is.character(z)); print(!anyDuplicated(z))
#---
print(names(attributes(OlsonNames())))
#---
print(Sys.getlocale("LC_COLLATE")); print(Sys.getlocale("LC_CTYPE"))
#---
print(tryCatch(Sys.getlocale("xx"), error = conditionMessage))
#---
print(tryCatch(Sys.getlocale(c("LC_CTYPE", "LC_TIME")), error = conditionMessage))
#---
old <- Sys.getlocale("LC_COLLATE")
print(Sys.setlocale("LC_COLLATE", "C")); print(Sys.getlocale("LC_COLLATE"))
print(Sys.setlocale("LC_COLLATE", old)); print(Sys.getlocale("LC_COLLATE"))
#---
old <- Sys.getlocale("LC_COLLATE"); print(withVisible(Sys.setlocale("LC_COLLATE", old))$visible)
#---
print(tryCatch(Sys.setlocale("LC_COLLATE", "xx_YY"), warning = conditionMessage))
#---
print(suppressWarnings(Sys.setlocale("LC_COLLATE", "xx_YY")))
#---
print(tryCatch(Sys.setlocale("LC_COLLATE", NA), error = conditionMessage))
#---
print(tryCatch(Sys.setlocale("LC_COLLATE", c("C", "C")), error = conditionMessage))
#---
print(tryCatch(Sys.setlocale("zz", "C"), error = conditionMessage))
#---
print(Sys.setlocale("LC_TIME", "C")); print(Sys.getlocale("LC_TIME"))
#---
print(names(Sys.localeconv())[1:6]); print(Sys.localeconv()[["decimal_point"]]); print(length(Sys.localeconv()))
#---
print(Sys.which("zz-nonexistent-prog"))
#---
print(Sys.which(c("zz-a", "zz-b")))
#---
print(Sys.which(character(0)))
#---
print(Sys.which(NA))
#---
print(names(Sys.which(1)))
#---
print(findRestart("abort"))
#---
print(length(computeRestarts()))
#---
print(computeRestarts()[[1]][[1]]); print(class(computeRestarts()[[1]]))
#---
withRestarts(print(sapply(computeRestarts(), function(r) r[[1]])), mine = function() 1)
#---
print(unclass(findRestart("abort")))
#---
print(tryCatch(invokeRestartInteractively("abort"), error = conditionMessage))
#---
print(tryCatch(withRestarts(invokeRestartInteractively("zz"), zz = function() 42), error = conditionMessage))
#---
print(.amatch_bounds(2)); print(.amatch_bounds(list(cost = 1, sub = 2)))
#---
print(.amatch_costs(NULL)); print(.amatch_costs(c(ins = 2, sub = 3)))
#---
print.restart(findRestart("abort"))
#---
w <- structure(list(quote(f(x)), quote(f(x))), names = c("a", "a"), class = "warnings")
print.summary.warnings(summary.warnings(w))
#---
x <- as.POSIXlt("2020-03-01 10:00:00", tz = "UTC"); b <- balancePOSIXlt(x); print(attr(b, "balanced"))
#---
f <- function() closeAllConnections(); print(is.function(f))
#---
x <- as.POSIXlt(c("2020-03-01 10:00:00", "2021-07-04 12:30:15"), tz = "UTC"); b <- balancePOSIXlt(x)
print(class(b)); print(names(attributes(b))); print(b$year); print(attr(b, "balanced"))
#---
x <- as.POSIXlt(c("2020-03-01 10:00:00", "2021-07-04 12:30:15"), tz = "UTC"); b <- balancePOSIXlt(x, classed = FALSE)
print(class(b)); print(names(attributes(b))); print(b$mon)
#---
x <- as.POSIXlt("2020-03-01 10:00:00", tz = "UTC"); b <- balancePOSIXlt(x, fill.only = TRUE)
print(attr(b, "balanced")); print(class(b)); print(length(b$sec))
#---
x <- as.POSIXlt("2020-03-01 10:00:00", tz = "UTC"); x$min <- c(5L, 6L, 7L); b <- balancePOSIXlt(x, fill.only = TRUE)
print(lengths(unclass(b))[1:6])
