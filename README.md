# Musgen

Musgen is a code generator for binary MUS format with validation support.

# Overview

Musgen generates for the specific type description and language 3 methods:
- MarshalMus(buf) - encodes data into MUS format. It will crash if the buffer 
  length < SizeMUS().
- UnmarshalMUS(buf) - decodes data from MUS format + performs validation. May 
  return error if the invalid format was met, or data is not valid.
- SizeMUS() - returns the size of the data in MUS format.

Type description consists of field descriptions. Each of which may contain:
- Skip flag
- Validator - name of the function, which receives a value of the field and may
  return error.
- MaxLength - positive number. Restricts length of the string, list, array, or 
  map field. Those data types are encoded with length and value.
  There is no need for further decoding if the length of the field is bigger 
  than MaxLength, validation error will be returned.
- ElemValidator - name of the function, which validates elements of the 
  list, array or map field.
- KeyValidator - name of the function, which validates keys of the map field.

Also, the type description has an unsafe flag. If a new buffer will be created 
for every Marshal(Unmarshaоl) of that type, unsafe flag could be set to true. 
In this case, generated code will be faster. Otherwise, if the same buffer will 
be used or if you want to modify it, the unsafe flag should be set to false.

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

  | Type            |     Value           |     Mus format (hex)                  |     In Parts          |
  |-----------------|---------------------|---------------------------------------|-----------------------|
  | int             | 500                 | e807                                  | <sub>e807 - value of the integer</sub> |
  | list of strings | {"hello", "world"}  | 040a68656c6c6f0a776f726c64            | <sub>04 - length of the list,<br>0a - length of the first elem,<br>68656c6c6f - value of the first elem,<br>0a - length of the second elem,<br>776f726c64 - value of the second elem.</sub> |
  | Person {<br>  Name string<br>  Age int<br>} | {<br>  Name: "Bill",<br>  Age: 35,<br>} | 0842696c6c46 | <sub>08 - length of the Name field,<br>42696c6c - value of the Name field,<br>46 - value of the Age field.</sub> |



  ## Versioning
  There is no explicit versioning support. But you can always do next

  ```
  // Add version field.
  TypeV1 {        
    version byte
    ...
  }

  TypeV2 {
    version byte
    ...
  }

  // Check version field before Unmarshal.
  if buff[0] == 1 {
    typeV1.UnmarshalMUS(buff)
  } else if buff[0] == 2 {
    typeV2.UnmarshalMUS(buff)
  }
  ```

 Moreover, it is highly recommended to have a `version` field. With it, you 
 will always be ready for changes in the MUS format of your type. 
