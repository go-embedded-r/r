## Snippets written for this repository; R's parse + deparse of each is the golden. A line holding only #--- separates them.
#---
1 + 2 * 3
#---
(1 + 2) * 3
#---
-2^2
#---
a %in% b || !c && d
#---
x <- y <- 5
#---
x = f(a, b = 2, ...)
#---
f(,)
#---
x[1]; x[[1]]; x$name; x@slot
#---
x[i, j, drop = FALSE]
#---
function(x, y = 2, ...) x + y
#---
\(x) x^2
#---
if (a) b else c
#---
if (a) {
  b
} else if (c) {
  d
} else {
  e
}
#---
for (i in 1:10) { if (i %% 2 == 0) next; print(i) }
#---
while (TRUE) break
#---
repeat { break }
#---
x |> f(y = 1) |> g()
#---
y ~ x + z | w
#---
~ x
#---
`my var` <- 1
#---
"a\tb\n"
#---
'single "quoted"'
#---
r"(raw \string)"
#---
0x10L
#---
1e5L
#---
1e-3
#---
100000
#---
123456789012
#---
TRUE & NA
#---
NULL
#---
Inf; -Inf; NaN
#---
1i
#---
pkg::fn(x); pkg:::hidden
#---
a -> b
#---
a ->> b
#---
a <<- b
#---
x[-1]
#---
!x
#---
-(1:3)
#---
function(x) {
  # a comment
  x + 1
}
#---
switch(x, a = , b = "ab", "other")
#---
f(g(h(1, 2), 3), 4)
