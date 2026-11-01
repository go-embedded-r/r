## Expressions written for this repository; each is evaluated by R and its serialization recorded. A line holding only #--- separates them. Literals must be exact in double precision: R's own parser rounds 1e300 differently on macOS arm64 and Linux x86-64 (R 4.6.2), which the CI oracle lane found.
#---
NULL
#---
TRUE
#---
c(TRUE, NA, FALSE)
#---
42L
#---
c(1L, NA, -3L)
#---
1:10
#---
3:1
#---
seq_len(0)
#---
1.5
#---
c(0.1, NA, NaN, Inf, -Inf, 1e15, -0)
#---
NA_real_
#---
1 + 2i
#---
c(1i, NA)
#---
"a"
#---
c("héllo", NA, "", "日本語")
#---
as.character(1:5)
#---
as.raw(c(0, 1, 255))
#---
list(1, "a", NULL)
#---
list(a = 1, b = list(c = "x"))
#---
c(a = 1, b = 2)
#---
matrix(1:6, 2, 3)
#---
matrix(c(1.5, 2.5, 3.5, 4.5), 2, dimnames = list(c("r1", "r2"), c("c1", "c2")))
#---
array(0, c(2, 2, 2))
#---
factor(c("lo", "hi", "lo"), levels = c("lo", "hi"))
#---
structure(list(x = 1:3, y = c("a", "b", "c")), class = "data.frame", row.names = c(NA, -3L))
#---
data.frame(n = 1:2, s = c("u", "v"))
#---
structure(1L, class = "myclass", note = "kept")
#---
as.Date("2026-11-01")
#---
quote(x)
#---
quote(f(x, y = 1))
#---
quote(if (a) b else c)
#---
expression(a + b, sin(x))
#---
function(x, y = 2) x + y
#---
globalenv()
#---
emptyenv()
#---
baseenv()
#---
asNamespace("stats")
#---
list(globalenv(), globalenv())
#---
local({ e <- new.env(); assign("v", 1:3, e); list(e, e) })
#---
local({ e <- new.env(); e$self <- e; e })
#---
sum
#---
`if`
#---
pairlist(a = 1, 2)
#---
as.name("a b")
#---
list(quote(x), quote(x), quote(y))
#---
numeric(0)
#---
character(0)
#---
list()
#---
c(1, 2, 3) + 0
#---
rev(1:5)
#---
as.numeric(1:3)
#---
1:3 + 0L
#---
tryCatch(stop("boom"), error = function(e) e)
