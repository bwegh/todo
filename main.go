package main

import (
  "os"
  "gopkg.in/alecthomas/kingpin.v2"
  "github.com/svenfuchs/todo.go/cmd"
  "github.com/svenfuchs/todo.go/item"
  "github.com/svenfuchs/todo.go/date"
  // "fmt"
)

type Opts struct {
  path string
  format string
  line string
  id int
  status string
  text string
  projects []string
  date string
  mode string
}

func (o Opts) filter() item.Filter {
  o.date = date.Normalize(o.date, date.Time)
  return item.Filter { o.id, o.status, o.text, o.projects, o.date, o.mode }
}

func (o *Opts) setDate(date string)   { o.mode, o.date = "date", date }
func (o *Opts) setBefore(date string) { o.mode, o.date = "before", date }
func (o *Opts) setSince(date string)  { o.mode, o.date = "since", date }
func (o *Opts) setAfter(date string)  { o.mode, o.date = "after", date  }

func (o *Opts) setStatus(status string) {
  if len(status) > 4 {
    status = status[0:4]
  }
  o.status = status
}

func (o *Opts) setLine(line string) {
  p := item.Parser { line }
  o.id = p.Id()
  o.status = p.Status()
  o.text = p.Text()
}

func (o *Opts) runList(c *kingpin.ParseContext) error {
  cmd.NewList(o.path, o.filter(), o.format).Run()
	return nil
}

func (o *Opts) runToggle(c *kingpin.ParseContext) error {
  cmd.NewToggle(o.path, o.filter()).Run()
	return nil
}

func main() {
  app := kingpin.New("todo", "A command-line todo.txt tool.")

	h := &Opts{}
	c := app.Command("list", "Filter and list todo items."      ).Action(h.runList)
	c.Flag("file",     "Todo.txt file to work with."            ).Short('f').StringVar(&h.path)
	c.Flag("format",   "Output format."                         ).Short('o').StringVar(&h.format)
	c.Flag("id",       "Filter by id."                          ).Short('i').IntVar(&h.id)
	c.Flag("status",   "Filter by status."                      ).Short('s').SetValue(&funcValue { h.setStatus })
	c.Flag("text",     "Filter by text."                        ).Short('t').StringVar(&h.text)
	c.Flag("projects", "Filter by projects (comma separated)."  ).Short('p').StringsVar(&h.projects)
	c.Flag("date",     "Filter by done date."                   ).Short('a').SetValue(&funcValue { h.setDate   })
	c.Flag("after",    "Filter by done after."                  ).Short('a').SetValue(&funcValue { h.setSince  })
	c.Flag("since",    "Filter by done since."                  ).Short('n').SetValue(&funcValue { h.setBefore })
	c.Flag("before",   "Filter by done before."                 ).Short('b').SetValue(&funcValue { h.setBefore })
	c.Arg("input",     "Filter by full line."                   ).SetValue(&funcValue { h.setLine })

	h = &Opts{}
	c = app.Command("toggle", "Toggle todo items."              ).Action(h.runToggle)
	c.Flag("file",     "Todo.txt file to work with."            ).Short('f').StringVar(&h.path)
	c.Flag("id",       "Filter by id."                          ).Short('i').IntVar(&h.id)
	c.Flag("status",   "Filter by status."                      ).Short('s').StringVar(&h.status)
	c.Flag("text",     "Filter by text."                        ).Short('t').StringVar(&h.text)
	c.Flag("projects", "Filter by projects (comma separated)."  ).Short('p').StringsVar(&h.projects)
	c.Flag("date",     "Filter by done date."                   ).Short('a').SetValue(&funcValue { h.setDate   })
	c.Flag("after",    "Filter by done after."                  ).Short('a').SetValue(&funcValue { h.setSince  })
	c.Flag("since",    "Filter by done since."                  ).Short('n').SetValue(&funcValue { h.setBefore })
	c.Flag("before",   "Filter by done before."                 ).Short('b').SetValue(&funcValue { h.setBefore })
	c.Arg("input",     "Filter by full line."                   ).SetValue(&funcValue { h.setLine })

  kingpin.MustParse(app.Parse(os.Args[1:]))
}


// This uses an internal (?) kingpin api to set values through function
// callbacks instead of directly forcing them onto a struct field

type funcValue struct{
  f func(string)
}

func (f *funcValue) Set(value string) error {
	value, err := value, error(nil)
	if err != nil { return err }
  f.f((string)(value))
  return nil
}

func (f *funcValue) String() string {
  return "no idea?"
}