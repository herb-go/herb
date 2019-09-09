package redispool

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

var defaultMaxIdle = 200
var defaultMaxAlive = 200
var defaultIdleTimeout = 60 * time.Second
var defualtConnectTimeout = 10 * time.Second

//New create a new redis pool
func New() *Pool {
	return &Pool{}
}

//Pool redis poll struct
type Pool struct {
	*redis.Pool
	// Network string of redis conn.
	Network string
	//Redis server address.
	Address string
	//Redis server password.
	Password string
	//ConnectTimeout redis connect timeout.
	ConnectTimeout time.Duration
	//ReadTimeout redis read timeout.
	ReadTimeout time.Duration
	//WriteTimeout redis write timeout
	WriteTimeout time.Duration
	//Redis server database id.
	Db int
	//Max idle conn in redis pool.
	MaxIdle int
	//Max Alive conn in redis pool.
	MaxAlive int
	//Idel conn time.
	IdleTimeout time.Duration
}

func (p *Pool) dial() (redis.Conn, error) {
	conn, err := redis.DialTimeout(p.Network, p.Address, p.ConnectTimeout, p.ReadTimeout, p.WriteTimeout)
	if err != nil {
		return nil, err
	}
	if p.Password != "" {
		_, err = conn.Do("auth", p.Password)
		if err != nil {
			conn.Close()
			return nil, err
		}
	}
	_, err = conn.Do("SELECT", p.Db)
	if err != nil {
		conn.Close()
		return nil, err
	}
	return conn, nil
}

//Open :open a redis poll then return it.
func (p *Pool) Open() *redis.Pool {
	maxIdle := p.MaxIdle
	if maxIdle == 0 {
		maxIdle = defaultMaxIdle
	}
	p.Pool = redis.NewPool(p.dial, maxIdle)
	p.Pool.MaxActive = p.MaxAlive
	if p.Pool.MaxActive == 0 {
		p.Pool.MaxActive = defaultMaxAlive
	}
	p.Pool.IdleTimeout = p.IdleTimeout
	if p.Pool.IdleTimeout == 0 {
		p.Pool.IdleTimeout = defaultIdleTimeout
	}
	p.Pool.Wait = true
	return p.Pool
}

//Config :redis pool config
type Config struct {
	//Network string of redis conn.
	Network string
	//Redis server address.
	Address string
	//Redis server password.
	Password string
	//Redis connect timeout in second
	ConnectTimeoutInSecond int64
	//Redis read timeout in second
	ReadTimeoutInSecond int64
	//Redis write timeout in second
	WriteTimeoutInSecond int64
	//Redis server database id.
	Db int
	//Max idle conn in redis pool.
	MaxIdle int
	//Max Alive conn in redis pool.
	MaxAlive int
	//Redis conn idle timeout in second
	IdleTimeoutInSecond int64
}

//NewConfig create new config
func NewConfig() *Config {
	return &Config{}
}

//ApplyTo apply confit to redis poll
func (c *Config) ApplyTo(p *Pool) error {
	p.Network = c.Network
	p.Address = c.Address
	p.Password = c.Password
	p.ConnectTimeout = time.Duration(c.ConnectTimeoutInSecond) * time.Second
	p.ReadTimeout = time.Duration(c.ReadTimeoutInSecond) * time.Second
	p.WriteTimeout = time.Duration(c.WriteTimeoutInSecond) * time.Second
	p.Db = c.Db
	p.MaxIdle = c.MaxIdle
	p.MaxAlive = c.MaxAlive
	p.IdleTimeout = time.Duration(c.IdleTimeoutInSecond) * time.Second
	return nil
}
