S <- function() c(.Random.seed[1:4], as.integer(sum(as.numeric(.Random.seed) %% 1000003)))
set.seed(1); print(.Random.seed[1:5]); print(length(.Random.seed)); print(typeof(.Random.seed)); print(runif(3)); print(S())
#---
S <- function() c(.Random.seed[1:4], as.integer(sum(as.numeric(.Random.seed) %% 1000003)))
for (s in c(42, 0, -1, 2147483647, -2147483647)) { set.seed(s); print(S()); print(runif(2)) }
#---
S <- function() c(.Random.seed[1:4], as.integer(sum(as.numeric(.Random.seed) %% 1000003)))
set.seed(1.9); a <- runif(1); set.seed(1); b <- runif(1); print(a == b); set.seed(-1.9); print(S()); set.seed(TRUE); print(S()); set.seed("7"); print(S()); set.seed(c(3, 9)); print(S())
#---
print(tryCatch(set.seed(NA), error = conditionMessage)); print(tryCatch(set.seed(NaN), error = conditionMessage)); print(tryCatch(set.seed(numeric(0)), error = conditionMessage))
#---
print(tryCatch(set.seed(2^31), warning = conditionMessage)); print(tryCatch(suppressWarnings(set.seed(Inf)), error = conditionMessage))
#---
print(tryCatch(set.seed("a"), warning = conditionMessage)); print(tryCatch(suppressWarnings(set.seed("a")), error = conditionMessage))
#---
print(tryCatch(set.seed(), error = conditionMessage)); print(is.null(set.seed(3))); print(withVisible(set.seed(3))$visible)
#---
set.seed(NULL); print(length(.Random.seed)); print(.Random.seed[1:2])
#---
set.seed(1); u <- runif(1); print(.Random.seed[2]); invisible(runif(623)); print(.Random.seed[2]); invisible(runif(1)); print(.Random.seed[2]); print(runif(1))
#---
if (exists(".Random.seed", envir = globalenv())) rm(list = ".Random.seed", envir = globalenv()); print(runif(0)); print(exists(".Random.seed")); print(length(runif(3))); print(exists(".Random.seed"))
#---
set.seed(1); print(runif(5, 10, 20)); print(runif(3, c(0, 10), 11)); print(runif(2, min = 2, max = 3)); print(runif(2, ma = 3))
#---
set.seed(1); print(tryCatch(runif(-1), error = conditionMessage)); print(tryCatch(runif(NA), error = conditionMessage)); print(tryCatch(runif(NULL), error = conditionMessage)); print(tryCatch(runif(list(1)), error = conditionMessage)); print(.Random.seed[2])
#---
set.seed(1); print(runif("3")); print(runif(c(5, 6))); print(runif(2.7)); print(runif(TRUE)); print(length(runif(numeric(0)))); print(length(runif(0.5)))
#---
set.seed(1); print(tryCatch(runif("a"), warning = conditionMessage)); print(tryCatch(suppressWarnings(runif("a")), error = conditionMessage)); print(tryCatch(runif(), error = conditionMessage))
#---
set.seed(1); print(suppressWarnings(runif(2, 2, 1))); print(tryCatch(runif(2, 2, 1), warning = conditionMessage)); print(suppressWarnings(runif(2, NA))); print(suppressWarnings(runif(1, -Inf, 1))); print(.Random.seed[2])
#---
set.seed(1); print(runif(1, 1, 1)); print(runif(3, 1, c(1, 5, 1))); print(.Random.seed[2]); print(suppressWarnings(runif(2, numeric(0)))); print(tryCatch(runif(2, "a"), error = conditionMessage)); print(tryCatch(runif(2, min = NULL), error = conditionMessage))
#---
set.seed(1); print(suppressWarnings(runif(2, 1:2))); print(.Random.seed[2])
#---
set.seed(1); x <- rnorm(5); print(signif(x, 12)); print(.Random.seed[2])
#---
set.seed(1); print(signif(rnorm(5, 10, 2), 12)); print(signif(rnorm(3, c(0, 100)), 12)); print(.Random.seed[2])
#---
set.seed(1); print(rnorm(3, 0, 0)); print(.Random.seed[2]); print(rnorm(1, Inf)); print(rnorm(2, -Inf, 0)); print(.Random.seed[2])
#---
set.seed(1); print(suppressWarnings(rnorm(2, 0, -1))); print(suppressWarnings(rnorm(1, 0, Inf))); print(suppressWarnings(rnorm(1, NA))); print(suppressWarnings(rnorm(1, 0, NaN))); print(.Random.seed[2]); print(tryCatch(rnorm(2, 0, -1), warning = conditionMessage))
#---
set.seed(1); print(signif(rnorm(3, 1, c(0, 1)), 12)); print(.Random.seed[2]); print(signif(rnorm(1, 0, 1e-320), 12)); print(.Random.seed[2])
#---
set.seed(1); print(tryCatch(rnorm(-1), error = conditionMessage)); print(tryCatch(rnorm(1, "a"), error = conditionMessage)); print(length(rnorm(0))); print(.Random.seed[2])
#---
set.seed(1); print(signif(rnorm(2, mean = 1:2, sd = 0.5), 12)); print(.Random.seed[2]); print(signif(rnorm(1000)[c(1, 500, 1000)], 12)); print(.Random.seed[2])
#---
set.seed(2024); x <- runif(1e5); print(c(min(x) > 0, max(x) < 1)); print(signif(mean(x), 12)); print(signif(sum(x), 12)); print(S <- .Random.seed[2])
#---
set.seed(2024); x <- rnorm(1e5); print(signif(mean(x), 10)); print(signif(sd(x), 10)); print(.Random.seed[2])
#---
set.seed(5); print(.Random.seed[1]); rm(.Random.seed); r <- runif(1); print(length(.Random.seed)); print(.Random.seed[1:2])
#---
set.seed(1); f <- function() { .Random.seed <- 1:3; runif(1) }; print(f()); print(.Random.seed[2])
#---
set.seed(1); a <- .Random.seed; u1 <- runif(3); .Random.seed <- a; u2 <- runif(3); print(identical(u1, u2)); print(u2)
#---
set.seed(1); s <- .Random.seed; s[2] <- 5L; .Random.seed <- s; print(runif(2)); print(.Random.seed[2]); s[2] <- 623L; .Random.seed <- s; print(runif(2)); s[2] <- 1L; .Random.seed <- s; print(runif(2))
#---
set.seed(1); s <- .Random.seed; s[2] <- 0L; .Random.seed <- s; print(runif(2)); s[2] <- 700L; .Random.seed <- s; print(runif(2)); s[2] <- 625L; .Random.seed <- s; print(runif(2)); print(.Random.seed[2:3])
#---
set.seed(1); s <- .Random.seed; s[2] <- 625L; .Random.seed <- s; print(runif(4)); print(.Random.seed[2:4])
#---
set.seed(1); s <- .Random.seed; s[3:5] <- c(NA, 0L, -1L); .Random.seed <- s; print(runif(3)); print(.Random.seed[3:5])
#---
set.seed(1); .Random.seed <- c(.Random.seed, 99L, 98L); print(runif(1)); print(length(.Random.seed))
#---
set.seed(1); .Random.seed <- .Random.seed[1:625]; print(tryCatch(runif(1), error = conditionMessage)); print(length(.Random.seed))
#---
set.seed(1); s <- .Random.seed; attr(s, "foo") <- "bar"; .Random.seed <- s; invisible(runif(1)); print(attributes(.Random.seed)); print(.Random.seed[2])
#---
set.seed(1); .Random.seed <- as.numeric(.Random.seed); print(tryCatch(runif(1), warning = conditionMessage)); print(typeof(.Random.seed)); print(length(.Random.seed))
#---
set.seed(1); .Random.seed <- NULL; print(tryCatch(runif(1), warning = conditionMessage)); print(length(.Random.seed))
#---
set.seed(1); .Random.seed <- "a"; print(tryCatch(runif(1), warning = conditionMessage))
#---
set.seed(1); s <- .Random.seed; s[1] <- NA; .Random.seed <- s; print(tryCatch(runif(1), warning = conditionMessage)); print(.Random.seed[1]); s[1] <- -3L; .Random.seed <- s; print(tryCatch(runif(1), warning = conditionMessage)); s[1] <- 11001L; .Random.seed <- s; print(tryCatch(runif(1), warning = conditionMessage))
#---
set.seed(1); s <- .Random.seed; s[1] <- 10407L; .Random.seed <- s; print(tryCatch(runif(1), warning = conditionMessage)); s[1] <- 10603L; .Random.seed <- s; print(tryCatch(runif(1), warning = conditionMessage)); s[1] <- 10403L; .Random.seed <- s; print(runif(1))
#---
set.seed(1); .Random.seed <- integer(0); print(tryCatch(runif(1), warning = conditionMessage)); print(length(.Random.seed))
#---
set.seed(1); s <- .Random.seed; s[1] <- 10401L; .Random.seed <- s; print(tryCatch(runif(1), error = function(e) "unsupported-or-error")); s[1] <- 10403L
#---
set.seed(1, kind = "Mersenne-Twister"); print(.Random.seed[1:4]); set.seed(1, kind = "default"); print(.Random.seed[1:4]); set.seed(1, kind = "Mers"); print(.Random.seed[1])
#---
set.seed(1, normal.kind = "Inversion", sample.kind = "Rejection"); print(.Random.seed[1:3]); print(tryCatch(set.seed(1, kind = "foo"), error = conditionMessage)); print(tryCatch(set.seed(1, normal.kind = "foo"), error = conditionMessage))
#---
print(RNGkind()); set.seed(1); a <- .Random.seed; print(RNGkind()); print(identical(a, .Random.seed)); print(withVisible(RNGkind())$visible)
#---
set.seed(1); print(.Random.seed[3:4]); r <- RNGkind("Mersenne-Twister"); print(.Random.seed[1:4]); print(r); print(withVisible(RNGkind("Mersenne-Twister"))$visible)
#---
set.seed(7); r <- RNGkind("Mersenne-Twister"); print(.Random.seed[1:5])
#---
set.seed(1); r <- RNGkind(normal.kind = "Inversion"); print(.Random.seed[1:4]); r <- RNGkind(sample.kind = "Rejection"); print(.Random.seed[1:4]); r <- RNGkind(normal.kind = "default"); print(.Random.seed[1:3])
#---
print(tryCatch(RNGkind("foo"), error = conditionMessage)); print(tryCatch(RNGkind(1), error = conditionMessage)); print(tryCatch(RNGkind(NA), error = conditionMessage)); print(tryCatch(RNGkind(c("a", "b")), error = conditionMessage)); print(tryCatch(RNGkind("M"), error = conditionMessage)); print(tryCatch(RNGkind("Knuth"), error = conditionMessage))
#---
print(tryCatch(RNGkind(normal.kind = "foo"), error = conditionMessage)); print(tryCatch(RNGkind(sample.kind = "foo"), error = conditionMessage)); print(tryCatch(RNGkind(normal.kind = 1), error = conditionMessage)); print(tryCatch(RNGkind(sample.kind = c("Rejection", "Rounding")), error = conditionMessage)); print(tryCatch(RNGkind(sample.kind = "R"), error = conditionMessage))
#---
set.seed(1); r <- RNGkind("Mers"); print(r); r <- RNGkind(normal.kind = "Inv"); r <- RNGkind(sample.kind = "Rej"); print(RNGkind())
#---
set.seed(1); print(tryCatch(RNGkind(sample.kind = "Rounding"), warning = conditionMessage)); print(tryCatch(RNGkind(sample.kind = "Rounding"), warning = function(w) deparse(conditionCall(w)))); r <- suppressWarnings(RNGkind(sample.kind = "Rounding")); print(RNGkind()); print(.Random.seed[1]); r <- RNGkind(sample.kind = "Rejection"); print(.Random.seed[1])
#---
print(tryCatch(set.seed(1, sample.kind = "Rounding"), warning = function(w) c(conditionMessage(w), deparse(conditionCall(w)))))
#---
print(tryCatch(RNGversion("3.5.0"), warning = function(w) c(conditionMessage(w), deparse(conditionCall(w))))); r <- suppressWarnings(RNGversion("3.5.0")); print(RNGkind()); r <- RNGversion("3.6.0"); print(RNGkind()); print(.Random.seed[1])
#---
r <- suppressWarnings(RNGversion(3.5)); print(RNGkind()); r <- RNGversion("4.0.0"); print(RNGkind()); r <- suppressWarnings(RNGversion("1.7.0")); print(RNGkind()); r <- RNGversion("3.6"); print(RNGkind()); print(withVisible(RNGversion("3.6.0"))$visible)
#---
print(tryCatch(suppressWarnings(RNGversion("foo")), error = conditionMessage)); print(tryCatch(RNGversion("foo"), warning = conditionMessage)); print(tryCatch(RNGversion("3"), error = conditionMessage))
#---
r <- suppressWarnings(RNGversion("3.5.0")); set.seed(1); print(sample(10)); print(sample.int(100, 5)); print(sample(5, 10, TRUE)); r <- RNGversion("4.0.0"); set.seed(1); print(sample(10))
#---
r <- suppressWarnings(RNGkind(sample.kind = "Rounding")); set.seed(1); print(sample(10, 3)); print(sample(20, 4, TRUE)); print(sample(5, 3, prob = c(.1, .2, .3, .2, .2))); print(runif(1)); print(.Random.seed[2]); r <- RNGkind("Mersenne-Twister", "Inversion", "Rejection")
#---
set.seed(1); r <- RNGkind(normal.kind = "Box-Muller"); print(RNGkind()); print(runif(1)); print(.Random.seed[1]); r <- RNGkind("Mersenne-Twister", "Inversion", "Rejection")
