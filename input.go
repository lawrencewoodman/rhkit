/*
 * Copyright (C) 2016 Lawrence Woodman <lwoodman@vlifesystems.com>
 */
package main

import "github.com/lawrencewoodman/dlit"

type Input interface {
	Read() (map[string]*dlit.Literal, error)
	// TODO: Add Close()
}
