// cargo rustc --target wasm32-wasi -- -Z wasi-exec-model=reactor
#![no_main]
mod storage;
use serde_derive::{Deserialize, Serialize};
use serde_json::json;
use std::alloc::{alloc, dealloc, Layout};
use std::{fs::File, io, path::Path, slice};
use storage::{ReadRequest, ReadResponse, WriteRequest, WriteResponse};

// staic variable to keep track of the message size
static mut RESPONSE_LEN: i32 = 0;

#[derive(Serialize, Deserialize, Debug)]
struct Content {
    nseconds: i32,
    seconds: i64,
    value: String,
}

fn write_to_file(file_path: String, data: serde_json::Value) -> Result<(), io::Error> {
    let file = File::create(file_path)?;
    let e = serde_json::to_writer_pretty(file, &data)?;
    Ok(e)
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
    let mut file_path = "./".to_owned();
    file_path.push_str(pb_message.get_FileName());
    file_path.push_str(".json");

    let time = pb_message.get_Timestamp();

    let data = json!({
    "seconds" : time.get_seconds(),
    "nseconds": time.get_nanos(),
    "value": pb_message.get_Value(),});

    // write to file
    let write_result = write_to_file(file_path, data);
    let write_result = match write_result {
        Ok(_result) => 1,
        Err(_e) => 0,
    };

    // return response
    let mut response: WriteResponse = protobuf::Message::new();
    response.set_Ok(write_result);
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

    let mut file_path = "./".to_owned();
    file_path.push_str(pb_message.get_FileName());
    file_path.push_str(".json");
    let pathf = Path::new(&file_path);
    let file = File::open(pathf);

    match file {
        Ok(file) => {
            let reader = io::BufReader::new(file);

            let file_content: Content =
                serde_json::from_reader(reader).expect("JSON was not well-formatted");

            let mut time = protobuf::well_known_types::Timestamp::new();
            time.set_nanos(file_content.nseconds);
            time.set_seconds(file_content.seconds);

            // return response
            let mut response: ReadResponse = protobuf::Message::new();
            response.set_Value(file_content.value);
            response.set_Ok(1);
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
        Err(_e) => {
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
    };
}
