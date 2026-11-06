grepl("a", c("abc", "xyz", NA)); grepl("^a", c("abc", "bac")); grepl("c$", "abc"); grepl("", "abc")
#---
grepl("A", "abc", ignore.case = TRUE); grepl("a.c", "a\nc"); grepl("a|b", c("a", "b", "c")); grepl("[[:digit:]]+", c("a1", "b")); grepl("\\d", c("a1", "b")); grepl("\\w+", "!"); grepl("\\s", "a b")
#---
grepl(".", "abc", fixed = TRUE); grepl("a.c", "a.c", fixed = TRUE); grepl("+", "a+b", fixed = TRUE)
#---
grepl("(a", "a")
#---
grepl("[a", "a")
#---
grepl("a{2,3}", c("a", "aa", "aaaa")); grepl("(ab)+", "ababab"); grepl("a*", "b"); grepl("colou?r", c("color", "colour"))
#---
grepl(NA, "a")
#---
grepl("a", NULL); grepl("a", character(0))
#---
grep("a", c("abc", "xyz", "bca")); grep("a", c("abc", "xyz", "bca"), value = TRUE); grep("a", c("abc", "xyz"), invert = TRUE); grep("a", c(x = "abc", y = "q"), value = TRUE); grep("z", "a")
#---
grep("a", c("a", NA, "ba"))
#---
sub("a", "X", "banana"); gsub("a", "X", "banana"); sub("x", "Y", "abc"); gsub("[aeiou]", "", "education"); gsub("(a)(n)", "\\2\\1", "banana")
#---
gsub("a", "\\\\", "banana"); gsub("a", "\\0\\0", "banana"); sub("(.*)", "<\\1>", "abc"); gsub(".", "-", "a.b", fixed = TRUE); gsub("\\.", "-", "a.b")
#---
gsub("", "-", "abc"); gsub("b*", "-", "abc"); gsub("x*", "-", "abc"); sub("", "-", "abc"); gsub("^", ">", "abc"); gsub("$", "<", "abc")
#---
gsub("a", "X", c(x = "aa", y = "b")); gsub("a", "X", NA); gsub("a", "X", c("a", NA)); gsub(NA, "X", "a"); sub("a", NA, "abc")
#---
gsub("a", "X", 1:3); gsub("1", "X", 111); gsub("A", "x", "aAa", ignore.case = TRUE)
#---
gsub("é", "e", "café"); gsub(".", "x", "日本"); sub("本", "X", "日本語")
#---
gsub("a.*?b", "X", "aXXbaYYb"); gsub("a.*b", "X", "aXXbaYYb"); sub("(a+?)", "<\\1>", "aaa")
#---
gsub("(a|ab)(c|bcd)", "[\\1,\\2]", "abcd"); gsub("a|ab", "X", "abab")
#---
gsub("[[:space:]]+", " ", "a   b \t c"); gsub("[^[:alnum:]]", "", "a-b_c!"); gsub("[[:upper:]]", "U", "aBcD"); gsub("[\\d]", "X", "a\\d1")
#---
gsub("a", "b", "x", fixed = TRUE, ignore.case = TRUE)
#---
regexpr("a", c("banana", "xyz", NA)); regexpr("an", "banana"); regexpr("", "abc")
#---
regexpr("é", "café"); regexpr("a+", "baaac"); regexpr("a", "banana", fixed = TRUE)
#---
gregexpr("an", "banana"); gregexpr("z", "banana"); gregexpr("", "abc"); gregexpr("b*", "abc"); gregexpr("a", c("aa", NA))
#---
gregexpr("a+", "baaacaa")
#---
x <- c("apple", "banana", "cherry"); regmatches(x, regexpr("[aeiou]", x)); regmatches(x, gregexpr("[aeiou]", x))
#---
x <- "a1b22c333"; regmatches(x, gregexpr("[0-9]+", x)); regmatches(x, regexpr("[0-9]+", x), invert = TRUE); regmatches(x, gregexpr("[0-9]+", x), invert = TRUE)
#---
x <- c(a = "x1", b = "y"); regmatches(x, regexpr("[0-9]", x))
#---
startsWith(c("apple", "banana", NA), "a"); endsWith("apple", "le"); startsWith("a", c("a", "b")); startsWith(1, "a")
#---
grepl("[[:alpha:]", "a")
#---
grepl("a\\", "a")
#---
grepl("*a", "a")
#---
grepl("a**", "a")
#---
grepl("a{2", "a")
#---
grepl("(", "a", fixed = TRUE)
#---
grepl("[[:foo:]]", "a")
#---
grepl("a)", "a)")
#---
gsub("(?i)A", "x", "aA")
#---
gsub("\\\\", "/", "a\\b")
#---
gsub("[.]", "", "a.b.c"); gsub("[]]", "", "a]b"); gsub("[^]]", "", "a]b")
#---
grepl("\\.", "a.b"); grepl("\\$", "a$"); grepl("[$]", "a$"); grepl("a\\|b", "a|b")
#---
sprintf("%s", gsub("\\s+$", "", "abc   "))
#---
grepl("a", 123); grepl("1", 123)
#---
gsub("\\<a", "[", "ab ab"); gsub("b\\>", "]", "ab ab")
#---
gsub("\\bab", "X", "ab cab ab")
#---
grepl("*a", "a"); grepl("*a", "*a"); grepl("^*a", "*a"); grepl("a{,2}", "aa"); grepl("{1}a", "a"); grepl("a||b", "a"); grepl("()", "a")
#---
grepl("a{x}", "a")
#---
grepl("a{2,1}", "a")
#---
grepl("[b-a]", "a")
#---
sort(c(3, 1, 2)); sort(c(3, 1, NA, 2), decreasing = TRUE); sort(c("b", "a")); sort(c(b = 2, a = 1)); sort(NULL); sort(c(TRUE, FALSE))
#---
strsplit("a1b2", "[0-9]"); strsplit("2020-01-02", "-"); strsplit("a  b", " +")
#---
strsplit("abc", "(?<=a)")
#---
gregexpr("é", "café au lait é")
#---
gregexpr("a", "café a")
#---
regexpr("a", "café")
#---
gsub("(?:a)+", "X", "caab"); gsub("(?i)b", "X", "aBb"); grepl("(?i)abc", "ABC")
#---
sub("a+?", "X", "caaab"); sub("a*?b", "X", "caaab"); regexpr("a+?", "caaab"); gsub("b??", "-", "abba")
#---
regmatches("aaa", regexpr("a{1,3}?", "aaa")); regmatches("foo bar", gregexpr("o+?", "foo bar"))
#---
grepl("(a)\\1", c("aa", "ab", "a")); sub("(a+)b\\1", "[\\1]", "caabaad"); gsub("(.)\\1", "<\\1>", "aabbcd")
#---
regexpr("(a*)\\1", "aaaa"); regexec <- NULL; grepl("([abc]+)\\w*\\1", c("banana", "foo bar"))
#---
gsub("\\<", "|", "ab cd"); gsub("\\>", "|", "ab cd"); gsub("\\bc", "X", "ab cd"); grepl("\\<ab\\>", c("ab", "cab", "ab c"))
#---
regmatches("abcd", regexpr("(a|ab)(c|bcd)(d*)", "abcd")); sub("(a*)*", "[\\1]", "aab"); sub("(a|ab)(c|bcd)(d*)", "\\1,\\2,\\3", "abcd")
#---
grepl("^a{,3}$", c("aaaa", "aaa")); grepl("a{256}", "a")
#---
strsplit(c("Abc", "héllo"), "\\d*"); grepl("b[^[:space:]]{1,3}", c("ab ab", "é"))
