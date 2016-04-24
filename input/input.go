/*
 * Copyright (C) 2016 Lawrence Woodman <lwoodman@vlifesystems.com>
 */
package input

import "github.com/lawrencewoodman/dlit_go"

type Input interface {
	Clone() (Input, error)
	Next() bool
	Err() error
	Read() (map[string]*dlit.Literal, error)
	Rewind() error
	Close() error
}
