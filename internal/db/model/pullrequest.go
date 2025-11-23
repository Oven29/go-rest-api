package model

type PrStatus string

const (
	PrStatusOpen   PrStatus = "OPEN"
	PrStatusMerged PrStatus = "MERGED"
)

type PullRequest struct {
	ID       uint   `gorm:"primaryKey"`
	Title    string `gorm:"size:255;not null"`
	AuthorID uint   `gorm:"not null"`

	Author    User     `gorm:"foreignKey:AuthorID;constraint:OnDelete:RESTRICT"`
	Status    PrStatus `gorm:"type:pr_status;default:OPEN;not null"`
	Reviewers []User   `gorm:"many2many:pull_request_reviewer;foreignKey:ID;joinForeignKey:PrID;References:ID;joinReferences:ReviewerID"`
}

type PullRequestReviewer struct {
	PrID       uint `gorm:"primaryKey"`
	ReviewerID uint `gorm:"primaryKey;index:idx_reviewer_id"`

	PullRequest PullRequest `gorm:"foreignKey:PrID;constraint:OnDelete:CASCADE"`
	Reviewer    User        `gorm:"constraint:OnDelete:RESTRICT"`
}

func (PullRequestReviewer) TableName() string {
	return "pull_request_reviewer"
}
