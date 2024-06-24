## Go grammar tip

### `*` and `&`

In Go (Golang), the `*` and `&` symbols are used for working with **pointer**. **Pointer** allow user to reference and manipulate memory addresses directly.

- `&` address-of Operator (get reference)
  The `&` symbol is used to get the memory address of a variable.
  When you place `&` before a variable, it returns the pointer to that variable.

- `*` Dereference Operator (get value)
  The `*` symbol is used to dereference a pointer, meaning it accesses the value stored at the memory address the pointer refers to.

When you place `*` before a pointer, it returns the value stored at the pointer's address.

### `label`

a `label` is used in conjunction with the `goto`, `break`, or `continue` statements to alter the flow of control within a program. `label` provides a way to jump to a different part of the code. They are typically used for error handling, breaking out of nested loops, or other advanced control flow mechanisms.

### `defer`

`defer` is a powerful keyword. It is used to ensure a function call is performed later in a program's execution, usually for purpose of cleanup. The deferred function call is executed after the surrounding function completes, but before it actually returns.

Deferred functions are executed in LIFO(last in, first out) order, meaning the last deferred function will be executed first.

Use cases: **Resource Cleanup** Closing files, releasing locks. **Post-Processing** Logging or handling errors.
