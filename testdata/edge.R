## Edge cases written for this repository (Unicode, literals, newlines, pipes, precedence); parsed by R and by Go and compared. A line holding only #--- separates them.
#---
x <- "héllo wörld"; y <- "日本語"; z <- 'é\U0001F600\u{e9}'
#---
naïve <- 1; `ça va` <- 2; café.au_lait <- 3
#---
r"(raw \n string)"; R"---[a ]--- b]---"; r'{x{y}z}'; r"|a|b|"
#---
x |> f() |> g(y = 2) |> h(z = _)
#---
x |> _$a$b[1]
#---
sapply(1:3, \(i) i^2)
#---
(\(x) x + 1)(2)
#---
a ->> b; a -> b -> c; 1 -> x
#---
0x10; 0xFF; 0x1p3; 0x1.8p1; 0XAbC; 0x10L; 1e3L; 100L; 1L:3L
#---
1.5; .5; 5.; 1e-5; 1E+5; 1e5; 123456789012345678901234567890
#---
2i; 1.5i; 0x10i; 1e3i; 3L
#---
TRUE; FALSE; NA; NA_integer_; NA_real_; NA_character_; NA_complex_; Inf; NaN; NULL; T; F
#---
if (a) b else c
#---
if (a) {
  b
} else {
  c
}
#---
{
  if (a) b
  else c
}
#---
f <- function(x,
              y = 2,
              ...) {
  x + y
}
#---
x <- c(1,
       2,
       3)
#---
x[1,
  2]
#---
a +
  b *
  c
#---
a %in% b %o% c %/% d %% e
#---
-x^2; (-x)^2; -x:y; !x == y; !!x; - - x; +-x
#---
a ? b ? c
#---
?help
#---
~ x + y; y ~ x | z; ~~x
#---
x = y = z <- w
#---
f(a = , b = 1)
#---
f(,)
#---
f(x)(y)(z)
#---
x$y$z; x@y@z; x$`a b`; x$"s"; x@"s"
#---
pkg::fun(1); pkg:::hidden; "pkg"::fun; stats::`+`(1, 2)
#---
x[[1]][[2]]; x[1][2]; x[[ "a" ]]; x[]
#---
function() NULL
#---
function(a, b) a + b
#---
function(a, b = c(1, 2), ...) NULL
#---
while (TRUE) { break; next }
#---
repeat break
#---
for (i in seq_len(10)) if (i %% 2) next else print(i)
#---
`if`(a, b, c); `+`(1, 2); `[`(x, 1)
#---
"f"(1); 'g'(2)
#---
f("a" = 1, 'b' = 2, `c` = 3, NULL = 4)
#---
x <- # comment here
  5 # trailing
# whole line
#---
f(a, # one
  b) # two
#---
"multi
line string"
#---
'it''s'
#---
"tab	inside"; "\t\n\r\\\"\'\`\a\b\f\v"
#---
"\x41\101\0101"
#---
x <- 1; y <- 2;
#---
;
#---

#---
a:b:c; a:b^c; 2^3^4; 2^-3
#---
a == b; a != b; a <= b; a >= b; a < b; a > b
#---
a && b || c & d | e
#---
-1:3
#---
x[-1]; x[-(1:2)]; x[-1:-2]
#---
`my fun`(1)
#---
...; ..1; ..2; .x; ._y; .5e3
#---
a$b(c); a$b$c(d); a@b(c)
#---
function(x) -x
#---
if (a) b else if (c) d else e
#---
f(if (a) b else c)
#---
f(function(x) x, 2)
#---
a = function(x) x
#---
(a = 1)
#---
{a = 1}
#---
if (a) b = 1 else c = 2
#---
x <- if (a) 1 else 2
#---
f(x = if (a) 1 else 2)
#---
switch(a, b = , c = "bc", "other")
#---
x <- y <- z <- 0
#---
a <- b ~ c
#---
a$b <- 1; a[[1]] <- 2; a@b <- 3; names(a)[2] <- "x"
#---
f(g(h(i(j(k(1))))))
#---
1:3 |> rev() |> head(n = 2)
#---
x |> (\(y) y + 1)()
#---
x |> f(y = _, ...)
#---
%%
#---
a %% b
#---
a %+% b
#---
"a" %in% "b"
#---
x <- 1e308 * 10
#---
0.1 + 0.2; 1/3; 2/3; 1e22; 1e23; 1.7976931348623157e308; 5e-324; 4.9e-324; 2.2250738585072014e-308
#---
123456789.123456789; 0.30000000000000004; 1234567890123456789
#---
1e-400; 1e400
#---
.Machine$integer.max; 2147483647L; 2147483648L
#---
x[["a", exact = FALSE]]
#---
{}
#---
{
}
#---
(
1
)
#---
function(x) {
}
#---
if (TRUE)
  1
#---
if (TRUE) 1 else
  2
#---
f <-
  function() 1
#---
x[1]$a[[2]]@b(3)
#---
a$"b"$'c'
