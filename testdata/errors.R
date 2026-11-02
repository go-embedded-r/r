## Inputs that do not parse, written for this repository; R's error message for each is the golden. A line holding only #--- separates them.
#---
a b
#---
x <-
y z
#---
f(1,
#---
a +
#---
1 +* 2
#---
if (a) b
else c
#---
"abc
#---
x <- 5 5
#---
)
#---
a $ 1
#---
"\q"
#---
0x
#---
1e
#---
x <- `a
#---
a %in b
#---
a <- <- b
#---
`
#---
{
1
2
3
4 5
}
#---
é é
#---
a @
#---
function(x, x) 1
#---
1 |> 2
#---
x |> f(y = _, z = _)
#---
x |> _
#---
\
#---
a = b =
#---
a ! b
#---
5 -> 6 7
#---
x <<- y <<
#---
a < b < c
#---
else
#---
f(a=1 b)
#---
x[[1]
#---
x[1]]
#---
if x
#---
for (i 1:3) 1
#---
function(1) 1
#---
"\x"
#---
"\u{zz}"
#---
"\0"
#---
'a' 'b'
#---
x ===  y
#---
@x
#---
a ->
#---
r"(abc"
#---
\(x)
#---
a?b?
#---
1.5L2
#---
f(,,
#---
x = = 2
#---
{ a; ; b }
#---
TRUE <- 1 2
