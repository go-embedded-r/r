class(Sys.time()); class(Sys.Date()); typeof(Sys.Date()); typeof(Sys.time())
#---
Sys.Date() > as.Date("2020-01-01"); Sys.time() > as.POSIXct("2020-01-01", tz = "UTC")
#---
d <- as.Date("2024-03-10"); d; unclass(d); class(d); typeof(unclass(d))
#---
as.Date("2024-03-10")
#---
as.Date("2024/03/10")
#---
as.Date(c("2024-01-01", NA, "2024-12-31"))
#---
as.Date(NA)
#---
as.Date(character(0))
#---
as.Date("2024-02-30")
#---
as.Date("2024-02-30", optional = TRUE)
#---
as.Date("garbage")
#---
as.Date(c(NA, "2024-02-03"))
#---
as.Date("01/02/2024", "%d/%m/%Y")
#---
as.Date("01/02/2024", format = "%m/%d/%Y")
#---
as.Date("2024-01-05 23:59:59")
#---
as.Date("3jan2024", "%d%b%Y"); as.Date("3 January 2024", "%d %B %Y")
#---
as.Date(c("1jan1960", "2jan1960", "31mar1960", "30jul1960"), "%d%b%Y")
#---
as.Date(c("02/27/92", "01/14/92"), "%m/%d/%y")
#---
as.Date("2024-366", "%Y-%j")
#---
as.Date(19000)
#---
as.Date(19000.7)
#---
as.Date(-1)
#---
as.Date(32768, origin = "1900-01-01")
#---
as.Date(35981, origin = "1899-12-30")
#---
as.Date(5, origin = as.Date("2024-01-01"))
#---
as.Date(factor("2024-03-10"))
#---
as.Date(TRUE)
#---
as.Date(NULL)
#---
as.Date(as.Date("2024-03-10"))
#---
as.Date(as.POSIXct("2024-01-05 23:59:59", tz = "America/New_York"))
#---
as.Date(as.POSIXct("2024-01-05 23:59:59", tz = "America/New_York"), tz = "America/New_York")
#---
as.Date(as.POSIXct("2024-01-05 01:00:00", tz = "UTC"), tz = "Asia/Tokyo")
#---
as.Date(ISOdate(2010, 4, 13, c(0, 12)))
#---
as.Date(as.POSIXlt("2024-03-10 23:00:00", tz = "UTC"))
#---
as.Date.character("2024-03-10"); as.Date.numeric(10); as.Date.default(NULL); as.Date.factor(factor("2024-01-02"))
#---
as.Date.POSIXct(as.POSIXct("2024-01-05 23:59:59", tz = "UTC")); as.Date.POSIXlt(as.POSIXlt("2024-01-05", tz = "UTC"))
#---
.Date(19000); .Date(c(1, 2))
#---
format(as.Date("2024-03-10"))
#---
format(as.Date("2024-03-10"), "%A %d %B %Y %j %U %u")
#---
format(as.Date("2024-03-10"), "%a %b %e %y %m/%d %H:%M:%S %p %Z")
#---
format(as.Date(c("2024-03-10", NA)))
#---
format(as.Date(c(a = "2024-03-10", b = "2025-01-02")))
#---
format(as.Date("2024-01-01") + 0:2, "%W %V %G %g %w")
#---
format(as.Date("2024-12-30"), "%G-%V-%u"); format(as.Date("2021-01-03"), "%G-%V-%u %U %W")
#---
format(as.Date("0005-03-10"), "%Y|%y|%C")
#---
format(as.Date(19000.7))
#---
format(as.Date("2024-03-10"), "%%|%n|%t|%h|%D|%F|%x|%c")
#---
format.Date(as.Date("2024-03-10"), "%d.%m."); as.character.Date(as.Date("2024-03-10"))
#---
as.character(as.Date("2024-03-10")); as.character(as.Date(c("2024-03-10", NA)))
#---
paste("d:", as.Date("2024-03-10")); paste0(as.Date("2024-03-10"), "x")
#---
print(as.Date("2024-03-10")); print.Date(as.Date(c("2024-03-10", "2024-03-11")))
#---
as.Date(character(0))[0]
#---
x <- as.Date(c("2024-03-10", "2024-03-11", "2024-03-12")); x[2]; x[-1]; x[c(TRUE, FALSE, TRUE)]; x[[3]]
#---
x <- as.Date(c("2024-03-10", "2024-03-11")); x[2] <- as.Date("2025-01-01"); x
#---
x <- as.Date(c("2024-03-10", "2024-03-11")); x[2] <- "2025-01-01"; x
#---
x <- as.Date("2024-03-10"); x[3] <- as.Date("2025-01-01"); x
#---
`[.Date`(as.Date(c("2024-03-10", "2024-03-11")), 1); `[[.Date`(as.Date(c("2024-03-10", "2024-03-11")), 2)
#---
x <- as.Date(c("2024-03-10", "2024-03-11")); `[<-.Date`(x, 1, as.Date("2000-01-01"))
#---
c(as.Date("2024-01-05"), as.Date("2024-01-06"))
#---
c(as.Date("2024-01-05"), "2024-01-07")
#---
c.Date(as.Date("2024-01-05"), as.Date("2024-01-06"))
#---
rev(as.Date(c("2024-01-05", "2024-01-06")))
#---
rep(as.Date("2024-01-05"), 3); rep.Date(as.Date("2024-01-05"), times = 2)
#---
range(as.Date(c("2024-01-06", "2024-01-05"))); min(as.Date(c("2024-01-06", "2024-01-05"))); max(as.Date(c("2024-01-06", "2024-01-05")))
#---
sum(as.Date("2024-01-05"))
#---
Summary.Date(as.Date(c("2024-01-06", "2024-01-05")))
#---
mean(as.Date(c("2024-01-05", "2024-01-06"))); mean.Date(as.Date(c("2024-01-05", "2024-01-08")))
#---
diff(as.Date(c("2024-01-05", "2024-01-08", "2024-02-01")))
#---
unique(as.Date(c("2024-01-05", "2024-01-05", "2024-01-06")))
#---
sort(as.Date(c("2024-01-06", "2024-01-05", "2024-01-07")))
#---
order(as.Date(c("2024-01-06", "2024-01-05", "2024-01-07")))
#---
as.Date("2024-01-05") + 1; 1 + as.Date("2024-01-05"); as.Date("2024-01-05") + 1:3
#---
as.Date("2024-01-05") - 1; as.Date("2024-01-05") - 1:2
#---
as.Date("2024-01-05") - as.Date("2023-12-25")
#---
as.Date("2024-01-05") + as.Date("2023-12-25")
#---
1 - as.Date("2024-01-05")
#---
-as.Date("2024-01-05")
#---
+as.Date("2024-01-05")
#---
as.Date("2024-01-05") * 2
#---
as.Date("2024-01-05") - "x"
#---
as.Date("2024-01-05") > as.Date("2023-12-25"); as.Date("2024-01-05") == as.Date("2024-01-05"); as.Date("2024-01-05") != as.Date("2024-01-06")
#---
as.Date("2024-01-05") > "2023-12-25"; "2024-01-05" <= as.Date("2024-01-05")
#---
as.Date(c("2024-01-05", "2024-01-06")) < as.Date("2024-01-06")
#---
as.Date("2024-01-05") + as.difftime(2, units = "days"); as.Date("2024-01-05") - as.difftime(1, units = "days")
#---
as.Date("2024-01-05") + as.difftime(48, units = "hours")
#---
`+.Date`(as.Date("2024-01-05"), 2); `-.Date`(as.Date("2024-01-05"), 2); is.function(Ops.Date)
#---
is.numeric(as.Date("2024-01-05")); is.numeric.Date(as.Date("2024-01-05")); is.numeric(Sys.time()); is.numeric.POSIXt(Sys.time())
#---
as.numeric(as.Date("2024-03-10")); as.integer(as.Date("2024-03-10")); as.double(as.Date("2024-03-10"))
#---
weekdays(as.Date("2024-03-10")); weekdays(as.Date("2024-03-10"), TRUE); weekdays.Date(as.Date("2024-03-11"))
#---
months(as.Date("2024-03-10")); months(as.Date("2024-03-10"), abbreviate = TRUE); months.Date(as.Date("2024-12-11"))
#---
quarters(as.Date("2024-03-10")); quarters(as.Date(c("2024-01-01", "2024-04-01", "2024-07-31", "2024-10-31"))); quarters.Date(as.Date("2024-11-11"))
#---
weekdays(as.POSIXct("2024-03-10 12:00:00", tz = "UTC")); months(as.POSIXct("2024-03-10 12:00:00", tz = "UTC")); quarters(as.POSIXct("2024-03-10 12:00:00", tz = "UTC"))
#---
weekdays.POSIXt(as.POSIXct("2024-03-10 12:00:00", tz = "UTC"), TRUE); months.POSIXt(as.POSIXct("2024-03-10 12:00:00", tz = "UTC"), TRUE); quarters.POSIXt(as.POSIXct("2024-03-10 12:00:00", tz = "UTC"))
#---
julian(as.Date("2024-03-10"))
#---
julian.Date(as.Date("2024-03-10"), origin = as.Date("2024-01-01"))
#---
as.integer(format(as.Date("2024-03-10"), "%Y")) + 1L
#---
seq(as.Date("2024-01-01"), as.Date("2024-01-10"), by = "3 days")
#---
seq(as.Date("2024-01-01"), by = "week", length.out = 3)
#---
seq(as.Date("2024-01-31"), by = "month", length.out = 4)
#---
seq(as.Date("2024-01-01"), as.Date("2025-01-01"), by = "quarter")
#---
seq(as.Date("2024-02-29"), by = "year", length.out = 3)
#---
seq(as.Date("2024-01-01"), as.Date("2024-01-05"), by = 2)
#---
seq(as.Date("2024-01-10"), as.Date("2024-01-01"), by = -3)
#---
seq(as.Date("2024-01-01"), as.Date("2024-01-05"), length.out = 3)
#---
seq(as.Date("2024-01-01"), by = "day", length.out = 3)
#---
seq(as.Date("2024-01-01"), as.Date("2024-03-01"), by = "month")
#---
seq(as.Date("2024-01-01"), by = as.difftime(2, units = "days"), length.out = 3)
#---
seq(as.Date("2024-01-01"), as.Date("2024-01-04"))
#---
seq(as.Date("2024-01-01"), as.Date("2024-01-02"), by = "fortnight")
#---
seq(as.Date("2024-01-01"))
#---
seq(as.Date("2024-01-01"), as.Date("2024-01-02"), by = "day", length.out = 2)
#---
seq.Date(as.Date("2024-01-01"), by = "2 weeks", length.out = 2)
#---
x <- as.POSIXct("2024-03-10 12:34:56", tz = "UTC"); x; unclass(x); attributes(x)
#---
as.POSIXct("2024-03-10 12:34:56", tz = "UTC")
#---
as.POSIXct("2024-03-10 12:34", tz = "UTC")
#---
as.POSIXct("2024-03-10", tz = "UTC")
#---
as.POSIXct("2024/03/10 01:02:03", tz = "UTC")
#---
as.POSIXct("2024-03-10T01:02:03", tz = "UTC")
#---
as.POSIXct("garbage", tz = "UTC")
#---
as.POSIXct(c("2024-03-10 01:02:03", "garbage"), tz = "UTC")
#---
as.POSIXct(c("2024-03-10 01:02:03", NA), tz = "UTC")
#---
as.POSIXct(NA)
#---
as.POSIXct(NA_character_, tz = "UTC")
#---
as.POSIXct("2024-03-10 12:34:56.789", tz = "UTC"); unclass(as.POSIXct("2024-03-10 12:34:56.5", tz = "UTC"))
#---
as.POSIXct("10/03/2024 08:09", format = "%d/%m/%Y %H:%M", tz = "UTC")
#---
as.POSIXct(1710064800, tz = "UTC")
#---
as.POSIXct(1710064800, origin = "1970-01-01", tz = "UTC")
#---
as.POSIXct(3600, origin = "2024-01-01", tz = "UTC")
#---
as.POSIXct(c(0, 86400), tz = "UTC")
#---
as.POSIXct(as.Date("2024-03-10"))
#---
attributes(as.POSIXct(as.Date("2024-03-10")))
#---
as.POSIXct(as.POSIXlt("2024-03-10 01:02:03", tz = "UTC"))
#---
as.POSIXct(as.POSIXct("2024-03-10 01:02:03", tz = "UTC"), tz = "Asia/Tokyo")
#---
as.POSIXct(NULL)
#---
as.POSIXct(list())
#---
as.POSIXct.numeric(3600, tz = "UTC"); as.POSIXct.default("2024-01-01", tz = "UTC"); as.POSIXct.Date(as.Date("2024-01-01")); as.POSIXct.POSIXlt(as.POSIXlt("2024-01-01", tz = "UTC"))
#---
.POSIXct(0, "UTC"); .POSIXct(c(0, 60), tz = "Asia/Tokyo")
#---
as.POSIXct("2024-07-01 12:00:00", tz = "America/New_York")
#---
unclass(as.POSIXct("2024-07-01 12:00:00", tz = "America/New_York"))
#---
as.POSIXct("2024-01-01 12:00:00", tz = "America/New_York")
#---
as.POSIXct("2024-03-31 02:30:00", tz = "Europe/Paris")
#---
as.POSIXct("2024-10-27 02:30:00", tz = "Europe/Paris")
#---
as.POSIXct("2024-06-01 12:00:00", tz = "Asia/Kolkata"); as.POSIXct("2024-06-01 12:00:00", tz = "Australia/Sydney")
#---
as.POSIXct("2024-06-01 12:00:00", tz = "GMT")
#---
as.POSIXct("2024-06-01 12:00:00", tz = "Europe/London")
#---
format(as.POSIXct("2024-03-10 12:34:56", tz = "UTC"))
#---
format(as.POSIXct("2024-03-10 12:34:00", tz = "UTC")); format(as.POSIXct("2024-03-10 00:00:00", tz = "UTC"))
#---
format(as.POSIXct(c("2024-03-10 00:00:00", "2024-03-10 00:00:01"), tz = "UTC"))
#---
format(as.POSIXct(c("2024-03-10 00:00:00", NA), tz = "UTC"))
#---
format(as.POSIXct("2024-03-10 12:34:56", tz = "UTC"), "%H:%M")
#---
format(as.POSIXct("2024-03-10 12:34:56", tz = "UTC"), "%c|%x|%X|%r|%p|%D|%e|%k|%l|%s|%n|%t|%h|%C|%y|%G|%g|%V|%u|%w|%W|%U|%j|%R|%T|%F")
#---
format(as.POSIXct("2024-03-10 00:05:06", tz = "UTC"), "%I %l %p %k %H")
#---
format(as.POSIXct("2024-03-10 12:05:06", tz = "UTC"), "%I %l %p")
#---
format(as.POSIXct("2024-07-01 12:00:00", tz = "America/New_York"), "%Z %z")
#---
format(as.POSIXct("2024-01-01 12:00:00", tz = "America/New_York"), "%Z %z")
#---
format(as.POSIXct("2024-01-01 12:00:00", tz = "Asia/Kolkata"), "%Z %z")
#---
format(as.POSIXct("2024-03-10 12:00:00", tz = "UTC"), tz = "Asia/Tokyo", usetz = TRUE)
#---
format(as.POSIXct("2024-03-10 12:00:00", tz = "UTC"), usetz = TRUE)
#---
format(as.POSIXct("2024-03-10 12:00:00", tz = "UTC"), tz = "America/Los_Angeles")
#---
format(as.POSIXct("2024-03-10 12:34:56.789", tz = "UTC"), "%H:%M:%OS3")
#---
format(as.POSIXct("2024-03-10 12:34:56.789", tz = "UTC"), "%OS")
#---
format(as.POSIXct("2024-03-10 12:34:56.789", tz = "UTC"), "%OS1 %OS6 %OS0")
#---
format(as.POSIXct("2024-03-10 12:34:56.789", tz = "UTC"), digits = 3)
#---
format(as.POSIXct("2024-03-10 12:34:56.789", tz = "UTC"), digits = 1)
#---
format(as.POSIXct("2024-03-10 12:34:56.5", tz = "UTC"), digits = 6)
#---
format(as.POSIXct("2024-03-10 12:34:56.5", tz = "UTC"), format = "%Y-%m-%d %H:%M:%OS6")
#---
format(as.POSIXct("2024-03-10 12:34:56.789", tz = "UTC"))
#---
op <- options(digits.secs = 3); r <- format(as.POSIXct("2024-03-10 12:34:56.789", tz = "UTC")); options(op); r
#---
op <- options(digits.secs = 3); r <- as.POSIXct("2024-03-10 12:34:56.789", tz = "UTC"); options(op); r
#---
format(as.POSIXct(c(a = "2024-03-10 12:34:56"), tz = "UTC"))
#---
format.POSIXct(as.POSIXct("2024-03-10 12:34:56", tz = "UTC"), "%S")
#---
format(as.POSIXlt("2024-03-10 12:34:56", tz = "UTC"))
#---
format.POSIXlt(as.POSIXlt("2024-03-10 12:34:56", tz = "UTC"), "%d", usetz = TRUE)
#---
format(as.POSIXlt("2024-03-10 12:34:56", tz = "America/New_York"), usetz = TRUE)
#---
strftime(as.POSIXct("2024-03-10 12:00:00", tz = "UTC"), "%H")
#---
strftime(as.Date("2024-03-10"), "%d/%m/%Y")
#---
strftime(as.POSIXct("2024-03-10 12:00:00", tz = "UTC"), "%H:%M", tz = "Asia/Tokyo")
#---
strftime("2024-03-10 12:00:00", "%H", tz = "UTC")
#---
as.character(as.POSIXct("2024-03-10 01:02:03", tz = "UTC")); as.character.POSIXt(as.POSIXct("2024-03-10", tz = "UTC"))
#---
as.character(as.POSIXct("2024-03-10 01:02:03.25", tz = "UTC"))
#---
paste(as.POSIXct("2024-03-10 01:02:03", tz = "UTC"))
#---
as.numeric(as.POSIXct("2024-03-10 01:02:03", tz = "UTC"))
#---
print(as.POSIXct("2024-03-10 01:02:03", tz = "UTC")); print.POSIXct(as.POSIXct("2024-03-10", tz = "UTC"))
#---
print(as.POSIXct("2024-03-10 01:02:03", tz = "UTC"), tz = "Asia/Tokyo")
#---
print(as.POSIXct("2024-03-10 01:02:03", tz = "UTC"), usetz = FALSE)
#---
as.POSIXct(character(0), tz = "UTC")
#---
as.POSIXlt(character(0), tz = "UTC")
#---
x <- as.POSIXct(c("2024-03-10 01:02:03", "2024-03-11 04:05:06"), tz = "UTC"); x[2]; x[-1]; x[[1]]; length(x)
#---
x <- as.POSIXct(c("2024-03-10 01:02:03", "2024-03-11 04:05:06"), tz = "UTC"); attributes(x[2])
#---
`[.POSIXct`(as.POSIXct(c("2024-03-10 01:02:03", "2024-03-11 04:05:06"), tz = "UTC"), 2); `[[.POSIXct`(as.POSIXct(c("2024-03-10 01:02:03", "2024-03-11 04:05:06"), tz = "UTC"), 1)
#---
c(as.POSIXct("2024-03-10 01:02:03", tz = "UTC"), as.POSIXct("2024-03-11", tz = "UTC"))
#---
attributes(c(as.POSIXct("2024-03-10 01:02:03", tz = "UTC"), as.POSIXct("2024-03-11", tz = "Asia/Tokyo")))
#---
c.POSIXct(as.POSIXct("2024-03-10 01:02:03", tz = "UTC"), as.POSIXct("2024-03-11", tz = "UTC"))
#---
rep(as.POSIXct("2024-03-10 01:02:03", tz = "UTC"), 2); rep.POSIXct(as.POSIXct("2024-03-10", tz = "UTC"), times = 2)
#---
range(as.POSIXct(c("2024-03-10 01:02:03", "2024-03-11 04:05:06"), tz = "UTC"))
#---
max(as.POSIXct(c("2024-03-10 01:02:03", "2024-03-11 04:05:06"), tz = "UTC")); min(as.POSIXct(c("2024-03-10 01:02:03", "2024-03-11 04:05:06"), tz = "UTC"))
#---
is.function(Summary.POSIXct)
#---
Summary.POSIXct(as.POSIXct(c("2024-03-10 01:02:03", "2024-03-11 04:05:06"), tz = "UTC"))
#---
mean(as.POSIXct(c("2024-03-10 00:00:00", "2024-03-11 00:00:00"), tz = "UTC")); mean.POSIXct(as.POSIXct(c("2024-03-10 00:00:00", "2024-03-10 02:00:00"), tz = "UTC"))
#---
as.POSIXct("2024-01-05 10:00", tz = "UTC") + 3600
#---
3600 + as.POSIXct("2024-01-05 10:00", tz = "UTC")
#---
as.POSIXct("2024-01-05 10:00", tz = "UTC") - 60
#---
as.POSIXct("2024-01-05 10:00", tz = "UTC") + as.difftime(2, units = "hours")
#---
as.POSIXct("2024-01-05 10:00", tz = "UTC") - as.difftime(1, units = "days")
#---
as.POSIXct("2024-01-05 10:00", tz = "UTC") + as.POSIXct("2024-01-05 10:00", tz = "UTC")
#---
-as.POSIXct("2024-01-05 10:00", tz = "UTC")
#---
60 - as.POSIXct("2024-01-05 10:00", tz = "UTC")
#---
as.POSIXct("2024-01-05 10:00", tz = "UTC") * 2
#---
as.POSIXct("2024-01-05 10:00", tz = "UTC") - "a"
#---
as.POSIXct("2024-01-05 10:00", tz = "UTC") - as.POSIXct("2024-01-05 09:59:30", tz = "UTC")
#---
as.POSIXct("2024-01-05 10:00", tz = "UTC") - as.POSIXct("2024-01-05 08:00", tz = "UTC")
#---
as.POSIXct("2024-01-05 10:00", tz = "UTC") - as.POSIXct("2024-01-04 08:00", tz = "UTC")
#---
as.POSIXct("2024-01-05 10:00", tz = "UTC") - as.POSIXct("2024-01-05 09:30", tz = "UTC")
#---
as.POSIXct("2024-01-05 10:00", tz = "UTC") > as.POSIXct("2024-01-05 09:30", tz = "UTC"); as.POSIXct("2024-01-05 10:00", tz = "UTC") == as.POSIXct("2024-01-05 10:00", tz = "UTC")
#---
as.POSIXct("2024-01-05 10:00", tz = "UTC") < "2024-01-05 11:00"
#---
as.POSIXct("2024-01-05 10:00", tz = "UTC") > as.POSIXlt("2024-01-05 09:00", tz = "UTC")
#---
`+.POSIXt`(as.POSIXct("2024-01-05 10:00", tz = "UTC"), 60); `-.POSIXt`(as.POSIXct("2024-01-05 10:00", tz = "UTC"), 60); is.function(Ops.POSIXt)
#---
as.POSIXct("2024-03-30 12:00", tz = "Europe/Paris") + 86400
#---
as.POSIXct("2024-03-30 12:00", tz = "Europe/Paris") + as.difftime(1, units = "days")
#---
as.POSIXlt("2024-03-10 12:34:56", tz = "UTC")
#---
unclass(as.POSIXlt("2024-03-10 12:34:56", tz = "UTC"))
#---
attributes(as.POSIXlt("2024-03-10 12:34:56", tz = "UTC"))
#---
x <- as.POSIXlt("2024-03-10 12:34:56", tz = "UTC"); c(x$sec, x$min, x$hour, x$mday, x$mon, x$year, x$wday, x$yday, x$isdst)
#---
x <- as.POSIXlt("2024-03-10 12:34:56", tz = "UTC"); x$zone; x$gmtoff
#---
unclass(as.POSIXlt("2024-03-10 12:34:56.789", tz = "Asia/Tokyo"))
#---
unclass(as.POSIXlt(as.POSIXct("2024-03-10 12:34:56", tz = "Asia/Tokyo")))
#---
unclass(as.POSIXlt(as.POSIXct("2024-07-01 12:34:56", tz = "America/New_York")))
#---
attributes(as.POSIXlt(as.POSIXct("2024-07-01 12:34:56", tz = "America/New_York")))
#---
attributes(as.POSIXlt("2024-03-10", tz = "Asia/Tokyo"))$tzone
#---
attributes(as.POSIXlt("2024-03-10", tz = "GMT"))$tzone; attributes(as.POSIXlt("2024-03-10", tz = "UTC"))$tzone
#---
attributes(as.POSIXlt(as.POSIXct("2024-03-10", tz = "Europe/London")))$tzone
#---
attributes(as.POSIXlt(as.POSIXct("2024-03-10", tz = "UTC")))$tzone
#---
as.POSIXlt(as.Date("2024-03-10"))
#---
unclass(as.POSIXlt(as.Date("2024-03-10")))
#---
as.POSIXlt(as.Date("2024-03-10"))$yday
#---
as.POSIXlt(1710064800, tz = "UTC")
#---
as.POSIXlt(1710064800, origin = "1970-01-01", tz = "UTC")$hour
#---
as.POSIXlt(as.POSIXct("2024-03-10 12:34:56", tz = "UTC"))
#---
as.POSIXlt(as.POSIXct("2024-03-10 12:34:56", tz = "UTC"), tz = "Asia/Tokyo")
#---
as.POSIXlt(as.POSIXlt("2024-03-10 12:34:56", tz = "UTC"))
#---
as.POSIXlt(as.POSIXlt("2024-03-10 12:34:56", tz = "UTC"), tz = "Asia/Tokyo")
#---
as.POSIXlt(factor("2024-03-10 01:02:03"), tz = "UTC")
#---
as.POSIXlt(NULL)
#---
as.POSIXlt(NA)
#---
as.POSIXlt(list())
#---
as.POSIXlt("garbage", tz = "UTC")
#---
as.POSIXlt("garbage", tz = "UTC", optional = TRUE)
#---
as.POSIXlt("2024-03-10 01:02:03", tz = "UTC", format = "%Y-%m-%d %H:%M:%S")
#---
as.POSIXlt.character("2024-03-10", tz = "UTC"); as.POSIXlt.numeric(0, tz = "UTC"); as.POSIXlt.Date(as.Date("2024-01-01")); as.POSIXlt.POSIXct(as.POSIXct("2024-01-01", tz = "UTC")); as.POSIXlt.factor(factor("2024-01-01"), tz = "UTC"); as.POSIXlt.default(NULL)
#---
x <- as.POSIXlt(c("2024-03-10 01:02:03", "2024-03-11 04:05:06"), tz = "UTC"); length(x); x$hour; format(x); is.na(x)
#---
length.POSIXlt(as.POSIXlt(c("2024-03-10", "2024-03-11"), tz = "UTC")); is.na.POSIXlt(as.POSIXlt(c("2024-03-10", "2024-03-11"), tz = "UTC"))
#---
x <- as.POSIXlt("2024-03-10 01:02:03", tz = "UTC"); x$mday <- x$mday + 30; as.POSIXct(x)
#---
x <- as.POSIXlt("2024-01-31", tz = "UTC"); x$mon <- x$mon + 1; as.Date(x)
#---
x <- as.POSIXlt("2024-03-10 01:02:03", tz = "UTC"); x$hour <- 25; format(x); as.POSIXct(x)
#---
x <- as.POSIXlt("2024-03-10 01:02:03", tz = "UTC"); x$min <- -1; as.POSIXct(x)
#---
print(as.POSIXlt("2024-03-10 01:02:03", tz = "UTC")); print.POSIXlt(as.POSIXlt("2024-03-10", tz = "UTC"))
#---
strptime("2024-03-10 01:02:03", "%Y-%m-%d %H:%M:%S", tz = "UTC")
#---
unclass(strptime("2024-03-10 01:02:03", "%Y-%m-%d %H:%M:%S", tz = "UTC"))
#---
attributes(strptime("2024-03-10 01:02:03", "%Y-%m-%d %H:%M:%S", tz = "UTC"))
#---
attributes(strptime("2024-03-10", "%Y-%m-%d", tz = "Asia/Tokyo"))$tzone
#---
strptime("2024-13-01", "%Y-%m-%d", tz = "UTC")
#---
is.na(strptime("2024-13-01", "%Y-%m-%d", tz = "UTC"))
#---
strptime("2024-02-30 08:00", "%Y-%m-%d %H:%M", tz = "UTC")
#---
unclass(strptime("2024-02-30 08:00", "%Y-%m-%d %H:%M", tz = "UTC"))
#---
strptime("10:20", "%H:%M", tz = "UTC")$hour
#---
strptime("10:20", "%H:%M", tz = "UTC")$min; strptime("10:20", "%H:%M", tz = "UTC")$sec
#---
strptime("Mar 5 2024", "%b %d %Y", tz = "UTC")
#---
strptime("mar 5 2024", "%b %d %Y", tz = "UTC"); strptime("MARCH 5 2024", "%B %d %Y", tz = "UTC")
#---
strptime("Tue, 23 Mar 2010 14:36:38 -0400", "%a, %d %b %Y %H:%M:%S %z", tz = "UTC")
#---
strptime("2024-03-10 12:34:56 +0530", "%Y-%m-%d %H:%M:%S %z", tz = "UTC")
#---
strptime("2024-03-10 12:34:56 -05:00", "%Y-%m-%d %H:%M:%S %z", tz = "UTC")
#---
strptime("2024-03-10 12:34:56 +9999", "%Y-%m-%d %H:%M:%S %z", tz = "UTC")
#---
strptime("3:04:05 PM", "%I:%M:%S %p", tz = "UTC")$hour
#---
strptime("12:04:05 AM", "%I:%M:%S %p", tz = "UTC")$hour; strptime("12:04:05 pm", "%I:%M:%S %p", tz = "UTC")$hour
#---
strptime("13:04:05 PM", "%I:%M:%S %p", tz = "UTC")
#---
strptime("2024-03-10 24:00:00", "%Y-%m-%d %H:%M:%S", tz = "UTC")
#---
strptime("2024-03-10 25:00:00", "%Y-%m-%d %H:%M:%S", tz = "UTC")
#---
strptime("2024-03-10 12:34:61", "%Y-%m-%d %H:%M:%S", tz = "UTC")
#---
strptime("20/2/06 11:16:16.683", "%d/%m/%y %H:%M:%OS", tz = "UTC")
#---
unclass(strptime("20/2/06 11:16:16.683", "%d/%m/%y %H:%M:%OS", tz = "UTC"))$sec
#---
format(strptime("20/2/06 11:16:16.683", "%d/%m/%y %H:%M:%OS", tz = "UTC"), digits = 3)
#---
strptime("69", "%y", tz = "UTC")$year; strptime("68", "%y", tz = "UTC")$year; strptime("00", "%y", tz = "UTC")$year; strptime("99", "%y", tz = "UTC")$year
#---
strptime("19 24", "%C %y", tz = "UTC")$year
#---
strptime("2024-366", "%Y-%j", tz = "UTC"); strptime("2023-366", "%Y-%j", tz = "UTC")
#---
strptime("2024-060", "%Y-%j", tz = "UTC")$mon; strptime("2024-060", "%Y-%j", tz = "UTC")$mday
#---
strptime("2024-03-10 trailing", "%Y-%m-%d", tz = "UTC")
#---
strptime("2024-03-10", "%Y-%m-%d %H", tz = "UTC")
#---
strptime("  2024-3-1", "%Y-%m-%d", tz = "UTC")
#---
strptime("2024 03 10", "%Y%m%d", tz = "UTC")
#---
strptime("20240310", "%Y%m%d", tz = "UTC")
#---
strptime("2024-03-10", "%F", tz = "UTC"); strptime("12:34:56", "%T", tz = "UTC")$hour; strptime("12:34", "%R", tz = "UTC")$min
#---
strptime("03/10/24", "%D", tz = "UTC"); strptime("24/03/10", "%x", tz = "UTC")
#---
strptime("Sun Mar 10 12:34:56 2024", "%c", tz = "UTC")
#---
strptime("12:34:56", "%X", tz = "UTC")$sec
#---
strptime("Sunday 10 March 2024", "%A %d %B %Y", tz = "UTC")
#---
strptime("2024-03-10", "%Y-%m-%d", tz = "America/New_York")
#---
unclass(strptime("2024-07-10 01:02:03", "%Y-%m-%d %H:%M:%S", tz = "America/New_York"))
#---
strptime(c("2024-03-10", "x", NA), "%Y-%m-%d", tz = "UTC")
#---
length(strptime(c("2024-03-10", "x"), "%Y-%m-%d", tz = "UTC")); class(strptime("2024-03-10", "%Y-%m-%d", tz = "UTC"))
#---
strptime(c("2024-03-10", "2024-03-11"), c("%Y-%m-%d", "%Y-%m-%d"), tz = "UTC")
#---
strptime(c("2024-03-10", "11/03/2024"), c("%Y-%m-%d", "%d/%m/%Y"), tz = "UTC")
#---
strptime("2024-03-10", c("%Y-%m-%d", "%d/%m/%Y"), tz = "UTC")
#---
strptime(character(0), "%Y", tz = "UTC")
#---
strptime(20240310, "%Y%m%d", tz = "UTC")
#---
strptime("2024-03-10", "%Y-%m-%d", tz = "UTC")$wday; strptime("2024-12-31", "%Y-%m-%d", tz = "UTC")$yday
#---
ISOdatetime(2024, 3, 10, 1, 2, 3, tz = "UTC")
#---
ISOdatetime(2024, 3, 10, 1, 2, 3.5, tz = "UTC"); unclass(ISOdatetime(2024, 3, 10, 1, 2, 3.5, tz = "UTC"))
#---
ISOdatetime(2024, 2, 30, 1, 2, 3, tz = "UTC")
#---
ISOdatetime(2024, 1:3, 1, 0, 0, 0, tz = "UTC")
#---
ISOdate(2024, 3, 10)
#---
ISOdate(2010, 4, 13, c(0, 12))
#---
ISOdate(2024, 13, 1)
#---
ISOdate(2024, 3, 10) - ISOdate(2024, 2, 26)
#---
difftime(as.Date("2024-03-10"), as.Date("2024-03-01"))
#---
difftime(as.Date("2024-03-10"), as.Date("2024-03-01"), units = "hours")
#---
difftime(as.POSIXct("2024-01-05", tz = "UTC"), as.POSIXct("2024-01-04 12:00", tz = "UTC"))
#---
difftime(as.POSIXct("2024-01-05", tz = "UTC"), as.POSIXct("2024-01-04 23:59:30", tz = "UTC"))
#---
difftime(as.POSIXct("2024-01-05", tz = "UTC"), as.POSIXct("2024-01-04 23:30", tz = "UTC"))
#---
difftime(as.POSIXct("2024-01-05", tz = "UTC"), as.POSIXct("2023-12-04", tz = "UTC"), units = "weeks")
#---
difftime(as.POSIXct("2024-01-05", tz = "UTC"), as.POSIXct("2024-01-04", tz = "UTC"), units = "mins")
#---
difftime(as.POSIXct("2024-01-05", tz = "UTC"), as.POSIXct("2024-01-04", tz = "UTC"), units = "secs")
#---
difftime(as.POSIXct("2024-01-05", tz = "UTC"), as.POSIXct("2024-01-04", tz = "UTC"), units = "w")
#---
difftime(as.POSIXct("2024-01-05", tz = "UTC"), as.POSIXct("2024-01-04", tz = "UTC"), units = "years")
#---
difftime("2024-01-05", "2024-01-04", tz = "UTC")
#---
difftime(as.POSIXct(c("2024-01-05", "2024-01-06"), tz = "UTC"), as.POSIXct("2024-01-04", tz = "UTC"))
#---
difftime(as.POSIXct(c("2024-01-05", NA), tz = "UTC"), as.POSIXct("2024-01-04", tz = "UTC"))
#---
difftime(as.POSIXct(NA), as.POSIXct("2024-01-04", tz = "UTC"))
#---
difftime(as.POSIXct("2024-01-05", tz = "Asia/Tokyo"), as.POSIXct("2024-01-05", tz = "UTC"))
#---
d <- difftime(as.Date("2024-03-10"), as.Date("2024-03-01")); d; unclass(d); class(d); attributes(d)
#---
d <- as.difftime(c(0, 30, 60), units = "mins"); d; units(d); as.numeric(d); as.numeric(d, units = "secs"); as.numeric(d, units = "hours")
#---
as.difftime(c("0:3:20", "11:23:15"))
#---
as.difftime(c("3:20", "23:15", "2:"), format = "%H:%M")
#---
as.difftime(30, units = "mins"); as.difftime(2, units = "hours") + as.difftime(30, units = "mins")
#---
as.difftime(30)
#---
as.difftime(30, units = "fortnights")
#---
as.difftime("x")
#---
as.difftime(TRUE)
#---
as.difftime(as.difftime(5, units = "secs"))
#---
d <- as.difftime(c(0, 30, 60), units = "mins"); format(d); format(d, digits = 3); format.difftime(d[1])
#---
d <- as.difftime(90, units = "mins"); units(d) <- "hours"; d
#---
d <- as.difftime(2, units = "days"); units(d) <- "secs"; d; `units<-`(d, "weeks")
#---
d <- as.difftime(2, units = "days"); units(d) <- "fortnights"
#---
d <- as.difftime(2, units = "days"); units(d) <- "days"; d
#---
units(as.difftime(2, units = "days")); units.difftime(as.difftime(2, units = "hours")); `units<-.difftime`(as.difftime(2, units = "hours"), "mins")
#---
units(1)
#---
as.difftime(c(30, 90), units = "mins")
#---
print(as.difftime(1.5, units = "hours")); print.difftime(as.difftime(c(1, 2), units = "weeks"))
#---
print(as.difftime(1.23456789, units = "secs"), digits = 3)
#---
print(as.difftime(numeric(0), units = "secs"))
#---
d <- as.difftime(c(1, 2, 3), units = "hours"); d[2]; d[-1]; d[[3]]; `[.difftime`(d, 1)
#---
d <- as.difftime(c(1, 2, 3), units = "hours"); attributes(d[2])
#---
c(as.difftime(1, units = "hours"), as.difftime(30, units = "mins"))
#---
c(as.difftime(1, units = "hours"), as.difftime(2, units = "hours"))
#---
c(as.difftime(1, units = "hours"), 5)
#---
c.difftime(as.difftime(1, units = "mins"), as.difftime(2, units = "mins")); c.difftime()
#---
rep(as.difftime(1, units = "hours"), 2)
#---
as.difftime(1, units = "hours") + as.difftime(30, units = "mins")
#---
as.difftime(1, units = "hours") + as.difftime(2, units = "hours")
#---
as.difftime(1, units = "hours") - as.difftime(30, units = "mins")
#---
as.difftime(1, units = "hours") + 1; 1 + as.difftime(1, units = "hours"); as.difftime(1, units = "hours") - 1
#---
-as.difftime(1, units = "hours"); +as.difftime(1, units = "hours")
#---
as.difftime(1, units = "hours") * 2; 2 * as.difftime(1, units = "hours"); as.difftime(4, units = "hours") / 2
#---
as.difftime(1, units = "hours") * as.difftime(1, units = "hours")
#---
2 / as.difftime(1, units = "hours")
#---
as.difftime(4, units = "hours") / as.difftime(2, units = "hours")
#---
as.difftime(1, units = "hours") ^ 2
#---
as.difftime(1, units = "hours") %% 2
#---
as.difftime(1, units = "hours") > as.difftime(30, units = "mins"); as.difftime(1, units = "hours") == as.difftime(60, units = "mins")
#---
as.difftime(1, units = "hours") > 0.5; as.difftime(1, units = "hours") < as.difftime(2, units = "hours")
#---
`*.difftime`(as.difftime(1, units = "hours"), 3); `/.difftime`(as.difftime(6, units = "hours"), 3); is.function(Ops.difftime)
#---
abs(as.difftime(-1.5, units = "hours")); round(as.difftime(1.567, units = "hours"), 1); floor(as.difftime(1.5, units = "hours")); ceiling(as.difftime(1.5, units = "hours")); trunc(as.difftime(1.5, units = "hours")); sign(as.difftime(-3, units = "hours"))
#---
signif(as.difftime(123.456, units = "secs"), 2)
#---
sqrt(as.difftime(4, units = "hours"))
#---
exp(as.difftime(1, units = "hours"))
#---
Math.difftime(as.difftime(-1, units = "mins"))
#---
sum(as.difftime(c(1, 2), units = "hours")); max(as.difftime(c(1, 2), units = "hours")); min(as.difftime(c(1, 2), units = "hours")); range(as.difftime(c(1, 2), units = "hours"))
#---
sum(as.difftime(1, units = "hours"), as.difftime(30, units = "mins"))
#---
prod(as.difftime(c(1, 2), units = "hours"))
#---
Summary.difftime(as.difftime(c(1, 2), units = "hours"), na.rm = FALSE)
#---
mean(as.difftime(c(1, 2), units = "hours")); mean.difftime(as.difftime(c(1, 4), units = "mins"))
#---
diff(as.POSIXct(c("2024-03-10 00:00", "2024-03-10 01:00", "2024-03-10 03:30"), tz = "UTC"))
#---
as.numeric(difftime(as.Date("2024-03-10"), as.Date("2024-03-01")), units = "mins")
#---
as.double(as.difftime(2, units = "hours"), units = "mins"); as.double.difftime(as.difftime(2, units = "hours"))
#---
as.numeric(as.difftime(2, units = "hours"))
#---
is.numeric(as.difftime(2, units = "hours")); is.numeric.difftime(as.difftime(2, units = "hours"))
#---
as.difftime(2, units = "hours") == 2
#---
unclass(as.difftime(c(a = 1, b = 2), units = "hours"))
#---
as.difftime(c(a = 1, b = 2), units = "hours")
#---
z <- as.POSIXct(c("2016-12-31 23:59:59", "2017-01-01 00:00:01"), tz = "UTC"); z[2] - z[1]
#---
.difftime(5, "mins"); .difftime(c(1, 2), units = "hours")
#---
x <- as.Date("2024-03-10"); x - as.Date("2024-03-10")
#---
as.Date("2024-03-10") - as.POSIXct("2024-03-09", tz = "UTC")
#---
as.POSIXct("2024-03-10", tz = "UTC") - as.Date("2024-03-09")
#---
as.Date("2024-03-10") > as.POSIXct("2024-03-09", tz = "UTC")
#---
inherits(as.Date("2024-03-10"), "Date"); oldClass(as.POSIXct("2024-03-10", tz = "UTC")); inherits(as.POSIXlt("2024-03-10", tz = "UTC"), "POSIXt")
#---
identical(as.Date("2024-03-10"), as.Date("2024-03-10")); identical(as.POSIXct("2024-03-10", tz = "UTC"), as.POSIXct("2024-03-10", tz = "UTC"))
#---
list(as.Date("2024-03-10"))
#---
list(a = as.POSIXct("2024-03-10", tz = "UTC"))
#---
structure(list(d = as.Date("2024-03-10")), class = "rec")
#---
sapply(c("2024-01-01", "2024-06-01"), function(s) format(as.Date(s), "%b"))
#---
lapply(as.Date(c("2024-01-01", "2024-06-01")), class)
#---
vapply(1:2, function(i) format(as.Date("2024-01-01") + i), "")
#---
as.Date(sapply(1:2, function(i) as.Date("2024-01-01") + i))
#---
for (d in as.Date("2024-01-01") + 0:1) print(d)
#---
x <- as.Date("2024-01-31"); y <- seq(x, by = "1 month", length.out = 3); format(y, "%d")
#---
getOption("digits"); getOption("nope"); getOption("nope", 5); getOption("OutDec"); getOption("digits.secs")
#---
op <- options(digits = 3); old <- options(op); getOption("digits")
#---
options(digits.secs = 2)$digits.secs; getOption("digits.secs"); options(digits.secs = NULL); getOption("digits.secs")
#---
withVisible(options(digits = 7))$visible; names(options("digits", "OutDec"))
#---
options(list(digits = 7))$digits
#---
options("digits")
#---
options(1)
#---
is.list(options()); "digits" %in% names(options())
#---
deparse1(quote(x + y)); deparse1(quote(f(a, b)), collapse = "|"); deparse1(1:3)
#---
gettext("abc"); gettext(c("a", "b")); gettextf("%d items", 3L); gettextf("%s-%s", "a", "b")
#---
rep.int(1:2, 2); rep.int("a", 3); rep.int(1:2, c(2, 1))
#---
is.factor(factor("a")); is.factor("a"); is.factor(1L)
#---
pmatch("me", c("mean", "median")); pmatch("mea", c("mean", "median")); pmatch("x", c("a", "b")); pmatch(c("a", "b"), c("abc", "b", "bcd"))
#---
pmatch("", ""); pmatch("m", c("mean", "median"), nomatch = 0); pmatch(c("a", "a"), c("abc", "abd")); pmatch(c("a", "a"), c("abc", "abd"), duplicates.ok = TRUE)
#---
pmatch("a", character(0)); pmatch(character(0), "a")
#---
unname(c(a = 1, b = 2)); unname(1:3); unname(list(x = 1)); unname(matrix(1:4, 2, dimnames = list(c("a", "b"), NULL)))
#---
lapply(as.Date(c("2024-01-01", "2024-06-01")), class); as.list(as.Date(c("2024-01-01", "2024-06-01"))); as.list.Date(as.Date("2024-01-01"))
#---
lapply(as.POSIXct(c("2024-01-01", "2024-06-01"), tz = "UTC"), format); as.list.POSIXct(as.POSIXct(c(a = "2024-01-01"), tz = "UTC"))
#---
lapply(as.difftime(c(1, 2), units = "mins"), format); as.list.difftime(as.difftime(1, units = "hours"))
#---
sapply(as.Date(c("2024-01-01", "2024-06-01")), function(d) format(d, "%m"))
#---
as.character(factor(c("b", "a", "b"))); paste(factor(c("x", "y"))); as.character(factor(c("b", NA)))
#---
Ops.Date(as.Date("2024-01-05"), as.Date("2024-01-06"))
#---
Ops.POSIXt(as.POSIXct("2024-01-05", tz = "UTC"), as.POSIXct("2024-01-06", tz = "UTC"))
#---
Ops.difftime(as.difftime(1, units = "hours"), as.difftime(2, units = "hours"))
#---
rep.difftime(as.difftime(1, units = "hours"), 2)
#---
rep.POSIXct(as.POSIXct("2024-01-05", tz = "UTC"), 2)
