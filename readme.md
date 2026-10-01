# errorx

golang error framework, support `with` , `stack` , `close` features.

- with

  ```go
  type Code int
  
  const CodeNotFound Code = 1
  
  err := errorx.WithT(sql.ErrNoRows, CodeNotFound) // attach a business code
  errorx.IsT(err, CodeNotFound)                    // true
  errorx.T[Code](err)                              // CodeNotFound
  ```

  

- stack

  ```go
  err := errorx.New("open config") // captures the stack here
  fmt.Printf("%+v\n", err)
  // open config:
  // /app/main.go:12
  // /app/main.go:30
  
  err = errorx.WithStack(err) // attach a stack to an existing error
  ```
  
  
  
- close

  ```go
  package main
  
  import (
  	"net"
  	"os"
  
  	"github.com/lysShub/errorx-go"
  )
  
  type Foo struct {
  	incoming, outgoing net.Conn
  	fh                 *os.File
  	closeErr           errorx.CloseErr
  }
  
  func New() (*Foo, error) {
  	var f = &Foo{}
  
  	if c, err := net.Dial("tcp", "aaa.com"); err != nil {
  		return nil, f.close(err)
  	} else {
  		f.incoming = c
  	}
  	if c, err := net.Dial("tcp", "bbb.com"); err != nil {
  		return nil, f.close(err)
  	} else {
  		f.outgoing = c
  	}
  	if fh, err := os.Open("snapshot.bin"); err != nil {
  		return nil, f.close(err)
  	} else {
  		f.fh = fh
  	}
  
  	go f.service()
  	return f, nil
  }
  
  func (f *Foo) close(cause error) error {
  	return f.closeErr.Close(func() (errs []error) {
  		errs = append(errs, cause)
  
  		if f.incoming != nil {
  			errs = append(errs, f.incoming.Close())
  		}
  		if f.outgoing != nil {
  			errs = append(errs, f.outgoing.Close())
  		}
  		if f.fh != nil {
  			errs = append(errs, f.fh.Close())
  		}
  		return errs
  	})
  }
  func (f *Foo) Close() error { return f.close(nil) }
  
  func (f *Foo) service() (_ error) {
  	var b = make([]byte, 0xffff)
  
  	for {
  		n, err := f.incoming.Read(b)
  		if err != nil {
  			return f.close(err)
  		}
  
  		_, err = f.fh.Write(b[:n])
  		if err != nil {
  			return f.close(err)
  		}
  
  		_, err = f.outgoing.Write(b[:])
  		if err != nil {
  			return f.close(err)
  		}
  	}
  }
  ```
  
  
  
  
  
