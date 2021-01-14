// cargo build --target wasm32-unknown-unknown
#[no_mangle]
// Simple function that return age given current year and date of birth.
extern "C" fn age(cy: i32, yob: i32) -> i32 {
    cy - yob
}
