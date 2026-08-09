package qasm

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/stephen-mcelhose/quantum-go/core"
)

// Parse converts an OpenQASM 2.0 string into a core.Program.
// This is a lightweight parser that supports the subset of QASM used by Strange.
func Parse(input string) (*core.Program, error) {
	scanner := bufio.NewScanner(strings.NewReader(input))
	var p *core.Program
	var numQubits int

	// Regex for basic declarations and gates
	qregRegex := regexp.MustCompile(`qreg\s+q\[(\d+)\]`)
	gateRegex := regexp.MustCompile(`([a-z0-9]+)(\((.*)\))?\s+q\[(\d+)\](,\s*q\[(\d+)\])?(,\s*q\[(\d+)\])?`)

	for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "OPENQASM") || strings.HasPrefix(line, "include") || strings.HasPrefix(line, "barrier") {
				continue
			}


		// Split by semicolon to handle multiple gates on one line
		segments := strings.Split(line, ";")
		for _, segment := range segments {
			segment = strings.TrimSpace(segment)
			if segment == "" {
				continue
			}

			// Handle qreg
			if matches := qregRegex.FindStringSubmatch(segment); matches != nil {
				n, _ := strconv.Atoi(matches[1])
				numQubits = n
				p = core.NewProgram(numQubits)
				continue
			}

			if p == nil {
				continue // Wait for qreg
			}

			// Handle measure specially
			if strings.HasPrefix(segment, "measure") {
				measureRegex := regexp.MustCompile(`measure\s+q\[(\d+)\]\s*->\s*c\[(\d+)\]`)
				if matches := measureRegex.FindStringSubmatch(segment); matches != nil {
					qIdx, _ := strconv.Atoi(matches[1])
					p.AddStep(core.NewStep(core.NewMeasurement(qIdx)))
					continue
				}
			}

			// Handle gates
			if matches := gateRegex.FindStringSubmatch(segment); matches != nil {
				gateName := matches[1]

			params := matches[3]
			q0, _ := strconv.Atoi(matches[4])
			
			var g core.Gate
				switch gateName {
					case "h":
						g = core.NewHadamard(q0)
					case "x":
						g = core.NewX(q0)
					case "y":
						g = core.NewY(q0)
					case "z":
						g = core.NewZ(q0)
						case "s":
							g = core.NewS(q0)
						case "sdg":
							g = core.NewS(q0)
							g.SetInverse(true)
						case "t":
							g = core.NewT(q0)
						case "tdg":
							g = core.NewT(q0)
							g.SetInverse(true)

					case "sx":
						g = core.NewV(q0)
					case "rx":
						theta, _ := strconv.ParseFloat(params, 64)
						g = core.NewRx(theta, q0)
					case "ry":
						theta, _ := strconv.ParseFloat(params, 64)
						g = core.NewRy(theta, q0)
					case "rz":
						theta, _ := strconv.ParseFloat(params, 64)
						g = core.NewRz(theta, q0)
					case "id":
						g = core.NewIdentity(q0)
					case "cx":
						q1, _ := strconv.Atoi(matches[6])
						g = core.NewCnot(q0, q1)
					case "cz":
						q1, _ := strconv.Atoi(matches[6])
						g = core.NewCz(q0, q1)
					case "swap":
						q1, _ := strconv.Atoi(matches[6])
						g = core.NewSwap(q0, q1)
					case "ccx":
						q1, _ := strconv.Atoi(matches[6])
						q2, _ := strconv.Atoi(matches[8])
						g = core.NewToffoli(q0, q1, q2)
					case "u1":
						theta, _ := strconv.ParseFloat(params, 64)
						g = core.NewPhaseShift(theta, q0)
					case "u3", "u":
						pList := strings.Split(params, ",")
						theta, _ := strconv.ParseFloat(strings.TrimSpace(pList[0]), 64)
						phi := 0.0
						lambda := 0.0
						if len(pList) > 1 {
							phi, _ = strconv.ParseFloat(strings.TrimSpace(pList[1]), 64)
						}
						if len(pList) > 2 {
							lambda, _ = strconv.ParseFloat(strings.TrimSpace(pList[2]), 64)
						}
						g = core.NewU(theta, phi, lambda, q0)
					case "cu1":
						theta, _ := strconv.ParseFloat(params, 64)
						q1, _ := strconv.Atoi(matches[6])
						g = core.NewCr(q0, q1, theta)
					case "measure":
						g = core.NewMeasurement(q0)

				default:
					fmt.Printf("Warning: Skipping unknown gate: %s\n", gateName)
					continue
				}
				p.AddStep(core.NewStep(g))
			}
		}
	}


	if p == nil {

		return nil, fmt.Errorf("no qreg found in QASM input")
	}

	return p, nil
}
