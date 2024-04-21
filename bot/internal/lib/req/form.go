package req

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func WithMultipartFormData(ctx context.Context, typeReq string, url string, reqHeaders map[string]string, body *bytes.Buffer, expBody any) error {
	req, err := http.NewRequestWithContext(ctx, typeReq, url, body)
	if err != nil {
		return err
	}

	for key, value := range reqHeaders {
		req.Header.Add(key, value)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if strings.Contains(resp.Header.Get("Content-Type"), "json") {
		err = json.NewDecoder(resp.Body).Decode(&expBody)
		if err != nil {
			return err
		}
		return nil
	}

	respText, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	htmlResp, ok := expBody.(*string)
	if !ok {
		return fmt.Errorf("expected the ptr to the string")
	}

	*htmlResp = string(respText)
	return nil
}
