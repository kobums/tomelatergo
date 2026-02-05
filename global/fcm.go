package global

type FcmType int

const (
	_ FcmType = iota
	FcmNotice
	FcmAlarm
	FcmComment
	FcmMatch
)

type Fcm struct {
	Code    string
	Target  []int64
	Message map[string]string
	Title   string
	Type    FcmType
}

var _fcmCh chan Fcm

func init() {
	_fcmCh = make(chan Fcm, 1000)
}

func SendFcm(item Fcm) {
	_fcmCh <- item
}

func GetFcm() chan Fcm {
	return _fcmCh
}
