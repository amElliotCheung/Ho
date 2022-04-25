## Progress 4.21 (develop a language)

## Finished
- the lexer part
- the parser part

```go
// test code
a = 0
b = 0
c = 0
if 1 > 2 {
    a = a + 2*3 + 1
    a = a - 3*2/1 + 1 - 1/1
} else if 2 > 3 {
    b = b + 1 + 0
} else {
    c = c + 1
}
even = 0
odd = 0
i = 0
while i < 10 {
    if i % 2 == 0 { // even number
        even = even + i
    } else { 
        odd = odd + i
    }
    i = i + 1
}
even + odd
```
``` python
 I have 9 children ==> ((a=0) (b=0) (c=0) (if (1>2)then I have 2 children ==> ((a=((a+(2*3))+1)) (a=(((a-((3*2)/1))+1)-(1/1))) )else if (2>3)then I have 1 children ==> ((b=((b+1)+0)) )else if truethen I have 1 children ==> ((c=(c+1)) )) (even=0) (odd=0) (i=0) (while (i<10) do I have 2 children ==> ((if ((i%2)==0)then I have 1 children ==> ((even=(even+i)) )else if truethen I have 1 children ==> ((odd=(odd+i)) )) (i=(i+1)) )) (even+odd) )
```

## Current 
- evaluation