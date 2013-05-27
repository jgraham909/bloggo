package app

import "github.com/robfig/revel"
import "github.com/jgraham909/revmgo"

func init() {
	revel.OnAppStart(revmgo.Init)
	revel.Filters = []revel.Filter{
		revel.PanicFilter,
		revel.RouterFilter,
		revel.FilterConfiguringFilter,
		revel.ParamsFilter,
		revel.SessionFilter,
		revel.FlashFilter,
		revel.ValidationFilter,
		revel.I18nFilter,
		revel.InterceptorFilter,
		revel.ActionInvoker,
	}
}
