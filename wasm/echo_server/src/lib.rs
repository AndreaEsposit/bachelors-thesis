use std::alloc::{alloc, dealloc, Layout};

#[no_mangle]
pub unsafe fn my_alloc(size: usize) -> *mut u8 {
    let align = std::mem::align_of::<usize>();
    let layout = Layout::from_size_align_unchecked(size, align);
    alloc(layout)
}

#[no_mangle]
pub unsafe fn my_dealloc(ptr: *mut u8, size: usize) {
    let align = std::mem::align_of::<usize>();
    let layout = Layout::from_size_align_unchecked(size, align);
    dealloc(ptr, layout);
}

#[no_mangle]
pub extern "C" fn echo(ptr: *mut u8, len: usize) -> *mut u8 {
    let s = std::str::from_utf8(unsafe { std::slice::from_raw_parts(ptr, len) }).unwrap();
    println!("Wasm has recived: {}, sending it back!", s);
    ptr
}
