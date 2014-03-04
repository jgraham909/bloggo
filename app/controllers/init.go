package controllers

import "github.com/revel/revel"
import "github.com/jgraham909/revmgo"

func init() {
	revmgo.ControllerInit()
	revel.InterceptMethod((*Application).Setup, revel.BEFORE)
}
