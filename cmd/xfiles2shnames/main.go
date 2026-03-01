package main

import (
	"bufio"
	"context"
	"errors"
	"log"
	"os"

	x2s "github.com/takanoriyanagitani/go-xlsx-files2sheet-names"
	xui "github.com/takanoriyanagitani/go-xlsx-files2sheet-names/xfile2sheets/xuri"
)

//nolint:gochecknoglobals
var conv x2s.XFilePathToSheetNames = xui.FilePathToSheetsDefault

func subBuf(ctx context.Context) error {
	var wtr *bufio.Writer = bufio.NewWriter(os.Stdout)

	swtr := x2s.Writer{Writer: wtr}
	var jwtr x2s.NamesWriterJSON = swtr.ToWriterJSON()

	var paths x2s.XFilePaths = x2s.StdinToPaths()

	err := conv.WriteAllJSON(
		ctx,
		paths,
		jwtr,
	)

	return errors.Join(err, wtr.Flush())
}

func sub(ctx context.Context) error {
	swtr := x2s.Writer{Writer: os.Stdout}
	var jwtr x2s.NamesWriterJSON = swtr.ToWriterJSON()

	var paths x2s.XFilePaths = x2s.StdinToPaths()

	err := conv.WriteAllJSON(
		ctx,
		paths,
		jwtr,
	)

	return err
}

func main() {
	err := sub(context.Background())
	if nil != err {
		log.Printf("%v\n", err)
	}
}
