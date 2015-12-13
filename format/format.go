package format

import (
  "fmt"
  "sort"
  "strings"
  "github.com/svenfuchs/todo.go/item"
)

var formats = map[string][]string{
  "short": []string{ "status", "done", "text" },
  "long":  []string{ "status", "done", "text", "tags" },
  "full":  []string{ "status", "text", "tags", "id" },
}

func New(format string) Format {
  if format == "" {
    format = "short"
  }
  f, ok := formats[format]
  if !ok {
    f = strings.Split(format, ",")
  }
  return Format { format: f }
}

type Format struct {
  format []string
}

func (f Format) Apply(items []item.Item) []string {
  var lines []string
  for _, item := range(items) {
    lines = append(lines, f.fmt(item))
  }
  return lines
}

func (f Format) fmt(item item.Item) string {
  if item.IsNone() {
    return item.Line
  } else {
    return f.fmtItem(item)
  }
}

func (f Format) fmtItem(item item.Item) string {
  var fields []string
  for _, format := range f.format {
    field := formatters[format](item)
    if field != "" {
      fields = append(fields, field)
    }
  }
  return strings.Join(fields, " ")
}

var formatters = map[string]func(item item.Item) string {
  "id":     fmtId,
  "status": fmtStatus,
  "text":   fmtText,
  "tags":   fmtTags,
  "done":   fmtDone,
  "due":    fmtDue,
}

var statuses = map[string]string {
  "done": "x",
  "pend": "-",
}

func fmtStatus(item item.Item) string {
  return statuses[string(item.Status)]
}

func fmtText(item item.Item) string {
  return item.Text
}

func fmtTags(item item.Item) string {
  var tags []string
  for _, key := range sortedKeys(item.Tags) {
    tags = append(tags, fmt.Sprintf("%s:%s", key, item.Tags[key]))
  }
  return strings.Join(tags, " ")
}

func fmtId(item item.Item) string {
  return fmt.Sprintf("[%d]", item.Id)
}

func fmtDone(item item.Item) string {
  return item.DoneDate()
}

func fmtDue(item item.Item) string {
  return item.DueDate()
}

func sortedKeys(m map[string]string) []string {
  var keys []string
  for key := range m {
    keys = append(keys, key)
  }
  sort.Strings(keys)
  return keys
}
