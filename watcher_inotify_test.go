// +build linux,!fsnotify

package notify

import (
	"io"
	"os"
	"testing"
	"time"
)

var fixtureos = Fixture{
	IN_ACCESS: func(p string) error {
		f, err := os.OpenFile(p, os.O_RDWR, 0755)
		if err != nil {
			return err
		}
		if _, err := f.WriteString(p); err != nil {
			f.Close()
			return err
		}
		f.Close()
		f, err = os.Open(p)
		if err != nil {
			return err
		}
		if _, err = f.Read([]byte{0x00}); err != nil && err != io.EOF {
			f.Close()
			return err
		}
		return f.Close()

	},
}

func init() {
	for e, f := range fixtureos {
		fixture[e] = f
	}
}

func TestEventMaskEvent(t *testing.T) {
	tests := []struct {
		passed   Event
		got      Event
		received Event
	}{
		{Create, IN_CREATE, Create},
		{Create, IN_MOVED_TO, Create},
		{Create | IN_CREATE, IN_CREATE, IN_CREATE},
		{Create | IN_MOVED_TO, IN_MOVED_TO, IN_MOVED_TO},
		{Create | IN_CREATE, IN_MOVED_TO, Create},
		{Create | IN_MOVED_TO, IN_CREATE, Create},
		{Create | IN_MOVED_TO | IN_CREATE, IN_CREATE, IN_CREATE},

		{Delete, IN_DELETE, Delete},
		{Delete, IN_DELETE_SELF, Delete},
		{Delete | IN_DELETE, IN_DELETE, IN_DELETE},
		{Delete | IN_DELETE_SELF, IN_DELETE_SELF, IN_DELETE_SELF},
		{Delete | IN_DELETE, IN_DELETE_SELF, Delete},
		{Delete | IN_DELETE_SELF, IN_DELETE, Delete},
		{Delete | IN_DELETE_SELF | IN_DELETE, IN_DELETE, IN_DELETE},

		{Write, IN_MODIFY, Write},
		{Write | IN_MODIFY, IN_MODIFY, IN_MODIFY},

		{Move, IN_MOVED_FROM, Move},
		{Move, IN_MOVE_SELF, Move},
		{Move | IN_MOVED_FROM, IN_MOVED_FROM, IN_MOVED_FROM},
		{Move | IN_MOVE_SELF, IN_MOVE_SELF, IN_MOVE_SELF},
		{Move | IN_MOVED_FROM, IN_MOVE_SELF, Move},
		{Move | IN_MOVE_SELF, IN_MOVED_FROM, Move},
		{Move | IN_MOVE_SELF | IN_MOVED_FROM, IN_MOVED_FROM, IN_MOVED_FROM},

		{IN_MOVE, IN_MOVED_FROM, IN_MOVED_FROM},
		{IN_MOVE, IN_MOVED_TO, IN_MOVED_TO},
		{Move | IN_MOVE, IN_MOVED_FROM, IN_MOVED_FROM},
		{Create | IN_MOVE, IN_MOVED_TO, IN_MOVED_TO},

		{IN_CLOSE, IN_CLOSE_WRITE, IN_CLOSE_WRITE},
		{IN_CLOSE, IN_CLOSE_NOWRITE, IN_CLOSE_NOWRITE},
		{Move | IN_CLOSE, IN_CLOSE_WRITE, IN_CLOSE_WRITE},
		{Create | IN_CLOSE, IN_CLOSE_NOWRITE, IN_CLOSE_NOWRITE},

		{All, IN_CREATE, Create},
		{All, IN_MOVED_TO, Create},
		{All, IN_DELETE, Delete},
		{All, IN_DELETE_SELF, Delete},
		{All, IN_MODIFY, Write},
		{All, IN_MOVED_FROM, Move},
		{All, IN_MOVE_SELF, Move},

		{IN_ALL_EVENTS, IN_ACCESS, IN_ACCESS},
		{IN_ALL_EVENTS, IN_MODIFY, IN_MODIFY},
		{IN_ALL_EVENTS, IN_ATTRIB, IN_ATTRIB},
		{IN_ALL_EVENTS, IN_CLOSE_WRITE, IN_CLOSE_WRITE},
		{IN_ALL_EVENTS, IN_CLOSE_NOWRITE, IN_CLOSE_NOWRITE},
		{IN_ALL_EVENTS, IN_OPEN, IN_OPEN},
		{IN_ALL_EVENTS, IN_MOVED_FROM, IN_MOVED_FROM},
		{IN_ALL_EVENTS, IN_MOVED_TO, IN_MOVED_TO},
		{IN_ALL_EVENTS, IN_CREATE, IN_CREATE},
		{IN_ALL_EVENTS, IN_DELETE, IN_DELETE},
		{IN_ALL_EVENTS, IN_DELETE_SELF, IN_DELETE_SELF},
		{IN_ALL_EVENTS, IN_MOVE_SELF, IN_MOVE_SELF},

		{IN_ACCESS, IN_ACCESS, IN_ACCESS},
		{IN_MODIFY, IN_MODIFY, IN_MODIFY},
		{IN_ATTRIB, IN_ATTRIB, IN_ATTRIB},
		{IN_CLOSE_WRITE, IN_CLOSE_WRITE, IN_CLOSE_WRITE},
		{IN_CLOSE_NOWRITE, IN_CLOSE_NOWRITE, IN_CLOSE_NOWRITE},
		{IN_OPEN, IN_OPEN, IN_OPEN},
		{IN_MOVED_FROM, IN_MOVED_FROM, IN_MOVED_FROM},
		{IN_MOVED_TO, IN_MOVED_TO, IN_MOVED_TO},
		{IN_CREATE, IN_CREATE, IN_CREATE},
		{IN_DELETE, IN_DELETE, IN_DELETE},
		{IN_DELETE_SELF, IN_DELETE_SELF, IN_DELETE_SELF},
		{IN_MOVE_SELF, IN_MOVE_SELF, IN_MOVE_SELF},
	}

	for i, test := range tests {
		if e := decodemask(uint32(test.passed), uint32(test.got)); e != test.received {
			t.Errorf("want event=%v, got %v (test: %d)", test.received, e, i)
		}
	}
}

func TestInotify(t *testing.T) {
	ei := []EventInfo{
		EI("github.com/rjeczalik/fs/fs.go", IN_ACCESS),
		// EI("github.com/rjeczalik/fs/binfs/", IN_MODIFY),
		// EI("github.com/rjeczalik/fs/binfs.go", IN_ATTRIB),
		// EI("github.com/rjeczalik/fs/binfs_test.go", IN_CLOSE_WRITE),
		// EI("github.com/rjeczalik/fs/binfs/", IN_CLOSE_NOWRITE),
		// EI("github.com/rjeczalik/fs/binfs/", IN_OPEN),
		// EI("github.com/rjeczalik/fs/fs_test.go", IN_MOVED_FROM),
		// EI("github.com/rjeczalik/fs/binfs/", IN_MOVED_TO),
		// EI("github.com/rjeczalik/fs/binfs.go", IN_CREATE),
		// EI("github.com/rjeczalik/fs/binfs_test.go", IN_DELETE),
		// EI("github.com/rjeczalik/fs/binfs/", IN_DELETE_SELF),
		// EI("github.com/rjeczalik/fs/binfs/", IN_MOVE_SELF),
		// EI("github.com/rjeczalik/fs/fs_test.go", IN_ACCESS),
		// EI("github.com/rjeczalik/fs/binfs/", Create),
		// EI("github.com/rjeczalik/fs/binfs.go", Create),
		// EI("github.com/rjeczalik/fs/binfs_test.go", Create),
		// EI("github.com/rjeczalik/fs/binfs/", Delete),
		// EI("github.com/rjeczalik/fs/binfs/", Create),
	}
	test(t, newWatcher(), IN_ACCESS, ei, time.Second)
}