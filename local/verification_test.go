package local_test

import (
	"math"
	"math/cmplx"
	"testing"

	"github.com/stephen-mcelhose/quantum-go/core"
	"github.com/stephen-mcelhose/quantum-go/local"
	smath "github.com/stephen-mcelhose/quantum-go/math"
)

const tolerance = 1e-6

func compareStateVectors(t *testing.T, name string, got, want []complex128) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("%s: dimension mismatch: got %d, want %d", name, len(got), len(want))
		return
	}

	for i := range got {
		diff := cmplx.Abs(got[i] - want[i])
		if diff > tolerance {
			t.Errorf("%s: value mismatch at index %d: got %v, want %v (diff %v)", name, i, got[i], want[i], diff)
		}
	}
}

func TestVerifyStandardStates(t *testing.T) {
	e := local.NewSimpleExecutionEnvironment()

	tests := []struct {
		name     string
		program  *core.Program
		expected []complex128
	}{
		{
			name: "Bell State",
			program: func() *core.Program {
				p := core.NewProgram(2)
				p.AddStep(core.NewStep(core.NewHadamard(0)))
				p.AddStep(core.NewStep(core.NewCnot(0, 1)))
				return p
			}(),
			expected: []complex128{
				complex(1.0/math.Sqrt(2), 0),
				0,
				0,
				complex(1.0/math.Sqrt(2), 0),
			},
		},
		{
			name: "GHZ State (3 qubits)",
			program: func() *core.Program {
				p := core.NewProgram(3)
				p.AddStep(core.NewStep(core.NewHadamard(0)))
				p.AddStep(core.NewStep(core.NewCnot(0, 1)))
				p.AddStep(core.NewStep(core.NewCnot(1, 2)))
				return p
			}(),
			expected: []complex128{
				complex(1.0/math.Sqrt(2), 0),
				0, 0, 0, 0, 0, 0,
				complex(1.0/math.Sqrt(2), 0),
			},
		},
		{
			name: "QFT on |00>",
			program: func() *core.Program {
				p := core.NewProgram(2)
				p.AddStep(core.NewStep(core.NewFourier(2, 0)))
				return p
			}(),
			expected: []complex128{
				0.5, 0.5, 0.5, 0.5,
			},
		},
			{
				name: "Time Evolution exp(-i*X*pi/4)",
				program: func() *core.Program {
					p := core.NewProgram(1)
					hx := smath.NewMatrix(2, 2)
					hx.Data = []complex128{0, 1, 1, 0}
					p.AddStep(core.NewStep(core.NewTimeEvolution(0, hx, math.Pi/4.0)))
					return p
				}(),
				expected: []complex128{
					complex(1.0/math.Sqrt(2), 0),
					complex(0, -1.0/math.Sqrt(2)),
				},
			},
			{
				name: "Rx(pi) equivalent to X (up to phase)",
				program: func() *core.Program {
					p := core.NewProgram(1)
					p.AddStep(core.NewStep(core.NewRx(math.Pi, 0)))
					return p
				}(),
				expected: []complex128{
					0,
					complex(0, -1),
				},
			},
			{
				name: "Ry(pi) equivalent to Y (up to phase)",
				program: func() *core.Program {
					p := core.NewProgram(1)
					p.AddStep(core.NewStep(core.NewRy(math.Pi, 0)))
					return p
				}(),
				expected: []complex128{
					0,
					1,
				},
			},
			{
				name: "V * V = X",
				program: func() *core.Program {
					p := core.NewProgram(1)
					p.AddStep(core.NewStep(core.NewV(0)))
					p.AddStep(core.NewStep(core.NewV(0)))
					return p
				}(),
				expected: []complex128{
					0,
					1,
				},
			},
			{
				name: "S Gate transition",
				program: func() *core.Program {
					p := core.NewProgram(1)
					p.AddStep(core.NewStep(core.NewHadamard(0)))
					p.AddStep(core.NewStep(core.NewS(0)))
					return p
				}(),
				expected: []complex128{
					complex(1.0/math.Sqrt(2), 0),
					complex(0, 1.0/math.Sqrt(2)),
				},
			},
			{
				name: "2-Qubit Rx(pi)",
				program: func() *core.Program {
					p := core.NewProgram(2)
					p.AddStep(core.NewStep(core.NewRx(math.Pi, 0), core.NewRx(math.Pi, 1)))
					return p
				}(),
				expected: []complex128{
					0, 0, 0, -1,
				},
			},
			{
				name: "2-Qubit Ry(pi)",
				program: func() *core.Program {
					p := core.NewProgram(2)
					p.AddStep(core.NewStep(core.NewRy(math.Pi, 0), core.NewRy(math.Pi, 1)))
					return p
				}(),
				expected: []complex128{
					0, 0, 0, 1,
				},
			},
			{
				name: "2-Qubit S gate",
				program: func() *core.Program {
					p := core.NewProgram(2)
					p.AddStep(core.NewStep(core.NewHadamard(0), core.NewHadamard(1)))
					p.AddStep(core.NewStep(core.NewS(0), core.NewS(1)))
					return p
				}(),
					expected: []complex128{
						0.5, complex(0, 0.5), complex(0, 0.5), -0.5,
					},
				},
				{
					name: "Toffoli Gate (110 -> 111)",
					program: func() *core.Program {
						p := core.NewProgram(3)
						p.AddStep(core.NewStep(core.NewX(0), core.NewX(1)))
						p.AddStep(core.NewStep(core.NewToffoli(0, 1, 2)))
						return p
					}(),
					expected: []complex128{
						0, 0, 0, 0, 0, 0, 0, 1,
					},
				},
				{
					name: "Fredkin Gate (110 -> 101)",
					program: func() *core.Program {
						p := core.NewProgram(3)
						p.AddStep(core.NewStep(core.NewX(0), core.NewX(1)))
						p.AddStep(core.NewStep(core.NewFredkin(0, 1, 2)))
						return p
					}(),
					expected: []complex128{
						0, 0, 0, 0, 0, 1, 0, 0,
					},
				},
				{
					name: "Superdense Coding (11)",
					program: func() *core.Program {
						p := core.NewSuperdenseCodingProgram()
						return p
					}(),
					expected: []complex128{
						0, 0, 0, 1,
					},
				},
			}



	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := e.RunProgram(tt.program)
			compareStateVectors(t, tt.name, res.GetProbability(), tt.expected)
		})
	}
}
