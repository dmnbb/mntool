package idmaker

import (
	"database/sql"
	"fmt"
	"sync"
	"time"
)

type BizId int32

func New(bizId BizId, db *sql.DB, step uint64) error {
	if idmakerMgr.get(bizId) != nil {
		return nil
	}

	if db == nil || step <= 0 {
		return fmt.Errorf("param error")
	}

	i := &idmaker{
		bizId: bizId,
		db:    db,
		step:  step,
	}

	//第一段
	if err := i.loadId(); err != nil {
		return err
	}

	//第二段
	if err := i.loadId2(false); err != nil {
		return err
	}

	idmakerMgr.add(i)

	return nil
}

func GetId(bizId BizId) (uint64, error) {
	i := idmakerMgr.get(bizId)
	if i == nil {
		return 0, fmt.Errorf("idmaker is nil")
	}

	return i.getId()
}

var idmakerMgr = &idmakerManager{
	idmakerMap: make(map[BizId]*idmaker),
}

type idmakerManager struct {
	idmakerMap map[BizId]*idmaker
	mutex      sync.RWMutex
}

func (im *idmakerManager) get(bizId BizId) *idmaker {
	im.mutex.RLock()
	defer im.mutex.RUnlock()

	return im.idmakerMap[bizId]
}

func (im *idmakerManager) add(i *idmaker) {
	im.mutex.Lock()
	defer im.mutex.Unlock()

	im.idmakerMap[i.bizId] = i
}

//第一段: [id, maxId)
//第二段: [id2, id2+step)
type idmaker struct {
	bizId BizId
	db    *sql.DB
	step  uint64
	id    uint64
	maxId uint64
	id2   uint64
	mutex sync.Mutex
}

//外面加锁
func (i *idmaker) loadId() error {
	id, err := i.getIdFromDb()
	if err != nil {
		return err
	}
	i.id = id
	i.maxId = i.id + i.step
	return nil
}

func (i *idmaker) loadId2(lock bool) error {
	if lock {
		i.mutex.Lock()
		defer i.mutex.Unlock()
	}

	if i.id2 > 0 {
		return nil
	}

	id, err := i.getIdFromDb()
	if err != nil {
		return err
	}
	i.id2 = id
	return nil
}

func (i *idmaker) getIdFromDb() (uint64, error) {
	tx, err := i.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("Begin err: %v", err)
	}
	defer func() {
		if err != nil {
			if er := tx.Rollback(); er != nil {
				fmt.Printf("Rollback failed err=%v\n", er)
			}
		}
	}()

	var id uint64

	err = tx.QueryRow(`select next_id from id_maker where biz_id = ? for update`, i.bizId).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("select id_maker err: %v", err)
	}

	_, err = tx.Exec(`update id_maker set next_id = next_id + ?, update_time = ? where biz_id = ?`,
		i.step, time.Now().Unix(), i.bizId)
	if err != nil {
		return 0, fmt.Errorf("update id_maker err: %v", err)
	}

	if er := tx.Commit(); er != nil {
		return 0, fmt.Errorf("Commit err: %v", er)
	}

	return id, nil
}

func (i *idmaker) getId() (uint64, error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	var id uint64

	if i.id2 <= 0 {
		go func() {
			_ = i.loadId2(true)
		}()
	}

	if i.id < i.maxId {
		id = i.id
		i.id++
		return id, nil
	}

	if i.id2 > 0 {
		i.id = i.id2
		i.maxId = i.id + i.step

		i.id2 = 0
		go func() {
			_ = i.loadId2(true)
		}()

		id = i.id
		i.id++
		return id, nil
	}

	if err := i.loadId(); err != nil {
		return 0, err
	}

	id = i.id
	i.id++
	return id, nil
}
