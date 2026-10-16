package listRepository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	listModel "sl-api/api/model/list"
)

type ListRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *ListRepository {
	return &ListRepository{
		db: db,
	}
}

func (r *ListRepository) Create(list *listModel.List) (*listModel.List, error) {
	if err := r.db.Create(list).Error; err != nil {
		return nil, err
	}

	return list, nil
}

func (r *ListRepository) Read(id uuid.UUID) (*listModel.List, error) {
	list := &listModel.List{}
	if err := r.db.Where("id = ?", id).First(&list).Error; err != nil {
		return nil, err
	}

	return list, nil
}

func (r *ListRepository) ListByGroup(groupID uuid.UUID) (listModel.Lists, error) {
	lists := make([]*listModel.List, 0)

	if err := r.db.Where("group_id = ?", groupID).Find(&lists).Error; err != nil {
		return nil, err
	}

	return lists, nil
}

func (r *ListRepository) Update(list *listModel.List) (int64, error) {
	result := r.db.Model(&listModel.List{}).
		Select("Name", "Status", "UpdatedAt").
		Where("id = ?", list.ID).
		Updates(list)

	return result.RowsAffected, result.Error
}

// UpdateStatusForMember sets a list's status, but only when the caller is a
// member of the list's group. It returns the updated row via RETURNING; updated
// is false when no row matched (not a member or no such list alike), so the
// endpoint never reveals a list the caller cannot see.
func (r *ListRepository) UpdateStatusForMember(listID, userID uuid.UUID, status listModel.Status) (*listModel.List, bool, error) {
	query := `
		UPDATE lists
		SET status = ?, updated_at = now()
		WHERE id = ? AND deleted_at IS NULL AND group_id IN (
			SELECT gm.group_id
			FROM group_members gm
			WHERE gm.user_id = ? AND gm.deleted_at IS NULL
		)
		RETURNING id, group_id, created_by, name, status
	`

	var list listModel.List
	result := r.db.Raw(query, string(status), listID, userID).Scan(&list)
	if result.Error != nil {
		return nil, false, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, false, nil
	}

	return &list, true, nil
}

// AddItemData carries the two independent existence lookups AddItem needs,
// fetched in a single round-trip. A field being Valid=false means that row was
// not found (NULL in the query), which lets the service tell a missing list
// apart from a missing product. ListGroupID also hands back the list's group so
// the membership check needs no extra read.
type AddItemData struct {
	ListGroupID uuid.NullUUID
	ProductID   uuid.NullUUID
}

// GetAddItemData looks up the list's group and the product's existence in one
// query. Scalar subqueries (rather than a join) are used so that a missing list
// and a missing product report independently instead of collapsing to zero rows.
func (r *ListRepository) GetAddItemData(listID, productID uuid.UUID) (*AddItemData, error) {
	var data AddItemData

	query := `
		SELECT
			(SELECT group_id FROM lists WHERE id = ? AND deleted_at IS NULL) AS list_group_id,
			(SELECT id FROM products WHERE id = ? AND deleted_at IS NULL) AS product_id
	`

	if err := r.db.Raw(query, listID, productID).Scan(&data).Error; err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *ListRepository) Delete(id uuid.UUID) (int64, error) {
	result := r.db.Where("id = ?", id).Delete(&listModel.List{})

	return result.RowsAffected, result.Error
}
