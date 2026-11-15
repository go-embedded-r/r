d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.table(d)
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.table(d, quote = FALSE)
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.table(d, sep = ",", row.names = FALSE)
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.table(d, col.names = FALSE)
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.csv(d)
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.csv(d, row.names = FALSE)
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.csv2(d)
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.table(d, na = "-", qmethod = "double")
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.table(d, quote = 2)
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.table(d, quote = c(1, 2), sep = "\t")
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.table(d, eol = "|\n")
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.table(d, dec = ",")
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.table(matrix(1:4, 2))
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.table(matrix(c("a", "b", "c", "d"), 2, dimnames = list(c("r1", "r2"), c("A", "B"))))
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.table(1:3)
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.table(c(a = 1, b = 2))
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.table(d, col.names = NA)
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.table(d, row.names = c("x", "y", "z"))
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.table(d, col.names = c("A", "B", "C", "D", "E"))
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.csv(d, col.names = FALSE)
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.csv(d, sep = ";")
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.csv(d, quote = FALSE, row.names = FALSE)
#---
f <- tempfile(); write.csv(d, f); writeLines(readLines(f)); write.table(d, f, append = TRUE, col.names = FALSE); writeLines(readLines(f))
#---
f <- tempfile(); con <- file(f, "w"); write.table(d[1:2, ], con); close(con); writeLines(readLines(f))
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.table(data.frame(x = c(0.1 + 0.2, 1/3, 1e10, 1e-10, 123456789.123456789, 100000, 1e5 + 0.1, 1e15, 1e16)), row.names = FALSE)
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.table(data.frame(d = as.Date("2020-01-02") + 0:1, t = as.POSIXct("2020-01-02 03:04:05", tz = "UTC") + 0:1))
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.table(list(a = 1:2, b = c("x", "y")))
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.table(d[0, ])
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.table(1:3, col.names = FALSE, row.names = FALSE)
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.table(d, sep = "")
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.table(data.frame(a = c("x y", "p,q"), b = 1:2), sep = ",", quote = FALSE)
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.table(data.frame(a = c("x\ny")), quote = TRUE)
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.table(data.frame(a = c("x\\y", "\"")), qmethod = "escape")
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.table(data.frame(a = c("x\\y", "\"")), qmethod = "double")
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.table(data.frame(a = complex(real = 1, imaginary = 2)))
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.table("a")
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.table(d, file = "")
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.table(d, qmethod = "bad")
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.table(data.frame(a = 1:2, row.names = c("a b", "c")))
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.table(data.frame(a = c(1.5, NaN, Inf, -Inf)))
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.table(d, row.names = 1:2)
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.table(d, col.names = c("A", "B"))
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.table(matrix(c(1.5, 2, 3, 4), 2), col.names = NA)
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.table(table(c("a", "b", "a")))
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.table(as.data.frame(table(c("a", "b", "a"))))
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.table(d, quote = TRUE, sep = ";", dec = ",", na = "")
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.csv(d, na = "")
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.csv(matrix(1:4, 2))
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.csv2(data.frame(x = 1.5, y = "a"))
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.csv(data.frame(x = 1.5), eol = "\r\n")
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.csv(data.frame(x = 1.5), qmethod = "escape")
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.table(data.frame(x = as.numeric(1:3) * 1e5))
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.table(data.frame(x = c(123456.7, 0.00001234)))
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.table(data.frame(x = c(TRUE, FALSE)))
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.table(data.frame(x = 1:2, y = I(c("a", "b"))))
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.table(d, fileEncoding = "UTF-8")
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.table(d, "/nonexistent/dir/f")
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.table(data.frame(a = c("é", "日本")), quote = TRUE)
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.table(d, row.names = FALSE, col.names = TRUE, quote = FALSE, sep = "|")
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.table(d, col.names = NA, row.names = FALSE)
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.csv(d, row.names = FALSE, na = "-")
#---
d <- data.frame(n = c(1.5, 2, NA), s = c("a", "b\"c", NA), f = factor(c("u", "v", "u")), l = c(TRUE, NA, FALSE), i = c(1L, NA, 3L)); write.csv(d[2, ], row.names = "only")
