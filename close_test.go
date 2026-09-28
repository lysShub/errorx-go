package errorx_test

import (
	"context"
	"errors"
	"math"
	"net"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/lysShub/errorx-go"
	"golang.org/x/sync/errgroup"
)

func Test_CloseErr(t *testing.T) {
	t.Run("close函数为nil", func(t *testing.T) {
		var closeErr errorx.CloseErr

		err := closeErr.Close(nil)
		if err != nil {
			t.Fatal(err)
		}

		err = closeErr.Close(nil)
		if !errors.Is(err, errorx.ErrClosed) {
			t.Fatalf("expect %v, got %v", errorx.ErrClosed, err)
		}
		if !closeErr.Closed() {
			t.Fatal("expect closed")
		}
	})

	t.Run("自定义ClosedErr", func(t *testing.T) {
		var closeErr = errorx.CloseErr{ErrClosed: net.ErrClosed}

		err := closeErr.Close(nil)
		if err != nil {
			t.Fatal(err)
		}

		err = closeErr.Close(nil)
		if !errors.Is(err, net.ErrClosed) {
			t.Fatalf("expect %v, got %v", net.ErrClosed, err)
		}
		if !closeErr.Closed() {
			t.Fatal("expect closed")
		}
	})

	t.Run("close没有发生错误", func(t *testing.T) {
		var closeErr errorx.CloseErr

		err := closeErr.Close(func() (errs []error) { return nil })
		if err != nil {
			t.Fatal(err)
		}

		err = closeErr.Close(nil)
		if !errors.Is(err, errorx.ErrClosed) {
			t.Fatalf("expect %v, got %v", errorx.ErrClosed, err)
		}
		if !closeErr.Closed() {
			t.Fatal("expect closed")
		}
	})

	t.Run("close中发生错误, 只记录第一个错误", func(t *testing.T) {
		var closeErr errorx.CloseErr

		e1 := closeErr.Close(func() (errs []error) {
			errs = append(errs, errors.New("1234"))
			errs = append(errs, errors.New("5678"))
			return
		})
		if !strings.Contains(e1.Error(), "1234") {
			t.Fatalf("expect contains %q, got %q", "1234", e1.Error())
		}
		if strings.Contains(e1.Error(), "5678") {
			t.Fatalf("expect not contains %q, got %q", "5678", e1.Error())
		}

		e2 := closeErr.Close(nil)
		if !strings.Contains(e2.Error(), "1234") {
			t.Fatalf("expect contains %q, got %q", "1234", e2.Error())
		}
		if strings.Contains(e2.Error(), "5678") {
			t.Fatalf("expect not contains %q, got %q", "5678", e2.Error())
		}

		// close出错, 重复close, 也会返回最初错误
		e3 := closeErr.Close(nil)
		if !strings.Contains(e3.Error(), "1234") {
			t.Fatalf("expect contains %q, got %q", "1234", e3.Error())
		}
		if strings.Contains(e3.Error(), "5678") {
			t.Fatalf("expect not contains %q, got %q", "5678", e3.Error())
		}
	})

	t.Run("直接调用Error()、Closed()", func(t *testing.T) {
		var closeErr errorx.CloseErr
		if closeErr.Closed() {
			t.Fatal("expect not closed")
		}
		if err := closeErr.Error(); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("closing时调用Error()", func(t *testing.T) {
		var closeErr errorx.CloseErr

		var wg = &sync.WaitGroup{}
		wg.Add(1)
		go func() {
			closeErr.Close(func() (errs []error) {
				wg.Done()
				time.Sleep(time.Second)
				return nil
			})
		}()
		wg.Wait() // 确保携程已执行

		start := time.Now()
		if !errors.Is(closeErr.Error(), errorx.ErrClosed) {
			t.Fatalf("expect %v, got %v", errorx.ErrClosed, closeErr.Error())
		}
		if d := time.Since(start); math.Abs(float64(d-time.Second)) > float64(time.Millisecond*100) {
			t.Fatalf("expect ~%v, got %v", time.Second, d)
		}
	})
	t.Run("closing时调用Closed()", func(t *testing.T) {
		var closeErr errorx.CloseErr

		var wg = &sync.WaitGroup{}
		wg.Add(1)
		go func() {
			closeErr.Close(func() (errs []error) {
				wg.Done()
				time.Sleep(time.Second)
				return nil
			})
		}()
		wg.Wait() // 确保携程已执行

		start := time.Now()
		if !closeErr.Closed() {
			t.Fatal("expect closed")
		}
		if d := time.Since(start); math.Abs(float64(d-time.Duration(0))) > float64(time.Millisecond*100) {
			t.Fatalf("expect ~%v, got %v", time.Duration(0), d)
		}
	})
	t.Run("close后调用Error()、Closed()", func(t *testing.T) {
		var closeErr errorx.CloseErr
		closeErr.Close(nil)

		if !closeErr.Closed() {
			t.Fatal("expect closed")
		}
		if !closeErr.Closed() {
			t.Fatal("expect closed")
		}
	})

	t.Run("closing时调用Close()", func(t *testing.T) {
		var closeErr errorx.CloseErr

		var wg = &sync.WaitGroup{}
		wg.Add(1)
		go func() {
			closeErr.Close(func() (errs []error) {
				wg.Done()
				time.Sleep(time.Second * 2)
				return nil
			})
		}()

		wg.Wait()

		start := time.Now()
		closeErr.Close(func() (errs []error) {
			errs = append(errs, os.ErrNotExist)
			return errs
		})
		if d := time.Since(start); math.Abs(float64(d-time.Second*2)) > float64(time.Millisecond*100) {
			t.Fatalf("expect ~%v, got %v", time.Second*2, d)
		}

		if !closeErr.Closed() {
			t.Fatal("expect closed")
		}
	})

	t.Run("Done", func(t *testing.T) {
		var err errorx.CloseErr
		go func() {
			time.Sleep(time.Second)
			err.Close(func() (errs []error) { return nil })
		}()
		start := time.Now()
		<-err.Done()
		if d := time.Since(start); math.Abs(float64(d-time.Second)) > float64(time.Millisecond*100) {
			t.Fatalf("expect ~%v, got %v", time.Second, d)
		}
	})

	t.Run("Done Closed", func(t *testing.T) {
		var err errorx.CloseErr
		err.Close(func() (errs []error) { return nil })
		start := time.Now()
		<-err.Done()
		if d := time.Since(start); math.Abs(float64(d-time.Duration(0))) > float64(time.Millisecond*100) {
			t.Fatalf("expect ~%v, got %v", time.Duration(0), d)
		}
	})

	t.Run("parallel", func(t *testing.T) {
		var eg, _ = errgroup.WithContext(context.Background())

		var closeErr errorx.CloseErr
		for i := 0; i < 0xff; i++ {
			if i == 0 {
				err := closeErr.Close(func() (errs []error) {
					time.Sleep(time.Second)
					errs = append(errs, errors.New("1234"))
					return
				})
				if err.Error() != "1234" {
					t.Errorf("expect %q, got %q", "1234", err.Error())
				}
			}
			if i%2 == 0 {
				eg.Go(func() error {
					err := closeErr.Close(func() (errs []error) {
						time.Sleep(time.Second)
						errs = append(errs, errors.New("1234"))
						return
					})
					if err.Error() != "1234" {
						t.Errorf("expect %q, got %q", "1234", err.Error())
					}
					return nil
				})
			} else {
				eg.Go(func() error {
					err := closeErr.Close(nil)
					if err.Error() != "1234" {
						t.Errorf("expect %q, got %q", "1234", err.Error())
					}
					return nil
				})
			}
		}
		eg.Wait()
	})
}
