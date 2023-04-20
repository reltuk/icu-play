Example of integrating with ICU4C C API using CGo.

Shows:

* Making a *UChar str for a golang string.

* Making a golang string from a (view of a) *UChar str

* Initializing a *URegularExpression and finding all the matches of it against a given *UChar str.

# Building and Running It

You need to install ICU into `third_party/icu_install`. It's simplest if you forego dynamic linking completely.

On a *-nix, something like:

```sh
$ mkdir third_party
$ cd third_party
$ curl -OL https://github.com/unicode-org/icu/releases/download/release-73-1/icu4c-73_1-src.tgz
$ tar zxvf icu4c-73_1-src.tgz
$ cd icu/source
$ ./configure --prefix=`pwd`/../../icu_install --enable-static=yes --enable-shared=no --disable-dyload --with-data-packaging=static --enable-tests=no --enable-samples=no
$ gmake -j10 CXXFLAGS=-std=c++11
$ gmake install
$ cd ../../../
$ go run .
```
