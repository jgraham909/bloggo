package app

import (
	"github.com/jgraham909/revmgo"
	"github.com/robfig/revel"
	"labix.org/v2/mgo/bson"
	"reflect"
)

func init() {
	revel.OnAppStart(revmgo.AppInit)
	revel.OnAppStart(AppInit)
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.ActionInvoker,           // Invoke the action.
	}

	revel.TemplateFuncs["nil"] = func(a interface{}) bool {
		return a == nil
	}

	// bson.ObjectId binder
	objId := bson.NewObjectId()
	revel.TypeBinders[reflect.TypeOf(objId)] = ObjectIdBinder
}

var ObjectIdBinder = revel.Binder{
	Bind: revel.ValueBinder(func(val string, typ reflect.Type) reflect.Value {
		// revel.WARN.Print("ObjectIdBinder.Bind - we need an ObjectId here!")
		if len(val) == 0 {
			// revel.WARN.Print("ObjectIdBinder.Bind - length zero - return Zero!")
			return reflect.Zero(typ)
		}
		if bson.IsObjectIdHex(val) {
			// revel.WARN.Print("ObjectIdBinder.Bind - we have a valid ObjectId here!")
			objId := bson.ObjectIdHex(val)
			return reflect.ValueOf(objId)
		} else {
			revel.ERROR.Print("ObjectIdBinder.Bind - Unsure how to handle invalid ObjectId..?")
			return reflect.Zero(typ)
		}
	}),
	Unbind: func(output map[string]string, name string, val interface{}) {
		revel.WARN.Print("ObjectIdBinder.Unbind called! - not sure when this is used, probably when a ObjectId needs to be serialized to a string?!..")
		output[name] = (val.(bson.ObjectId)).Hex()
	},
}
