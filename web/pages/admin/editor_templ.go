// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.731
// web/pages/admin/editor.templ

package admin

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"blog-portfolio/internal/models"
	"blog-portfolio/web/layouts"
	"fmt"
)

// PostEditorData holds all the data needed for the post editor
type PostEditorData struct {
	Post  *models.Post // Can be nil for new posts
	Tags  []models.Tag // All available tags
	IsNew bool         // True if creating new post
	Error string       // Any error message to display
}

func PostEditor(data PostEditorData) templ.Component {
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
		templ_7745c5c3_Var2 := templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
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
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-10\"><div class=\"md:flex md:items-center md:justify-between\"><div class=\"flex-1 min-w-0\"><h2 class=\"text-2xl font-bold leading-7 text-neutral-900 dark:text-white sm:text-3xl sm:truncate\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if data.IsNew {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("New Post")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			} else {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("Edit Post")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</h2></div><div class=\"mt-4 flex md:mt-0 md:ml-4\"><button type=\"submit\" form=\"post-form\" name=\"action\" value=\"draft\" class=\"inline-flex items-center px-4 py-2 border border-neutral-300 dark:border-neutral-600 rounded-md shadow-sm text-sm font-medium text-neutral-700 dark:text-neutral-200 bg-white dark:bg-neutral-800 hover:bg-neutral-50 dark:hover:bg-neutral-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500\">Save as Draft</button> <button type=\"submit\" form=\"post-form\" name=\"action\" value=\"publish\" class=\"inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500\">Publish</button></div></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if data.Error != "" {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"mt-6 rounded-md bg-red-50 dark:bg-red-900 p-4\"><div class=\"flex\"><div class=\"flex-shrink-0\"><svg class=\"h-5 w-5 text-red-400\" xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 20 20\" fill=\"currentColor\" aria-hidden=\"true\"><path fill-rule=\"evenodd\" d=\"M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z\" clip-rule=\"evenodd\"></path></svg></div><div class=\"ml-3\"><h3 class=\"text-sm font-medium text-red-800 dark:text-red-200\">Error</h3><div class=\"mt-2 text-sm text-red-700 dark:text-red-300\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var3 string
				templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(data.Error)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/pages/admin/editor.templ`, Line: 61, Col: 22}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></div></div></div>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<form id=\"post-form\" class=\"mt-6 space-y-8\" method=\"POST\" action=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var4 templ.SafeURL = templ.SafeURL(getFormAction(data))
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(string(templ_7745c5c3_Var4)))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if !data.IsNew {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<input type=\"hidden\" name=\"_method\" value=\"PUT\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"space-y-6\"><div><label for=\"title\" class=\"block text-sm font-medium text-neutral-700 dark:text-neutral-300\">Title</label><div class=\"mt-1\"><input type=\"text\" name=\"title\" id=\"title\" required value=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var5 string
			templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(getPostTitle(data))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/pages/admin/editor.templ`, Line: 77, Col: 88}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" class=\"shadow-sm focus:ring-primary-500 focus:border-primary-500 block w-full sm:text-sm border-neutral-300 dark:border-neutral-600 rounded-md dark:bg-neutral-800 dark:text-white\" placeholder=\"Post title\"></div></div><div><label for=\"description\" class=\"block text-sm font-medium text-neutral-700 dark:text-neutral-300\">Description</label><div class=\"mt-1\"><textarea id=\"description\" name=\"description\" rows=\"3\" class=\"shadow-sm focus:ring-primary-500 focus:border-primary-500 block w-full sm:text-sm border-neutral-300 dark:border-neutral-600 rounded-md dark:bg-neutral-800 dark:text-white\" placeholder=\"A brief description of your post\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if data.Post != nil {
				var templ_7745c5c3_Var6 string
				templ_7745c5c3_Var6, templ_7745c5c3_Err = templ.JoinStringErrs(data.Post.Description)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/pages/admin/editor.templ`, Line: 91, Col: 35}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var6))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</textarea></div></div><div><label for=\"content\" class=\"block text-sm font-medium text-neutral-700 dark:text-neutral-300\">Content</label><div class=\"mt-1\"><textarea id=\"content\" name=\"content\" rows=\"20\" required class=\"shadow-sm focus:ring-primary-500 focus:border-primary-500 block w-full sm:text-sm border-neutral-300 dark:border-neutral-600 rounded-md dark:bg-neutral-800 dark:text-white font-mono\" placeholder=\"Write your post content here...\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if data.Post != nil {
				var templ_7745c5c3_Var7 string
				templ_7745c5c3_Var7, templ_7745c5c3_Err = templ.JoinStringErrs(data.Post.Content)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/pages/admin/editor.templ`, Line: 105, Col: 31}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var7))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</textarea></div><p class=\"mt-2 text-sm text-neutral-500 dark:text-neutral-400\">Write your post content using Markdown formatting.</p></div><div><label for=\"cover_image\" class=\"block text-sm font-medium text-neutral-700 dark:text-neutral-300\">Cover Image URL</label><div class=\"mt-1\"><input type=\"url\" name=\"cover_image\" id=\"cover_image\" value=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var8 string
			templ_7745c5c3_Var8, templ_7745c5c3_Err = templ.JoinStringErrs(getPostCoverImage(data))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/pages/admin/editor.templ`, Line: 118, Col: 95}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var8))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" class=\"shadow-sm focus:ring-primary-500 focus:border-primary-500 block w-full sm:text-sm border-neutral-300 dark:border-neutral-600 rounded-md dark:bg-neutral-800 dark:text-white\" placeholder=\"https://example.com/image.jpg\"></div></div><div><label class=\"block text-sm font-medium text-neutral-700 dark:text-neutral-300\">Tags</label><div class=\"mt-2 space-y-2\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			for _, tag := range data.Tags {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"relative flex items-start\"><div class=\"flex items-center h-5\"><input id=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var9 string
				templ_7745c5c3_Var9, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("tag-%d", tag.ID))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/pages/admin/editor.templ`, Line: 131, Col: 55}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var9))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" name=\"tags[]\" type=\"checkbox\" value=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var10 string
				templ_7745c5c3_Var10, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("%d",
					tag.ID))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/pages/admin/editor.templ`, Line: 132, Col: 23}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var10))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				if data.Post != nil && hasTag(data.Post.Tags, tag) {
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" checked")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" class=\"focus:ring-primary-500 h-4 w-4 text-primary-600 border-neutral-300 dark:border-neutral-600 rounded\"></div><div class=\"ml-3 text-sm\"><label for=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var11 string
				templ_7745c5c3_Var11, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("tag-%d", tag.ID))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/pages/admin/editor.templ`, Line: 136, Col: 56}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var11))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" class=\"text-neutral-700 dark:text-neutral-300\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var12 string
				templ_7745c5c3_Var12, templ_7745c5c3_Err = templ.JoinStringErrs(tag.Name)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/pages/admin/editor.templ`, Line: 137, Col: 26}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var12))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</label></div></div>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></div></div></form></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			return templ_7745c5c3_Err
		})
		templ_7745c5c3_Err = layouts.Admin(layouts.PageData{
			Title:       getEditorTitle(data),
			Description: "Create or edit a blog post",
		}).Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

// Helper functions to handle conditional logic outside the template
func getEditorTitle(data PostEditorData) string {
	if data.IsNew {
		return "New Post"
	}
	return "Edit Post: " + data.Post.Title
}

func getFormAction(data PostEditorData) string {
	if data.IsNew {
		return "/admin/posts"
	}
	return fmt.Sprintf("/admin/posts/%d", data.Post.ID)
}

func getPostTitle(data PostEditorData) string {
	if data.Post != nil {
		return data.Post.Title
	}
	return ""
}

func getPostCoverImage(data PostEditorData) string {
	if data.Post != nil {
		return data.Post.CoverImage
	}
	return ""
}

// hasTag checks if a post has a specific tag
func hasTag(postTags []models.Tag, tag models.Tag) bool {
	for _, t := range postTags {
		if t.ID == tag.ID {
			return true
		}
	}
	return false
}
