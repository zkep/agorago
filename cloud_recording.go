package agora

import (
	"fmt"
	"net/http"
)

// 云端录制
// https://docs.agora.io/cn/cloud-recording/restfulapi/#/%E4%BA%91%E7%AB%AF%E5%BD%95%E5%88%B6

const (
	AGORA_API                         = "https://api.agora.io/"
	CLOUD_RECORDING_ACQUIRE_URL       = AGORA_API + "v1/apps/%s/cloud_recording/acquire"
	CLOUD_RECORDING_START_URL         = AGORA_API + "v1/apps/%s/cloud_recording/resourceid/%s/mode/%s/start"
	CLOUD_RECORDING_STOP_URL          = AGORA_API + "v1/apps/%s/cloud_recording/resourceid/%s/sid/%s/mode/%s/stop"
	CLOUD_RECORDING_QUERY_URL         = AGORA_API + "v1/apps/%s/cloud_recording/resourceid/%s/sid/%s/mode/%s/query"
	CLOUD_RECORDING_UPDATE_URL        = AGORA_API + "v1/apps/%s/cloud_recording/resourceid/%s/sid/%s/mode/%s/update"
	CLOUD_RECORDING_UPDATE_LAYOUT_URL = AGORA_API + "v1/apps/%s/cloud_recording/resourceid/%s/sid/%s/mode/%s/updateLayout"
)

type CloudRecording struct {
	req *Request
}

type RecordOption func(c *CloudRecording)

func AddRequest(req *Request) RecordOption {
	return func(c *CloudRecording) {
		c.req = req
	}
}

func NewCloudRecording(opts ...RecordOption) *CloudRecording {
	recording := &CloudRecording{}
	for _, opt := range opts {
		opt(recording)
	}
	return recording
}

// 获取resource ID
func (self *CloudRecording) Acquire(req CommonRequest, ret *AcquireResponse) error {
	uri := fmt.Sprintf(CLOUD_RECORDING_ACQUIRE_URL, self.req.appid)
	err := self.req.Do(uri, http.MethodPost, req, nil, nil, ret)
	if err != nil {
		return err
	}
	return nil
}

// 开启云端录制
func (self *CloudRecording) Start(resourceId, mode string, req StartRequest, ret *StartResponse) error {
	uri := fmt.Sprintf(CLOUD_RECORDING_START_URL, self.req.appid, resourceId, mode)
	err := self.req.Do(uri, http.MethodPost, req, nil, nil, ret)
	if err != nil {
		return err
	}
	return nil
}

// 停止云端录制 sid 通过 start 请求获取的录制 ID
func (self *CloudRecording) Stop(resourceId, sid, mode string, req CommonRequest, ret *StopResponse) error {
	uri := fmt.Sprintf(CLOUD_RECORDING_STOP_URL, self.req.appid, resourceId, sid, mode)
	req.ClientRequest = struct{}{}
	err := self.req.Do(uri, http.MethodPost, req, nil, nil, ret)
	if err != nil {
		return err
	}
	return nil
}

// 更新订阅名单
func (self *CloudRecording) Update(resourceId, sid, mode string, req UpdateRequest, ret *UpdateResponse) error {
	uri := fmt.Sprintf(CLOUD_RECORDING_UPDATE_URL, self.req.appid, resourceId, sid, mode)
	err := self.req.Do(uri, http.MethodPost, req, nil, nil, ret)
	if err != nil {
		return err
	}
	return nil
}

// 更新合流布局
func (self *CloudRecording) UpdateLayOut(resourceId, sid, mode string, req UpdateLayOutRequest, ret *UpdateResponse) error {
	uri := fmt.Sprintf(CLOUD_RECORDING_UPDATE_LAYOUT_URL, self.req.appid, resourceId, sid, mode)
	err := self.req.Do(uri, http.MethodPost, req, nil, nil, ret)
	if err != nil {
		return err
	}
	return nil
}

// 查询云端录制状态
func (self *CloudRecording) Query(resourceId, sid, mode string, ret *QueryResponse) error {
	uri := fmt.Sprintf(CLOUD_RECORDING_QUERY_URL, self.req.appid, resourceId, sid, mode)
	err := self.req.Do(uri, http.MethodGet, nil, nil, nil, ret)
	if err != nil {
		return err
	}
	return nil
}
