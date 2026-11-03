## Coercion, math and summaries. A line holding only #--- separates cells.
#---
as.integer(3.9); as.integer(-3.9); as.integer(TRUE); as.integer("12"); as.integer(" 7 "); as.integer("0x10")
#---
as.integer("a")
#---
as.integer(c("1", "x", "3"))
#---
as.integer(3e9)
#---
as.integer(c(1.5, NA, NaN, Inf))
#---
as.integer(NULL); as.integer(as.character(NULL))
#---
as.integer(c(a = 1.5, b = 2.5))
#---
as.numeric("1.5"); as.numeric("1e3"); as.numeric(" 2 "); as.numeric("Inf"); as.numeric("-inf"); as.numeric("NaN"); as.numeric("NA"); as.numeric(TRUE); as.numeric(3L)
#---
as.numeric("abc")
#---
as.double("0x1A"); as.double("1e-2"); as.double(".5"); as.double("5.")
#---
as.numeric(c(a = 1, b = 2))
#---
as.as.logical(NULL); as.logical(2); as.logical(NA); as.logical("TRUE"); as.logical("T"); as.logical("true"); as.logical("yes"); as.logical("F"); as.logical("false"); as.logical(0L); as.logical(NaN)
#---
as.logical(c(1, 0, NA))
#---
as.character(1); as.character(1.5); as.character(100000); as.character(123456); as.character(1e5 + 0.5); as.character(TRUE); as.character(10L); as.character("a"); as.character(NA); as.character(0.1)
#---
as.character(c(1.5, NA, 3)); as.character(1e15); as.character(1e-5); as.character(0.5)
#---
as.character(quote(x))
#---
as.character(NULL)
#---
abs(-3); abs(3L); abs(-3L); abs(TRUE); abs(c(-1.5, 2, NA)); abs(NA_integer_)
#---
abs(c(a = -1, b = 2))
#---
abs("a")
#---
sqrt(16); sqrt(2.25); sqrt(0); sqrt(Inf); sqrt(NA)
#---
sqrt(-1)
#---
sqrt(c(4, -4, 9))
#---
exp(0); exp(-Inf); exp(NA)
#---
floor(2.7); floor(-2.7); ceiling(2.1); ceiling(-2.1); floor(3L); ceiling(3L); floor(NA); floor(Inf)
#---
sum(1:10); sum(1, 2, 3); sum(1L, 2L); sum(1L, 2.5); sum(); sum(NULL); sum(TRUE, TRUE, NA, na.rm = TRUE); sum(c(1, NA), na.rm = TRUE)
#---
sum(c(1, NA)); sum(NA_integer_, 1L); sum(1:3, NA)
#---
sum(2147483647L, 1L)
#---
sum(1:100000)
#---
prod(1:5); prod(); prod(2, 3.5); prod(c(2, NA), na.rm = TRUE); prod(as.integer(NULL)); prod(1:20)
#---
min(3, 1, 2); max(3, 1, 2); min(1:5); max(c(1.5, 2.5)); min(2L, 1.5); max(TRUE, FALSE); min(c(1, NA), na.rm = TRUE); max(NA, 1)
#---
min(as.numeric(NULL))
#---
max(as.integer(NULL))
#---
max(NULL)
#---
min(c(2L, NA)); max(c(2L, NA), na.rm = TRUE); min(NaN, 1); max(c(NaN, NA))
#---
any(TRUE, FALSE); any(FALSE, FALSE); all(TRUE, TRUE); all(TRUE, FALSE); any(); all(); any(NA, TRUE); any(NA, FALSE); all(NA, FALSE); all(NA, TRUE); all(as.logical(NULL)); any(as.logical(NULL))
#---
any(c(0, 1)); all(1:3); any(NA, na.rm = TRUE); all(NA, na.rm = TRUE)
#---
any("a")
#---
seq_len(5); seq_len(0); seq_along(c("a", "b", "c")); seq_along(NULL); seq_along(list(1, 2))
#---
seq_len(-1)
#---
seq_len(NA)
#---
seq_len(c(1, 2))
#---
seq_len(2.9)
