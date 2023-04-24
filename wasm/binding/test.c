#include <unistd.h>
#include <stdlib.h>
#include <emscripten.h>
#include "unicode/ustring.h"
#include "unicode/utext.h"
#include "unicode/uregex.h"

#ifdef __cplusplus
extern "C" {
#endif
  
typedef struct icu_UCharString {
  UChar *ptr;
  int len;
} icu_UCharString;

void EMSCRIPTEN_KEEPALIVE icu_ucharstring_free(icu_UCharString *p) {
  free(p->ptr);
  free(p);
}

char *EMSCRIPTEN_KEEPALIVE icu_ucharstring_substr_toUTF8(icu_UCharString *str, int start, int end, int *len) {
  UErrorCode uerr = U_ZERO_ERROR;
  u_strToUTF8(NULL, 0, len, str->ptr + start, end - start, &uerr);
  if (uerr != U_BUFFER_OVERFLOW_ERROR) {
    exit(uerr);
  }

  char *ret = (char *)malloc(*len);
  uerr = U_ZERO_ERROR;
  u_strToUTF8(ret, *len, NULL, str->ptr + start, end - start, &uerr);
  if (uerr > U_ZERO_ERROR){
    exit(uerr);
  }

  return ret;
}

char *EMSCRIPTEN_KEEPALIVE icu_ucharstring_toUTF8(icu_UCharString *str, int *len) {
  return icu_ucharstring_substr_toUTF8(str, 0, str->len, len);
}

icu_UCharString *EMSCRIPTEN_KEEPALIVE icu_ucharstring_fromUTF8(void *chars, int len) {
  UErrorCode uerr = U_ZERO_ERROR;
  int retlen;
  u_strFromUTF8(NULL, 0, &retlen, (const char *)chars, len, &uerr);
  if (uerr != U_BUFFER_OVERFLOW_ERROR) {
    exit(uerr);
  }

  icu_UCharString *ret = (icu_UCharString *)malloc(sizeof(icu_UCharString));
  ret->len = retlen;
  ret->ptr = (UChar *)malloc(retlen * sizeof(UChar));
  uerr = U_ZERO_ERROR;
  u_strFromUTF8(ret->ptr, retlen, NULL, (const char *)chars, len, &uerr);
  if (uerr > U_ZERO_ERROR) {
    exit(uerr);
  }

  return ret;
}

URegularExpression *EMSCRIPTEN_KEEPALIVE icu_uregex_open(icu_UCharString *str, uint32_t flags) {
  UErrorCode uerr = U_ZERO_ERROR;
  URegularExpression *ret = uregex_open(str->ptr, str->len, flags, NULL, &uerr);
  if (uerr > U_ZERO_ERROR) {
    exit(uerr);
  }
  return ret;
}

void EMSCRIPTEN_KEEPALIVE icu_uregex_setText(URegularExpression *regex, icu_UCharString *str) {
  UErrorCode uerr = U_ZERO_ERROR;
  uregex_setText(regex, str->ptr, str->len, &uerr);
  if (uerr > U_ZERO_ERROR) {
    exit(uerr);
  }
}

UBool EMSCRIPTEN_KEEPALIVE icu_uregex_findNext(URegularExpression *regex) {
  UErrorCode uerr = U_ZERO_ERROR;
  UBool res = uregex_findNext(regex, &uerr);
  if (uerr > U_ZERO_ERROR) {
    exit(uerr);
  }
  return res;
}

int32_t EMSCRIPTEN_KEEPALIVE icu_uregex_start(URegularExpression *regex, int32_t group) {
  UErrorCode uerr = U_ZERO_ERROR;
  int32_t res = uregex_start(regex, group, &uerr);
  if (uerr > U_ZERO_ERROR) {
    exit(uerr);
  }
  return res;
}

int32_t EMSCRIPTEN_KEEPALIVE icu_uregex_end(URegularExpression *regex, int32_t group) {
  UErrorCode uerr = U_ZERO_ERROR;
  int32_t res = uregex_end(regex, group, &uerr);
  if (uerr > U_ZERO_ERROR) {
    exit(uerr);
  }
  return res;
}

#ifdef __cplusplus
} // extern "C"
#endif
