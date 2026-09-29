package listItemRepository

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	listItemModel "sl-api/api/model/list_item"
)

type ListItemRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *ListItemRepository {
	return &ListItemRepository{
		db: db,
	}
}

func (r *ListItemRepository) Create(item *listItemModel.ListItem) (*listItemModel.ListItem, error) {
	if err := r.db.Create(item).Error; err != nil {
		return nil, err
	}

	return item, nil
}

func (r *ListItemRepository) Read(id uuid.UUID) (*listItemModel.ListItem, error) {
	item := &listItemModel.ListItem{}
	if err := r.db.Where("id = ?", id).First(&item).Error; err != nil {
		return nil, err
	}

	return item, nil
}

func (r *ListItemRepository) ListByList(listID uuid.UUID) (listItemModel.ListItems, error) {
	items := make([]*listItemModel.ListItem, 0)

	if err := r.db.Where("list_id = ?", listID).Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

// FindByListAndProduct looks up the row for a (list, product) pair including
// soft-deleted ones. The unique constraint spans deleted rows, so callers adding
// a product need to see a previously removed row in order to revive it instead
// of inserting a duplicate.
func (r *ListItemRepository) FindByListAndProduct(listID, productID uuid.UUID) (*listItemModel.ListItem, error) {
	item := &listItemModel.ListItem{}
	if err := r.db.Unscoped().
		Where("list_id = ? AND product_id = ?", listID, productID).
		First(&item).Error; err != nil {
		return nil, err
	}

	return item, nil
}

// IncrementQuantity adds to an item's quantity in a single statement, so two
// concurrent adds of the same product cannot lose an increment.
func (r *ListItemRepository) IncrementQuantity(id uuid.UUID, by int) (int64, error) {
	result := r.db.Model(&listItemModel.ListItem{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"quantity":   gorm.Expr("quantity + ?", by),
			"updated_at": time.Now(),
		})

	return result.RowsAffected, result.Error
}

// Restore revives a soft-deleted row and resets the fields a fresh add would
// set. It writes through a map so empty values (an unset unit or note) are
// persisted rather than skipped as zero values.
func (r *ListItemRepository) Restore(item *listItemModel.ListItem) (int64, error) {
	result := r.db.Unscoped().
		Model(&listItemModel.ListItem{}).
		Where("id = ?", item.ID).
		Updates(map[string]any{
			"deleted_at": nil,
			"added_by":   item.AddedBy,
			"quantity":   item.Quantity,
			"unit":       item.Unit,
			"status":     item.Status,
			"note":       item.Note,
			"updated_at": time.Now(),
		})

	return result.RowsAffected, result.Error
}

func (r *ListItemRepository) Update(item *listItemModel.ListItem) (int64, error) {
	result := r.db.Model(&listItemModel.ListItem{}).
		Select("Quantity", "Unit", "Status", "Note", "UpdatedAt").
		Where("id = ?", item.ID).
		Updates(item)

	return result.RowsAffected, result.Error
}

func (r *ListItemRepository) Delete(id uuid.UUID) (int64, error) {
	result := r.db.Where("id = ?", id).Delete(&listItemModel.ListItem{})

	return result.RowsAffected, result.Error
}
