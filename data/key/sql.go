/*
 * Copyright (c) 2018. Abstrium SAS <team (at) pydio.com>
 * This file is part of Pydio Cells.
 *
 * Pydio Cells is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Pydio Cells is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Pydio Cells.  If not, see <http://www.gnu.org/licenses/>.
 *
 * The latest code can be found at <https://pydio.com>.
 */

package key

import (
	"errors"
	"sync/atomic"

	sqldb "database/sql"
	"github.com/gobuffalo/packr"
	"github.com/pydio/cells/common"
	"github.com/pydio/cells/common/proto/encryption"
	"github.com/pydio/cells/common/sql"
	"github.com/rubenv/sql-migrate"
)

var (
	queries = map[string]interface{}{
		"node_select":                `SELECT * FROM enc_nodes WHERE node_id=?;`,
		"node_insert":                `INSERT INTO enc_nodes VALUES (?, ?);`,
		"node_update":                `UPDATE enc_nodes SET legacy=? WHERE node_id=?;`,
		"node_delete":                `DELETE FROM enc_nodes WHERE node_id=?;`,
		"node_key_insert":            `INSERT INTO enc_node_keys (node_id,owner_id,user_id,key_data) VALUES (?,?,?,?)`,
		"node_key_select":            `SELECT * FROM enc_node_keys WHERE node_id=? AND user_id=?;`,
		"node_key_delete":            `DELETE FROM enc_node_keys WHERE node_id=? AND user_id=?;`,
		"node_shared_key_delete":     `DELETE FROM enc_node_keys WHERE user_id<>owner_id AND node_id=? AND owner_id=? AND user_id=?`,
		"node_shared_key_delete_all": `DELETE FROM enc_node_keys WHERE  user_id<>owner_id AND node_id=? AND owner_id=?`,
		"node_block_insert":          `INSERT INTO enc_node_blocks (?, ?, ?, ?, ?, ?) ;`,
		"node_block_select":          `SELECT * FROM enc_node_blocks WHERE node_id=? order by part_id, block_position;`,
		"node_block_delete":          `DELETE FROM enc_node_blocks WHERE node_id=?;`,
		"node_block_part_delete":     `DELETE FROM enc_node_blocks WHERE node_id=? and part_id=?;`,
	}
	mu atomic.Value
)

type sqlimpl struct {
	sql.DAO
}

// Init handler for the SQL DAO
func (s *sqlimpl) Init(options common.ConfigValues) error {
	// super
	s.DAO.Init(options)

	// Doing the database migrations
	migrations := &sql.PackrMigrationSource{
		Box:         packr.NewBox("../../data/key/migrations"),
		Dir:         s.Driver(),
		TablePrefix: s.Prefix(),
	}

	_, err := sql.ExecMigration(s.DB(), s.Driver(), migrations, migrate.Up, "data_key_")
	if err != nil {
		return err
	}

	// Preparing the db statements
	if options.Bool("prepare", true) {
		for key, query := range queries {
			if err := s.Prepare(key, query); err != nil {
				return err
			}
		}
	}
	return nil
}

func (h *sqlimpl) ListEncryptedBlockInfo(nodeUuid string) (QueryResultCursor, error) {
	stmt := h.GetStmt("node_block_select")
	if stmt == nil {
		return nil, errors.New("internal error: 'node_block_select' statement not found")
	}

	rows, err := stmt.Query(nodeUuid)
	if err != nil {
		return nil, err
	}
	return NewDBCursor(rows, scanBlock), nil
}

func (h *sqlimpl) SaveEncryptedBlockInfo(nodeUuid string, b *encryption.Block) error {
	stmt := h.GetStmt("node_block_insert")
	if stmt == nil {
		return errors.New("internal error: 'node_block_insert' statement not found")
	}
	_, err := stmt.Exec(nodeUuid, b.PartId, b.Position, b.HeaderSize, b.Nonce, b.OwnerId)
	return err
}

func (h *sqlimpl) ClearNodeEncryptedBlockInfo(nodeUuid string) error {
	stmt := h.GetStmt("node_block_delete")
	if stmt == nil {
		return errors.New("internal error: 'node_block_delete' statement not found")
	}
	_, err := stmt.Exec(nodeUuid)
	return err
}

func (h *sqlimpl) SaveNode(node *encryption.Node) error {
	stmt := h.GetStmt("node_insert")
	if stmt == nil {
		return errors.New("internal error: 'node_insert' statement not found")
	}
	var intLegacy int
	if node.Legacy {
		intLegacy = 1
	}
	_, err := stmt.Exec(node.NodeId, intLegacy)
	return err
}

func (h *sqlimpl) DeleteNode(nodeUuid string) error {
	stmt := h.GetStmt("node_delete")
	if stmt == nil {
		return errors.New("internal error: 'node_delete' statement not found")
	}
	_, err := stmt.Exec(nodeUuid)
	return err
}

func (h *sqlimpl) SaveNodeKey(key *encryption.NodeKey) error {
	stmt := h.GetStmt("node_key_insert")
	if stmt == nil {
		return errors.New("internal error: 'node_key_insert' statement not found")
	}

	_, err := stmt.Exec(key.NodeId, key.OwnerId, key.UserId, key.KeyData)
	return err
}

func (h *sqlimpl) GetNodeKey(nodeUuid string, user string) (*encryption.NodeKey, error) {
	stmt := h.GetStmt("node_key_select")
	if stmt == nil {
		return nil, errors.New("internal error: 'node_key_select' statement not found")
	}

	rows, err := stmt.Query(nodeUuid, user)
	if err != nil {
		return nil, err
	}

	c := NewDBCursor(rows, scanNodeKey)

	if !c.HasNext() {
		_ = c.Close()
		return nil, errors.New("not found")
	}

	k, err := c.Next()
	if err != nil {
		_ = c.Close()
		return nil, err
	}

	return k.(*encryption.NodeKey), c.Close()
}

func (h *sqlimpl) DeleteNodeKey(key *encryption.NodeKey) error {
	stmt := h.GetStmt("node_key_delete")
	if stmt == nil {
		return errors.New("internal error: 'node_delete' statement not found")
	}
	_, err := stmt.Exec(key.NodeId, key.UserId)
	return err
}

// dbRowScanner
type dbRowScanner func(rows *sqldb.Rows) (interface{}, error)

// DBCursor
type DBCursor struct {
	err  error
	scan dbRowScanner
	rows *sqldb.Rows
}

func NewDBCursor(rows *sqldb.Rows, scanner dbRowScanner) QueryResultCursor {
	return &DBCursor{
		scan: scanner,
		rows: rows,
	}
}

func (c *DBCursor) Close() error {
	return c.rows.Close()
}

func (c *DBCursor) HasNext() bool {
	return c.rows.Next()
}

func (c *DBCursor) Next() (interface{}, error) {
	return c.scan(c.rows)
}

// scanBlock
func scanBlock(rows *sqldb.Rows) (interface{}, error) {
	b := new(encryption.Block)
	var nodeId string
	err := rows.Scan(&nodeId, &b.Position, &b.BlockSize, &b.HeaderSize, &b.Nonce, &b.OwnerId)
	return b, err
}

// scanNode
func scanNode(rows *sqldb.Rows) (interface{}, error) {
	n := new(encryption.Node)
	var legacy int

	err := rows.Scan(&n.NodeId, &legacy)
	n.Legacy = legacy == 1

	return n, err
}

// scanNodeKey
func scanNodeKey(rows *sqldb.Rows) (interface{}, error) {
	k := new(encryption.NodeKey)
	err := rows.Scan(&k.NodeId, &k.OwnerId, &k.UserId, &k.KeyData)
	return k, err
}
