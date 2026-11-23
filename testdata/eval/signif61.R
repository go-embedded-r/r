## signif of a subnormal, where R agrees with the rounding of the exact binary value (issue #61). The numbers are k * 2^-1074, built exactly. A line holding only #--- separates cells.
#---
u <- .Machine$double.xmin / 2^52; sprintf("%.17g", signif(u * c(1, 2, 3, 7, 10, 99, 100, 101, 1000, 4096, 123456), 1))
#---
u <- .Machine$double.xmin / 2^52; sprintf("%.17g", signif(u * c(1, 2, 3, 7, 10, 99, 100, 101, 1000, 4096, 123456), 2))
#---
u <- .Machine$double.xmin / 2^52; sprintf("%.17g", signif(u * c(1, 2, 3, 7, 10, 99, 100, 101, 1000, 4096, 123456), 3))
#---
u <- .Machine$double.xmin / 2^52; sprintf("%.17g", signif(u * c(1, 2, 3, 7, 10, 99, 100, 101, 1000, 4096, 123456), 5))
#---
u <- .Machine$double.xmin / 2^52; sprintf("%.17g", signif(u * c(1, 2, 3, 7, 10, 99, 100, 101, 1000, 4096, 123456), 8))
#---
u <- .Machine$double.xmin / 2^52; sprintf("%.17g", signif(u * c(1, 2, 3, 7, 10, 99, 100, 101, 1000, 4096, 123456), 12))
#---
u <- .Machine$double.xmin / 2^52; sprintf("%.17g", signif(u * c(1, 2, 3, 7, 10, 99, 100, 101, 1000, 4096, 123456), 17))
#---
u <- .Machine$double.xmin / 2^52; sprintf("%.17g", signif(-u * c(3, 99, 123456, 7654321), 3))
#---
u <- .Machine$double.xmin / 2^52; sprintf("%.17g", signif(u * c(2^20 + 12345, 2^30 + 987654321, 2^36 + 5), 4))
#---
u <- .Machine$double.xmin / 2^52; sprintf("%.17g", signif(u * c(2^20 + 12345, 2^30 + 987654321, 2^36 + 5), 9))
#---
u <- .Machine$double.xmin / 2^52; signif(u * 999, 1) == signif(u * 1000, 1)
#---
u <- .Machine$double.xmin / 2^52; signif(u * 3, 22) == u * 3
