median(c(3, 1, 2))
#---
median(1:4)
#---
median(c(1L, 3L))
#---
median(1:3)
#---
median(c(1, NA, 3))
#---
median(c(1, NA, 3), na.rm = TRUE)
#---
median(numeric(0))
#---
median(integer(0))
#---
median("a")
#---
median(factor(c("a", "b")))
#---
median(as.Date("2020-01-01") + c(0, 3))
#---
median(c(TRUE, FALSE, TRUE))
#---
median(c(a = 1, b = 5, c = 3))
#---
median(c(1, Inf, -Inf, 2))
#---
median(c(.Machine$double.xmax, .Machine$double.xmax))
#---
median(NULL)
#---
median(c(NaN, 1))
#---
median(c(5, 3, 8, 1, 9, 2, 7))
#---
median(c(5, 3, 8, 1, 9, 2, 7, 4))
#---
var(c(1, 2, 3, 4)) # ulp
#---
var(1) # ulp
#---
var(c(1, NA, 3)) # ulp
#---
var(c(1, NA, 3), na.rm = TRUE) # ulp
#---
var(1:10, (1:10)^2) # ulp
#---
var(cbind(a = 1:5, b = c(2, 4, 5, 4, 5))) # ulp
#---
var(data.frame(a = 1:5, b = c(2, 4, 5, 4, 5))) # ulp
#---
var(cbind(1:5, c(2, 4, 5, NA, 5)), use = "complete.obs") # ulp
#---
var(cbind(1:5, c(2, 4, 5, NA, 5)), use = "pairwise.complete.obs") # ulp
#---
var(cbind(1:5, c(2, 4, 5, NA, 5)), na.rm = TRUE) # ulp
#---
var(NULL) # ulp
#---
var(factor(1:3)) # ulp
#---
var(c(TRUE, FALSE, TRUE)) # ulp
#---
var(1:3, 1:4) # ulp
#---
var(numeric(0)) # ulp
#---
var(c(1e10 + 1, 1e10 + 2, 1e10 + 3)) # ulp
#---
var(matrix(1:6, 3), matrix(c(1, 3, 2, 5, 4, 6), 3)) # ulp
#---
var(c(1, 2, Inf)) # ulp
#---
sd(c(1, 2, 3, 4)) # ulp
#---
sd(1) # ulp
#---
sd(c(1, NA, 3), na.rm = TRUE) # ulp
#---
sd(1:10) # ulp
#---
cov(1:5, c(2, 4, 5, 4, 5)) # ulp
#---
cov(cbind(a = 1:5, b = c(2, 4, 5, 4, 5))) # ulp
#---
cov(1:5, c(2, 4, 5, 4, 5), method = "spearman") # ulp
#---
cov(1:5, c(2, 4, 5, 4, 5), method = "kendall") # ulp
#---
cor(1:5, c(2, 4, 5, 4, 5)) # ulp
#---
cor(1:5, c(2, 4, 5, 4, 5), method = "spearman") # ulp
#---
cor(1:5, c(2, 4, 5, 4, 5), method = "kendall") # ulp
#---
cor(cbind(a = 1:5, b = c(2, 4, 5, 4, 5), c = c(5, 3, 4, 1, 2))) # ulp
#---
cor(cbind(a = 1:5, b = c(2, 4, 5, 4, 5), c = c(5, 3, 4, 1, 2)), method = "spearman") # ulp
#---
cor(cbind(a = 1:5, b = c(2, 4, 5, 4, 5), c = c(5, 3, 4, 1, 2)), method = "kendall") # ulp
#---
cor(1:5, rep(1, 5)) # ulp
#---
cor(1:3, 1:4) # ulp
#---
cor(c(1, NA, 3, 4), c(2, 3, NA, 5), use = "complete.obs") # ulp
#---
cor(c(1, NA, 3, 4), c(2, 3, NA, 5)) # ulp
#---
cor(c(1, NA, 3, 4, 7), c(2, 3, NA, 5, 1), use = "pairwise.complete.obs") # ulp
#---
cor(c(1, NA, 3, 4), c(2, 3, NA, 5), use = "all.obs") # ulp
#---
cor(c(1, NA, 3, 4), c(2, 3, NA, 5), use = "na.or.complete") # ulp
#---
cor(c(1, 2, 3, 4), c(2, 3, 1, 5), use = "everything") # ulp
#---
cor(data.frame(a = 1:4, b = c(2, 3, 1, 5))) # ulp
#---
cor(1:4, c(2, 3, 1, 5), method = "bad") # ulp
#---
cor("a", 1) # ulp
#---
cor(1:5) # ulp
#---
cor(cbind(1:5, c(2, 4, 5, 4, 5)), 1:5) # ulp
#---
cor(c(1, 1, 2, 2, 3, 3), c(1, 2, 1, 2, 2, 3), method = "kendall") # ulp
#---
cor(c(1, 1, 2, 2, 3, 3), c(1, 2, 1, 2, 2, 3), method = "spearman") # ulp
#---
cor(c(1, 2, 3), c(1, 2, 3), use = "bad") # ulp
#---
cor(1, 1) # ulp
#---
cor(c(1e10 + 1, 1e10 + 2, 1e10 + 3), c(1, 2, 3.5)) # ulp
#---
quantile(1:10) # ulp
#---
quantile(c(1.5, 2, 10, 7), c(.1, .5, .9)) # ulp
#---
quantile(c(1, NA, 2))
#---
quantile(c(1, 2, NA), na.rm = TRUE)
#---
quantile(1:10, names = FALSE) # ulp
#---
quantile(1:10, c(1/3, 2/3, 0.025, 0.9999)) # ulp
#---
quantile(1:10, 1.5) # ulp
#---
quantile(1:10, NA)
#---
quantile(1:10, c(0.5, NA))
#---
quantile(numeric(0)) # ulp
#---
quantile(c(a = 1, b = 2, c = 3), .5) # ulp
#---
quantile(factor(c("a", "b", "c"), ordered = TRUE), .5, type = 1)
#---
quantile(factor(c("a", "b", "c"), ordered = TRUE), .5) # ulp
#---
quantile(as.Date("2020-01-01") + 0:10, c(.25, .75), type = 1)
#---
quantile(1:10, 0.5, type = 10)
#---
quantile(c(1, Inf, 3)) # ulp
#---
quantile(1:10, digits = 2, probs = 1/3) # ulp
#---
quantile(1:10, 1e-10) # ulp
#---
quantile(c(TRUE, FALSE, TRUE)) # ulp
#---
quantile(1:10, probs = c(0.1, 0.2), names = TRUE, type = 6) # ulp
#---
quantile(c(1, NaN, 3), na.rm = TRUE) # ulp
#---
quantile(5) # ulp
#---
quantile(c(1, 2), 0.5, type = 3)
#---
quantile(1:10, c(0.1, 0.3, 0.7, 0.9), type = 1)
#---
quantile(1:20, c(0.1, 0.3, 0.7, 0.9, 0.35, 0.55), type = 2)
#---
quantile(1:10, seq(0, 1, 0.1), type = 3)
#---
quantile(c(2, 4, 4, 5, 7, 9, 12), seq(0, 1, 0.05)) # ulp
#---
fivenum(c(1, 2, 3, 4, 5, 6))
#---
fivenum(c(1, 2, 3, 4, 5))
#---
fivenum(c(1, NA, 3))
#---
fivenum(c(1, NA, 3), na.rm = FALSE)
#---
fivenum(numeric(0))
#---
fivenum(1:10)
#---
fivenum(c(5))
#---
fivenum(c(1, Inf))
#---
fivenum(c(7, 1, 9, 3, 3, 8, 2))
#---
IQR(1:10) # ulp
#---
IQR(c(1, 2, 3, 10), type = 6) # ulp
#---
IQR(c(1, NA))
#---
IQR(c(1, NA, 5, 9), na.rm = TRUE)
#---
IQR(factor(1:3)) # ulp
#---
mad(c(1, 2, 3, 4, 100)) # ulp
#---
mad(c(1, 2, 3, 4, 100), constant = 1) # ulp
#---
mad(c(1, 2, 3, 4, 100), center = 2) # ulp
#---
mad(c(1, 2, 3, 4, 100, 5), low = TRUE) # ulp
#---
mad(c(1, 2, 3, 4, 100, 5), high = TRUE) # ulp
#---
mad(c(1, NA, 3))
#---
mad(c(1, NA, 3), na.rm = TRUE)
#---
mad(numeric(0)) # ulp
#---
mad(c(1, 2), low = TRUE, high = TRUE) # ulp
#---
weighted.mean(c(1, 2, 3), c(3, 2, 1)) # ulp
#---
weighted.mean(c(1, 2, 3)) # ulp
#---
weighted.mean(c(1, NA, 3), c(3, 2, 1))
#---
weighted.mean(c(1, NA, 3), c(3, 2, 1), na.rm = TRUE)
#---
weighted.mean(c(1, 2, 3), c(3, NA, 1), na.rm = TRUE)
#---
weighted.mean(1:3, 1:2) # ulp
#---
weighted.mean(c(1, 2, 3), c(0, 0, 0)) # ulp
#---
weighted.mean(c(1L, 2L, 3L), c(1L, 1L, 1L)) # ulp
#---
weighted.mean(c(1, 2, 3), c(-1, 2, 1)) # ulp
#---
weighted.mean(c(TRUE, FALSE), c(1, 3)) # ulp
#---
weighted.mean(c(1, 2, 3), c(1, 1, 1), na.rm = TRUE) # ulp
#---
weighted.mean(numeric(0), numeric(0)) # ulp
#---
weighted.mean(c(1, Inf), c(1, 1)) # ulp
#---
weighted.mean(1:3, c(a = 1, b = 2, c = 3)) # ulp
#---
quantile(c(1, 3, 4, 8, 10, 15), c(0, .1, .25, 1/3, .5, .62, .9, 1), type = 1)
#---
quantile(c(1L, 3L, 4L, 8L, 10L), c(0, .3, .5, .77, 1), type = 1)
#---
quantile(c(1, 3, 4, 8, 10, 15), c(0, .1, .25, 1/3, .5, .62, .9, 1), type = 2)
#---
quantile(c(1L, 3L, 4L, 8L, 10L), c(0, .3, .5, .77, 1), type = 2)
#---
quantile(c(1, 3, 4, 8, 10, 15), c(0, .1, .25, 1/3, .5, .62, .9, 1), type = 3)
#---
quantile(c(1L, 3L, 4L, 8L, 10L), c(0, .3, .5, .77, 1), type = 3)
#---
quantile(c(1, 3, 4, 8, 10, 15), c(0, .1, .25, 1/3, .5, .62, .9, 1), type = 4) # ulp
#---
quantile(c(1L, 3L, 4L, 8L, 10L), c(0, .3, .5, .77, 1), type = 4) # ulp
#---
quantile(c(1, 3, 4, 8, 10, 15), c(0, .1, .25, 1/3, .5, .62, .9, 1), type = 5) # ulp
#---
quantile(c(1L, 3L, 4L, 8L, 10L), c(0, .3, .5, .77, 1), type = 5) # ulp
#---
quantile(c(1, 3, 4, 8, 10, 15), c(0, .1, .25, 1/3, .5, .62, .9, 1), type = 6) # ulp
#---
quantile(c(1L, 3L, 4L, 8L, 10L), c(0, .3, .5, .77, 1), type = 6) # ulp
#---
quantile(c(1, 3, 4, 8, 10, 15), c(0, .1, .25, 1/3, .5, .62, .9, 1), type = 7) # ulp
#---
quantile(c(1L, 3L, 4L, 8L, 10L), c(0, .3, .5, .77, 1), type = 7) # ulp
#---
quantile(c(1, 3, 4, 8, 10, 15), c(0, .1, .25, 1/3, .5, .62, .9, 1), type = 8) # ulp
#---
quantile(c(1L, 3L, 4L, 8L, 10L), c(0, .3, .5, .77, 1), type = 8) # ulp
#---
quantile(c(1, 3, 4, 8, 10, 15), c(0, .1, .25, 1/3, .5, .62, .9, 1), type = 9) # ulp
#---
quantile(c(1L, 3L, 4L, 8L, 10L), c(0, .3, .5, .77, 1), type = 9) # ulp
#---
median.default(c(3, 1, 2, 10))
