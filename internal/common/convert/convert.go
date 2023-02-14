package convert

import (
	"Open_IM/internal/common/check"
	"Open_IM/pkg/common/db/table/relation"
	discoveryRegistry "Open_IM/pkg/discoveryregistry"
	sdk "Open_IM/pkg/proto/sdkws"
	"context"
	utils "github.com/OpenIMSDK/open_utils"
	"time"
)

type DBFriend struct {
	*relation.FriendModel
	userCheck *check.UserCheck
}

func NewDBFriend(friend *relation.FriendModel, zk discoveryRegistry.SvcDiscoveryRegistry) *DBFriend {
	return &DBFriend{FriendModel: friend, userCheck: check.NewUserCheck(zk)}
}

type PBFriend struct {
	*sdk.FriendInfo
}

func NewPBFriend(friendInfo *sdk.FriendInfo) *PBFriend {
	return &PBFriend{FriendInfo: friendInfo}
}

func (*PBFriend) PB2DB(friends []*sdk.FriendInfo) (DBFriends []*relation.FriendModel, err error) {
	for _, v := range friends {
		u, err := NewPBFriend(v).Convert()
		if err != nil {
			return nil, err
		}
		DBFriends = append(DBFriends, u)
	}
	return
}

func (*DBFriend) DB2PB(ctx context.Context, zk discoveryRegistry.SvcDiscoveryRegistry, friends []*relation.FriendModel) (PBFriends []*sdk.FriendInfo, err error) {
	for _, v := range friends {
		u, err := NewDBFriend(v, zk).Convert(ctx)
		if err != nil {
			return nil, err
		}
		PBFriends = append(PBFriends, u)
	}
	return
}

func (db *DBFriend) Convert(ctx context.Context) (*sdk.FriendInfo, error) {
	pbFriend := &sdk.FriendInfo{FriendUser: &sdk.UserInfo{}}
	utils.CopyStructFields(pbFriend, db)
	user, err := db.userCheck.GetUsersInfo(ctx, db.FriendUserID)
	if err != nil {
		return nil, err
	}
	utils.CopyStructFields(pbFriend.FriendUser, user)
	pbFriend.CreateTime = db.CreateTime.Unix()

	pbFriend.FriendUser.CreateTime = db.CreateTime.Unix()
	return pbFriend, nil
}

func (pb *PBFriend) Convert() (*relation.FriendModel, error) {
	dbFriend := &relation.FriendModel{}
	utils.CopyStructFields(dbFriend, pb)
	dbFriend.FriendUserID = pb.FriendUser.UserID
	dbFriend.CreateTime = utils.UnixSecondToTime(pb.CreateTime)
	return dbFriend, nil
}

type DBFriendRequest struct {
	*relation.FriendRequestModel
	userCheck *check.UserCheck
}

func NewDBFriendRequest(friendRequest *relation.FriendRequestModel, zk discoveryRegistry.SvcDiscoveryRegistry) *DBFriendRequest {
	return &DBFriendRequest{FriendRequestModel: friendRequest, userCheck: check.NewUserCheck(zk)}
}

type PBFriendRequest struct {
	*sdk.FriendRequest
}

func NewPBFriendRequest(friendRequest *sdk.FriendRequest) *PBFriendRequest {
	return &PBFriendRequest{FriendRequest: friendRequest}
}

func (*PBFriendRequest) PB2DB(friendRequests []*sdk.FriendRequest) (DBFriendRequests []*relation.FriendRequestModel, err error) {
	for _, v := range friendRequests {
		u, err := NewPBFriendRequest(v).Convert()
		if err != nil {
			return nil, err
		}
		DBFriendRequests = append(DBFriendRequests, u)
	}
	return
}

func (*DBFriendRequest) DB2PB(ctx context.Context, zk discoveryRegistry.SvcDiscoveryRegistry, friendRequests []*relation.FriendRequestModel) (PBFriendRequests []*sdk.FriendRequest, err error) {
	for _, v := range friendRequests {
		u, err := NewDBFriendRequest(v, zk).Convert(ctx)
		if err != nil {
			return nil, err
		}
		PBFriendRequests = append(PBFriendRequests, u)
	}
	return
}

func (pb *PBFriendRequest) Convert() (*relation.FriendRequestModel, error) {
	dbFriendRequest := &relation.FriendRequestModel{}
	utils.CopyStructFields(dbFriendRequest, pb)
	dbFriendRequest.CreateTime = utils.UnixSecondToTime(int64(pb.CreateTime))
	dbFriendRequest.HandleTime = utils.UnixSecondToTime(int64(pb.HandleTime))
	return dbFriendRequest, nil
}
func (db *DBFriendRequest) Convert(ctx context.Context) (*sdk.FriendRequest, error) {
	pbFriendRequest := &sdk.FriendRequest{}
	utils.CopyStructFields(pbFriendRequest, db)
	user, err := db.userCheck.GetUsersInfo(ctx, db.FromUserID)
	if err != nil {
		return nil, err
	}
	pbFriendRequest.FromNickname = user.Nickname
	pbFriendRequest.FromFaceURL = user.FaceURL
	pbFriendRequest.FromGender = user.Gender
	user, err = db.userCheck.GetUsersInfo(ctx, db.ToUserID)
	if err != nil {
		return nil, err
	}
	pbFriendRequest.ToNickname = user.Nickname
	pbFriendRequest.ToFaceURL = user.FaceURL
	pbFriendRequest.ToGender = user.Gender
	pbFriendRequest.CreateTime = db.CreateTime.Unix()
	pbFriendRequest.HandleTime = db.HandleTime.Unix()
	return pbFriendRequest, nil
}

type DBBlack struct {
	*relation.BlackModel
	userCheck *check.UserCheck
}

func (*PBBlack) PB2DB(blacks []*sdk.BlackInfo) (DBBlacks []*relation.BlackModel, err error) {
	for _, v := range blacks {
		u, err := NewPBBlack(v).Convert()
		if err != nil {
			return nil, err
		}
		DBBlacks = append(DBBlacks, u)
	}
	return
}

func (*DBBlack) DB2PB(ctx context.Context, zk discoveryRegistry.SvcDiscoveryRegistry, blacks []*relation.BlackModel) (PBBlacks []*sdk.BlackInfo, err error) {
	for _, v := range blacks {
		u, err := NewDBBlack(v, zk).Convert(ctx)
		if err != nil {
			return nil, err
		}
		PBBlacks = append(PBBlacks, u)
	}
	return
}

func NewDBBlack(black *relation.BlackModel, zk discoveryRegistry.SvcDiscoveryRegistry) *DBBlack {
	return &DBBlack{BlackModel: black, userCheck: check.NewUserCheck(zk)}
}

type PBBlack struct {
	*sdk.BlackInfo
}

func NewPBBlack(blackInfo *sdk.BlackInfo) *PBBlack {
	return &PBBlack{BlackInfo: blackInfo}
}

func (pb *PBBlack) Convert() (*relation.BlackModel, error) {
	dbBlack := &relation.BlackModel{}
	dbBlack.BlockUserID = pb.BlackUserInfo.UserID
	dbBlack.CreateTime = utils.UnixSecondToTime(int64(pb.CreateTime))
	return dbBlack, nil
}
func (db *DBBlack) Convert(ctx context.Context) (*sdk.BlackInfo, error) {
	pbBlack := &sdk.BlackInfo{}
	utils.CopyStructFields(pbBlack, db)
	pbBlack.CreateTime = db.CreateTime.Unix()
	user, err := db.userCheck.GetUsersInfo(ctx, db.BlockUserID)
	if err != nil {
		return nil, err
	}
	utils.CopyStructFields(pbBlack.BlackUserInfo, user)
	return pbBlack, nil
}

type DBGroup struct {
	*relation.GroupModel
	zk         discoveryRegistry.SvcDiscoveryRegistry
	groupCheck *check.GroupChecker
}

func (*PBGroup) PB2DB(groups []*sdk.GroupInfo) (DBGroups []*relation.GroupModel, err error) {
	for _, v := range groups {
		u, err := NewPBGroup(v).Convert()
		if err != nil {
			return nil, err
		}
		DBGroups = append(DBGroups, u)
	}
	return
}

func (db *DBGroup) DB2PB(ctx context.Context, zk discoveryRegistry.SvcDiscoveryRegistry, groups []*relation.GroupModel) (PBGroups []*sdk.GroupInfo, err error) {
	for _, v := range groups {
		u, err := NewDBGroup(v, zk).Convert(ctx)
		if err != nil {
			return nil, err
		}
		PBGroups = append(PBGroups, u)
	}
	return
}

func NewDBGroup(groupModel *relation.GroupModel, zk discoveryRegistry.SvcDiscoveryRegistry) *DBGroup {
	return &DBGroup{GroupModel: groupModel, groupCheck: check.NewGroupChecker(zk)}
}

type PBGroup struct {
	*sdk.GroupInfo
}

func NewPBGroup(groupInfo *sdk.GroupInfo) *PBGroup {
	return &PBGroup{GroupInfo: groupInfo}
}

func (pb *PBGroup) Convert() (*relation.GroupModel, error) {
	dst := &relation.GroupModel{}
	err := utils.CopyStructFields(dst, pb)
	return dst, err
}
func (db *DBGroup) Convert(ctx context.Context) (*sdk.GroupInfo, error) {
	dst := &sdk.GroupInfo{}
	utils.CopyStructFields(dst, db)
	user, err := db.groupCheck.GetOwnerInfo(ctx, db.GroupID)
	if err != nil {
		return nil, err
	}
	dst.OwnerUserID = user.UserID

	g, err := db.groupCheck.GetGroupInfo(ctx, db.GroupID)
	if err != nil {
		return nil, err
	}
	dst.MemberCount = g.MemberCount
	dst.CreateTime = db.CreateTime.Unix()
	dst.NotificationUpdateTime = db.NotificationUpdateTime.Unix()
	if db.NotificationUpdateTime.Unix() < 0 {
		dst.NotificationUpdateTime = 0
	}
	return dst, nil
}

type DBGroupMember struct {
	*relation.GroupMemberModel
	userCheck *check.UserCheck
}

func (*PBGroupMember) PB2DB(groupMembers []*sdk.GroupMemberFullInfo) (DBGroupMembers []*relation.GroupMemberModel, err error) {
	for _, v := range groupMembers {
		u, err := NewPBGroupMember(v).Convert()
		if err != nil {
			return nil, err
		}
		DBGroupMembers = append(DBGroupMembers, u)
	}
	return
}

func (*DBGroupMember) DB2PB(ctx context.Context, groupMembers []*relation.GroupMemberModel) (PBGroupMembers []*sdk.GroupMemberFullInfo, err error) {
	for _, v := range groupMembers {
		u, err := NewDBGroupMember(v).Convert(ctx)
		if err != nil {
			return nil, err
		}
		PBGroupMembers = append(PBGroupMembers, u)
	}
	return
}

func NewDBGroupMember(groupMember *relation.GroupMemberModel) *DBGroupMember {
	return &DBGroupMember{GroupMemberModel: groupMember}
}

type PBGroupMember struct {
	*sdk.GroupMemberFullInfo
}

func NewPBGroupMember(groupMemberFullInfo *sdk.GroupMemberFullInfo) *PBGroupMember {
	return &PBGroupMember{GroupMemberFullInfo: groupMemberFullInfo}
}

func (pb *PBGroupMember) Convert() (*relation.GroupMemberModel, error) {
	dst := &relation.GroupMemberModel{}
	utils.CopyStructFields(dst, pb)
	dst.JoinTime = utils.UnixSecondToTime(int64(pb.JoinTime))
	dst.MuteEndTime = utils.UnixSecondToTime(int64(pb.MuteEndTime))
	return dst, nil
}
func (db *DBGroupMember) Convert(ctx context.Context) (*sdk.GroupMemberFullInfo, error) {
	dst := &sdk.GroupMemberFullInfo{}
	utils.CopyStructFields(dst, db)

	user, err := db.userCheck.GetUsersInfo(ctx, db.UserID)
	if err != nil {
		return nil, err
	}
	dst.AppMangerLevel = user.AppMangerLevel

	dst.JoinTime = db.JoinTime.Unix()
	if db.JoinTime.Unix() < 0 {
		dst.JoinTime = 0
	}
	dst.MuteEndTime = db.MuteEndTime.Unix()
	if dst.MuteEndTime < time.Now().Unix() {
		dst.MuteEndTime = 0
	}
	return dst, nil
}

type DBGroupRequest struct {
	*relation.GroupRequestModel
}

func (*PBGroupRequest) PB2DB(groupRequests []*sdk.GroupRequest) (DBGroupRequests []*relation.GroupRequestModel, err error) {
	for _, v := range groupRequests {
		u, err := NewPBGroupRequest(v).Convert()
		if err != nil {
			return nil, err
		}
		DBGroupRequests = append(DBGroupRequests, u)
	}
	return
}

func (*DBGroupRequest) DB2PB(groupRequests []*relation.GroupRequestModel) (PBGroupRequests []*sdk.GroupRequest, err error) {
	for _, v := range groupRequests {
		u, err := NewDBGroupRequest(v).Convert()
		if err != nil {
			return nil, err
		}
		PBGroupRequests = append(PBGroupRequests, u)
	}
	return
}

func NewDBGroupRequest(groupRequest *relation.GroupRequestModel) *DBGroupRequest {
	return &DBGroupRequest{GroupRequestModel: groupRequest}
}

type PBGroupRequest struct {
	*sdk.GroupRequest
}

func NewPBGroupRequest(groupRequest *sdk.GroupRequest) *PBGroupRequest {
	return &PBGroupRequest{GroupRequest: groupRequest}
}

func (pb *PBGroupRequest) Convert() (*relation.GroupRequestModel, error) {
	dst := &relation.GroupRequestModel{}
	utils.CopyStructFields(dst, pb)
	dst.ReqTime = utils.UnixSecondToTime(int64(pb.ReqTime))
	dst.HandledTime = utils.UnixSecondToTime(int64(pb.HandleTime))
	return dst, nil
}
func (db *DBGroupRequest) Convert() (*sdk.GroupRequest, error) {
	dst := &sdk.GroupRequest{}
	utils.CopyStructFields(dst, db)
	dst.ReqTime = db.ReqTime.Unix()
	dst.HandleTime = db.HandledTime.Unix()
	return dst, nil
}

type DBUser struct {
	*relation.UserModel
}

func NewDBUser(user *relation.UserModel) *DBUser {
	return &DBUser{UserModel: user}
}

type PBUser struct {
	*sdk.UserInfo
}

func NewPBUser(userInfo *sdk.UserInfo) *PBUser {
	return &PBUser{UserInfo: userInfo}
}

func (*PBUser) PB2DB(users []*sdk.UserInfo) (DBUsers []*relation.UserModel, err error) {
	for _, v := range users {
		u, err := NewPBUser(v).Convert()
		if err != nil {
			return nil, err
		}
		DBUsers = append(DBUsers, u)
	}
	return
}

func (*DBUser) DB2PB(users []*relation.UserModel) (PBUsers []*sdk.UserInfo, err error) {
	for _, v := range users {
		u, err := NewDBUser(v).Convert()
		if err != nil {
			return nil, err
		}
		PBUsers = append(PBUsers, u)
	}
	return
}

func (pb *PBUser) Convert() (*relation.UserModel, error) {
	dst := &relation.UserModel{}
	utils.CopyStructFields(dst, pb)
	dst.Birth = utils.UnixSecondToTime(pb.Birthday)
	dst.CreateTime = utils.UnixSecondToTime(pb.CreateTime)
	return dst, nil
}

func (db *DBUser) Convert() (*sdk.UserInfo, error) {
	dst := &sdk.UserInfo{}
	utils.CopyStructFields(dst, db)
	dst.CreateTime = db.CreateTime.Unix()
	dst.Birthday = db.Birth.Unix()
	return dst, nil
}

func (db *DBUser) ConvertPublic() (*sdk.PublicUserInfo, error) {
	dst := &sdk.PublicUserInfo{}
	utils.CopyStructFields(dst, db)
	return dst, nil
}
