isTRUE(TRUE); isTRUE(c(TRUE, TRUE)); isTRUE(NA); isTRUE("a"); isFALSE(FALSE); isFALSE(NA); isTRUE(c(a = TRUE))
#---
xor(TRUE, FALSE); xor(c(TRUE, TRUE), c(TRUE, FALSE)); xor(NA, TRUE)
#---
f <- Negate(is.na); f(c(1, NA)); Negate("isTRUE")(FALSE)
#---
setNames(1:3, c("a", "b", "c")); setNames(nm = c("x", "y")); setNames(list(1, 2), c("p", "q")); setNames(1:2, NULL)
#---
is.element(2, 1:3); is.element(c(1, 5), 1:3)
#---
union(c(1, 2, 3), c(2, 5)); union(c("a", "b"), "b"); union(1:3, NULL); intersect(1:5, 4:8); intersect(c("a", "b", "a"), c("a")); setdiff(1:5, 2:3); setdiff(c("a", "b"), character(0)); setdiff(c(1, 1, 2), 3)
#---
append(1:3, 99); append(1:3, 99, after = 1); append(1:3, 99, after = 0); append(list(1), list(2)); append(1:3, 4:5, after = 2)
#---
head(1:10); head(1:10, 3); head(1:10, -2); head(letters, 0); head(1:3, 10); head(matrix(1:20, 5), 2); head(c(a = 1, b = 2, c = 3), 2)
#---
tail(1:10); tail(1:10, 3); tail(1:10, -7); tail(letters, 2); tail(1:3, 10); tail(matrix(1:20, 5), 2); tail(1:5, 0)
#---
ifelse(c(TRUE, FALSE, NA), "y", "n"); ifelse(1:6 > 3, 1:6, -1); ifelse(c(a = TRUE, b = FALSE), 1, 2); ifelse(matrix(c(TRUE, FALSE, FALSE, TRUE), 2), 1, 0); ifelse(logical(0), 1, 2); ifelse(c(1, 0, 2), "t", "f"); ifelse(TRUE, c(1, 2), 3)
#---
nrow(matrix(1:6, 2)); ncol(matrix(1:6, 2)); nrow(1:3); NROW(1:3); NCOL(1:3); NROW(matrix(1:6, 2)); NCOL(matrix(1:6, 2))
#---
lengths(list(1:3, "a", NULL)); lengths(list(a = 1:2, b = 1)); lengths(1:3); lengths(list())
#---
rep_len(1:3, 7); rep_len("a", 2)
#---
diff(c(1, 4, 9, 16)); diff(1:10, lag = 2); diff(c(1, 4, 9, 16), differences = 2); diff(5); diff(c(a = 1, b = 3))
#---
as.vector(matrix(1:4, 2)); as.vector(c(a = 1, b = 2)); as.vector(list(a = 1, b = 2)); as.vector(factor(c("x", "y"))); as.vector(1:3, "character"); as.vector("1", "numeric"); as.vector(c(1.7, 2), "integer"); as.vector(1:2, "list"); as.vector(NULL, "list"); as.vector(c(a = 1), "list"); as.vector(0:1, "logical")
#---
as.vector(list(1, 2), "numeric")
#---
as.vector(1, "bogus")
#---
as.vector(sum)
#---
order(c(3, 1, 2)); order(c(3, 1, 2), decreasing = TRUE); order(c("b", "A", "a", "B")); order(c(2, 1, NA, 3)); order(c(2, 1, NA, 3), na.last = FALSE); order(c(2, 1, NA, 3), na.last = NA); order(c(1, 1, 2), c(3, 2, 1)); order(c(1, 1, 2), c(3, 2, 1), decreasing = TRUE); order(c(TRUE, FALSE)); order(factor(c("b", "a"))); order(integer(0)); order(c(1, 1, 1))
#---
order(list(1, 2))
#---
order(1:3, 1:2)
#---
x <- c(5, 2, 9); x[order(x)]; sort(x, decreasing = TRUE)
#---
pmax(c(1, 5, 3), c(4, 2, 6)); pmin(c(1, 5, 3), 2); pmax(c(1, NA, 3), 2); pmax(c(1, NA, 3), 2, na.rm = TRUE); pmax(1:3, c(2.5, 1, 1)); pmax("a", c("b", "A")); pmin(c(a = 1, b = 5), 3); pmax(integer(0), 1); pmin(1:6, 1:2)
#---
pmax(1:3, 1:2)
#---
pmax()
#---
which.max(c(1, 5, 3, 5)); which.min(c(4, 1, 1)); which.max(c(a = 1, b = 9)); which.max(c(NA, 2, NA)); which.max(numeric(0)); which.min(c(TRUE, FALSE)); which.max(NULL); which.max("3")
#---
tabulate(c(1, 2, 2, 5)); tabulate(c(1, 2, 2), nbins = 4); tabulate(c(-1, 0, 3, NA)); tabulate(factor(c("a", "b", "a"))); tabulate(integer(0)); tabulate(c(2, 2), 1)
#---
tabulate("a")
#---
m <- matrix(1:6, 2, dimnames = list(c("a", "b"), c("x", "y", "z"))); rowSums(m); colSums(m); rowMeans(m); colMeans(m); rowSums(matrix(c(1, NA, 3, 4), 2), na.rm = TRUE); colMeans(matrix(c(1, NA, 3, 4), 2), na.rm = TRUE); colSums(matrix(c(TRUE, FALSE, TRUE, TRUE), 2)); rowSums(matrix(1:4, 2))
#---
rowSums(1:3)
#---
colSums(matrix(c("a", "b"), 1))
#---
m <- matrix(1:6, 2, dimnames = list(c("a", "b"), c("x", "y", "z"))); t(m); t(1:3); t(c(a = 1, b = 2)); dim(t(matrix(1:6, 2))); t(t(m)); t(matrix(1:4, 2, dimnames = list(NULL, c("p", "q"))))
#---
t("a")
#---
t(array(1:8, c(2, 2, 2)))
#---
diag(3); diag(2); diag(c(1, 2, 3)); diag(matrix(1:9, 3)); diag(matrix(1:6, 2)); diag(5, 2); diag(1, 2, 3); diag(0)
#---
diag("a")
#---
outer(1:3, 1:3); outer(1:2, 1:3, "+"); outer(c(a = 1, b = 2), c(x = 1, y = 2)); outer(1:2, 1:2, function(x, y) x * 10 + y); outer(1:2, numeric(0)); dim(outer(matrix(1:4, 2), 1:3))
#---
mean(c(1, 2, 3, 10)); mean(1:10); mean(c(1, NA, 3)); mean(c(1, NA, 3), na.rm = TRUE); mean(c(TRUE, FALSE)); mean(numeric(0)); mean("a"); mean(c(0.1, 0.2, 0.3)); mean(c(1e308, 1e308)); mean(c(Inf, 1)); mean(NULL)
#---
mean.money <- function(x, ...) "dispatched"; mean(structure(1, class = "money"))

#---
stopifnot(TRUE, 1 == 1)
#---
stopifnot(1 == 2)
#---
stopifnot(c(TRUE, FALSE))
#---
stopifnot("must be positive" = -1 > 0)
#---
f <- function(x) { stopifnot(is.numeric(x), x > 0); x }; f(2); f(-1)
#---
f <- function(x) stopifnot(is.character(x)); f(1)
#---
stopifnot(NA)
#---
stopifnot(1)
#---
rank(c(10, 30, 20, 20)); rank(c(10, 30, 20, 20), ties.method = "min"); rank(c(10, 30, 20, 20), ties.method = "max"); rank(c(10, 30, 20, 20), ties.method = "first"); rank(c("b", "a", "c")); rank(c(a = 3, b = 1)); rank(numeric(0)); rank(c(2, NA, 1))