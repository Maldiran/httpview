// httpview is a helper library to increase code reusability in other projects
// It is useful for streamlining rendering HTML and text templates.
// Inspired by https://www.calhoun.io/intro-to-templates-p4-v-in-mvc/
package httpview

import (
	ht "html/template"
	"io"
	"log/slog"
	"net/http"
	"path/filepath"
	tt "text/template"
)

// View holds template configuration
type View struct {
	HTemplate *ht.Template
	TTemplate *tt.Template
	Layout    string
}

// Type for holding template configurations (Views)
type Tmpl map[string]*View

// Global variable for holding template configurations (Views)
var Templates Tmpl

var layoutDir string

var newViewData func(*http.Request, string, any) any

// Sets layout directory where the template base files are stored.
// All of this template files are included in each template.
func SetLayoutDir(dir string) {
	layoutDir = dir
}

// Sets a user-defined function that returns standard struct passed to every template.
// It should accept HTTP request, HTML title and data for additional template values.
func NewViewData(f func(*http.Request, string, any) any) {
	newViewData = f
}

// Executes chosen named template and writes it to provided io.Writer.
// This is the default way of executing templates in this library.
func Execute(r *http.Request, w io.Writer, name string, title string, data any) error {
	viewData := newViewData(r, title, data)
	return Templates[name].RenderHtml(w, viewData)
}

// Creates new View named layoutName as HTML Template
func NewViewHtml(layoutName string, files ...string) (*View, error) {
	v := &View{
		Layout: layoutName,
	}
	v.SetViewHtml(files...)
	return v, nil
}

// Creates new View named layoutName as Text Template
func NewViewText(layoutName string, files ...string) (*View, error) {
	v := &View{
		Layout: layoutName,
	}
	v.SetViewText(files...)
	return v, nil
}

// Sets Text Template to existing View
func (v *View) SetViewText(files ...string) error {
	files = append(layoutFiles("gotxt"), files...)
	t, err := tt.ParseFiles(files...)
	if err != nil {
		slog.Error("Text Template error",
			slog.Any("err", err),
		)
		return err
	}

	v.TTemplate = t
	return nil
}

// Sets HTML Template to existing View
func (v *View) SetViewHtml(files ...string) error {
	files = append(layoutFiles("gohtml"), files...)
	t, err := ht.ParseFiles(files...)
	if err != nil {
		slog.Error("Html Template error",
			slog.Any("err", err),
		)
		return err
	}

	v.HTemplate = t
	return nil
}

// Renders HTML Template from View to io.Writer using custom data struct
func (v *View) RenderHtml(w io.Writer, data any) error {
	return v.HTemplate.ExecuteTemplate(w, v.Layout, data)
}

// Renders Text Template from View to io.Writer using custom data struct
func (v *View) RenderText(w io.Writer, data any) error {
	return v.TTemplate.ExecuteTemplate(w, v.Layout, data)
}

func layoutFiles(ext string) []string {
	files, err := filepath.Glob(layoutDir + "/*." + ext)
	if err != nil {
		panic(err)
	}
	return files
}
