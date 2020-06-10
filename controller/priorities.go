package controller

import (
	"log"
	"math/rand"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/api"
)

const (
	luckyPrioMsg = "pod %v/%v is lucky to get score %v\n"
)


// 打分的分数 + 默认调度器给出的分数 = 最后的分数
// 对已过滤节点进行打分
func prioritize(args schedulerapi.ExtenderArgs) *schedulerapi.HostPriorityList {
	pod := args.Pod 
	nodes := args.Nodes.Items // 已过滤节点

	hostPriorityList := make(schedulerapi.HostPriorityList, len(nodes))
	for i, node := range nodes {
		// score := rand.Intn(schedulerapi.MaxPriority + 1)
                score := (len(pod.Name) + len(pod.Namespace) ) % rand.Intn(schedulerapi.MaxPriority) // 打分
		log.Printf(luckyPrioMsg, pod.Name, pod.Namespace, score)
		hostPriorityList[i] = schedulerapi.HostPriority{
			Host:  node.Name,
			Score: score,
		}
	}

	return &hostPriorityList
}
