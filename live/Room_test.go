package live

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/tursom/GoCollections/exceptions"
)

func getRoomImpl() (Room, exceptions.Exception) {
	cookieBytes, err := os.ReadFile("Room_test.cookie.txt")
	if err != nil {
		return nil, exceptions.NewPackageException("read cookie file failed", exceptions.Cfg().SetCause(err))
	}

	room := NewRoom(917818)
	room.SetCookie(string(cookieBytes))

	return room, nil
}

func Test_roomImpl_Send(t *testing.T) {
	room, err := getRoomImpl()
	if err != nil {
		t.Fatalf("get room impl failed: %s", err)
	}

	resp, err := room.Send("弹幕测试")
	if err != nil {
		t.Fatalf("test send danmu failed: %s", exceptions.GetStackTraceString(err.(exceptions.Exception)))
	}

	t.Logf("send danmu resp: %v", resp)
}

func Test_roomImpl_GetDanmuColors(t *testing.T) {
	room, err := getRoomImpl()
	if err != nil {
		t.Fatalf("get room impl failed: %s", err)
	}

	colors, exception := room.GetDanmuColors()
	if exception != nil {
		t.Fatalf("get danmu colors failed: %s", exceptions.GetStackTraceString(exception))
	}

	marshal, _ := json.Marshal(colors)
	t.Logf("get colors: %s", string(marshal))
}
