package feed

import (
	"errors"
	"reflect"
	"testing"
)

type captureJSONWriter struct {
	writes []any
	err    error
}

func (w *captureJSONWriter) WriteJSON(v any) error {
	w.writes = append(w.writes, v)
	return w.err
}

func TestWSSSubscribeAssetsQueuesAndDeduplicates(t *testing.T) {
	w := NewWSSClient([]string{"a", "a", ""})
	added, err := w.SubscribeAssets("a", "b", "", "b")
	if err != nil || added != 1 {
		t.Fatalf("added=%d err=%v", added, err)
	}
	if got := w.subscribedAssets(); !reflect.DeepEqual(got, []string{"a", "b"}) {
		t.Fatalf("assets=%v", got)
	}
}

func TestWSSSubscribeAssetsWritesDynamicUpdate(t *testing.T) {
	w := NewWSSClient([]string{"a"})
	writer := &captureJSONWriter{}
	w.writer = writer
	added, err := w.SubscribeAssets("b", "c")
	if err != nil || added != 2 {
		t.Fatalf("added=%d err=%v", added, err)
	}
	if len(writer.writes) != 1 {
		t.Fatalf("writes=%v", writer.writes)
	}
	msg, ok := writer.writes[0].(map[string]any)
	if !ok || !reflect.DeepEqual(msg["assets_ids"], []string{"b", "c"}) || msg["operation"] != "subscribe" {
		t.Fatalf("message=%#v", writer.writes[0])
	}
}

func TestWSSSubscribeAssetsRetainsAssetOnWriteFailure(t *testing.T) {
	w := NewWSSClient([]string{"a"})
	w.writer = &captureJSONWriter{err: errors.New("closed")}
	added, err := w.SubscribeAssets("b")
	if added != 1 || err == nil {
		t.Fatalf("added=%d err=%v", added, err)
	}
	if got := w.subscribedAssets(); !reflect.DeepEqual(got, []string{"a", "b"}) {
		t.Fatalf("assets=%v", got)
	}
}
