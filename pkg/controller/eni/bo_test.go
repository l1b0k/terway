package eni

import (
	"errors"
	"sync"
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
)

func TestBackoffManager(t *testing.T) {
	manager := &BackoffManager{}
	key := "test-key"

	// 创建一个简单的退避策略
	bo := wait.Backoff{
		Steps:    3,
		Duration: 10 * time.Millisecond,
		Factor:   2.0,
		Jitter:   0.1,
	}

	// 第一次调用应该创建一个新的 ResourceBackoff 并返回一个持续时间
	d1, err := manager.Get(key, bo)
	if err != nil {
		t.Fatalf("First Get failed: %v", err)
	}
	if d1 <= 0 {
		t.Errorf("Expected positive duration, got %v", d1)
	}

	// 验证 NextTS 已经被设置
	nextTS, ok := manager.GetNextTS(key)
	if !ok {
		t.Fatalf("Failed to get NextTS")
	}
	expectedTS := time.Now().Add(d1)
	if nextTS.Sub(expectedTS).Abs() > time.Second {
		t.Errorf("NextTS not set correctly. Got %v, expected around %v", nextTS, expectedTS)
	}

	// 等待一段时间，确保下一次调用不会因为太快而跳过退避
	time.Sleep(d1 + 5*time.Millisecond)

	// 第二次调用应该使用同一个 ResourceBackoff 并执行下一步退避
	d2, err := manager.Get(key, bo)
	if err != nil {
		t.Fatalf("Second Get failed: %v", err)
	}
	if d2 <= 0 {
		t.Errorf("Expected positive duration, got %v", d2)
	}

	// d2 应该大约是 d1 的两倍（考虑抖动）
	if float64(d2) < float64(d1)*1.5 || float64(d2) > float64(d1)*2.5 {
		t.Errorf("Expected d2 to be about twice d1. Got d1=%v, d2=%v", d1, d2)
	}

	// 测试删除功能
	manager.Del(key)
	_, ok = manager.GetNextTS(key)
	if ok {
		t.Error("Key should have been deleted")
	}

	// 测试步骤耗尽的情况
	boExhausted := wait.Backoff{
		Steps:    1,
		Duration: 10 * time.Millisecond,
	}

	// 第一次调用应该成功
	_, err = manager.Get("exhausted", boExhausted)
	if err != nil {
		t.Fatalf("Get with exhaustible backoff failed: %v", err)
	}

	time.Sleep(20 * time.Millisecond)

	// 第二次调用应该返回超时错误
	_, err = manager.Get("exhausted", boExhausted)
	if !errors.Is(err, errTimeOut) {
		t.Errorf("Expected timeout error, got: %v", err)
	}
}

func TestResourceBackoffConcurrency(t *testing.T) {
	manager := &BackoffManager{}
	key := "concurrent-key"

	// 创建一个简单的退避策略
	bo := wait.Backoff{
		Steps:    100, // 足够多的步骤
		Duration: 10 * time.Millisecond,
		Factor:   1.1,
	}

	// 并发调用 Get 方法
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				_, err := manager.Get(key, bo)
				if err != nil && !errors.Is(err, errTimeOut) {
					t.Errorf("Unexpected error: %v", err)
				}
				time.Sleep(time.Millisecond)
			}
		}()
	}

	wg.Wait()

	// 验证最终的 Steps 值已经减少
	v, ok := manager.store.Load(key)
	if !ok {
		t.Fatal("Key not found after concurrent operations")
	}
	actual := v.(*ResourceBackoff)

	// 因为多次调用，Steps 应该减少
	if actual.Bo.Steps >= 100 {
		t.Errorf("Steps not decremented properly: %d", actual.Bo.Steps)
	}
}
