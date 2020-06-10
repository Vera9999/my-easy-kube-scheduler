package controller

import (
	"log"
	"strings"
	"math/rand"
	"k8s.io/api/core/v1"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/api"
)

const (
	LuckyPred        = "Lucky"
	LuckyPredFailMsg = "Sorry, you're not lucky"
)

var predicatesFuncs = map[string]FitPredicate{
	LuckyPred: LuckyPredicate,
}

type FitPredicate func(pod *v1.Pod, node v1.Node) (bool, []string, error)

var predicatesSorted = []string{LuckyPred}

// filter 根据扩展程序定义的预选规则来过滤节点
func filter(args schedulerapi.ExtenderArgs) *schedulerapi.ExtenderFilterResult {
	var filteredNodes []v1.Node // 存放符合过滤规则的节点
	failedNodes := make(schedulerapi.FailedNodesMap) //用来记录被过滤掉的失败的节点
	pod := args.Pod //被调度的pod
	for _, node := range args.Nodes.Items {
		fits, failReasons, _ := podFitsOnNode(pod, node) // 判断这个node是否满足pod的过滤条件
		if fits {
			filteredNodes = append(filteredNodes, node) // 满足
		} else {
			failedNodes[node.Name] = strings.Join(failReasons, ",") // 不满足
		}
	}

	result := schedulerapi.ExtenderFilterResult{
		Nodes: &v1.NodeList{
			Items: filteredNodes,
		},
		FailedNodes: failedNodes,
		Error:       "",
	}

	return &result
}

// 判断node是否与pod相匹配
func podFitsOnNode(pod *v1.Pod, node v1.Node) (bool, []string, error) {
	fits := true
	var failReasons []string
	for _, predicateKey := range predicatesSorted { // 多次进行预选规则的判断
		fit, failures, err := predicatesFuncs[predicateKey](pod, node)
		if err != nil {
			return false, nil, err
		}
		fits = fits && fit
		failReasons = append(failReasons, failures...)
	}
	return fits, failReasons, nil
}

// 预选规则
func LuckyPredicate(pod *v1.Pod, node v1.Node) (bool, []string, error) {
	lucky := ( len(pod.Name) > 10) && ( rand.Intn(100) > 50) // 判断pod的名字是否大于10 且 随机数是否满足条件
	if lucky {// 满足条件
		log.Printf("pod %v/%v is lucky to fit on node %v\n", pod.Name, pod.Namespace, node.Name)
		return true, nil, nil
	}
	log.Printf("pod %v/%v is unlucky to fit on node %v\n", pod.Name, pod.Namespace, node.Name) // 不满足
	return false, []string{LuckyPredFailMsg}, nil
}
