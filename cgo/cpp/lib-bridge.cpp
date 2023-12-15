#include <iostream>

#include "namegen.hpp"
#include "lib-bridge.h"

using ng = dasmig::ng;

// Utility function local to the bridge's implementation
const char* wstring_to_char(const std::wstring& wstr);

// Manually load the resources folder if necessary.
void LIB_Load(const char* resource_path) {
  ng::instance().load(resource_path);
}

const char* LIB_GetName(int sex, int cult) {
  ng::gender s = static_cast<ng::gender> (sex);
  ng::culture c = static_cast<ng::culture> (cult);

  std::wstring name = ng::instance().get_name(s, c).append_surname();
  return wstring_to_char(name);
}

const char* wstring_to_char(const std::wstring& wstr) {
  // Get the length of the std::wstring
  size_t wlen = wstr.length();

  // Allocate a buffer to store the converted string
  // The buffer size should be at least wlen + 1
  char * buffer = new char[wlen + 1];

  // Convert the std::wstring to a multibyte character string
  // The return value is the number of bytes written to the buffer
  size_t mlen = std::wcstombs(buffer, wstr.c_str(), wlen + 1);

  // Check if the conversion was successful
  if (mlen == static_cast<size_t>(-1)) {
    // Delete the buffer
    delete[] buffer;
    return NULL;
  }
  
  // Return the converted string
  return buffer;
}
