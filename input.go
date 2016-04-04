/*
 * Copyright (C) 2016 Lawrence Woodman <lwoodman@vlifesystems.com>
 */
package rulehunter

import "github.com/lawrencewoodman/dlit_go"

type Input interface {
	Clone() (Input, error)
	Read() (map[string]*dlit.Literal, error)
	Rewind() error
	// TODO: Add Close()
}
