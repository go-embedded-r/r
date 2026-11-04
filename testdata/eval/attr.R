## Attributes, names, class, dim. A line holding only #--- separates cells.
#---
x <- 1:3; names(x) <- c("a", "b", "c"); x
#---
x <- 1:3; names(x) <- c("a", "b"); x
#---
x <- 1:3; names(x) <- c("a", "b", "c", "d")
#---
x <- 1:3; names(x) <- 1:3; names(x)
#---
x <- c(a = 1, b = 2); names(x) <- NULL; x
#---
x <- c(a = 1); names(x)[1] <- "z"; x
#---
x <- 1:3; names(x)[2] <- "b"; x
#---
x <- NULL; names(x) <- "a"
#---
x <- c(a = 1, b = 2, c = 3); names(x)
#---
x <- 1:3; attr(x, "foo") <- "bar"; x
#---
x <- 1:3; attr(x, "foo") <- "bar"; attr(x, "foo")
#---
x <- 1:3; attr(x, "foo") <- "bar"; attr(x, "fo"); attr(x, "fo", exact = TRUE); attr(x, "nope")
#---
x <- 1:3; attr(x, "foo") <- 1; attr(x, "foo") <- NULL; x
#---
x <- 1:3; attr(x, "a") <- 1; attr(x, "b") <- 2; attributes(x)
#---
x <- c(n = 1); attr(x, "a") <- 1; attributes(x)
#---
attributes(1:3); attributes(NULL)
#---
x <- 1:3; attributes(x) <- list(a = 1, b = "z"); x
#---
x <- 1:3; attributes(x) <- NULL; x
#---
x <- 1:3; attr(x, "names") <- c("p", "q", "r"); x
#---
x <- 1:6; attr(x, "dim") <- c(2, 3); dim(x)
#---
x <- 1:6; dim(x) <- c(2, 3); attributes(x)
#---
x <- 1:6; dim(x) <- c(4, 2)
#---
x <- 1:6; dim(x) <- NULL; x
#---
x <- c(a = 1, b = 2, c = 3, d = 4); dim(x) <- c(2, 2); attributes(x)
#---
dim(1:3)
#---
x <- 1:3; class(x) <- "foo"; class(x)
#---
x <- 1:3; class(x) <- "foo"; x
#---
x <- 1:3; class(x) <- c("a", "b"); inherits(x, "b"); inherits(x, "z"); inherits(x, c("z", "b")); inherits(x, c("z", "b"), which = TRUE)
#---
x <- 1:3; class(x) <- "foo"; class(x) <- NULL; x
#---
class(1); class(1L); class("a"); class(TRUE); class(NULL); class(list()); class(sum); class(function(x) x); class(quote(x)); class(quote(x + 1)); class(quote(if (a) b))
#---
x <- 1:6; dim(x) <- c(2, 3); class(x); class(1:3)
#---
x <- 1:6; dim(x) <- 6; class(x)
#---
oldClass(1:3); x <- 1:3; oldClass(x) <- "bar"; oldClass(x)
#---
x <- 1:3; class(x) <- "foo"; unclass(x)
#---
unclass(1:3)
#---
x <- structure(1:3, class = "k", extra = "e"); attributes(x)
#---
structure(1:3, names = c("a", "b", "c"))
#---
structure(list(a = 1), class = "rec")
#---
structure(1:4, dim = c(2, 2))
#---
structure(1:3, foo = NULL)
#---
structure(NULL, class = "a")
#---
x <- 1:3; class(x) <- 5
#---
numeric(3); integer(2); logical(2); character(2); numeric(0); vector("list", 2); vector("numeric", 2); vector("character", 1); vector()
#---
vector("foo", 1)
#---
numeric(-1)
#---
identical(1, 1); identical(1L, 1); identical(c(a = 1), c(a = 1)); identical(list(1, "a"), list(1, "a")); identical(NULL, NULL); identical(NaN, NA_real_); identical(0, -0)
#---
identical(1:3, c(1L, 2L, 3L)); identical("a", "a"); identical(quote(x + y), quote(x + y)); identical(sum, sum)
#---
paste("a", "b"); paste("a", "b", sep = ""); paste0("a", 1:3); paste(c("a", "b"), collapse = "+"); paste("x", 1:2, c("p", "q", "r")); paste(1.5, TRUE, NA, NULL); paste(character(0), collapse = "|"); paste0()
#---
paste("a", character(0)); paste(1:3, collapse = NULL)
#---
paste(quote(x))
