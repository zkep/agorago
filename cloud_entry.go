package agora

import "encoding/json"

// api文档
// https://docs.agora.io/cn/cloud-recording/cloud_recording_api_rest?platform=RESTful#storageConfig

type AgoraProjects struct {
	Success  bool       `json:"success"`
	Projects []Projects `json:"projects"`
}

type Projects struct {
	ID              string      `json:"id"`
	Name            string      `json:"name"`
	Status          int         `json:"status"`
	SignKey         string      `json:"sign_key"`
	VendorKey       string      `json:"vendor_key"`
	RecordingServer interface{} `json:"recording_server"`
	Created         int         `json:"created"`
}

type CommonRequest struct {
	UID           string      `json:"uid"`
	Cname         string      `json:"cname"`
	ClientRequest interface{} `json:"clientRequest"`
}

type AcquireClientRequest struct {
	ResourceExpiredHour int    `json:"resourceExpiredHour"`
	Scene               int    `json:"scene"`
	Region              string `json:"region"`
}

type AcquireResponse struct {
	ResourceId string `json:"resourceId"`
}

type StartRequest struct {
	UID           string             `json:"uid"`
	Cname         string             `json:"cname"`
	ClientRequest StartClientRequest `json:"clientRequest"`
}

type StartClientRequest struct {
	Token               string              `json:"token"`
	RecordingConfig     RecordingConfig     `json:"recordingConfig"`
	RecordingFileConfig RecordingFileConfig `json:"recordingFileConfig"`
	StorageConfig       StorageConfig       `json:"storageConfig"`
}

type TranscodingConfig struct {
	Height           int    `json:"height"`
	Width            int    `json:"width"`
	Bitrate          int    `json:"bitrate"`
	Fps              int    `json:"fps"`
	MixedVideoLayout int    `json:"mixedVideoLayout"`
	BackgroundColor  string `json:"backgroundColor"`
}

type RecordingFileConfig struct {
	AvFileType []string `json:"avFileType"`
}

type RecordingConfig struct {
	MaxIdleTime        int               `json:"maxIdleTime"`
	StreamTypes        int               `json:"streamTypes"`
	AudioProfile       int               `json:"audioProfile"`
	ChannelType        int               `json:"channelType"`
	VideoStreamType    int               `json:"videoStreamType"`
	TranscodingConfig  TranscodingConfig `json:"transcodingConfig"`
	SubscribeVideoUids []string          `json:"subscribeVideoUids"`
	SubscribeAudioUids []string          `json:"subscribeAudioUids"`
	SubscribeUIDGroup  int               `json:"subscribeUidGroup"`
}

type StorageConfig struct {
	AccessKey      string   `json:"accessKey"`
	Region         int      `json:"region"`
	Bucket         string   `json:"bucket"`
	SecretKey      string   `json:"secretKey"`
	Vendor         int      `json:"vendor"`
	FileNamePrefix []string `json:"fileNamePrefix"`
}

//录制模式，支持以下几种录制模式：
//
//单流模式 individual：分开录制频道内每个 UID 的音频流和视频流，每个 UID 均有其对应的音频文件和视频文件。
//合流模式 mix ：（默认模式）频道内所有 UID 的音视频混合录制为一个音视频文件。
//页面录制模式 web：将指定网页的页面内容和音频混合录制为一个音视频文件
const (
	RECORDING_MODE_INDIVIDUAL = "individual"
	RECORDING_MODE_MIX        = "mix"
	RECORDING_MODE_WEB        = "web"
)

type StartResponse struct {
	Sid        string `json:"sid"`
	ResourceID string `json:"resourceId"`
}

type StopResponse struct {
	ResourceID     string             `json:"resourceId"`
	Sid            string             `json:"sid"`
	ServerResponse StopServerResponse `json:"serverResponse"`
}

type StopServerResponse struct {
	FileListMode    string          `json:"fileListMode"`
	UploadingStatus string          `json:"uploadingStatus"`
	FileList        json.RawMessage `json:"fileList"`
}

type FileList struct {
	Filename       string `json:"filename"`
	TrackType      string `json:"trackType"`
	UID            string `json:"uid"`
	MixedAllUser   bool   `json:"mixedAllUser"`
	IsPlayable     bool   `json:"isPlayable"`
	SliceStartTime int64  `json:"sliceStartTime"`
}

type QueryServerResponse struct {
	FileListMode    string     `json:"fileListMode"`
	FileList        []FileList `json:"fileList"`
	UploadingStatus string     `json:"uploadingStatus"`
}

type QueryResponse struct {
	ResourceID     string              `json:"resourceId"`
	Sid            string              `json:"sid"`
	ServerResponse QueryServerResponse `json:"serverResponse"`
}

type UpdateRequest struct {
	Cname         string              `json:"cname"`
	UID           string              `json:"uid"`
	ClientRequest UpdateClientRequest `json:"clientRequest"`
}

type AudioUIDList struct {
	SubscribeAudioUids []string `json:"subscribeAudioUids"`
}
type VideoUIDList struct {
	UnSubscribeVideoUids []string `json:"unSubscribeVideoUids"`
}
type StreamSubscribe struct {
	AudioUIDList AudioUIDList `json:"audioUidList"`
	VideoUIDList VideoUIDList `json:"videoUidList"`
}
type UpdateClientRequest struct {
	StreamSubscribe StreamSubscribe `json:"streamSubscribe"`
}

type UpdateResponse struct {
	Sid string `json:"sid"`
}

type UpdateLayOutRequest struct {
	Cname         string                    `json:"cname"`
	UID           string                    `json:"uid"`
	ClientRequest UpdateLayOutClientRequest `json:"clientRequest"`
}

type UpdateLayOutClientRequest struct {
	StreamSubscribe StreamSubscribe `json:"streamSubscribe"`
}
