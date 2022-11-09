package web

import (
	"k8s.io/klog/v2"
	"strings"
)

const (
	devopsTag    = "+devops"
	devopsAppTag = "app"
	//devopsDeployTag  = "deploy"
	//devopsCommandTag = "command"
)

type DevOpsOption struct {
	Enable bool
	Server []string
	//Deploy  *string
	//Command bool
}

func ExtractDevOpsOption(lines []string) *DevOpsOption {
	out := &DevOpsOption{
		Enable: false,
		//Command: false,
	}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, devopsTag) {
			continue
		}
		out.Enable = true
		comments := strings.Split(strings.TrimSpace(line[len(devopsTag):]), " ")
		for _, comment := range comments {
			comment = strings.TrimSpace(comment)
			kv := strings.Split(comment, "=")
			if len(kv) > 2 {
				klog.Warningf("Invalid comment format: %s", kv)
			}
			switch kv[0] {
			case devopsAppTag:
				if len(kv) == 1 {
					continue
				}
				for _, s := range strings.Split(kv[1], ",") {
					out.Server = append(out.Server, strings.TrimSpace(s))
				}
				//case devopsDeployTag:
				//	if len(kv) == 1 {
				//		continue
				//	}
				//	out.Deploy = &kv[1]
				//case devopsCommandTag:
				//	if len(kv) == 1 {
				//		out.Command = true
				//		continue
				//	}
				//	if kv[1] == "true" {
				//		out.Command = true
				//	}
			}
		}
		break
	}
	return out
}
