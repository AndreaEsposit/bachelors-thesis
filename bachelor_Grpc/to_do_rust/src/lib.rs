use std::alloc::{alloc, dealloc, Layout};

/// Allocate memory into the module's linear memory
/// and return the offset to the start of the block.
// #[no_mangle]
// pub fn alloc(len: usize) -> *mut u8 {
//     // create a new mutable buffer with capacity `len`
//     let mut buf = Vec::with_capacity(len);
//     // take a mutable pointer to the buffer
//     let ptr = buf.as_mut_ptr();
//     // take ownership of the memory block and
//     // ensure that its destructor is not
//     // called when the object goes out of scope
//     // at the end of the function
//     std::mem::forget(buf);
//     // return the pointer so the runtime
//     // can write data at this offset
//     return ptr;
// }

// #[no_mangle]
// pub unsafe fn dealloc(ptr: *mut u8, size: usize) {
//     let data = Vec::from_raw_parts(ptr, size, size);

//     std::mem::drop(data);
// }

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

/// Given a pointer to the start of a byte array and
/// its length, return the sum of its elements.
#[no_mangle]
pub unsafe fn array_sum(ptr: *mut u8, len: usize) -> u8 {
    // create a Vec<u8> from the pointer to the
    // linear memory and the length
    let data = Vec::from_raw_parts(ptr, len, len);
    println!("{:?}", data);
    // actually compute the sum and return it
    data.iter().sum()
}

#[no_mangle]
pub extern "C" fn string() -> *const u8 {
    b"Hello, World!\0".as_ptr()
}

#[no_mangle]
pub extern "C" fn greet(s: *mut u8, len: usize) {
    let s = std::str::from_utf8(unsafe { std::slice::from_raw_parts(s, len) }).unwrap();
    println!("Hello, {}!", s)
}

#[no_mangle]
pub extern "C" fn deposit(amount: i32, total: i32) -> i32 {
    let _sum = amount + total;
    println!("This is the new total: {}", _sum);
    _sum
}

#[no_mangle]
pub extern "C" fn withdraw(amount: i32, total: i32) -> i32 {
    let new_total = total - amount;
    println!("This is the new total: {}", new_total);
    new_total
}

// #[no_mangle]
// pub extern "C" fn withdrawAcc(amount: i32, account_number: usize, accounts: *mut i32) -> i32 {
//     let slice = unsafe { std::slice::from_raw_parts_mut(accounts, account_number + 1) };
//     slice[account_number] -= amount;
//     slice[account_number]
// }
