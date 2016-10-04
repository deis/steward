package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/deis/k8s-claimer/client"
	conf "github.com/deis/steward/config"
	"github.com/juju/loggo"
)

const (
	clusterVersion = ""
	clusterRegex   = ""
	leaseLengthSec = 1200 // 20 minutes
	kubeConfigFile = "kubeconfig.yaml"
)

var logger = loggo.GetLogger("")

func init() {
	logger.SetLogLevel(loggo.TRACE)
}

type config struct {
	Server    string `envconfig:"K8S_CLAIMER_SERVER" default:"k8s-claimer.champagne.deis.com"`
	AuthToken string `envconfig:"K8S_CLAIMER_AUTH_TOKEN" require:"true"`
}

func main() {
	cfg := new(config)
	if err := conf.Load(cfg); err != nil {
		logger.Errorf("Error parsing test driver configuration: %s", err)
		os.Exit(1)
	}
	logger.Infof("Leasing k8s cluster from %s for %d seconds...", cfg.Server, leaseLengthSec)
	resp, err := client.CreateLease(cfg.Server, cfg.AuthToken, clusterVersion, clusterRegex, leaseLengthSec)
	if err != nil {
		logger.Errorf("Error leasing k8s cluster: %s", err)
		os.Exit(1)
	}
	logger.Infof("Acquired lease on cluster %s; token is %s", resp.ClusterName, resp.Token)
	logger.Infof("Writing cluster configuration to %s", kubeConfigFile)
	var exitCode int
	if kubeConfigBytes, err := resp.KubeConfigBytes(); err == nil {
		if err := ioutil.WriteFile(kubeConfigFile, kubeConfigBytes, 0644); err == nil {
			argsStr := strings.Join(os.Args[1:], " ")
			logger.Infof("Executing `%s`", argsStr)
			cmd := exec.Command(os.Args[1], os.Args[2:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err == nil {
				logger.Infof("`%s` completed successfully", argsStr)
			} else {
				logger.Errorf("Errors or test failures detected; cleaning up before exiting")
				exitCode = 1
			}
		} else {
			logger.Errorf("Errors writing cluster configuration to %s; cleaning up before exiting", kubeConfigFile)
			exitCode = 1
		}
	} else {
		logger.Errorf("Errors getting cluster configuration bytes")
		exitCode = 1
	}
	logger.Infof("Deleting cluster configuration at %s", kubeConfigFile)
	os.Remove(kubeConfigFile)
	// Deliberately ignore any errors from deleting kubeconfig file... just move on to deleting the
	// lease-- which is more important.
	logger.Infof("Giving up lease on cluster %s", resp.ClusterName)
	if err := client.DeleteLease(cfg.Server, cfg.AuthToken, resp.Token); err != nil {
		logger.Errorf("Error giving up lease on cluster %s: %s", resp.ClusterName, err)
		os.Exit(1)
	}
	logger.Infof("Successfully gave up lease on cluster %s", resp.ClusterName)
	os.Exit(exitCode)
}
