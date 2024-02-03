## Compiling and running code
rustc main.rs
./main
## Rust macros
println! -> Rust macro
println -> normal function
## Using cargo
```
cargo new hello_cargo
cargo build
cargo run
```
Check code, but do not compile
```
cargo check
```
Build release version
```
cargo build --release
```

Cargo.lock file should be added to source control!

update a crate to get a new version
```
cargo update
```

## Variables
```
let apples = 5; // immutable
let mut bananas = 5; // mutable
```
## Reading user input
```
io::stdin() // returns instance of std::io::Stdin
    .read_line(&mut guess) // & means reference
    .expect("Failed to read line"); // handle errors
```
## Printing values
```
let x = 5;
let y = 10;

println!("x = {} and y = {}", x, y);
// or
println!("x = {x} and y = {y}");
```

## Depenedency versioning
The number 0.8.3 is actually shorthand for ^0.8.3, which means any version that is at least 0.8.3 but below 0.9.0.

0.8.3 -> means e.g. 0.8.4 can be used, but not 0.9.0 or higher

## Documentation
```
cargo doc --open
```
