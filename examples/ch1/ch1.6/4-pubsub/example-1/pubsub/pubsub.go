// Package pubsub implements a simple multi-topic pub-sub library.
package pubsub

import (
    "sync"
    "time"
)

type (
    subscriber chan interface{}         // subscriber kiểu channel
    topicFunc  func(v interface{}) bool // topic là một filter
)

type Publisher struct {
    m           sync.RWMutex             // khóa đọc ghi
    buffer      int                      // kích thước  hàng đợi subscribe
    timeout     time.Duration            // hết thời gian publish
    subscribers map[subscriber]topicFunc // thông tin subscriber
}

// constructor với timeout và độ dài hàng đợi
func NewPublisher(publishTimeout time.Duration, buffer int) *Publisher {
    return &Publisher{
        buffer:      buffer,
        timeout:     publishTimeout,
        subscribers: make(map[subscriber]topicFunc),
    }
}

// Thêm subscriber mới, đăng ký hết tất cả topic
func (p *Publisher) Subscribe() chan interface{} {
    return p.SubscribeTopic(nil)
}

// Thêm subscriber mới, đăng ký các filter được lọc theo topic
func (p *Publisher) SubscribeTopic(topic topicFunc) chan interface{} {
    ch := make(chan interface{}, p.buffer)
    p.m.Lock()
    p.subscribers[ch] = topic
    p.m.Unlock()
    return ch
}

// hủy đăng ký
func (p *Publisher) Evict(sub chan interface{}) {
    p.m.Lock()
    defer p.m.Unlock()

    delete(p.subscribers, sub)
    close(sub)
}

// đăng 1 topic
func (p *Publisher) Publish(v interface{}) {
    p.m.RLock()
    defer p.m.RUnlock()

    var wg sync.WaitGroup
    for sub, topic := range p.subscribers {
        wg.Add(1)
        go p.sendTopic(sub, topic, v, &wg)
    }
    wg.Wait()
}

// Đóng 1 đối tượng publisher và đóng tất cả các subscriber
func (p *Publisher) Close() {
    p.m.Lock()
    defer p.m.Unlock()

    for sub := range p.subscribers {
        delete(p.subscribers, sub)
        close(sub)
    }
}

// Gửi 1 topic có thể duy trì trong thời gian chờ wg
func (p *Publisher) sendTopic(
    sub subscriber, topic topicFunc, v interface{}, wg *sync.WaitGroup,
) {
    defer wg.Done()
    if topic != nil && !topic(v) {
        return
    }

    select {
    case sub <- v:
    case <-time.After(p.timeout):
    }
}