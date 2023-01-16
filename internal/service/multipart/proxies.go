package multipart

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"

	datasetsservice "github.com/inst-api/parser/gen/datasets_service"
	"github.com/inst-api/parser/pkg/logger"
	"go.uber.org/multierr"
)

func readProxiesList(ctx context.Context, reader io.Reader) ([]*datasetsservice.ProxyRecord, error) {
	csvReader := csv.NewReader(reader)

	csvReader.Comma = ':'
	csvReader.FieldsPerRecord = 4

	var botAccounts []*datasetsservice.ProxyRecord
	var errs []error

	var line int
	for {
		record, err := csvReader.Read()
		line, _ = csvReader.FieldPos(0)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break // readed all file
			}

			errs = append(errs, fmt.Errorf("failed to read from file at line %d: %v", line, err))
			continue
		}

		botAccounts = append(botAccounts, &datasetsservice.ProxyRecord{
			Record:     record,
			LineNumber: line,
		})
	}

	logger.Debugf(ctx, "read %d lines\n", line)

	return botAccounts, multierr.Combine(errs...)
}
