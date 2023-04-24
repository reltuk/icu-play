#!/bin/bash

set -e

exec emcc \
  -s EXPORTED_FUNCTIONS=_malloc,_free,_uregex_close_68 \
  -s USE_ICU=1 \
  -s WASM=1 \
  --no-entry -o binding/test.wasm \
  binding/test.c \
  -Wl,--whole-archive -licu_i18n -Wl,--no-whole-archive
