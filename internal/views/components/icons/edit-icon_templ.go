// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.747
package icons

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func EditIcon() templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 24 24\" width=\"15\" height=\"15\"><defs><style>\n                .a{fill:none;stroke:#ca8a04;stroke-linecap:round;stroke-linejoin:round;}\n            </style></defs><title>pencil-2</title><polygon class=\"a\" points=\"7 21.5 0.5 23.5 2.5 17 15.33 4.169 19.83 8.669 7 21.5\"></polygon> <path class=\"a\" d=\"M15.33,4.169l3.086-3.086a2.007,2.007,0,0,1,2.828,0l1.672,1.672a2,2,0,0,1,0,2.828L19.83,8.669\"></path> <line class=\"a\" x1=\"17.58\" y1=\"6.419\" x2=\"6\" y2=\"18\"></line><polyline class=\"a\" points=\"2.5 17 3.5 18 6 18 6 20.5 7 21.5\"></polyline> <line class=\"a\" x1=\"1.5\" y1=\"20.5\" x2=\"3.5\" y2=\"22.5\"></line><line class=\"a\" x1=\"16.83\" y1=\"2.669\" x2=\"21.33\" y2=\"7.169\"></line></svg>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}