module mining-monitoring

go 1.14

require (
	github.com/aliyun/alibaba-cloud-sdk-go v1.61.774
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fastly/go-utils v0.0.0-20180712184237-d95a45783239 // indirect
	github.com/gin-gonic/gin v1.6.3
	github.com/go-gomail/gomail v0.0.0-20160411212932-81ebce5c23df
	github.com/googollee/go-socket.io v1.4.4
	github.com/gorilla/websocket v1.4.2
	github.com/graarh/golang-socketio v0.0.0-20170510162725-2c44953b9b5f
	github.com/jehiah/go-strftime v0.0.0-20171201141054-1d33003b3869 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/lestrrat-go/strftime v1.0.3 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5
	github.com/robfig/cron/v3 v3.0.1
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.7.0
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
	github.com/stretchr/testify v1.6.1 // indirect
	github.com/tebeka/strftime v0.1.5 // indirect
	github.com/tidwall/gjson v1.6.4
	github.com/zhouhui8915/engine.io-go v0.0.0-20150910083302-02ea08f0971f
	github.com/zhouhui8915/go-socket.io-client v0.0.0-20200925034401-83ee73793ba4
	golang.org/x/text v0.3.3
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df // indirect
)

//replace github.com/googollee/go-socket.io => ./extern/go-socket.io
