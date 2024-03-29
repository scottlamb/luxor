// Package client is a concrete client implementation of the FX Luminaire
// Luxor ZD wi-fi module's protocol using JSON-over-HTTP.
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/scottlamb/luxor/protocol"
	"io/ioutil"
	"net/http"
)

// *Controller implements protocol.Controller
type Controller struct {
	BaseURL string
}

// request issues a request for method with prefilled request and ready-to-fill
// response. It returns error on JSON- or HTTP-level problems; it does not
// check the Status field in the response.
func (c *Controller) request(ctx context.Context, method string, request interface{}, response interface{}) (err error) {
	serializedReq, err := json.Marshal(request)
	if err != nil {
		return err
	}
	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL+"/"+method+".json", bytes.NewReader(serializedReq))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()
	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return err
	}
	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected HTTP status: %q with body: %q", httpResp.Status, body)
	}
	if contentType := httpResp.Header.Get("Content-Type"); contentType != "application/json" {
		return fmt.Errorf("Unexpected response content type: %q with body: %q", contentType, body)
	}

	err = json.Unmarshal(body, response)
	if err != nil {
		return fmt.Errorf("JSON error: %v while parsing body: %q", err, body)
	}
	return nil
}

// The methods below are all boilerplate.

func (c *Controller) AssignLight(ctx context.Context, req *protocol.AssignLightRequest) (*protocol.AssignLightResponse, error) {
	resp := &protocol.AssignLightResponse{}
	if err := c.request(ctx, "AssignLight", req, resp); err != nil {
		return nil, err
	}
	return resp, protocol.ErrorForStatus(resp.Status)
}

func (c *Controller) ControllerName(ctx context.Context, req *protocol.ControllerNameRequest) (*protocol.ControllerNameResponse, error) {
	resp := &protocol.ControllerNameResponse{}
	if err := c.request(ctx, "ControllerName", req, resp); err != nil {
		return nil, err
	}
	return resp, protocol.ErrorForStatus(resp.Status)
}

func (c *Controller) ExtinguishAll(ctx context.Context, req *protocol.ExtinguishAllRequest) (*protocol.ExtinguishAllResponse, error) {
	resp := &protocol.ExtinguishAllResponse{}
	if err := c.request(ctx, "ExtinguishAll", req, resp); err != nil {
		return nil, err
	}
	return resp, protocol.ErrorForStatus(resp.Status)
}

func (c *Controller) FlashLights(ctx context.Context, req *protocol.FlashLightsRequest) (*protocol.FlashLightsResponse, error) {
	resp := &protocol.FlashLightsResponse{}
	if err := c.request(ctx, "FlashLights", req, resp); err != nil {
		return nil, err
	}
	return resp, protocol.ErrorForStatus(resp.Status)
}

func (c *Controller) GroupListAdd(ctx context.Context, req *protocol.GroupListAddRequest) (*protocol.GroupListAddResponse, error) {
	resp := &protocol.GroupListAddResponse{}
	if err := c.request(ctx, "GroupListAdd", req, resp); err != nil {
		return nil, err
	}
	return resp, protocol.ErrorForStatus(resp.Status)
}

func (c *Controller) GroupListClear(ctx context.Context, req *protocol.GroupListClearRequest) (*protocol.GroupListClearResponse, error) {
	resp := &protocol.GroupListClearResponse{}
	if err := c.request(ctx, "GroupListClear", req, resp); err != nil {
		return nil, err
	}
	return resp, protocol.ErrorForStatus(resp.Status)
}

func (c *Controller) GroupListDelete(ctx context.Context, req *protocol.GroupListDeleteRequest) (*protocol.GroupListDeleteResponse, error) {
	resp := &protocol.GroupListDeleteResponse{}
	if err := c.request(ctx, "GroupListDelete", req, resp); err != nil {
		return nil, err
	}
	return resp, protocol.ErrorForStatus(resp.Status)
}

func (c *Controller) GroupListGet(ctx context.Context, req *protocol.GroupListGetRequest) (*protocol.GroupListGetResponse, error) {
	resp := &protocol.GroupListGetResponse{}
	if err := c.request(ctx, "GroupListGet", req, resp); err != nil {
		return nil, err
	}
	return resp, protocol.ErrorForStatus(resp.Status)
}

func (c *Controller) GroupListRename(ctx context.Context, req *protocol.GroupListRenameRequest) (*protocol.GroupListRenameResponse, error) {
	resp := &protocol.GroupListRenameResponse{}
	if err := c.request(ctx, "GroupListRename", req, resp); err != nil {
		return nil, err
	}
	return resp, protocol.ErrorForStatus(resp.Status)
}

func (c *Controller) GroupListReorder(ctx context.Context, req *protocol.GroupListReorderRequest) (*protocol.GroupListReorderResponse, error) {
	resp := &protocol.GroupListReorderResponse{}
	if err := c.request(ctx, "GroupListReorder", req, resp); err != nil {
		return nil, err
	}
	return resp, protocol.ErrorForStatus(resp.Status)
}

func (c *Controller) IlluminateAll(ctx context.Context, req *protocol.IlluminateAllRequest) (*protocol.IlluminateAllResponse, error) {
	resp := &protocol.IlluminateAllResponse{}
	if err := c.request(ctx, "IlluminateAll", req, resp); err != nil {
		return nil, err
	}
	return resp, protocol.ErrorForStatus(resp.Status)
}

func (c *Controller) IlluminateGroup(ctx context.Context, req *protocol.IlluminateGroupRequest) (*protocol.IlluminateGroupResponse, error) {
	resp := &protocol.IlluminateGroupResponse{}
	if err := c.request(ctx, "IlluminateGroup", req, resp); err != nil {
		return nil, err
	}
	return resp, protocol.ErrorForStatus(resp.Status)
}

func (c *Controller) IlluminateTheme(ctx context.Context, req *protocol.IlluminateThemeRequest) (*protocol.IlluminateThemeResponse, error) {
	resp := &protocol.IlluminateThemeResponse{}
	if err := c.request(ctx, "IlluminateTheme", req, resp); err != nil {
		return nil, err
	}
	return resp, protocol.ErrorForStatus(resp.Status)
}

func (c *Controller) ThemeClear(ctx context.Context, req *protocol.ThemeClearRequest) (*protocol.ThemeClearResponse, error) {
	resp := &protocol.ThemeClearResponse{}
	if err := c.request(ctx, "ThemeClear", req, resp); err != nil {
		return nil, err
	}
	return resp, protocol.ErrorForStatus(resp.Status)
}

func (c *Controller) ThemeGet(ctx context.Context, req *protocol.ThemeGetRequest) (*protocol.ThemeGetResponse, error) {
	resp := &protocol.ThemeGetResponse{}
	if err := c.request(ctx, "ThemeGet", req, resp); err != nil {
		return nil, err
	}
	return resp, protocol.ErrorForStatus(resp.Status)
}

func (c *Controller) ThemeListAdd(ctx context.Context, req *protocol.ThemeListAddRequest) (*protocol.ThemeListAddResponse, error) {
	resp := &protocol.ThemeListAddResponse{}
	if err := c.request(ctx, "ThemeListAdd", req, resp); err != nil {
		return nil, err
	}
	return resp, protocol.ErrorForStatus(resp.Status)
}

func (c *Controller) ThemeListClear(ctx context.Context, req *protocol.ThemeListClearRequest) (*protocol.ThemeListClearResponse, error) {
	resp := &protocol.ThemeListClearResponse{}
	if err := c.request(ctx, "ThemeListClear", req, resp); err != nil {
		return nil, err
	}
	return resp, protocol.ErrorForStatus(resp.Status)
}

func (c *Controller) ThemeListDelete(ctx context.Context, req *protocol.ThemeListDeleteRequest) (*protocol.ThemeListDeleteResponse, error) {
	resp := &protocol.ThemeListDeleteResponse{}
	if err := c.request(ctx, "ThemeListDelete", req, resp); err != nil {
		return nil, err
	}
	return resp, protocol.ErrorForStatus(resp.Status)
}

func (c *Controller) ThemeListGet(ctx context.Context, req *protocol.ThemeListGetRequest) (*protocol.ThemeListGetResponse, error) {
	resp := &protocol.ThemeListGetResponse{}
	if err := c.request(ctx, "ThemeListGet", req, resp); err != nil {
		return nil, err
	}
	return resp, protocol.ErrorForStatus(resp.Status)
}

func (c *Controller) ThemeListRename(ctx context.Context, req *protocol.ThemeListRenameRequest) (*protocol.ThemeListRenameResponse, error) {
	resp := &protocol.ThemeListRenameResponse{}
	if err := c.request(ctx, "ThemeListRename", req, resp); err != nil {
		return nil, err
	}
	return resp, protocol.ErrorForStatus(resp.Status)
}

func (c *Controller) ThemeListReorder(ctx context.Context, req *protocol.ThemeListReorderRequest) (*protocol.ThemeListReorderResponse, error) {
	resp := &protocol.ThemeListReorderResponse{}
	if err := c.request(ctx, "ThemeListReorder", req, resp); err != nil {
		return nil, err
	}
	return resp, protocol.ErrorForStatus(resp.Status)
}

func (c *Controller) ThemeSet(ctx context.Context, req *protocol.ThemeSetRequest) (*protocol.ThemeSetResponse, error) {
	resp := &protocol.ThemeSetResponse{}
	if err := c.request(ctx, "ThemeSet", req, resp); err != nil {
		return nil, err
	}
	return resp, protocol.ErrorForStatus(resp.Status)
}

// Ensure *Controller implements protocol.Controller.
var _ protocol.Controller = (*Controller)(nil)
