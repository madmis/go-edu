package session

import (
    "sync"
    "fmt"
    "net/http"
)

// Store contains all data for one session process with specific id.
type Store interface {
    Set(key, value interface{}) error     //set session value
    Get(key interface{}) interface{}      //get session value
    Delete(key interface{}) error         //delete session value
    SessionID() string                    //back current sessionID
    SessionRelease(w http.ResponseWriter) // release the resource & save data to provider & return the data
    Flush() error                         //delete all data
}

// Provider contains global session methods and saved SessionStores.
// it can operate a SessionStore by its id.
type Provider interface {
    SessionInit(gclifetime int64, config string) error
    SessionRead(sid string) (Store, error)
    SessionExist(sid string) bool
    SessionRegenerate(oldsid, sid string) (Store, error)
    SessionDestroy(sid string) error
    SessionAll() int //get all active session
    SessionGC()
}

var providers = make(map[string]Provider)

type Manager struct {
    cookieName  string     //private cookiename
    lock        sync.Mutex // protects session
    provider    Provider
    maxlifetime int64
}

func NewManager(provideName, cookieName string, maxlifetime int64) (*Manager, error) {
    provider, ok := providers[provideName]
    if !ok {
        return nil, fmt.Errorf("session: unknown provide %q (forgotten import?)", provideName)
    }
    return &Manager{provider: provider, cookieName: cookieName, maxlifetime: maxlifetime}, nil
}