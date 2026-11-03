## Comparison and sequences. A line holding only #--- separates cells.
#---
1 == 1; 1 != 1; 1 < 2; 1 <= 1; 2 > 1; 2 >= 3
#---
1L == 1; 1L < 2.5; TRUE == 1L; TRUE > FALSE
#---
c(1, 2, 3) == c(1, 5, 3)
#---
1:6 > 3
#---
1:5 < 1:2
#---
NA == 1; NA_integer_ > 0L; NaN == NaN; NA_real_ < 1
#---
c(a = 1, b = 5) > 2
#---
"a" == "a"; "a" < "b"; "b" < "a"; "abc" != "abd"; "a" <= "a"
#---
"1" == 1; 1 == "1"; TRUE == "TRUE"; 1L == "1"; "a" == NA
#---
NULL == 1; 1 == NULL; as.logical(NULL) == 1
#---
as.numeric(NULL) < 1
#---
1 == c(1, 2, 1)
#---
Inf > 1e308; -Inf < 0; Inf == Inf
#---
1:3 == c(1L, 2L)
#---
1:10
#---
5:1
#---
-2:2
#---
1.5:4
#---
1:1
#---
3:3.5
#---
0:0
#---
(-2):2; 2:-2
#---
1:0
#---
TRUE:3
#---
1:NA
#---
1:NULL
#---
c(5, 6):7
#---
1:3.9
#---
2147483647L:2147483647L
#---
1e10:1e10
#---
c(1, 2, 3)
#---
c(1L, 2L); c(1L, 2.5); c(TRUE, 2L); c(TRUE, "a"); c(1, "a", TRUE)
#---
c()
#---
c(NULL, NULL)
#---
c(1, NULL, 2)
#---
c(a = 1, b = 2)
#---
c(a = 1, 2, c = 3)
#---
c(x = c(1, 2), y = 3)
#---
c(x = c(p = 1, q = 2))
#---
c(c(a = 1), b = 2)
#---
c(list(1, 2), 3)
#---
c(list(a = 1), b = 2)
#---
c(1:2, list("x"))
#---
c("a", NA, "b"); c(NA, 1L); c(NA, "a"); c(NA_real_, 1)
#---
c(1.5, 2L, TRUE)
#---
c(1e5, 1)
#---
c(0.1, "x"); c(1e5, "x"); c(100000L, "x"); c(1e15, "x"); c(123456.7, "x"); c(1/3, "x")
#---
length(1:10); length(NULL); length(list(1, 2)); length("a"); length(c(a = 1))
#---
length(c(1, 2, 3)) + 1L
#---
list(1, "a")
#---
list(a = 1, b = "x")
#---
list(1, b = 2)
#---
list()
#---
list(list(1), list(2, 3))
#---
list(NULL)
#---
names(c(a = 1, b = 2)); names(1:3); names(list(x = 1))
#---
is.na(c(1, NA, 3)); is.na(NA); is.na(NaN); is.na("a"); is.na(c("a", NA)); is.na(NULL)
#---
is.na(c(a = 1, b = NA))
#---
is.na(list(1, NA, c(1, 2), "a", NA_character_))
#---
is.null(NULL); is.null(1); is.null(list()); is.null(c())
#---
is.numeric(1); is.numeric(1L); is.numeric("a"); is.numeric(TRUE); is.numeric(NULL)
#---
is.character("a"); is.character(1); is.logical(TRUE); is.logical(NA); is.integer(1L); is.integer(1); is.double(1); is.double(1L)
#---
is.list(list()); is.list(1); is.atomic(1); is.atomic(list()); is.atomic(NULL)
#---
is.function(sum); is.function(function(x) x); is.function(1)
#---
is.symbol(quote(x)); is.name(quote(x)); is.call(quote(f(x))); is.call(quote(x)); is.symbol(1)
#---
is.environment(1); is.environment(NULL); is.environment(sum)
