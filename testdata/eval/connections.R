tc <- textConnection("a\nb"); r <- readLines(tc); close(tc); r
#---
tc <- textConnection(c("one", "two", "three")); x <- readLines(tc, 2); y <- readLines(tc, 1); z <- readLines(tc); close(tc); list(x, y, z)
#---
tc <- textConnection(c("one", "two", "three")); x <- readLines(tc, n = 0); close(tc); x
#---
tc <- textConnection(character(0)); x <- readLines(tc); close(tc); x
#---
tc <- textConnection(c("a", NA, "c")); x <- readLines(tc); close(tc); x
#---
tc <- textConnection("a\r\nb\rc\n\nd"); x <- readLines(tc); close(tc); x
#---
tc <- textConnection(""); x <- readLines(tc); close(tc); x
#---
tc <- textConnection("x\ny\nz"); a <- readLines(tc, 1); b <- readLines(tc, 5); c <- readLines(tc); close(tc); list(a, b, c)
#---
tc <- textConnection("a"); k <- class(tc); close(tc); k
#---
tc <- textConnection("a"); r <- close.connection(tc); r
#---
f <- tempfile(); con <- file(f); r <- open.connection(con, "w"); r2 <- isOpen(con); close(con); unlink(f); list(r, r2)
#---
f <- tempfile(); con <- file(f, "w"); r <- flush.connection(con); close(con); unlink(f); r
#---
tc <- textConnection("a"); r <- summary.connection(tc)$class; close(tc); r
#---
tc <- textConnection("abc"); r <- print.connection(tc); close(tc); class(r)
#---
tc <- textConnection("abc"); r <- print(tc); close(tc); class(r)
#---
tc <- textConnection("abc"); s <- summary(tc); close(tc); s
#---
tc <- textConnection("abc"); print(tc); close(tc)
#---
tc <- textConnection("abc", name = "mytext"); s <- summary(tc)$description; close(tc); s
#---
x <- c("p", "q"); tc <- textConnection(x); s <- summary(tc)$description; close(tc); s
#---
tc <- textConnection("abc"); r <- c(isOpen(tc), isOpen(tc, "read"), isOpen(tc, "write")); close(tc); r
#---
tc <- textConnection("abc"); close(tc); isOpen(tc)
#---
tc <- textConnection("abc"); close(tc); close(tc)
#---
tc <- textConnection("abc"); close(tc); readLines(tc)
#---
tc <- textConnection("abc"); r <- close(tc); r
#---
tc <- textConnection("abc"); r <- isIncomplete(tc); close(tc); r
#---
tc <- textConnection("abc"); tryCatch(writeLines("x", tc), finally = close(tc))
#---
textConnection(1:3)
#---
tc <- textConnection(NULL, "w"); writeLines(c("a", "b"), tc); r <- textConnectionValue(tc); close(tc); r
#---
tc <- textConnection(NULL, "w"); writeLines("a", tc); cat("par", file = tc); r <- list(textConnectionValue(tc), isIncomplete(tc)); close(tc); r
#---
tc <- textConnection("out", "w"); writeLines("hi", tc); a <- out; cat("x", file = tc); b <- out; close(tc); list(a, b, out)
#---
tc <- textConnection("out", "w"); r <- out; close(tc); r
#---
assign("app_out", c("old1", "old2"), envir = globalenv()); tc <- textConnection("app_out", "a"); writeLines("new", tc); close(tc); app_out
#---
assign("app_out2", c("old1", "old2"), envir = globalenv()); tc <- textConnection("app_out2", "w"); writeLines("new", tc); close(tc); app_out2
#---
f <- function() { tc <- textConnection("loc", "w", local = TRUE); writeLines("inside", tc); close(tc); loc }; list(f(), exists("loc"))
#---
tc <- textConnection("out2", "w"); writeLines(c("a", "b", "c"), tc, sep = "|"); close(tc); out2
#---
tc <- textConnection("out3", "w"); cat("a\nb\nc", file = tc); close(tc); out3
#---
tc <- textConnection("out4", "w"); cat("x", "y", sep = "\n", file = tc); close(tc); out4
#---
tc <- textConnection(NULL, "w"); r <- try(readLines(tc), silent = TRUE); close(tc); class(r)
#---
tc <- textConnection(NULL, "w"); tryCatch(readLines(tc), finally = close(tc))
#---
tc <- textConnection("a"); tryCatch(textConnectionValue(tc), finally = close(tc))
#---
tc <- textConnection(NULL, "w"); s <- summary(tc); close(tc); s[c("class", "mode", "opened", "can read", "can write")]
#---
tc <- textConnection("abc"); s <- unlist(summary(tc)[c("class", "mode", "text", "opened", "can read", "can write")]); close(tc); s
#---
f <- tempfile(); writeLines(c("a", "b"), f); r <- readLines(f); unlink(f); r
#---
f <- tempfile(); writeLines(c("a", "b"), f); writeLines("c", f); r <- readLines(f); unlink(f); r
#---
f <- tempfile(); writeLines(c("a", "b"), f, sep = "\r\n"); r <- readLines(f); unlink(f); r
#---
f <- tempfile(); writeLines(c("a", NA), f); r <- readLines(f); unlink(f); r
#---
f <- tempfile(); writeLines(character(0), f); r <- list(file.exists(f), readLines(f)); unlink(f); r
#---
f <- tempfile(); cat("alpha\nbeta\ngamma\n", file = f); r <- readLines(f, n = 2); unlink(f); r
#---
f <- tempfile(); cat("alpha\n", file = f); cat("beta\n", file = f, append = TRUE); r <- readLines(f); unlink(f); r
#---
f <- tempfile(); cat("alpha\n", file = f); cat("beta\n", file = f); r <- readLines(f); unlink(f); r
#---
f <- tempfile(); cat("x", "y", "z", file = f, sep = "\n"); r <- readLines(f); unlink(f); r
#---
f <- tempfile(); cat("x", "y", file = f, sep = "\n"); r <- file.size(f); unlink(f); r
#---
f <- tempfile(); cat("x", "y", file = f, sep = ", "); r <- file.size(f); unlink(f); r
#---
f <- tempfile(); cat("x", "y", file = f, sep = " \n"); r <- file.size(f); unlink(f); r
#---
f <- tempfile(); cat("x", "y", "z", file = f, sep = c(",", "\n")); r <- file.size(f); unlink(f); r
#---
f <- tempfile(); cat(character(0), file = f, sep = "\n"); r <- file.size(f); unlink(f); r
#---
f <- tempfile(); cat("x\n", file = f, sep = "\n"); r <- file.size(f); unlink(f); r
#---
f <- tempfile(); cat(1:3, 4.5, TRUE, "s", file = f); r <- readLines(f, warn = FALSE); unlink(f); r
#---
cat("a", "b", sep = "\n")
#---
cat("a", "b", "c", sep = c(",", "\n"))
#---
cat("a\n", sep = "\n")
#---
cat(character(0), sep = "\n")
#---
cat("a", "b", sep = ", ")
#---
cat("a", NULL, "b", sep = "-")
#---
cat(NULL, "b", sep = "-")
#---
cat("a", character(0), "b", sep = "-")
#---
cat(character(0), "", sep = "-")
#---
cat("a", "b", NULL, NULL, sep = "-")
#---
cat("a", "b", character(0), sep = "-")
#---
cat(c("a", "b"), "c", sep = c("-", "+"))
#---
cat(character(0), "")
#---
cat(NULL, 2.5, sep = "\n\n")
#---
cat("a", "b", sep = " \n")
#---
tc <- textConnection("append_nonexist_zz", "a"); close(tc); append_nonexist_zz
#---
list(1, structure(list(a = 1), class = "kk"))
#---
f <- tempfile(); cat("123\nabc", file = f); r <- tryCatch(readLines(f), warning = function(w) sub(f, "FILE", conditionMessage(w), fixed = TRUE)); unlink(f); r
#---
f <- tempfile(); cat("123\nabc", file = f); r <- suppressWarnings(readLines(f)); unlink(f); r
#---
f <- tempfile(); cat("123\nabc", file = f); r <- readLines(f, warn = FALSE); unlink(f); r
#---
f <- tempfile(); cat("123\nabc", file = f); w <- NULL; r <- withCallingHandlers(readLines(f), warning = function(c) { w <<- conditionCall(c); invokeRestart("muffleWarning") }); unlink(f); w
#---
f <- tempfile(); cat("123\nabc\n", file = f); r <- tryCatch(readLines(f, n = 5, ok = FALSE), error = function(e) conditionMessage(e)); unlink(f); r
#---
f <- tempfile(); cat("123\nabc\n", file = f); r <- readLines(f, n = 5); unlink(f); r
#---
f <- tempfile(); cat("123\nabc\n", file = f); readLines(f, n = -1L)
#---
f <- tempfile(); cat("a\nb\n", file = f); con <- file(f); r <- list(isOpen(con), readLines(con, 1), isOpen(con), readLines(con, 1)); close(con); unlink(f); r
#---
f <- tempfile(); cat("a\nb\nc\n", file = f); con <- file(f, "r"); r <- list(isOpen(con), readLines(con, 1), readLines(con, 1), readLines(con), readLines(con)); close(con); unlink(f); r
#---
f <- tempfile(); cat("a\nb\n", file = f); con <- file(f); open(con, "r"); r <- list(isOpen(con), readLines(con, 1)); close(con); unlink(f); r
#---
f <- tempfile(); cat("a\nb\n", file = f); con <- file(f, "r"); r <- open(con); close(con); r
#---
f <- tempfile(); cat("a\nb\n", file = f); con <- file(f, "r"); r <- tryCatch(open(con), warning = function(w) conditionMessage(w)); close(con); r
#---
f <- tempfile(); cat("a\nb\n", file = f); con <- file(f); r <- open(con); r2 <- isOpen(con); close(con); list(r, r2)
#---
f <- tempfile(); con <- file(f, "w"); writeLines("one", con); writeLines("two", con); close(con); r <- readLines(f); unlink(f); r
#---
f <- tempfile(); con <- file(f, "w"); cat("a", file = con); cat("b\n", file = con); close(con); r <- readLines(f); unlink(f); r
#---
f <- tempfile(); writeLines("one", f); con <- file(f, "a"); writeLines("two", con); close(con); r <- readLines(f); unlink(f); r
#---
f <- tempfile(); writeLines("one", f); con <- file(f); writeLines("two", con); r <- readLines(f); close(con); unlink(f); r
#---
f <- tempfile(); writeLines("one", f); con <- file(f, "w"); close(con); r <- readLines(f); unlink(f); r
#---
f <- tempfile(); con <- file(f, "w"); r <- try(readLines(con), silent = TRUE); close(con); unlink(f); class(r)
#---
f <- tempfile(); writeLines("one", f); con <- file(f, "r"); r <- tryCatch(writeLines("x", con), error = function(e) conditionMessage(e)); close(con); unlink(f); r
#---
f <- tempfile(); con <- file(f); s <- summary(con); close(con); s[c("class", "mode", "text", "opened", "can read", "can write")]
#---
f <- tempfile(); writeLines("x", f); con <- file(f, "r"); s <- summary(con); close(con); unlink(f); s[c("class", "mode", "text", "opened", "can read", "can write")]
#---
f <- tempfile(); con <- file(f, "w"); s <- summary(con); close(con); unlink(f); s[c("class", "mode", "text", "opened", "can read", "can write")]
#---
f <- tempfile(); con <- file(f); r <- identical(summary(con)$description, f); close(con); r
#---
f <- tempfile(); con <- file(f); class(con) -> k; close(con); k
#---
readLines("/nonexistent/dir/file.txt")
#---
tryCatch(readLines("/nonexistent/dir/file.txt"), warning = function(w) conditionMessage(w))
#---
w <- NULL; withCallingHandlers(try(readLines("/nonexistent/dir/file.txt"), silent = TRUE), warning = function(c) { w <<- deparse(conditionCall(c)); invokeRestart("muffleWarning") }); w
#---
file("/nonexistent/dir/file.txt", "r")
#---
tryCatch(file("/nonexistent/dir/file.txt", "r"), warning = function(w) conditionMessage(w))
#---
file("/nonexistent/dir/file.txt", "w")
#---
writeLines("x", "/nonexistent/dir/file.txt")
#---
cat("x", file = "/nonexistent/dir/file.txt")
#---
con <- file("/nonexistent/dir/file.txt"); tryCatch(readLines(con), finally = close(con))
#---
con <- file("/nonexistent/dir/file.txt"); w <- NULL; r <- withCallingHandlers(try(readLines(con), silent = TRUE), warning = function(c) { w <<- deparse(conditionCall(c)); invokeRestart("muffleWarning") }); close(con); w
#---
tryCatch(file("/nonexistent/dir/file.txt", "q"), error = function(e) conditionMessage(e))
#---
f <- tempfile(); tryCatch(file(f, "q"), warning = function(w) sub(f, "FILE", conditionMessage(w), fixed = TRUE))
#---
f <- tempfile(); r <- suppressWarnings(try(file(f, "q"), silent = TRUE)); class(r)
#---
writeLines(1:3)
#---
writeLines("a", con = 5)
#---
readLines(5)
#---
close(5)
#---
writeLines("hello")
#---
writeLines(c("a", "b"), sep = "|")
#---
writeLines(c("a", NA, "c"))
#---
writeLines("a", stdout())
#---
writeLines("to stderr", stderr())
#---
cat("to stderr\n", file = stderr())
#---
cat("to stdout\n", file = stdout())
#---
cat("a", "b", file = "")
#---
r <- readLines(stdin(), 0); r
#---
tc <- textConnection("x\ny"); r <- readLines(tc, 1); close(tc); r
#---
readline("prompt> ")
#---
x <- readline(); x
#---
nullfile()
#---
isatty(stdout())
#---
flush(stdout())
#---
f <- tempfile(); con <- file(f, "w"); r <- flush(con); close(con); r
#---
f <- tempfile(); r <- c(file.exists(f), file.create(f), file.exists(f), file.remove(f), file.exists(f)); r
#---
f <- tempfile(); writeLines("x", f); r <- file.remove(f); r2 <- suppressWarnings(file.remove(f)); list(r, r2)
#---
tryCatch(file.remove(tempfile()), warning = function(w) sub("'.*', reason", "'FILE', reason", conditionMessage(w)))
#---
f <- tempfile(); writeLines("x", f); r <- unlink(f); list(r, file.exists(f))
#---
unlink(tempfile())
#---
unlink("/nonexistent/dir/x")
#---
d <- tempfile(); dir.create(d); r <- list(unlink(d), dir.exists(d), unlink(d, recursive = TRUE), dir.exists(d)); r
#---
d <- tempfile(); dir.create(d); writeLines("x", file.path(d, "a.txt")); unlink(d, recursive = TRUE); file.exists(file.path(d, "a.txt"))
#---
d <- tempfile(); dir.create(d); writeLines("x", file.path(d, "a.txt")); writeLines("y", file.path(d, "b.txt")); unlink(file.path(d, "*.txt")); list.files(d)
#---
d <- tempfile(); dir.create(d); for (n in c("a1.txt", "a2.txt", "b1.csv")) writeLines("x", file.path(d, n)); unlink(file.path(d, "a?.txt")); list.files(d)
#---
f <- tempfile(); file.create(f)
#---
f <- tempfile(); r <- file.create(f); unlink(f); r
#---
file.create("/nonexistent/dir/x")
#---
suppressWarnings(file.create("/nonexistent/dir/x"))
#---
file.create("/nonexistent/dir/x", showWarnings = FALSE)
#---
f1 <- tempfile(); f2 <- tempfile(); r <- file.create(f1, f2); unlink(c(f1, f2)); r
#---
file.exists(c("/", "/nonexistent", NA))
#---
f <- tempfile(); g <- tempfile(); writeLines("x", f); r <- list(file.exists(f, g), file.exists(c(f, g))); unlink(f); r
#---
file.exists(1)
#---
file.exists()
#---
a <- tempfile(); b <- tempfile(); writeLines("data", a); r <- list(file.copy(a, b), readLines(b), file.copy(a, b), file.copy(a, b, overwrite = TRUE)); unlink(c(a, b)); r
#---
a <- tempfile(); d <- tempfile(); writeLines("data", a); dir.create(d); r <- list(file.copy(a, d), readLines(file.path(d, basename(a)))); unlink(c(a, d), recursive = TRUE); r
#---
file.copy("/nonexistent/x", tempfile())
#---
a <- tempfile(); b <- tempfile(); writeLines("data", a); r <- list(file.rename(a, b), file.exists(a), file.exists(b), readLines(b)); unlink(b); r
#---
tryCatch(file.rename(tempfile(), tempfile()), warning = function(w) sub("'[^']*' to '[^']*'", "'A' to 'B'", conditionMessage(w)))
#---
suppressWarnings(file.rename(tempfile(), tempfile()))
#---
a <- tempfile(); b <- tempfile(); writeLines("one", a); writeLines("two", b); r <- list(file.append(a, b), readLines(a)); unlink(c(a, b)); r
#---
a <- tempfile(); writeLines("one", a); file.append(a, tempfile())
#---
f <- tempfile(); writeLines(c("abc", "de"), f); r <- file.size(f); unlink(f); r
#---
file.size(tempfile())
#---
f <- tempfile(); writeLines("x", f); r <- file.size(f, tempfile()); unlink(f); r
#---
d <- tempfile(); dir.create(d)
#---
d <- tempfile(); r <- list(dir.create(d), dir.exists(d), suppressWarnings(dir.create(d))); unlink(d, recursive = TRUE); r
#---
d <- tempfile(); dir.create(d); tryCatch(dir.create(d), warning = function(w) sub(d, "DIR", conditionMessage(w), fixed = TRUE))
#---
d <- tempfile(); dir.create(d); dir.create(d, showWarnings = FALSE)
#---
d <- file.path(tempfile(), "a", "b"); r <- list(suppressWarnings(dir.create(d)), dir.create(d, recursive = TRUE), dir.exists(d)); unlink(dirname(dirname(d)), recursive = TRUE); r
#---
tryCatch(dir.create("/nonexistent/dir/sub"), warning = function(w) conditionMessage(w))
#---
dir.exists(c("/", "/nonexistent", tempdir()))
#---
dir.exists(1)
#---
d <- tempfile(); dir.create(d); for (n in c("b.txt", "a.txt", "C.txt", "c.csv", ".hidden")) writeLines("x", file.path(d, n)); list.files(d)
#---
d <- tempfile(); dir.create(d); for (n in c("b.txt", "a.txt", "C.txt", "c.csv", ".hidden")) writeLines("x", file.path(d, n)); list.files(d, all.files = TRUE)
#---
d <- tempfile(); dir.create(d); for (n in c("b.txt", "a.txt", "C.txt", "c.csv", ".hidden")) writeLines("x", file.path(d, n)); list.files(d, all.files = TRUE, no.. = TRUE)
#---
d <- tempfile(); dir.create(d); for (n in c("b.txt", "a.txt", "C.txt", "c.csv")) writeLines("x", file.path(d, n)); list.files(d, pattern = "\\.txt$")
#---
d <- tempfile(); dir.create(d); for (n in c("b.txt", "a.txt", "C.TXT", "c.csv")) writeLines("x", file.path(d, n)); list.files(d, pattern = "txt$", ignore.case = TRUE)
#---
d <- tempfile(); dir.create(d); for (n in c("b.txt", "a.txt", "ab.txt")) writeLines("x", file.path(d, n)); list.files(d, pattern = "^a")
#---
d <- tempfile(); dir.create(d); for (n in c("a.txt", "a_txt")) writeLines("x", file.path(d, n)); list.files(d, pattern = ".txt", fixed = TRUE)
#---
d <- tempfile(); dir.create(d); for (n in c("b.txt", "a.txt")) writeLines("x", file.path(d, n)); r <- list.files(d, full.names = TRUE); identical(r, file.path(d, c("a.txt", "b.txt")))
#---
d <- tempfile(); dir.create(d); dir.create(file.path(d, "sub")); writeLines("x", file.path(d, "top.txt")); writeLines("y", file.path(d, "sub", "deep.txt")); list.files(d)
#---
d <- tempfile(); dir.create(d); dir.create(file.path(d, "sub")); writeLines("x", file.path(d, "top.txt")); writeLines("y", file.path(d, "sub", "deep.txt")); list.files(d, recursive = TRUE)
#---
d <- tempfile(); dir.create(d); dir.create(file.path(d, "sub")); writeLines("x", file.path(d, "top.txt")); writeLines("y", file.path(d, "sub", "deep.txt")); list.files(d, recursive = TRUE, include.dirs = TRUE)
#---
d <- tempfile(); dir.create(d); dir.create(file.path(d, "sub")); writeLines("x", file.path(d, "sub", "deep.txt")); list.files(d, recursive = TRUE, pattern = "deep")
#---
d <- tempfile(); dir.create(d); dir.create(file.path(d, "sub")); writeLines("x", file.path(d, "sub", "deep.txt")); r <- list.files(d, recursive = TRUE, full.names = TRUE); identical(r, file.path(d, "sub", "deep.txt"))
#---
d <- tempfile(); dir.create(d); dir.create(file.path(d, "sub")); writeLines("x", file.path(d, "sub", ".h")); list.files(d, recursive = TRUE, all.files = TRUE)
#---
d <- tempfile(); dir.create(d); dir <- c(basename(d)); r <- dir(d); r
#---
d <- tempfile(); dir.create(d); writeLines("x", file.path(d, "f.R")); dir(d, pattern = "R$")
#---
list.files(tempfile())
#---
list.files("/nonexistent")
#---
d <- tempfile(); dir.create(d); dir.create(file.path(d, "b")); dir.create(file.path(d, "a")); dir.create(file.path(d, "a", "x")); writeLines("f", file.path(d, "f.txt")); sub(d, "D", list.dirs(d), fixed = TRUE)
#---
d <- tempfile(); dir.create(d); dir.create(file.path(d, "b")); dir.create(file.path(d, "a")); dir.create(file.path(d, "a", "x")); sub(d, "D", list.dirs(d, recursive = FALSE), fixed = TRUE)
#---
d <- tempfile(); dir.create(d); dir.create(file.path(d, "b")); dir.create(file.path(d, "a")); dir.create(file.path(d, "a", "x")); list.dirs(d, full.names = FALSE)
#---
d <- tempfile(); dir.create(d); dir.create(file.path(d, "b")); dir.create(file.path(d, "a")); list.dirs(d, full.names = FALSE, recursive = FALSE)
#---
nchar(basename(tempfile())) > 4
#---
startsWith(basename(tempfile("abc")), "abc")
#---
endsWith(tempfile(fileext = ".txt"), ".txt")
#---
startsWith(basename(tempfile("abc", fileext = ".csv")), "abc") && endsWith(tempfile("abc", fileext = ".csv"), ".csv")
#---
length(tempfile(c("a", "b", "c")))
#---
substr(basename(tempfile(c("a", "b"))), 1, 1)
#---
length(unique(sapply(1:20, function(i) tempfile())))
#---
dirname(tempfile()) == tempdir()
#---
dirname(tempfile(tmpdir = "/foo/bar")) == "/foo/bar"
#---
dir.exists(tempdir())
#---
tempdir() == tempdir()
#---
file.exists(tempfile())
#---
basename(tempdir()) != ""
#---
old <- getwd(); d <- tempfile(); dir.create(d); setwd(d); writeLines("rel", "r.txt"); r <- list(readLines("r.txt"), file.exists(file.path(d, "r.txt")), list.files()); setwd(old); r
#---
old <- getwd(); d <- tempfile(); dir.create(d); r <- setwd(d); setwd(old); identical(r, old)
#---
setwd("/nonexistent/dir")
#---
old <- getwd(); setwd(tempdir()); r <- identical(getwd(), normalizePath(tempdir())); setwd(old); r
#---
is.character(getwd())
#---
normalizePath("/nonexistent/dir/zz")
#---
tryCatch(normalizePath("/nonexistent/dir/zz"), warning = function(w) conditionMessage(w))
#---
normalizePath("/nonexistent/dir/zz", mustWork = TRUE)
#---
normalizePath("/nonexistent/dir/zz", mustWork = FALSE)
#---
normalizePath(c("/nonexistent/a", "/nonexistent/b"), mustWork = FALSE)
#---
tryCatch(normalizePath(c("/", "/nonexistent/b")), warning = function(w) conditionMessage(w))
#---
normalizePath("/")
#---
path.expand("/a/b")
#---
path.expand(c("x", "/y"))
#---
path.expand(1)
#---
Sys.getenv("RGO_NOPE_VAR")
#---
Sys.getenv("RGO_NOPE_VAR", unset = NA)
#---
Sys.getenv("RGO_NOPE_VAR", "dflt")
#---
Sys.getenv(c("RGO_NOPE_A", "RGO_NOPE_B"))
#---
Sys.getenv(c("RGO_NOPE_A", "RGO_NOPE_B"), names = FALSE)
#---
Sys.getenv("RGO_NOPE_A", names = TRUE)
#---
Sys.setenv(RGO_A = "1", RGO_B = "two")
#---
Sys.setenv(RGO_A = "1"); Sys.getenv("RGO_A")
#---
r <- Sys.setenv(RGO_A = "1", RGO_B = "two"); r
#---
Sys.setenv(RGO_A = "1", RGO_B = "two"); Sys.getenv(c("RGO_A", "RGO_B"))
#---
Sys.setenv(RGO_A = "1", RGO_B = "two"); Sys.getenv(c("RGO_A", "RGO_B"), names = FALSE)
#---
Sys.setenv(RGO_A = "1"); Sys.setenv(RGO_A = "2"); Sys.getenv("RGO_A")
#---
Sys.setenv(RGO_A = 5); Sys.getenv("RGO_A")
#---
Sys.setenv(RGO_A = "1"); r <- Sys.unsetenv("RGO_A"); list(r, Sys.getenv("RGO_A"), Sys.getenv("RGO_A", NA))
#---
Sys.unsetenv("RGO_NOPE_VAR")
#---
Sys.setenv("a")
#---
Sys.setenv(RGO_A = "1", "b")
#---
Sys.setenv(RGO_A = c("1", "2"))
#---
Sys.setenv(RGO_E = ""); Sys.getenv("RGO_E", unset = NA)
#---
Sys.setenv(RGO_ONLY = "z"); "RGO_ONLY" %in% names(Sys.getenv())
#---
Sys.sleep(0)
#---
Sys.sleep(0.01)
#---
Sys.sleep(-1)
#---
Sys.sleep("a")
#---
capture.output(print(1:3))
#---
capture.output(1 + 1, 2 + 2)
#---
capture.output({1 + 1; 2 + 2})
#---
capture.output(cat("a\nb\n"), cat("c"))
#---
capture.output(invisible(5))
#---
capture.output(x <- 5)
#---
capture.output()
#---
capture.output(NULL)
#---
capture.output(print("a"), print(c(a = 1, b = 2)))
#---
capture.output(cat("no newline"))
#---
capture.output(cat("a\n\nb\n"))
#---
capture.output(matrix(1:4, 2))
#---
capture.output(list(1, "a"))
#---
capture.output(factor(c("u", "v")))
#---
x <- capture.output(print(letters[1:3])); class(x)
#---
capture.output(writeLines(c("w1", "w2")))
#---
capture.output(message("to stderr"))
#---
capture.output(message("m1"), type = "message")
#---
capture.output(print("both"), message("err"), type = "message")
#---
f <- tempfile(); r <- capture.output(print(1:2), file = f); list(r, readLines(f));
#---
f <- tempfile(); capture.output(print(1:2), file = f); capture.output(print(3), file = f, append = TRUE); readLines(f)
#---
f <- tempfile(); capture.output(print(1:2), file = f); capture.output(print(3), file = f); readLines(f)
#---
f <- tempfile(); con <- file(f, "w"); capture.output(print(1:2), file = con); writeLines("after", con); close(con); readLines(f)
#---
n0 <- sink.number(); f <- function() capture.output(stop("boom")); try(f(), silent = TRUE); sink.number() - n0
#---
n0 <- sink.number(); r <- tryCatch(capture.output(stop("boom")), error = function(e) conditionMessage(e)); list(r, sink.number() - n0)
#---
capture.output(capture.output(print(1)))
#---
x <- capture.output(print(1)); y <- capture.output(print(2)); c(x, y)
#---
sink.number() >= 0
#---
f <- tempfile(); sink(f); print(1:3); cat("x\n"); sink(); readLines(f)
#---
n0 <- sink.number(); f <- tempfile(); sink(f); n <- sink.number(); sink(); c(n - n0, sink.number() - n0)
#---
f <- tempfile(); sink(f); print(1); sink(); sink(f, append = TRUE); print(2); sink(); readLines(f)
#---
f <- tempfile(); sink(f); print("in"); sink(); print("out"); readLines(f)
#---
f <- tempfile(); con <- file(f, "w"); sink(con); print(1); sink(); print(2); close(con); readLines(f)
#---
f <- tempfile(); sink(f, split = TRUE); print("both"); sink(); readLines(f)
#---
tc <- textConnection("sunk", "w"); sink(tc); print(1:3); sink(); close(tc); sunk
#---
tc <- textConnection("sunk", "w"); sink(tc); print(1:3); cat("partial"); sink(); close(tc); sunk
#---
n0 <- sink.number(); f1 <- tempfile(); f2 <- tempfile(); sink(f1); print("one"); sink(f2); print("two"); r <- sink.number() - n0; sink(); print("three"); sink(); list(r, readLines(f1), readLines(f2))
#---
sink("/nonexistent/dir/x")
#---
sink(5)
#---
f <- tempfile(); sink(f); print(1:3); invisible(sink()); length(readLines(f))
#---
f <- tempfile(); saveRDS(1:3, f); r <- readRDS(f); unlink(f); r
#---
f <- tempfile(); saveRDS(c(a = 1.5, b = NA, c = Inf), f); r <- readRDS(f); unlink(f); r
#---
f <- tempfile(); x <- list(1, "a", TRUE, NULL, list(2, 3)); saveRDS(x, f); r <- identical(readRDS(f), x); unlink(f); r
#---
f <- tempfile(); x <- c("a", NA, "é中"); saveRDS(x, f); r <- identical(readRDS(f), x); unlink(f); r
#---
f <- tempfile(); x <- factor(c("u", "v", "u")); saveRDS(x, f); r <- readRDS(f); unlink(f); r
#---
f <- tempfile(); x <- matrix(1:6, 2, dimnames = list(c("a", "b"), NULL)); saveRDS(x, f); r <- identical(readRDS(f), x); unlink(f); r
#---
f <- tempfile(); x <- structure(list(a = 1), class = "foo"); saveRDS(x, f); r <- readRDS(f); unlink(f); class(r)
#---
f <- tempfile(); x <- as.Date("2024-03-05"); saveRDS(x, f); r <- readRDS(f); unlink(f); r
#---
f <- tempfile(); x <- quote(a + b * 2); saveRDS(x, f); r <- readRDS(f); unlink(f); r
#---
f <- tempfile(); saveRDS(NULL, f); r <- readRDS(f); unlink(f); r
#---
f <- tempfile(); saveRDS(1:10, f, compress = FALSE); r <- readRDS(f); unlink(f); r
#---
f <- tempfile(); saveRDS(1:10, f, version = 2); r <- readRDS(f); unlink(f); r
#---
f <- tempfile(); saveRDS("x", f); r <- saveRDS("y", f); unlink(f); r
#---
f <- tempfile(); saveRDS(1:3, f); file.exists(f)
#---
f <- tempfile(); saveRDS(1:3, f); con <- file(f); r <- readRDS(con); close(con); unlink(f); r
#---
readRDS("/nonexistent/dir/x.rds")
#---
tryCatch(readRDS("/nonexistent/dir/x.rds"), warning = function(w) conditionMessage(w))
#---
f <- tempfile(); writeLines("not an rds file", f); readRDS(f)
#---
saveRDS(1:3, "/nonexistent/dir/x.rds")
#---
f <- tempfile(); saveRDS(globalenv(), f); r <- readRDS(f); unlink(f); identical(r, globalenv())
#---
tc <- textConnection("x_src <- 5; y_src <- x_src * 2"); source(tc); close(tc); y_src
#---
tc <- textConnection("1 + 1\n2 + 2"); r <- source(tc); close(tc); r
#---
tc <- textConnection("invisible(3)"); r <- source(tc); close(tc); r
#---
tc <- textConnection("1 + 1"); r <- source(tc); close(tc); r$visible
#---
tc <- textConnection("1 + 1"); source(tc, print.eval = TRUE); close(tc)
#---
tc <- textConnection("print('p')\ncat('c\\n')"); source(tc); close(tc)
#---
tc <- textConnection("cat('hi\\n')"); source(tc, echo = FALSE); close(tc)
#---
f <- tempfile(); writeLines(c("a_src <- 1", "b_src <- a_src + 1"), f); source(f); r <- b_src; unlink(f); r
#---
f <- tempfile(); writeLines("z_src <- 9", f); g <- function() { source(f, local = TRUE); z_src }; r <- list(g(), exists("z_src")); unlink(f); r
#---
f <- tempfile(); writeLines("z_src <- 9", f); e <- new.env(); source(f, local = e); r <- list(e$z_src, exists("z_src")); unlink(f); r
#---
f <- tempfile(); writeLines("z_src <- 9", f); source(f, local = FALSE); z_src
#---
f <- tempfile(); writeLines("v_src <- 7", f); e <- new.env(); sys.source(f, envir = e); r <- e$v_src; unlink(f); r
#---
f <- tempfile(); writeLines("v_src <- 7", f); e <- new.env(); r <- sys.source(f, envir = e); unlink(f); r
#---
source("/nonexistent/dir/x.R")
#---
tc <- textConnection("stop('in source')"); tryCatch(source(tc), finally = close(tc))
#---
tc <- textConnection("stop('in source')"); r <- tryCatch(source(tc), error = function(e) conditionMessage(e)); close(tc); r
#---
tc <- textConnection("x <-"); tryCatch(source(tc), finally = close(tc))
#---
tc <- textConnection(character(0)); tryCatch(source(tc), finally = close(tc))
#---
tc <- textConnection(character(0)); r <- source(tc); close(tc); r
#---
f <- tempfile(); writeLines(c("sq <- function(x) x^2", "sq(4)"), f); r <- source(f); unlink(f); r
#---
f <- tempfile(); writeLines("1 + 1", f); con <- file(f); r <- source(con); close(con); unlink(f); r$value
#---
tc <- textConnection("a1 <- 1\nb1 <- 2"); source(tc); close(tc); a1 + b1
