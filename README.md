# qsms

A golang implementation of [@qcloudsms](https://github.com/qcloudsms) 's SMS API

## Usage

```go
    package main

    import (
        "github.com/mysqto/qsms"
        "log"
        "time"
    )

    func main() {
        // Get appID & appKey from the tencent cloud console
        appID := ""
        appKey := ""
        receiver := "+8618888888888"
        now := time.Now().Unix()
        client := qsms.NewClient(appID, appKey) // create a new SMS Client

        // sending a sms using the Client
        result, err := client.SendSMS("+8618888888888", "hello, world")
        if err != nil || result.Status != 0 {
            log.Printf("error sending sms : [err = %v, result = %v]", err, result)
        } else {
            // result.MessageID is server-side messageId, can be paired with result.LocalMessageID to identify single message
            log.Printf("sms send with : status = %s, messageID = %s", result.ErrMsg, result.MessageID)
        }

        // Pull the send result
        rptResult, err := client.MobileStatusPull(receiver, now-1*60*60, now)

        if err != nil {
            log.Printf("error pulling sms status: [err = %v, multiResult = %v]", err, result)
        } else {
            for _, report := range rptResult.Reports {
                // report.MessageID is same with mtResult.MessageID in a SMS
                log.Printf("sms deliver status : %v", report)
            }
        }
    }
```

see [example](https://github.com/mysqto/qsms/blob/master/example/) for more detailed usage
