# Base library, part 3: conditions. Written from the help page of conditions
# in R 4.6.2 (errorCondition, warningCondition) and checked against the oracle
# (docs/m3-plan.md). BSD licence.

# a condition is a list: message, call, then the extra fields; its class
# vector starts with the caller's classes
..newCondition <- function(message, call, extra, classes) {
  cond <- c(list(message = as.character(message), call = call), extra)
  class(cond) <- classes
  cond
}

errorCondition <- function(message, ..., class = character(), call = NULL)
  ..newCondition(message, call, list(...), c(class, "error", "condition"))

warningCondition <- function(message, ..., class = character(), call = NULL)
  ..newCondition(message, call, list(...), c(class, "warning", "condition"))
