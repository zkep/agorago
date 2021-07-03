package agorago

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

type Role uint16

const (
	Role_Attendee   = 1
	Role_Publisher  = 2
	Role_Subscriber = 3
	Role_Admin      = 4
)

var attendeePrivileges = map[uint16]uint32{
	KJoinChannel:        0,
	KPublishAudioStream: 0,
	KPublishVideoStream: 0,
	KPublishDataStream:  0,
}
var publisherPrivileges = map[uint16]uint32{
	KJoinChannel:              0,
	KPublishAudioStream:       0,
	KPublishVideoStream:       0,
	KPublishDataStream:        0,
	KPublishAudiocdn:          0,
	KPublishVideoCdn:          0,
	KInvitePublishAudioStream: 0,
	KInvitePublishVideoStream: 0,
	KInvitePublishDataStream:  0,
}

var subscriberPrivileges = map[uint16]uint32{
	KJoinChannel:               0,
	KRequestPublishAudioStream: 0,
	KRequestPublishVideoStream: 0,
	KRequestPublishDataStream:  0,
}

var adminPrivileges = map[uint16]uint32{
	KJoinChannel:         0,
	KPublishAudioStream:  0,
	KPublishVideoStream:  0,
	KPublishDataStream:   0,
	KAdministrateChannel: 0,
}

var RolePrivileges = map[uint16](map[uint16]uint32){
	Role_Attendee:   attendeePrivileges,
	Role_Publisher:  publisherPrivileges,
	Role_Subscriber: subscriberPrivileges,
	Role_Admin:      adminPrivileges,
}

type TokenBuilder struct {
	Token AccessToken
}

func CreateTokenBuilder(appID, appCertificate, channelName string, uid uint32) TokenBuilder {
	return TokenBuilder{CreateAccessToken(appID, appCertificate, channelName, uid)}
}

func (builder *TokenBuilder) InitPrivileges(role Role) {
	rolepri := uint16(role)
	for key, value := range RolePrivileges[rolepri] {
		builder.Token.Message[key] = value
	}
}

func (builder *TokenBuilder) InitTokenBuilder(originToken string) bool {
	return builder.Token.FromString(originToken)
}

func (builder *TokenBuilder) SetPrivilege(privilege Privileges, expireTimestamp uint32) {
	pri := uint16(privilege)
	builder.Token.Message[pri] = expireTimestamp
}

func (builder *TokenBuilder) RemovePrivilege(privilege Privileges) {
	pri := uint16(privilege)
	delete(builder.Token.Message, pri)
}

func (builder *TokenBuilder) BuildToken() (string, error) {
	return builder.Token.Build()
}

func GenerateSignalingToken(account, appID, appCertificate string, expiredTsInSeconds uint32) string {
	version := "1"
	expired := fmt.Sprint(expiredTsInSeconds)
	content := account + appID + appCertificate + expired
	hasher := md5.New()
	hasher.Write([]byte(content))
	md5sum := hex.EncodeToString(hasher.Sum(nil))
	result := fmt.Sprintf("%s:%s:%s:%s", version, appID, expired, md5sum)
	return result
}
