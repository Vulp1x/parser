package multipart

import (
	"context"
	"fmt"
	"mime/multipart"

	datasetsservice "github.com/inst-api/parser/gen/datasets_service"
	"github.com/inst-api/parser/pkg/logger"
)

const (
	botsPartName    = "bots"
	proxiesPartName = "proxies"
)

// DatasetsServiceUploadFilesDecoderFunc implements the multipart decoder for
// service "auth_service" endpoint "upload file". The decoder must populate the
// argument p after encoding.
func DatasetsServiceUploadFilesDecoderFunc(mr *multipart.Reader, p **datasetsservice.UploadFilesPayload) error {
	// Add multipart request decoder logic here
	payload := &datasetsservice.UploadFilesPayload{}

	ctx := context.Background()

	for i := 0; i < 4; i++ {
		part, err := mr.NextPart()
		if err != nil {
			return fmt.Errorf("failed to get next part: %v", err)
		}

		switch part.FormName() {
		case botsPartName:
			payload.Bots, err = readBotsList(ctx, part)
			if err != nil {
				return fmt.Errorf("failed to read bots list: %v", err)
			}
			payload.BotsFilename = part.FileName()
		case proxiesPartName:
			payload.Proxies, err = readProxiesList(ctx, part)
			if err != nil {
				return fmt.Errorf("failed to read proxies list: %v", err)
			}
			payload.ProxiesFilename = part.FileName()
		default:
			return fmt.Errorf("unknown part '%s' expected one of: [%s, %s]", part.FormName(),
				botsPartName, proxiesPartName)
		}
	}

	logger.Infof(ctx, "read %d bots and %d proxies", len(payload.Bots), len(payload.Proxies))

	*p = payload

	return nil
}
