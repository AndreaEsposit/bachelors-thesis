// cargo rustc --target wasm32-wasi -- -Z wasi-exec-model=reactor
#![no_main]
use std::alloc::{alloc, dealloc, Layout};
use std::fs;
//use std::io::prelude::*;
use std::slice;
use std::str;

// staic variable to keep track of the message
static mut MESSAGE_LEN: i32 = 0;

#[no_mangle]
pub unsafe fn new_alloc(length: usize) -> *mut u8 {
    let align = std::mem::align_of::<usize>();
    let layout = Layout::from_size_align_unchecked(length, align);
    let b = alloc(layout);
    println!("This is the pointer with alloc {:?}", b);
    b
}

#[no_mangle]
pub unsafe fn new_dealloc(ptr: *mut u8, length: usize) {
    let align = std::mem::align_of::<usize>();
    let layout = Layout::from_size_align_unchecked(length, align);
    dealloc(ptr, layout);
}

#[no_mangle]
pub extern "C" fn n_dealloc2(ptr: *mut u8, length: usize) {
    unsafe {
        let data = Vec::from_raw_parts(ptr, length, length);
        std::mem::drop(data);
    }
}

#[no_mangle]
pub extern "C" fn get_message_len() -> i32 {
    unsafe { MESSAGE_LEN }
}

#[no_mangle]
pub extern "C" fn store_data(ptr: *mut u8, length: usize) {
    let slice = unsafe { slice::from_raw_parts(ptr, length) };
    let strin = str::from_utf8(slice).unwrap();
    let _file = fs::File::create("alloc.txt").unwrap();
    println!("{}", strin);
    // Write the text to the file we created
    fs::write("/alloc.txt", strin).expect("Unable to write to file");
}

#[no_mangle]
pub extern "C" fn retrive_data() -> *mut u8 {
    let contents =
        fs::read_to_string("testData.txt").expect("Something went wrong reading the file");
    let mut data = contents.into_bytes();
    unsafe {
        MESSAGE_LEN = data.len() as i32;
    }

    println! {"This is the message size: {}, and this is the capacity: {}", data.len(), data.capacity()};
    let new_ptr = data.as_mut_ptr();
    // take ownership of the memory block where the new message is written and esure its
    // destryctuion is not called when the object goes out of scope at the end of the func
    std::mem::forget(data);
    println!("{:?}", new_ptr);

    new_ptr
}

#[no_mangle]
pub extern "C" fn retrive_data2() -> *mut u8 {
    let contents =
        fs::read_to_string("testData2.txt").expect("Something went wrong reading the file");
    let mut data = contents.into_bytes();
    unsafe {
        MESSAGE_LEN = data.len() as i32;
    }
    let new_ptr = data.as_mut_ptr();
    // take ownership of the memory block where the new message is written and esure its
    // destryctuion is not called when the object goes out of scope at the end of the func
    println! {"This is the message size: {}, and this is the capacity: {}", data.len(), data.capacity()};
    std::mem::forget(data);
    println!("{:?}", new_ptr);
    new_ptr
}
