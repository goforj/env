//go:build ignore
// +build ignore

package main

import (
	"log"
	"time"

	"github.com/goforj/env"
)

type snapshot struct {
	APP_ENV           string
	IsDev             bool
	IsLocal           bool
	IsStaging         bool
	IsProduction      bool
	IsTesting         bool
	IsTestingOrLocal  bool
	IsLocalOrStaging  bool
	Addr              string
	Debug             bool
	Port              int
	Limit             int64
	Workers           uint
	MaxItems          uint64
	Threshold         float64
	Timeout           time.Duration
	Peers             []string
	Limits            map[string]string
	AppEnvEnum        string
	OS                string
	Arch              string
	IsLinux           bool
	IsMac             bool
	IsWindows         bool
	IsBSD             bool
	IsUnix            bool
	IsContainerOS     bool
	IsDocker          bool
	IsDockerInDocker  bool
	IsDockerHost      bool
	IsContainer       bool
	IsKubernetes      bool
	IsHostEnvironment bool
}

func main() {
	// Load env files if present
	if err := env.LoadEnvFileIfExists(); err != nil {
		log.Fatalf("load env: %v", err)
	}

	s := snapshot{
		APP_ENV:           env.GetAppEnv(),
		IsDev:             env.IsAppEnvDev(),
		IsLocal:           env.IsAppEnvLocal(),
		IsStaging:         env.IsAppEnvStaging(),
		IsProduction:      env.IsAppEnvProduction(),
		IsTesting:         env.IsAppEnvTesting(),
		IsTestingOrLocal:  env.IsAppEnvTestingOrLocal(),
		IsLocalOrStaging:  env.IsAppEnvLocalOrStaging(),
		Addr:              env.Get("ADDR", "127.0.0.1:8080"),
		Debug:             env.GetBool("DEBUG", "false"),
		Port:              env.GetInt("PORT", "3000"),
		Limit:             env.GetInt64("LIMIT", "1024"),
		Workers:           env.GetUint("WORKERS", "4"),
		MaxItems:          env.GetUint64("MAX_ITEMS", "5000"),
		Threshold:         env.GetFloat("THRESHOLD", "0.75"),
		Timeout:           env.GetDuration("HTTP_TIMEOUT", "5s"),
		Peers:             env.GetSlice("PEERS", "10.0.0.1,10.0.0.2"),
		Limits:            env.GetMap("LIMITS", "read=10,write=5"),
		AppEnvEnum:        env.GetEnum("APP_ENV", "dev", []string{"dev", "staging", "prod"}),
		OS:                env.OS(),
		Arch:              env.Arch(),
		IsLinux:           env.IsLinux(),
		IsMac:             env.IsMac(),
		IsWindows:         env.IsWindows(),
		IsBSD:             env.IsBSD(),
		IsUnix:            env.IsUnix(),
		IsContainerOS:     env.IsContainerOS(),
		IsDocker:          env.IsDocker(),
		IsDockerInDocker:  env.IsDockerInDocker(),
		IsDockerHost:      env.IsDockerHost(),
		IsContainer:       env.IsContainer(),
		IsKubernetes:      env.IsKubernetes(),
		IsHostEnvironment: env.IsHostEnvironment(),
	}

	env.Dump(s)

	// Example dump output:
	// #main.snapshot {
	//  +APP_ENV          => "local" #string
	//  +IsDev            => false #bool
	//  +IsLocal          => true #bool
	//  +IsStaging        => false #bool
	//  +IsProduction     => false #bool
	//  +IsTesting        => false #bool
	//  +IsTestingOrLocal => true #bool
	//  +IsLocalOrStaging => true #bool
	//  +Addr             => "127.0.0.1:8080" #string
	//  +Debug            => false #bool
	//  +Port             => 3000 #int
	//  +Limit            => 1024 #int64
	//  +Workers          => 4 #uint
	//  +MaxItems         => 5000 #uint64
	//  +Threshold        => 0.750000 #float64
	//  +Timeout          => 5s #time.Duration
	//  +Peers            => #[]string [
	//    0 => "10.0.0.1" #string
	//    1 => "10.0.0.2" #string
	//  ]
	//  +Limits => #map[string]string {
	//     read => "10" #string
	//     write => "5" #string
	//  }
	//  +AppEnvEnum        => "dev" #string
	//  +OS                => "darwin" #string
	//  +Arch              => "arm64" #string
	//  +IsLinux           => false #bool
	//  +IsMac             => true #bool
	//  +IsWindows         => false #bool
	//  +IsBSD             => false #bool
	//  +IsUnix            => true #bool
	//  +IsContainerOS     => false #bool
	//  +IsDocker          => false #bool
	//  +IsDockerInDocker  => false #bool
	//  +IsDockerHost      => false #bool
	//  +IsContainer       => false #bool
	//  +IsKubernetes      => false #bool
	//  +IsHostEnvironment => true #bool
	// }

	time.Sleep(10 * time.Millisecond)
}
