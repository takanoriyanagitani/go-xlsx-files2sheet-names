package xfiles2sheets

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"iter"
	"os"
)

type FilePath = string

type SheetNames = []string

type XFilePathToSheetNames func(context.Context, FilePath) (SheetNames, error)

type XFilePaths iter.Seq[FilePath]

type NamesWriter func(context.Context, SheetNames) error

type NamesWriterJSON struct{ *json.Encoder }

func (j NamesWriterJSON) ToWriter() NamesWriter {
	return func(_ context.Context, names SheetNames) error {
		return j.Encoder.Encode(names)
	}
}

type SheetNamesInfo struct {
	XlsxFilePath string   `json:"xpath"`
	SheetNames   []string `json:"names"`
}

func (j NamesWriterJSON) ToWriterWithPath(xpath string) NamesWriter {
	return func(_ context.Context, names SheetNames) error {
		info := SheetNamesInfo{
			XlsxFilePath: xpath,
			SheetNames:   names,
		}
		return j.Encoder.Encode(info)
	}
}

type Writer struct{ io.Writer }

func (w Writer) ToWriterJSON() NamesWriterJSON {
	return NamesWriterJSON{json.NewEncoder(w.Writer)}
}

func ReaderToPaths(rdr io.Reader) XFilePaths {
	return func(yield func(string) bool) {
		var scn *bufio.Scanner = bufio.NewScanner(rdr)
		for scn.Scan() {
			var filename string = scn.Text()
			if !yield(filename) {
				return
			}
		}
	}
}

func StdinToPaths() XFilePaths { return ReaderToPaths(os.Stdin) }

func (x XFilePathToSheetNames) WriteAllJSON(
	ctx context.Context,
	paths XFilePaths,
	jwtr NamesWriterJSON,
) error {
	for xpat := range paths {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		names, err := x(ctx, xpat)
		if nil != err {
			return err
		}

		wtr := jwtr.ToWriterWithPath(xpat)

		err = wtr(ctx, names)
		if nil != err {
			return err
		}
	}
	return nil
}

type PathsToNamesWriter struct {
	NamesWriter
	XFilePathToSheetNames
}

func (w PathsToNamesWriter) WriteAll(
	ctx context.Context,
	paths XFilePaths,
) error {
	for xpat := range paths {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		names, err := w.XFilePathToSheetNames(ctx, xpat)
		if nil != err {
			return err
		}

		err = w.NamesWriter(ctx, names)
		if nil != err {
			return err
		}
	}
	return nil
}
