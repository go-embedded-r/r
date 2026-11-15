rs <- function(o) { if (is.list(o)) { for (i in seq_along(o)) o[[i]] <- rs(o[[i]]); o } else if (is.double(o)) signif(o, 9) else o }; x <- matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4, dimnames = list(letters[1:4], c("A","B","C"))); rs(svd(x))
#---
x <- matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4); signif(svd(x)$d, 9)
#---
rs <- function(o) { if (is.list(o)) { for (i in seq_along(o)) o[[i]] <- rs(o[[i]]); o } else if (is.double(o)) signif(o, 9) else o }; x <- matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4); rs(svd(x, nu = 0, nv = 0))
#---
rs <- function(o) { if (is.list(o)) { for (i in seq_along(o)) o[[i]] <- rs(o[[i]]); o } else if (is.double(o)) signif(o, 9) else o }; x <- matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4); rs(svd(x, nu = 4))
#---
rs <- function(o) { if (is.list(o)) { for (i in seq_along(o)) o[[i]] <- rs(o[[i]]); o } else if (is.double(o)) signif(o, 9) else o }; x <- matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4); rs(svd(x, nu = 2, nv = 2))
#---
rs <- function(o) { if (is.list(o)) { for (i in seq_along(o)) o[[i]] <- rs(o[[i]]); o } else if (is.double(o)) signif(o, 9) else o }; x <- matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4); rs(svd(x, nv = 1))
#---
rs <- function(o) { if (is.list(o)) { for (i in seq_along(o)) o[[i]] <- rs(o[[i]]); o } else if (is.double(o)) signif(o, 9) else o }; x <- matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4); rs(La.svd(x))
#---
x <- matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4); signif(La.svd(x, nu = 0, nv = 0)$d, 9)
#---
rs <- function(o) { if (is.list(o)) { for (i in seq_along(o)) o[[i]] <- rs(o[[i]]); o } else if (is.double(o)) signif(o, 9) else o }; x <- matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4); rs(La.svd(x, nu = 0, nv = 2))
#---
rs <- function(o) { if (is.list(o)) { for (i in seq_along(o)) o[[i]] <- rs(o[[i]]); o } else if (is.double(o)) signif(o, 9) else o }; x <- matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4); rs(La.svd(x, nu = 5))
#---
x <- matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4); s <- svd(x); max(abs(s$u %*% diag(s$d) %*% t(s$v) - x)) < 1e-13
#---
x <- matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4); s <- svd(t(x)); max(abs(s$u %*% diag(s$d) %*% t(s$v) - t(x))) < 1e-13
#---
x <- matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4); s <- svd(x); max(abs(crossprod(s$u) - diag(3))) < 1e-13
#---
x <- matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4); s <- svd(t(x)); signif(abs(s$u), 9)
#---
signif(svd(1:3)$d, 9)
#---
svd(matrix(0, 2, 2))$d
#---
svd(matrix(c(1,NA,3,4), 2))
#---
svd(matrix(numeric(0), 0, 3))
#---
svd(matrix(1:4, 2), nu = -1)
#---
s <- matrix(c(2,1,0,1,3,1,0,1,4), 3, dimnames = list(c("a","b","c"), NULL)); signif(eigen(s)$values, 9)
#---
s <- matrix(c(2,1,0,1,3,1,0,1,4), 3); signif(abs(eigen(s)$vectors), 9)
#---
s <- matrix(c(2,1,0,1,3,1,0,1,4), 3); e <- eigen(s); max(abs(s %*% e$vectors - e$vectors %*% diag(e$values))) < 1e-13
#---
s <- matrix(c(2,1,0,1,3,1,0,1,4), 3); signif(eigen(s, only.values = TRUE)$values, 9)
#---
s <- matrix(c(2,1,0,1,3,1,0,1,4), 3); eigen(s, only.values = TRUE)$vectors
#---
s <- matrix(c(2,1,0,1,3,1,0,1,4), 3); class(eigen(s))
#---
s <- matrix(c(2,1,0,1,3,1,0,1,4), 3); class(eigen(s, only.values = TRUE))
#---
rs <- function(o) { if (is.list(o)) { for (i in seq_along(o)) o[[i]] <- rs(o[[i]]); o } else if (is.double(o)) signif(o, 9) else o }; rs(eigen(matrix(c(0,1,-1,0), 2)))$values
#---
g <- matrix(c(1,2,3,4,5,6,7,8,10), 3); signif(eigen(g)$values, 9)
#---
g <- matrix(c(1,2,3,4,5,6,7,8,10), 3); signif(eigen(g, symmetric = FALSE, only.values = TRUE)$values, 9)
#---
g <- matrix(c(1,2,3,4,5,6,7,8,10), 3); e <- eigen(g); max(abs(g %*% e$vectors - e$vectors %*% diag(e$values))) < 1e-12
#---
g <- matrix(c(1,2,3,4,5,6,7,8,10), 3); signif(eigen(g)$vectors, 9)
#---
eigen(matrix(c(2,0,0,3), 2))$values
#---
eigen(matrix(c(1,2,3,4), 2), symmetric = TRUE)$values
#---
signif(eigen(matrix(c(5,4,1,2), 2))$vectors, 9)
#---
signif(eigen(matrix(c(4,1,1,3), 2))$values, 9)
#---
eigen(matrix(1:6, 2))
#---
eigen(matrix(c(NA,1,1,1), 2))
#---
eigen(matrix(numeric(0), 0, 0))
#---
e <- eigen(matrix(c(1+1i,2,2,3), 2)); class(e$values)
#---
h <- matrix(c(2,1-1i,1+1i,3), 2); e <- eigen(h); max(Mod(h %*% e$vectors - e$vectors %*% diag(e$values))) < 1e-13
#---
h <- matrix(c(2,1-1i,1+1i,3), 2); signif(eigen(h)$values, 9)
