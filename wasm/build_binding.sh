#!/bin/bash

set -e

EXPORTED_FUNCTIONS=_malloc,_free,_uregex_close_68,_uregex_open_68,_uregex_findNext_68,_uregex_start_68,_uregex_end_68,_uregex_setText_68,_u_strFromUTF8_68,_u_strToUTF8_68

exec emcc \
  -s EXPORTED_FUNCTIONS="$EXPORTED_FUNCTIONS" \
  -s USE_ICU=1 \
  -s WASM=1 \
  --no-entry -o binding/test.wasm \
  binding/test.c \
  -Wl,--whole-archive -licu_i18n -Wl,--no-whole-archive
