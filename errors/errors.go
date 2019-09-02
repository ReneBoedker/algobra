package errors

import (
	"fmt"
	"strings"
)

type Op string

type Kind uint8

const (
	Other                Kind = iota // General error type
	Input                            // General input error
	InputIncompatible                // Inputs are incompatible
	InputTooLarge                    // Given input exceeds some upper bound
	Parsing                          // General parsing error
	ParsingVariableNames             // Parsing encountered unexpected variable names
	Overflow                         // Overflow error
	Internal                         // Internal error
)

type Error struct {
	Op   Op    // The operation causing the error
	Kind Kind  // The kind of error
	Err  error // The underlying error
}

func New(op Op, kind Kind, message string, formatArgs ...interface{}) *Error {
	return &Error{
		Op:   op,
		Kind: kind,
		Err:  fmt.Errorf(message, formatArgs...),
	}
}

func Wrap(op Op, kind Kind, err error) *Error {
	return &Error{
		Op:   op,
		Kind: kind,
		Err:  err,
	}
}

func Is(kind Kind, err error) bool {
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	if e.Kind != Other {
		return e.Kind == kind
	}
	if e.Err != nil {
		return Is(kind, e.Err)
	}
	return false
}

func (e *Error) Error() string {
	var sb strings.Builder

	if e.Op != "" {
		fmt.Fprintf(&sb, "%s: ", e.Op)
	}

	if e.Err != nil {
		fmt.Fprint(&sb, e.Err.Error())
	}

	return sb.String()
}

/* Copyright 2019 René Bødker Christensen
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice, this
 *    list of conditions and the following disclaimer.
 *
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *
 * 3. Neither the name of the copyright holder nor the names of its contributors
 *    may be used to endorse or promote products derived from this software
 *    without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
 * CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
 * OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */
