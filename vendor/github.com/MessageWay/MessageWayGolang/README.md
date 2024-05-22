![messageWay](logo.png)

![Swagger][ico-swagger]
[![MessageWay][ico-MSGWay]][link-MSGWay]

# MessageWay Golang SDK

A Golang SDK for the MessageWay API.

----

# Installation

```shell
$ go get github.com/MessageWay/MessageWayGolang
```

---

## Send Message

### Via SMS

```go
package main

import (
	"fmt"
	MessageWay "github.com/MessageWay/MessageWayGolang"
)

func main() {
	app := MessageWay.New(MessageWay.Config{
		ApiKey: "YOUR_APIKEY",
	})
	res, err := app.Send(MessageWay.Message{
		Method:     "sms",
		Mobile:     "09123456789",
		TemplateID: 3,
	})
	if err != nil {
		return
	}
	fmt.Println(res.ReferenceID)
}
```

### Via Whatsapp Messenger

```go
package main

import (
	"fmt"
	MessageWay "github.com/MessageWay/MessageWayGolang"
)

func main() {
	app := MessageWay.New(MessageWay.Config{
		ApiKey: "YOUR_APIKEY",
	})
	message := MessageWay.NewBuilder().SetMobile("09123456789").SetParams("foo", "doo", "loo").ViaWhatsapp().SetTemplateID(3).Build()
	res, err := app.Send(message)
	if err != nil {
		return
	}
	fmt.Println(res.ReferenceID)
}
```

### Via Gap Messenger

```go
package main

import (
	"fmt"
	MessageWay "github.com/MessageWay/MessageWayGolang"
)

func main() {
	app := MessageWay.New(MessageWay.Config{
		ApiKey: "YOUR_APIKEY",
	})
	message := MessageWay.NewBuilder().SetMobile("09123456789").ViaGap().SetTemplateID(3).Build()
	res, err := app.Send(message)
	if err != nil {
		return
	}
	fmt.Println(res.ReferenceID)
}
```

### Via IVR

```go
package main

import (
	"fmt"
	MessageWay "github.com/MessageWay/MessageWayGolang"
)

func main() {
	app := MessageWay.New(MessageWay.Config{
		ApiKey: "YOUR_APIKEY",
	})
	message := MessageWay.NewBuilder().SetMobile("09123456789").ViaIVR().SetTemplateID(2).Build()
	res, err := app.Send(message)
	if err != nil {
		return
	}
	fmt.Println(res.ReferenceID)
}
```

---

## Get Status

```go
package main

import (
	"fmt"
	MessageWay "github.com/MessageWay/MessageWayGolang"
)

func main() {
	app := MessageWay.New(MessageWay.Config{
		ApiKey: "YOUR_APIKEY",
	})
	res, err := app.GetStatus(MessageWay.StatusRequest{
		ReferenceID: "1234567890",
	})
	if err != nil {
		return
	}
	fmt.Println(res)
}
```

---

## Verify OTP

```go
package main

import (
	"fmt"
	MessageWay "github.com/MessageWay/MessageWayGolang"
)

func main() {
	app := MessageWay.New(MessageWay.Config{
		ApiKey: "YOUR_APIKEY",
	})
	res, err := app.Verify(MessageWay.OTPVerifyRequest{
		Mobile: "9123456789",
		OTP:    "123456",
	})
	if err != nil {
		return
	}
	fmt.Println(res)
}
```

---

## License

MIT

[ico-MSGWay]: https://img.shields.io/badge/-MSGWay.com-critical?link=https://MSGWay.com&style=for-the-badge

[ico-swagger]: https://img.shields.io/swagger/valid/3.0?specUrl=https%3A%2F%2Fdoc.msgway.com%2Fswagger.json&style=for-the-badge

[link-MSGWay]: https://MSGWay.com/

