# Condition constructors, written from the help page and the oracle.

..newCondition <- function(message, ..., class, call, kind) {
  obj <- list()
  obj["message"] <- list(as.character(message))
  obj["call"] <- list(call)
  obj <- c(obj, list(...))
  attr(obj, "class") <- c(class, kind, "condition")
  obj
}

errorCondition <- function(message, ..., class = NULL, call = NULL) {
  ..newCondition(message, ..., class = class, call = call, kind = "error")
}

warningCondition <- function(message, ..., class = NULL, call = NULL) {
  ..newCondition(message, ..., class = class, call = call, kind = "warning")
}
