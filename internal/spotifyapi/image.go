package spotifyapi

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/topvennie/sortifyr/pkg/image"
)

func (c *client) ImageGet(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("new http request %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get image data %s | %w", url, err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("unexpected status %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read image data %s | %w", url, err)
	}

	webp, err := image.ToWebp(data)
	if err != nil {
		return nil, fmt.Errorf("convert image to webp %s | %w", url, err)
	}

	return webp, nil
}
