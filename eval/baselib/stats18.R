# Base library, slice 18 (3): the generics of the descriptive statistics.
# Written from the help pages of R 4.6.2 (median, quantile, weighted.mean);
# the methods are Go natives (eval/stats18.go, eval/stats18b.go).

median <- function(x, na.rm = FALSE, ...) UseMethod("median")

quantile <- function(x, ...) UseMethod("quantile")

weighted.mean <- function(x, w, ...) UseMethod("weighted.mean")
