/*
 * Copyright (C) 2016 Lawrence Woodman <lwoodman@vlifesystems.com>
 */
package main

import "github.com/lawrencewoodman/dlit"

type Expression interface {
	IsTrue([]*dlit.Literal) bool
}
