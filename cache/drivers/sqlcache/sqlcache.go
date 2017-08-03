//Package sqlcache provides cache driver uses sqlite or mysql to store cache data.
//Using database/sql as driver.
//You should create data table with sql file in "sql" folder first.
package sqlcache

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/herb-go/herb/cache"
)

const modelSet = 0
const modelUpdate = 1

var defaultGCPeriod = 5 * time.Minute
var tokenMask = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var defaultGcLimit = int64(100)

//Cache The sql cache Driver.
type Cache struct {
	DB           *sql.DB
	table        string
	name         string
	ticker       *time.Ticker
	quit         chan int
	gcErrHandler func(err error)
	gcLimit      int64
}

//SetGCErrHandler Set callback to handler error raised when gc.
func (c *Cache) SetGCErrHandler(f func(err error)) {
	c.gcErrHandler = f
	return
}
func (c *Cache) start() error {
	err := c.gc()
	return err
}
func (c *Cache) getVersionTx(tx *sql.Tx) ([]byte, error) {
	var version []byte
	stmt, err := tx.Prepare(`Select version from ` + c.table + ` WHERE cache_key="" AND cache_name = ?`)
	if err != nil {
		return version, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(c.name)
	err = row.Scan(&version)
	if err == sql.ErrNoRows {
		return []byte{}, nil
	}
	return version, err
}

//SearchByPrefix Search All key start with given prefix.
//Return All matched key and any error raised.
func (c *Cache) SearchByPrefix(prefix string) ([]string, error) {
	return nil, cache.ErrFeatureNotSupported
}
func (c *Cache) gc() error {
	var keys []string

	tx, err := c.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	version, err := c.getVersionTx(tx)
	if err != nil {
		return err
	}

	stmtExpired, err := tx.Prepare(`Select cache_key FROM ` + c.table + ` Where cache_name = ? AND expired > -1  AND expired < ? limit ?`)
	if err != nil {
		return err
	}
	defer stmtExpired.Close()

	rows, err := stmtExpired.Query(c.name, time.Now().Unix(), c.gcLimit)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var key string
		err = rows.Scan(&key)
		if err != nil {
			return err
		}
		keys = append(keys, key)

	}
	stmtVersionWrong, err := tx.Prepare(`Select cache_key FROM ` + c.table + ` Where cache_name = ? AND version != ? limit ?`)
	if err != nil {
		return err
	}
	defer stmtVersionWrong.Close()
	rows2, err := stmtVersionWrong.Query(c.name, version, c.gcLimit)
	if err != nil {
		return err
	}
	defer rows2.Close()
	for rows2.Next() {
		var key string
		err = rows2.Scan(&key)
		if err != nil {
			return err
		}
		keys = append(keys, key)

	}

	stmt, err := tx.Prepare(`DELETE FROM ` + c.table + ` Where cache_name=? and cache_key = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, v := range keys {
		_, err = stmt.Exec(c.name, v)
		if err != nil {
			return err
		}
	}
	if err == nil {
		tx.Commit()
	}
	return err
}

//IncrCounter Increase int val in cache by given key.Count cache and data cache are in two independent namespace.
//Return int data value and any error raised.
func (c *Cache) IncrCounter(key string, increment int64, ttl time.Duration) (int64, error) {
	var v int64
	tx, err := c.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	version, err := c.getVersionTx(tx)
	if err != nil {
		return 0, err
	}
	stmt, err := tx.Prepare(`Select cache_value from ` + c.table + ` WHERE ( expired < 0 OR expired > ?) AND cache_name =? AND cache_key = ? AND version=?`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	r := stmt.QueryRow(time.Now().Unix(), c.name, key, version)
	var j string
	err = r.Scan(&j)
	if err == sql.ErrNoRows {
		v = 0
	} else if err != nil {
		return 0, err
	} else {
		err = json.Unmarshal([]byte(j), &v)
		if err != nil {
			v = 0
		}
	}
	v = v + increment
	val, err := json.Marshal(v)
	if err != nil {
		return 0, err
	}
	var expired int64
	if ttl < 0 {
		expired = -1
	} else {
		expired = time.Now().Add(ttl).Unix()
	}
	stmtset, err := tx.Prepare(`update ` + c.table + ` set
	 cache_value=?,
	 version=?,
	 expired=?
	 Where cache_name=? 
	 and cache_key=?
	 `)

	defer stmtset.Close()
	row, err := stmtset.Exec(
		val,
		version,
		expired,
		c.name,
		key)
	if err != nil {
		return 0, err
	}
	affected, err := row.RowsAffected()
	if err != nil {
		return v, err
	}
	if affected == 0 {
		stmt2, err := tx.Prepare(`insert into ` + c.table + ` (cache_name,cache_key,cache_value,version,expired) values (?,?,?,?,?)`)
		if err != nil {
			return v, err
		}
		defer stmt2.Close()
		_, err = stmt2.Exec(c.name, key, string(val), version, expired)
	}
	if err != nil {
		return v, err
	}
	tx.Commit()
	return v, nil
}

//SetCounter Set int val in cache by given key.Count cache and data cache are in two independent namespace.
//Return any error raised.
func (c *Cache) SetCounter(key string, v int64, ttl time.Duration) error {
	return c.Set(key, v, ttl)
}

//GetCounter Get int val from cache by given key.Count cache and data cache are in two independent namespace.
//Return int data value and any error raised.
func (c *Cache) GetCounter(key string) (int64, error) {
	var v int64
	err := c.Get(key, &v)
	return v, err
}

//DelCounter Delete int val in cache by given key.Count cache and data cache are in two independent namespace.
//Return any error raised.
func (c *Cache) DelCounter(key string) error {
	return c.DelCounter(key)
}

//Set Set data model to cache by given key.
//Return any error raised.
func (c *Cache) Set(key string, v interface{}, ttl time.Duration) error {
	return c.doSet(key, v, ttl, modelSet)
}

//Update Update data model to cache by given key only if the cache exist.
//Return any error raised.
func (c *Cache) Update(key string, v interface{}, ttl time.Duration) error {
	return c.doSet(key, v, ttl, modelUpdate)
}
func (c *Cache) doSet(key string, v interface{}, ttl time.Duration, mode int) error {
	tx, err := c.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	version, err := c.getVersionTx(tx)
	if err != nil {
		return err
	}
	val, err := json.Marshal(v)
	if err != nil {
		return err
	}
	var expired int64
	if ttl < 0 {
		expired = -1
	} else {
		expired = time.Now().Add(ttl).Unix()
	}
	stmt, err := tx.Prepare(`update ` + c.table + ` set
	 cache_value=?,
	 version=?,
	 expired=?
	 Where cache_name=? 
	 and cache_key=?
	 `)

	defer stmt.Close()
	r, err := stmt.Exec(
		val,
		version,
		expired,
		c.name,
		key)
	if err != nil {
		return err
	}
	affected, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 && mode != modelUpdate {
		stmt2, err := tx.Prepare(`insert into ` + c.table + ` (cache_name,cache_key,cache_value,version,expired) values (?,?,?,?,?)`)
		if err != nil {
			return err
		}
		defer stmt2.Close()
		_, err = stmt2.Exec(c.name, key, string(val), version, expired)
	}
	if err != nil {
		return err
	}
	tx.Commit()
	return err
}

//Get Get data model from cache by given key.
//Parameter v should be pointer to empty data model which data filled in.
//Return any error raised.
func (c *Cache) Get(key string, v interface{}) error {

	tx, err := c.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	version, err := c.getVersionTx(tx)
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`Select cache_value from ` + c.table + ` WHERE ( expired < 0 OR expired > ?) AND cache_name =? AND cache_key = ? AND version=?`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	r := stmt.QueryRow(time.Now().Unix(), c.name, key, version)
	var j string
	err = r.Scan(&j)
	if err == sql.ErrNoRows {
		return cache.ErrNotFound
	}
	if err != nil {
		return err
	}
	tx.Commit()
	err = json.Unmarshal([]byte(j), &v)
	return err
}

//Flush Delete all data in cache.
//Return any error if raised
func (c *Cache) Flush() error {
	tx, err := c.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	version, err := c.getVersionTx(tx)
	if err != nil {
		return err
	}
	newversion, err := cache.NewRandMaskedBytes(tokenMask, 16, version)
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`update ` + c.table + ` set
	 cache_value=?,
	 version=?,
	 expired=?
	 Where cache_name=? and cache_key=""
	 `)
	if err != nil {
		return err
	}
	defer stmt.Close()
	r, err := stmt.Exec(
		"",
		string(newversion),
		-1,
		c.name)
	if err != nil {
		return err
	}
	affected, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		stmt2, err := tx.Prepare(`insert into ` + c.table + ` (cache_name,cache_value,version,expired,cache_key) 
		values (?,?,?,?,"")`)
		if err != nil {
			return err
		}
		defer stmt2.Close()
		_, err = stmt2.Exec(c.name, string(newversion), newversion, -1)

	}
	if err != nil {
		return err
	}
	tx.Commit()
	err = c.gc()
	if err != nil {
		return err
	}
	return err
}

//Del Delete data in cache by given key.
//Return any error raised.
func (c *Cache) Del(key string) error {
	tx, err := c.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare(`DELETE FROM ` + c.table + ` WHERE cache_name= ? and cache_key = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(c.name, key)
	if err == nil {
		tx.Commit()
	}
	return err
}

//SetBytesValue Set bytes data to cache by given key.
//Return any error raised.
func (c *Cache) SetBytesValue(key string, bytes []byte, ttl time.Duration) error {
	b, err := json.Marshal(bytes)
	if err != nil {
		return err
	}
	return c.Set(key, b, ttl)
}

//UpdateBytesValue Update bytes data to cache by given key only if the cache exist.
//Return any error raised.
func (c *Cache) UpdateBytesValue(key string, bytes []byte, ttl time.Duration) error {
	b, err := json.Marshal(bytes)
	if err != nil {
		return err
	}
	return c.Update(key, b, ttl)
}

//GetBytesValue Get bytes data from cache by given key.
//Return data bytes and any error raised.
func (c *Cache) GetBytesValue(key string) ([]byte, error) {
	b := []byte{}
	err := c.Get(key, &b)
	if err != nil {
		return nil, err
	}

	v := []byte{}
	err = json.Unmarshal(b, &v)
	return v, err
}

//Close Close cache.
//Return any error if raised
func (c *Cache) Close() error {
	err := c.gc()
	if err != nil {
		return nil
	}
	close(c.quit)
	return c.DB.Close()
}

//Config Cache driver config.
type Config struct {
	Driver   string //Registered sql driver.
	Conn     string //Conn string of database.``
	Table    string //Database table name.
	Name     string //Database cache name.
	GCPeriod int64  //Period of gc.Default value is 5 minute.
	GCLimit  int64  ////Max delete limit in every gc call.Default value is 100.
}

//New Create new cache driver with given json bytes.
//Return new driver and any error raised.
func (c *Cache) New(config json.RawMessage) (cache.Driver, error) {
	cf := Config{}
	err := json.Unmarshal(config, &cf)
	if err != nil {
		return nil, err
	}
	cache := Cache{}
	cache.DB, err = sql.Open(cf.Driver, cf.Conn)
	if err != nil {
		return &cache, err
	}
	cache.table = cf.Table
	cache.quit = make(chan int)
	period := time.Duration(cf.GCPeriod)
	if period == 0 {
		period = defaultGCPeriod
	}
	cache.ticker = time.NewTicker(period)
	gcLimit := cf.GCLimit
	if gcLimit == 0 {
		gcLimit = defaultGcLimit
	}
	cache.gcLimit = gcLimit
	cache.name = cf.Name
	go func() {
		for {
			select {
			case <-cache.ticker.C:
				err := cache.gc()
				if err != nil {
					if cache.gcErrHandler != nil {
						cache.gcErrHandler(err)
					}
				}
			case <-cache.quit:
				cache.ticker.Stop()
				return
			}
		}

	}()
	err = cache.start()
	if err != nil {
		return &cache, err
	}
	return &cache, nil
}

func init() {
	cache.Register("sqlcache", &Cache{})
}
