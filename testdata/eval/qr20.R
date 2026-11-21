x <- matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4, dimnames = list(letters[1:4], c("A","B","C"))); q <- qr(x); q$qr <- round(q$qr, 7); q$qraux <- round(q$qraux, 7); q # ulp
#---
x <- matrix(c(1,2,3,4,5,7,7,8,10), 3); q <- qr(x); list(signif(q$qr, 6), q$rank, signif(q$qraux, 6), q$pivot)
#---
x <- matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4, dimnames = list(letters[1:4], c("A","B","C"))); attributes(qr(x[, 1:2]))$class
#---
x <- matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4, dimnames = list(letters[1:4], c("A","B","C"))); round(qr.Q(qr(x)), 8) + 0
#---
x <- matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4, dimnames = list(letters[1:4], c("A","B","C"))); round(qr.R(qr(x)), 7) + 0 # ulp
#---
x <- matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4, dimnames = list(letters[1:4], c("A","B","C"))); round(qr.X(qr(x)), 8) + 0
#---
x <- matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4, dimnames = list(letters[1:4], c("A","B","C"))); round(qr.Q(qr(x), complete = TRUE), 8) + 0
#---
x <- matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4, dimnames = list(letters[1:4], c("A","B","C"))); round(qr.R(qr(x), complete = TRUE), 7) + 0 # ulp
#---
x <- matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4); round(qr.Q(qr(x), Dvec = c(1,-1,1)), 8) + 0
#---
x <- matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4, dimnames = list(letters[1:4], c("A","B","C"))); signif(qr.coef(qr(x), c(1,3,2,5)), 9)
#---
x <- matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4, dimnames = list(letters[1:4], c("A","B","C"))); signif(qr.coef(qr(x), cbind(y = c(1,3,2,5), y2 = c(2,6,4,10))), 9)
#---
x <- matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4); round(qr.qy(qr(x), c(1,3,2,5)), 8) + 0
#---
x <- matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4); round(qr.qty(qr(x), c(1,3,2,5)), 8) + 0
#---
x <- matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4, dimnames = list(letters[1:4], c("A","B","C"))); round(qr.qy(qr(x), matrix(1:8, 4, dimnames = list(letters[1:4], c("p","q")))), 8) + 0
#---
x <- matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4); round(qr.fitted(qr(x), c(1,3,2,5)), 8) # ulp
#---
x <- matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4); round(qr.resid(qr(x), c(1,3,2,5)), 8) + 0
#---
x <- matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4); round(qr.fitted(qr(x), 1:4, k = 2), 8) # ulp
#---
x <- matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4, dimnames = list(letters[1:4], c("A","B","C"))); signif(qr.solve(x, c(1,3,2,5)), 9)
#---
x <- matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4, dimnames = list(letters[1:4], c("A","B","C"))); signif(qr.solve(x[1:3, ]), 6)
#---
x <- matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4, dimnames = list(letters[1:4], c("A","B","C"))); signif(solve.qr(qr(x), c(1,3,2,5)), 9)
#---
round(qr.solve(matrix(c(1,2,3,4), 2), c(1, 1)), 10) # ulp
#---
x <- matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4); all.equal(solve(qr(x), 1:4), qr.coef(qr(x), 1:4))
#---
round(solve(qr(matrix(c(1,2,3,5), 2))), 10) + 0 # ulp
#---
x <- matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4); qr.solve(x)
#---
solve.qr(1)
#---
qr.coef(1, 1)
#---
x <- matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4); qr.coef(qr(x), 1:3)
#---
x <- matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4); qr.qty(qr(x), 1:3)
#---
is.qr(qr(matrix(1:4, 2)))
#---
is.qr(1)
#---
as.qr(qr(matrix(1:4, 2)))
#---
z <- cbind(c(1,2,3,4), c(2,4,6,8), c(10,1,0,1)); qr(z)$rank
#---
z <- cbind(c(1,2,3,4), c(2,4,6,8), c(10,1,0,1)); qr(z)$pivot
#---
z <- cbind(c(1,2,3,4), c(2,4,6,8), c(10,1,0,1)); q <- qr(z); signif(c(q$qr[1, ], q$qr[2, 1:2]), 6)
#---
z <- cbind(c(1,2,3,4), c(2,4,6,8), c(10,1,0,1)); signif(qr.coef(qr(z), c(1,3,2,5)), 9)
#---
z <- cbind(c(1,2,3,4), c(2,4,6,8), c(10,1,0,1)); qr.solve(z, c(1,3,2,5))
#---
z <- cbind(c(1,2,3,4), c(2,4,6,8), c(10,1,0,1)); round(qr.X(qr(z)), 8) + 0
#---
x <- matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4, dimnames = list(letters[1:4], c("A","B","C"))); q <- qr(x, LAPACK = TRUE); q$qr <- round(q$qr, 7); q$qraux <- round(q$qraux, 7); q # ulp
#---
x <- matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4, dimnames = list(letters[1:4], c("A","B","C"))); signif(qr.coef(qr(x, LAPACK = TRUE), c(1,3,2,5)), 9)
#---
x <- matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4, dimnames = list(letters[1:4], c("A","B","C"))); round(qr.Q(qr(x, LAPACK = TRUE)), 8) + 0
#---
x <- matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4, dimnames = list(letters[1:4], c("A","B","C"))); round(qr.R(qr(x, LAPACK = TRUE)), 7) + 0 # ulp
#---
x <- matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4, dimnames = list(letters[1:4], c("A","B","C"))); all.equal(unname(solve(qr(x, LAPACK = TRUE), c(1,3,2,5))), unname(qr.coef(qr(x), c(1,3,2,5))))
#---
z <- cbind(c(1,2,3,4), c(2,4,6,8), c(10,1,0,1)); qr(z, LAPACK = TRUE)$rank
#---
z <- cbind(c(1,2,3,4), c(2,4,6,8), c(10,1,0,1)); qr(z, LAPACK = TRUE)$pivot
#---
x <- matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4); qr.fitted(qr(x, LAPACK = TRUE), 1:4)
#---
q <- qr(1:3); q$qr <- round(q$qr, 7); q$qraux <- round(q$qraux, 7); q # ulp
#---
qr(matrix(0, 2, 2))
#---
qr("a")
#---
qr(matrix(numeric(0), 0, 3))
#---
x <- matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4, dimnames = list(letters[1:4], c("A","B","C"))); q <- qr(x, tol = 1e-1); q$qr <- round(q$qr, 7); q$qraux <- round(q$qraux, 7); q # ulp
#---
qr(matrix(c(1,NA,3,4), 2))
#---
x <- matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4); q <- qr(x); max(abs(qr.Q(q) %*% qr.R(q) - x)) < 1e-13
#---
x <- matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4); q <- qr(x); max(abs(crossprod(qr.Q(q)) - diag(3))) < 1e-13
#---
x <- matrix(c(1,2,3,4,5,6,7,8,10,1,0,1), 4); q <- qr(x, LAPACK = TRUE); max(abs(qr.Q(q) %*% qr.R(q) - x[, q$pivot])) < 1e-13
