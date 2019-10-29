[![Go Report Card](https://goreportcard.com/badge/github.com/ReneBoedker/algobra)](https://goreportcard.com/report/github.com/ReneBoedker/algobra)
![coverage-badge](https://img.shields.io/badge/coverage-100.0%25-brightgreen?cacheSeconds=86400&style=flat)
[![GoDoc](https://godoc.org/github.com/ReneBoedker/algobra/errors?status.svg)](https://godoc.org/github.com/ReneBoedker/algobra/errors)
# Algobra: Errors
This package implements the error handling used in Algobra. Some users may find the `Is` function useful to test which kind of error occurred. Apart from this, the package is mainly intended for internal use.

## Third party licenses
Parts of the error package are heavily inspired by the corresponding package in [Upspin](https://github.com/upspin/upspin), which is published under the following license:
```
Copyright (c) 2016 The Upspin Authors. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google Inc. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
```
