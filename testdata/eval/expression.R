## Expression vectors: expression(), parse(text =), as.expression, eval, subsetting, printing. A line holding only #--- separates cells.
#---
expression(a+b, 2); expression(); expression(a=x+1, y); length(expression(a,b,c)); typeof(expression(a)); class(expression(a)); mode(expression(a))
#---
is.expression(expression(a)); is.expression(quote(a)); is.list(expression(a)); is.vector(expression(a)); is.atomic(expression(a)); is.recursive(expression(a)); is.language(expression(a)); is.null(expression()); inherits(expression(a), "expression"); is.function(expression(a))
#---
expression(a+b, 2)[[1]]; expression(a+b, 2)[1]; expression(a+b, 2)[2]; expression(a+b, 2)[[2]]; class(expression(a+b)[[1]]); expression(a+b, 2)[-1]; expression(a,b)[3]
#---
expression(a=1, b=2)["b"]; expression(a, b)[c(TRUE, FALSE)]; expression(a,b)[NA]; expression(a, b)[0]; attributes(expression(a=1)[1]); class(expression(a)[1])
#---
expression(a,b)[[3]]
#---
expression(a,b)[["z"]]
#---
x <- expression(a=1, b=2); x$a; x[["b"]]; names(x)
#---
eval(expression(1+2, 3*4)); eval(expression(1+2, 3*4)[1]); eval(expression(x+1), list(x=10)); eval(expression()); eval(expression(a, 5), list(a = 1))
#---
as.expression(quote(a+b)); as.expression(1:2); as.expression("a"); as.expression(list(1, quote(x))); as.expression(expression(a)); as.expression(c(a=1, b=2)); as.expression(NULL)
#---
as.list(expression(a, b+1)); as.character(expression(a, b+1)); as.vector(expression(a, 1), "list"); unlist(expression(a, 1))
#---
c(expression(a), expression(b+1)); c(expression(a), quote(b)); c(expression(a), 3); c(expression(a), list(1)); c(list(1), expression(a)); c(1, expression(a))
#---
rev(expression(a, b)); rep(expression(a), 2); names(expression(a=1, b=2)); length(expression(a, b))
#---
identical(expression(a), expression(a)); identical(expression(a), list(quote(a))); identical(expression(a, b), expression(a, c))
#---
deparse(expression(a+b, 2)); deparse(expression(a)); deparse(expression()); deparse(expression(a=x+1)); format(expression(a+b, 3)); toString(expression(a, b)); paste(expression(a+b, 3)); nchar(expression(abc))
#---
parse(text="1+2"); parse(text=c("a", "b+c")); parse(text=""); parse(text=character(0)); parse(text=NULL); parse(text=quote(a))
#---
typeof(parse(text="1")); class(parse(text="1")); length(parse(text="a; b; c")); attributes(parse(text="1+2")); attributes(expression(a))
#---
parse(text="x <- 1; y = 2"); parse(text="a\nb+c"); length(parse(text="a\n\nb"))
#---
eval(parse(text="1+2")); eval(parse(text="x <- 5; x*2")); sapply(parse(text="a; b+c"), class)
#---
parse(text="function(x) x+1"); parse(text="if (a) b else c"); parse(text="1i"); parse(text="'a'"); parse(text="x[1]"); parse(text="a %in% b"); parse(text="-x^2")
#---
parse(text="1+2", keep.source=FALSE); parse(text="1+2", n=1); parse(text="1;2;3", n=2); parse(text="1;2;3", n=0); parse(text="1;2;3", n=-1)
#---
parse(text="1 +")
#---
parse(text="a b")
#---
parse(text="x <- ")
#---
parse(text="'abc")
#---
parse(text="1 +\n2 3")
#---
str2expression("a+b"); str2lang("a+b"); class(str2lang("a+b")); str2lang("x"); class(str2lang("1")); str2expression(c("a", "b"))
#---
str2lang("a;b")
#---
is.na(expression(a, NA))
#---
lapply(expression(a, b), class); vapply(expression(a, 2), is.name, TRUE); Map(function(e) class(e), expression(a, 2)); sapply(expression(a, 2), class); for (e in expression(a, 2)) print(e)
#---
sapply(expression(a+b, c*d), function(e) as.character(e[[1]])); expression(x)[[1]] == quote(x); expression(f(x, y = 2))[[1]][[3]]; as.list(expression(f(x)))[[1]]
#---
x <- expression(a, b); x[[2]] <- quote(z); x
#---
x <- expression(a, b); x[3] <- expression(q); x
#---
x <- expression(a, b); x[[3]] <- 5; x
#---
x <- expression(a, b); x$q <- 1; x
#---
x <- expression(a, b); names(x) <- c("p", "q"); x
#---
x <- expression(a, b); x[[1]] <- NULL; x
#---
x <- expression(a, b); length(x) <- 3; x
#---
x <- expression(a, b); x[2] <- list(quote(w)); x
#---
as.numeric(expression(1, 2))
#---
as.character(expression(1, "a")); sapply(expression(1, 2), function(e) e + 1); do.call(expression, list(1, 2)); do.call("expression", list(quote(a)))
#---
print(expression(a)); structure(expression(a), class = "foo"); unclass(structure(expression(a), foo = 1)); attr(expression(a), "srcref")
#---
vector("expression", 2); is.expression(vector("expression", 0)); as.vector(1:2, "expression")
#---
mapply(function(e, i) i, expression(a, b), 1:2); Reduce(function(a, b) b, expression(1, 2, 3)); unlist(lapply(expression(a, b), as.character))
#---
quote(expression(a)); deparse(quote(expression(a, b))); eval(quote(expression(a+1, b)))
#---
f <- function(x) substitute(x); f(a + b); expression(f(a + b))[[1]]; eval(expression(f(1 + 2))[[1]]); eval(parse(text = "f(q)"))
#---
eval(expression(stop("boom")))
#---
tryCatch(eval(parse(text = "1 +")), error = function(e) "parse error"); tryCatch(eval(expression(stop("boom"))), error = function(e) conditionMessage(e))
#---
e <- expression(a * b); a <- 3; b <- 4; eval(e[[1]]); eval(e); evalq(a + b)
