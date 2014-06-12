package main

import (
	"code.google.com/p/go.net/websocket"
	"errors"
	"github.com/golang/glog"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// var homeTempl = template.Must(template.ParseFiles("home.html"))

// func homeHandle(w http.ResponseWriter, r *http.Request) {
// 	homeTempl.Execute(w, r.Host)
// }

const (
	killedByOtherDevice  = "Another device login %d"
	wronggloginParams     = "Wrong login params %s"
	wronggloginType       = "Wrong login type %s"
	wronggloginDevice     = "Wrong login deviceId %s"
	wronggloginTimestamp  = "Wrong login timestamp %s"
	usergloginTimeout     = "User %d login timeout"
	userReconnectTimeout = "Reconnect timeout %d"
	wrongMd5Check        = "User %d has wrong md5"
	wronggloginTimeout    = "Wrong login %d timeout %s"

	LOGIN_PARAM_COUNT = 6
	READ_TIMEOUT      = 10
	PING_MSG          = "p"
	PONG_MSG          = "P"
	TIME_OUT          = 3 * 60         // 3 mins
	EXPIRE_TIME       = uint64(1 * 60) // 1 mins

	publicKey = "BlackCrystal"
)

var (
	LOGIN_PARAM_ERROR = errors.New("glogin params parse error!")
	ParamsError       = &ErrorCode{2001, "登陆参数错误"}
	gloginFailed       = &ErrorCode{2002, "登陆失败"}

	AckgloginOK             = []byte{byte(0)} // 登陆成功
	AckWrongParams         = []byte{byte(1)} // 错误的登陆参数
	AckWronggloginType      = []byte{byte(2)} // 登陆类型解析错误
	AckWronggloginDevice    = []byte{byte(3)} // 登陆设备ID解析错误
	AckWronggloginTimestamp = []byte{byte(4)} // 登陆时间戳解析错误
	AckgloginTimeout        = []byte{byte(5)} // 登陆超时
	AckWrongMD5            = []byte{byte(6)} // 错误的md5
	AckOtherglogoned        = []byte{byte(7)} // 您已在别处登陆
	AckWronggloginTimeout   = []byte{byte(8)} // 超时解析错误
)

type ErrorCode struct {
	ErrorId   int
	ErrorDesc string
}

// StartHttp start http listen.
func StartHttp(bindAddrs []string) {
	for _, bind := range bindAddrs {
		glog.Infof("start websocket listen addr:[%s]\n", bind)
		go websocketListen(bind)
	}
}

func websocketListen(bindAddr string) {
	httpServeMux := http.NewServeMux()
	//httpServeMux.HandleFunc("/", homeHandle)
	httpServeMux.Handle("/ws", websocket.Handler(WsHandler))
	server := &http.Server{
		Addr: bindAddr,
		Handler: httpServeMux,
		ReadTimeout: READ_TIMEOUT * time.Second,
	}
	err := server.ListenAndServe()

	//l, err := net.Listen("tcp", bindAddr)
	//if err != nil {
	//	glog.Error("net.Listen(\"tcp\", \"%s\") error(%v)", bindAddr, err)
	//	panic(err)
	//}
	//err := server.Serve(l)
	if err != nil {
		glog.Error("server.Serve(\"%s\") error(%v)", bindAddr, err)
		panic(err)
	}
}

func getgloginParams(req string) (id int64, mac, alias, expire, hmac string, err error) {
	args := strings.Split(req, "|")
	if len(args) != LOGIN_PARAM_COUNT {
		err = LOGIN_PARAM_ERROR
		return
	}
	// skip the 0
	id, err = strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return
	}
	mac = args[2]
	alias = args[3]
	expire = args[4]
	hmac = args[5]
	return
}

func isAuth(id int64, mac, alias, expire, hmac string) bool {
	expireTime, err := strconv.ParseUint(expire, 10, 64)
	if err != nil || uint64(time.Now().Unix())-expireTime >= EXPIRE_TIME {
		return true
		// return false
	}
	// TODO: check hmac is equal
	return true
}

func setReadTimeout(conn net.Conn, delaySec int) error {
	return conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(delaySec)))
}

func WsHandler(ws *websocket.Conn) {
	addr := ws.Request().RemoteAddr
	var err error
	if err = setReadTimeout(ws, 10); err != nil {
		glog.Errorf("[%s] websocket.SetReadDeadline() error(%s)\n", addr, err)
		return
	}
	reply := ""
	if err = websocket.Message.Receive(ws, &reply); err != nil {
		glog.Errorf("[%s] websocket.Message.Receive() error(%s)\n", addr, err)
		return
	}
	// 旧程序需要成功登陆后的一次回复
	websocket.Message.Send(ws, 0)

	glog.Infof("Recv glogin %s\n", reply)
	// parse login params
	id, mac, alias, expire, hmac, loginErr := getgloginParams(reply)
	if loginErr != nil {
		glog.Errorf("[%s] params error (%s)\n", addr, reply)
		websocket.Message.Send(ws, ParamsError)
		return
	}
	// check login
	if !isAuth(id, mac, alias, expire, hmac) {
		glog.Errorf("[%s] auth failed:\"%s\" error(%s)\n", addr, reply)
		websocket.Message.Send(ws, gloginFailed)
		return
	}
	localHost := ws.LocalAddr().String()
	s := NewSession(id, alias, mac, ws)
	gSessionList.AddSession(s)
	SetUserOnline(id, localHost)

	start := time.Now().UnixNano()
	end := int64(start + int64(time.Second))
	for {
		// more then 1 sec, reset the timer
		if end-start >= int64(time.Second) {
			if err = setReadTimeout(ws, TIME_OUT); err != nil {
				glog.Errorf("<%s> user_id:\"%s\" websocket.SetReadDeadline() error(%s)\n", addr, id, err)
				break
			}
			start = end
		}

		if err = websocket.Message.Receive(ws, &reply); err != nil {
			glog.Errorf("<%s> user_id:\"%s\" websocket.Message.Receive() error(%s)\n", addr, id, err)
			break
		}
		if reply == PING_MSG {
			if err = websocket.Message.Send(ws, PONG_MSG); err != nil {
				glog.Errorf("<%s> user_id:\"%s\" write heartbeat to client error(%s)\n", addr, id, err)
				break
			}
			//glog.Debugf("<%s> user_id:\"%s\" receive heartbeat\n", addr, id)
		} else {
			//glog.Debugf("<%s> user_id:\"%s\" recv msg %s\n", addr, id, reply)
			// Send to Message Bus

		}
		end = time.Now().UnixNano()
	}
	SetUserOffline(id, localHost)
	// remove exists conn
	gSessionList.RemoveSession(s)
	return
}