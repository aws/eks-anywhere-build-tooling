package utils

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	systemd "github.com/coreos/go-systemd/v22/dbus"
	"github.com/godbus/dbus/v5"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
)

const (
	KubeletService  = "kubelet.service"
	MultiUserTarget = "multi-user.target"
)

func WaitForSystemdService(service string, timeout time.Duration) error {
	fmt.Printf("Waiting for %s to come up\n", service)

	conn, err := systemd.NewConnection(func() (*dbus.Conn, error) {
		dbusConn, err := dbus.Dial("unix:path=/.bottlerocket/rootfs/run/dbus/system_bus_socket")
		if err != nil {
			return nil, errors.Wrap(err, "Error dialing br systemd")
		}
		err = dbusConn.Auth([]dbus.Auth{dbus.AuthExternal(strconv.Itoa(os.Getuid()))})
		if err != nil {
			dbusConn.Close()
			return nil, errors.Wrap(err, "Error running auth on dbus connection")
		}
		err = dbusConn.Hello()
		if err != nil {
			dbusConn.Close()
			return nil, errors.Wrap(err, "Error running hello handshake on dbus connection")
		}
		return dbusConn, nil
	})
	if err != nil {
		return errors.Wrap(err, "Error creating systemd dbus connection")
	}
	fmt.Println("Created dbus connection to talk to systemd")
	defer conn.Close()

	// The filter function here is an inverse filter, it will filter any included units and hence nil is provided
	statusChan, errChan := conn.SubscribeUnitsCustom(time.Second*1, 1, func(u1, u2 *systemd.UnitStatus) bool {
		return *u1 != *u2
	}, nil)

	for {
		select {
		case unitStatus := <-statusChan:
			fmt.Printf("Received status change: %+v\n", unitStatus)
			if _, ok := unitStatus[service]; ok {
				if unitStatus[service].ActiveState == "active" {
					if strings.HasSuffix(service, ".service") {
						if unitStatus[service].SubState == "running" {
							fmt.Printf("%s service is active and running\n", service)
							return nil
						}
					} else if strings.HasSuffix(service, ".target") {
						if unitStatus[service].SubState == "active" {
							fmt.Printf("%s service is active and running\n", service)
							return nil
						}
					}
				}
			}
		case err = <-errChan:
			fmt.Printf("Error received while checking for unit status: %v\n", err)
			return errors.Wrap(err, "Error while checking for kubelet status")
		// Timeout after timeout duration
		case <-time.After(timeout):
			return errors.New("Timeout checking for kubelet status")
		}
	}
}

func WaitFor200(url string, timeout time.Duration) error {
	fmt.Printf("Waiting for 200: OK on url %s\n", url)
	counter := 0

	timeoutSignal := time.After(timeout)

	for {
		counter++
		select {
		case <-timeoutSignal:
			return errors.New("Timeout occurred while waiting for 200 OK")
		default:
			fmt.Printf("******  Try %d, hitting url %s ****** \n", counter, url)
			tr := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
			client := &http.Client{
				Transport: tr,
				// Each call attempt will have a timeout of 1 minute
				Timeout: 1 * time.Minute,
			}
			resp, err := client.Get(url)
			if err != nil {
				fmt.Printf("Error occured while hitting url: %v\n", err)
				time.Sleep(time.Second * 10)
				continue
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return errors.Wrap(err, "Error reading body from response of http get call")
			}
			fmt.Println(string(body))
			if resp.StatusCode != 200 {
				time.Sleep(time.Second * 10)
			} else {
				return nil
			}
		}
	}
}

func WaitForManifestAndOptionallyKillCmd(cmd *exec.Cmd, checkFiles []string, shouldKill bool) error {
	// Check for files written and send ok signal back on channel
	// ctx is used here to cancel the goroutine if the timeout has occurred
	okChan := make(chan bool)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func(ctx context.Context, okChan chan bool) {
		for {
			time.Sleep(2 * time.Second)

			select {
			case <-ctx.Done():
				return
			default:
				allFilesExist := true
				for _, file := range checkFiles {
					if fileInfo, err := os.Stat(file); err == nil {
						fmt.Printf("File %s exists with size: %d\n", file, fileInfo.Size())
						if fileInfo.Size() == 0 {
							fmt.Printf("File %s doesnt not have any size yet\n", file)
						}
					} else if os.IsNotExist(err) {
						fmt.Printf("File %s doest not exist yet\n", file)
						allFilesExist = false
					}
				}

				// Send ok on the channel
				if allFilesExist {
					okChan <- true
					return
				}
			}
		}
	}(ctx, okChan)

	timeout := time.After(40 * time.Second)
	select {
	case <-okChan:
		fmt.Println("All files were created, exiting kubeadm command")
		if shouldKill {
			cmd.Process.Kill()
		}
		return nil
	case <-timeout:
		cmd.Process.Kill()
		cancel()
		return errors.New(fmt.Sprintf("command: %s killed after timeout", strings.Join(cmd.Args, " ")))
	}
}

func waitForPodLiveness(podDefinition *v1.Pod) error {
	for _, container := range podDefinition.Spec.Containers {
		// Validate if liveness probe exists on the definition
		if container.LivenessProbe != nil {
			livenessProbeHandler := container.LivenessProbe.HTTPGet
			scheme := ""

			if livenessProbeHandler.Scheme != "" {
				scheme = string(livenessProbeHandler.Scheme)
			} else {
				scheme = "http"
			}

			port, err := ResolveContainerPort(livenessProbeHandler.Port, &container)
			if err != nil {
				return errors.Wrap(err, "Error resolving liveness probe port")
			}

			url := fmt.Sprintf("%s://%s:%d%s", scheme, livenessProbeHandler.Host, port, livenessProbeHandler.Path)
			fmt.Printf("Waiting for probe check on pod: %s\n", podDefinition.Name)
			err = WaitFor200(url, 5*time.Minute)
			if err != nil {
				return errors.Wrap(err, "Error waiting for 200 OK")
			}
		}
	}
	return nil
}

func WaitForPods(podDefinitions []*v1.Pod) error {
	for _, pod := range podDefinitions {
		err := waitForPodLiveness(pod)
		if err != nil {
			return errors.Wrapf(err, "Error checking liveness probe for pod: %s", pod.Name)
		}
	}
	return nil
}
