// cargo build --target wasm32-unknown-unknown
// Simple function that return age given current year and date of birth.
#[no_mangle]
pub extern "C" fn age(cy: i32, yob: i32) -> i32 {
    cy - yob
}

#[no_mangle]
pub extern "C" fn print_hello() {
    println!("Hello, world!");
}
