d <- data.frame(a = c("x", "y", "x", "y", "x"), b = c("u", "u", "v", "v", "v"), c = c("p", "q", "p", "p", "q")); ftable(table(d))
#---
d <- data.frame(a = c("x", "y", "x", "y", "x"), b = c("u", "u", "v", "v", "v"), c = c("p", "q", "p", "p", "q")); ftable(table(d), row.vars = 1)
#---
d <- data.frame(a = c("x", "y", "x", "y", "x"), b = c("u", "u", "v", "v", "v"), c = c("p", "q", "p", "p", "q")); ftable(table(d), row.vars = c("c", "a"))
#---
d <- data.frame(a = c("x", "y", "x", "y", "x"), b = c("u", "u", "v", "v", "v"), c = c("p", "q", "p", "p", "q")); ftable(table(d), col.vars = c(1, 3))
#---
d <- data.frame(a = c("x", "y", "x", "y", "x"), b = c("u", "u", "v", "v", "v"), c = c("p", "q", "p", "p", "q")); ftable(a ~ b + c, data = d)
#---
d <- data.frame(a = c("x", "y", "x", "y", "x"), b = c("u", "u", "v", "v", "v"), c = c("p", "q", "p", "p", "q")); ftable(b ~ a, data = d)
#---
d <- data.frame(a = c("x", "y", "x", "y", "x"), b = c("u", "u", "v", "v", "v"), c = c("p", "q", "p", "p", "q")); ftable(d)
#---
d <- data.frame(a = c("x", "y", "x", "y", "x"), b = c("u", "u", "v", "v", "v")); ftable(d$a, d$b)
#---
TT <- as.table(array(c(0,0,35,0,0,0,17,0,118,154,387,670,4,13,89,3,5,11,13,0,1,13,14,0,57,14,75,192,140,80,76,20), c(4,2,2,2), dimnames = list(Class = c("1st","2nd","3rd","Crew"), Sex = c("Male","Female"), Age = c("Child","Adult"), Survived = c("No","Yes")))); ftable(TT, row.vars = 1:2)
#---
TT <- as.table(array(c(0,0,35,0,0,0,17,0,118,154,387,670,4,13,89,3,5,11,13,0,1,13,14,0,57,14,75,192,140,80,76,20), c(4,2,2,2), dimnames = list(Class = c("1st","2nd","3rd","Crew"), Sex = c("Male","Female"), Age = c("Child","Adult"), Survived = c("No","Yes")))); ftable(TT, row.vars = 1:3)
#---
TT <- as.table(array(c(0,0,35,0,0,0,17,0,118,154,387,670,4,13,89,3,5,11,13,0,1,13,14,0,57,14,75,192,140,80,76,20), c(4,2,2,2), dimnames = list(Class = c("1st","2nd","3rd","Crew"), Sex = c("Male","Female"), Age = c("Child","Adult"), Survived = c("No","Yes")))); ftable(TT, row.vars = 1:2, col.vars = "Survived")
#---
TT <- as.table(array(c(0,0,35,0,0,0,17,0,118,154,387,670,4,13,89,3,5,11,13,0,1,13,14,0,57,14,75,192,140,80,76,20), c(4,2,2,2), dimnames = list(Class = c("1st","2nd","3rd","Crew"), Sex = c("Male","Female"), Age = c("Child","Adult"), Survived = c("No","Yes")))); x <- ftable(TT, row.vars = 1:2); dim(x); attributes(x)$row.vars; attributes(x)$col.vars
#---
TT <- as.table(array(c(0,0,35,0,0,0,17,0,118,154,387,670,4,13,89,3,5,11,13,0,1,13,14,0,57,14,75,192,140,80,76,20), c(4,2,2,2), dimnames = list(Class = c("1st","2nd","3rd","Crew"), Sex = c("Male","Female"), Age = c("Child","Adult"), Survived = c("No","Yes")))); x <- ftable(TT, row.vars = 1:2); class(x); is.matrix(x); unclass(x)[1:3, ]
#---
TT <- as.table(array(c(0,0,35,0,0,0,17,0,118,154,387,670,4,13,89,3,5,11,13,0,1,13,14,0,57,14,75,192,140,80,76,20), c(4,2,2,2), dimnames = list(Class = c("1st","2nd","3rd","Crew"), Sex = c("Male","Female"), Age = c("Child","Adult"), Survived = c("No","Yes")))); x <- ftable(TT, row.vars = 1:2); format(x)
#---
TT <- as.table(array(c(0,0,35,0,0,0,17,0,118,154,387,670,4,13,89,3,5,11,13,0,1,13,14,0,57,14,75,192,140,80,76,20), c(4,2,2,2), dimnames = list(Class = c("1st","2nd","3rd","Crew"), Sex = c("Male","Female"), Age = c("Child","Adult"), Survived = c("No","Yes")))); x <- ftable(TT, row.vars = 1:2); format(x, quote = FALSE)
#---
TT <- as.table(array(c(0,0,35,0,0,0,17,0,118,154,387,670,4,13,89,3,5,11,13,0,1,13,14,0,57,14,75,192,140,80,76,20), c(4,2,2,2), dimnames = list(Class = c("1st","2nd","3rd","Crew"), Sex = c("Male","Female"), Age = c("Child","Adult"), Survived = c("No","Yes")))); x <- ftable(TT, row.vars = 1:2); format(x, method = "col.compact")
#---
TT <- as.table(array(c(0,0,35,0,0,0,17,0,118,154,387,670,4,13,89,3,5,11,13,0,1,13,14,0,57,14,75,192,140,80,76,20), c(4,2,2,2), dimnames = list(Class = c("1st","2nd","3rd","Crew"), Sex = c("Male","Female"), Age = c("Child","Adult"), Survived = c("No","Yes")))); x <- ftable(TT, row.vars = 1:2); format(x, method = "row.compact")
#---
TT <- as.table(array(c(0,0,35,0,0,0,17,0,118,154,387,670,4,13,89,3,5,11,13,0,1,13,14,0,57,14,75,192,140,80,76,20), c(4,2,2,2), dimnames = list(Class = c("1st","2nd","3rd","Crew"), Sex = c("Male","Female"), Age = c("Child","Adult"), Survived = c("No","Yes")))); x <- ftable(TT, row.vars = 1:2); format(x, method = "compact")[, 1:5]
#---
TT <- as.table(array(c(0,0,35,0,0,0,17,0,118,154,387,670,4,13,89,3,5,11,13,0,1,13,14,0,57,14,75,192,140,80,76,20), c(4,2,2,2), dimnames = list(Class = c("1st","2nd","3rd","Crew"), Sex = c("Male","Female"), Age = c("Child","Adult"), Survived = c("No","Yes")))); x <- ftable(TT, row.vars = 1:2); write.ftable(x, quote = FALSE)
#---
TT <- as.table(array(c(0,0,35,0,0,0,17,0,118,154,387,670,4,13,89,3,5,11,13,0,1,13,14,0,57,14,75,192,140,80,76,20), c(4,2,2,2), dimnames = list(Class = c("1st","2nd","3rd","Crew"), Sex = c("Male","Female"), Age = c("Child","Adult"), Survived = c("No","Yes")))); x <- ftable(TT, row.vars = 1:2); format(x, method = "compact")[, 5:6]
#---
TT <- as.table(array(c(0,0,35,0,0,0,17,0,118,154,387,670,4,13,89,3,5,11,13,0,1,13,14,0,57,14,75,192,140,80,76,20), c(4,2,2,2), dimnames = list(Class = c("1st","2nd","3rd","Crew"), Sex = c("Male","Female"), Age = c("Child","Adult"), Survived = c("No","Yes")))); x <- ftable(TT, row.vars = 1:2); as.table(x)[, , 1, 1]
#---
TT <- as.table(array(c(0,0,35,0,0,0,17,0,118,154,387,670,4,13,89,3,5,11,13,0,1,13,14,0,57,14,75,192,140,80,76,20), c(4,2,2,2), dimnames = list(Class = c("1st","2nd","3rd","Crew"), Sex = c("Male","Female"), Age = c("Child","Adult"), Survived = c("No","Yes")))); x <- ftable(TT, row.vars = 1:2); as.data.frame(x)[1:4, ]
#---
TT <- as.table(array(c(0,0,35,0,0,0,17,0,118,154,387,670,4,13,89,3,5,11,13,0,1,13,14,0,57,14,75,192,140,80,76,20), c(4,2,2,2), dimnames = list(Class = c("1st","2nd","3rd","Crew"), Sex = c("Male","Female"), Age = c("Child","Adult"), Survived = c("No","Yes")))); x <- ftable(TT, row.vars = 1:2); t(x)[1:2, ]
#---
TT <- as.table(array(c(0,0,35,0,0,0,17,0,118,154,387,670,4,13,89,3,5,11,13,0,1,13,14,0,57,14,75,192,140,80,76,20), c(4,2,2,2), dimnames = list(Class = c("1st","2nd","3rd","Crew"), Sex = c("Male","Female"), Age = c("Child","Adult"), Survived = c("No","Yes")))); x <- ftable(TT, row.vars = 1:2); x[1:2, ]
#---
TT <- as.table(array(c(0,0,35,0,0,0,17,0,118,154,387,670,4,13,89,3,5,11,13,0,1,13,14,0,57,14,75,192,140,80,76,20), c(4,2,2,2), dimnames = list(Class = c("1st","2nd","3rd","Crew"), Sex = c("Male","Female"), Age = c("Child","Adult"), Survived = c("No","Yes")))); x <- ftable(TT, row.vars = 1:2); x[1:2, 1]
#---
TT <- as.table(array(c(0,0,35,0,0,0,17,0,118,154,387,670,4,13,89,3,5,11,13,0,1,13,14,0,57,14,75,192,140,80,76,20), c(4,2,2,2), dimnames = list(Class = c("1st","2nd","3rd","Crew"), Sex = c("Male","Female"), Age = c("Child","Adult"), Survived = c("No","Yes")))); x <- ftable(TT, row.vars = 1:2); x[1:2, , drop = FALSE]
#---
TT <- as.table(array(c(0,0,35,0,0,0,17,0,118,154,387,670,4,13,89,3,5,11,13,0,1,13,14,0,57,14,75,192,140,80,76,20), c(4,2,2,2), dimnames = list(Class = c("1st","2nd","3rd","Crew"), Sex = c("Male","Female"), Age = c("Child","Adult"), Survived = c("No","Yes")))); x <- ftable(TT, row.vars = 1:2); x[, 1:2]
#---
TT <- as.table(array(c(0,0,35,0,0,0,17,0,118,154,387,670,4,13,89,3,5,11,13,0,1,13,14,0,57,14,75,192,140,80,76,20), c(4,2,2,2), dimnames = list(Class = c("1st","2nd","3rd","Crew"), Sex = c("Male","Female"), Age = c("Child","Adult"), Survived = c("No","Yes")))); x <- ftable(TT, row.vars = 1:2); is.na(x[1:2, 1])
#---
ftable(1:3)
#---
ftable(c("a", "b"), c("u", "v"))
#---
ftable(matrix(1:4, 2))
#---
x <- table(a = c(1, 2, 1), b = c("u", "v", "v")); ftable(x, row.vars = 3)
#---
x <- table(a = c(1, 2, 1), b = c("u", "v", "v")); ftable(x, row.vars = "z")
#---
x <- table(a = c(1, 2, 1), b = c("u", "v", "v")); ftable(x, row.vars = 1, col.vars = 1)
#---
x <- table(c(1, 2, 1), c("u", "v", "v")); ftable(x)
#---
x <- table(a = c(1, 2, 1), b = c("u", "v", "v")); r <- ftable(x, exclude = "u"); r
#---
TT <- as.table(array(c(0,0,35,0,0,0,17,0,118,154,387,670,4,13,89,3,5,11,13,0,1,13,14,0,57,14,75,192,140,80,76,20), c(4,2,2,2), dimnames = list(Class = c("1st","2nd","3rd","Crew"), Sex = c("Male","Female"), Age = c("Child","Adult"), Survived = c("No","Yes")))); ftable(TT, row.vars = 1, col.vars = 2:4)
#---
TT <- as.table(array(c(0,0,35,0,0,0,17,0,118,154,387,670,4,13,89,3,5,11,13,0,1,13,14,0,57,14,75,192,140,80,76,20), c(4,2,2,2), dimnames = list(Class = c("1st","2nd","3rd","Crew"), Sex = c("Male","Female"), Age = c("Child","Adult"), Survived = c("No","Yes")))); x <- ftable(TT, row.vars = 1:2); x2 <- ftable(x, row.vars = 1); x2
#---
TT <- as.table(array(c(0,0,35,0,0,0,17,0,118,154,387,670,4,13,89,3,5,11,13,0,1,13,14,0,57,14,75,192,140,80,76,20), c(4,2,2,2), dimnames = list(Class = c("1st","2nd","3rd","Crew"), Sex = c("Male","Female"), Age = c("Child","Adult"), Survived = c("No","Yes")))); x <- ftable(TT, row.vars = 1:2); as.matrix(x)[1:2, ]
#---
TT <- as.table(array(c(0,0,35,0,0,0,17,0,118,154,387,670,4,13,89,3,5,11,13,0,1,13,14,0,57,14,75,192,140,80,76,20), c(4,2,2,2), dimnames = list(Class = c("1st","2nd","3rd","Crew"), Sex = c("Male","Female"), Age = c("Child","Adult"), Survived = c("No","Yes")))); x <- ftable(TT, row.vars = 1:2); as.matrix(x, sep = "_")[1:2, ]
#---
TT <- as.table(array(c(0,0,35,0,0,0,17,0,118,154,387,670,4,13,89,3,5,11,13,0,1,13,14,0,57,14,75,192,140,80,76,20), c(4,2,2,2), dimnames = list(Class = c("1st","2nd","3rd","Crew"), Sex = c("Male","Female"), Age = c("Child","Adult"), Survived = c("No","Yes")))); x <- ftable(TT, row.vars = 1:2); rbind(x[1, ], x[2, ])
#---
TT <- as.table(array(c(0,0,35,0,0,0,17,0,118,154,387,670,4,13,89,3,5,11,13,0,1,13,14,0,57,14,75,192,140,80,76,20), c(4,2,2,2), dimnames = list(Class = c("1st","2nd","3rd","Crew"), Sex = c("Male","Female"), Age = c("Child","Adult"), Survived = c("No","Yes")))); x <- ftable(TT, row.vars = 1:2); ftable(x[1:3, ])
#---
ftable(table(a = c(1, 2, 1, 2), b = c("u", "u", "v", "v")), row.vars = 1, col.vars = 2)
#---
TT <- as.table(array(c(0,0,35,0,0,0,17,0,118,154,387,670,4,13,89,3,5,11,13,0,1,13,14,0,57,14,75,192,140,80,76,20), c(4,2,2,2), dimnames = list(Class = c("1st","2nd","3rd","Crew"), Sex = c("Male","Female"), Age = c("Child","Adult"), Survived = c("No","Yes")))); x <- ftable(TT); print(x, quote = TRUE)
#---
TT <- as.table(array(c(0,0,35,0,0,0,17,0,118,154,387,670,4,13,89,3,5,11,13,0,1,13,14,0,57,14,75,192,140,80,76,20), c(4,2,2,2), dimnames = list(Class = c("1st","2nd","3rd","Crew"), Sex = c("Male","Female"), Age = c("Child","Adult"), Survived = c("No","Yes")))); x <- ftable(TT); print(x, digits = 2)
#---
TT <- as.table(array(c(0,0,35,0,0,0,17,0,118,154,387,670,4,13,89,3,5,11,13,0,1,13,14,0,57,14,75,192,140,80,76,20), c(4,2,2,2), dimnames = list(Class = c("1st","2nd","3rd","Crew"), Sex = c("Male","Female"), Age = c("Child","Adult"), Survived = c("No","Yes")))); x <- ftable(prop.table(TT, 1)); format(x, digits = 2)[1:4, ]
#---
TT <- as.table(array(c(0,0,35,0,0,0,17,0,118,154,387,670,4,13,89,3,5,11,13,0,1,13,14,0,57,14,75,192,140,80,76,20), c(4,2,2,2), dimnames = list(Class = c("1st","2nd","3rd","Crew"), Sex = c("Male","Female"), Age = c("Child","Adult"), Survived = c("No","Yes")))); x <- ftable(prop.table(TT, 1), row.vars = 1:2); format(x, digits = 2, justify = "right")[1:4, ]
#---
TT <- as.table(array(c(0,0,35,0,0,0,17,0,118,154,387,670,4,13,89,3,5,11,13,0,1,13,14,0,57,14,75,192,140,80,76,20), c(4,2,2,2), dimnames = list(Class = c("1st","2nd","3rd","Crew"), Sex = c("Male","Female"), Age = c("Child","Adult"), Survived = c("No","Yes")))); x <- ftable(TT, row.vars = 1:2); format(x, justify = "left")[1:3, ]
#---
TT <- as.table(array(c(0,0,35,0,0,0,17,0,118,154,387,670,4,13,89,3,5,11,13,0,1,13,14,0,57,14,75,192,140,80,76,20), c(4,2,2,2), dimnames = list(Class = c("1st","2nd","3rd","Crew"), Sex = c("Male","Female"), Age = c("Child","Adult"), Survived = c("No","Yes")))); x <- ftable(TT, row.vars = 1:2); format(x, lsep = "|")[1:3, ]
#---
TT <- as.table(array(c(0,0,35,0,0,0,17,0,118,154,387,670,4,13,89,3,5,11,13,0,1,13,14,0,57,14,75,192,140,80,76,20), c(4,2,2,2), dimnames = list(Class = c("1st","2nd","3rd","Crew"), Sex = c("Male","Female"), Age = c("Child","Adult"), Survived = c("No","Yes")))); x <- ftable(TT, row.vars = 1:2); format(x, lsep = "")[1:3, ]
