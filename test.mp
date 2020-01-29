# A simple any function
let any = fn(arr, f) {
  let iter = fn(arr, accumulated) {
    if (len(arr) == 0) {
      accumulated
    } else {
      iter(rest(arr), accumulated || f(first(arr)))
    }
  };

  iter(arr, false);
};

let a = [1, 3, 5, 7, 9];
let b = push(a, 10);

puts(any(a, fn(x) { x % 2 == 0 })) # false
puts(any(b, fn(x) { x % 2 == 0 })) # true