package main

import (
	"fmt"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
	"os"
	"time"
)

func main() {
	m, err := mgr.Connect()
	if err != nil {
		fmt.Printf("connect to service manager error = %s", err.Error())
	}
	defer m.Disconnect()

	s, err := m.OpenService("cassandra")
	if err != nil {
		fmt.Printf("connect to cassandra error = %s", err.Error())
	}
	defer s.Close()

	status, err := s.Query()
	if err != nil {
		fmt.Printf("query status from cassandra error = %s", err.Error())
	}
	if status.State == svc.Running {
		_, err := s.Control(svc.Stop)
		if err != nil {
			fmt.Printf("stop cassandra service error = %s", err.Error())
			os.Exit(1)
		}

		for {
			time.Sleep(5 * time.Second)

			status, err := s.Query()
			if err != nil {
				fmt.Printf("query status from cassandra error = %s", err.Error())
				os.Exit(1)
			}
			if status.State == svc.Stopped {
				fmt.Printf("cassandra service is stopped\n")

				break
			}

			fmt.Printf(".")
		}
	}

	fmt.Println("do sth...")

	if err = s.Start([]string{}...); err != nil {
		fmt.Printf("start cassandra error = %s", err.Error())
	}

	for {
		time.Sleep(5 * time.Second)

		status, err := s.Query()
		if err != nil {
			fmt.Printf("query status from cassandra error = %s", err.Error())
			os.Exit(1)
		}
		if status.State == svc.Running {
			fmt.Printf("cassandra service is running\n")

			break
		}

		fmt.Printf(".")
	}

	args := os.Args[1:]
	for i, v := range args {
		code, err := MsgBox(fmt.Sprintf("[%d]", i), v, MsgBoxBtnOk)
		if err != nil {
			fmt.Printf("[%d] error = %s", i, err.Error())
			continue
		}

		fmt.Printf("[%d] code = %d", i, code)
	}
}
