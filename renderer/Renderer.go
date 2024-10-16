package renderer

import (
	"html/template"
	"io"
	"strings"
)

type Renderer struct {
    Folder string
    Ext    string 
}

func NewRenderer(folder string, ext string) *Renderer {
   return &Renderer{
        Folder: folder,
        Ext: ext,
    } 
}

func (r *Renderer) Render(w io.Writer, tmplName string, data map[string]any, layouts ...string) {
    tmplNames := make([]string, 0, 16)
    tmplNames = append(tmplNames, r.Folder + tmplName + r.Ext)
    name := tmplName
    if layouts != nil {
        splitterName := strings.Split(layouts[0], "/")
        name = splitterName[len(splitterName) - 1]
        for _, l := range layouts {
            tmplNames = append(tmplNames, r.Folder + l + r.Ext)
        }
    }    
    tmpl := template.Must(template.ParseFiles(tmplNames...))
    tmpl.ExecuteTemplate(w, name, data)
}
