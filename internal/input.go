/*
 * Copyright (C) 2016 Lawrence Woodman <lwoodman@vlifesystems.com>
 */
package internal

import "github.com/lawrencewoodman/dlit_go"

type Input interface {
	Read() (map[string]*dlit.Literal, error)
	Rewind() error
	// TODO: Add Close()
}
