# Base library, part 3: conditions. Written from the help pages of R 4.6.2
# and checked against the oracle (docs/m3-plan.md). BSD licence.

errorCondition <- function(message, ..., class = character(), call = NULL)
  structure(class = c(class, "error", "condition"),
            list(message = as.character(message), call = call, ...))

warningCondition <- function(message, ..., class = character(), call = NULL)
  structure(class = c(class, "warning", "condition"),
            list(message = as.character(message), call = call, ...))
