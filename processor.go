package main

// Libraries needed:
// fmt: to print output / statements.
// ioutil: file operations like open file, read file and close file.
// bufio: read input from standard input using the Scanner class.
// os: used to read command line arguments.
// encoding/json: used to marshal / unmarshal json input as needed.


import (  
	"fmt"
	"io/ioutil"
	"bufio"
	"os"
	"encoding/json"
)

// A structure that defines the a Push Down Automata
// This class contains the attributes of a PDA
type Pda struct {
	Name string
	States [] string
	Input_alphabet [] string
	Stack_alphabet [] string
	Accepting_states [] string
	Start_state string
	Transitions [][]string
	Eos string
}

// This class simulates a PDA processor that is it runs the PDA for teh provided input.
type PDAProcessor struct{
	Stack [] string
	Pda Pda
	Current_State string
}

// Function to push data on to the stack when executing the PDA. It modifies the stack.
func push(p *PDAProcessor, val string) {
	p.Stack = append(p.Stack, val)
}

// Function to pop data from the stack when executing the PDA. It modifies the stack.
func pop(p *PDAProcessor) {
	p.Stack = p.Stack[:len(p.Stack) -1]
}

// Function to obtain the top n elements of the stack. This function does not modify the stack.
func peek(p *PDAProcessor, k int) []string {
	top := [] string{}
	l := len(p.Stack)
	if (l <= k) {
		top = p.Stack
	} else if ( k == 1) {
		top = append(top, p.Stack[l-1])
	} else {
		top = p.Stack[l-k:l-1]
	}
	return top
}

// Function to reset the PDA and the stack. This deletes everything from the stack 
// and sets the current state to the start state so that we can start new.
func reset(p *PDAProcessor) {
	p.Stack = make([]string, 0)
	p.Current_State = p.Pda.Start_state
}

// Function to open the grammar file. This function unmarshal's the data input into the PDA structure.
func open(data []byte, p *Pda) bool {	
	err := json.Unmarshal(data, &p)
	if err != nil {
		fmt.Print(err)
		return false
	}
	return true
}

// Function to check if the input string has been accepted by the pda 
func is_accepted(proc PDAProcessor) bool{
	flag := false
	accepting_states := proc.Pda.Accepting_states
	cs := proc.Current_State

	if len(proc.Stack) == 0 {
		for i:= 0; i < len(accepting_states); i++ {
			if cs == accepting_states[i] {
				flag = true
				fmt.Println("Input token Accepted")
				break
			}
		}
	}
	if !flag {
		fmt.Println("Input token Rejected")
	}
	return flag
}

// The done returns the final status of the current state and the stack after the input string is processed.
func done(proc PDAProcessor, is_accepted bool, transition_count int){
	fmt.Println("")
	fmt.Println("Pda = ", proc.Pda.Name,":method=is_accepted:: ", is_accepted)
	fmt.Println("Current State: ", current_state(proc))
	var top = peek(&proc, 5)
	if (len(top) >= 5){
		fmt.Println("Stack Symbols", top)
	}
	fmt.Println("Number of transitions = ", transition_count)
}
// Returns the current state of the PDA
func current_state(proc PDAProcessor) string{
	return proc.Current_State
}


// This function accepts the input string and performs the necessary transitions and 
// stack operations for every token,
func put(proc *PDAProcessor, s string) int {
	
	var p Pda = proc.Pda
	transitions := p.Transitions
	tran_len := len(transitions)
	transition_count := 0
	for j := 0; j < tran_len; j++ {
		var allowed_current_state = transitions[j][0]
		var input = transitions[j][1]
		var allowed_top_of_stack = transitions[j][2]
		var target_state = transitions[j][3]
		var action_item = transitions[j][4]
		var currentStackSymbol = ""
		var top = peek(proc, 1)
		if(len(top)>=1){
			currentStackSymbol = top[0]
		}
		
		// PDA is deterministic. It jumps from current state to target state in the specified conditions
		if(input == "null" && allowed_current_state == proc.Current_State && allowed_top_of_stack == "null" && action_item == "null") {
			fmt.Println("No push/pop performed...... Processed dead transition")
			fmt.Println("Old Current State ",proc.Current_State)
			fmt.Println("New Current State ",target_state)
			proc.Current_State = target_state
			transition_count = transition_count + 1
		}
		if allowed_current_state == proc.Current_State && input == s {
			//Perform Push action
			if action_item != "null" {
				fmt.Println("push ", action_item, " on the stack")
				fmt.Println("Old Current State ",proc.Current_State)
				fmt.Println("New Current State ",target_state)
				transition_count = transition_count + 1
				proc.Current_State = target_state
				push(proc, action_item)
				break
				//performs Pop action
			} else if (allowed_top_of_stack == currentStackSymbol){
				pop(proc)
				fmt.Println("pop top of the stack")
				fmt.Println("Old Current State ",proc.Current_State)
				fmt.Println("New Current State ",target_state)
				transition_count = transition_count + 1
				proc.Current_State = target_state
				break
				//Neither push nor pop action required
			} else if(allowed_top_of_stack == "null"){
				fmt.Println("No push/pop performed...... Consumed input token")
				fmt.Println("Old Current State ",proc.Current_State)
				fmt.Println("New Current State ",target_state)
				proc.Current_State = target_state
				transition_count = transition_count + 1
				break
			}
		}	       
	}

	return transition_count
}

// Performs the last transition to move the Automata to accepting state after the input
// string has been successfully parsed. 
func eos(proc *PDAProcessor) {
	length_of_stack := len(proc.Stack)
	allowed_transitions := proc.Pda.Transitions
	target_state := ""
	allowed_top_of_stack := ""
	var currentStackSymbol = ""
	var top = peek(proc, 1)
	if(len(top)>=1){
		currentStackSymbol = top[0]
	}
	for j := 0; j < len(allowed_transitions); j++ {	
		var allowed_current_state = allowed_transitions[j][0]
		allowed_top_of_stack = allowed_transitions[j][2]
		
		if allowed_current_state == proc.Current_State && allowed_top_of_stack == currentStackSymbol{
			target_state = allowed_transitions[j][3]
			break
		}
	}
	if currentStackSymbol == proc.Pda.Eos {
		fmt.Println("")
		fmt.Println("Popping last $ from the stack")
		fmt.Println("Old Current State ",proc.Current_State)
		fmt.Println("New Current State ",target_state)
		proc.Current_State = target_state
		if length_of_stack > 0 {
			pop(proc)
		}
	}
}

// Pushes initial EOS token into the stack and moves to the next state indicating
// the start of transitions
func check_for_first_move(proc *PDAProcessor, transition_count int) int{
	allowed_transitions := proc.Pda.Transitions
	target_state := ""
	input := ""
	allowed_top_of_stack := ""
	action_item := ""
	
	for j := 0; j < len(allowed_transitions); j++ {
		var allowed_current_state = allowed_transitions[j][0]
		if allowed_current_state == proc.Current_State {
			input = allowed_transitions[j][1]
			allowed_top_of_stack = allowed_transitions[j][2]
			target_state = allowed_transitions[j][3]
			action_item = allowed_transitions[j][4]
			break
		}
	}
	
	if  input == "null" && allowed_top_of_stack == "null"{
		fmt.Println("Pushing first $ in the stack")
		fmt.Println("Old Current State ",proc.Current_State)
		fmt.Println("New Current State ",target_state)
        proc.Current_State = target_state
        push(proc, action_item)
        transition_count = transition_count + 1
	} 
	return transition_count
}


// main function to start input processing. It reads the grammar json file
// to create an instance of PDA. The input is obtained either from a text file or
// standard input.
func main(){

	fn := os.Args[1]

	var p Pda
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		fmt.Print(err)
    }
	open(data, &p)

	proc := PDAProcessor {
		Pda: p,
	}
	reset(&proc)

	transition_count := 0
	input_string := ""

	if len(os.Args) < 3 {

		fmt.Print("Enter input string: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input_string = scanner.Text()
	} else {
		inp := os.Args[2]
		file, err := os.Open(inp)

		input := make([]byte, 32)
		n, err := file.Read(input)
		if err != nil {
			fmt.Print(err)
		}
		input_string = string(input[:n])
	}
	fmt.Println("")
	isAccepted := false
	if input_string != "" {
		transition_count = check_for_first_move(&proc, transition_count)
		inp_len := len(input_string)
		i := 0
		for ; i < inp_len; i++ {
			count:=0
			char := string(input_string[i])
			fmt.Println(" ")
			fmt.Println("************Consuming Input Token***************", char)
			count = put(&proc, char)
			if  count == 0 {
				fmt.Println("Input token not processed ", char)
				break
			} else {
				transition_count = transition_count + count
			}
		}
		//End of input string reached
		if i == inp_len {
			eos(&proc)			
		} 
	}
	isAccepted = is_accepted(proc)
	done(proc, isAccepted, transition_count)
}