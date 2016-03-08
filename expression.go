/*
 * Copyright (C) 2016 Lawrence Woodman <lwoodman@vlifesystems.com>
 */
package main

import "github.com/lawrencewoodman/dlit_go"

type Expression interface {
	IsTrue([]*dlit.Literal) bool
}
