T <- function(e) tryCatch(e, error = function(e) c("ERR", conditionMessage(e), deparse(conditionCall(e))[1]))
Q <- function(e) { set.seed(1); r <- T(e); list(r, .Random.seed[2]) }
print(Q(sample(10))); print(Q(sample(10, 3))); print(Q(sample(5, 10, TRUE))); print(Q(sample(5, 10))); print(Q(sample(c("a", "b", "c")))); print(Q(sample(c("a", "b", "c"), 2)))
#---
T <- function(e) tryCatch(e, error = function(e) c("ERR", conditionMessage(e), deparse(conditionCall(e))[1]))
Q <- function(e) { set.seed(1); r <- T(e); list(r, .Random.seed[2]) }
print(Q(sample(c(a = 1, b = 2, c = 3)))); print(Q(sample(list(1, 2, 3)))); print(Q(sample(factor(c("u", "v", "w"))))); print(Q(sample(TRUE))); print(Q(sample(c(TRUE, FALSE)))); print(Q(sample(1))); print(Q(sample(1, 1)))
#---
T <- function(e) tryCatch(e, error = function(e) c("ERR", conditionMessage(e), deparse(conditionCall(e))[1]))
Q <- function(e) { set.seed(1); r <- T(e); list(r, .Random.seed[2]) }
print(Q(sample(0))); print(Q(sample(0.5))); print(Q(sample(2.5))); print(Q(sample(-1))); print(Q(sample(NA))); print(Q(sample(NULL))); print(Q(sample(integer(0)))); print(Q(sample(character(0)))); print(Q(sample(Inf)))
#---
T <- function(e) tryCatch(e, error = function(e) c("ERR", conditionMessage(e), deparse(conditionCall(e))[1]))
Q <- function(e) { set.seed(1); r <- T(e); list(r, .Random.seed[2]) }
print(Q(sample(10, -1))); print(Q(sample(10, NA))); print(Q(sample(10, 2.9))); print(Q(sample(10, "2"))); print(Q(sample(10, c(1, 2)))); print(Q(sample(10, 0))); print(Q(sample(10, replace = NA))); print(Q(sample(10, replace = "a"))); print(Q(sample(10, 3, replace = 1)))
#---
T <- function(e) tryCatch(e, error = function(e) c("ERR", conditionMessage(e), deparse(conditionCall(e))[1]))
Q <- function(e) { set.seed(1); r <- T(e); list(r, .Random.seed[2]) }
print(Q(sample(1e10, 2))); print(Q(sample(2^40, 3))); print(Q(sample(2^53, 2))); print(Q(sample(2^60, 2))); print(Q(sample(4.5e15, 2))); print(Q(sample(5e15, 2)))
#---
T <- function(e) tryCatch(e, error = function(e) c("ERR", conditionMessage(e), deparse(conditionCall(e))[1]))
Q <- function(e) { set.seed(1); r <- T(e); list(r, .Random.seed[2]) }
print(Q(sample.int(10))); print(Q(sample.int(10, 3))); print(Q(sample.int(0))); print(Q(sample.int(0, 0))); print(Q(sample.int(-1))); print(Q(sample.int(10, 20))); print(Q(sample.int(10, 20, TRUE))); print(Q(sample.int(2.5)))
#---
T <- function(e) tryCatch(e, error = function(e) c("ERR", conditionMessage(e), deparse(conditionCall(e))[1]))
Q <- function(e) { set.seed(1); r <- T(e); list(r, .Random.seed[2]) }
print(Q(sample.int(10, useHash = TRUE))); print(Q(sample.int(10, 3, useHash = TRUE))); print(Q(sample.int(10, 3, useHash = FALSE))); print(Q(sample.int(20000000, 3))); print(Q(sample.int(20000000, 3, useHash = FALSE))); print(Q(sample.int(10, 6, useHash = TRUE))); print(Q(sample.int(10, 5, useHash = TRUE)))
#---
T <- function(e) tryCatch(e, error = function(e) c("ERR", conditionMessage(e), deparse(conditionCall(e))[1]))
Q <- function(e) { set.seed(1); r <- T(e); list(r, .Random.seed[2]) }
print(Q(sample.int(10, 3, TRUE, useHash = TRUE))); print(Q(sample.int(10, 3, prob = 1:10, useHash = TRUE))); print(Q(sample.int(c(1, 2)))); print(Q(sample.int(10, c(1, 2)))); print(Q(sample.int(10, 3, prob = 1:9)))
#---
T <- function(e) tryCatch(e, error = function(e) c("ERR", conditionMessage(e), deparse(conditionCall(e))[1]))
Q <- function(e) { set.seed(1); r <- T(e); list(r, .Random.seed[2]) }
print(Q(sample(5, 3, prob = c(.1, .2, .3, .2, .2)))); print(Q(sample(5, 3, TRUE, prob = c(.1, .2, .3, .2, .2)))); print(Q(sample(5, 3, prob = 1:5))); print(Q(sample(5, 3, prob = 1:4))); print(Q(sample(5, 3, prob = c(NA, 1, 1, 1, 1)))); print(Q(sample(5, 3, prob = c(-1, 1, 1, 1, 1))))
#---
T <- function(e) tryCatch(e, error = function(e) c("ERR", conditionMessage(e), deparse(conditionCall(e))[1]))
Q <- function(e) { set.seed(1); r <- T(e); list(r, .Random.seed[2]) }
print(Q(sample(5, 3, prob = c(0, 0, 0, 0, 0)))); print(Q(sample(5, 5, prob = c(1, 0, 0, 0, 0)))); print(Q(sample(5, 2, prob = c(1, 0, 0, 0, 0)))); print(Q(sample(5, 3, TRUE, prob = c(0, 0, 0, 0, 1)))); print(Q(sample.int(3, 3, prob = c(1, 1e-300, 1))))
#---
T <- function(e) tryCatch(e, error = function(e) c("ERR", conditionMessage(e), deparse(conditionCall(e))[1]))
Q <- function(e) { set.seed(1); r <- T(e); list(r, .Random.seed[2]) }
print(Q(sample.int(100, 5, prob = 1:100))); print(Q(sample.int(100, 5, TRUE, prob = 1:100))); print(Q(sample.int(300, 5, TRUE, prob = 1:300))); print(Q(sample.int(300, 5, prob = 1:300)))
#---
T <- function(e) tryCatch(e, error = function(e) c("ERR", conditionMessage(e), deparse(conditionCall(e))[1]))
Q <- function(e) { set.seed(1); r <- T(e); list(r, .Random.seed[2]) }
print(Q(sample.int(250, 3, TRUE, prob = rep(1, 250)))); print(Q(sample.int(199, 3, TRUE, prob = rep(1, 199)))); print(Q(sample.int(200, 3, TRUE, prob = rep(1, 200)))); print(Q(sample.int(5, 3, prob = c(.1, .2, .3, .2, .2) * 1e-20)))
#---
T <- function(e) tryCatch(e, error = function(e) c("ERR", conditionMessage(e), deparse(conditionCall(e))[1]))
Q <- function(e) { set.seed(1); r <- T(e); list(r, .Random.seed[2]) }
print(Q(sample.int(5, 3, prob = c(Inf, 1, 1, 1, 1)))); print(Q(suppressWarnings(sample.int(5, 3, prob = c(1, 1, 1, 1, "a"))))); print(Q(sample.int(5, 3, prob = list(1, 1, 1, 1, 1)))); print(Q(sample.int(5, 3, prob = numeric(0)))); print(Q(sample.int(5, 3, prob = NULL))); print(Q(sample.int(5, prob = matrix(1, 5, 1))))
#---
T <- function(e) tryCatch(e, error = function(e) c("ERR", conditionMessage(e), deparse(conditionCall(e))[1]))
Q <- function(e) { set.seed(1); r <- T(e); list(r, .Random.seed[2]) }
print(Q(sample.int(2^31 - 1, 2))); print(Q(sample.int(2^31, 2))); print(Q(sample.int(2^31 + 1, 2))); print(Q(sample.int(1e15, 3))); print(Q(sample.int(2^50, 4))); print(Q(sample.int(2^50 + 1, 4)))
#---
T <- function(e) tryCatch(e, error = function(e) c("ERR", conditionMessage(e), deparse(conditionCall(e))[1]))
Q <- function(e) { set.seed(1); r <- T(e); list(r, .Random.seed[2]) }
print(Q(sample.int(2^52, 4))); print(Q(sample.int(2^52 - 1, 4))); print(Q(sample.int(4e15, 2, TRUE))); print(Q(sample.int(2^50, 4, TRUE))); print(Q(sample.int(2^40, 1)))
#---
T <- function(e) tryCatch(e, error = function(e) c("ERR", conditionMessage(e), deparse(conditionCall(e))[1]))
Q <- function(e) { set.seed(1); r <- T(e); list(r, .Random.seed[2]) }
print(Q(sample.int(1, 1))); print(Q(sample.int(1, 2, TRUE))); print(Q(sample.int(1, 2))); print(Q(sample.int(2, 2))); print(Q(sample.int(3, 3))); print(Q(sample.int(17, 17))); print(Q(sample.int(65, 40)))
#---
set.seed(42); x <- sample(1000, 50); print(x); print(.Random.seed[1:4]); y <- sample(1000, 50, TRUE); print(y); print(sum(as.numeric(.Random.seed) %% 1000003))
#---
set.seed(7); print(sample(letters, 8)); print(sample(100, 10, prob = (1:100)^2)); print(sample(3, 20, TRUE, prob = c(.5, .3, .2))); print(sum(as.numeric(.Random.seed) %% 1000003))
#---
set.seed(3); p <- c(.2, .2, .2, .2, .2); print(sample(5, 5, prob = p)); print(sample(5, 10, TRUE, prob = p)); p2 <- c(.1, .1, .4, .4); print(sample(4, 4, prob = p2)); print(sample(4, 12, TRUE, prob = p2)); print(sum(as.numeric(.Random.seed) %% 1000003))
#---
set.seed(11); print(sample(600, 8, TRUE, prob = rep(c(1, 2, 3), 200))); print(sample(260, 6, TRUE, prob = 1:260)); print(sum(as.numeric(.Random.seed) %% 1000003))
#---
set.seed(1); print(sample(2e7, 4)); print(sample(3e7, 3)); print(.Random.seed[2])
#---
set.seed(1); print(sample.int(2^33, 3)); print(sample(2^45, 2)); print(.Random.seed[2])
#---
set.seed(1); a <- jitter(1:5); print(signif(a, 12)); set.seed(2); print(signif(jitter(c(1, 2, 3.5), amount = 0), 12)); set.seed(3); print(signif(jitter(c(5, 5, 5)), 12)); print(.Random.seed[2])
#---
set.seed(1); print(signif(jitter(c(0, 0, 0)), 12)); print(signif(jitter(1:3, factor = 3), 12)); print(signif(jitter(1:3, amount = 0.5), 12)); print(signif(jitter(10), 12)); print(jitter(numeric(0))); print(tryCatch(jitter("a"), error = conditionMessage))
#---
set.seed(1); x <- c(a = 1.5, b = 2.5, c = 10); print(signif(jitter(x), 12)); print(signif(jitter(c(1, NA, 3)), 12)); print(.Random.seed[2])
#---
set.seed(10); x <- rep(1:3, each = 4); print(signif(jitter(x), 12)); print(signif(jitter(x, 2), 12)); print(sum(as.numeric(.Random.seed) %% 1000003))
