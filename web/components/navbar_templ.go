// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.731
// web/components/navbar.templ

package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

type NavbarProps struct {
	IsAdmin bool
}

func Navbar(props NavbarProps) templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<nav class=\"bg-gradient-to-r from-pastel-base to-pastel-warmGray border-b border-pastel-warmGray/50 dark:from-neutral-800 dark:to-neutral-900 dark:border-neutral-700\"><div class=\"container mx-auto px-4\"><div class=\"flex justify-between h-16\"><div class=\"flex\"><div class=\"flex-shrink-0 flex items-center\"><a href=\"/\" class=\"text-xl font-bold text-pastel-text dark:text-white font-mono\">Blog & Portfolio</a></div><div class=\"hidden sm:ml-6 sm:flex sm:space-x-8\"><a href=\"/\" class=\"inline-flex items-center px-1 pt-1 text-sm font-medium text-pastel-text dark:text-white hover:text-neutral-800 dark:hover:text-white transition-colors duration-200\">Home</a> <a href=\"/blog\" class=\"inline-flex items-center px-1 pt-1 text-sm font-medium text-pastel-text/80 dark:text-neutral-400 hover:text-neutral-800 dark:hover:text-white transition-colors duration-200\">Blog</a> <a href=\"/portfolio\" class=\"inline-flex items-center px-1 pt-1 text-sm font-medium text-pastel-text/80 dark:text-neutral-400 hover:text-neutral-800 dark:hover:text-white transition-colors duration-200\">Portfolio</a></div></div><div class=\"flex items-center space-x-4\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if props.IsAdmin {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<a href=\"/admin/dashboard\" class=\"inline-flex items-center px-3 py-1 text-sm font-medium bg-gradient-to-r from-pastel-blue to-pastel-purple text-pastel-text rounded-md hover:opacity-90 transition-opacity duration-200\">Admin Dashboard</a> ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<button type=\"button\" x-data=\"{ darkMode: localStorage.theme === &#39;dark&#39; }\" @click=\"darkMode = !darkMode; localStorage.theme = darkMode ? &#39;dark&#39; : &#39;light&#39;; document.documentElement.classList.toggle(&#39;dark&#39;)\" class=\"p-2 text-pastel-text/70 hover:text-pastel-text dark:text-neutral-400 dark:hover:text-white transition-colors duration-200\"><span x-show=\"!darkMode\" class=\"w-5 h-5\">🌙</span> <span x-show=\"darkMode\" class=\"w-5 h-5\">☀️</span></button></div></div></div></nav>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}
