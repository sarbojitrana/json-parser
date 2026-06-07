package utils

import (
	"fmt"
	"strconv"

	"github.com/xuri/excelize/v2"
)

type Attribute struct {
	SectionName string
	SectionID   int64

	AttributeID string
	Label       string

	InputType string

	FieldSetID *int64
}

type Spreadsheet struct {
	Attributes map[string]Attribute
}

func LoadSpreadsheet(path string) (*Spreadsheet, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, fmt.Errorf("open spreadsheet: %w", err)
	}
	defer f.Close()

	sheet := f.GetSheetName(0)

	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("read rows: %w", err)
	}

	spreadsheet := &Spreadsheet{
		Attributes: make(map[string]Attribute),
	}

	for i, row := range rows {

		if i == 0 {
			continue
		}

		if len(row) < 6 {
			continue
		}

		var sectionID int64
		if row[2] != "" {
			sectionID, _ = strconv.ParseInt(row[2], 10, 64)
		}

		var fieldSetID *int64
		if len(row) > 6 && row[6] != "" {

			id, err := strconv.ParseInt(row[6], 10, 64)
			if err == nil {
				fieldSetID = &id
			}
		}

		attr := Attribute{
			SectionName: row[1],
			SectionID:   sectionID,

			AttributeID: row[3],
			Label:       row[4],

			InputType: row[5],

			FieldSetID: fieldSetID,
		}

		if attr.AttributeID == "" {
			continue
		}

		spreadsheet.Attributes[attr.AttributeID] = attr
	}

	return spreadsheet, nil
}

func (s *Spreadsheet) Get(attributeID string) (Attribute, bool) {
	attr, ok := s.Attributes[attributeID]
	return attr, ok
}