Here is how you would translate that Python code into Go. We're going to create an equivalent function called `fibonacci` that uses channels to achieve similar generator-like behavior. Note that Go does not have generator functions like Python does.

```go
package main

import "fmt"

func fibonacci(n int, c chan int) {
    x, y := 0, 1
    for i := 0; i < n; i++ {
        c <- x
        x, y = y, x+y
    }
    close(c)
}

func main() {
    c := make(chan int, 10)
    go fibonacci(cap(c), c)
    for i := range c {
        fmt.Println(i)
    }
}
```
In the Go code above, we create a Go routine with the function `fibonacci` and use a channel `c` to communicate between the main function and the `fibonacci` function. The fibonacci function generates numbers in the fibonacci series and sends them on the channel `c`. In the main function, we range over the channel `c` to receive and print each generated fibonacci number. 

Remember that channels in Go are a way for goroutines to communicate with each other and synchronize their execution. In this case, we're using the channel to demonstrate Python's generator-like behavior, not necessarily for its synchronization capabilities. 

Feel free to replace `cap(c)` in `go fibonacci(cap(c), c)` with the exact number of fibonacci numbers you'd like to generate.
