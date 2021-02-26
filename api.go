package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/xerrors"
)

func call(ctx context.Context, token, method string, dst interface{}) error {
	uri := fmt.Sprintf("https://api.telegram.org/bot%s/%s", token, method)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return xerrors.Errorf("build request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return xerrors.Errorf("execute request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return xerrors.Errorf("read body: %w", err)
	}

	var response struct {
		Ok          bool            `json:"ok"`
		Description string          `json:"description"`
		Result      json.RawMessage `json:"result"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return xerrors.Errorf("unmarshal body: %w", err)
	}

	if !response.Ok {
		return xerrors.Errorf("call error: %s (%d)", response.Description, res.StatusCode)
	}

	if dst != nil {
		if err := json.Unmarshal(response.Result, dst); err != nil {
			return xerrors.Errorf("unmarshal result: %w", err)
		}
	}

	return nil
}
