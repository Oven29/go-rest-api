package model

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"size:255;not null"`
	IsActive bool   `gorm:"not null;default:true"`

	Teams                []Team        `gorm:"many2many:user_team"`
	AuthoredPullRequests []PullRequest `gorm:"foreignKey:AuthorID"`
	ReviewedPullRequests []PullRequest `gorm:"many2many:pull_request_reviewer;foreignKey:ID;joinForeignKey:ReviewerID;References:ID;joinReferences:PrID"`
}

type UserTeam struct {
	UserID uint `gorm:"primaryKey"`
	TeamID uint `gorm:"primaryKey"`

	User User
	Team Team
}

func (UserTeam) TableName() string {
	return "user_team"
}
