str(1:100)
#---
str(as.numeric(1:100))
#---
str(c(1.2345, 2))
#---
str(c(1.5, NA, 3))
#---
str(c(1234567.891, 0.00012))
#---
str(seq(0.5, 50, by = 0.5))
#---
str(rep(c("alpha", "beta"), 20))
#---
str(c(TRUE, NA, FALSE))
#---
str(rep(TRUE, 30))
#---
str(factor(c("u", "v", "u", "w", NA)))
#---
str(factor(letters))
#---
str(factor(c("a", "b"), levels = c("b", "a")))
#---
str(ordered(c("lo", "hi"), levels = c("lo", "hi")))
#---
str(as.Date("2024-01-01") + 0:20)
#---
str(as.POSIXct("2024-01-01 10:00:00", tz = "UTC"))
#---
str(array(1:24, c(2, 3, 4)))
#---
str(matrix(c("a", "b", "c", "d"), 2))
#---
str(3L)
#---
str(2.5)
#---
str(NA)
#---
str(c(2+3i))
#---
str(as.raw(1:3))
#---
str(character(0))
#---
str(integer(0))
#---
str(list())
#---
str("a long string that goes on and on and on and on and on and on and on and on and on and on and on")
#---
str(1e10)
#---
str(c(0.1, 0.123456789))
#---
str(100000L)
#---
str(123456789)
#---
str(c(1e-5, 1))
#---
str(-1.5)
#---
str(1:3, vec.len = 1)
#---
str(letters, vec.len = 2)
#---
str(as.numeric(1:20), digits.d = 5)
#---
str(pi, digits.d = 10)
#---
str(structure(1:3, class = "foo"))
#---
str(quote(a + b))
#---
str(as.name("x"))
#---
str(1:3, give.attr = FALSE)
#---
str(-5:5)
#---
str(c(10.5, 200.25, 3000.125))
#---
str(NULL)
#---
str(letters)
#---
str(1:3, give.head = FALSE)
#---
str(TRUE)
#---
str(Sys.Date() > as.Date("2000-01-01"))
#---
str(matrix(1:6, 2))
#---
str(c(a = 1, b = 2), give.attr = FALSE)
#---
x <- str(1:3); is.null(x)
#---
invisible(str(1:3))
