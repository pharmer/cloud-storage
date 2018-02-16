## ext-storage init

Initializes the driver.

### Synopsis


Initializes the driver.

```
ext-storage init
```

### Options

```
      --kubeconfig string    Absolute path to the kubeconfig file. Either this or master needs to be set if the provisioner is being run out of cluster.
      --master string        Master URL to build a client config from. Either this or kubeconfig needs to be set if the provisioner is being run out of cluster.
      --provisioner string   Name of the provisioner. The provisioner will only provision volumes for claims that request a StorageClass with a provisioner field set equal to this name. (default "external/pharmer")
```

### Options inherited from parent commands

```
      --alsologtostderr                  log to standard error as well as files
      --analytics                        Send analytical events to Google Analytics (default true)
      --log_backtrace_at traceLocation   when logging hits line file:N, emit a stack trace (default :0)
      --log_dir string                   If non-empty, write log files in this directory
      --logtostderr                      log to standard error instead of files
      --stderrthreshold severity         logs at or above this threshold go to stderr (default 2)
  -v, --v Level                          log level for V logs
      --vmodule moduleSpec               comma-separated list of pattern=N settings for file-filtered logging
```

### SEE ALSO
* [ext-storage](ext-storage.md)	 - Pharm external storage by Appscode - Start farms

