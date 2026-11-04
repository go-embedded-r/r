## Replacement functions and copy on modify. A line holding only #--- separates cells.
#---
x <- 1:5; x[2] <- 10L; x
#---
x <- 1:5; x[2] <- 10; x
#---
x <- 1:5; x[2] <- "a"; x
#---
x <- 1:5; x[2] <- TRUE; x
#---
x <- c(1.5, 2.5); x[1] <- 1L; x
#---
x <- 1:5; x[c(1, 3)] <- c(100L, 300L); x
#---
x <- 1:5; x[c(1, 3)] <- 0L; x
#---
x <- 1:6; x[c(1, 3, 5)] <- c(7L, 8L); x
#---
x <- 1:5; x[7] <- 9L; x
#---
x <- 1:3; x[10] <- 1L; x
#---
x <- 1:3; x[-1] <- 0L; x
#---
x <- 1:5; x[x > 2] <- 0L; x
#---
x <- 1:5; x[c(TRUE, FALSE)] <- 0L; x
#---
x <- 1:3; x[NA] <- 1L; x
#---
x <- 1:3; x[c(1, NA)] <- 5L; x
#---
x <- 1:3; x[c(1, NA)] <- c(5L, 6L)
#---
x <- 1:3; x[1] <- integer(0)
#---
x <- 1:3; x[integer(0)] <- 1L; x
#---
x <- 1:3; x[] <- 0L; x
#---
x <- 1:3; x[1:2] <- 1:3
#---
x <- 1:3; x[0] <- 5L; x
#---
x <- c(a = 1, b = 2); x["b"] <- 20; x
#---
x <- c(a = 1, b = 2); x["c"] <- 3; x
#---
x <- c(a = 1, b = 2); x[c("a", "z")] <- c(10, 30); x
#---
x <- 1:2; x["a"] <- 5L; x
#---
x <- c(a = 1); x[3] <- 3; x
#---
x <- NULL; x[3] <- 1; x
#---
x <- NULL; x["a"] <- 1; x
#---
x <- NULL; x[1] <- "s"; x
#---
x <- NULL; x[1] <- NULL; x
#---
x <- list(1, 2, 3); x[2] <- list("a"); x
#---
x <- list(1, 2, 3); x[2] <- "a"; x
#---
x <- list(1, 2, 3); x[2] <- NULL; x
#---
x <- list(1, 2, 3); x[c(1, 3)] <- NULL; x
#---
x <- list(a = 1, b = 2); x["a"] <- NULL; x
#---
x <- list(1, 2); x[4] <- list(9); x
#---
x <- 1:3; x[2] <- list("q"); x
#---
x <- 1:3; x[2] <- NULL
#---
x <- 1:3; x[[2]] <- 10L; x
#---
x <- 1:3; x[[2]] <- "a"; x
#---
x <- 1:3; x[[5]] <- 9L; x
#---
x <- 1:3; x[[1]] <- 1:2
#---
x <- 1:3; x[[1]] <- integer(0)
#---
x <- 1:3; x[[1]] <- list(1)
#---
x <- 1:3; x[["a"]] <- 4L; x
#---
x <- list(1, 2); x[[2]] <- "b"; x
#---
x <- list(1, 2); x[[3]] <- 3; x
#---
x <- list(1, 2); x[[1]] <- NULL; x
#---
x <- list(1, 2); x[[5]] <- NULL; x
#---
x <- list(a = 1); x[["b"]] <- 2; x
#---
x <- list(a = 1, b = 2); x[["a"]] <- NULL; x
#---
x <- list(a = list(b = 1)); x[["a"]][["b"]] <- 5; x
#---
x <- list(a = list(b = 1)); x[[c("a", "b")]] <- 6; x
#---
x <- list(list(1, 2), 3); x[[c(1, 2)]] <- "z"; x
#---
x <- NULL; x[[1]] <- 5; x
#---
x <- NULL; x[["a"]] <- "q"; x
#---
x <- NULL; x[[1]] <- c(1, 2); x
#---
x <- NULL; x[[1]] <- NULL; x
#---
x <- list(1); x[[2]] <- list(2); x
#---
x <- list(a = 1, b = 2); x$a <- 10; x
#---
x <- list(a = 1); x$b <- "new"; x
#---
x <- list(a = 1, b = 2); x$a <- NULL; x
#---
x <- list(a = 1); x$z <- NULL; x
#---
x <- list(abc = 1); x$ab <- 2; x
#---
x <- NULL; x$a <- 1; x
#---
x <- NULL; x$a <- NULL; x
#---
x <- list(a = list(b = 1)); x$a$b <- 2; x$a$c <- 3; x
#---
x <- list(); x$a$b$c <- 1; x
#---
x <- 1:3; x$a <- 1; x
#---
x <- c(a = 1); x$b <- 2; x
#---
f <- function(x) x; f$a <- 1
#---
x <- 1:6; dim(x) <- c(2, 3); x[1, 2] <- 99L; x
#---
x <- 1:6; dim(x) <- c(2, 3); x[2, ] <- 0L; x
#---
x <- 1:6; dim(x) <- c(2, 3); x[, 2] <- c(7L, 8L); x
#---
x <- 1:6; dim(x) <- c(2, 3); x[, 2] <- 1:3
#---
x <- 1:6; dim(x) <- c(2, 3); x[3, 1] <- 0L
#---
x <- 1:6; dim(x) <- c(2, 3); x[[1, 3]] <- 0L; x
#---
x <- 1:6; dim(x) <- c(2, 3); x[5] <- 0L; x
#---
x <- 1:6; dim(x) <- c(2, 3); x[7] <- 0L; x
#---
x <- 1:6; dim(x) <- c(2, 3); dimnames(x) <- list(c("a", "b"), c("p", "q", "r")); x["a", "q"] <- 0L; x
#---
x <- 1:6; dim(x) <- c(2, 3); x[1, 1] <- "s"; x
#---
x <- 1:3; y <- x; y[1] <- 100L; x; y
#---
x <- list(a = 1); y <- x; y$a <- 2; x$a; y$a
#---
f <- function(v) { v[1] <- 0L; v }; x <- 1:3; f(x); x
#---
f <- function() { x[2] <- 99L; x }; x <- 1:3; f(); x
#---
x <- 1:3; f <- function() { x[2] <<- 99L; x }; f(); x
#---
x <- 1:3; f <- function() { x[[1]] <<- 0L }; f(); x
#---
x <- c(a = 1, b = 2); names(x)[2] <- "B"; x
#---
x <- list(p = 1); names(x)[1] <- "q"; x
#---
x <- 1:3; levels(x) <- "a"
#---
`second<-` <- function(x, value) { x[2] <- value; x }; x <- 1:3; second(x) <- 10L; x
#---
`mid<-` <- function(x, i, value) { x[i] <- value; x }; x <- 1:5; mid(x, 3) <- 0L; x
#---
`second<-` <- function(x, value) { x[2] <- value; x }; l <- list(v = 1:3); second(l$v) <- 99L; l
#---
`second<-` <- function(x, value) { x[2] <- value; x }; l <- list(a = list(v = 1:3)); second(l$a$v) <- 7L; l
#---
nosuch(x) <- 1
#---
x <- 1:3; nosuch(x) <- 1
#---
undefined_var[1] <- 1
#---
undefined_var$a <- 1
#---
undefined_var[[1]] <- 1
#---
names(undefined_var) <- "a"
#---
1[1] <- 2
#---
NULL[1] <- 1
#---
x <- 1:3; (x[2] <- 5L)
#---
x <- 1:3; y <- (x[2] <- 5L); y
#---
x <- 1:3; names(x) <- c("a", "b", "c"); x[["b"]] <- 20L; x
#---
a <- b <- c(1, 2); a[1] <- 100; a; b
#---
x <- 1:3; attr(x, "k") <- "v"; x[2] <- 9L; attributes(x)
#---
x <- 1:3; class(x) <- "foo"; x[2] <- 9L; class(x)
#---
x <- c(1, 2, 3); x[2] <- NA; x
#---
x <- c("a", "b"); x[2] <- NA; x
#---
x <- 1:3; x[2] <- NA; x
#---
x <- list(1, 2); x[[2]] <- NA; x
#---
i <- 2; x <- 1:3; x[i] <- 0L; x[i + 1] <- 5L; x
#---
x <- numeric(3); for (i in 1:3) x[i] <- i * 2; x
#---
x <- list(); for (i in 1:3) x[[i]] <- i; x
#---
x <- c(); for (i in 1:3) x <- c(x, i); x
#---
x <- 1:3; x[2:3] <- c("a", "b"); x

#---
x <- 0; x[[2]] <- list(7); x
#---
x <- 0; x[[1]] <- list()
#---
x <- 0; x[[1]] <- list(1, 2)
#---
x <- c(a = 1); x[["b"]] <- list(2); x
