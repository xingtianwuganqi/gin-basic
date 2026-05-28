package models

import (
	"encoding/json"
	"time"
)

type SysDictType struct {
	ID        uint      `json:"id" form:"id" gorm:"primaryKey"`
	DictType  string    `json:"dict_type" form:"dict_type" gorm:"column:dict_type;size:100;not null;uniqueIndex;comment:字典类型"`
	DictName  string    `json:"dict_name" form:"dict_name" gorm:"column:dict_name;size:100;not null;comment:字典名称"`
	Status    int8      `json:"status" form:"status" gorm:"column:status;not null;default:1;comment:1启用 0禁用"`
	Remark    string    `json:"remark" form:"remark" gorm:"column:remark;size:255;default:''"`
	CreatedAt time.Time `json:"created_at" form:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" form:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (SysDictType) TableName() string {
	return "sys_dict_type"
}

type SysDictItem struct {
	ID        uint            `json:"id" form:"id" gorm:"primaryKey"`
	DictType  string          `json:"dict_type" form:"dict_type" gorm:"column:dict_type;size:100;not null;index:idx_dict_type;uniqueIndex:uk_type_value;comment:字典类型"`
	ItemLabel string          `json:"item_label" form:"item_label" gorm:"column:item_label;size:100;not null;comment:显示名称"`
	ItemValue string          `json:"item_value" form:"item_value" gorm:"column:item_value;size:100;not null;uniqueIndex:uk_type_value;comment:实际值"`
	Sort      int             `json:"sort" form:"sort" gorm:"column:sort;not null;default:0;comment:排序"`
	Status    int8            `json:"status" form:"status" gorm:"column:status;not null;default:1;comment:1启用 0禁用"`
	IsDefault int8            `json:"is_default" form:"is_default" gorm:"column:is_default;not null;default:0;comment:是否默认"`
	Extra     json.RawMessage `json:"extra,omitempty" gorm:"column:extra;type:json;comment:扩展字段"`
	Remark    string          `json:"remark" form:"remark" gorm:"column:remark;size:255;default:''"`
	CreatedAt time.Time       `json:"created_at" form:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time       `json:"updated_at" form:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (SysDictItem) TableName() string {
	return "sys_dict_item"
}

type CreateSysDictTypeRequest struct {
	DictType string `json:"dict_type" form:"dict_type" binding:"required"`
	DictName string `json:"dict_name" form:"dict_name" binding:"required"`
	Remark   string `json:"remark" form:"remark"`
}

type UpdateSysDictTypeRequest struct {
	DictType string `json:"dict_type" form:"dict_type" binding:"required"`
	DictName string `json:"dict_name" form:"dict_name" binding:"required"`
	Remark   string `json:"remark" form:"remark"`
}

type CreateSysDictItemRequest struct {
	DictType  string          `json:"dict_type" form:"dict_type" binding:"required"`
	ItemLabel string          `json:"item_label" form:"item_label" binding:"required"`
	ItemValue string          `json:"item_value" form:"item_value" binding:"required"`
	Sort      int             `json:"sort" form:"sort"`
	Status    *int8           `json:"status" form:"status"`
	IsDefault *int8           `json:"is_default" form:"is_default"`
	Extra     json.RawMessage `json:"extra"`
	Remark    string          `json:"remark" form:"remark"`
}

type UpdateSysDictItemRequest struct {
	ID        uint            `json:"id" form:"id" binding:"required"`
	DictType  string          `json:"dict_type" form:"dict_type" binding:"required"`
	ItemLabel string          `json:"item_label" form:"item_label" binding:"required"`
	ItemValue string          `json:"item_value" form:"item_value" binding:"required"`
	Sort      int             `json:"sort" form:"sort"`
	Status    *int8           `json:"status" form:"status"`
	IsDefault *int8           `json:"is_default" form:"is_default"`
	Extra     json.RawMessage `json:"extra"`
	Remark    string          `json:"remark" form:"remark"`
}

type SysDictTypeListResponse struct {
	Total    int64         `json:"total"`
	PageNum  int           `json:"pageNum"`
	PageSize int           `json:"pageSize"`
	List     []SysDictType `json:"list"`
}

type SysDictItemListResponse struct {
	Total    int64         `json:"total"`
	PageNum  int           `json:"pageNum"`
	PageSize int           `json:"pageSize"`
	List     []SysDictItem `json:"list"`
}
