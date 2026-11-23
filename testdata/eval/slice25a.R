print(1.837924e-321); print(-6.265e-321); print(4.9406564584124654e-324)
#---
print(c(1.837924e-321, 1)); print(c(2.5e-320, 1e-310, 3)); format(2.5e-320); format(4.9406564584124654e-324, digits = 3)
#---
x <- 2^-1074 * c(1, 3, 1000, 123456789); print(x); print(signif(x, 3)); cat(x, "\n"); as.character(signif(x, 3))
#---
print(2.2250738585072014e-308); print(2.225073858507201e-308); print(c(1e-320, 1e-323)); print(.Machine$double.xmin)
