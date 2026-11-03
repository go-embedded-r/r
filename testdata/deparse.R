## Inputs chosen to exercise deparse(): parentheses, precedence, line breaking, indentation, constants, names. Written for this repository. A line holding only #--- separates them.
#---
(a + b) * c; a + (b * c); (a * b) + c; a - (b - c); (a - b) - c; a / (b / c); (a^b)^c; a^(b^c)
#---
-(1 + 2); (-1)^2; -1^2; -(-x); !(a && b); (!a) && b; !a == b; (!a) == b
#---
(a <- 1) + 2; a <- (b <- 2); (a = 1); f((a = 1)); f(a = (b = 1)); x <- (y = 2)
#---
a %in% b; (a %in% b) %in% c; a %in% (b %in% c); a + b %in% c; (a + b) %in% c; `%+%`(a, b); `%+%`(a, b = 1)
#---
a < b; (a < b) < c; (a == b) == c; a < (b < c)
#---
x$a$b; (x$a)$b; x$(a); x$"a b"; x$`a b`; x@a; (a + b)$c; f(x)$y; x[1]$a; (x[1])$a; x$a[1]; (x$a)[1]
#---
x[1]; x[[1]]; x[1, 2]; x[, 2]; x[1, ]; x[]; x[[1]][[2]]; (a + b)[1]; (-a)[1]; (function(x) x)[1]
#---
f(x)(y); (f + g)(x); (function(x) x)(1); (\(x) x + 1)(2); x$f(1); x[[1]](2); (x)(1); `f g`(1); "f"(1)
#---
function(x) x; function(x, y = 2, ...) NULL; function() 1; function(`a b` = 1) 1; function(x) function(y) x + y
#---
if (a) b; if (a) b else c; if (a) b else if (c) d else e; (if (a) b else c) + 1; 1 + if (a) b else c; -if (a) b
#---
for (i in x) y; while (a) b; repeat break; (for (i in x) y); function(x) for (i in x) i
#---
{
  a
  b
}
#---
{}
#---
{
  if (a) b
  c
}
#---
{
  if (a) b else c
  if (a) {
    b
  } else {
    c
  }
  if (a) {
    b
  }
  else if (d) e
}
#---
f <- function(x) {
  if (x) {
    1
  } else if (y) {
    2
  } else {
    3
  }
}
#---
f <- function(x, y) {
  for (i in 1:10) {
    while (TRUE) {
      repeat {
        break
      }
    }
  }
}
#---
f(if (a) b else c, d)
#---
{ f(if (a) b else c) ; g(function(x) { if (x) 1 else 2 }) }
#---
c(if (a) b, 1)
#---
list(a = if (a) b else c)
#---
x <- if (a) b else c
#---
{ x <- if (a) b else c; y <- 1 }
#---
aaaaaaaaaaaaaaaaaaaaaa + bbbbbbbbbbbbbbbbbbbbbbbbbb + cccccccccccccccccccccc + dddddddddddddddddddd + eeeeeeeeeeeeee
#---
aaaaaaaaaaaaaaaaaaaaaa(bbbbbbbbbbbbbbbbbbbbbbbbbb, cccccccccccccccccccccc, dddddddddddddddddddd, eeeeeeeeeeeeee, ffffffffff)
#---
f(aaaaaaaaaaaaaaaaaaaaaa = bbbbbbbbbbbbbbbbbbbbbbbbbb, ccccccccccccccccccccc = dddddddddddddddddddd, eeeeeeeeeeeeee = ffffffffff, g = 1)
#---
function(aaaaaaaaaaaaaaaaaaaaaa = 1, bbbbbbbbbbbbbbbbbbbbbbbbbb = 2, cccccccccccccccccccccc = 3, dddddddddddddddddddd = 4) NULL
#---
x[aaaaaaaaaaaaaaaaaaaaaa, bbbbbbbbbbbbbbbbbbbbbbbbbb, cccccccccccccccccccccc, dddddddddddddddddddd, eeeeeeeeeeeeee, ]
#---
{
  aaaaaaaaaaaaaaaaaaaaaa <- bbbbbbbbbbbbbbbbbbbbbbbbbb(cccccccccccccccccccccc, dddddddddddddddddddd, eeeeeeeeeeeeee(fffffffffffff, gggggggggggg))
}
#---
a && b || c && d || eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee && ffffffffffffffffffffffffffff || gggggggggggggggggggg
#---
1; 1L; 1.5; 1e5; 1e-5; 100000; 123456; 0.1; 1/3; 1e15; 1e16
#---
2i; 1.5i; 0x10; 0xFFFFFFFFFF; 1e5L; 0x7FFFFFFFL; 100L
#---
TRUE; FALSE; NA; NA_integer_; NA_real_; NA_character_; NA_complex_; Inf; NaN; NULL
#---
"a"; 'b'; "it's"; 'say "hi"'; "tab\there"; "nl\nhere"; "back\\slash"; "é"; "日本語"; "​"; "\x01"; "\001"; "\a\b\f\v\r"
#---
`a b`; `if`; `a b` <- 1; f(`a b` = 1); x$`a b`; `+`(1, 2); `-`(1); `[`(x, 1); `{`(1, 2); `(`(1); `if`(a, b)
#---
`-`(1, 2, 3); `+`(); `!`(a, b); `!`(); `[`(); `$`(a); `$`(a, 1); `@`(a, "b"); `::`(a, b); `::`("a", "b")
#---
a::b; a:::b; a::b(1); a::b$c; a$b::c
#---
~x; y ~ x; ~x + y; (~x) + y; y ~ x | z; a ~ b ~ c; ~~x
#---
?help; `?`(a, b); a ? b
#---
a:b; (a:b):c; a:(b:c); -a:b; (-a):b; a:-b; 1:3 + 1; 1:(3 + 1); (1:3)[2]
#---
a && b; a || b; a & b; a | b; (a || b) && c; a || (b && c)
#---
x |> f(); x |> f(y = 1) |> g(); x |> f(a = _)
#---
f(, ); f(a, ); f(, a); f(a = , b = 1); f(a = )
#---
f(...); f(..1); f(...)(2); function(...) list(...); function(x, ...) NULL
#---
a$b <- 1; a[[1]] <- 2; a[1] <- 3; a@b <- 4; names(a)[2] <- "x"; `names<-`(a, b); a <<- 1; 1 -> a; 1 ->> a
#---
switch(a, b = , c = 1, 2)
#---
(function(x, y) x + y)(1, 2)
#---
f(g(h(i(j(k(l(m(n(o(p(q(r(s(t(u(v(w(x(y(z(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13)))))))))))))))))))))
#---
function(x) {
  # comment
  x + 1  # trailing
}
#---
local({
  a <- 1
  b <- function() {
    a + 1
  }
  b()
})
#---
tryCatch({
  stop("x")
}, error = function(e) {
  conditionMessage(e)
}, finally = {
  cat("done\n")
})
#---
x <- c(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30)
#---
paste("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", "cccccccccccccccccccccccccccccccc")
#---
if (a) b else c -> d
#---
function(x) if (x) 1 else 2
#---
a <- function() NULL; b <- function(x) -x; d <- function(x = c(1, 2)) x
#---
`my var`$`other var`[[1]]
#---
(1)(2); (1 + 2)(3); (a$b)(c); a$b(c); (a[1])(2)
#---
if (a) (b) else (c)
#---
if (a) -1 else +1
#---
x[[1]][1]$a@b
#---
-x$a; (-x)$a; -x[1]; !x$a; !x[1]; -x^2; (-x)^2; x^-2; x^(-2); 2^-x
#---
a <- b <- c; (a <- b) <- c; a <- (b <- c)
#---
a = b = c
#---
a ~ b = c
