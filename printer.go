package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

type Printer interface {
	Print(interface{}) error
}

type JSONPrinter json.Encoder

func NewJSONPrinter(w io.Writer) *JSONPrinter {
	return (*JSONPrinter)(json.NewEncoder(w))
}

func (p *JSONPrinter) Print(v interface{}) error {
	return (*json.Encoder)(p).Encode(v)
}

type TablePrinter tabwriter.Writer

func NewTablePrinter(w io.Writer) *TablePrinter {
	return (*TablePrinter)(tabwriter.NewWriter(w, 4, 0, 2, ' ', tabwriter.Debug))
}

func (p *TablePrinter) Print(v interface{}) error {
	w := (*tabwriter.Writer)(p)
	defer w.Flush()

	summarables, ok := v.(Summarables)
	if !ok {
		return errors.New("can not print value that don't satisfy Tabler interface")
	}

	fmt.Fprintln(w, strings.Join(summarables.Header(), "\t"))

	for _, s := range summarables.Summaries() {
		fmt.Fprintln(w, strings.Join(s.Summary(), "\t"))
	}

	return nil
}

type Summarable interface {
	Summary() []string
}

type Summarables interface {
	Header() []string
	Summaries() []Summarable
}
