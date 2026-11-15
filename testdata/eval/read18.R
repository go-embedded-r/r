read.table(text = "a b\n1 2\n3 4", header = TRUE)
#---
read.table(text = "1 2\n3 4")
#---
read.table(text = "a b\nr1 1 2\nr2 3 4")
#---
read.csv(text = "a,b\n1,x\n2,y")
#---
read.csv(text = "a,b\n1,x\n2,y", stringsAsFactors = TRUE)
#---
read.table(text = "a b\n1 NA\n- 4", header = TRUE, na.strings = c("NA", "-"))
#---
read.table(text = "a b\n1 2\n3 4", header = TRUE, colClasses = c("character", "numeric"))
#---
read.table(text = "a b\n1 2\n3 4", header = TRUE, colClasses = c(b = "character"))
#---
read.table(text = "a b\n1 2\n3 4", header = TRUE, colClasses = c("NULL", NA))
#---
read.table(text = "a b\n1 x\n3 y", header = TRUE, colClasses = c("integer", "factor"))
#---
read.table(text = "a b\n1 2\n3 4", header = TRUE, row.names = 1)
#---
read.table(text = "a b\n1 2\n3 4", header = TRUE, row.names = "b")
#---
read.table(text = "a b\nr1 1 2\nr2 3 4", header = TRUE, row.names = NULL)
#---
read.table(text = "a b\n1 2\n3 4", header = TRUE, row.names = c("x", "y"))
#---
read.table(text = "1 2\n3 4", col.names = c("u", "v"))
#---
read.table(text = "a a 1x a b\n1 2 3 4 5", header = TRUE)
#---
read.table(text = "a a 1x a b\n1 2 3 4 5", header = TRUE, check.names = FALSE)
#---
read.table(text = "junk\na b\n1 2", skip = 1, header = TRUE)
#---
read.table(text = "a b\n1 2\n3 4\n5 6", header = TRUE, nrows = 2)
#---
read.table(text = "a b\n1 2 # c\n# full\n3 4", header = TRUE)
#---
read.csv(text = "a,b\n1,2 # c\n3,4")
#---
read.csv(text = "a,b,c\n1,2\n3,4,5")
#---
read.table(text = "a b c\n1 2\n3 4 5", header = TRUE)
#---
read.table(text = "a b c\n1 2\n3 4 5", header = TRUE, fill = TRUE)
#---
read.csv(text = "a,b\n\n1,2\n\n3,4")
#---
read.csv(text = " a , b \n 1 , x ", strip.white = TRUE)
#---
read.csv(text = "a,b\n1,\"x,y\"\n2,\"p\nq\"")
#---
read.csv2(text = "a;b\n1,5;x\n2,5;y")
#---
read.delim(text = "a\tb\n1\tx y\n2\tz")
#---
read.delim2(text = "a\tb\n1,5\tx y\n2,5\tz")
#---
read.table(text = "a b\nTRUE 1\nFALSE 2", header = TRUE)
#---
read.table(text = "a b\nT x\nF y", header = TRUE)
#---
read.table(text = "a b\nT x\nF y", header = TRUE, tryLogical = FALSE)
#---
read.table(text = "", header = TRUE)
#---
read.table(text = "a b", header = TRUE)
#---
read.table(text = "a b\n1 2 3\n4 5 6")
#---
read.table(text = "a b c\n1 2 3\n4 5", header = TRUE)
#---
read.table(text = "a b\nr1 1 2\nr1 3 4")
#---
read.table(text = "1 2 3\n4 5 6", row.names = 1)
#---
read.table(text = "1 2\n3 4", as.is = FALSE)
#---
read.table(text = "a b\n1 x\n2 y", header = TRUE, as.is = c(FALSE, FALSE))
#---
read.table(text = "a b\n1 x\n2 y", header = TRUE, as.is = 1)
#---
read.table(text = "a b\n1 x\n2 y", header = TRUE, stringsAsFactors = TRUE)
#---
read.table(text = "a\n1\n2", header = TRUE)
#---
read.table(text = "'a b' c\n1 2", header = TRUE)
#---
read.table(text = "a b\n'x y' 2", header = TRUE)
#---
read.table(text = "a b\n1 2\n3", header = TRUE)
#---
read.table(text = c("a b", "1 2", "3 4"), header = TRUE)
#---
read.table(text = "1,2;3,4", sep = ";", dec = ",")
#---
read.table(text = "a b\n1 2", header = TRUE, colClasses = "character")
#---
read.table(text = "a b\n1 x", header = TRUE, colClasses = c("numeric", "numeric"))
#---
read.table(text = "a b\n2020-01-02 2", header = TRUE, colClasses = c("Date", "integer"))
#---
read.table(text = "a b\n1.5 2", header = TRUE, colClasses = c("integer", "integer"))
#---
read.table("/nonexistent")
#---
read.table(text = "a b\n1 2", header = TRUE, nrows = 0)
#---
read.table(text = "a b\n1 2", header = TRUE, nrows = -1)
#---
read.table(text = "a b\n1 2", header = TRUE, skip = 5)
#---
read.table(text = "\n\n\na b\n1 2", header = TRUE)
#---
read.table(text = "a b\n1 2", header = TRUE, blank.lines.skip = FALSE)
#---
read.table(text = "a b\n1 2\n\n3 4", header = TRUE, blank.lines.skip = FALSE)
#---
read.csv(text = "a,b\n1,2\n3,4", row.names = 1)
#---
read.csv(text = "a,a\n1,2")
#---
read.csv(text = ",b\n1,2")
#---
read.csv(text = "a,b\n1,2", header = FALSE)
#---
read.csv(text = "a,b\n1,\n,2")
#---
read.csv(text = "a,b\n1,\n,2", na.strings = "")
#---
read.csv(text = "a,b\n\"1\",\"x\"\n\"NA\",\"NA\"")
#---
read.csv(text = "x\n1e3\n0x10")
#---
read.csv(text = "x,y\n1,a\n2,b\n3,c", nrows = 2)
#---
read.csv(text = "x\n2147483648\n1")
#---
read.csv(text = "x\n007\n1")
#---
read.table(text = "a b\n1 2 9\n3 4 8", header = TRUE, flush = TRUE)
#---
read.table(text = "a b\n1 2 9\n3 4 8", header = TRUE)
#---
read.table(text = "a b\n1 2\n", header = TRUE, comment.char = "")
#---
read.table(text = "a\tb\n1\t2", header = TRUE, sep = "\t")
#---
f <- tempfile(); writeLines(c("a,b", "1,2", "3,4"), f); read.csv(f); read.csv(file(f)); con <- file(f, "r"); read.csv(con, nrows = 1); close(con)
#---
sapply(read.csv(text = "a,b,c,d\n1,x,2.5,TRUE\n2,y,NA,FALSE"), class)
#---
print(c(a = "p\nq", b = "t\tu", c = "x\\y"), quote = FALSE); print(matrix(c("a\"b", "c\nd")), quote = FALSE); noquote("a\tb")
#---
read.table(text = "x y\n1 2", col.names = c("a", "b"), header = TRUE)
#---
read.csv(text = "a,b\n1,2", colClasses = "numeric")
#---
read.table(text = "a b\n1 2", header = TRUE, colClasses = c(a = "character"))
#---
type.convert(c("1", "2", "3"), as.is = TRUE)
#---
type.convert(c("1.5", "2", "NA"), as.is = TRUE)
#---
type.convert(c("a", "b", "a"), as.is = FALSE)
#---
type.convert(c("a", "b", "a"), as.is = TRUE)
#---
type.convert(c("T", "F", "NA"), as.is = TRUE)
#---
type.convert(c("TRUE", "false", "NA"), as.is = TRUE)
#---
type.convert(c("T", "F"), as.is = TRUE, tryLogical = FALSE)
#---
type.convert(c("1", "2.5e3", "0x10"), as.is = TRUE)
#---
type.convert(c("2147483647", "1"), as.is = TRUE)
#---
type.convert(c("2147483648", "1"), as.is = TRUE)
#---
type.convert(c("1+2i", "3"), as.is = TRUE)
#---
type.convert(c(" 1 ", "2"), as.is = TRUE)
#---
type.convert(c("", "1"), as.is = TRUE)
#---
type.convert(c("", ""), as.is = TRUE)
#---
type.convert(c("NA", "NA"), as.is = TRUE)
#---
type.convert(character(0), as.is = TRUE)
#---
type.convert(c("1,5", "2"), dec = ",", as.is = TRUE)
#---
type.convert(c("a", "-"), na.strings = "-", as.is = TRUE)
#---
type.convert(list(a = c("1", "2"), b = c("x", "y")), as.is = TRUE)
#---
type.convert(data.frame(a = c("1", "2"), b = c("x", "y")), as.is = TRUE)
#---
type.convert(matrix(c("1", "2", "3", "4"), 2), as.is = TRUE)
#---
type.convert(c(a = "1", b = "2"), as.is = TRUE)
#---
type.convert(factor(c("1", "2")), as.is = TRUE)
#---
type.convert(1:3, as.is = TRUE)
#---
type.convert(c("1", "2"))
#---
type.convert(c("1e5", "1E-2", "1d5"), as.is = TRUE)
#---
type.convert(c("Inf", "-inf", "NaN"), as.is = TRUE)
#---
type.convert(c("0.1234567890123456789", "2"), numerals = "no.loss", as.is = TRUE)
#---
type.convert(c("0.1234567890123456789", "2"), numerals = "allow.loss", as.is = TRUE) # ulp
#---
type.convert(c("b", "a"), as.is = FALSE)
#---
type.convert(c("B", "a", "A"), as.is = FALSE)
#---
type.convert(c("1", "a"), as.is = TRUE)
#---
type.convert(list("1", 2), as.is = TRUE)
#---
type.convert("1", as.is = TRUE, na.strings = NULL)
#---
type.convert(c("1", " NA"), as.is = TRUE)
#---
type.convert(c("TRUE", "1"), as.is = TRUE)
#---
type.convert(c("1", "2", NA), as.is = TRUE)
#---
type.convert(c("-5", "+5"), as.is = TRUE)
#---
type.convert("1L", as.is = TRUE)
#---
type.convert(c(".5", "5."), as.is = TRUE)
#---
type.convert(c("007", "1"), as.is = TRUE)
#---
type.convert(NULL, as.is = TRUE)
#---
type.convert(c(TRUE, NA), as.is = TRUE)
#---
type.convert("2", as.is = TRUE, numerals = "bad")
#---
f <- tempfile(); writeLines(c("a b", "1 x", "2 y"), f); count.fields(f)
#---
f <- tempfile(); writeLines(c("a,b,c", "1,'x,y',3", "", "# c", "1,2"), f); count.fields(f, sep = ",")
#---
f <- tempfile(); writeLines(c("a,b,c", "1,'x,y',3", "", "# c", "1,2"), f); count.fields(f, sep = ",", blank.lines.skip = FALSE)
#---
f <- tempfile(); writeLines(c("a,b,c", "1,'x,y',3", "", "# c", "1,2"), f); count.fields(f, sep = ",", comment.char = "#")
#---
f <- tempfile(); writeLines(c("a,b,c", "1,'x,y',3", "", "# c", "1,2"), f); count.fields(f, sep = ",", quote = "")
#---
f <- tempfile(); writeLines(c("a,b,c", "1,'x,y',3", "", "# c", "1,2"), f); count.fields(f, sep = ",", skip = 1)
#---
f <- tempfile(); writeLines(c("a,b,c", "1,'x,y',3", "", "# c", "1,2"), f); count.fields(f, sep = ",", skip = 1, comment.char = "#", blank.lines.skip = FALSE)
#---
f <- tempfile(); writeLines(c("a \"b c\" d", "e", "f \"g", "h"), f); count.fields(f)
#---
count.fields(textConnection("1 2 3\n4 5"))
#---
f <- tempfile(); writeLines(c("a \"b c\" d", "e", "f \"g", "h"), f); count.fields(f, quote = "'")
#---
count.fields("/nonexistent")
#---
f <- tempfile(); writeLines(c("1 2 3", "x", "4 5 6 7"), f); count.fields(f, skip = 1)
#---
f <- tempfile(); writeLines(c("1 2 3", "x", "4 5 6 7"), f); count.fields(f, comment.char = "")
