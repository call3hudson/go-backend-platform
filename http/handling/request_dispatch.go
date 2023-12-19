package handling

import (
	"fmt"
	"io"
	"net/http"
	"platform/http/actionresults"
	"platform/http/handling/params"
	"platform/pipeline"
	"platform/services"
	"reflect"
	"strings"
)

func NewRouter(handlers ...HandlerEntry) *RouterComponent {
	return &RouterComponent{generateRoutes(handlers...)}
}

type RouterComponent struct {
	routes []Route
}

func (router *RouterComponent) Init() {}

func (router *RouterComponent) ProcessRequest(context *pipeline.ComponentContext, next func(*pipeline.ComponentContext)) {
	for _, route := range router.routes {
		if strings.EqualFold(context.Request.Method, route.httpMethod) {
			matches := route.expression.FindAllStringSubmatch(context.URL.Path, -1)
			if len(matches) > 0 {
				rawParamsVals := []string{}
				if len(matches[0]) > 1 {
					rawParamsVals = matches[0][1:]
				}
				err := router.invokeHandler(route, rawParamsVals, context)
				if err == nil {
					next(context)
				} else {
					context.Error(err)
				}
				return
			}
		}
	}
	context.ResponseWriter.WriteHeader(http.StatusNotFound)
}

func (router *RouterComponent) invokeHandler(route Route, rawParams []string, context *pipeline.ComponentContext) error {
	paramVals, err := params.GetParametersFromRequest(context.Request, route.handlerMethod, rawParams)
	if err == nil {
		structVal := reflect.New(route.handlerMethod.Type.In(0))
		services.PopulateForContext(context.Context(), structVal.Interface())
		paramVals = append([]reflect.Value{structVal.Elem()}, paramVals...)
		result := route.handlerMethod.Func.Call(paramVals)
		if len(result) > 0 {
			if action, ok := result[0].Interface().(actionresults.ActionResult); ok {
				invoker := createInvokehandlerFunc(context.Context(), router.routes)
				err = services.PopulateForContextWithExtras(context.Context(), action, map[reflect.Type]reflect.Value{
					reflect.TypeOf(invoker): reflect.ValueOf(invoker),
				})
				if err == nil {
					err = action.Execute(&actionresults.ActionContext{
						Context: context.Context(), ResponseWriter: context.ResponseWriter})
				}
			} else {
				io.WriteString(context.ResponseWriter, fmt.Sprint(result[0].Interface()))
			}
		}
	}
	return err
}
