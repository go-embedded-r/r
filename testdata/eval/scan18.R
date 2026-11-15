scan(text = "1 2 3")
#---
scan(text = "1,2,,4", sep = ",")
#---
scan(text = "a 'b c' d", what = "")
#---
scan(text = "a,'b,c',d", what = "", sep = ",")
#---
scan(text = "a,\"b\"\"c\",d", what = "", sep = ",")
#---
scan(text = "a 'b\\'c' d", what = "")
#---
scan(text = "1 2 NA 4 x", what = "")
#---
scan(text = "1 2 NA 4")
#---
scan(text = "1 2 abc 4")
#---
scan(text = "1 a\n2 b\n3 c", what = list(0, ""))
#---
scan(text = "1 a\n2 b\n3", what = list(0, ""))
#---
scan(text = "1 a\n2 b\n3", what = list(0, ""), fill = TRUE)
#---
scan(text = "1 a x\n2 b y\n3 c z", what = list(0, ""), flush = TRUE)
#---
scan(text = "# c\n1 2 # d\n3", comment.char = "#")
#---
scan(text = "x\ny\n1 2\n\n3", skip = 2)
#---
scan(text = "1 2\n3 4\n5 6", nlines = 2)
#---
scan(text = "1 2 3 4 5", n = 3)
#---
scan(text = "1 2 3 4 5", nmax = 2)
#---
scan(text = "1,2,3\n4,5,6", sep = ",", what = list(a = 0, b = 0, c = 0))
#---
scan(text = " a , b ", what = "", sep = ",", strip.white = TRUE)
#---
scan(text = " a , b ", what = "", sep = ",")
#---
scan(text = "1,5 2,5", dec = ",")
#---
scan(text = "a\\tb c", what = "", allowEscapes = TRUE)
#---
scan(text = "a\\tb c", what = "")
#---
scan(text = "TRUE F T false", what = TRUE)
#---
scan(text = "1 2 3", what = integer())
#---
scan(text = "1.5 2", what = integer())
#---
scan(text = "1+2i 3", what = complex())
#---
scan(text = "a b c", what = "", quiet = TRUE)
#---
scan(text = "a,b\n\nc,d", what = "", sep = ",", blank.lines.skip = FALSE)
#---
scan(text = "a,b\n\nc,d", what = "", sep = ",")
#---
scan(text = "a b\nc d", what = "", sep = "\n")
#---
scan(text = "0x1F 1e3 Inf -inf NaN", what = 0)
#---
scan(text = "'abc", what = "")
#---
scan(text = "a 'b c\nd' e", what = "")
#---
scan(text = "a 'b c\nd' e", what = "", multi.line = FALSE)
#---
scan(text = "1 a\n2 b\n3 c", what = list(0, ""), multi.line = FALSE)
#---
scan(text = "1 a\n2\n3 c", what = list(0, ""), multi.line = FALSE)
#---
scan(text = "1 2 3 4 5 6", what = list(0, 0, 0))
#---
scan(text = "1 2 3 4 5", what = list(0, 0))
#---
scan(text = "x\"y z", what = "")
#---
scan(text = "a,b,c", what = "", sep = ",", quote = "")
#---
scan(text = "  ", what = 0)
#---
scan(text = "", what = 0)
#---
scan(text = c("1 2", "3"), what = 0)
#---
scan(text = "NA 1", what = "", na.strings = c("NA", "1"))
#---
scan(text = "'NA' NA", what = "")
#---
scan(text = "1;2;3", sep = ";")
#---
scan(text = "1\t2\t\t3", sep = "\t")
#---
scan(text = "a 1\nb 2", what = list(name = "", val = 0))
#---
scan(text = "1 2", what = NULL)
#---
scan(text = "1 2", what = list(NULL, 0))
#---
scan(text = "1,2", sep = ",,")
#---
scan(text = "1 2", sep = 1)
#---
scan(text = "1 2 3", nmax = 0)
#---
scan(text = "1 2 3", n = 0)
#---
scan("/nonexistent/file")
#---
f <- tempfile(); writeLines(c("1 2", "3 4"), f); scan(f); scan(f, nlines = 1); scan(file(f), n = 3)
#---
f <- tempfile(); writeLines(c("1 2", "3 4"), f); con <- file(f, "r"); scan(con, n = 1); scan(con, n = 2); scan(con); close(con)
#---
tc <- textConnection("5 6\n7 8"); scan(tc, nlines = 1); scan(tc); close(tc)
#---
scan(text = "1 2 3", what = list(0, 0, 0), nmax = 1)
#---
scan(text = "\"a b\"c d", what = "")
#---
scan(text = "'a''b'", what = "")
#---
scan(text = "'a''b'", what = "", sep = ",")
#---
scan(text = "1 2\r\n3 4\r\n")
#---
scan(text = "a b", what = "", comment.char = "##")
#---
scan(text = "a b # x", what = "", comment.char = "#")
#---
scan(text = "a,'b # x',c", what = "", sep = ",", comment.char = "#")
#---
scan(text = "TRUE T yes", what = TRUE)
#---
scan(text = "1e", what = 0)
#---
scan(text = " 1 ", what = 0, sep = ",")
#---
scan(text = "1,,3", what = 0L, sep = ",")
#---
scan(text = "a,,c", what = "", sep = ",")
#---
scan(text = "1 2 3\n4 5", what = list(0, 0, 0), fill = TRUE)
#---
scan(text = "\"a\"b,c", what = "", sep = ",", quiet = TRUE)
#---
scan(text = "x\"y,z\"w,c", what = "", sep = ",", quiet = TRUE)
#---
scan(text = "  \"a b\"  ,c", what = "", sep = ",", quiet = TRUE)
#---
scan(text = "  'a b'  ,c", what = "", sep = ",", strip.white = TRUE, quiet = TRUE)
#---
scan(text = "a\\nb c", what = "", allowEscapes = TRUE, quiet = TRUE)
#---
scan(text = "'a\\nb' c", what = "", allowEscapes = TRUE, quiet = TRUE)
#---
scan(text = "'a\\\\b' c", what = "", allowEscapes = TRUE, quiet = TRUE)
#---
scan(text = "'a\\x41b' c", what = "", allowEscapes = TRUE, quiet = TRUE)
#---
scan(text = "'a\\101b' c", what = "", allowEscapes = TRUE, quiet = TRUE)
#---
scan(text = "1,2\n3", what = list(0, 0), sep = ",", quiet = TRUE)
#---
scan(text = "1,2,3\n4,5,6", what = list(0, 0), sep = ",", quiet = TRUE)
#---
scan(text = "1,2,3\n4,5,6", what = list(0, 0), sep = ",", flush = TRUE, quiet = TRUE)
#---
scan(text = "1 2 3\n4 5 6", what = list(0, 0), quiet = TRUE)
#---
scan(text = "1 2 3\n4 5", what = list(0, 0), quiet = TRUE)
#---
scan(text = "1,2\n\n3,4", what = list(0, 0), sep = ",", quiet = TRUE)
#---
scan(text = "1,2\n\n3,4", what = list(0, 0), sep = ",", blank.lines.skip = FALSE, quiet = TRUE)
#---
scan(text = "1,2\n\n3,4", what = list(0, 0), sep = ",", blank.lines.skip = FALSE, fill = TRUE, quiet = TRUE)
#---
scan(text = "a b\n   \nc", what = "", blank.lines.skip = FALSE, quiet = TRUE)
#---
scan(text = "a,b\n   \nc", what = "", sep = ",", blank.lines.skip = FALSE, quiet = TRUE)
#---
scan(text = "a b c d", what = list("", ""), quiet = TRUE)
#---
scan(text = "TRUE 1 0 NA T", what = logical(), quiet = TRUE)
#---
scan(text = "1 2", what = list(a = 0, b = 0), nmax = 1, quiet = TRUE)
#---
scan(text = "1 2 3 4 5 6", nmax = 3, quiet = TRUE)
#---
scan(text = "a,b;c", what = "", sep = ";", quiet = TRUE)
#---
scan(text = " 12 ,13 ", sep = ",", what = 1L, quiet = TRUE)
#---
scan(text = "1e5L", quiet = TRUE)
#---
scan(text = "1d5 1D5", quiet = TRUE)
#---
scan(text = "+5 -5 .5 5. +.5e-1", quiet = TRUE)
#---
scan(text = "2147483648", what = 0L, quiet = TRUE)
#---
scan(text = "-2147483647 2147483647", what = 0L, quiet = TRUE)
#---
scan(text = "1.5", dec = ",", quiet = TRUE)
#---
scan(text = "1,5e2", dec = ",", quiet = TRUE)
#---
scan(text = "1+2i 1,5-1i", dec = ",", what = complex(), quiet = TRUE)
#---
scan(text = "1\n2\n3\n4", nlines = 2, skip = 1, quiet = TRUE)
#---
scan(text = "1\n2", skip = 5, quiet = TRUE)
#---
scan(text = "# a\n1", comment.char = "#", skip = 1, quiet = TRUE)
#---
scan(text = "a#b c", what = "", comment.char = "#", quiet = TRUE)
#---
scan(text = "'a#b' c", what = "", comment.char = "#", quiet = TRUE)
#---
scan(text = "1 2\n#x\n3", comment.char = "#", quiet = TRUE)
#---
scan(text = c("a b", "c"), what = "", sep = "\n", quiet = TRUE)
#---
scan(text = "a\n\nb", what = "", sep = "\n", quiet = TRUE)
#---
scan(text = "a\n\nb", what = "", sep = "\n", blank.lines.skip = FALSE, quiet = TRUE)
#---
scan(text = "a 'b\nc", what = "", sep = "\n", quiet = TRUE)
#---
scan(text = character(0), quiet = TRUE)
#---
scan(text = NA_character_, what = "", quiet = TRUE)
#---
scan(text = "1 2", what = "", na.strings = NULL, quiet = TRUE)
#---
scan(text = "NA", what = "", na.strings = NULL, quiet = TRUE)
#---
scan(text = "1 2", what = list(a = 0, b = ""), quiet = TRUE)
#---
scan(text = "a 'b c' d", what = list("", "", ""), quiet = TRUE)
#---
scan(text = "a b c\nd e", what = list("", "", ""), fill = TRUE, quiet = TRUE)
#---
scan(text = "a b c\nd e", what = list("", "", ""))
#---
scan(text = "a b c\nd e f g", what = list("", "", ""))
#---
scan(text = "a b c\nd e f g", what = list("", "", ""), quiet = TRUE, flush = TRUE)
#---
scan(text = "a,b,c\nd,e", what = list("", "", ""), sep = ",", quiet = TRUE)
#---
scan(text = "1\t\t3", what = list(0, 0, 0), sep = "\t", quiet = TRUE)
#---
scan(text = "1 2 3", what = 0L, n = 2, quiet = TRUE)
#---
scan(text = "a b c d e", what = list("", "", ""))
#---
scan(text = "a b c\nd e f\ng", what = list("", "", ""))
#---
scan(text = "a b c\nd", what = list("", "", ""))
#---
scan(text = "a b c d", what = list("", "", ""))
#---
scan(text = "a,b,c,d", what = list("", "", ""), sep = ",")
#---
scan(text = "1 2\n3\n5 6", what = list(0, 0), nlines = 2)
#---
scan(text = "a\n\n\nb", what = "", blank.lines.skip = FALSE, nlines = 2)
#---
scan(text = "\"a\nb\" c", what = "", nlines = 1)
#---
scan(text = "1 2 3", what = list(0, ""), fill = TRUE)
#---
scan(text = "a b c", what = list("", NULL, ""))
#---
scan(text = "1 b\n2 c", what = list(x = 0, y = ""), nlines = 1)
#---
scan(text = "1 2 3\n4 5 6", what = list(0, 0, 0), flush = TRUE, nmax = 1)
#---
scan(text = "a, b ,c", what = list("", "", ""), sep = ",", strip.white = TRUE)
#---
scan(text = " a ,\" b \" ,c", what = "", sep = ",", strip.white = TRUE)
#---
scan(text = "a,b", what = "", sep = ",", comment.char = "b")
#---
scan(text = "T F TRUE FALSE true false True False", what = TRUE)
#---
scan(text = "1 2", what = list(0, 0), nmax = -1, n = -1)
#---
scan(text = "1\n2\n3", what = list(a = 0))
#---
scan(text = " NA  1", what = 0L)
#---
scan(text = "1,2\n3,4", what = list(0, 0), sep = ",", skip = 1)
#---
scan(text = "a b\tc", what = "", sep = " ")
#---
scan(text = "a  b", what = "", sep = " ")
#---
scan(text = "1 2 3", what = list(0, 0, 0), multi.line = FALSE)
#---
scan(text = "1 2\n3", what = list(0, 0), multi.line = FALSE, fill = TRUE)
#---
scan(text = "1 2 3 4\n5 6", what = list(0, 0), multi.line = FALSE)
#---
scan(text = "1 2 ", what = list(0, 0, 0), multi.line = FALSE)
