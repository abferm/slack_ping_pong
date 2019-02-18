package main

import (
  "fmt"
  "os"
  "strings"
  "runtime"
  "time"

  "github.com/nlopes/slack"
)

func getenv(name string) string {
  v := os.Getenv(name)
  if v == "" {
    panic("missing required environment variable " + name)
  }
  return v
}

func main() {
  startTime := time.Now()
  token := getenv("SLACKTOKEN")
  api := slack.New(token)
  rtm := api.NewRTM()
  go rtm.ManageConnection()

Loop:
  for {
    select {
    case msg := <-rtm.IncomingEvents:
      fmt.Print("Event Received: ")
      switch ev := msg.Data.(type) {

      case *slack.MessageEvent:
        info := rtm.GetInfo()

        text := ev.Text
        text = strings.TrimSpace(text)
        text = strings.ToLower(text)
        var response string

        switch text {
        case "os":
          response = runtime.GOOS
        case "arch":
          response = runtime.GOARCH
        case "uptime":
          response = time.Since(startTime).String()
        case "env":
          environ := os.Environ()
          response = strings.Join(environ, "\n")
        case "help":
          response = "Accepted Commands:\n \tos : Prints runtime.GOOS\n\tarch : Prints runtime.GOARCH\n\tuptime : prints application uptime\n\tenv : prints environment variables\n\thelp : prints this message"
        default:
          response = "Unknown command: " + text
        }

        if ev.User != info.User.ID {
          rtm.SendMessage(rtm.NewOutgoingMessage(response, ev.Channel))
        }

      case *slack.RTMError:
        fmt.Printf("Error: %s\n", ev.Error())

      case *slack.InvalidAuthEvent:
        fmt.Printf("Invalid credentials")
        break Loop

      default:
        // Take no action
      }
    }
  }
}
