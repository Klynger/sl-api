package listItemRepository

import (
	"database/sql"
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

// FindByListAndProduct looks up the live row for a (list, product) pair, used to
// detect whether a product is already on a list. Returns gorm.ErrRecordNotFound
// when the pair is not present.
func (r *ListItemRepository) FindByListAndProduct(listID, productID uuid.UUID) (*listItemModel.ListItem, error) {
	item := &listItemModel.ListItem{}
	if err := r.db.
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
// ListForMember returns a list's items (each with its product name) only if the
// caller is a member of the list's group. The inner join on group_members
// enforces membership, and the LEFT JOIN on list_items keeps one row for an
// authorized-but-empty list, so the caller can tell "member, empty list" (rows
// returned, no item data) from "not a member / no such list" (zero rows).
// authorized is false exactly when zero rows come back.
func (r *ListItemRepository) ListForMember(listID, userID uuid.UUID) (listItemModel.ItemsWithProduct, bool, error) {
	type row struct {
		ID          uuid.NullUUID
		ProductID   uuid.NullUUID
		ProductName sql.NullString
		AddedBy     uuid.NullUUID
		Quantity    sql.NullInt64
		Unit        sql.NullString
		Status      sql.NullString
		Note        sql.NullString
	}

	query := `
		SELECT
			li.id, li.product_id, p.product_name AS product_name,
			li.added_by, li.quantity, li.unit, li.status, li.note
		FROM lists l
		JOIN group_members gm
			ON gm.group_id = l.group_id AND gm.user_id = ? AND gm.deleted_at IS NULL
		LEFT JOIN list_items li ON li.list_id = l.id
		LEFT JOIN products p ON p.id = li.product_id
		WHERE l.id = ? AND l.deleted_at IS NULL
	`

	var rows []row
	if err := r.db.Raw(query, userID, listID).Scan(&rows).Error; err != nil {
		return nil, false, err
	}

	if len(rows) == 0 {
		return nil, false, nil
	}

	items := make(listItemModel.ItemsWithProduct, 0, len(rows))
	for _, rw := range rows {
		if !rw.ID.Valid {
			continue // the placeholder row for an authorized but empty list
		}

		items = append(items, &listItemModel.ItemWithProduct{
			ID:          rw.ID.UUID,
			ProductID:   rw.ProductID.UUID,
			ProductName: rw.ProductName.String,
			AddedBy:     rw.AddedBy.UUID,
			Quantity:    int(rw.Quantity.Int64),
			Unit:        rw.Unit.String,
			Status:      listItemModel.Status(rw.Status.String),
			Note:        rw.Note.String,
		})
	}

	return items, true, nil
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
