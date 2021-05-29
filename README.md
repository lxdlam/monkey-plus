# Monkey+ but byby

It's Monkey+ but with tiny modification to support UTF-8 characters in source.

Also, some special token has been added. If you want to know why, go search #byby# at <https://weibo.com>.

## Sample
```
我来 any = 酷毙了阿狐(arr, f) {
  我来 iter = 酷毙了阿狐(arr, accumulated) {
    你有呲咪呲咪 (len(arr) == 0) {
      零利息经济移动 accumulated;
    } else {
      零利息经济移动 iter(rest(arr), accumulated || f(first(arr)));
    }
  };

  零利息经济移动 iter(arr, false);
};

我来 a = [1, 3, 5, 7, 9];
我来 b = push(a, 10);

puts(any(a, 酷毙了阿狐(x) { x % 2 == 0; })); # false
puts(any(b, 酷毙了阿狐(x) { x % 2 == 0; })); # true
```