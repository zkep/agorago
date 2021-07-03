package agorago

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	APP_ID          = ""
	APP_CERTIFICATE = ""

	TOKEN_APP_ID          = ""
	TOKEN_APP_CERTIFICATE = ""

	Kodo_AccessId  = ""
	Kodo_SecretKey = ""
)

var request *Request
var cloudRecording *CloudRecording

func init() {
	credentials := base64.StdEncoding.EncodeToString([]byte(TOKEN_APP_ID + ":" + TOKEN_APP_CERTIFICATE))
	request = NewRequest(SetCredentials(credentials), SetKeyAndSecret(APP_ID, APP_CERTIFICATE))
	cloudRecording = NewCloudRecording(AddRequest(request))
}

func TestProjects(t *testing.T) {
	uri := "https://api.agora.io/dev/v1/projects"
	var ret AgoraProjects
	payload := strings.NewReader(``)
	err := request.Do(uri, http.MethodGet, payload, nil, nil, &ret)
	assert.NoError(t, err)
	fmt.Println(ret)
}

// 获取ResourceId
func TestAcquire(t *testing.T) {
	req := CommonRequest{
		Cname: "httpClient463224",
		UID:   "527841",
		ClientRequest: AcquireClientRequest{
			Region:              "CN",
			ResourceExpiredHour: 24,
		},
	}
	var ret AcquireResponse
	err := cloudRecording.Acquire(req, &ret)
	fmt.Println(ret)
	assert.NoError(t, err)
	assert.NotEmpty(t, ret.ResourceId)
}

// 开启录制
func TestStart(t *testing.T) {
	req := StartRequest{
		Cname:         "httpClient463224",
		UID:           "527841",
		ClientRequest: StartClientRequest{},
	}
	resourceId := ""
	mode := RECORDING_MODE_MIX
	var ret StartResponse
	err := cloudRecording.Start(resourceId, mode, req, &ret)
	fmt.Println(ret)
	assert.NoError(t, err)
	assert.NotEmpty(t, ret.ResourceID)
}

func TestCloudRecording(t *testing.T) {
	// 获取resource ID
	uid := "123456"
	req := CommonRequest{
		Cname: uid,
		UID:   uid,
		ClientRequest: AcquireClientRequest{
			Region:              "CN",
			ResourceExpiredHour: 24,
		},
	}
	var ret AcquireResponse
	err := cloudRecording.Acquire(req, &ret)
	assert.NoError(t, err)
	assert.NotEmpty(t, ret.ResourceId)
	t.Logf("%+v", ret)
	// 开启云端录制
	remoteId := "2222"
	startReq := StartRequest{
		Cname: uid,
		UID:   uid,
		ClientRequest: StartClientRequest{
			RecordingConfig: RecordingConfig{
				MaxIdleTime:     30,
				StreamTypes:     2,
				ChannelType:     0,
				VideoStreamType: 0,
				TranscodingConfig: TranscodingConfig{
					Height:           640,
					Width:            360,
					Bitrate:          500,
					Fps:              15,
					MixedVideoLayout: 1,
					BackgroundColor:  "#FF0000",
				},
				SubscribeAudioUids: []string{
					uid,
					remoteId,
				},
				SubscribeVideoUids: []string{
					uid,
					remoteId,
				},
				SubscribeUIDGroup: 0,
			},
			RecordingFileConfig: RecordingFileConfig{
				AvFileType: []string{"hls"},
			},
			StorageConfig: StorageConfig{
				AccessKey:      Kodo_AccessId,
				SecretKey:      Kodo_SecretKey,
				Region:         1,
				Bucket:         "uneed-file",
				Vendor:         0,
				FileNamePrefix: []string{"recording"},
			},
		},
	}
	mode := RECORDING_MODE_MIX
	var startret StartResponse
	err = cloudRecording.Start(ret.ResourceId, mode, startReq, &startret)
	assert.NoError(t, err)
	t.Logf("%+v", startret)
}

// 结束录制
func TestStop(t *testing.T) {
	req := CommonRequest{
		Cname: "",
		UID:   "",
	}
	sid := ""
	resourceId := ""
	mode := RECORDING_MODE_MIX
	var ret StopResponse
	err := cloudRecording.Stop(resourceId, sid, mode, req, &ret)
	assert.NoError(t, err)
	assert.NotEmpty(t, ret.ResourceID)
}
