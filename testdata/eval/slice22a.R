print(numToBits(1))
#---
print(numToInts(c(1, -2.5)))
#---
print(numToBits(c(NA, Inf))[1:8])
#---
print(numToBits(c(0, -0, 5e-324)))
#---
print(charmatch("", ""))
#---
print(charmatch("", c("a","b")))
#---
print(charmatch(NA, c("NA","b")))
#---
print(charmatch(c("m","me","mea","x"), c("mean","median","mode")))
#---
print(charmatch("me", c("me","mean","me")))
#---
print(charmatch(character(), "a"))
#---
print(charmatch("a", character()))
#---
print(charmatch("a", "b", nomatch=0))
#---
print(charmatch(1, "b"))
#---
print(charmatch("a", 1))
#---
print(charmatch(factor("m"), "mean"))
#---
print(char.expand("me", c("mean","median")))
#---
print(char.expand("mea", c("mean","median")))
#---
e <- tryCatch(char.expand("x", c("mean","median")), error = function(e) e)
print(conditionCall(e)); print(conditionMessage(e))
#---
print(char.expand("x", c("mean","median"), nomatch=NULL))
#---
print(tryCatch(char.expand(c("x","y"), "mean"), error = conditionMessage))
#---
print(tryCatch(char.expand(1, "mean"), error = conditionMessage))
#---
print(tryCatch(browserText(), error = function(e) c(conditionMessage(e), deparse(conditionCall(e)))))
#---
print(tryCatch(browserCondition(), error = function(e) c(conditionMessage(e), deparse(conditionCall(e)))))
#---
print(tryCatch(browserSetDebug(), error = function(e) c(conditionMessage(e), deparse(conditionCall(e)))))
#---
print(validEnc(c("a", "\xff", NA, "é")))
#---
print(encodeString(c("a\nb", NA, "x"), quote='"'))
#---
print(encodeString(c("a","bbb"), width=NA))
#---
print(encodeString(c("a","bbb"), width=5, justify="right"))
#---
print(encodeString("a", width=5, justify="centre"))
#---
print(encodeString(NA_character_, na.encode=FALSE))
#---
print(encodeString(c(a="x", b="y\\z")))
#---
print(encodeString("it's", quote="'"))
#---
print(encodeString(1))
#---
print(format.info(c(1,10,100)))
#---
print(format.info(123.456))
#---
print(format.info(1e10))
#---
print(format.info(1e-300))
#---
print(format.info(c(1.5,NA)))
#---
print(format.info(1:100))
#---
print(format.info(c(TRUE,NA)))
#---
print(format.info("abc"))
#---
print(format.info(pi, digits=3))
#---
print(format.info(1, nsmall=2))
#---
print(format.info(c(-1.5, 22.25, NaN, Inf)))
#---
print(Sys.getpid() > 0)
#---
print(is.integer(Sys.getpid()))
#---
print(is.null(reg.finalizer(new.env(), function(e) 1)))
#---
print(tryCatch(reg.finalizer(1, function(e) 1), error = conditionMessage))
#---
print(findInterval(c(0, 1, 1.5, 3, 5), c(1, 2, 3)))
#---
print(findInterval(c(1, 3), c(1, 2, 3), rightmost.closed = TRUE))
#---
print(findInterval(c(0, 1, 4), c(1, 2, 3), all.inside = TRUE))
#---
print(findInterval(c(1, 2, 3.5), c(1, 2, 3), left.open = TRUE))
#---
print(findInterval(c(NA, 2), c(1, 2, 3)))
#---
print(try(findInterval(1, c(2, 1)), silent = TRUE))
#---
print(findInterval(1, numeric()))
#---
print(choose(5, 2))
#---
print(choose(50, 25))
#---
print(choose(c(4, 5, 6), 0:2))
#---
print(choose(-3, 2))
#---
print(choose(-3, 3))
#---
print(choose(0.5, 3))
#---
print(choose(5, 2.2))
#---
print(choose(5, 7))
#---
print(choose(5, -1))
#---
print(choose(c(a = 5, b = 6), 2))
#---
print(choose(NA, 2))
#---
print(signif(choose(60, 30), 12))
#---
print(try(choose("a", 2), silent = TRUE))
#---
print(choose(Inf, 2))
#---
print(signif(choose(2000, 1000), 12))
#---
print(.makeMessage(1:3, "a"))
#---
cat(.makeMessage("a", appendLF = TRUE), "|", sep = "")
#---
print(.mapply(function(x, y) x + y, list(1:3, 4:6), NULL))
#---
print(.mapply(function(a, b) a + b, list(a = 1:3, b = 1), NULL))
#---
print(.mapply(function(x, y) x + y, list(1:3, integer()), NULL))
#---
print(try(default.stringsAsFactors(), silent = TRUE))
#---
print(dontCheck(3))
#---
print(.valid.factor(factor(c("a", "b"))))
#---
f <- factor(c("a", "b")); levels(f) <- c("x", "x"); print(.valid.factor(f))
#---
r <- withCallingHandlers(.handleSimpleError(function(e) 99, "msg", quote(g())), error = function(e) 1)
print(r)
#---
print(as.expression.default(1:2))
#---
print(as.expression.default(quote(x)))
#---
print(as.expression.default("a"))
#---
x <- .POSIXlt(unclass(as.POSIXlt("2020-01-01", tz = "UTC")), "UTC"); print(x); print(attributes(x)$tzone)
#---
print(.S3method("print", "zzz", function(x, ...) cat("zzz\n")))
print(structure(1, class = "zzz"))
#---
x <- 1:3; comment(x) <- "hello"; print(comment(x)); print(x)
#---
x <- 1; r <- try(comment(x) <- 1, silent = TRUE); cat(conditionMessage(attr(r, "condition")), "\n")
#---
x <- 1; comment(x) <- "a"; comment(x) <- NULL; print(attributes(x))
#---
d <- structure(c(a = "x", b = "yy"), class = "Dlist")
print(unclass(d[1])); print(class(d[1])); print(class(d[2:1]))
#---
print(names(R.Version())[1:3])
#---
print(as.character(factor(c("b", "a", "b"))))
#---
print(as.character.default(structure(1:2, class = "foo")))
#---
f <- function() 1
g <- remove
print(identical(g, rm))
#---
print(base::kronecker(diag(2), matrix(1:4, 2)))
#---
print(base::kronecker(1:2, 1:3, "+"))
#---
print(encodeString(c("a", NA)))
#---
print(encodeString(c("a", NA), width=NA, justify="right"))
#---
print(findInterval(c(6, 1.5, NA, 4), 6, left.open=TRUE, all.inside=TRUE))
#---
print(findInterval(c(0, 6, 9), 6, all.inside=TRUE))
#---
print(findInterval(c(-1, 2.5, NA, 9), 0, all.inside=TRUE))
#---
print(findInterval(c(7, 3.5, NA, 0), 5, left.open=TRUE, all.inside=TRUE))
#---
print(findInterval(c(0, 3, 9), c(2, 4), all.inside=TRUE))
#---
d <- structure(c(a = "first", bb = "second"), class = "Dlist"); print(class(`[.Dlist`(d, 2)))
#---
d <- structure(c(a = "first", bb = "second"), class = "Dlist"); invisible(print.Dlist(d))
#---
f <- factor(c("u", "v", "u")); print(as.character.factor(f))
#---
zz <- 1; remove(zz); print(exists("zz"))
#---
e <- quote(base:::`[<-.POSIXct`(p, 2, value = v)); print(length(e))
#---
dd <- as.difftime(c(1, 2), units = "mins"); dd <- `[<-.difftime`(dd, 1, value = as.difftime(1, units = "hours")); print(dd)
#---
print(choose(numeric(), 2)); print(choose(3, integer()))
#---
print(tryCatch(format.info(1, digits = NA), error = conditionMessage))
#---
print(tryCatch(format.info(1, digits = 0), error = conditionMessage))
#---
print(tryCatch(format.info(1, digits = integer()), error = conditionMessage))
#---
print(choose(c(3, 4), c(5, 0)))
