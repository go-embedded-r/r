## Arithmetic: types, recycling, NA, integer overflow, the modulus family. Written for this repository. A line holding only #--- separates cells.
#---
1 + 2
#---
1L + 2L
#---
1L + 2
#---
TRUE + TRUE
#---
5 - 3; 5L - 3L
#---
2 * 3; 2L * 3L; 2L * 3
#---
7 / 2; 7L / 2L; 6L / 3L; 1 / 0; -1 / 0; 0 / 0
#---
2 ^ 10; 2L ^ 3L; 2 ^ -1; 0 ^ 0; (-8) ^ 2; 4 ^ 0.5; Inf ^ 0; 1 ^ NA; NA ^ 0
#---
7 %% 3; -7 %% 3; 7 %% -3; -7 %% -3; 5.5 %% 2; -5.5 %% 2; 7L %% 3L; -7L %% 3L; 7L %% -3L
#---
7 %/% 3; -7 %/% 3; 7 %/% -3; 5.5 %/% 2; 7L %/% 3L; -7L %/% 3L
#---
5 %% 0; 5L %% 0L; 5 %/% 0; 5L %/% 0L; -5 %/% 0; 0 %/% 0
#---
Inf %% 3; 3 %% Inf; -3 %% Inf; Inf %/% 3
#---
1:3 + 1
#---
1:3 + 1:3
#---
1:6 + 1:2
#---
1:5 + 1:2
#---
as.numeric(NULL) + 1
#---
NULL + 1
#---
NULL + 1L
#---
c(1, 2) + NULL
#---
c(a = 1, b = 2) + 1
#---
1 + c(a = 1, b = 2)
#---
c(a = 1, b = 2) + c(x = 10, y = 20)
#---
NA + 1; NA_integer_ + 1L; NA_real_ * 0; NA_integer_ * 0L; NA_real_ ^ 0
#---
.Machine_integer_max <- 2147483647L; .Machine_integer_max + 1L
#---
2147483647L + 1L
#---
-2147483647L - 1L
#---
46341L * 46341L
#---
2147483647L * 2L; 2147483647L + 0L
#---
1e308 * 10; -1e308 * 10
#---
5 %/% 2 * 2 + 5 %% 2
#---
-(1:3)
#---
-TRUE; +TRUE; -c(a = 1L, b = NA)
#---
+c(1.5, -2)
#---
1 - 2 - 3; 2 ^ 3 ^ 2; 1 + 2 * 3 - 4 / 2
#---
(1 + 2) * 3
#---
x <- 5; y <- x * 2; y + 1
#---
a <- b <- 3; a + b
#---
x = 4; x ^ 2
#---
10 -> z; z / 4
#---
5L / 2L; 5L %/% 2L; class_test <- 5L %% 2L; class_test
#---
0.5 + 0.25; 0.125 * 8; 1 / 4
#---
100000 * 100000; 123456789 * 1000
#---
TRUE & NA; FALSE & NA; TRUE | NA; FALSE | NA; NA & NA; NA | NA
#---
c(TRUE, FALSE, NA) & c(TRUE, TRUE, TRUE); c(TRUE, FALSE, NA) | c(FALSE, FALSE, FALSE)
#---
!TRUE; !c(TRUE, NA, FALSE); !0; !1:3; !NULL
#---
1:3 & c(1, 0, 1); 0 | 0; 2 & 3
#---
c(1, 2, 3) & c(TRUE, FALSE)
#---
x <- c(a = TRUE, b = FALSE); !x
#---
TRUE && FALSE; TRUE && TRUE; FALSE && 1 + "a"; TRUE || 1 + "a"; FALSE || TRUE; NA && FALSE; NA && TRUE; NA || TRUE; NA || FALSE
#---
1 && 2; 0 || 0; 0L && 1L
#---
as.logical(NULL) && TRUE
#---
c(TRUE, FALSE) && TRUE
#---
"a" && TRUE
#---
TRUE && "a"
