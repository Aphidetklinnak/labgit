package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	NumResources  = 3  
	TotalResource = 10
)

type Process struct {
	Name       string
	Allocation [NumResources]int
	Need       [NumResources]int
	Max        [NumResources]int
}

var (
	available = [NumResources]int{TotalResource, TotalResource, TotalResource}
	processes = []Process{}
)

func main() {
	for {
		showProcess()
		command, err := getCommand()
		if err != nil {
			fmt.Println("เกิดข้อผิดพลาดในการอ่านข้อมูล:", err)
			return
		}
		switch strings.TrimSpace(command) {
		case "exit":
			fmt.Println("---Exit Progarm---")
			return
		case "new":
			fmt.Print("Enter name -> ")
			name, _ := getCommand()
			fmt.Print("Enter Max -> ")
			Max, _ := getCommand()
			MaxP := strings.Split(Max, ",")
			new_p(name, MaxP)
		case "call":
			fmt.Print("Enter process name -> ")
			name, _ := getCommand()
			fmt.Print("Enter number of resources (Ex => a,b,c) -> ")
			all, _ := getCommand()
			AllP := strings.Split(all, ",")
			found := false
			for l := range processes {
				if strings.TrimSpace(name) == processes[l].Name {
					found = true
					allocateResources(&processes[l], AllP)
					break
				}
			}
			if !found {
				fmt.Print("Cannot find the processes")
			}
		default:
			fmt.Printf("\n!!-Command Error-!!\n")
		}
	}
}

func showProcess() {
	fmt.Printf("\n")
	fmt.Println("______________________________________________________________________________________")
	fmt.Println("| Process |    Allocation   |       Need      |       Max       |      Available     |")
	fmt.Println("|         |  A  |  B  |  C  |  A  |  B  |  C  |  A  |  B  |  C  |   A  |   B  |   C  |")
	fmt.Println("|---------|-----------------|-----------------|-----------------|--------------------|")
	// ส่วน show โปรเซส
	for i := range processes {
		fmt.Printf("|   %s    ", processes[i].Name)
		for k := range processes[i].Allocation {
			fmt.Printf("|  %d  ", processes[i].Allocation[k])
		}
		for k := range processes[i].Need {
			fmt.Printf("|  %d  ", processes[i].Need[k])
		}
		for k := range processes[i].Max {
			fmt.Printf("|  %d  ", processes[i].Max[k])
		}
		if i == 0 {
			for k := range available {
				fmt.Printf("|  %d  ", available[k])
			}
		} else {
			fmt.Printf("|                    ")
		}
		fmt.Printf("|")
		fmt.Printf("\n")
	}
	// ส่วน show โปรเซส
	fmt.Print("\n[Command] -> ")
}

func getCommand() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	data, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	data = strings.Trim(data, "\n")
	return strings.TrimSpace(strings.ToLower(data)), nil
}

func addNeed(n *Process) {
	for i := range n.Need {
		n.Need[i] = n.Max[i] - n.Allocation[i]
	}
}

func new_p(name string, max []string) {
	// แปลง max เป็น int
	maxInt := make([]int, len(max))
	for i, str := range max {
		num, err := strconv.Atoi(str)
		if err != nil {
			fmt.Println("Error:maxInt", err)
			return
		}
		maxInt[i] = num
	}

	p := Process{
		Name:       name,
		Allocation: [NumResources]int{0, 0, 0},
		Max:        [NumResources]int{maxInt[0], maxInt[1], maxInt[2]},
	}

	addNeed(&p)
	processes = append(processes, p)
}

func isSafeState(p *Process, all []int) bool {
	canAllocate := true
	for j := 0; j < NumResources; j++ {
		if p.Need[j] > available[j] {
			canAllocate = false
			break
		}
		if all[j] > p.Need[j] {
			canAllocate = false
			break
		}
	}

	return canAllocate // Safe
}

func allocateResources(p *Process, all []string) {
	//แปลง all เป็น int
	allInt := make([]int, len(all))
	for i, str := range all {
		num, err := strconv.Atoi(str)
		if err != nil {
			fmt.Println("Error:allInt", err)
			return
		}
		allInt[i] = num
	}

	if !isSafeState(p, allInt) {
		fmt.Printf("\n!!-Error (It's not safe!!!)-!!\n")
		return
	}

	for i := 0; i < NumResources; i++ {
		available[i] -= allInt[i]
		p.Allocation[i] += allInt[i]
	}
	addNeed(p)
	if releaseResources(p) {
		fmt.Printf("\nEnd process\n")
		for i := 0; i < NumResources; i++ {
			available[i] += p.Allocation[i]
		}
		removeProcessByName(p.Name)
	}
}

func releaseResources(p *Process) bool {
	for i := 0; i < NumResources; i++ {
		if p.Need[i] == 0 && p.Allocation == p.Max {
			continue
		} else {
			return false
		}
	}
	return true
}

func removeProcessByName(nameToRemove string) {
	for i, p := range processes {
		if p.Name == nameToRemove {
			processes = append(processes[:i], processes[i+1:]...)
			return
		}
	}
}
