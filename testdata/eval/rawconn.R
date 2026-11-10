## rawConnection, readBin, writeBin, readChar, writeChar, seek. A line holding only #--- separates cells.
#---
readBin(as.raw(1:8), "integer", n=2); readBin(as.raw(1:8), "integer", n=10); readBin(as.raw(1:8), "integer", n=0); readBin(raw(0), "integer"); readBin(as.raw(1:3), "integer"); readBin(as.raw(1:8), "int"); readBin(as.raw(1:4), 1L)
#---
readBin(as.raw(1:8), "integer", n=-1)
#---
readBin(as.raw(1:8), "integer", n=NA)
#---
readBin(as.raw(1:8), "integer", n=TRUE)
#---
readBin(as.raw(1:8), "integer", n=1.9); readBin(as.raw(1:8), "integer", n="2"); readBin(as.raw(1:8), "integer", n=c(1,2))
#---
readBin(as.raw(1:8), "integer", n=1, size=1); readBin(as.raw(1:8), "integer", n=8, size=1); readBin(as.raw(c(255,254)), "integer", n=2, size=1); readBin(as.raw(c(255,254)), "integer", n=2, size=1, signed=FALSE); readBin(as.raw(c(255,255)), "integer", n=1, size=2); readBin(as.raw(c(255,255)), "integer", n=1, size=2, signed=FALSE)
#---
readBin(as.raw(1:8), "integer", n=1, size=8); readBin(as.raw(c(255,255,255,255,255,255,255,255)), "integer", size=8); readBin(as.raw(c(1,2,3,4,5,6,7,8)), "integer", n=3, size=2); readBin(as.raw(c(1,2,3,4,5,6,7,8)), "integer", n=3, size=2, endian="big")
#---
readBin(as.raw(1:8), "integer", n=1, size=3)
#---
readBin(as.raw(1:4), "integer", endian="big"); readBin(as.raw(1:4), "integer", endian="little"); readBin(as.raw(1:4), "integer", endian="swap")
#---
readBin(as.raw(1:4), "integer", endian="x")
#---
readBin(as.raw(1:4), "integer", endian=c("big","little"))
#---
readBin(as.raw(c(0,0,0,0,0,0,240,63)), "double"); readBin(as.raw(c(0,0,128,63)), "double", size=4); readBin(as.raw(c(0,0,0,0,0,0,240,63)), "numeric"); readBin(as.raw(c(0,0,128,63)), "numeric", size=4)
#---
readBin(as.raw(c(0,0,0,0,0,0,248,127)), "double"); readBin(as.raw(c(0,0,192,127)), "double", size=4); readBin(as.raw(c(0,0,128,127)), "double", size=4); readBin(as.raw(c(0,0,128,255)), "double", size=4); readBin(as.raw(c(162,7,0,0,0,0,240,127)), "double")
#---
readBin(as.raw(c(1,0,0,0)), "logical"); readBin(as.raw(c(2,0,0,0)), "logical"); readBin(as.raw(c(0,0,0,128)), "logical"); readBin(as.raw(1:2), "logical", size=1, n=2)
#---
readBin(as.raw(1:16), "complex")
#---
readBin(as.raw(1:16), "complex", size=8)
#---
readBin(as.raw(c(97,98,0,99,0)), "character", n=2); readBin(as.raw(c(97,98,0,99)), "character", n=2); readBin(as.raw(c(97,98)), "character"); readBin(as.raw(c(0,0)), "character", n=2); readBin(raw(0), "character", 2)
#---
readBin(as.raw(1:4), "raw", n=2); readBin(as.raw(1:4), "raw", n=10); readBin(as.raw(1:4), raw(), n=2); readBin(as.raw(1:2), raw(0), 1); readBin(as.raw(1:4), "foo"); readBin(as.raw(1:4), "list"); readBin(as.raw(1:2), TRUE)
#---
readBin(as.raw(1:4), "raw", size=2)
#---
readBin(as.raw(1:2), list())
#---
readBin(1:3, "integer")
#---
readBin("no-such-file.bin", "integer")
#---
readBin(as.raw(1:4), "integer", signed=NA)
#---
readBin(as.raw(1:4), "integer", endian=NA)
#---
writeBin(1L, raw()); writeBin(c(1L, 256L), raw()); writeBin(1L, raw(), size=1); writeBin(1L, raw(), size=2); writeBin(1L, raw(), size=8); writeBin(1L, raw(), endian="big"); writeBin(1L, raw(), endian="swap"); writeBin(300L, raw(), size=1); writeBin(-1L, raw(), size=2); writeBin(1:2, raw(), size=8)
#---
writeBin(1.5, raw()); writeBin(1.5, raw(), size=4); writeBin(1.5, raw(), endian="big"); writeBin(1e10, raw(), size=4); writeBin(2^31, raw(), size=8)
#---
writeBin(TRUE, raw()); writeBin(NA, raw()); writeBin(NA_integer_, raw()); writeBin(NA_real_, raw()); writeBin(c(a=1L), raw())
#---
writeBin("ab", raw()); writeBin(c("a", "bc"), raw()); writeBin(NA_character_, raw()); writeBin("", raw()); writeBin("é", raw())
#---
writeBin(as.raw(1:3), raw()); writeBin(1i, raw()); writeBin(1+2i, raw(), endian="big"); writeBin(raw(0), raw()); writeBin(integer(0), raw())
#---
writeBin(1L, raw(), size=3)
#---
writeBin(1.5, raw(), size=2)
#---
writeBin(1L, raw(), size=0)
#---
writeBin(NULL, raw())
#---
writeBin(list(1), raw())
#---
writeBin(quote(a), raw())
#---
writeBin(factor("a"), raw())
#---
writeBin(1L, raw(), endian="x")
#---
x <- c(1.5, NA, Inf, -Inf, NaN, 0, -0); r <- writeBin(x, raw()); identical(readBin(r, "double", 7), x); readBin(r, "double", 7)
#---
x <- c(1L, NA, -1L, 2147483647L); r <- writeBin(x, raw()); identical(readBin(r, "integer", 4), x); r <- writeBin(x, raw(), endian="big"); identical(readBin(r, "integer", 4, endian="big"), x)
#---
x <- c(1+2i, NA, 3i); r <- writeBin(x, raw()); readBin(r, "complex", 3)
#---
x <- c("a", "bc", ""); r <- writeBin(x, raw()); readBin(r, "character", 3); readBin(r, "character", 5)
#---
x <- c(TRUE, FALSE, NA); r <- writeBin(x, raw()); readBin(r, "logical", 3); readBin(r, "logical", 3, size=1)
#---
r <- writeBin(c(0.1, 0.2), raw(), size=4); readBin(r, "double", 2, size=4)
#---
writeChar("abc", raw()); writeChar("abc", raw(), eos=NULL); writeChar("abc", raw(), nchars=2); writeChar("abc", raw(), eos="XY"); writeChar("é", raw()); writeChar("é", raw(), useBytes=TRUE); writeChar(character(0), raw()); writeChar(c("a","b"), raw(), nchars=1); writeChar("abc", raw(), nchars=2.9)
#---
writeChar(c("ab", "c"), raw(), nchars=c(1, 5))
#---
writeChar("abc", raw(), nchars=-1)
#---
writeChar(1, raw())
#---
writeChar("abc", raw(), nchars=c(1,2))
#---
writeChar(NA_character_, raw())
#---
writeChar("abc", raw(), eos=1)
#---
readChar(as.raw(c(97,98,99)), 2); readChar(as.raw(c(97,98,99)), c(1,1,5)); readChar(as.raw(c(97,98,99)), 0); readChar(as.raw(c(97,98,99)), 10); readChar(as.raw(c(0xc3,0xa9,97)), 1); readChar(as.raw(c(0xc3,0xa9,97)), 2, useBytes=TRUE); readChar(raw(0), 1); readChar(as.raw(97:99), numeric(0)); readChar(as.raw(97:99), 2.9); readChar(as.raw(97:99), TRUE)
#---
readChar(as.raw(97:100), c(2, 0, 2)); readChar(as.raw(97:100), c(2, 5, 2)); readChar(as.raw(97:100), c(5, 2)); nchar(readChar(as.raw(c(0xc3, 0xa9, 0xc3, 0xa9)), 3))
#---
readChar(as.raw(c(97,0,98)), 3)
#---
readChar(as.raw(c(97,98)), -1)
#---
readChar(as.raw(c(97,98)), NA)
#---
readChar(as.raw(97), 1e12)
#---
readChar(as.raw(97), "a")
#---
readChar(as.raw(c(0xff, 97)), 1)
#---
readChar(as.raw(97:99), 2, useBytes=NA)
#---
readChar(1, 1)
#---
readChar("no-such-file", 1)
#---
con <- rawConnection(as.raw(1:8)); k <- class(con); s <- summary(con); r <- list(k, s$description, s$class, s$mode, s$text, s$opened, s[["can read"]], s[["can write"]], isOpen(con)); close(con); r
#---
con <- rawConnection(as.raw(1:8)); a <- readBin(con, "integer"); b <- readBin(con, "integer"); d <- readBin(con, "integer"); close(con); list(a, b, d)
#---
con <- rawConnection(as.raw(1:8)); a <- readBin(con, "raw", 3); b <- readBin(con, "raw", 3); d <- readBin(con, "raw", 3); e <- readBin(con, "raw", 3); close(con); list(a, b, d, e)
#---
con <- rawConnection(as.raw(1:8), "rb"); a <- readBin(con, "raw", 3); close(con); a
#---
con <- rawConnection(raw(0), "w"); s <- summary(con); writeBin(1:2, con); writeBin("hi", con); a <- rawConnectionValue(con); writeChar("abc", con, eos=NULL); b <- rawConnectionValue(con); close(con); list(s$mode, s$text, s[["can read"]], s[["can write"]], a, b)
#---
con <- rawConnection(raw(0), "wb"); writeBin(1L, con); a <- rawConnectionValue(con); close(con); a
#---
con <- rawConnection(as.raw(1:3), "a"); writeBin(as.raw(9), con); a <- rawConnectionValue(con); close(con); a
#---
con <- rawConnection(as.raw(1:3), "w"); a <- rawConnectionValue(con); close(con); a
#---
con <- rawConnection(as.raw(1:3)); writeBin(1L, con)
#---
con <- rawConnection(raw(0), "w"); readBin(con, "integer")
#---
con <- rawConnection(raw(0), "w"); readChar(con, 1)
#---
con <- rawConnection(as.raw(1:3)); rawConnectionValue(con)
#---
con <- rawConnection(as.raw(1:3)); close(con); readBin(con, "integer")
#---
con <- rawConnection(as.raw(1:3)); close(con); isOpen(con)
#---
rawConnectionValue(1)
#---
rawConnection(1:3)
#---
rawConnection(as.raw(1), open=1)
#---
con <- rawConnection(charToRaw("a\nb\nc\n")); a <- readLines(con, 1); b <- readLines(con, 1); d <- readLines(con); close(con); list(a, b, d)
#---
con <- rawConnection(raw(0), "w"); writeLines(c("a", "b"), con); a <- rawConnectionValue(con); cat("x\n", file=con); b <- rawToChar(rawConnectionValue(con)); close(con); list(a, b)
#---
con <- rawConnection(raw(0), "w"); sink(con); print(1); sink(); r <- rawToChar(rawConnectionValue(con)); close(con); r
#---
con <- rawConnection(as.raw(1:8)); a <- seek(con); b <- readBin(con, "raw", 2); d <- seek(con); e <- seek(con, 0); f <- seek(con, 6, rw="read"); g <- readBin(con, "raw", 5); close(con); list(a, b, d, e, f, g)
#---
con <- rawConnection(as.raw(1:8)); r <- list(isIncomplete(con), isSeekable(con), flush(con)); close(con); r
#---
con <- rawConnection(as.raw(1:8)); a <- readBin(con, "integer", 1, endian="big"); b <- readBin(con, "integer", 1, size=2); close(con); list(a, b)
#---
con <- rawConnection(as.raw(c(1,0,0,0,2,0,0,0,3,0,0,0))); a <- readBin(con, "integer", 2); b <- readBin(con, "integer", 5); close(con); list(a, b)
#---
con <- rawConnection(as.raw(1:8)); a <- readChar(con, 3); b <- readChar(con, 3); d <- readChar(con, 3); close(con); list(a, b, d)
#---
con <- rawConnection(as.raw(c(97, 98, 0, 99, 0, 100))); a <- readBin(con, "character", 2); b <- readBin(con, "character", 2); close(con); list(a, b)
#---
con <- rawConnection(raw(0), "w"); writeBin(1:3, con, size=2); writeBin(c(1.5, 2), con, size=4); writeBin(c(TRUE, NA), con); writeBin(1i, con); writeBin(c("a", "b"), con); r <- rawConnectionValue(con); close(con); r
#---
f <- tempfile(); con <- file(f, "wb"); writeBin(1:3, con); writeBin(2.5, con); writeBin("s", con); writeChar("cd", con); close(con); a <- file.size(f); con <- file(f, "rb"); b <- readBin(con, "integer", 3); d <- readBin(con, "double"); e <- readBin(con, "character"); g <- readChar(con, 2); h <- readBin(con, "integer"); close(con); unlink(f); list(a, b, d, e, g, h)
#---
f <- tempfile(); writeBin(1:3, f); a <- file.exists(f); b <- readBin(f, "integer", 10); d <- file.size(f); e <- readBin(f, "raw", 100); unlink(f); list(a, b, d, e)
#---
f <- tempfile(); writeBin(as.raw(1:3), f); a <- readBin(f, "raw", n=2); b <- readChar(f, c(1, 2)); unlink(f); list(a, b)
#---
f <- tempfile(); con <- file(f, "w"); writeBin(1L, con)
#---
f <- tempfile(); con <- file(f, "wb"); close(con); con <- file(f, "r"); readBin(con, "integer")
#---
f <- tempfile(); con <- file(f, "wb"); writeBin(1:2, con); close(con); con <- file(f, "ab"); writeBin(3L, con); close(con); r <- readBin(f, "integer", 5); unlink(f); r
#---
f <- tempfile(); con <- file(f, "wb"); writeBin(1:2, con); close(con); con <- file(f, "rb"); a <- seek(con); b <- readBin(con, "integer"); d <- seek(con); close(con); unlink(f); list(a, b, d)
#---
con <- textConnection("abc"); readBin(con, "raw", 1)
#---
con <- textConnection("abc", "w"); writeBin(1L, con)
#---
con <- stdin(); readBin(con, "raw", 1)
#---
f <- tempfile(); con <- file(f, "wb"); writeBin(1L, con); close(con); con <- file(f, "rb"); r <- readChar(con, 2); close(con); unlink(f); r
