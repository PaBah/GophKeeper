package server

import (
	"errors"
	"sync"

	pb "github.com/PaBah/GophKeeper/internal/gen/proto/gophkeeper/v1"
)

var (
	ErrUserSyncExists = errors.New("this sync for this user already exists")
	ErrSyncExists     = errors.New("sync for this user already exists")
	ErrSyncNotFound   = errors.New("sync not found")
	ErrSyncUnexpected = errors.New("unexpected sync error")
)

// SyncStatus manages the synchronization status of users and their operations.
type SyncStatus struct {
	syncClients map[string]pb.GophKeeperService_SubscribeToChangesServer
	rwMutex     *sync.RWMutex
}

// NewSyncStatus initializes a new SyncStatus instance.
func NewSyncStatus() *SyncStatus {
	instance := &SyncStatus{
		syncClients: make(map[string]pb.GophKeeperService_SubscribeToChangesServer),
		rwMutex:     &sync.RWMutex{},
	}

	return instance
}

//// AddSync adds a synchronization session for a user.
//func (ss *SyncStatus) AddSync(userID string, syncID uuid.UUID) (bool, error) {
//	ss.rwMutex.RLock()
//	// fmt.Println("start lock r")
//	val, ok := ss.syncClients[userID]
//	// fmt.Println("try unlock r RemoveClientSync")
//	ss.rwMutex.RUnlock()
//
//	if ok {
//		if val == syncID {
//			// fmt.Println("sync already exist for user")
//			return false, fmt.Errorf("%w", ErrSyncExists)
//		}
//
//		return false, fmt.Errorf("%w", ErrUserSyncExists)
//	}
//
//	// fmt.Println("try lock for w AddSync")
//	ss.rwMutex.Lock()
//	// fmt.Println("locked for w AddSync")
//
//	ss.syncClients[userID] = syncID
//	// fmt.Println("try unlock for w AddSync")
//	ss.rwMutex.Unlock()
//	// fmt.Println("unlock for w AddSync")
//	// fmt.Println("add val for sync AddSync")
//
//	return true, nil
//}
//
//// RemoveClientSync removes a synchronization session for a user.
//func (ss *SyncStatus) RemoveClientSync(userID string, syncID uuid.UUID) error {
//	ss.rwMutex.RLock()
//	// fmt.Println("start lock r RemoveClientSync")
//
//	val, ok := ss.syncClients[userID]
//	// fmt.Println("try unlock r RemoveClientSync")
//	ss.rwMutex.RUnlock()
//	if !ok || val != syncID {
//		return fmt.Errorf("%w", ErrSyncNotFound)
//	}
//
//	// fmt.Println("try lock for w RemoveClientSync")
//	ss.rwMutex.Lock()
//	// fmt.Println("locked for w RemoveClientSync")
//
//	delete(ss.syncClients, userID)
//	// fmt.Println("try unlock for w RemoveClientSync")
//	ss.rwMutex.Unlock()
//	// fmt.Println("unlock for w RemoveClientSync")
//
//	return nil
//}

//// IsSyncExists checks if a synchronization session exists for a user.
//func (ss *SyncStatus) IsSyncExists(userID string, syncID uuid.UUID) (bool, error) {
//	ss.rwMutex.RLock()
//	// fmt.Println("start lock r IsSyncExists")
//	val, ok := ss.syncClients[userID]
//	// fmt.Println("try unlock r IsSyncExists")
//	ss.rwMutex.RUnlock()
//	if ok && val == syncID {
//		// fmt.Println("sync already exist for user")
//
//		return true, nil
//	}
//	// fmt.Println("this sync is not exist")
//
//	return false, nil
//}
