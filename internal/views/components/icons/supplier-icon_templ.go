// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.543
package icons

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func SupplierIcon() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<svg version=\"1.1\" width=\"15\" height=\"15\" id=\"Layer_1_1_\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" x=\"0px\" y=\"0px\" viewBox=\"0 0 64 64\" style=\"enable-background:new 0 0 64 64;\" xml:space=\"preserve\" fill=\"#90a4ae\"><rect x=\"29\" y=\"35\" width=\"6\" height=\"6\"></rect> <rect x=\"50\" y=\"22\" width=\"6\" height=\"6\"></rect> <rect x=\"8\" y=\"22\" width=\"6\" height=\"6\"></rect> <path d=\"M11,41c0.314,0,0.611-0.148,0.8-0.4l3-4l-1.6-1.2L12,37v-6h-2v6l-1.2-1.6l-1.6,1.2l3,4C10.389,40.852,10.686,41,11,41z\"></path> <path d=\"M53,41c0.314,0,0.611-0.148,0.8-0.4l3-4l-1.6-1.2L54,37v-6h-2v6l-1.2-1.6l-1.6,1.2l3,4C52.389,40.852,52.686,41,53,41z\"></path> <path d=\"M42,16l-8,5v-5l-8,5V10h-4v21h20V16z M28,26h-2v-2h2V26z M33,26h-2v-2h2V26z M38,26h-2v-2h2V26z\"></path> <path d=\"M28,4h7V2h-7c-2.757,0-5,2.243-5,5h2C25,5.346,26.346,4,28,4z\"></path> <rect x=\"37\" y=\"2\" width=\"7\" height=\"2\"></rect> <path d=\"M11,61c4.963,0,9-4.038,9-9s-4.037-9-9-9s-9,4.038-9,9S6.037,61,11,61z M11,45c3.859,0,7,3.14,7,7\n    c0,2.373-1.189,4.47-3,5.736V56c0-1.1-0.9-2-2-2h-2H9c-1.1,0-2,0.9-2,2v1.736C5.189,56.47,4,54.372,4,52C4,48.14,7.141,45,11,45z\"></path> <circle cx=\"11\" cy=\"51\" r=\"3\"></circle> <path d=\"M23,53c0,4.962,4.037,9,9,9s9-4.038,9-9s-4.037-9-9-9S23,48.038,23,53z M32,46c3.859,0,7,3.14,7,7\n    c0,2.373-1.189,4.47-3,5.736V57c0-1.1-0.9-2-2-2h-2h-2c-1.1,0-2,0.9-2,2v1.736c-1.811-1.267-3-3.364-3-5.736\n    C25,49.14,28.141,46,32,46z\"></path> <circle cx=\"32\" cy=\"52\" r=\"3\"></circle> <path d=\"M53,61c4.963,0,9-4.038,9-9s-4.037-9-9-9s-9,4.038-9,9S48.037,61,53,61z M53,45c3.859,0,7,3.14,7,7\n    c0,2.373-1.189,4.47-3,5.736V56c0-1.1-0.9-2-2-2h-2h-2c-1.1,0-2,0.9-2,2v1.736c-1.811-1.267-3-3.364-3-5.736\n    C46,48.14,49.141,45,53,45z\"></path> <circle cx=\"53\" cy=\"51\" r=\"3\"></circle></svg>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

//fill:#90a4ae