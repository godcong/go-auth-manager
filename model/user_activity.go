package model

import (
	"github.com/go-xorm/xorm"
	"golang.org/x/xerrors"
)

// UserActivity ...
type UserActivity struct {
	Model        `xorm:"extends" json:",inline"`
	PropertyID   string `xorm:"notnull default('') comment(配置ID) property_id" json:"property_id"`
	ActivityID   string `xorm:"notnull unique(user_activity) default('') comment(活动ID) activity_id" json:"activity_id"`
	UserID       string `xorm:"notnull unique(user_activity) default('') comment(参加活动的用户ID) user_id" json:"user_id"`
	IsFavorite   bool   `xorm:"notnull default(false) comment(是否收藏) is_favorite " json:"is_favorite"`
	SpreadCode   string `xorm:"notnull unique default('') comment(参加活动的用户推广码) spread_code"  json:"spread_code"`
	IsVerified   bool   `xorm:"notnull default(false)  comment(校验通过) is_verified" json:"is_verified"`
	SpreadNumber int64  `xorm:"notnull default(0) comment(推广数) spread_number" json:"spread_number"`
}

// NewUserActivity ...
func NewUserActivity(id string) *UserActivity {
	return &UserActivity{
		Model: Model{
			ID: id,
		},
	}
}

// Get ...
func (obj *UserActivity) Get() (bool, error) {
	return Get(nil, obj)
}

// Update ...
func (obj *UserActivity) Update(cols ...string) (int64, error) {
	return Update(nil, obj.ID, obj)
}

// CodeSpread ...
func (obj *UserActivity) CodeSpread(session *xorm.Session) (*Spread, error) {
	var info struct {
		UserActivity UserActivity `xorm:"extends"`
		Spread       Spread       `xorm:"extends"`
	}
	b, e := MustSession(session).Table(obj).Join("left", info.Spread, "user_activity.user_id = spread.user_id").
		Where("user_activity.spread_code = ?", obj.SpreadCode).
		Get(&info)
	if e != nil {
		return nil, e
	}
	if !b {
		e = xerrors.New("spread not found")
		return nil, e
	}
	*obj = info.UserActivity

	return &info.Spread, nil
}

// Property ...
func (obj *UserActivity) Property(session *xorm.Session) (*Property, error) {
	var info struct {
		UserActivity UserActivity `xorm:"extends"`
		Property     Property     `xorm:"extends"`
	}
	b, e := MustSession(session).Table(obj).Join("left", info.Property, "user_activity.property_id = property.id").
		Where("user_activity.id = ?", obj.ID).
		Get(&info)
	if e != nil {
		return nil, e
	}
	if !b {
		e = xerrors.New("property not found")
		return nil, e
	}
	*obj = info.UserActivity
	return &info.Property, nil
}

// UserActivityActivity ...
type UserActivityActivity struct {
	UserActivity UserActivity `xorm:"extends"`
	Activity     Activity     `xorm:"extends"`
}

// Activities ...
func (obj *UserActivity) Activities(session *xorm.Session) ([]*UserActivityActivity, error) {
	var activities []*UserActivityActivity
	e := MustSession(session).Table(obj).Join("left", &Activity{}, "user_activity.activity_id = activity.id").
		Where("user_activity.user_id = ?", obj.UserID).
		Find(&activities)
	if e != nil {
		return nil, e
	}
	return activities, nil
}
