// This file is generated by rust-protobuf 2.20.0. Do not edit
// @generated

// https://github.com/rust-lang/rust-clippy/issues/702
#![allow(unknown_lints)]
#![allow(clippy::all)]

#![allow(unused_attributes)]
#![rustfmt::skip]

#![allow(box_pointers)]
#![allow(dead_code)]
#![allow(missing_docs)]
#![allow(non_camel_case_types)]
#![allow(non_snake_case)]
#![allow(non_upper_case_globals)]
#![allow(trivial_casts)]
#![allow(unused_imports)]
#![allow(unused_results)]
//! Generated file from `storage.proto`

/// Generated files are compatible only with the same version
/// of protobuf runtime.
// const _PROTOBUF_VERSION_CHECK: () = ::protobuf::VERSION_2_20_0;

#[derive(PartialEq,Clone,Default)]
pub struct WriteRequest {
    // message fields
    pub FileName: ::std::string::String,
    pub Value: ::std::string::String,
    pub Timestamp: ::protobuf::SingularPtrField<::protobuf::well_known_types::Timestamp>,
    // special fields
    pub unknown_fields: ::protobuf::UnknownFields,
    pub cached_size: ::protobuf::CachedSize,
}

impl<'a> ::std::default::Default for &'a WriteRequest {
    fn default() -> &'a WriteRequest {
        <WriteRequest as ::protobuf::Message>::default_instance()
    }
}

impl WriteRequest {
    pub fn new() -> WriteRequest {
        ::std::default::Default::default()
    }

    // string FileName = 1;


    pub fn get_FileName(&self) -> &str {
        &self.FileName
    }
    pub fn clear_FileName(&mut self) {
        self.FileName.clear();
    }

    // Param is passed by value, moved
    pub fn set_FileName(&mut self, v: ::std::string::String) {
        self.FileName = v;
    }

    // Mutable pointer to the field.
    // If field is not initialized, it is initialized with default value first.
    pub fn mut_FileName(&mut self) -> &mut ::std::string::String {
        &mut self.FileName
    }

    // Take field
    pub fn take_FileName(&mut self) -> ::std::string::String {
        ::std::mem::replace(&mut self.FileName, ::std::string::String::new())
    }

    // string Value = 2;


    pub fn get_Value(&self) -> &str {
        &self.Value
    }
    pub fn clear_Value(&mut self) {
        self.Value.clear();
    }

    // Param is passed by value, moved
    pub fn set_Value(&mut self, v: ::std::string::String) {
        self.Value = v;
    }

    // Mutable pointer to the field.
    // If field is not initialized, it is initialized with default value first.
    pub fn mut_Value(&mut self) -> &mut ::std::string::String {
        &mut self.Value
    }

    // Take field
    pub fn take_Value(&mut self) -> ::std::string::String {
        ::std::mem::replace(&mut self.Value, ::std::string::String::new())
    }

    // .google.protobuf.Timestamp Timestamp = 3;


    pub fn get_Timestamp(&self) -> &::protobuf::well_known_types::Timestamp {
        self.Timestamp.as_ref().unwrap_or_else(|| <::protobuf::well_known_types::Timestamp as ::protobuf::Message>::default_instance())
    }
    pub fn clear_Timestamp(&mut self) {
        self.Timestamp.clear();
    }

    pub fn has_Timestamp(&self) -> bool {
        self.Timestamp.is_some()
    }

    // Param is passed by value, moved
    pub fn set_Timestamp(&mut self, v: ::protobuf::well_known_types::Timestamp) {
        self.Timestamp = ::protobuf::SingularPtrField::some(v);
    }

    // Mutable pointer to the field.
    // If field is not initialized, it is initialized with default value first.
    pub fn mut_Timestamp(&mut self) -> &mut ::protobuf::well_known_types::Timestamp {
        if self.Timestamp.is_none() {
            self.Timestamp.set_default();
        }
        self.Timestamp.as_mut().unwrap()
    }

    // Take field
    pub fn take_Timestamp(&mut self) -> ::protobuf::well_known_types::Timestamp {
        self.Timestamp.take().unwrap_or_else(|| ::protobuf::well_known_types::Timestamp::new())
    }
}

impl ::protobuf::Message for WriteRequest {
    fn is_initialized(&self) -> bool {
        for v in &self.Timestamp {
            if !v.is_initialized() {
                return false;
            }
        };
        true
    }

    fn merge_from(&mut self, is: &mut ::protobuf::CodedInputStream<'_>) -> ::protobuf::ProtobufResult<()> {
        while !is.eof()? {
            let (field_number, wire_type) = is.read_tag_unpack()?;
            match field_number {
                1 => {
                    ::protobuf::rt::read_singular_proto3_string_into(wire_type, is, &mut self.FileName)?;
                },
                2 => {
                    ::protobuf::rt::read_singular_proto3_string_into(wire_type, is, &mut self.Value)?;
                },
                3 => {
                    ::protobuf::rt::read_singular_message_into(wire_type, is, &mut self.Timestamp)?;
                },
                _ => {
                    ::protobuf::rt::read_unknown_or_skip_group(field_number, wire_type, is, self.mut_unknown_fields())?;
                },
            };
        }
        ::std::result::Result::Ok(())
    }

    // Compute sizes of nested messages
    #[allow(unused_variables)]
    fn compute_size(&self) -> u32 {
        let mut my_size = 0;
        if !self.FileName.is_empty() {
            my_size += ::protobuf::rt::string_size(1, &self.FileName);
        }
        if !self.Value.is_empty() {
            my_size += ::protobuf::rt::string_size(2, &self.Value);
        }
        if let Some(ref v) = self.Timestamp.as_ref() {
            let len = v.compute_size();
            my_size += 1 + ::protobuf::rt::compute_raw_varint32_size(len) + len;
        }
        my_size += ::protobuf::rt::unknown_fields_size(self.get_unknown_fields());
        self.cached_size.set(my_size);
        my_size
    }

    fn write_to_with_cached_sizes(&self, os: &mut ::protobuf::CodedOutputStream<'_>) -> ::protobuf::ProtobufResult<()> {
        if !self.FileName.is_empty() {
            os.write_string(1, &self.FileName)?;
        }
        if !self.Value.is_empty() {
            os.write_string(2, &self.Value)?;
        }
        if let Some(ref v) = self.Timestamp.as_ref() {
            os.write_tag(3, ::protobuf::wire_format::WireTypeLengthDelimited)?;
            os.write_raw_varint32(v.get_cached_size())?;
            v.write_to_with_cached_sizes(os)?;
        }
        os.write_unknown_fields(self.get_unknown_fields())?;
        ::std::result::Result::Ok(())
    }

    fn get_cached_size(&self) -> u32 {
        self.cached_size.get()
    }

    fn get_unknown_fields(&self) -> &::protobuf::UnknownFields {
        &self.unknown_fields
    }

    fn mut_unknown_fields(&mut self) -> &mut ::protobuf::UnknownFields {
        &mut self.unknown_fields
    }

    fn as_any(&self) -> &dyn (::std::any::Any) {
        self as &dyn (::std::any::Any)
    }
    fn as_any_mut(&mut self) -> &mut dyn (::std::any::Any) {
        self as &mut dyn (::std::any::Any)
    }
    fn into_any(self: ::std::boxed::Box<Self>) -> ::std::boxed::Box<dyn (::std::any::Any)> {
        self
    }

    fn descriptor(&self) -> &'static ::protobuf::reflect::MessageDescriptor {
        Self::descriptor_static()
    }

    fn new() -> WriteRequest {
        WriteRequest::new()
    }

    fn descriptor_static() -> &'static ::protobuf::reflect::MessageDescriptor {
        static descriptor: ::protobuf::rt::LazyV2<::protobuf::reflect::MessageDescriptor> = ::protobuf::rt::LazyV2::INIT;
        descriptor.get(|| {
            let mut fields = ::std::vec::Vec::new();
            fields.push(::protobuf::reflect::accessor::make_simple_field_accessor::<_, ::protobuf::types::ProtobufTypeString>(
                "FileName",
                |m: &WriteRequest| { &m.FileName },
                |m: &mut WriteRequest| { &mut m.FileName },
            ));
            fields.push(::protobuf::reflect::accessor::make_simple_field_accessor::<_, ::protobuf::types::ProtobufTypeString>(
                "Value",
                |m: &WriteRequest| { &m.Value },
                |m: &mut WriteRequest| { &mut m.Value },
            ));
            fields.push(::protobuf::reflect::accessor::make_singular_ptr_field_accessor::<_, ::protobuf::types::ProtobufTypeMessage<::protobuf::well_known_types::Timestamp>>(
                "Timestamp",
                |m: &WriteRequest| { &m.Timestamp },
                |m: &mut WriteRequest| { &mut m.Timestamp },
            ));
            ::protobuf::reflect::MessageDescriptor::new_pb_name::<WriteRequest>(
                "WriteRequest",
                fields,
                file_descriptor_proto()
            )
        })
    }

    fn default_instance() -> &'static WriteRequest {
        static instance: ::protobuf::rt::LazyV2<WriteRequest> = ::protobuf::rt::LazyV2::INIT;
        instance.get(WriteRequest::new)
    }
}

impl ::protobuf::Clear for WriteRequest {
    fn clear(&mut self) {
        self.FileName.clear();
        self.Value.clear();
        self.Timestamp.clear();
        self.unknown_fields.clear();
    }
}

impl ::std::fmt::Debug for WriteRequest {
    fn fmt(&self, f: &mut ::std::fmt::Formatter<'_>) -> ::std::fmt::Result {
        ::protobuf::text_format::fmt(self, f)
    }
}

impl ::protobuf::reflect::ProtobufValue for WriteRequest {
    fn as_ref(&self) -> ::protobuf::reflect::ReflectValueRef {
        ::protobuf::reflect::ReflectValueRef::Message(self)
    }
}

#[derive(PartialEq,Clone,Default)]
pub struct ReadRequest {
    // message fields
    pub FileName: ::std::string::String,
    // special fields
    pub unknown_fields: ::protobuf::UnknownFields,
    pub cached_size: ::protobuf::CachedSize,
}

impl<'a> ::std::default::Default for &'a ReadRequest {
    fn default() -> &'a ReadRequest {
        <ReadRequest as ::protobuf::Message>::default_instance()
    }
}

impl ReadRequest {
    pub fn new() -> ReadRequest {
        ::std::default::Default::default()
    }

    // string FileName = 1;


    pub fn get_FileName(&self) -> &str {
        &self.FileName
    }
    pub fn clear_FileName(&mut self) {
        self.FileName.clear();
    }

    // Param is passed by value, moved
    pub fn set_FileName(&mut self, v: ::std::string::String) {
        self.FileName = v;
    }

    // Mutable pointer to the field.
    // If field is not initialized, it is initialized with default value first.
    pub fn mut_FileName(&mut self) -> &mut ::std::string::String {
        &mut self.FileName
    }

    // Take field
    pub fn take_FileName(&mut self) -> ::std::string::String {
        ::std::mem::replace(&mut self.FileName, ::std::string::String::new())
    }
}

impl ::protobuf::Message for ReadRequest {
    fn is_initialized(&self) -> bool {
        true
    }

    fn merge_from(&mut self, is: &mut ::protobuf::CodedInputStream<'_>) -> ::protobuf::ProtobufResult<()> {
        while !is.eof()? {
            let (field_number, wire_type) = is.read_tag_unpack()?;
            match field_number {
                1 => {
                    ::protobuf::rt::read_singular_proto3_string_into(wire_type, is, &mut self.FileName)?;
                },
                _ => {
                    ::protobuf::rt::read_unknown_or_skip_group(field_number, wire_type, is, self.mut_unknown_fields())?;
                },
            };
        }
        ::std::result::Result::Ok(())
    }

    // Compute sizes of nested messages
    #[allow(unused_variables)]
    fn compute_size(&self) -> u32 {
        let mut my_size = 0;
        if !self.FileName.is_empty() {
            my_size += ::protobuf::rt::string_size(1, &self.FileName);
        }
        my_size += ::protobuf::rt::unknown_fields_size(self.get_unknown_fields());
        self.cached_size.set(my_size);
        my_size
    }

    fn write_to_with_cached_sizes(&self, os: &mut ::protobuf::CodedOutputStream<'_>) -> ::protobuf::ProtobufResult<()> {
        if !self.FileName.is_empty() {
            os.write_string(1, &self.FileName)?;
        }
        os.write_unknown_fields(self.get_unknown_fields())?;
        ::std::result::Result::Ok(())
    }

    fn get_cached_size(&self) -> u32 {
        self.cached_size.get()
    }

    fn get_unknown_fields(&self) -> &::protobuf::UnknownFields {
        &self.unknown_fields
    }

    fn mut_unknown_fields(&mut self) -> &mut ::protobuf::UnknownFields {
        &mut self.unknown_fields
    }

    fn as_any(&self) -> &dyn (::std::any::Any) {
        self as &dyn (::std::any::Any)
    }
    fn as_any_mut(&mut self) -> &mut dyn (::std::any::Any) {
        self as &mut dyn (::std::any::Any)
    }
    fn into_any(self: ::std::boxed::Box<Self>) -> ::std::boxed::Box<dyn (::std::any::Any)> {
        self
    }

    fn descriptor(&self) -> &'static ::protobuf::reflect::MessageDescriptor {
        Self::descriptor_static()
    }

    fn new() -> ReadRequest {
        ReadRequest::new()
    }

    fn descriptor_static() -> &'static ::protobuf::reflect::MessageDescriptor {
        static descriptor: ::protobuf::rt::LazyV2<::protobuf::reflect::MessageDescriptor> = ::protobuf::rt::LazyV2::INIT;
        descriptor.get(|| {
            let mut fields = ::std::vec::Vec::new();
            fields.push(::protobuf::reflect::accessor::make_simple_field_accessor::<_, ::protobuf::types::ProtobufTypeString>(
                "FileName",
                |m: &ReadRequest| { &m.FileName },
                |m: &mut ReadRequest| { &mut m.FileName },
            ));
            ::protobuf::reflect::MessageDescriptor::new_pb_name::<ReadRequest>(
                "ReadRequest",
                fields,
                file_descriptor_proto()
            )
        })
    }

    fn default_instance() -> &'static ReadRequest {
        static instance: ::protobuf::rt::LazyV2<ReadRequest> = ::protobuf::rt::LazyV2::INIT;
        instance.get(ReadRequest::new)
    }
}

impl ::protobuf::Clear for ReadRequest {
    fn clear(&mut self) {
        self.FileName.clear();
        self.unknown_fields.clear();
    }
}

impl ::std::fmt::Debug for ReadRequest {
    fn fmt(&self, f: &mut ::std::fmt::Formatter<'_>) -> ::std::fmt::Result {
        ::protobuf::text_format::fmt(self, f)
    }
}

impl ::protobuf::reflect::ProtobufValue for ReadRequest {
    fn as_ref(&self) -> ::protobuf::reflect::ReflectValueRef {
        ::protobuf::reflect::ReflectValueRef::Message(self)
    }
}

#[derive(PartialEq,Clone,Default)]
pub struct WriteResponse {
    // message fields
    pub Ok: i32,
    // special fields
    pub unknown_fields: ::protobuf::UnknownFields,
    pub cached_size: ::protobuf::CachedSize,
}

impl<'a> ::std::default::Default for &'a WriteResponse {
    fn default() -> &'a WriteResponse {
        <WriteResponse as ::protobuf::Message>::default_instance()
    }
}

impl WriteResponse {
    pub fn new() -> WriteResponse {
        ::std::default::Default::default()
    }

    // int32 Ok = 1;


    pub fn get_Ok(&self) -> i32 {
        self.Ok
    }
    pub fn clear_Ok(&mut self) {
        self.Ok = 0;
    }

    // Param is passed by value, moved
    pub fn set_Ok(&mut self, v: i32) {
        self.Ok = v;
    }
}

impl ::protobuf::Message for WriteResponse {
    fn is_initialized(&self) -> bool {
        true
    }

    fn merge_from(&mut self, is: &mut ::protobuf::CodedInputStream<'_>) -> ::protobuf::ProtobufResult<()> {
        while !is.eof()? {
            let (field_number, wire_type) = is.read_tag_unpack()?;
            match field_number {
                1 => {
                    if wire_type != ::protobuf::wire_format::WireTypeVarint {
                        return ::std::result::Result::Err(::protobuf::rt::unexpected_wire_type(wire_type));
                    }
                    let tmp = is.read_int32()?;
                    self.Ok = tmp;
                },
                _ => {
                    ::protobuf::rt::read_unknown_or_skip_group(field_number, wire_type, is, self.mut_unknown_fields())?;
                },
            };
        }
        ::std::result::Result::Ok(())
    }

    // Compute sizes of nested messages
    #[allow(unused_variables)]
    fn compute_size(&self) -> u32 {
        let mut my_size = 0;
        if self.Ok != 0 {
            my_size += ::protobuf::rt::value_size(1, self.Ok, ::protobuf::wire_format::WireTypeVarint);
        }
        my_size += ::protobuf::rt::unknown_fields_size(self.get_unknown_fields());
        self.cached_size.set(my_size);
        my_size
    }

    fn write_to_with_cached_sizes(&self, os: &mut ::protobuf::CodedOutputStream<'_>) -> ::protobuf::ProtobufResult<()> {
        if self.Ok != 0 {
            os.write_int32(1, self.Ok)?;
        }
        os.write_unknown_fields(self.get_unknown_fields())?;
        ::std::result::Result::Ok(())
    }

    fn get_cached_size(&self) -> u32 {
        self.cached_size.get()
    }

    fn get_unknown_fields(&self) -> &::protobuf::UnknownFields {
        &self.unknown_fields
    }

    fn mut_unknown_fields(&mut self) -> &mut ::protobuf::UnknownFields {
        &mut self.unknown_fields
    }

    fn as_any(&self) -> &dyn (::std::any::Any) {
        self as &dyn (::std::any::Any)
    }
    fn as_any_mut(&mut self) -> &mut dyn (::std::any::Any) {
        self as &mut dyn (::std::any::Any)
    }
    fn into_any(self: ::std::boxed::Box<Self>) -> ::std::boxed::Box<dyn (::std::any::Any)> {
        self
    }

    fn descriptor(&self) -> &'static ::protobuf::reflect::MessageDescriptor {
        Self::descriptor_static()
    }

    fn new() -> WriteResponse {
        WriteResponse::new()
    }

    fn descriptor_static() -> &'static ::protobuf::reflect::MessageDescriptor {
        static descriptor: ::protobuf::rt::LazyV2<::protobuf::reflect::MessageDescriptor> = ::protobuf::rt::LazyV2::INIT;
        descriptor.get(|| {
            let mut fields = ::std::vec::Vec::new();
            fields.push(::protobuf::reflect::accessor::make_simple_field_accessor::<_, ::protobuf::types::ProtobufTypeInt32>(
                "Ok",
                |m: &WriteResponse| { &m.Ok },
                |m: &mut WriteResponse| { &mut m.Ok },
            ));
            ::protobuf::reflect::MessageDescriptor::new_pb_name::<WriteResponse>(
                "WriteResponse",
                fields,
                file_descriptor_proto()
            )
        })
    }

    fn default_instance() -> &'static WriteResponse {
        static instance: ::protobuf::rt::LazyV2<WriteResponse> = ::protobuf::rt::LazyV2::INIT;
        instance.get(WriteResponse::new)
    }
}

impl ::protobuf::Clear for WriteResponse {
    fn clear(&mut self) {
        self.Ok = 0;
        self.unknown_fields.clear();
    }
}

impl ::std::fmt::Debug for WriteResponse {
    fn fmt(&self, f: &mut ::std::fmt::Formatter<'_>) -> ::std::fmt::Result {
        ::protobuf::text_format::fmt(self, f)
    }
}

impl ::protobuf::reflect::ProtobufValue for WriteResponse {
    fn as_ref(&self) -> ::protobuf::reflect::ReflectValueRef {
        ::protobuf::reflect::ReflectValueRef::Message(self)
    }
}

#[derive(PartialEq,Clone,Default)]
pub struct ReadResponse {
    // message fields
    pub Value: ::std::string::String,
    pub Timestamp: ::protobuf::SingularPtrField<::protobuf::well_known_types::Timestamp>,
    // special fields
    pub unknown_fields: ::protobuf::UnknownFields,
    pub cached_size: ::protobuf::CachedSize,
}

impl<'a> ::std::default::Default for &'a ReadResponse {
    fn default() -> &'a ReadResponse {
        <ReadResponse as ::protobuf::Message>::default_instance()
    }
}

impl ReadResponse {
    pub fn new() -> ReadResponse {
        ::std::default::Default::default()
    }

    // string Value = 1;


    pub fn get_Value(&self) -> &str {
        &self.Value
    }
    pub fn clear_Value(&mut self) {
        self.Value.clear();
    }

    // Param is passed by value, moved
    pub fn set_Value(&mut self, v: ::std::string::String) {
        self.Value = v;
    }

    // Mutable pointer to the field.
    // If field is not initialized, it is initialized with default value first.
    pub fn mut_Value(&mut self) -> &mut ::std::string::String {
        &mut self.Value
    }

    // Take field
    pub fn take_Value(&mut self) -> ::std::string::String {
        ::std::mem::replace(&mut self.Value, ::std::string::String::new())
    }

    // .google.protobuf.Timestamp Timestamp = 3;


    pub fn get_Timestamp(&self) -> &::protobuf::well_known_types::Timestamp {
        self.Timestamp.as_ref().unwrap_or_else(|| <::protobuf::well_known_types::Timestamp as ::protobuf::Message>::default_instance())
    }
    pub fn clear_Timestamp(&mut self) {
        self.Timestamp.clear();
    }

    pub fn has_Timestamp(&self) -> bool {
        self.Timestamp.is_some()
    }

    // Param is passed by value, moved
    pub fn set_Timestamp(&mut self, v: ::protobuf::well_known_types::Timestamp) {
        self.Timestamp = ::protobuf::SingularPtrField::some(v);
    }

    // Mutable pointer to the field.
    // If field is not initialized, it is initialized with default value first.
    pub fn mut_Timestamp(&mut self) -> &mut ::protobuf::well_known_types::Timestamp {
        if self.Timestamp.is_none() {
            self.Timestamp.set_default();
        }
        self.Timestamp.as_mut().unwrap()
    }

    // Take field
    pub fn take_Timestamp(&mut self) -> ::protobuf::well_known_types::Timestamp {
        self.Timestamp.take().unwrap_or_else(|| ::protobuf::well_known_types::Timestamp::new())
    }
}

impl ::protobuf::Message for ReadResponse {
    fn is_initialized(&self) -> bool {
        for v in &self.Timestamp {
            if !v.is_initialized() {
                return false;
            }
        };
        true
    }

    fn merge_from(&mut self, is: &mut ::protobuf::CodedInputStream<'_>) -> ::protobuf::ProtobufResult<()> {
        while !is.eof()? {
            let (field_number, wire_type) = is.read_tag_unpack()?;
            match field_number {
                1 => {
                    ::protobuf::rt::read_singular_proto3_string_into(wire_type, is, &mut self.Value)?;
                },
                3 => {
                    ::protobuf::rt::read_singular_message_into(wire_type, is, &mut self.Timestamp)?;
                },
                _ => {
                    ::protobuf::rt::read_unknown_or_skip_group(field_number, wire_type, is, self.mut_unknown_fields())?;
                },
            };
        }
        ::std::result::Result::Ok(())
    }

    // Compute sizes of nested messages
    #[allow(unused_variables)]
    fn compute_size(&self) -> u32 {
        let mut my_size = 0;
        if !self.Value.is_empty() {
            my_size += ::protobuf::rt::string_size(1, &self.Value);
        }
        if let Some(ref v) = self.Timestamp.as_ref() {
            let len = v.compute_size();
            my_size += 1 + ::protobuf::rt::compute_raw_varint32_size(len) + len;
        }
        my_size += ::protobuf::rt::unknown_fields_size(self.get_unknown_fields());
        self.cached_size.set(my_size);
        my_size
    }

    fn write_to_with_cached_sizes(&self, os: &mut ::protobuf::CodedOutputStream<'_>) -> ::protobuf::ProtobufResult<()> {
        if !self.Value.is_empty() {
            os.write_string(1, &self.Value)?;
        }
        if let Some(ref v) = self.Timestamp.as_ref() {
            os.write_tag(3, ::protobuf::wire_format::WireTypeLengthDelimited)?;
            os.write_raw_varint32(v.get_cached_size())?;
            v.write_to_with_cached_sizes(os)?;
        }
        os.write_unknown_fields(self.get_unknown_fields())?;
        ::std::result::Result::Ok(())
    }

    fn get_cached_size(&self) -> u32 {
        self.cached_size.get()
    }

    fn get_unknown_fields(&self) -> &::protobuf::UnknownFields {
        &self.unknown_fields
    }

    fn mut_unknown_fields(&mut self) -> &mut ::protobuf::UnknownFields {
        &mut self.unknown_fields
    }

    fn as_any(&self) -> &dyn (::std::any::Any) {
        self as &dyn (::std::any::Any)
    }
    fn as_any_mut(&mut self) -> &mut dyn (::std::any::Any) {
        self as &mut dyn (::std::any::Any)
    }
    fn into_any(self: ::std::boxed::Box<Self>) -> ::std::boxed::Box<dyn (::std::any::Any)> {
        self
    }

    fn descriptor(&self) -> &'static ::protobuf::reflect::MessageDescriptor {
        Self::descriptor_static()
    }

    fn new() -> ReadResponse {
        ReadResponse::new()
    }

    fn descriptor_static() -> &'static ::protobuf::reflect::MessageDescriptor {
        static descriptor: ::protobuf::rt::LazyV2<::protobuf::reflect::MessageDescriptor> = ::protobuf::rt::LazyV2::INIT;
        descriptor.get(|| {
            let mut fields = ::std::vec::Vec::new();
            fields.push(::protobuf::reflect::accessor::make_simple_field_accessor::<_, ::protobuf::types::ProtobufTypeString>(
                "Value",
                |m: &ReadResponse| { &m.Value },
                |m: &mut ReadResponse| { &mut m.Value },
            ));
            fields.push(::protobuf::reflect::accessor::make_singular_ptr_field_accessor::<_, ::protobuf::types::ProtobufTypeMessage<::protobuf::well_known_types::Timestamp>>(
                "Timestamp",
                |m: &ReadResponse| { &m.Timestamp },
                |m: &mut ReadResponse| { &mut m.Timestamp },
            ));
            ::protobuf::reflect::MessageDescriptor::new_pb_name::<ReadResponse>(
                "ReadResponse",
                fields,
                file_descriptor_proto()
            )
        })
    }

    fn default_instance() -> &'static ReadResponse {
        static instance: ::protobuf::rt::LazyV2<ReadResponse> = ::protobuf::rt::LazyV2::INIT;
        instance.get(ReadResponse::new)
    }
}

impl ::protobuf::Clear for ReadResponse {
    fn clear(&mut self) {
        self.Value.clear();
        self.Timestamp.clear();
        self.unknown_fields.clear();
    }
}

impl ::std::fmt::Debug for ReadResponse {
    fn fmt(&self, f: &mut ::std::fmt::Formatter<'_>) -> ::std::fmt::Result {
        ::protobuf::text_format::fmt(self, f)
    }
}

impl ::protobuf::reflect::ProtobufValue for ReadResponse {
    fn as_ref(&self) -> ::protobuf::reflect::ReflectValueRef {
        ::protobuf::reflect::ReflectValueRef::Message(self)
    }
}

static file_descriptor_proto_data: &'static [u8] = b"\
    \n\rstorage.proto\x12\x05proto\x1a\x1fgoogle/protobuf/timestamp.proto\"z\
    \n\x0cWriteRequest\x12\x1a\n\x08FileName\x18\x01\x20\x01(\tR\x08FileName\
    \x12\x14\n\x05Value\x18\x02\x20\x01(\tR\x05Value\x128\n\tTimestamp\x18\
    \x03\x20\x01(\x0b2\x1a.google.protobuf.TimestampR\tTimestamp\")\n\x0bRea\
    dRequest\x12\x1a\n\x08FileName\x18\x01\x20\x01(\tR\x08FileName\"\x1f\n\r\
    WriteResponse\x12\x0e\n\x02Ok\x18\x01\x20\x01(\x05R\x02Ok\"^\n\x0cReadRe\
    sponse\x12\x14\n\x05Value\x18\x01\x20\x01(\tR\x05Value\x128\n\tTimestamp\
    \x18\x03\x20\x01(\x0b2\x1a.google.protobuf.TimestampR\tTimestamp2n\n\x07\
    Storage\x12/\n\x04Read\x12\x12.proto.ReadRequest\x1a\x13.proto.ReadRespo\
    nse\x122\n\x05Write\x12\x13.proto.WriteRequest\x1a\x14.proto.WriteRespon\
    seB\x16Z\x14storage_server/protob\x06proto3\
";

static file_descriptor_proto_lazy: ::protobuf::rt::LazyV2<::protobuf::descriptor::FileDescriptorProto> = ::protobuf::rt::LazyV2::INIT;

fn parse_descriptor_proto() -> ::protobuf::descriptor::FileDescriptorProto {
    ::protobuf::Message::parse_from_bytes(file_descriptor_proto_data).unwrap()
}

pub fn file_descriptor_proto() -> &'static ::protobuf::descriptor::FileDescriptorProto {
    file_descriptor_proto_lazy.get(|| {
        parse_descriptor_proto()
    })
}