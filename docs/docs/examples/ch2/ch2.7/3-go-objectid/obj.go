package main

import "sync"

type ObjectId int32

// quản lý ánh xạ giữa các đối tượng ngôn ngữ Go
// và objectID
var refs struct {
    sync.Mutex
    objs map[ObjectId]interface{}
    next ObjectId
}

func init() {
    refs.Lock()
    defer refs.Unlock()

    refs.objs = make(map[ObjectId]interface{})
    refs.next = 1000
}

// NewObjectId được sử dụng để tạo ObjectId
// liên kết với đối tượng
func NewObjectId(obj interface{}) ObjectId {
    refs.Lock()
    defer refs.Unlock()

    id := refs.next
    refs.next++

    refs.objs[id] = obj
    return id
}

func (id ObjectId) IsNil() bool {
    return id == 0
}

// Get được sử dụng để decode đối tượng Go ban đầu
func (id ObjectId) Get() interface{} {
    refs.Lock()
    defer refs.Unlock()

    return refs.objs[id]
}

// Free được sử dụng để giải phóng liên kết
// của ObjectId và đối tượng Go ban đầu
func (id *ObjectId) Free() interface{} {
    refs.Lock()
    defer refs.Unlock()

    obj := refs.objs[*id]
    delete(refs.objs, *id)
    *id = 0

    return obj
}