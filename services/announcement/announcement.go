package announcement

import (
	"bunker-web/models"
	"database/sql"
	"time"
)

func Create(user *models.User, title, content string, isPinned bool) error {
	model := &models.Announcement{
		Title:    title,
		Content:  content,
		AuthorID: &user.ID,
	}
	if isPinned {
		model.PinnedAt = sql.NullTime{Time: time.Now(), Valid: true}
	}
	return models.DBCreate(model)
}

func GetTotal() (int64, error) {
	var total int64
	result := models.DB.Model(&models.Announcement{}).Count(&total)
	return total, result.Error
}

func QueryByPage(pageNum, pageSize int) ([]*models.Announcement, error) {
	var announcements []*models.Announcement
	result := models.DB.Preload("Author").
		Order("pinned_at desc").
		Order("created_at desc").
		Offset((pageNum - 1) * pageSize).
		Limit(pageSize).
		Find(&announcements)
	if result.Error != nil {
		return nil, result.Error
	}
	// Set author name and create time
	for _, announcement := range announcements {
		if announcement.Author != nil {
			announcement.AuthorName = announcement.Author.Username
		}
	}
	return announcements, nil
}

func QueryByID(id uint) (*models.Announcement, error) {
	var announcement models.Announcement
	result := models.DB.Where("id = ?", id).First(&announcement)
	return &announcement, result.Error
}

func DeleteByID(id uint) error {
	return models.DB.Delete(&models.Announcement{}, id).Error
}
