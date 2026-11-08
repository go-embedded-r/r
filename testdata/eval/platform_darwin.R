format(as.Date("0005-03-10"), "%Y|%y|%C")
#---
as.POSIXct("2024-03-31 02:30:00", tz = "Europe/Paris")
#---
as.POSIXct("2024-10-27 02:30:00", tz = "Europe/Paris")
#---
format(as.POSIXct("2024-07-08 07:32:13", tz = "Asia/Tokyo"), "%s %Z")
#---
format(as.POSIXct("2024-07-15 12:00:00", tz = "Europe/Paris"), "%s")
#---
format(as.POSIXct("2024-07-15 12:00:00", tz = "Europe/Paris"), "%s %Z|%s")
#---
format(as.POSIXct("2024-01-15 12:00:00", tz = "America/New_York"), "%s %z")
#---
format(as.POSIXct("2024-07-15 12:00:00", tz = "Australia/Sydney"), "%Z %s")
#---
format(as.POSIXct("2024-07-15 12:00:00", tz = "Asia/Tokyo"), "%s")
