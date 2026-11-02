## Numeric literals whose value depends on the platform's floating-point arithmetic in R (80-bit long double on x86-64 Linux, fused multiply-add on arm64): compared by `rparity snippets` within 16 ulp, kept out of the goldens. (Not here: 1.7976931348623157e308, DBL_MAX, which 64-bit-double R rounds to Inf and 80-bit R keeps.) A line holding only #--- separates them.
#---
1.5; .5; 5.; 1e-5; 1E+5; 1e5; 123456789012345678901234567890
#---
0.1 + 0.2; 1/3; 2/3; 1e22; 1e23; 5e-324; 4.9e-324; 2.2250738585072014e-308
#---
123456789.123456789; 0.30000000000000004; 1234567890123456789
#---
1e-400; 1e400
#---
x <- 1e308 * 10
