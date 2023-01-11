package im_mysql_model

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/utils"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

var (
	BlackListDB *gorm.DB
	UserDB      *gorm.DB
)

type BlackList struct {
	UserId           string    `gorm:"column:uid"`
	BeginDisableTime time.Time `gorm:"column:begin_disable_time"`
	EndDisableTime   time.Time `gorm:"column:end_disable_time"`
}

type User struct {
	UserID           string    `gorm:"column:user_id;primary_key;size:64"`
	Nickname         string    `gorm:"column:name;size:255"`
	FaceURL          string    `gorm:"column:face_url;size:255"`
	Gender           int32     `gorm:"column:gender"`
	PhoneNumber      string    `gorm:"column:phone_number;size:32"`
	Birth            time.Time `gorm:"column:birth"`
	Email            string    `gorm:"column:email;size:64"`
	Ex               string    `gorm:"column:ex;size:1024"`
	CreateTime       time.Time `gorm:"column:create_time;index:create_time"`
	AppMangerLevel   int32     `gorm:"column:app_manger_level"`
	GlobalRecvMsgOpt int32     `gorm:"column:global_recv_msg_opt"`

	status int32 `gorm:"column:status"`
}

func UserRegister(user User) error {
	user.CreateTime = time.Now()
	if user.AppMangerLevel == 0 {
		user.AppMangerLevel = constant.AppOrdinaryUsers
	}
	if user.Birth.Unix() < 0 {
		user.Birth = utils.UnixSecondToTime(0)
	}
	err := UserDB.Table("users").Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func GetAllUser() ([]User, error) {
	var userList []User
	err := UserDB.Table("users").Find(&userList).Error
	return userList, err
}

func TakeUserByUserID(userID string) (*User, error) {
	var user User
	err := UserDB.Table("users").Where("user_id=?", userID).Take(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByUserID(userID string) (*User, error) {
	var user User
	err := UserDB.Table("users").Where("user_id=?", userID).Take(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUsersByUserIDList(userIDList []string) ([]*User, error) {
	var userList []*User
	err := UserDB.Table("users").Where("user_id in (?)", userIDList).Find(&userList).Error
	return userList, err
}

func GetUserNameByUserID(userID string) (string, error) {
	var user User
	err := UserDB.Table("users").Select("name").Where("user_id=?", userID).First(&user).Error
	if err != nil {
		return "", err
	}
	return user.Nickname, nil
}

func UpdateUserInfo(user User) error {
	return UserDB.Where("user_id=?", user.UserID).Updates(&user).Error
}

func UpdateUserInfoByMap(user User, m map[string]interface{}) error {
	err := UserDB.Where("user_id=?", user.UserID).Updates(m).Error
	return err
}

func SelectAllUserID() ([]string, error) {
	var resultArr []string
	err := UserDB.Pluck("user_id", &resultArr).Error
	if err != nil {
		return nil, err
	}
	return resultArr, nil
}

func SelectSomeUserID(userIDList []string) ([]string, error) {
	var resultArr []string
	err := UserDB.Pluck("user_id", &resultArr).Error
	if err != nil {
		return nil, err
	}
	return resultArr, nil
}

func GetUsers(showNumber, pageNumber int32) ([]User, error) {
	var users []User
	err := UserDB.Limit(int(showNumber)).Offset(int(showNumber * (pageNumber - 1))).Find(&users).Error
	if err != nil {
		return users, err
	}
	return users, err
}

func AddUser(userID string, phoneNumber string, name string, email string, gender int32, faceURL string, birth string) error {
	_birth, err := utils.TimeStringToTime(birth)
	if err != nil {
		return err
	}
	user := User{
		UserID:      userID,
		Nickname:    name,
		FaceURL:     faceURL,
		Gender:      gender,
		PhoneNumber: phoneNumber,
		Birth:       _birth,
		Email:       email,
		Ex:          "",
		CreateTime:  time.Now(),
	}
	result := UserDB.Create(&user)
	return result.Error
}

func UserIsBlock(userId string) (bool, error) {
	var user BlackList
	rows := BlackListDB.Where("uid=?", userId).First(&user).RowsAffected
	if rows >= 1 {
		return user.EndDisableTime.After(time.Now()), nil
	}
	return false, nil
}

func UsersIsBlock(userIDList []string) (inBlockUserIDList []string, err error) {
	err = BlackListDB.Where("uid in (?) and end_disable_time > now()", userIDList).Pluck("uid", &inBlockUserIDList).Error
	return inBlockUserIDList, err
}

func BlockUser(userID, endDisableTime string) error {
	user, err := GetUserByUserID(userID)
	if err != nil || user.UserID == "" {
		return err
	}
	end, err := time.Parse("2006-01-02 15:04:05", endDisableTime)
	if err != nil {
		return err
	}
	if end.Before(time.Now()) {
		return errors.New("endDisableTime is before now")
	}
	var blockUser BlackList
	BlackListDB.Where("uid=?", userID).First(&blockUser)
	if blockUser.UserId != "" {
		BlackListDB.Where("uid=?", blockUser.UserId).Update("end_disable_time", end)
		return nil
	}
	blockUser = BlackList{
		UserId:           userID,
		BeginDisableTime: time.Now(),
		EndDisableTime:   end,
	}
	err = BlackListDB.Create(&blockUser).Error
	return err
}

func UnBlockUser(userID string) error {
	return BlackListDB.Where("uid=?", userID).Delete(&BlackList{}).Error
}

type BlockUserInfo struct {
	User             User
	BeginDisableTime time.Time
	EndDisableTime   time.Time
}

func GetBlockUserByID(userId string) (BlockUserInfo, error) {
	var blockUserInfo BlockUserInfo
	blockUser := BlackList{
		UserId: userId,
	}
	if err := BlackListDB.Table("black_lists").Where("uid=?", userId).Find(&blockUser).Error; err != nil {
		return blockUserInfo, err
	}
	user := User{
		UserID: blockUser.UserId,
	}
	if err := BlackListDB.Find(&user).Error; err != nil {
		return blockUserInfo, err
	}
	blockUserInfo.User.UserID = user.UserID
	blockUserInfo.User.FaceURL = user.FaceURL
	blockUserInfo.User.Nickname = user.Nickname
	blockUserInfo.User.Birth = user.Birth
	blockUserInfo.User.PhoneNumber = user.PhoneNumber
	blockUserInfo.User.Email = user.Email
	blockUserInfo.User.Gender = user.Gender
	blockUserInfo.BeginDisableTime = blockUser.BeginDisableTime
	blockUserInfo.EndDisableTime = blockUser.EndDisableTime
	return blockUserInfo, nil
}

func GetBlockUsers(showNumber, pageNumber int32) ([]BlockUserInfo, error) {
	var blockUserInfos []BlockUserInfo
	var blockUsers []BlackList
	if err := BlackListDB.Limit(int(showNumber)).Offset(int(showNumber * (pageNumber - 1))).Find(&blockUsers).Error; err != nil {
		return blockUserInfos, err
	}
	for _, blockUser := range blockUsers {
		var user User
		if err := UserDB.Table("users").Where("user_id=?", blockUser.UserId).First(&user).Error; err == nil {
			blockUserInfos = append(blockUserInfos, BlockUserInfo{
				User: User{
					UserID:      user.UserID,
					Nickname:    user.Nickname,
					FaceURL:     user.FaceURL,
					Birth:       user.Birth,
					PhoneNumber: user.PhoneNumber,
					Email:       user.Email,
					Gender:      user.Gender,
				},
				BeginDisableTime: blockUser.BeginDisableTime,
				EndDisableTime:   blockUser.EndDisableTime,
			})
		}
	}
	return blockUserInfos, nil
}

func GetUserByName(userName string, showNumber, pageNumber int32) ([]User, error) {
	var users []User
	err := UserDB.Where(" name like ?", fmt.Sprintf("%%%s%%", userName)).Limit(int(showNumber)).Offset(int(showNumber * (pageNumber - 1))).Find(&users).Error
	return users, err
}

func GetUsersByNameAndID(content string, showNumber, pageNumber int32) ([]User, int64, error) {
	var users []User
	var count int64
	db := UserDB.Where(" name like ? or user_id = ? ", fmt.Sprintf("%%%s%%", content), content)
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	err := db.Limit(int(showNumber)).Offset(int(showNumber * (pageNumber - 1))).Find(&users).Error
	return users, count, err
}

func GetUserIDsByEmailAndID(phoneNumber, email string) ([]string, error) {
	if phoneNumber == "" && email == "" {
		return nil, nil
	}
	db := UserDB
	if phoneNumber != "" {
		db = db.Where("phone_number = ? ", phoneNumber)
	}
	if email != "" {
		db = db.Where("email = ? ", email)
	}
	var userIDList []string
	err := db.Pluck("user_id", &userIDList).Error
	return userIDList, err
}

func GetUsersCount(userName string) (int32, error) {
	var count int64
	if err := UserDB.Where(" name like ? ", fmt.Sprintf("%%%s%%", userName)).Count(&count).Error; err != nil {
		return 0, err
	}
	return int32(count), nil
}

func GetBlockUsersNumCount() (int32, error) {
	var count int64
	if err := BlackListDB.Count(&count).Error; err != nil {
		return 0, err
	}
	return int32(count), nil
}
