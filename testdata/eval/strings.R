nchar("hello"); nchar(c("a", "bb", NA)); nchar(NA); nchar(123); nchar("")
#---
nchar(c(a = "xx", b = "y"))
#---
nchar("héllo"); nchar("héllo", type = "bytes"); nchar("日本語", type = "width"); nchar(NA_character_, type = "width")
#---
nchar(list("a"))
#---
nchar("abc", type = "nope")
#---
substr("hello", 2, 4); substr("hello", 0, 2); substr("hello", 4, 100); substr("hello", 3, 2); substr(c("abc", "defg"), 2, 3); substr(NA, 1, 2)
#---
substr(c(a = "hello", b = "world"), 1, 3)
#---
substr("héllo", 2, 3)
#---
substr("abc", 1)
#---
substring("hello", 1:3, 3:5); substring("hello", 2); substring(c("ab", "cd"), 1, 1)
#---
toupper("abc"); tolower("ABC"); toupper(c(a = "x")); toupper(NA); toupper(1)
#---
casefold("AbC"); casefold("AbC", upper = TRUE)
#---
chartr("abc", "xyz", "aabbcc"); chartr("a-c", "A-C", "abcdef")
#---
chartr("abc", "x", "abc")
#---
trimws("  hi  "); trimws("  hi  ", which = "left"); trimws("  hi  ", which = "right"); trimws(c(" a", "b "), whitespace = "[ ]")
#---
strsplit("a,b,c", ","); strsplit("a b  c", " "); strsplit("abc", ""); strsplit(c(x = "a-b", y = "c"), "-")
#---
strsplit("a.b.c", ".", fixed = TRUE); strsplit("a1b22c", "[0-9]+"); strsplit(",a,,b,", ","); strsplit("", ",")
#---
strsplit(NA_character_, ","); strsplit("ab", "x*")
#---
strsplit(1, "a")
#---
Encoding("abc"); Encoding("é"); enc2utf8("é"); enc2native("a"); validUTF8(c("a", "é"))
#---
utf8ToInt("A"); utf8ToInt("é"); intToUtf8(c(72, 105)); intToUtf8(c(72, 105), multiple = TRUE); intToUtf8(NULL)
#---
sprintf("%d items", 3L); sprintf("%s is %d", "x", 5L); sprintf("%5.2f", pi); sprintf("%-6s|", "ab"); sprintf("%06.2f", 3.14159); sprintf("%e", 12345.6789); sprintf("%g", 1e-5); sprintf("%g", 123456789)
#---
sprintf("%d", 3); sprintf("%i", 7L); sprintf("%x", 255L); sprintf("%X", 255L); sprintf("%o", 8L); sprintf("%%"); sprintf("%5s", "ab"); sprintf("%.2s", "abcdef")
#---
sprintf("%d", 3.5)
#---
sprintf("%s", 1.5); sprintf("%s", TRUE); sprintf("%s", 1L); sprintf("%d", TRUE); sprintf("%s", 1e6); sprintf("%s", 123456.7)
#---
sprintf("%d", NA_integer_); sprintf("%f", NA_real_); sprintf("%s", NA_character_); sprintf("%5.1f", NaN); sprintf("%f", Inf); sprintf("%f", -Inf)
#---
sprintf(c("%d a", "%d b"), 1:2); sprintf("%d", 1:3); sprintf("%s-%s", c("a", "b"), "z"); sprintf(character(0)); sprintf("x", integer(0))
#---
sprintf("%2$s %1$s", "world", "hello"); sprintf("%*d", 5, 42L); sprintf("%+d", 5L); sprintf("% d", 5L); sprintf("%#x", 255L)
#---
sprintf("%d")
#---
sprintf("%d", 1L, 2L)
#---
sprintf("%z", 1L)
#---
sprintf("%s", list("a"))
#---
sprintf(1)
#---
format(3.14159); format(1e6); format(1e-10); format(123456789); format(0.1 + 0.2); format(100); format(c(1, 10, 100))
#---
format(3.14159, nsmall = 2); format(2, nsmall = 2); format(1234.5678, digits = 3); format(0.5, digits = 1)
#---
format(TRUE); format(c(TRUE, FALSE, NA)); format(NA); format(5L); format(c(1L, 100L))
#---
format("a"); format(c("a", "bbb")); format(c("a", "bbb"), justify = "right"); format(c("a", "bbb"), justify = "centre"); format("a", width = 5); format(NA_character_)
#---
format(c(a = 1.5, b = 22.25)); format(c(1.5, NA, 3)); format(NULL); format(-1.5)
#---
format(1e5, scientific = TRUE); format(123456, scientific = FALSE); format(0.00001234, scientific = FALSE); format(12345.6789, nsmall = 1)
#---
format(42L, width = 6); format(3.5, width = 8); format(c(1, 22, 333), trim = TRUE)
#---
format(matrix(1:4, 2))
#---
format(list(1))