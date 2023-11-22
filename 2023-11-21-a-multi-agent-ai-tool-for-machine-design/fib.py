Sure, here is a simple Python function that generates fibonacci series. It's implemented as a generator.

```python
def fibonacci():
    a, b = 0, 1
    while True:
        yield a
        b, a = a+b, b

# Test the function
fib_series = fibonacci()
for _ in range(10):
    print(next(fib_series))
```
In the above code, the `fibonacci` function will yield the next number in the Fibonacci series each time it is called. It will start with 0 and 1, and then continue the series. The `for` loop at the end is simply to test the function and print the first 10 numbers of the Fibonacci series.
