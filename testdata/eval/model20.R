contr.treatment(3)
#---
contr.treatment(c("a","b","c"), base = 2)
#---
contr.treatment(3, contrasts = FALSE)
#---
contr.treatment(1)
#---
contr.treatment(0)
#---
contr.treatment(-1)
#---
contr.treatment("a")
#---
contr.treatment(3, base = 4)
#---
contr.treatment(2.5)
#---
contr.sum(3)
#---
contr.sum(c("a","b","c"))
#---
contr.sum(3, contrasts = FALSE)
#---
contr.sum(1)
#---
contr.sum(0)
#---
contr.helmert(4)
#---
contr.helmert(c("a","b","c"))
#---
contr.helmert(3, contrasts = FALSE)
#---
contr.helmert(1)
#---
contr.SAS(3)
#---
contr.SAS(c("a","b","c"))
#---
round(contr.poly(3), 10) + 0
#---
round(contr.poly(4), 10) + 0
#---
round(contr.poly(5), 10) + 0
#---
round(contr.poly(5, scores = c(1,2,4,8,16)), 10) + 0
#---
round(contr.poly(3, contrasts = FALSE), 10) + 0
#---
contr.poly(1)
#---
contr.poly(3, scores = 1:2)
#---
contr.poly(3, scores = c(1,1,2))
#---
max(abs(crossprod(contr.poly(6)) - diag(5))) < 1e-14
#---
options("contrasts")
#---
f <- factor(c("a","b","c","a","b","c")); contrasts(f)
#---
f <- factor(c("a","b","c","a","b","c")); contrasts(f, contrasts = FALSE)
#---
contrasts(c(TRUE, FALSE, TRUE))
#---
contrasts(1:3)
#---
f <- factor(c("a","b","c")); contrasts(f) <- "contr.sum"; contrasts(f)
#---
f <- factor(c("a","b","c")); contrasts(f) <- "contr.sum"; attributes(f)
#---
f <- factor(c("a","b","c")); contrasts(f, 1) <- contr.helmert(3)[, 1, drop = FALSE]; contrasts(f)
#---
f <- factor(c("a","b","c")); contrasts(f) <- contr.helmert; contrasts(f)
#---
f <- factor(c("a","b","c","d")); contrasts(f, 2) <- contr.helmert(4)[, 1, drop = FALSE]; round(contrasts(f), 10) + 0
#---
f <- factor(c("a","b","c")); contrasts(f) <- TRUE
#---
f <- factor("a"); contrasts(f) <- "contr.sum"
#---
o <- ordered(c("lo","mid","hi","lo"), levels = c("lo","mid","hi")); round(contrasts(o), 10) + 0
#---
d <- data.frame(y = c(1,3,2,5,4,6), x = c(1.5,2,3.5,4,5,7), f = factor(c("a","b","c","a","b","c")), g = factor(c("u","v","u","v","u","v")), l = c(TRUE,FALSE,TRUE,TRUE,FALSE,FALSE), s = c("p","q","p","q","p","p")); model.matrix(y ~ x, d)
#---
d <- data.frame(y = c(1,3,2,5,4,6), f = factor(c("a","b","c","a","b","c"))); model.matrix(y ~ f, d)
#---
d <- data.frame(f = factor(c("a","b","c","a","b","c")), g = factor(c("u","v","u","v","u","v"))); model.matrix(~ f + g, d)
#---
d <- data.frame(f = factor(c("a","b","c","a","b","c")), g = factor(c("u","v","u","v","u","v"))); model.matrix(~ f * g, d)
#---
d <- data.frame(f = factor(c("a","b","c","a","b","c")), g = factor(c("u","v","u","v","u","v"))); model.matrix(~ f:g, d)
#---
d <- data.frame(f = factor(c("a","b","c","a","b","c")), g = factor(c("u","v","u","v","u","v"))); model.matrix(~ 0 + f, d)
#---
d <- data.frame(f = factor(c("a","b","c","a","b","c")), g = factor(c("u","v","u","v","u","v"))); model.matrix(~ f - 1 + g, d)
#---
d <- data.frame(f = factor(c("a","b","c","a","b","c")), g = factor(c("u","v","u","v","u","v"))); model.matrix(~ 0 + f:g, d)
#---
d <- data.frame(x = c(1.5,2,3.5,4,5,7), l = c(TRUE,FALSE,TRUE,TRUE,FALSE,FALSE), s = c("p","q","p","q","p","p")); model.matrix(~ x + l + s, d)
#---
d <- data.frame(x = c(1.5,2,3.5,4,5,7), f = factor(c("a","b","c","a","b","c"))); model.matrix(~ x:f, d)
#---
d <- data.frame(x = c(1.5,2,3.5,4,5,7), f = factor(c("a","b","c","a","b","c"))); model.matrix(~ f/x, d)
#---
d <- data.frame(x = c(1.5,2,3.5,4,5,7), f = factor(c("a","b","c","a","b","c"))); model.matrix(~ x %in% f, d)
#---
d <- data.frame(f = factor(c("a","b","c","a","b","c"))); model.matrix(~ f, d, contrasts.arg = list(f = "contr.sum"))
#---
d <- data.frame(f = factor(c("a","b","c","a","b","c"))); model.matrix(~ f, d, contrasts.arg = list(f = contr.helmert))
#---
d <- data.frame(f = factor(c("a","b","c","a","b","c"))); model.matrix(~ f, d, contrasts.arg = list(f = "nonexistent"))
#---
d <- data.frame(f = factor(c("a","b","c","a","b","c"))); model.matrix(~ f, d, contrasts.arg = list(g = "contr.sum"))
#---
d <- data.frame(x = c(1.5,2,3.5,4,5,7)); model.matrix(~ log(x) + I(x^2), d)
#---
d <- data.frame(x = c(1.5,2,3.5,4,5,7), f = factor(c("a","b","c","a","b","c"))); attributes(model.matrix(~ f + x, d))
#---
d <- data.frame(x = c(1.5,2,3.5,4,5,7)); model.matrix(~ 1, d)
#---
d <- data.frame(x = c(1.5,2,3.5,4,5,7)); model.matrix(~ 0, d)
#---
d <- data.frame(y = c(1,3,2,5,4,6), f = factor(c("a","b","c","a","b","c"))); m <- model.matrix(y ~ ordered(f), d); m[] <- signif(m, 9); m[abs(m) < 1e-10] <- 0; m
#---
d <- data.frame(y = c(1,3,2,5,4,6), x = c(1.5,2,3.5,4,5,7), f = factor(c("a","b","c","a","b","c"))); model.matrix(terms(y ~ x + f), d)
#---
d <- data.frame(f = factor(c("a","b","c","a","b","c"))); model.matrix(~ f, d[0, , drop = FALSE])
#---
d <- data.frame(x = c(1.5,2,3.5)); model.matrix(~ nofactor, d)
#---
d <- data.frame(y = c(1,3,2), f = factor(c("a","b","c"))); b <- model.matrix(y ~ f, d); dim(b)
#---
d <- data.frame(f = factor(c("a","b","c","a","b","c"), levels = c("c","b","a"))); model.matrix(~ f, d)
#---
d <- data.frame(y = c(1,3,2,5,4,6), f = factor(c("a","b","c","a","b","c")), g = factor(c("u","v","u","v","u","v"))); X <- model.matrix(y ~ f + g, d); b <- qr.coef(qr(X), d$y); round(unname(b), 8) + 0
