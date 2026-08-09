# Parameterized Rotations and the Bloch Sphere

In quantum computing, the state of a single qubit can be visualized as a point on the surface of a unit sphere, known as the **Bloch Sphere**. While basic gates like Hadamard and Pauli-X provide fixed transformations, **Parameterized Rotation Gates** allow for arbitrary precision in positioning the qubit's state.

## 1. The Bloch Sphere Representation

Any single-qubit state $|\psi\rangle$ can be written as:
$$|\psi\rangle = \cos(\theta/2)|0\rangle + e^{i\phi}\sin(\theta/2)|1\rangle$$
where:
- $\theta$ is the angle from the Z-axis ($0 \leq \theta \leq \pi$).
- $\phi$ is the phase angle around the Z-axis ($0 \leq \phi < 2\pi$).

## 2. Rotation Gates (Rx, Ry, Rz)

quantum-go provides three primary rotation gates, each corresponding to a rotation around one of the principal axes of the Bloch Sphere.

### Rx: Rotation around the X-axis
$$Rx(\theta) = \begin{pmatrix} \cos(\theta/2) & -i\sin(\theta/2) \\ -i\sin(\theta/2) & \cos(\theta/2) \end{pmatrix}$$
- **Effect**: Rotates the state vector around the X-axis.
- **Identity**: $Rx(\pi)$ is equivalent to the Pauli-X gate (up to a global phase).

### Ry: Rotation around the Y-axis
$$Ry(\theta) = \begin{pmatrix} \cos(\theta/2) & -\sin(\theta/2) \\ \sin(\theta/2) & \cos(\theta/2) \end{pmatrix}$$
- **Effect**: Rotates the state vector around the Y-axis.
- **Example**: $Ry(\pi/2)$ moves the state from $|0\rangle$ to the equator ($|+\rangle$ state).

### Rz: Rotation around the Z-axis
$$Rz(\theta) = \begin{pmatrix} e^{-i\theta/2} & 0 \\ 0 & e^{i\theta/2} \end{pmatrix}$$
- **Effect**: Rotates the state vector around the Z-axis, changing the relative phase between $|0\rangle$ and $|1\rangle$.
- **Identity**: $Rz(\pi)$ is equivalent to the Pauli-Z gate (up to a global phase).

## 3. Practical Usage in quantum-go

You can use these gates to prepare specific states or as part of variational algorithms (like VQE).

```go
p := core.NewProgram(1)
// Rotate 45 degrees around the Y axis
p.AddStep(core.NewStep(core.NewRy(math.Pi/4, 0)))
```

## 4. Summary Table

| Gate | Axis | Matrix Effect | Captions |
| :--- | :--- | :--- | :--- |
| **Rx** | X | Bit-flip like rotation | `Rx` |
| **Ry** | Y | Superposition like rotation | `Ry` |
| **Rz** | Z | Phase-shift rotation | `Rz` |

---
**Next Step**: Explore [Universality and the U Gate](./universality.md) to see how these rotations combine into a single powerful operation.
