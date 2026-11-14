table(c("a", "b", "a", "c", "a"))
#---
x <- c(3, 1, 3, 2, 3); table(x)
#---
x <- c(3, 1, 3, 2, 3); attributes(table(x))
#---
x <- c(3, 1, NA, 2, 3); table(x)
#---
x <- c(3, 1, NA, 2, 3); table(x, useNA = "ifany")
#---
x <- c(3, 1, 2, 3); table(x, useNA = "ifany")
#---
x <- c(3, 1, 2, 3); table(x, useNA = "always")
#---
x <- c("a", "b", "a"); table(x, exclude = "a")
#---
x <- c(3, 1, NA, 2, 3); table(x, exclude = NULL)
#---
x <- c("a", "b", "a"); y <- c("u", "u", "v"); table(x, y)
#---
x <- c("a", "b", "a"); y <- c("u", "u", "v"); attributes(table(x, y))
#---
table(c("a", "b", "a"), c("u", "u", "v"))
#---
table(g = c("a", "b", "a"), h = c("u", "u", "v"))
#---
x <- c("a", "b", "a"); table(x, dnn = "letters")
#---
x <- c("a", "b", "a"); y <- c(1, 1, 2); z <- c("p", "q", "p"); table(x, y, z)
#---
f <- factor(c("a", "b", "a"), levels = c("a", "b", "z")); table(f)
#---
f <- factor(c("a", "b", NA), levels = c("a", "b")); table(f, useNA = "ifany")
#---
table(data.frame(a = c(1, 1, 2), b = c("u", "v", "u")))
#---
table(character(0))
#---
table(c(TRUE, FALSE, TRUE))
#---
t1 <- table(c("a", "b", "a")); names(t1); t1[["a"]]; t1["b"]; sum(t1); dim(t1); dimnames(t1)
#---
t1 <- table(c("a", "b", "a")); as.vector(t1)
#---
t1 <- table(c("a", "b", "a")); as.data.frame(t1)
#---
t1 <- table(c("a", "b", "a")); attributes(as.data.frame(t1))
#---
t2 <- table(x = c("a", "b", "a"), y = c("u", "u", "v")); as.data.frame(t2)
#---
t2 <- table(x = c("a", "b", "a"), y = c("u", "u", "v")); as.data.frame(t2, responseName = "n", stringsAsFactors = FALSE)
#---
t2 <- table(x = c("a", "b", "a"), y = c("u", "u", "v")); t2[2, ]
#---
t2 <- table(x = c("a", "b", "a"), y = c("u", "u", "v")); prop.table(t2, 1)
#---
t2 <- table(x = c("a", "b", "a"), y = c("u", "u", "v")); margin.table(t2, 1)
#---
tabulate(c(1, 2, 2, 5), 5)
#---
expand.grid(a = 1:2, b = c("x", "y"))
#---
attributes(expand.grid(a = 1:2, b = c("x", "y")))
#---
expand.grid(1:2, c("x", "y"), KEEP.OUT.ATTRS = FALSE, stringsAsFactors = FALSE)
#---
nrow(expand.grid(a = 1:3, b = 1:4, c = 1:2))
#---
interaction(c("a", "b", "a"), c("u", "u", "v"))
#---
interaction(c("a", "b", "a"), c("u", "u", "v"), lex.order = TRUE)
#---
interaction(c("a", "b", "a"), c("u", "u", "v"), drop = TRUE, sep = "_")
#---
cut(c(1, 5, 10), breaks = c(0, 5, 10))
#---
cut(c(1, 5, 10), breaks = c(0, 5, 10), right = FALSE)
#---
cut(c(0, 5, 10), breaks = c(0, 5, 10), include.lowest = TRUE)
#---
cut(c(0, 5, 10), breaks = c(0, 5, 10), include.lowest = TRUE, right = FALSE)
#---
cut(c(1, 5, 10), breaks = c(0, 5, 10), labels = c("lo", "hi"))
#---
cut(c(1, 5, 10), breaks = c(0, 5, 10), labels = FALSE)
#---
cut(c(1, 2, 3, 4, 5, 6), 3)
#---
cut(c(1, 2, 3, 4, 5, 6), breaks = 2, dig.lab = 2)
#---
cut(c(3, 3, 3), 2)
#---
cut(c(0.5, 2.5, 100000, 1e6), breaks = c(0, 1, 1e5, 1e6, 1e7))
#---
cut(c(1, NA, 7), breaks = c(0, 5, 10), ordered_result = TRUE)
#---
cut(c(1, 5), breaks = c(0, 5, 5, 10))
#---
cut("a", 2)
#---
table(cut(c(1, 2, 3, 4, 5, 6), 3))
#---
cut(c(-Inf, 1, Inf), breaks = c(-Inf, 0, Inf))
#---
nlevels(interaction(factor(c("a", "b")), factor(c("x", "y", "z"))))
#---
x <- table(c("a", "b")); class(x)
