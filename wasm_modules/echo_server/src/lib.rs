// cargo wasi build
mod echo;
use echo::echo_message;
use std::alloc::{alloc, dealloc, Layout};

// staic variable to keep track of the
static mut MESSAGE_LEN: i32 = 0;

#[no_mangle]
pub unsafe fn new_alloc(length: usize) -> *mut u8 {
    let align = std::mem::align_of::<usize>();
    let layout = Layout::from_size_align_unchecked(length, align);
    alloc(layout)
}

#[no_mangle]
pub unsafe fn new_dealloc(ptr: *mut u8, length: usize) {
    let align = std::mem::align_of::<usize>();
    let layout = Layout::from_size_align_unchecked(length, align);
    dealloc(ptr, layout);
}

#[no_mangle]
pub extern "C" fn get() -> i32 {
    unsafe { MESSAGE_LEN }
}

#[no_mangle]
pub extern "C" fn echo(ptr: *mut u8, length: usize) -> *mut u8 {
    let slice = unsafe { std::slice::from_raw_parts(ptr, length) };
    // unmarshal the rawfiles to make a message
    let recived_message: echo_message = protobuf::Message::parse_from_bytes(slice).unwrap();
    println!(
        "Wasm has recived this message: {:?}, sending it back!",
        recived_message.get_content()
    );
    let mut new_message: echo_message = protobuf::Message::new();
    new_message.set_content(recived_message.get_content().into());
    // use the command down to write a new message back to the client
    //new_message.set_content("back".to_string().into());

    // marshal the message to bytes
    let mut new_bytes = protobuf::Message::write_to_bytes(&new_message).unwrap();
    unsafe {
        MESSAGE_LEN = new_bytes.len() as i32;
    }
    let new_ptr = new_bytes.as_mut_ptr();
    // take ownership of the memory block where the result string is written and esure its
    // destryctuion is not called when the object goes out of scope at the end of the func
    std::mem::forget(new_bytes);
    new_ptr
}
