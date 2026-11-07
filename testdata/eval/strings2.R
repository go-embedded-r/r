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
