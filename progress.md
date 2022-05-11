## Progress 4.28 (develop a language)

## Finished
- the lexer part
- the parser part
- the evaluater part

```go
// compute 5!
n = 5
f = 1
while n {
    f = f * n
    n = n - 1
}
f
```
``` sh
--  1 INTEGER 5 --
--  2 INTEGER 1 --
--  3 INTEGER 0 --
--  4 INTEGER 120 --
=============  evaluater final result  ============
--  INTEGER 120 --
```

``` go
// compute 0 + 1 + ... + 10
while i <= 10 {
    sum = sum + i
    i = i + 1
}
```
``` go
--  1 INTEGER 11 --
--  2 INTEGER 55 --
=============  evaluater final result  ============
--  INTEGER 55 --
```
```go
i = 1
while i <= 10 {
    i % 2 ? odd = odd + i : even = even + i
    i = i + 1
}
odd
even
odd + even
```

```go
--  1 INTEGER 1 --
--  2 INTEGER 11 --
--  3 INTEGER 25 --
--  4 INTEGER 30 --
--  5 INTEGER 55 --
```

## Next 
- add functions