package service

import (
	"errors"

	"gin-basic/db"
	"gin-basic/models"

	"gorm.io/gorm"
)

var (
	ErrSysDictTypeExists    = errors.New("sys dict type exists")
	ErrSysDictTypeNotFound  = errors.New("sys dict type not found")
	ErrSysDictTypeDisabled  = errors.New("sys dict type disabled")
	ErrSysDictItemExists    = errors.New("sys dict item exists")
	ErrSysDictItemNotFound  = errors.New("sys dict item not found")
	ErrSysDictItemAmbiguous = errors.New("sys dict item ambiguous")
)

func normalizePage(page models.Page) models.Page {
	if page.PageNum <= 0 {
		page.PageNum = 1
	}
	if page.PageSize <= 0 {
		page.PageSize = 10
	}
	return page
}

func dictTypeExists(dictType string) error {
	var count int64
	if err := db.DB.Model(&models.SysDictType{}).Where("dict_type = ?", dictType).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return ErrSysDictTypeNotFound
	}
	return nil
}

func getEnabledSysDictType(dictType string) (*models.SysDictType, error) {
	record := &models.SysDictType{}
	if err := db.DB.Where("dict_type = ?", dictType).First(record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSysDictTypeNotFound
		}
		return nil, err
	}
	if record.Status != 1 {
		return nil, ErrSysDictTypeDisabled
	}
	return record, nil
}

func CreateSysDictType(req models.CreateSysDictTypeRequest) (*models.SysDictType, error) {
	dictType := &models.SysDictType{}
	err := db.DB.Where("dict_type = ?", req.DictType).First(dictType).Error
	if err == nil {
		return nil, ErrSysDictTypeExists
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	dictType = &models.SysDictType{
		DictType: req.DictType,
		DictName: req.DictName,
		Status:   1,
		Remark:   req.Remark,
	}
	if err := db.DB.Create(dictType).Error; err != nil {
		return nil, err
	}
	return dictType, nil
}

func UpdateSysDictType(req models.UpdateSysDictTypeRequest) (*models.SysDictType, error) {
	dictType := &models.SysDictType{}
	if err := db.DB.Where("dict_type = ?", req.DictType).First(dictType).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSysDictTypeNotFound
		}
		return nil, err
	}

	dictType.DictName = req.DictName
	dictType.Remark = req.Remark
	if err := db.DB.Save(dictType).Error; err != nil {
		return nil, err
	}
	return dictType, nil
}

func ListSysDictTypes(page models.Page) (*models.SysDictTypeListResponse, error) {
	page = normalizePage(page)

	query := db.DB.Model(&models.SysDictType{})
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	list := make([]models.SysDictType, 0)
	if err := query.Order("created_at desc").
		Limit(page.PageSize).
		Offset((page.PageNum - 1) * page.PageSize).
		Find(&list).Error; err != nil {
		return nil, err
	}

	return &models.SysDictTypeListResponse{
		Total:    total,
		PageNum:  page.PageNum,
		PageSize: page.PageSize,
		List:     list,
	}, nil
}

func UpdateSysDictTypeStatus(dictType string, status int8) (*models.SysDictType, error) {
	record := &models.SysDictType{}
	if err := db.DB.Where("dict_type = ?", dictType).First(record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSysDictTypeNotFound
		}
		return nil, err
	}

	record.Status = status
	if err := db.DB.Save(record).Error; err != nil {
		return nil, err
	}
	return record, nil
}

func CreateSysDictItem(req models.CreateSysDictItemRequest) (*models.SysDictItem, error) {
	if err := dictTypeExists(req.DictType); err != nil {
		return nil, err
	}

	duplicate := &models.SysDictItem{}
	err := db.DB.Where("dict_type = ? AND item_value = ?", req.DictType, req.ItemValue).First(duplicate).Error
	if err == nil {
		return nil, ErrSysDictItemExists
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	item := &models.SysDictItem{
		DictType:  req.DictType,
		ItemLabel: req.ItemLabel,
		ItemValue: req.ItemValue,
		Sort:      req.Sort,
		Status:    1,
		IsDefault: 0,
		Extra:     req.Extra,
		Remark:    req.Remark,
	}
	if req.Status != nil {
		item.Status = *req.Status
	}
	if req.IsDefault != nil {
		item.IsDefault = *req.IsDefault
	}

	if err := saveSysDictItem(item); err != nil {
		return nil, err
	}
	return item, nil
}

func UpdateSysDictItem(req models.UpdateSysDictItemRequest) (*models.SysDictItem, error) {
	if err := dictTypeExists(req.DictType); err != nil {
		return nil, err
	}

	item := &models.SysDictItem{}
	if err := db.DB.Where("id = ?", req.ID).First(item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSysDictItemNotFound
		}
		return nil, err
	}

	var count int64
	if err := db.DB.Model(&models.SysDictItem{}).
		Where("dict_type = ? AND item_value = ? AND id <> ?", req.DictType, req.ItemValue, req.ID).
		Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, ErrSysDictItemExists
	}

	item.DictType = req.DictType
	item.ItemLabel = req.ItemLabel
	item.ItemValue = req.ItemValue
	item.Sort = req.Sort
	item.Extra = req.Extra
	item.Remark = req.Remark
	if req.Status != nil {
		item.Status = *req.Status
	}
	if req.IsDefault != nil {
		item.IsDefault = *req.IsDefault
	}

	if err := saveSysDictItem(item); err != nil {
		return nil, err
	}
	return item, nil
}

func saveSysDictItem(item *models.SysDictItem) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if item.IsDefault == 1 {
			if err := tx.Model(&models.SysDictItem{}).
				Where("dict_type = ? AND id <> ?", item.DictType, item.ID).
				Update("is_default", 0).Error; err != nil {
				return err
			}
		}
		return tx.Save(item).Error
	})
}

func ListSysDictItems(page models.Page, dictType string) (*models.SysDictItemListResponse, error) {
	page = normalizePage(page)

	query := db.DB.Model(&models.SysDictItem{})
	if dictType != "" {
		query = query.Where("dict_type = ?", dictType)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	list := make([]models.SysDictItem, 0)
	if err := query.Order("sort asc, id asc").
		Limit(page.PageSize).
		Offset((page.PageNum - 1) * page.PageSize).
		Find(&list).Error; err != nil {
		return nil, err
	}

	return &models.SysDictItemListResponse{
		Total:    total,
		PageNum:  page.PageNum,
		PageSize: page.PageSize,
		List:     list,
	}, nil
}

func UpdateSysDictItemStatus(itemValue string, status int8, dictType string) (*models.SysDictItem, error) {
	query := db.DB.Model(&models.SysDictItem{}).Where("item_value = ?", itemValue)
	if dictType != "" {
		query = query.Where("dict_type = ?", dictType)
	}

	list := make([]models.SysDictItem, 0)
	if err := query.Find(&list).Error; err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, ErrSysDictItemNotFound
	}
	if len(list) > 1 {
		return nil, ErrSysDictItemAmbiguous
	}

	item := &list[0]
	item.Status = status
	if err := db.DB.Save(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

func ListClientSysDictItems(dictType string) ([]models.SysDictItem, error) {
	if _, err := getEnabledSysDictType(dictType); err != nil {
		return nil, err
	}

	list := make([]models.SysDictItem, 0)
	if err := db.DB.Model(&models.SysDictItem{}).
		Where("dict_type = ? AND status = ?", dictType, 1).
		Order("sort asc, id asc").
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
