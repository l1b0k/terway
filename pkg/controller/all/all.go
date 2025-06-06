/*
Copyright 2021 Terway Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package all register all controllers
package all

import (
	// register all controllers
	_ "github.com/AliyunContainerService/terway/pkg/controller/eni"
	_ "github.com/AliyunContainerService/terway/pkg/controller/multi-ip/node"
	_ "github.com/AliyunContainerService/terway/pkg/controller/multi-ip/pod"
	_ "github.com/AliyunContainerService/terway/pkg/controller/node"
	_ "github.com/AliyunContainerService/terway/pkg/controller/pod"
	_ "github.com/AliyunContainerService/terway/pkg/controller/pod-eni"
	_ "github.com/AliyunContainerService/terway/pkg/controller/pod-networking"
)
