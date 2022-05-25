# Progress 5.18 (develop a language)


# Current
- add functions to compiler and VM... difficult
  
## Interpreter
- integer, array, string and functions
- static type
- first-class function

```go
	// find the 20 fibonacci numbers
	// sure, in real life, we don't do it like this

	// getFibArray : a function which returns the fibonacci sequence of size "size".
	getFibArray := func(size) {
		// fib: a function, returns the n-th fibonacci number
		fib := func (n) {n <= 2 ? n : fib(n-1) + fib(n-2)} 
		n := 1
		fibArray := [] // empty array
		while n <= size {
			fibArray = append(fibArray, fib(n))
			n = n + 1
		}
		fibArray // return the last expression
	}
	getFibArray(20)
```
output
```python
[1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144, 233, 377, 610, 987, 1597, 2584, 4181, 6765, 10946]
```


## map in java and golang
- Java
  - map is an interface
  - hashmap is a class
- Golang
  - map is a keyword, to build a hard coded data structure
  - hashmap under the hood
  - it supports index opertor [] like array and slice in golang
  
  
```java
		HashMap<Integer, String> Sites = new HashMap<Integer, String>();
        Sites.put(1, "Google");
        Sites.get(1);
        Sites.remove(1);
		for (String i : people.keySet()) {
			//...
    	}
```
```go
	sites := make(map[int]string) // or sites := map[int]string{}
	sites[1] = "Google"
	sites[1]
	delete(sites, 1)
	for k, v := range sites {
		//...
	}
```

if we want a linkedHashMap in golang 

```go
// interface{} is something like a point to any type
type linkedHashMap struct {
	table    map[interface{}]interface{}
	ordering DoubleLinkedList // type "DoubleLinkedList" should be defined by user..
}
func (m *linkedHashMap) Put(key interface{}, value interface{}) {
	//...
}
func (m *linkedHashMap) Get(key interface{}, value interface{}) {
	//...
}
//...
```
[], range, make don't support user-defined structure, which means 
```go
	// impossible!!	
	sites := make(linkedHashMap[int]string) 
	...
```
```go
	// the only way
	sites := NewLinkedHashMap()
	sites.put(1, "Google")
	sites.get(1)

```
