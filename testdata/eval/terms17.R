`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ a * b + I(x^2) + offset(w) - 1)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ a + b %in% a + (c + d)^2)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(~ a + b)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ a:b + a, keep.order = TRUE)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ a:b + a)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ a / b)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ a * b * c - a:b:c)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ . , data = data.frame(y = 1, u = 2, v = 3))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
terms(y ~ ., data = data.frame(y = 1, u = 2, v = 3))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ a + b - a)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ (a + b + c)^2 - a:b)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ 0 + a)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ a + 0)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ a - 1 + 1)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ b %in% a)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ a:b)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ 1)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ 0)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(~ 1)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ a + b:a)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl((y ~ (a + b) / c))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ a / (b + c))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(log(y) ~ sqrt(x) + a)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ y + x)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ s(x) + s(z) + a, specials = "s")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
attr(terms(y ~ s(x) + s(z) + a, specials = "s"), "specials")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
attr(terms(y ~ a + b, specials = c("a", "zz")), "specials")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
names(attributes(terms(y ~ a + b, specials = "b")))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ a * b - a)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ a + b + a)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ x^2)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ .^2, data = data.frame(y = 1, u = 2, v = 3, w = 4))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ . - u, data = data.frame(y = 1, u = 2, v = 3, w = 4))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ a + offset(o) + b)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ a + (b + c) %in% a)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ a * b, keep.order = TRUE)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ a:b:c + a + b:c, keep.order = TRUE)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ -1 + a)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ a + I(b + c))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ a + (1))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ b + a:b)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(cbind(y1, y2) ~ a)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ a + b^2 + a:b:b)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ (a + b + c)^3 - a:b:c)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y + z ~ a)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ a | b)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ a + log(b, base = 2))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ -a)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ a + (b - b))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ a - a)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(~ a - 1)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ a:b - 1)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ a + b - 1)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ a + a:b:c + b)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ a + b + c + a:b:c + b:c, keep.order = TRUE)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ (a + b)^2 + a %in% b)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ a + I(a^2) + poly(b, 2))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ a + offset(log(z)) + offset(w))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
tl(y ~ a * offset(w))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tl <- function(f, ...) { t <- terms(f, ...); list(labels = attr(t, "term.labels"), order = attr(t, "order"), intercept = attr(t, "intercept"), response = attr(t, "response"), offset = attr(t, "offset"), factors = attr(t, "factors"), vars = attr(t, "variables")) }
attr(terms(y ~ 1 + 0), "intercept")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
attr(terms(y ~ 0 + 1), "intercept")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
attr(terms(y ~ a - 0), "intercept")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
attr(terms(y ~ a + (0)), "intercept")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
attr(terms(y ~ a - (b + 1)), "intercept")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
attr(terms(y ~ a + (b - 1)), "intercept")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
terms(y ~ a, specials = "s")$foo
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
is.empty.model(y ~ 0)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
is.empty.model(y ~ a)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
is.empty.model(y ~ 1)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
is.empty.model(y ~ a - 1)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
is.empty.model(y ~ 0 + offset(o))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
terms(y ~ a * b + I(x^2) + offset(w) - 1)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
terms(y ~ a + b)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
terms(~ a)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
terms(y ~ 1)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
print(terms(y ~ a))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
class(terms(y ~ a))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
names(attributes(terms(y ~ a * b + offset(w))))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
names(attributes(terms(y ~ a, specials = "s")))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
names(attributes(terms(~ 1)))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
attributes(terms(y ~ a))$.Environment
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
identical(attr(terms(y ~ a), ".Environment"), globalenv())
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
t <- terms(y ~ a); identical(terms(t), t)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
terms(y ~ .^2, data = data.frame(y = 1, u = 2, v = 3))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
terms(y ~ u + ., data = data.frame(y = 1, u = 2, v = 3, w = 4))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
terms(y ~ . - u, data = data.frame(y = 1, u = 2, v = 3, w = 4))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
terms(y ~ ., data = data.frame(y = 1, u = 2))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
terms(log(y) ~ ., data = data.frame(y = 1, u = 2, v = 3))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
terms(~ ., data = data.frame(u = 2, v = 3))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
terms(y ~ ., data = list(y = 1, a = 2, b = 3))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tryCatch(terms(y ~ .), error = function(e) conditionMessage(e))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
attr(terms(y ~ ., allowDotAsName = TRUE), "term.labels")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tryCatch(terms(y ~ a^b), error = function(e) conditionMessage(e))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tryCatch(terms(y ~ a^0), error = function(e) conditionMessage(e))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tryCatch(terms(y ~ a^-1), error = function(e) conditionMessage(e))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tryCatch(terms(y ~ a^1.5), error = function(e) conditionMessage(e))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tryCatch(terms(y ~ 2), error = function(e) conditionMessage(e))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tryCatch(terms(y ~ a + 2), error = function(e) conditionMessage(e))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tryCatch(terms(1), error = function(e) conditionMessage(e))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tryCatch(terms("y ~ x"), error = function(e) conditionMessage(e))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tryCatch(terms(y ~ a^b), error = function(e) deparse(conditionCall(e)))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
tryCatch(terms(y ~ a + 2), error = function(e) deparse(conditionCall(e)))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
terms(y ~ a^b)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
attr(terms(y ~ a ~ b), "term.labels")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
terms(y ~ a * b + c, simplify = TRUE)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
terms(y ~ a - 1, simplify = TRUE)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
terms(y ~ a:b:c + a - 1 + offset(o), simplify = TRUE)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
terms(y ~ 0 + a, simplify = TRUE)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
terms(y ~ 1, simplify = TRUE)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
terms(y ~ a + offset(o) + b, simplify = TRUE)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
formula(terms(y ~ a * b + c, simplify = TRUE))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
attr(terms(y ~ a, abb = TRUE), "term.labels")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
nchar(paste(deparse(terms(y ~ a + b)), collapse = ""))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
length(attr(terms(y ~ (a + b + c + d)^4), "term.labels"))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
attr(terms(y ~ (a + b + c + d)^4), "order")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
dim(attr(terms(y ~ (a + b + c)^3), "factors"))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
dimnames(attr(terms(y ~ a * b), "factors"))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
typeof(attr(terms(y ~ a * b), "factors"))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
typeof(attr(terms(y ~ a * b), "order"))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
typeof(attr(terms(y ~ a * b), "intercept"))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
typeof(attr(terms(y ~ a + offset(o)), "offset"))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
attr(terms(y ~ a * b), "variables")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
class(attr(terms(y ~ a * b), "variables"))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
as.list(attr(terms(y ~ a * b), "variables"))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
length(attr(terms(y ~ a * b), "variables"))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
labels(terms(y ~ a * b))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
all.vars(terms(y ~ a * b))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
all.vars(attr(terms(y ~ a * log(b)), "variables"))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
update(terms(y ~ a), . ~ . + b)
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
class(update(terms(y ~ a), . ~ . + b))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
formula(terms(y ~ a + b))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
as.formula(terms(y ~ a + b))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
terms(y ~ I(a + b) + I(a - b))
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
attr(terms(y ~ I(a + b) + I(a - b)), "term.labels")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
attr(terms(y ~ `a b` + c), "term.labels")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
attr(terms(y ~ `a b` * c), "factors")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
attr(terms(`y z` ~ a), "factors")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
attr(terms(y ~ x1 + x2 + x1:x2 + I(x1^2)), "factors")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
attr(terms(y ~ a * b * c * d), "term.labels")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
attr(terms(y ~ (a + b) * (c + d)), "term.labels")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
attr(terms(y ~ (a + b) %in% (c + d)), "term.labels")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
attr(terms(y ~ a/b/c), "term.labels")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
attr(terms(y ~ a/b/c), "factors")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
attr(terms(y ~ (a + b + c)^2 - a - b), "factors")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
attr(terms(y ~ a + b + a:b - a:b), "term.labels")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
attr(terms(y ~ a - b), "term.labels")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
attr(terms(y ~ a:b - a), "term.labels")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
attr(terms(y ~ - 1 - a + b), "term.labels")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
attr(terms(y ~ +a), "term.labels")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
attr(terms(y ~ +a + +b), "term.labels")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
attr(terms(y ~ (a)), "term.labels")
#---
`~` <- function(...) { f <- eval(sys.call(), baseenv()); environment(f) <- globalenv(); f }
attr(terms(y ~ ((a + b))), "term.labels")
