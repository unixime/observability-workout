package model

type Tenant struct {
	ID       int32  `json:"id" gorm:"primaryKey;autoIncrement"`
	Name     string `json:"name" gorm:"not null;uniqueIndex:idx_tenant_identity" binding:"required"`
	Lastname string `json:"lastname" gorm:"not null;uniqueIndex:idx_tenant_identity" binding:"required"`
	Age      int    `json:"age"`
}
