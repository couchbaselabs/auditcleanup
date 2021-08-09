/*
 *  Copyright 2021 Couchbase, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file  except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the  License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"os"
	"strings"

	"github.com/couchbase/audit-cleanup/pkg/logging"
	"github.com/couchbase/audit-cleanup/pkg/version"
	"github.com/fsnotify/fsnotify"
)

var (
	log = logging.Log
)

// Create a filesystem watcher that removes rotated audit logs as soon as they are created.
// Event-driven rather than polling.

const (
	LogDirEnvVar = "AUDIT_LOG_DIR"
)

func main() {
	log.Infow("Starting up Couchbase Audit Cleanup", "version", version.WithBuildNumber(), "revision", version.GitRevision(), "environment", os.Environ())

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalw("Unable to create watcher", "error", err)
	}
	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Op&fsnotify.Create != fsnotify.Create {
					log.Debugw("Skipping event as not file creation", "event", event)

					return
				}

				log.Debugw("File creation event", "event", event, "file", event.Name)

				if strings.HasSuffix(event.Name, "-audit.log") {
					err := os.Remove(event.Name)
					if err != nil {
						log.Errorw("Unable to remove file", "error", err, "event", event, "file", event.Name)
					} else {
						log.Infow("Deleted file successfully", "file", event.Name)
					}
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}

				log.Errorw("Error during watch", "error", err)
			}
		}
	}()

	directoryToWatch := os.Getenv(LogDirEnvVar)
	if directoryToWatch == "" {
		directoryToWatch = "/opt/couchbase/var/lib/couchbase/logs/"
	}

	log.Infow("Monitoring directory, override with "+LogDirEnvVar, "dir", directoryToWatch)

	err = watcher.Add(directoryToWatch)
	if err != nil {
		log.Fatalw("Unable to add directory to watcher", "error", err, "dir", directoryToWatch)
	}

	<-done

	log.Info("Exiting Couchbase Audit Cleanup")
}
