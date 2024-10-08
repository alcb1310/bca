// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.778
package details

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"fmt"
	"github.com/alcb1310/bca/internal/types"
	"github.com/alcb1310/bca/internal/views/components"
)

func EditDetails(budgetItems []types.Select, invoiceId string) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<form hx-post=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("/bca/partials/invoices/%s/details", invoiceId))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/views/bca/transaction/partials/details/editdetails.templ`, Line: 10, Col: 75}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-trigger=\"submit\" hx-target=\"#invoice-details\" hx-swap=\"innerHTML\" hx-on=\"htmx:afterRequest: htmxHandleDetailsError(event)\"><div class=\"flex flex-col h-full gap-8\"><div id=\"details-error\" class=\"text-red-500 text-sm\"></div><div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = components.DrawerTitle("Agregar Detalles").Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = components.SelectComponent(budgetItems, "Seleccione una Partida", "item", "item", "", "Partida").Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = components.Input("text", "Cantidad", "quantity", "quantity", "").Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = components.Input("text", "Costo", "cost", "cost", "").Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = components.Input("text", "Total", "detailtotal", "detailtotal", "").Render(ctx, templ_7745c5c3_Buffer)
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></form><script>\n  var quantity = document.getElementById(\"quantity\")\n  var cost = document.getElementById(\"cost\")\n  var detailTotal = document.getElementById(\"detailtotal\")\n  var b = document.getElementById(\"save-button\")\n\n  detailTotal.disabled = true\n\n  quantity.addEventListener(\"input\", calculateTotal)\n  cost.addEventListener(\"input\", calculateTotal)\n\n  function htmxHandleDetailsError(event) {\n    document.getElementById(\"details-error\").innerHTML = \"\"\n    if (event.detail.xhr.status >= 400) {\n      document.getElementById(\"details-error\").innerHTML = event.detail.xhr.response\n    }\n  }\n\n  function calculateTotal() {\n    let q = quantity.value === \"\" ? 0 : parseFloat(quantity.value)\n    let c = cost.value === \"\" ? 0 : parseFloat(cost.value)\n\n    if (isNaN(q)) {\n      detailTotal.classList.add(\"error-border\")\n      quantity.classList.add(\"error-border\")\n      q = 0\n      detailTotal.value = \"0.00\"\n      b.disabled = true\n      return\n    }\n    if (isNaN(c)) {\n      detailTotal.classList.add(\"error-border\")\n      cost.classList.add(\"error-border\")\n      c = 0\n      detailTotal.value = \"0.00\"\n      b.disabled = true\n      return\n    }\n\n\n    detailTotal.classList.remove(\"error-border\")\n    quantity.classList.remove(\"error-border\")\n    cost.classList.remove(\"error-border\")\n    b.disabled = false\n\n    detailTotal.value = (q * c).toLocaleString(2)\n  }\n\n  function resetClose() {\n    closeDrawer()\n  }\n</script>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate
