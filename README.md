# Description

This project is my first try at building an interpreter. This `rinha` tree-walking interpreter was built in the span of a few hours for a competition/hackathon (see [github repo](https://github.com/aripiprazole/rinha-de-compiler) - in portuguese). This interpreter reads a JSON Abstract Syntax Tree (AST) that gets generated from the lexer/parser that the creators of the competition and language created.

# AST Format

```json
{
  "name": "ata.rinha",
  "expression": {
    "kind": "Binary",
    "lhs": {
      "kind": "Int",
      "value": 1,
      "location": ..
    },
    "op": "Add",
    "rhs": {
      "kind": "Int",
      "value": 2,
      "location": ..
    },
    "location": ..
  },
  "location": ..
}
```

# Steps to run interpreter

## Locally:

1. Create a new .rinha file (see examples directory for some examples)
2. Parse the content of that file into the json AST using `rinha examples/int.rinha > examples/int.json` (install rust/rinha before this)
3. Run this interpreter using `go run . ./examples/<filename>.json`

## Using docker:

1. Build image: `docker build -t rinha .`
2. Run container: `docker run -it rinha examples/<file>.json`

# Example

### `fib.rinha` file:

```go
let fib = fn (n) => {
  if (n < 2) {
    n
  } else {
    fib(n - 1) + fib(n - 2)
  }
};

print("fib: " + fib(11))
```

### Output:

`fib: 89`

# Conclusion

This was a very fun project that helped me learn more about how a lexer, parser, and interpreter work. This topic interests me since I was introduced to it in college, and this was the perfect opportunity to learn, and, at the same time, compete in a hackathon style competition.
Obviously, this interpreter could be optimized in many ways, but due to lack of time (I worked on it for a few hours after I was done with my full-time job) I wasn't able to make many changes after making all operations work.
