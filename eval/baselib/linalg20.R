# The S3 generics of the linear algebra of slice 20; their methods are Go natives
# (eval/linalg20*.go, eval/qr20.go). Written from the help pages and checked
# against the oracle.

solve <- function(a, b, ...) UseMethod("solve")
qr <- function(x, ...) UseMethod("qr")
chol <- function(x, ...) UseMethod("chol")
kappa <- function(z, ...) UseMethod("kappa")
determinant <- function(x, logarithm = TRUE, ...) UseMethod("determinant")
kappa.qr <- function(z, ...) .kappa_tri(qr.R(z), ...)
kappa.lm <- function(z, ...) kappa.qr(z$qr, ...)
