d <- data.frame(g = c("b","a","b","a","c"), h = c("x","x","y","y","x"), v = 1:5, w = c(2.5, 3.5, NA, 1, 2), stringsAsFactors = FALSE)
aggregate(d$v, list(d$g), sum)
#---
d <- data.frame(g = c("b","a","b","a","c"), h = c("x","x","y","y","x"), v = 1:5, w = c(2.5, 3.5, NA, 1, 2), stringsAsFactors = FALSE)
aggregate(d$v, list(grp = d$g), sum)
#---
d <- data.frame(g = c("b","a","b","a","c"), h = c("x","x","y","y","x"), v = 1:5, w = c(2.5, 3.5, NA, 1, 2), stringsAsFactors = FALSE)
aggregate(d[c("v","w")], list(g = d$g), sum)
#---
d <- data.frame(g = c("b","a","b","a","c"), h = c("x","x","y","y","x"), v = 1:5, w = c(2.5, 3.5, NA, 1, 2), stringsAsFactors = FALSE)
aggregate(d[c("v","w")], list(g = d$g), sum, na.rm = TRUE)
#---
d <- data.frame(g = c("b","a","b","a","c"), h = c("x","x","y","y","x"), v = 1:5, w = c(2.5, 3.5, NA, 1, 2), stringsAsFactors = FALSE)
aggregate(d[c("v","w")], list(g = d$g, h = d$h), sum)
#---
d <- data.frame(g = c("b","a","b","a","c"), h = c("x","x","y","y","x"), v = 1:5, w = c(2.5, 3.5, NA, 1, 2), stringsAsFactors = FALSE)
aggregate(d[c("v","w")], list(g = d$g, h = d$h), sum, drop = FALSE)
#---
d <- data.frame(g = c("b","a","b","a","c"), h = c("x","x","y","y","x"), v = 1:5, w = c(2.5, 3.5, NA, 1, 2), stringsAsFactors = FALSE)
aggregate(v ~ g, d, sum)
#---
d <- data.frame(g = c("b","a","b","a","c"), h = c("x","x","y","y","x"), v = 1:5, w = c(2.5, 3.5, NA, 1, 2), stringsAsFactors = FALSE)
aggregate(cbind(v, w) ~ g, d, sum)
#---
d <- data.frame(g = c("b","a","b","a","c"), h = c("x","x","y","y","x"), v = 1:5, w = c(2.5, 3.5, NA, 1, 2), stringsAsFactors = FALSE)
aggregate(. ~ g, d[c("g","v","w")], mean)
#---
d <- data.frame(g = c("b","a","b","a","c"), h = c("x","x","y","y","x"), v = 1:5, w = c(2.5, 3.5, NA, 1, 2), stringsAsFactors = FALSE)
aggregate(v ~ g + h, d, sum)
#---
d <- data.frame(g = c("b","a","b","a","c"), h = c("x","x","y","y","x"), v = 1:5, w = c(2.5, 3.5, NA, 1, 2), stringsAsFactors = FALSE)
aggregate(v ~ ., d[c("g","h","v")], sum)
#---
d <- data.frame(g = c("b","a","b","a","c"), h = c("x","x","y","y","x"), v = 1:5, w = c(2.5, 3.5, NA, 1, 2), stringsAsFactors = FALSE)
aggregate(w ~ g, d, sum)
#---
d <- data.frame(g = c("b","a","b","a","c"), h = c("x","x","y","y","x"), v = 1:5, w = c(2.5, 3.5, NA, 1, 2), stringsAsFactors = FALSE)
aggregate(w ~ g, d, sum, na.action = na.pass)
#---
d <- data.frame(g = c("b","a","b","a","c"), h = c("x","x","y","y","x"), v = 1:5, w = c(2.5, 3.5, NA, 1, 2), stringsAsFactors = FALSE)
aggregate(v ~ g, d, length)
#---
d <- data.frame(g = c("b","a","b","a","c"), h = c("x","x","y","y","x"), v = 1:5, w = c(2.5, 3.5, NA, 1, 2), stringsAsFactors = FALSE)
aggregate(v ~ 1, d, sum)
#---
d <- data.frame(g = c("b","a","b","a","c"), h = c("x","x","y","y","x"), v = 1:5, w = c(2.5, 3.5, NA, 1, 2), stringsAsFactors = FALSE)
aggregate(x = d$v, by = list(d$g), FUN = mean)
#---
d <- data.frame(g = c("b","a","b","a","c"), h = c("x","x","y","y","x"), v = 1:5, w = c(2.5, 3.5, NA, 1, 2), stringsAsFactors = FALSE)
aggregate(v ~ g, d, sum, subset = v > 1)
#---
d <- data.frame(g = c("b","a","b","a","c"), h = c("x","x","y","y","x"), v = 1:5, w = c(2.5, 3.5, NA, 1, 2), stringsAsFactors = FALSE)
aggregate(d["v"], list(f = factor(d$g, levels = c("c","b","a"))), sum)
#---
d <- data.frame(g = c("b","a","b","a","c"), h = c("x","x","y","y","x"), v = 1:5, w = c(2.5, 3.5, NA, 1, 2), stringsAsFactors = FALSE)
aggregate(d["v"], list(f = factor(d$g, levels = c("c","b","a","z"))), sum)
#---
d <- data.frame(g = c("b","a","b","a","c"), h = c("x","x","y","y","x"), v = 1:5, w = c(2.5, 3.5, NA, 1, 2), stringsAsFactors = FALSE)
aggregate(d["v"], list(f = factor(d$g, levels = c("c","b","a","z"))), sum, drop = FALSE)
#---
d <- data.frame(g = c("b","a","b","a","c"), h = c("x","x","y","y","x"), v = 1:5, w = c(2.5, 3.5, NA, 1, 2), stringsAsFactors = FALSE)
aggregate(1:5, list(c(1,1,2,2,NA)), sum)
#---
d <- data.frame(g = c("b","a","b","a","c"), h = c("x","x","y","y","x"), v = 1:5, w = c(2.5, 3.5, NA, 1, 2), stringsAsFactors = FALSE)
aggregate(d$v, list(d$g), sum, simplify = FALSE)
#---
d <- data.frame(g = c("b","a","b","a","c"), h = c("x","x","y","y","x"), v = 1:5, w = c(2.5, 3.5, NA, 1, 2), stringsAsFactors = FALSE)
class(aggregate(v ~ g, d, sum))
#---
d <- data.frame(g = c("b","a","b","a","c"), h = c("x","x","y","y","x"), v = 1:5, w = c(2.5, 3.5, NA, 1, 2), stringsAsFactors = FALSE)
names(aggregate(cbind(vv = v, ww = w) ~ g, d, sum))
#---
d <- data.frame(g = c("b","a","b","a","c"), h = c("x","x","y","y","x"), v = 1:5, w = c(2.5, 3.5, NA, 1, 2), stringsAsFactors = FALSE)
aggregate(d$v, list(d$g), paste, collapse = "")
#---
d <- data.frame(g = c("b","a","b","a","c"), h = c("x","x","y","y","x"), v = 1:5, w = c(2.5, 3.5, NA, 1, 2), stringsAsFactors = FALSE)
aggregate(d$g, list(d$h), function(x) x[1])
#---
d <- data.frame(g = c("b","a","b","a","c"), h = c("x","x","y","y","x"), v = 1:5, w = c(2.5, 3.5, NA, 1, 2), stringsAsFactors = FALSE)
aggregate(d["v"], list(d$g), max)
#---
d <- data.frame(g = c("b","a","b","a","c"), h = c("x","x","y","y","x"), v = 1:5, w = c(2.5, 3.5, NA, 1, 2), stringsAsFactors = FALSE)
aggregate(d[c("v","w")], list(d$g), mean)
#---
d <- data.frame(g = c("b","a","b","a","c"), h = c("x","x","y","y","x"), v = 1:5, w = c(2.5, 3.5, NA, 1, 2), stringsAsFactors = FALSE)
tryCatch(aggregate(d$v, list(1:3), sum), error = function(e) conditionMessage(e))
#---
d <- data.frame(g = c("b","a","b","a","c"), h = c("x","x","y","y","x"), v = 1:5, w = c(2.5, 3.5, NA, 1, 2), stringsAsFactors = FALSE)
tryCatch(aggregate(d$v, d$g, sum), error = function(e) conditionMessage(e))
#---
d <- data.frame(g = c("b","a","b","a","c"), h = c("x","x","y","y","x"), v = 1:5, w = c(2.5, 3.5, NA, 1, 2), stringsAsFactors = FALSE)
tryCatch(aggregate(v ~ zz, d, sum), error = function(e) conditionMessage(e))
#---
d <- data.frame(g = c("b","a","b","a","c"), h = c("x","x","y","y","x"), v = 1:5, w = c(2.5, 3.5, NA, 1, 2), stringsAsFactors = FALSE)
tryCatch(aggregate(d$v, list(d$g)), error = function(e) conditionMessage(e))
#---
d <- data.frame(g = c("b","a","b","a","c"), h = c("x","x","y","y","x"), v = 1:5, w = c(2.5, 3.5, NA, 1, 2), stringsAsFactors = FALSE)
aggregate(d$v, list(as.Date("2020-01-01") + c(0,0,1,1,1)), sum)
#---
d <- data.frame(g = c("b","a","b","a","c"), h = c("x","x","y","y","x"), v = 1:5, w = c(2.5, 3.5, NA, 1, 2), stringsAsFactors = FALSE)
aggregate(matrix(1:6, 3), list(c("a","b","a")), sum)
#---
d <- data.frame(g = c("b","a","b","a","c"), h = c("x","x","y","y","x"), v = 1:5, w = c(2.5, 3.5, NA, 1, 2), stringsAsFactors = FALSE)
aggregate(1:6, list(rep(1:2, 3)), sum)
#---
d <- data.frame(g = c("b","a","b","a","c"), h = c("x","x","y","y","x"), v = 1:5, w = c(2.5, 3.5, NA, 1, 2), stringsAsFactors = FALSE)
aggregate(c(a=1,b=2,c=3), list(c("x","x","y")), sum)
#---
d <- data.frame(g = c("b","a","b","a","c"), h = c("x","x","y","y","x"), v = 1:5, w = c(2.5, 3.5, NA, 1, 2), stringsAsFactors = FALSE)
aggregate(d$w, list(d$g), mean)
#---
d <- data.frame(g = c("b","a","b","a","c"), h = c("x","x","y","y","x"), v = 1:5, w = c(2.5, 3.5, NA, 1, 2), stringsAsFactors = FALSE)
aggregate(d$w, list(d$g), mean, na.rm = TRUE)
#---
d <- data.frame(g = c("b","a","b","a","c"), h = c("x","x","y","y","x"), v = 1:5, w = c(2.5, 3.5, NA, 1, 2), stringsAsFactors = FALSE)
aggregate(d[c("v","w")], list(d$g, d$h), sum)
#---
d <- data.frame(g = c("b","a","b","a","c"), v = c(3,1,2,5,4), w = c(1,1,2,2,1))
sort_by(d, ~ v)
#---
d <- data.frame(g = c("b","a","b","a","c"), v = c(3,1,2,5,4), w = c(1,1,2,2,1))
sort_by(d, ~ g + v)
#---
d <- data.frame(g = c("b","a","b","a","c"), v = c(3,1,2,5,4), w = c(1,1,2,2,1))
sort_by(d, ~ g + v, decreasing = TRUE)
#---
d <- data.frame(g = c("b","a","b","a","c"), v = c(3,1,2,5,4), w = c(1,1,2,2,1))
sort_by(d, ~ w + v, decreasing = c(TRUE, FALSE))
#---
d <- data.frame(g = c("b","a","b","a","c"), v = c(3,1,2,5,4), w = c(1,1,2,2,1))
sort_by(d, ~ g)
#---
d <- data.frame(g = c("b","a","b","a","c"), v = c(3,1,2,5,4), w = c(1,1,2,2,1))
sort_by(d, ~ -v)
#---
d <- data.frame(g = c("b","a","b","a","c"), v = c(3,1,2,5,4), w = c(1,1,2,2,1))
sort_by(d, ~ rev(v))
#---
d <- data.frame(g = c("b","a","b","a","c"), v = c(3,1,2,5,4), w = c(1,1,2,2,1))
sort_by(d, ~ v, na.last = FALSE)
#---
d <- data.frame(g = c("b","a","b","a","c"), v = c(3,1,2,5,4), w = c(1,1,2,2,1))
tryCatch(sort_by(d, ~ zz), error = function(e) conditionMessage(e))
#---
d <- data.frame(g = c("b","a","b","a","c"), v = c(3,1,2,5,4), w = c(1,1,2,2,1))
tryCatch(sort_by(d, v ~ g), error = function(e) conditionMessage(e))
#---
d <- data.frame(g = c("b","a","b","a","c"), v = c(3,1,2,5,4), w = c(1,1,2,2,1))
sort_by(d, ~ w, method = "radix")
#---
d <- data.frame(g = c("b","a","b","a","c"), v = c(3,1,2,5,4), w = c(1,1,2,2,1))
sort_by(d, ~ g, decreasing = TRUE)
#---
d <- data.frame(g = c("b","a","b","a","c"), v = c(3,1,2,5,4), w = c(1,1,2,2,1))
sort_by(d, list(d$w, d$v))
#---
d <- data.frame(g = c("b","a","b","a","c"), v = c(3,1,2,5,4), w = c(1,1,2,2,1))
s <- stack(list(a = 1:2, b = 3:5)); s
#---
d <- data.frame(g = c("b","a","b","a","c"), v = c(3,1,2,5,4), w = c(1,1,2,2,1))
stack(data.frame(x = 1:2, y = 3:4))
#---
d <- data.frame(g = c("b","a","b","a","c"), v = c(3,1,2,5,4), w = c(1,1,2,2,1))
stack(data.frame(x = 1:2, y = c("a","b")))
#---
d <- data.frame(g = c("b","a","b","a","c"), v = c(3,1,2,5,4), w = c(1,1,2,2,1))
stack(list(a = 1, b = 2), drop = TRUE)
#---
d <- data.frame(g = c("b","a","b","a","c"), v = c(3,1,2,5,4), w = c(1,1,2,2,1))
stack(list(a = 1:2, b = 3), select = a)
#---
d <- data.frame(g = c("b","a","b","a","c"), v = c(3,1,2,5,4), w = c(1,1,2,2,1))
stack(data.frame(x = 1:2, y = 3:4, z = 5:6), select = c(x, z))
#---
d <- data.frame(g = c("b","a","b","a","c"), v = c(3,1,2,5,4), w = c(1,1,2,2,1))
stack(data.frame(x = 1:2, y = 3:4, z = 5:6), select = -y)
#---
d <- data.frame(g = c("b","a","b","a","c"), v = c(3,1,2,5,4), w = c(1,1,2,2,1))
stack(c(a = 1, b = 2))
#---
d <- data.frame(g = c("b","a","b","a","c"), v = c(3,1,2,5,4), w = c(1,1,2,2,1))
levels(stack(list(b = 1, a = 2))$ind)
#---
d <- data.frame(g = c("b","a","b","a","c"), v = c(3,1,2,5,4), w = c(1,1,2,2,1))
unstack(stack(list(a = 1:2, b = 3:4)))
#---
d <- data.frame(g = c("b","a","b","a","c"), v = c(3,1,2,5,4), w = c(1,1,2,2,1))
unstack(data.frame(values = 1:4, ind = factor(c("a","a","b","b"))))
#---
d <- data.frame(g = c("b","a","b","a","c"), v = c(3,1,2,5,4), w = c(1,1,2,2,1))
unstack(data.frame(values = 1:3, ind = factor(c("a","a","b"))))
#---
d <- data.frame(g = c("b","a","b","a","c"), v = c(3,1,2,5,4), w = c(1,1,2,2,1))
unstack(data.frame(values = 1:4, g = factor(c("a","a","b","b"))), values ~ g)
#---
d <- data.frame(g = c("b","a","b","a","c"), v = c(3,1,2,5,4), w = c(1,1,2,2,1))
unstack(data.frame(v = 1:4, f = factor(c("a","b","a","b"))), v ~ f)
#---
d <- data.frame(g = c("b","a","b","a","c"), v = c(3,1,2,5,4), w = c(1,1,2,2,1))
unstack(data.frame(v = 1:6, f = factor(c("a","b","a","b","a","b")), g = factor(c("x","x","x","x","y","y"))), v ~ f)
#---
d <- data.frame(g = c("b","a","b","a","c"), v = c(3,1,2,5,4), w = c(1,1,2,2,1))
unstack(data.frame(values = 1:4, ind = factor(c("a","a","b","b"))), form = values ~ ind)
#---
d <- data.frame(g = c("b","a","b","a","c"), v = c(3,1,2,5,4), w = c(1,1,2,2,1))
tryCatch(unstack(data.frame(a = 1)), error = function(e) conditionMessage(e))
#---
d <- data.frame(g = c("b","a","b","a","c"), v = c(3,1,2,5,4), w = c(1,1,2,2,1))
stack(list(a = 1:2, b = letters[1:2]))
#---
d <- data.frame(g = c("b","a","b","a","c"), v = c(3,1,2,5,4), w = c(1,1,2,2,1))
stack(data.frame(x = c(a = 1, b = 2)))
#---
rep.factor(factor(c("a", "b")), 2)
#---
rep(factor(c("a", "b")), c(2, 1))
#---
x <- factor(c("a", "b"), levels = c("b", "a")); rep(x, each = 2)
#---
attributes(rep(factor(c("a", "b")), 2))
#---
rep(factor(c("u", NA)), 2)
#---
rep(factor(c("a", "b"), ordered = TRUE), 2)
#---
unname(matrix(1:4, 2, dimnames = list(c("a", "b"), NULL)))
#---
deparse(formula(data.frame(a = 1)))
#---
d <- data.frame(g = c("b","a","b","a","c"), v = 1:5)
dim(aggregate(v ~ g, d, sum))
#---
d <- data.frame(g = c("b","a","b","a","c"), v = 1:5)
aggregate.data.frame(d["v"], list(g = d$g), sum)
#---
d <- data.frame(g = c("b","a"), v = 1:2)
registerS3method("stack", "zzz", function(x, ...) 1); stack(structure(list(), class = "zzz"))
#---
aggregate(cbind(mpg = c(1, 2, 3, 4), hp = c(10, 20, 30, 40)) ~ g, data.frame(g = c("a", "b", "a", "b")), mean)
#---
aggregate(x ~ g, data.frame(x = c(1, 2, 3), g = factor(c("u", "v", "u"), levels = c("v", "u", "w"))), sum)
#---
aggregate(x ~ g, data.frame(x = c(1, 2, 3), g = factor(c("u", "v", "u"), levels = c("v", "u", "w"))), sum, drop = FALSE)
#---
aggregate(x ~ g, data.frame(x = c(1, NA, 3), g = c("u", "v", "u")), sum)
#---
aggregate(x ~ g, data.frame(x = c(1, 2, 3), g = c("u", NA, "u")), sum)
#---
aggregate(x ~ g + h, data.frame(x = 1:4, g = c("a", "a", "b", "b"), h = c(TRUE, FALSE, TRUE, FALSE)), sum)
#---
aggregate(x ~ g, data.frame(x = c(TRUE, FALSE, TRUE), g = c("a", "a", "b")), sum)
#---
aggregate(x ~ g, data.frame(x = c(1.5, 2.5, 3.5), g = c(1, 1, 2)), mean)
#---
aggregate(x ~ g, data.frame(x = c(1.5, 2.5, 3.5), g = c(1, 1, 2)), function(z) sum(z) / 2)
#---
aggregate(x ~ g, data.frame(x = c(1.5, 2.5, 3.5), g = c(1, 1, 2)), max)
#---
aggregate(x ~ g, data.frame(x = c("p", "q", "r"), g = c(1, 1, 2)), paste, collapse = "+")
#---
tryCatch(aggregate(x ~ g, data.frame(x = 1:2, g = 1:2)), error = function(e) conditionMessage(e))
#---
tryCatch(aggregate(~ g, data.frame(x = 1:2, g = 1:2), sum), error = function(e) conditionMessage(e))
