// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.543
package partials

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import "bca-go-final/internal/views/components"
import "bca-go-final/internal/types"
import "bca-go-final/internal/utils"

func quantityValText(quantity *types.Quantity, text string) string {
	if quantity == nil {
		return ""
	}
	switch text {
	case "project":
		return quantity.Project.Name
	case "item":
		return quantity.Rubro.Name
	case "quantity":
		return utils.PrintFloat(quantity.Quantity)
	default:
		return ""
	}
}

func concat(s1, s2 string) string {
	return s1 + s2
}

func EditCantidades(quantities *types.Quantity, projects, items []types.Select) templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<form id=\"edit-quantity\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if quantities == nil {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" hx-post=\"/bca/partials/cantidades/add\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" hx-put=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(concat("/bca/partials/cantidades/", quantities.Id.String())))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" autocomplete=\"off\" hx-target=\"#cantidades-table\" hx-target-error=\"#error\" hx-swap=\"innerHTML\" hx-trigger=\"submit\" hx-on=\"htmx:afterOnLoad: handleHtmxError(event)\"><div class=\"flex h-full flex-col gap-4\"><div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if quantities == nil {
			templ_7745c5c3_Err = components.DrawerTitle("Crear Cantidad").Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = components.DrawerTitle("Editar Cantidad").Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div id=\"error\" class=\"text-red-600 text-sm\"></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = components.SelectComponent(projects, "Seleccione un Proyecto", "project", "project", quantityValText(quantities, "project"), "Proyecto").Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = components.SelectComponent(items, "Seleccione un Rubro", "item", "item", quantityValText(quantities, "item"), "Rubro").Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = components.Input("text", "Cantidad", "quantity", "quantity", quantityValText(quantities, "quantity")).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = components.ButtonGroup().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></form><script>\n    function handleHtmxError(event) {\n      document.getElementById(\"error\").innerHTML = \"\"\n\n      if (event.detail.xhr.status === 200) {\n        resetClose()\n        return\n      }\n\n      document.getElementById(\"error\").innerHTML = event.detail.xhr.responseText\n    }\n\n    function resetClose() {\n      closeDrawer()\n    }\n  </script>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}