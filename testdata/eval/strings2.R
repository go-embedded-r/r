strrep("ab", 3); strrep(c("a", "b"), 1:2); strrep("x", 0); strrep(NA, 2); strrep("a", c(1, NA)); strrep(c(a = "z"), 2)
#---
strrep("a", -1)
#---
x <- "abcdef"; substr(x, 2, 3) <- "ZZ"; x; substr(x, 1, 2) <- "Q"; x; substr(x, 5, 10) <- "123456"; x; y <- c("aaa", "bbb"); substr(y, 2, 2) <- c("X", "Y"); y; z <- NA_character_; substr(z, 1, 1) <- "a"; z
#---
x <- 1; substr(x, 1, 1) <- "a"
#---
strtoi("77", 8L); strtoi("ff", 16L); strtoi("0x1A", 16L); strtoi("12"); strtoi("zz", 36L); strtoi("9", 8L); strtoi(c("1", NA, "x")); strtoi("123456789012"); strtoi(" 12")
#---
make.unique(c("a", "a", "b", "a")); make.unique(c("a", "a"), sep = "_"); make.unique(c("a", "a", "a.1")); make.unique(character(0))
#---
basename("a/b/c.txt"); basename("a/b/"); basename("c"); basename("/"); dirname("a/b/c.txt"); dirname("a"); dirname("/a"); dirname("/"); dirname("a/b/"); basename(c("x/y", NA))
#---
formatC(3.14159); formatC(3.14159, digits = 3, format = "f"); formatC(42L); formatC(42, width = 8, flag = "0"); formatC(1234567.891, format = "f", big.mark = ",", digits = 1); formatC(c(1, 10, 100), width = 5); formatC("a", width = 5); formatC("a", width = -5); formatC("a", flag = "-", width = 5); formatC(0.000012345, format = "e", digits = 2); formatC(123456789, format = "d", big.mark = ","); formatC(1e6, format = "g"); formatC(c(a = 1.5, b = 2))
#---
formatC(1.5, format = "z")
#---
prettyNum(1234567, big.mark = ","); prettyNum(c(123, 1234, 12345678), big.mark = ","); prettyNum("1234567.891", big.mark = ","); prettyNum(123); prettyNum(NA, big.mark = ",")
#---
toString(1:5); toString(c("a", "b")); toString(character(0)); toString(NULL)
#---
sQuote("a"); dQuote("a"); sQuote("a", FALSE); dQuote(c("a", "b"), FALSE); sQuote(character(0)); sQuote(1:2)
#---
shQuote("abc"); shQuote("it's"); shQuote(c("a b", "c"))
#---
x <- noquote(c("a", "b")); x; class(x); print(c("a", "b"), quote = FALSE); print("q", quote = TRUE); noquote(c(a = "x")); noquote(1:3)
#---
file.path("a", "b", "c.txt"); file.path("a", c("b", "c")); file.path(); file.path("a", character(0)); file.path("a", "b", fsep = "\\")

#---
formatC(5, flag = "+"); formatC(-5, width = 6, flag = "0"); formatC(5.5, flag = " "); formatC(c(1.5, NA)); formatC(Inf)
#---
formatC(TRUE)
#---
x <- noquote("a"); print.noquote(x)

#---
formatC(NA_real_, digits = 2, format = "f"); formatC(NA_real_, digits = 5, format = "f"); formatC(c(1, NA), format = "f", digits = 1, width = 8); formatC(c(1.5, NA), format = "d"); formatC(2.5, format = "d"); formatC(-1.5, format = "d")
#---
prettyNum(1234567.891, big.mark = ","); strtoi(c(a = "q"), 36L); strrep("abc", NA)

#---
formatC("a", format = "f"); formatC("a", format = "d"); formatC("1.5", format = "e"); formatC(1.5, format = "s"); formatC(1L, format = "f"); formatC("a", digits = 2); formatC(c("a", "bb"), format = "s", width = 4)
#---
formatC(1L, format = "s")

#---
shQuote("a b", "cmd"); shQuote("a\"b", "cmd"); shQuote("a", "csh"); shQuote("it's", "csh"); shQuote("a!b", "csh"); shQuote("a b", "cmd2"); shQuote("a\"b", "cmd2"); shQuote("x", type = "sh")
#---
shQuote("x", "bogus")
#---
formatC(3.14159, format = "fg", digits = 3); formatC(1234.5678, format = "fg", digits = 2); formatC(0.000123456, format = "fg", digits = 3); formatC(1.5, format = "fg", digits = 5); formatC(1.5, format = "fg", digits = 5, flag = "#"); formatC(100, format = "fg"); formatC(2, mode = "integer"); formatC(2.7, mode = "integer"); formatC(2L, mode = "double"); formatC(1e10, format = "fg", digits = 3, big.mark = ","); formatC("a", mode = "character")
#---
formatC(1, mode = "bogus")
#---
format(1234567.891, big.mark = ","); format(c(1, 1234, 1234567), big.mark = ","); format(-1234.5, big.mark = ","); format(1234L, big.mark = ","); format("1234567", big.mark = ","); format(1e6, big.mark = ",", scientific = FALSE); format(1234.5678, nsmall = 2, big.mark = ","); format(c(a = 1000, b = 10), big.mark = ","); format(NA, big.mark = ","); format(123456789.123, big.mark = ",", nsmall = 1)
#---
f <- function(type = c("one", "two", "three")) match.arg(type); f(); f("two"); f("th"); f(NULL)
#---
f <- function(type = c("one", "two", "three")) match.arg(type); f("x")
#---
f <- function(type = c("one", "two")) match.arg(type); f(c("one", "two")); f(1)
#---
match.arg("me", c("mean", "median")); match.arg("mea", c("mean", "median")); match.arg(c("a", "b"), c("a", "b", "c"), several.ok = TRUE)
#---
match.arg("m", c("mean", "median"))
#---
match.arg(NULL, c("x", "y"))
