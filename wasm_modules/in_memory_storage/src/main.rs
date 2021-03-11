// cargo rustc --target wasm32-wasi -- -Z wasi-exec-model=reactor
#![no_main]
mod storage;
use lazy_static::lazy_static;
use protobuf;
use std::alloc::{alloc, dealloc, Layout};
use std::slice;
use std::sync::RwLock;
use storage::{ReadRequest, ReadResponse, WriteRequest, WriteResponse};

// staic variable to keep track of the message size
static mut RESPONSE_LEN: i32 = 0;

// represents the value of the message we are storing
lazy_static! {
    static ref VALUE: RwLock<String> = RwLock::new("".to_string());
}

static mut TIME: ContentTime = ContentTime {
    nseconds: 0,
    seconds: 0,
};

#[derive(Debug)]
struct ContentTime {
    nseconds: i32,
    seconds: i64,
}

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
pub extern "C" fn get_response_len() -> i32 {
    unsafe { RESPONSE_LEN }
}

#[no_mangle]
pub extern "C" fn store_data(ptr: *mut u8, length: usize) -> *mut u8 {
    let slice = unsafe { slice::from_raw_parts(ptr, length) };
    let pb_message: WriteRequest = protobuf::Message::parse_from_bytes(slice)
        .expect("Something went wrong unmarshalling the pb-message");

    let time = pb_message.get_Timestamp();
    let val = pb_message.get_Value();

    // Assign the new value to store
    unsafe {
        let mut myval = VALUE.write().unwrap();
        *myval = val.to_string();
        TIME.nseconds = time.get_nanos();
        TIME.seconds = time.get_seconds();
    }

    // return response
    let mut response: WriteResponse = protobuf::Message::new();
    response.set_Ok(1);
    let mut new_bytes = protobuf::Message::write_to_bytes(&response)
        .expect("Something went wrong marshalling the pb-message");

    unsafe {
        RESPONSE_LEN = new_bytes.capacity() as i32;
    }
    let new_ptr = new_bytes.as_mut_ptr();
    // take ownership of the memory block where the new message is written and esure its
    // destryctuion is not called when the object goes out of scope at the end of the func
    std::mem::forget(new_bytes);
    new_ptr
}

#[no_mangle]
pub extern "C" fn read_data(ptr: *mut u8, length: usize) -> *mut u8 {
    let slice = unsafe { slice::from_raw_parts(ptr, length) };
    let pb_message: ReadRequest = protobuf::Message::parse_from_bytes(slice)
        .expect("Something went wrong unmarshalling the pb-message");

    if pb_message.get_FileName() == "test" {
        let mut time = protobuf::well_known_types::Timestamp::new();
        unsafe {
            time.set_nanos(TIME.nseconds);
            time.set_seconds(TIME.seconds);
        }

        // return response
        let mut response: ReadResponse = protobuf::Message::new();

        let val = VALUE.read().unwrap();
        response.set_Ok(1);
        response.set_Timestamp(time);
        response.set_Value(val.to_string());

        let mut new_bytes = protobuf::Message::write_to_bytes(&response)
            .expect("Something went wrong marshalling the pb-message");

        unsafe {
            RESPONSE_LEN = new_bytes.capacity() as i32;
        }
        let new_ptr = new_bytes.as_mut_ptr();
        // take ownership of the memory block where the new message is written and esure its
        // destryctuion is not called when the object goes out of scope at the end of the func
        std::mem::forget(new_bytes);
        return new_ptr;
    } else {
        // FileName is not "test"
        // return response
        let mut response: ReadResponse = protobuf::Message::new();
        let time = protobuf::well_known_types::Timestamp::new();
        response.set_Value("".to_string());
        response.set_Ok(0);
        response.set_Timestamp(time);

        let mut new_bytes = protobuf::Message::write_to_bytes(&response)
            .expect("Something went wrong marshalling the pb-message");

        unsafe {
            RESPONSE_LEN = new_bytes.capacity() as i32;
        }
        let new_ptr = new_bytes.as_mut_ptr();
        // take ownership of the memory block where the new message is written and esure its
        // destryctuion is not called when the object goes out of scope at the end of the func
        std::mem::forget(new_bytes);
        return new_ptr;
    }
}
