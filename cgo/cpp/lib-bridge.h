#pragma once
#ifdef __cplusplus
extern "C" {
#endif

void LIB_Load(const char* resource_path);
const char* LIB_GetName(int sex, int culture);

#ifdef __cplusplus
}  // extern "C"
#endif
