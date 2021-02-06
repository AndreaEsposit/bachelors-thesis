#![no_main]
use std::fs;
use std::io::prelude::*;

#[no_mangle]
pub extern "C" fn foo() {
    println!("Before creating");
    let mut file = fs::File::create("helloworld.txt").unwrap();

    // Write the text to the file we created
    write!(file, "Hello world!\n").unwrap();
}
