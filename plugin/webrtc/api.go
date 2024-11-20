package plugin_webrtc

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	. "github.com/pion/webrtc/v3"
	. "m7s.live/v5/plugin/webrtc/pkg"
)

// https://datatracker.ietf.org/doc/html/draft-ietf-wish-whip
func (conf *WebRTCPlugin) Push_(w http.ResponseWriter, r *http.Request) {
	streamPath := r.URL.Path[len("/push/"):]
	rawQuery := r.URL.RawQuery
	auth := r.Header.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		auth = auth[len("Bearer "):]
		if rawQuery != "" {
			rawQuery += "&bearer=" + auth
		} else {
			rawQuery = "bearer=" + auth
		}
		conf.Info("push", "stream", streamPath, "bearer", auth)
	}
	w.Header().Set("Content-Type", "application/sdp")
	w.Header().Set("Location", "/webrtc/api/stop/push/"+streamPath)
	w.Header().Set("Access-Control-Allow-Private-Network", "true")
	if rawQuery != "" {
		streamPath += "?" + rawQuery
	}
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	conn := Connection{
		PLI: conf.PLI,
		SDP: string(bytes),
	}
	conn.Logger = conf.Logger
	if conn.PeerConnection, err = conf.api.NewPeerConnection(Configuration{
		ICEServers: conf.ICEServers,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if conn.Publisher, err = conf.Publish(conf.Context, streamPath); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	conn.Publisher.RemoteAddr = r.RemoteAddr
	conf.AddTask(&conn)
	if err = conn.WaitStarted(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := conn.SetRemoteDescription(SessionDescription{Type: SDPTypeOffer, SDP: conn.SDP}); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if answer, err := conn.GetAnswer(); err == nil {
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, answer.SDP)
	} else {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (conf *WebRTCPlugin) Play_(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/sdp")
	streamPath := r.URL.Path[len("/play/"):]
	rawQuery := r.URL.RawQuery
	var conn Connection
	conn.EnableDC = conf.EnableDC
	bytes, err := io.ReadAll(r.Body)
	defer func() {
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}()
	if err != nil {
		return
	}
	conn.SDP = string(bytes)
	if conn.PeerConnection, err = conf.api.NewPeerConnection(Configuration{
		ICEServers: conf.ICEServers,
	}); err != nil {
		return
	}
	if rawQuery != "" {
		streamPath += "?" + rawQuery
	}
	if conn.Subscriber, err = conf.Subscribe(conn.Context, streamPath); err != nil {
		return
	}
	conn.Subscriber.RemoteAddr = r.RemoteAddr
	conf.AddTask(&conn)
	if err = conn.WaitStarted(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = conn.SetRemoteDescription(SessionDescription{Type: SDPTypeOffer, SDP: conn.SDP}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if sdp, err := conn.GetAnswer(); err == nil {
		w.Write([]byte(sdp.SDP))
	} else {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
