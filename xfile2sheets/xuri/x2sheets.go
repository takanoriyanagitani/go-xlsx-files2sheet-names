package x2sheets

import (
	"context"
	"log"

	x2s "github.com/takanoriyanagitani/go-xlsx-files2sheet-names"
	xpkg "github.com/xuri/excelize/v2"
)

type Xfile struct{ *xpkg.File }

func (x Xfile) Close() error { return x.File.Close() }

func (x Xfile) Sheets() []string { return x.File.GetSheetList() }

func FilepathToSheets(_ context.Context, filepath string) ([]string, error) {
	file, err := xpkg.OpenFile(filepath)
	if nil != err {
		return nil, err
	}
	xfile := Xfile{File: file}
	defer func() {
		err := xfile.Close()
		if nil != err {
			log.Printf("%v\n", err)
		}
	}()

	return xfile.Sheets(), nil
}

//nolint:gochecknoglobals
var FilePathToSheetsDefault x2s.XFilePathToSheetNames = FilepathToSheets
