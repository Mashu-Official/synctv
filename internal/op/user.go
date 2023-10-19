package op

import (
	"errors"
	"hash/crc32"
	"sync/atomic"

	"github.com/synctv-org/synctv/internal/db"
	"github.com/synctv-org/synctv/internal/model"
	"github.com/zijiren233/stream"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	model.User
	version uint32
}

func (u *User) Version() uint32 {
	return atomic.LoadUint32(&u.version)
}

func (u *User) CheckVersion(version uint32) bool {
	return atomic.LoadUint32(&u.version) == version
}

func (u *User) CreateRoom(name, password string, conf ...db.CreateRoomConfig) (*model.Room, error) {
	return db.CreateRoom(name, password, append(conf, db.WithCreator(&u.User))...)
}

func (u *User) NewMovie(movie model.MovieInfo) model.Movie {
	return model.Movie{
		MovieInfo: movie,
		CreatorID: u.ID,
	}
}

func (u *User) HasPermission(room *Room, permission model.Permission) bool {
	return room.HasPermission(&u.User, permission)
}

func (u *User) DeleteRoom(room *Room) error {
	if !u.HasPermission(room, model.CanDeleteRoom) {
		return errors.New("no permission")
	}
	return DeleteRoom(room)
}

func (u *User) NeedPassword() bool {
	return len(u.HashedPassword) != 0
}

func (u *User) SetPassword(password string) error {
	if u.CheckPassword(password) && u.NeedPassword() {
		return errors.New("password is the same")
	}
	var hashedPassword []byte
	if password != "" {
		var err error
		hashedPassword, err = bcrypt.GenerateFromPassword(stream.StringToBytes(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
	}
	u.HashedPassword = hashedPassword

	atomic.StoreUint32(&u.version, crc32.ChecksumIEEE(u.HashedPassword))
	return db.SetUserHashedPassword(u.ID, hashedPassword)
}
