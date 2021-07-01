package main

import (
  "fmt"
  "math"
  "os"
  "time"

  "github.com/gen2brain/dlgs"
  "github.com/getlantern/systray"
)

var secondsUntilNextTrigger int = 60 * 20
var ticker *time.Ticker = time.NewTicker(time.Second)
var timeUpdateChannel = make(chan int, 1)
var err error;

func main() {
  go mainLoop()
  go func() {
    _, err = dlgs.Info("320 Active", "320 has been activated and will remind you every 20 minutes to look away from the screen for a bit.")
    handle(err)
  }()
  systray.Run(trayOnReady, func() {})
}

func trigger320() {
  _, err = dlgs.Info("320 Start", "When you click OK, start looking at something 20 feet (around 6 metres) away. I'll tell you when to stop.")
  handle(err)
  time.Sleep(20 * time.Second)
  _, err = dlgs.Info("320 End", "You can stop looking now. I'll remind you again in 20 minutes!")
  handle(err)
}

func mainLoop() {
  for {
    <-ticker.C
    secondsUntilNextTrigger--
    timeUpdateChannel <- secondsUntilNextTrigger
    if secondsUntilNextTrigger <= 0 {
      trigger320()
      secondsUntilNextTrigger = 60 * 20
    }
  }
}

func trayOnReady() {
  //TODO: Add an icon
  systray.SetTitle("320")
  time.Sleep(500 * time.Millisecond)
  systray.SetTooltip("Three Twenty")
  trayTrigger320 := systray.AddMenuItem("Trigger 320 Now", "There's no such thing as too much eye resting...")
  trayQuit := systray.AddMenuItem("Quit", "Quits Three Twenty")
  systray.AddSeparator()
  trayTime := systray.AddMenuItem("Getting time until next break...", "Time until next break")
  trayTime.Disable()

  // Update time left until next break
  go func() {
    for {
      <-timeUpdateChannel
      trayTime.SetTitle(fmt.Sprintf("Time until next break: %dm %ds", int(math.Floor(float64(secondsUntilNextTrigger/60))), secondsUntilNextTrigger%60))
    }
  }()

  // Handle buttons
  go func() {
    for {
      select {
      case <-trayQuit.ClickedCh:
        systray.Quit()
        os.Exit(0)
      case <-trayTrigger320.ClickedCh:
        trigger320()
      }
    }
  }()
}

func handle(err error) {
  if err != nil {
    panic(err)
  }
}
