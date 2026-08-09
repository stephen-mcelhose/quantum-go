// Package core defines the domain models for quantum programs, steps, and gates.
// It provides the foundational structures used to build and execute quantum circuits.
package core

import (
	"fmt"
	"math"

	smath "github.com/stephen-mcelhose/quantum-go/math"
)

// Gate describes an operation on one or more qubits.
type Gate interface {
	GetMatrix() smath.Matrix
	GetAffectedQubitIndexes() []int
	GetHighestAffectedQubitIndex() int
	GetCaption() string
	GetName() string
	GetGroup() string
	GetSize() int
	SetInverse(inv bool)
	IsInverse() bool
	HasOptimization() bool
	ApplyOptimize(v []complex128) []complex128
}

// BaseGate provides a common implementation for most gates.
// It implements the Gate interface with default behaviors that can be overridden
// by specific gate types.
type BaseGate struct {
	// AffectedQubits contains the indices of qubits affected by this gate.
	AffectedQubits []int
	// Inverse indicates whether this gate is inverted (adjoint/conjugate transpose).
	Inverse bool
	// Caption is a short display string for the gate (e.g., "H" for Hadamard).
	Caption string
	// Name is the full name of the gate (e.g., "Hadamard").
	Name string
	// Group is an optional grouping identifier for the gate.
	Group string
}

// GetAffectedQubitIndexes returns the indices of all qubits this gate operates on.
func (g *BaseGate) GetAffectedQubitIndexes() []int {
	return g.AffectedQubits
}

// GetHighestAffectedQubitIndex returns the highest qubit index affected by this gate.
// Returns -1 if no qubits are affected.
func (g *BaseGate) GetHighestAffectedQubitIndex() int {
	max := -1
	for _, idx := range g.AffectedQubits {
		if idx > max {
			max = idx
		}
	}
	return max
}

// GetCaption returns a short display string for this gate.
func (g *BaseGate) GetCaption() string {
	return g.Caption
}

// GetName returns the full name of this gate.
func (g *BaseGate) GetName() string {
	return g.Name
}

// GetGroup returns the grouping identifier for this gate.
func (g *BaseGate) GetGroup() string {
	return g.Group
}

// GetSize returns the number of qubits this gate operates on.
func (g *BaseGate) GetSize() int {
	return len(g.AffectedQubits)
}

// SetInverse sets whether this gate should be applied as its inverse (adjoint).
func (g *BaseGate) SetInverse(inv bool) {
	g.Inverse = inv
}

// IsInverse returns whether this gate is currently set to be applied as its inverse.
func (g *BaseGate) IsInverse() bool {
	return g.Inverse
}

// HasOptimization returns whether this gate has an optimized state vector application.
// The default implementation returns false.
func (g *BaseGate) HasOptimization() bool {
	return false
}

// ApplyOptimize applies this gate to a state vector using an optimized algorithm.
// The default implementation returns the input unchanged.
// Gates with optimization should override this method.
func (g *BaseGate) ApplyOptimize(v []complex128) []complex128 {
	return v
}

// Block represents a sequence of steps that can be treated as a single unit.
// A block can be used to create composite gates like QFT or arithmetic circuits,
// allowing complex operations to be encapsulated and reused.
type Block struct {
	// Steps contains the sequence of quantum operations in this block.
	Steps []*Step
	// NQubits is the number of qubits this block operates on.
	NQubits int
	// Name is a descriptive identifier for this block.
	Name string
}

// NewBlock creates a new Block with the specified name and number of qubits.
func NewBlock(name string, nqubits int) *Block {
	return &Block{
		Name:    name,
		NQubits: nqubits,
		Steps:   make([]*Step, 0),
	}
}

// StepExecutor is a function that calculates the next state vector given a list of gates.
// This function type is used to decouple the core definitions from the local simulator.
type StepExecutor func(gates []Gate, vector []complex128, numQubits int) []complex128

// GlobalStepExecutor is used to avoid circular dependencies between core and local packages.
// The local package sets this during initialization.
var GlobalStepExecutor StepExecutor

// AddStep appends a step to this block's sequence of operations.
func (b *Block) AddStep(s *Step) {
	b.Steps = append(b.Steps, s)
}

// ApplyOptimize applies all steps in this block to a state vector in sequence.
// If inverse is true, the steps are applied in reverse order with their inverse operations.
// Returns the transformed state vector.
func (b *Block) ApplyOptimize(v []complex128, inverse bool) []complex128 {
	if GlobalStepExecutor == nil {
		return v
	}
	res := make([]complex128, len(v))
	copy(res, v)

	if inverse {
		for i := len(b.Steps) - 1; i >= 0; i-- {
			s := b.Steps[i]
			s.SetInverse(true)
			res = GlobalStepExecutor(s.Gates, res, b.NQubits)
			s.SetInverse(false) // Restore
		}
	} else {
		for _, s := range b.Steps {
			res = GlobalStepExecutor(s.Gates, res, b.NQubits)
		}
	}
	return res
}

// Result describes the result of a quantum program execution.
// It contains both the full quantum state vector and derived qubit probabilities.
type Result interface {
	GetNumQubits() int
	GetProbability() []complex128
	GetQubits() []*Qubit
	PrintBinary()
	GetStateVectorReference() *StateVectorReference
}

// StateVectorReference represents the JSON schema for state vector exchange.
type StateVectorReference struct {
	NumQubits  int `json:"num_qubits"`
	Amplitudes []struct {
		Re float64 `json:"re"`
		Im float64 `json:"im"`
	} `json:"amplitudes"`
}

// InstrumentedResult extends Result with methods to access intermediate states.
type InstrumentedResult interface {
	Result
	GetIntermediateProbability(step int) []complex128
	GetIntermediateQubits(step int) []*Qubit
}

// CompactResult implements the Result interface for a single final state.
type CompactResult struct {
	// NumQubits is the number of qubits in the quantum system.
	NumQubits int
	// Probability is the full state vector representing the quantum state.
	// For n qubits, this has 2^n complex amplitudes.
	Probability []complex128
	// Qubits contains individual qubit measurement probabilities and values.
	Qubits []*Qubit
	// QubitStatesComputed tracks whether individual qubit states have been calculated.
	QubitStatesComputed bool
}

// GetNumQubits returns the number of qubits in the system.
func (r *CompactResult) GetNumQubits() int {
	return r.NumQubits
}

// GetProbability returns the full state vector.
func (r *CompactResult) GetProbability() []complex128 {
	return r.Probability
}

// GetQubits returns the individual qubit states, computing them from the state vector
// if they haven't been calculated yet. Each qubit contains its probability of being |1⟩.
func (r *CompactResult) GetQubits() []*Qubit {
	if !r.QubitStatesComputed {
		probs := CalculateQubitStatesFromVector(r.Probability)
		r.Qubits = make([]*Qubit, len(probs))
		for i, p := range probs {
			r.Qubits[i] = &Qubit{Probability: p}
		}
		r.QubitStatesComputed = true
	}
	return r.Qubits
}

// PrintBinary prints the states with non-zero probability in binary format.
// This is the preferred console output format for multi-qubit results.
func (r *CompactResult) PrintBinary() {
	fmt.Printf("Quantum Result (%d qubits):\n", r.NumQubits)
	for i, amp := range r.Probability {
		prob := real(amp)*real(amp) + imag(amp)*imag(amp)
		if prob > 0.0001 {
			fmt.Printf("|%0*b>: %.4f\n", r.NumQubits, i, prob)
		}
	}
}

// GetStateVectorReference returns a serializable reference to the state vector.
func (r *CompactResult) GetStateVectorReference() *StateVectorReference {
	ref := &StateVectorReference{
		NumQubits:  r.NumQubits,
		Amplitudes: make([]struct {
			Re float64 `json:"re"`
			Im float64 `json:"im"`
		}, len(r.Probability)),
	}
	for i, amp := range r.Probability {
		ref.Amplitudes[i].Re = real(amp)
		ref.Amplitudes[i].Im = imag(amp)
	}
	return ref
}

// InstrumentedResultImpl implements the InstrumentedResult interface.
type InstrumentedResultImpl struct {
	CompactResult
	// IntermediateProbability stores the state vector after each step.
	IntermediateProbability [][]complex128
	// IntermediateQubits stores the calculated qubit states after each step.
	IntermediateQubits map[int][]*Qubit
}

// GetIntermediateProbability returns the state vector after a specific step.
func (r *InstrumentedResultImpl) GetIntermediateProbability(step int) []complex128 {
	if step < 0 || step >= len(r.IntermediateProbability) {
		return nil
	}
	return r.IntermediateProbability[step]
}

// GetIntermediateQubits returns the individual qubit states after a specific step.
func (r *InstrumentedResultImpl) GetIntermediateQubits(step int) []*Qubit {
	if qubits, ok := r.IntermediateQubits[step]; ok {
		return qubits
	}
	probs := r.GetIntermediateProbability(step)
	if probs == nil {
		return nil
	}
	qubitProbs := CalculateQubitStatesFromVector(probs)
	qubits := make([]*Qubit, len(qubitProbs))
	for i, p := range qubitProbs {
		qubits[i] = &Qubit{Probability: p}
	}
	if r.IntermediateQubits == nil {
		r.IntermediateQubits = make(map[int][]*Qubit)
	}
	r.IntermediateQubits[step] = qubits
	return qubits
}

// Qubit represents a single qubit's state and measurement information.
type Qubit struct {
	// Probability is the probability of this qubit being in the |1⟩ state.
	Probability float64
	// Measured indicates whether this qubit has been measured.
	Measured bool
	// MeasuredValue is the measured value (false for |0⟩, true for |1⟩).
	MeasuredValue bool
}

// Measure returns the measured value of this qubit as an integer (0 or 1).
func (q *Qubit) Measure() int {
	if q.MeasuredValue {
		return 1
	}
	return 0
}

// CalculateQubitStatesFromVector computes the probability of each qubit being in the |1⟩ state.
// For each qubit, it sums the squared magnitudes of all amplitudes where that qubit's bit is 1.
// Returns a slice where element i is the probability of qubit i being |1⟩.
func CalculateQubitStatesFromVector(vector []complex128) []float64 {
	if len(vector) == 0 {
		return nil
	}
	nq := int(math.Round(math.Log2(float64(len(vector)))))
	answer := make([]float64, nq)
	for i := 0; i < nq; i++ {
		div := 1 << i
		for j := 0; j < len(vector); j++ {
			if (j/div)%2 == 1 {
				val := vector[j]
				answer[i] += real(val)*real(val) + imag(val)*imag(val)
			}
		}
	}
	return answer
}

// StepType defines the type of step in a quantum circuit.
type StepType int

const (
	// StepNormal represents a standard quantum gate operation.
	StepNormal StepType = iota
	// StepPseudo represents a pseudo-operation (e.g., for visualization).
	StepPseudo
	// StepProbability represents a probability measurement operation.
	StepProbability
)

// Step represents a single step in a quantum circuit.
// A step contains one or more gates that can be applied in parallel
// (gates operating on disjoint sets of qubits).
type Step struct {
	// Type specifies the kind of operation this step performs.
	Type StepType
	// Gates contains all gates applied in this step.
	Gates []Gate
	// Name is a descriptive identifier for this step.
	Name string
	// Program is a back-reference to the program containing this step.
	Program *Program
	// Informal indicates whether this is an informal step (not part of execution).
	Informal bool
}

// NewStep creates a new Step with the given gates.
// Gates are verified to operate on disjoint qubits.
func NewStep(gates ...Gate) *Step {
	s := &Step{
		Type:  StepNormal,
		Name:  "unknown",
		Gates: make([]Gate, 0),
	}
	for _, g := range gates {
		s.AddGate(g)
	}
	return s
}

// AddGate adds a gate to this step after verifying it doesn't conflict with existing gates.
// Panics if the gate operates on a qubit already used by another gate in this step.
func (s *Step) AddGate(g Gate) {
	s.verifyUnique(g)
	s.Gates = append(s.Gates, g)
}

// verifyUnique checks that the gate doesn't operate on any qubits already used in this step.
// This ensures all gates in a step operate on disjoint qubit sets and can be applied in parallel.
func (s *Step) verifyUnique(gate Gate) {
	affected := gate.GetAffectedQubitIndexes()
	for _, g := range s.Gates {
		existing := g.GetAffectedQubitIndexes()
		for _, a := range affected {
			for _, e := range existing {
				if a == e {
					panic(fmt.Sprintf("Adding gate that affects qubit %d already involved in this step", a))
				}
			}
		}
	}
}

// SetInverse sets the inverse flag for all gates in this step.
func (s *Step) SetInverse(inv bool) {
	for _, g := range s.Gates {
		g.SetInverse(inv)
	}
}

// Program represents a complete quantum circuit.
// It defines the number of qubits, initial states, and sequence of quantum operations.
type Program struct {
	// NumQubits is the number of qubits in this quantum circuit.
	NumQubits int
	// Steps contains the ordered sequence of quantum operations.
	Steps []*Step
	// InitAlpha specifies the initial probability amplitudes for each qubit.
	// A value of 1.0 means the qubit starts in the |0⟩ state.
	InitAlpha []float64
	// Result holds the execution result after the program is run.
	Result Result
}


// NewProgram creates a new quantum Program with the specified number of qubits.
// All qubits are initialized to the |0⟩ state.
// Additional steps can be provided to initialize the circuit.
func NewProgram(numQubits int, steps ...*Step) *Program {
	p := &Program{
		NumQubits: numQubits,
		Steps:     make([]*Step, 0),
		InitAlpha: make([]float64, numQubits),
	}
	for i := range p.InitAlpha {
		p.InitAlpha[i] = 1.0
	}
	p.AddSteps(steps...)
	return p
}

// AddStep appends a step to this program and sets the step's program reference.
func (p *Program) AddStep(s *Step) {
	s.Program = p
	p.Steps = append(p.Steps, s)
}

// AddSteps appends multiple steps to this program.
func (p *Program) AddSteps(steps ...*Step) {
	for _, s := range steps {
		p.AddStep(s)
	}
}
