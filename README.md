# Monkey+

[The Monkey programming language](https://monkeylang.org/) implemented by Ramen.

## Overview

The Monkey programming language is from the book series [*Writing An Interpreter In Go*](https://interpreterbook.com/) and [*Writing A Compiler In Go*](https://compilerbook.com/), the author is @mrnugget. 

I must say this series is the best books for the beginner of whom want to start their journey of implementing a programming language. If you are also interested in, go and grab your own copy now!

The plus sign means that I have done some additional works to support more basic feature for Monkey. For detailed information, just scroll down and see the *Feature* section.

**IMPORTANT**: This language is neither ready to use nor with reliable performance. It's an educational purpose language. Use it at your own risk.

## Features

### Integer arithmetic

Monkey supports basic integer arithmetic operations:

```
>> 5 + 5
10
>> 7 - 6
1
>> 5 * 9
45
>> 7 / 3
2
>> 1 / 0
ERROR: the right operand of / is 0
```

Monkey+ also supports `%` operator.

```
>> 149 % 22
17
>> 1 % 0
ERROR: the right operand of % is 0
```

### Boolean operation

Monkey support integer, string and boolean compare operations:

```
>> 5 > 1
true
>> 4 == 5
false
>> "a" == "b"
false
>> true != false
true
```

Monkey+ also supports basic logic arithmetic and string compare operations:

```
>> "aa" > "b"
false
>> "apple" < "appletart"
true
>> true && false
false
>> true && false || true
true
```

You may noticed, the strings are compared by their lexicographical order.

### Variable binding

Monkey support `let` keyword to bind a value to a variable

```
>> let a = 5;
>> puts(a);
5
```

### Comments

Monkey+ supports single line comments starts with `#`.

```
# This is a single comment
puts(5); # This is an end of line comment
```

### Control Flow

Monkey supports `if-else` control flow.

```
let test = fn(x) {
  if (x > 10) {
    return "Large than 10";
  } else {
    return "Less than 10";
  }
};

puts(test(11), test(5));
```

### Function and Closure

Monkey supports functions:

```
>> let add = fn(a, b) { a + b; };
>> add(5, 6);
11
```

And high order function and closures!

```
>> let twice = fn(f, x) { return f(f(x)); };
>> twice(fn(x) { x * x; }, 4);
256
>> let adder = fn(x) { return fn(y) { return x + y; } };
>> let add_five = adder(5);
>> add_five(6);
11
```

### Built-in functions

- `len(x)`: return the length of `x`. `x` should be a string, an array or a hash.
- `puts(a, b, ...)`: prints each variable in lines.
- `eval(c)`: eval a code snippet `c`, the environment will not be exported to current env.
- `load(f)`: load a file `f` into the global environment.
- `type(x)`: report `x`'s type.

### Built-in Data Structures

#### String

The string is totally re-designed to support escaped characters.

In Monkey+ you can also using `[]` to access the character in the string.

```
>> let a = "hello";
>> let b = "world!";
>> a[1]; # It's 0-indexed
e
>> b[10]; # If not in the correct range...
null
>> let c = a + "\n" + b; # Using + to do concatenation!
>> c; # Escaped charater is support too!
hello
world!
>> len(c); # The length is correct too!
12
```

#### Array

Monkey support array literal, `[]` random access. The item in it can be different, like Python's `list`.

```
>> let a = [15, 20, 25];
>> a[0];
15
>> a[2];
25
>> a[-1]; # If out of bound...
null
>> len(a);
3
```

These built-in functions will help you:

- `first(a)`: returns the first item in array `a`, `null` if the array is empty.
- `last(a)`: returns the last item in array `a`, `null` if the array is empty.
- `rest(a)`: returns the array that except the first element of `a`, `null` if the array has no more than one item.
- `push(a, el)`: append `el` to the end of `a` and produce a new array. `a` will stay unmodified after calling that.

```
>> let a = [15, 20, 25];
>> first(a);
15
>> last(a);
25
>> rest(a);
[20, 25]
>> push(a, 30);
[15, 20, 25, 30]
>> a;
[15, 20, 25]
```

#### Hash

Hash is hash map or dictionary in other languages. Like Array, Monkey also supports Hash literal and `[]` random access.

The key of Hash can be integer, boolean and string and the value can be any valid type.

```
>> let h = {"a": "b", 1: fn(x) { x + x; }, false: [123]};
>> h["a"];
b
>> h[1](2);
4
>> h[false][0];
123
>> h["b"]; # If the key is not found...
null
>> len(h);
3
```

Monkey+ extended the built-in functions for Hash so you can do more operations:

- `set(h, k, v)`: set `k` to `v` in `h`. If `k` exists the value will be replaced to `v`, otherwise create. It will return a new Hash instead modify the original one.
- `contains(h, k)`: test if `k` is in `h`.
- `delete(h, k)`: delete the entry which key is `k`. It will also return a new Hash instead modify the original one.

```
>> let h = {"a": "b", 1: 2, false: [123]};
>> set(h, "new", "year");
{1: 2, false: [123], a: b, new: year}
>> set(h, "a", "c");
{a: c, 1: 2, false: [123]}
>> contains(h, "a");
true
>> delete(h, "a");
{1: 2, false: [123]}
>> delete(h, 456);
{a: c, 1: 2, false: [123]}
>> h; # Remain unmodified
{a: c, 1: 2, false: [123]}
```

## Example

```
# A simple any function
let any = fn(arr, f) {
  let iter = fn(arr, accumulated) {
    if (len(arr) == 0) {
      accumulated;
    } else {
      iter(rest(arr), accumulated || f(first(arr)));
    }
  };

  iter(arr, false);
};

let a = [1, 3, 5, 7, 9];
let b = push(a, 10);

puts(any(a, fn(x) { x % 2 == 0; })); # false
puts(any(b, fn(x) { x % 2 == 0; })); # true
```

## Usage

**Currently, only interpreter is supported.** The compiler support will be added when I finish the next book :).

Be sure you have installed go. My version is `go version go1.13.5 darwin/amd64`, but I'm not using any fancy feature of go, so it should works for go 1.7 and later.

Then just clone this repo:

```bash
$ git clone https://github.com/lxdlam/monkey-plus
$ cd monkey-plus
```

There are different running mode:

```bash
$ go run main.go # You're entering the REPL
$ go run main.go foo.mp # Running file foo.mp
$ go run main.go -f foo.mp # Same as above
$ go run main.go -c "let a = 5; puts(a)" # Running a code snippet
```

The tests are also be extended for the new feature, so you can try:

```bash
$ go test ./...
```

Nothing should fail.

## License

The original works are licensed under MIT License. Thanks Thorsten Ball!

My works are also licensed under MIT License.