#!/usr/bin/env bash

go install covid/vendor/k8s.io/code-generator/cmd/openapi-gen
go install covid/vendor/k8s.io/code-generator/cmd/deepcopy-gen
go install covid/vendor/k8s.io/code-generator/cmd/conversion-gen
go install covid/vendor/k8s.io/code-generator/cmd/defaulter-gen
go install covid/vendor/k8s.io/code-generator/cmd/client-gen
go install covid/vendor/k8s.io/code-generator/cmd/lister-gen
go install covid/vendor/k8s.io/code-generator/cmd/informer-gen