package session

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type cacheSession struct {
	Sess *Session
	T    string
}

type DataBase struct {
	Memory

	db      *sql.DB
	driver  string
	address string
	cache   map[string]cacheSession
}

// 连接数据库
func (d *DataBase) Connect() error {
	var err error

	d.db, err = sql.Open(d.driver, d.address)
	if err != nil {
		return err
	}

	ctx := context.Background()

	// 测试数据库
	if err := d.db.PingContext(ctx); err != nil {
		return err
	}

	return nil
}

// 关闭连接
func (d *DataBase) Disconnected() error {
	return d.db.Close()
}

// 从数据库拉取数据
func (d *DataBase) Pull() error {
	sqlCmd := "SELECT SESS_ID,SESS_DATA,LIFE_CYCLE_START_TIME FROM SESSION_DATA"

	rows, err := d.db.Query(sqlCmd)
	if err != nil {
		return err
	}

	defer rows.Close()

	// 清空数据
	d.Empty()

	for rows.Next() {
		var id string
		var data []byte
		var lifeCycleStart int64

		err := rows.Scan(&id, &data, &lifeCycleStart)
		if err != nil {
			return err
		}

		// 判断生命周期是否结束
		if (lifeCycleStart + d.opt.LifeCycle) < time.Now().Unix() {
			continue
		}

		// 转成session对象
		sess, err := binaryToSession(data)
		if err != nil {
			return err
		}

		sess.LifeCycleStart = lifeCycleStart

		d.Data[id] = sess

	}

	return nil
}

// 推送数据
func (d *DataBase) Push() error {
	if len(d.Data) == 0 {
		_, err := d.db.Exec("DELETE FROM SESSION_DATA")
		if err != nil {
			return err
		}
	}

	defer d.cleanCache()

	for id, sess := range d.Data {
		cacheSess, ok := d.cache[id]
		if !ok {
			if !sess.change {
				continue
			}

			sess.change = false

			d.RemoteUpdate(id)
			continue
		}

		cacheT := cacheSess.T

		if cacheT == "REMOVE" {
			err := d.RemoteRemove(id)

			if err != nil {
				return err
			}

			continue
		}

		if cacheT == "CREATE" {
			err := d.RemoteCreate(id, cacheSess.Sess)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

// 清理数据库生命周期以结束的数据
func (d *DataBase) RemoteClean() error {
	_, err := d.db.Exec("DELETE FROM SESSION_DATA WHERE LIFE_CYCLE_END_TIME<?",
		time.Now().Unix())

	return err
}

// 删除数据库指定 session id 的数据
func (d *DataBase) RemoteRemove(idArr ...string) error {
	stmt, err := d.db.Prepare("DELETE FROM SESSION_DATA WHERE SESS_ID=?")
	if err != nil {
		return err
	}

	for _, id := range idArr {
		_, err = stmt.Exec(id)
		if err != nil {
			return fmt.Errorf("%s: %v", id, err)
		}
	}

	return nil
}

// 更新数据库指定 session 数据
func (d *DataBase) RemoteUpdate(id string) error {
	sess := d.Get(id)
	if sess == nil {
		return fmt.Errorf("there is no session object corresponding to id")
	}

	data, err := sess.ToBinary()
	if err != nil {
		return err
	}

	_, err = d.db.Exec("UPDATE SESSION_DATA SET SESS_DATA=? WHERE SESS_ID=?", data, id)

	return err
}

// 追加新创建数据(数据库内无记录)
func (d *DataBase) RemoteCreate(id string, sess *Session) error {

	data, err := sess.ToBinary()
	if err != nil {
		return err
	}

	_, err = d.db.Exec("INSERT INTO SESSION_DATA VALUES(?,?,?,?)",
		id, data, sess.LifeCycleStart, sess.LifeCycleStart+d.opt.LifeCycle)

	return err
}

// 创建一个新的session id 和 对应的 session 对象
func (d *DataBase) Create() (id string, sess *Session) {
	id, sess = d.Memory.Create()

	d.cache[id] = cacheSession{
		T:    "CREATE",
		Sess: sess,
	}

	return
}

// 删除本地数据
// 可以通过 Push 函数 推送到数据库
func (d *DataBase) Remove(id string) {
	d.Memory.Remove(id)

	d.cache[id] = cacheSession{
		T: "REMOVE",
	}
}

// 清除缓存
func (d *DataBase) cleanCache() {
	for k := range d.cache {
		delete(d.cache, k)
	}
}

// 垃圾回收
func (d *DataBase) gc() {

	if d.opt.CleanInterval <= 0 {
		return
	}

	time.AfterFunc(time.Duration(d.opt.CleanInterval)*time.Second, func() {
		err := d.RemoteClean()
		if err != nil {
			d.opt.HandlerGCError(err)
		}

		err = d.Pull()
		if err != nil {
			d.opt.HandlerGCError(err)
		}

		d.gc()
	})
}

func NewDataBase(driver string, address string, opt *Options) (d *DataBase, err error) {
	d = &DataBase{
		Memory:  *NewMemory(opt),
		cache:   make(map[string]cacheSession),
		address: address,
		driver:  driver,
	}

	if err = d.Connect(); err != nil {
		return
	}

	err = d.Pull()

	d.gc()

	return
}
