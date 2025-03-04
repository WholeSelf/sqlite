// Copyright (c) 2019 David Crawshaw <david@zentus.com>
//
// Permission to use, copy, modify, and distribute this software for any
// purpose with or without fee is hereby granted, provided that the above
// copyright notice and this permission notice appear in all copies.
//
// THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
// WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
// MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
// ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
// WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
// ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
// OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.

package sqlitex

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/wholeself/sqlite"
)

// InsertRandID executes stmt with a random value in the range [min, max) for $param.
func InsertRandID(stmt *sqlite.Stmt, param string, min, max int64) (id int64, err error) {
	if min < 0 {
		return 0, fmt.Errorf("sqlitex.InsertRandID: min (%d) is negative", min)
	}

	for i := 0; ; i++ {
		v, err := rand.Int(rand.Reader, big.NewInt(max-min))
		if err != nil {
			return 0, fmt.Errorf("sqlitex.InsertRandID: %v", err)
		}
		id := v.Int64() + min

		stmt.Reset()
		stmt.SetInt64(param, id)
		_, err = stmt.Step()
		if err == nil {
			return id, err
		}
		if sqErr, _ := err.(*sqlite.Error); sqErr != nil {
			if sqErr.Code == sqlite.SQLITE_CONSTRAINT_PRIMARYKEY {
				if i < 100 {
					continue
				}
			}
			sqErr.Loc = "sqlitex.InsertRandID: " + sqErr.Loc
			return 0, sqErr
		}
		return 0, fmt.Errorf("sqlitex.InsertRandID: %v", err)
	}
}
