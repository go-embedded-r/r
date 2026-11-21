a <- matrix(c(2,1,1,3), 2, dimnames = list(c("r1","r2"), c("c1","c2"))); solve(a) # ulp
#---
a <- matrix(c(2,1,1,3), 2, dimnames = list(c("r1","r2"), c("c1","c2"))); solve(a, c(1, 2)) # ulp
#---
a <- matrix(c(2,1,1,3), 2, dimnames = list(c("r1","r2"), c("c1","c2"))); solve(a, c(x = 1, y = 2)) # ulp
#---
a <- matrix(c(2,1,1,3), 2, dimnames = list(c("r1","r2"), c("c1","c2"))); solve(a, matrix(1:4, 2, dimnames = list(NULL, c("u","v")))) # ulp
#---
solve(2) # ulp
#---
solve(matrix(1:4, 2)) # ulp
#---
round(solve(matrix(c(4,2,7,6), 2), matrix(c(1,0,0,1), 2)), 10) + 0 # ulp
#---
round(solve(matrix(c(4,2,1,3,5,1,2,2,6), 3)), 10)
#---
solve(matrix(c(1,2,2,4), 2))
#---
solve(matrix(c(1,2,3,4,5,6,7,8,9), 3))
#---
solve(matrix(1:6, 2))
#---
solve(matrix(c(1,1,1,1+1e-17), 2))
#---
solve(matrix(c(1,2,3,5), 2), tol = 1)
#---
solve(matrix(1:4, 2), 1:3)
#---
solve("a")
#---
solve(matrix(numeric(0), 0, 0))
#---
solve(matrix(c(2,0,0,2), 2), tol = 1e-20) # ulp
#---
round(solve(matrix(c(1+1i,2,3,4), 2)), 10)
#---
x <- matrix(c(4,2,1,3,5,1,2,2,6), 3); max(abs(x %*% solve(x) - diag(3))) < 1e-14
#---
det(matrix(c(2,1,1,3), 2)) # ulp
#---
determinant(matrix(c(2,1,1,3), 2)) # ulp
#---
determinant(matrix(c(2,1,1,3), 2), logarithm = FALSE) # ulp
#---
det(matrix(c(1,2,3,4), 2)) # ulp
#---
det(matrix(numeric(0), 0, 0))
#---
det(matrix(1:6, 2))
#---
determinant(matrix(0, 2, 2))
#---
det(matrix(c(1,2,2,4), 2))
#---
signif(det(matrix(c(3,1,2,4,5,6,7,8,10), 3)), 12) # R differs from itself by 27 ulp between macOS and Linux (exact value 2)
#---
d <- determinant(matrix(c(-3,1,2,4,5,6,7,8,10), 3)); c(signif(as.vector(d$modulus), 9), d$sign)
#---
chol(matrix(c(4,2,2,3), 2, dimnames = list(c("a","b"), c("a","b")))) # ulp
#---
chol(matrix(c(1,2,2,1), 2))
#---
chol(matrix(c(4,2,0,3), 2)) # ulp
#---
chol(matrix(1:6, 2))
#---
chol(diag(0))
#---
chol(2) # ulp
#---
chol(matrix(c(4,2,2,3), 2), pivot = TRUE) # ulp
#---
chol(matrix(c(1,1,1,1), 2), pivot = TRUE)
#---
chol(matrix(c(1,2,3,2,5,7,3,7,14), 3)) # ulp
#---
x <- matrix(c(1,2,3,2,5,7,3,7,14), 3); max(abs(crossprod(chol(x)) - x)) < 1e-14
#---
chol2inv(chol(matrix(c(4,2,2,3), 2))) # ulp
#---
round(chol2inv(chol(matrix(c(1,2,3,2,5,7,3,7,14), 3))), 10)
#---
chol2inv(matrix(c(2,0,1,0), 2))
#---
chol2inv(matrix(c(2,0,1,3), 2), size = 3)
#---
backsolve(matrix(c(2,0,1,3), 2), c(1, 2)) # ulp
#---
backsolve(matrix(c(2,0,1,3), 2), matrix(1:4, 2)) # ulp
#---
backsolve(matrix(c(2,0,1,3), 2), c(1, 2), k = 1) # ulp
#---
backsolve(matrix(c(2,0,1,3), 2), c(1, 2), transpose = TRUE) # ulp
#---
backsolve(matrix(c(2,0,1,3), 2), c(1, 2, 3))
#---
backsolve(matrix(c(2,0,1,3), 2), c(1, 2), k = 3)
#---
backsolve(matrix(c(1,0,0,0,0,0,0,0,1), 3), c(1, 2, 3))
#---
forwardsolve(matrix(c(2,1,0,3), 2), c(1, 2)) # ulp
#---
forwardsolve(matrix(c(1,1,0,0), 2), c(1, 2))
#---
forwardsolve(matrix(c(2,1,0,3), 2), c(1, 2), transpose = TRUE) # ulp
#---
norm(matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4)) # ulp
#---
norm(matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4), "I") # ulp
#---
norm(matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4), "F") # ulp
#---
norm(matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4), "E") # ulp
#---
norm(matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4), "M") # ulp
#---
round(norm(matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4), "2"), 9) # ulp
#---
norm(matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4), "o") # ulp
#---
norm(matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4), "1") # ulp
#---
norm(1:3)
#---
norm(matrix(1:4, 2), "z")
#---
norm("a")
#---
norm(matrix(c(1,NA,3,4), 2))
#---
norm(matrix(numeric(0), 0, 0))
#---
norm(matrix(c(1+1i,2,3,4), 2)) # ulp
#---
rcond(matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4)[1:3, ]) # ulp
#---
rcond(matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4)[1:3, ], "I") # ulp
#---
rcond(matrix(c(4,2,1,3,5,1,2,2,6), 3)) # ulp
#---
rcond(matrix(c(4,2,1,3,5,1,2,2,6), 3), "I") # ulp
#---
rcond(matrix(c(1,2,2,4), 2))
#---
rcond(matrix(c(2,0,1,3), 2), triangular = TRUE) # ulp
#---
rcond(1:3)
#---
rcond(matrix(c(1,2,3,4), 2), "2")
#---
kappa(diag(3))
#---
kappa(matrix(c(2,0,1,3), 2), triangular = TRUE) # ulp
#---
round(kappa(matrix(c(1,2,3,4), 2), exact = TRUE), 9) # ulp
#---
round(kappa(matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4), exact = TRUE), 9) # ulp
#---
isSymmetric(matrix(c(2,1,1,3), 2))
#---
crossprod(matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4, dimnames = list(NULL, c("A","B","C"))))
#---
tcrossprod(1:3)
#---
outer(1:2, 1:3)
#---
.kappa_tri(matrix(c(2,0,1,3), 2))
#---
chol.default(matrix(c(4,2,2,3), 2)) # ulp
#---
determinant.matrix(matrix(c(2,1,1,3), 2)) # ulp
#---
kappa.default(diag(3))
#---
kappa.qr(qr(diag(3)))
#---
kappa.lm(list(qr = qr(diag(2))))
#---
solve.default(matrix(c(2,1,1,3), 2)) # ulp
#---
signif(qr.default(matrix(c(1,2,3,5), 2))$qraux, 6)
#---
print.eigen(structure(list(values = c(2, 1), vectors = NULL), class = "eigen"))
