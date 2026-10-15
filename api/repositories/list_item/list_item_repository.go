package listItemRepository

import (
	"database/sql"
	"strings"
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

// IncrementQuantityForMember adds to an item's quantity in a single atomic
// statement, but only when the item belongs to the given list and the caller is
// a member of that list's group. The membership subquery gates the update, and
// RETURNING yields the new quantity. updated is false when no row matched (not a
// member, no such list, or no such item alike), so the caller cannot tell them
// apart.
func (r *ListItemRepository) IncrementQuantityForMember(listID, itemID, userID uuid.UUID, by int) (int, bool, error) {
	query := `
		UPDATE list_items
		SET quantity = quantity + ?, updated_at = now()
		WHERE id = ? AND list_id = ? AND list_id IN (
			SELECT l.id
			FROM lists l
			JOIN group_members gm
				ON gm.group_id = l.group_id AND gm.user_id = ? AND gm.deleted_at IS NULL
			WHERE l.deleted_at IS NULL
		)
		RETURNING quantity
	`

	var newQuantity int
	result := r.db.Raw(query, by, itemID, listID, userID).Scan(&newQuantity)
	if result.Error != nil {
		return 0, false, result.Error
	}

	if result.RowsAffected == 0 {
		return 0, false, nil
	}

	return newQuantity, true, nil
}

// DecrementQuantityForMember lowers an item's quantity by `by`, but only when the
// caller is a member of the item's group and the result stays at least 1. One
// statement answers three questions so the caller can map them to the right
// status: `visible` is whether the member can see the item at all (false means a
// 404), and `decremented` is whether the update actually happened (false while
// visible means the floor of 1 would be breached, a 422). The floor lives in the
// UPDATE's own WHERE, so under READ COMMITTED a concurrent decrement re-checks it
// and cannot push the quantity below 1.
func (r *ListItemRepository) DecrementQuantityForMember(listID, itemID, userID uuid.UUID, by int) (newQuantity int, visible bool, decremented bool, err error) {
	query := `
		WITH visible AS (
			SELECT li.id
			FROM list_items li
			JOIN lists l ON l.id = li.list_id AND l.deleted_at IS NULL
			JOIN group_members gm
				ON gm.group_id = l.group_id AND gm.user_id = ? AND gm.deleted_at IS NULL
			WHERE li.id = ? AND li.list_id = ?
		),
		updated AS (
			UPDATE list_items
			SET quantity = quantity - ?, updated_at = now()
			WHERE id IN (SELECT id FROM visible) AND quantity - ? >= 1
			RETURNING quantity
		)
		SELECT
			EXISTS (SELECT 1 FROM visible) AS visible,
			(SELECT quantity FROM updated) AS new_quantity
	`

	var row struct {
		Visible     bool
		NewQuantity sql.NullInt64
	}
	if err := r.db.Raw(query, userID, itemID, listID, by, by).Scan(&row).Error; err != nil {
		return 0, false, false, err
	}

	return int(row.NewQuantity.Int64), row.Visible, row.NewQuantity.Valid, nil
}

// ItemPatch carries the mutable fields of a list item for a partial update. A
// nil field is left unchanged. Modelling the updatable columns as fixed struct
// fields (rather than a map of string keys) means a column name can never come
// from outside this package: the set is closed at compile time.
type ItemPatch struct {
	Quantity *int
	Unit     *string
	Note     *string
}

// UpdateForMember applies a patch to one item, but only when the item belongs to
// the list and the caller is a member of its group. The membership subquery
// gates the write; it returns the updated row via RETURNING, and updated is
// false when no row matched (not a member, no such list, or no such item alike).
// Each column maps to a hardcoded literal here, so nothing stringly-typed from a
// request can reach the SQL.
func (r *ListItemRepository) UpdateForMember(listID, itemID, userID uuid.UUID, patch ItemPatch) (*listItemModel.ListItem, bool, error) {
	setClauses := make([]string, 0, 4)
	args := make([]any, 0, 6)
	if patch.Quantity != nil {
		setClauses = append(setClauses, "quantity = ?")
		args = append(args, *patch.Quantity)
	}
	if patch.Unit != nil {
		setClauses = append(setClauses, "unit = ?")
		args = append(args, *patch.Unit)
	}
	if patch.Note != nil {
		setClauses = append(setClauses, "note = ?")
		args = append(args, *patch.Note)
	}
	setClauses = append(setClauses, "updated_at = now()")

	query := "UPDATE list_items SET " + strings.Join(setClauses, ", ") + `
		WHERE id = ? AND list_id = ? AND list_id IN (
			SELECT l.id
			FROM lists l
			JOIN group_members gm
				ON gm.group_id = l.group_id AND gm.user_id = ? AND gm.deleted_at IS NULL
			WHERE l.deleted_at IS NULL
		)
		RETURNING id, list_id, product_id, added_by, quantity, unit, status, note`
	args = append(args, itemID, listID, userID)

	var item listItemModel.ListItem
	result := r.db.Raw(query, args...).Scan(&item)
	if result.Error != nil {
		return nil, false, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, false, nil
	}

	return &item, true, nil
}

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

// DeleteForMember hard-deletes one item from a list, but only when the item
// belongs to that list and the caller is a member of the list's group. The
// membership subquery is what gates the delete; it reports whether a row was
// removed, and false covers "not a member", "no such list", and "no such item"
// alike, so the caller cannot tell them apart.
func (r *ListItemRepository) DeleteForMember(listID, itemID, userID uuid.UUID) (bool, error) {
	query := `
		DELETE FROM list_items
		WHERE id = ? AND list_id = ? AND list_id IN (
			SELECT l.id
			FROM lists l
			JOIN group_members gm
				ON gm.group_id = l.group_id AND gm.user_id = ? AND gm.deleted_at IS NULL
			WHERE l.deleted_at IS NULL
		)
	`

	result := r.db.Exec(query, itemID, listID, userID)
	if result.Error != nil {
		return false, result.Error
	}

	return result.RowsAffected > 0, nil
}
