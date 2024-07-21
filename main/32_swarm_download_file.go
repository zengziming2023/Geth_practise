package main

import (
	bzzclient "github.com/ethersphere/swarm/api/client"
	"github.com/golang/glog"
	"io"
	"main/main/util"
)

func main() {
	util.GlogInit()
	defer glog.Flush()

	client := bzzclient.NewClient("http://127.0.0.1:8500")
	manifestHash := "437553621d2c6839e8c966037d9677f8b13de411729f56fad5cbd6b2286c8e60"
	manifest, isEncrypted, err := client.DownloadManifest(manifestHash)
	if err != nil {
		glog.Fatal("Failed to download manifest: ", err)
	}
	glog.Info("isEncrypted: ", isEncrypted)
	for _, entry := range manifest.Entries {
		glog.Info("hash: ", entry.Hash)
		glog.Info("contentType: ", entry.ContentType)
		glog.Info("size: ", entry.Size)
		glog.Info("path: ", entry.Path)
	}

	file, err := client.Download(manifestHash, "")
	if err != nil {
		glog.Fatal("Failed to download manifest: ", err)
	}

	contentBytes, err := io.ReadAll(file)
	if err != nil {
		glog.Fatal("Failed to read file: ", err)
	}

	glog.Info("content: ", string(contentBytes))
}
