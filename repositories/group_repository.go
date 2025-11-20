package repositories

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"mms-backend/models"
)

// GroupRepository handles database operations for groups
type GroupRepository struct {
	db *gorm.DB
}

// NewGroupRepository creates a new group repository
func NewGroupRepository(db *gorm.DB) *GroupRepository {
	return &GroupRepository{db: db}
}

// Create creates a new group
func (r *GroupRepository) Create(group *models.Group) error {
	return r.db.Create(group).Error
}

// FindByID finds a group by ID
func (r *GroupRepository) FindByID(id uuid.UUID) (*models.Group, error) {
	var group models.Group
	err := r.db.Preload("Creator").Preload("Members.User").Where("id = ?", id).First(&group).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("group not found")
		}
		return nil, err
	}
	return &group, nil
}

// Update updates a group
func (r *GroupRepository) Update(group *models.Group) error {
	return r.db.Save(group).Error
}

// Delete deletes a group
func (r *GroupRepository) Delete(id uuid.UUID) error {
	// Delete all group members first
	if err := r.db.Where("group_id = ?", id).Delete(&models.GroupMember{}).Error; err != nil {
		return err
	}
	
	// Delete all group messages
	if err := r.db.Where("group_id = ?", id).Delete(&models.GroupMessage{}).Error; err != nil {
		return err
	}
	
	// Delete the group
	return r.db.Delete(&models.Group{}, id).Error
}

// List returns all groups (with pagination)
func (r *GroupRepository) List(limit, offset int) ([]models.Group, error) {
	var groups []models.Group
	err := r.db.Preload("Creator").Limit(limit).Offset(offset).Find(&groups).Error
	return groups, err
}

// GetUserGroups returns all groups a user belongs to
func (r *GroupRepository) GetUserGroups(userID uuid.UUID) ([]models.Group, error) {
	var groups []models.Group
	err := r.db.Preload("Creator").
		Joins("JOIN group_members ON group_members.group_id = groups.id").
		Where("group_members.user_id = ?", userID).
		Find(&groups).Error
	return groups, err
}

// GetPublicGroups returns all public groups
func (r *GroupRepository) GetPublicGroups(limit, offset int) ([]models.Group, error) {
	var groups []models.Group
	err := r.db.Preload("Creator").
		Where("type = ?", models.GroupTypePublic).
		Limit(limit).
		Offset(offset).
		Find(&groups).Error
	return groups, err
}

// AddMember adds a member to a group
func (r *GroupRepository) AddMember(member *models.GroupMember) error {
	return r.db.Create(member).Error
}

// RemoveMember removes a member from a group
func (r *GroupRepository) RemoveMember(groupID, userID uuid.UUID) error {
	return r.db.Where("group_id = ? AND user_id = ?", groupID, userID).
		Delete(&models.GroupMember{}).Error
}

// GetMember gets a specific group member
func (r *GroupRepository) GetMember(groupID, userID uuid.UUID) (*models.GroupMember, error) {
	var member models.GroupMember
	err := r.db.Where("group_id = ? AND user_id = ?", groupID, userID).First(&member).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("member not found")
		}
		return nil, err
	}
	return &member, nil
}

// UpdateMemberRole updates a member's role in a group
func (r *GroupRepository) UpdateMemberRole(groupID, userID uuid.UUID, role models.MemberRole) error {
	return r.db.Model(&models.GroupMember{}).
		Where("group_id = ? AND user_id = ?", groupID, userID).
		Update("role", role).Error
}

// IsMember checks if a user is a member of a group
func (r *GroupRepository) IsMember(groupID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&models.GroupMember{}).
		Where("group_id = ? AND user_id = ?", groupID, userID).
		Count(&count).Error
	return count > 0, err
}

// IsAdmin checks if a user is an admin of a group
func (r *GroupRepository) IsAdmin(groupID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&models.GroupMember{}).
		Where("group_id = ? AND user_id = ? AND role = ?", groupID, userID, models.MemberRoleAdmin).
		Count(&count).Error
	return count > 0, err
}

// IsCreator checks if a user is the creator of a group
func (r *GroupRepository) IsCreator(groupID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&models.Group{}).
		Where("id = ? AND created_by = ?", groupID, userID).
		Count(&count).Error
	return count > 0, err
}

// GetGroupMembers returns all members of a group
func (r *GroupRepository) GetGroupMembers(groupID uuid.UUID) ([]models.GroupMember, error) {
	var members []models.GroupMember
	err := r.db.Preload("User").Where("group_id = ?", groupID).Find(&members).Error
	return members, err
}

