# MusGen
MusGen is a code generator for the binary MUS format with validation support.

# Overview
MusGen generates 3 methods for a type description and language:
- Marshal(buf) - encodes data into the MUS format. It will crash if the buffer
  length < Size().
- Unmarshal(buf) - decodes data from the MUS format + performs validation. May 
  return an error if encounters an invalid format or data is not valid.
- Size() - returns the size of the data in the MUS format.

The type description consists of field descriptions, each of which may contain:
- Skip flag
- Validator - name of the function, that validates the field.
- MaxLength - positive number. Restricts length of the string, list, array, or
  map field. Those data types are encoded with a length and value.
  There is no need for further decoding if the length of the field is bigger
  than MaxLength, validation error will be returned.
- ElemValidator - name of the function, that validates elements of the list, 
  array or map field.
- KeyValidator - name of the function, that validates keys of the map field.

Also, the type description has:
- Unsafe flag - if it's true, every string will be decoded from the buffer 
  using the faster, unsafe method.
- Suffix - defines a suffix for Marshal/Unmarhsal/Size methods.

# Supported languages
Go

# MUS format
MUS is an acronym for "Marshal, Unmarshal, Size". The emphasis here is made on 
the presence of the Size function, which has become an important part of the 
Marshal/Unmarshal process.
It's a super simple binary format.

## Features
- Fields are encoded/decoded by order, so there are no field names, only field
  values.
- Varint encoding is used for integers and floats. It is safe, because
  there is a maximum number of bytes limit for each of these types.
- Binary representation of the float is reversed before encoding.
- ZigZag encoding is used for signed to unsigned integer mapping.
- Strings, lists, arrays, and maps are encoded with length (int type) and 
  value.
- Booleans and bytes are encoded by a single byte.
- Pointers are value-encoded.
- Nil pointers are not supported. Encoding will crash if meets a nil pointer.

## Samples
| Type            |     Value           |     MUS format (hex)                  |     In Parts          |
|-----------------|---------------------|---------------------------------------|-----------------------|
| int             | 500                 | e807                                  | <sub>e807 - value of the integer</sub> |
| list of strings | {"hello", "world"}  | 040a68656c6c6f0a776f726c64            | <sub>04 - length of the list,<br>0a - length of the first elem,<br>68656c6c6f - value of the first elem,<br>0a - length of the second elem,<br>776f726c64 - value of the second elem.</sub> |
| Person {<br>??????Name string<br>??????Age int<br>} | {<br>??????Name: "Bill",<br>??????Age: 35,<br>} | 0842696c6c46 | <sub>08 - length of the Name field,<br>42696c6c - value of the Name field,<br>46 - value of the Age field.</sub> |

## Versioning
There is no explicit versioning support. But you can always do next:
```
// Add version field.
TypeV1 {        
  Version byte
  ...
}

TypeV2 {
  Version byte
  ...
}

// Check version field before Unmarshal.
if buff[0] == 1 {
  typeV1.Unmarshal(buff)
} else if buff[0] == 2 {
  typeV2.Unmarshal(buff)
} else {
  return UnsupportedVersionError
}
```

Moreover, it is highly recommended to have a `Version` field. With it, you 
will always be ready for changes in the MUS format of your type.
