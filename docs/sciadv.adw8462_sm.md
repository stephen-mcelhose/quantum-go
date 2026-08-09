Below is the conversion of the provided source material into Markdown format with LaTeX for the mathematical expressions.

## Derivation of the Generalized First Law of Thermodynamics

In this section, the generalized first law of thermodynamics (Eq. 1 of the main text) is derived:

$$W = \sum_{j} Q_{j} + \Delta U \tag{S1}$$

Where the components are defined as follows:
*   **$W = \int_{t}^{t+\tau} dt' \langle \partial_{t'} H_{tot}(t') \rangle$**: The total work extracted during a cycle of duration $\tau$.
*   **$Q_{j} = \langle H_{j} \rangle (t+\tau) - \langle H_{j} \rangle (t)$**: The heat absorbed by reservoir $R_j$.
*   **$\Delta U = \sum_{i} \Delta U_i = \sum_{i} [\langle H_{i} + \sum_{j} H_{ij} \rangle (t+\tau) - \langle H_{i} + \sum_{j} H_{ij} \rangle (t)]$**: The change in energy of the compound system, including the **interaction energy**.

The **total Hamiltonian** is given by:
$$H_{tot}(t) = \sum_{i} H_i(t) + \sum_{j} H_j + \sum_{i,j} H_{ij}(t) \tag{S2}$$
where $H_i(t)$ is the Hamiltonian of subsystem $S_i$, $H_j$ is the Hamiltonian of reservoir $R_j$, and $H_{ij}(t)$ is the interaction Hamiltonian between $S_i$ and $R_j$.

Using the **generalized version of the Ehrenfest theorem**:
$$\frac{d}{dt} \langle H_{tot} \rangle = \langle \partial_t H_{tot}(t) \rangle \tag{S3}$$

By the linearity of differentiation:
$$\frac{d}{dt} \langle H_{tot} \rangle = \sum_{i} \frac{d}{dt} \langle H_i \rangle + \sum_{i,j} \frac{d}{dt} \langle H_{ij} \rangle + \sum_{j} \frac{d}{dt} \langle H_j \rangle \tag{S4}$$

Combining Eqs. (S3) and (S4):
$$\langle \partial_t H_{tot}(t) \rangle = \sum_{i} \frac{d}{dt} \langle H_i + \sum_{j} H_{ij} \rangle + \sum_{j} \frac{d}{dt} \langle H_j \rangle \tag{S5}$$

Integrating with respect to $t$ yields the **generalized first law**:
$$\underbrace{\int_{t}^{t+\tau} dt' \langle \partial_{t'} H_{tot}(t') \rangle}_{W} = \underbrace{\sum_{i} [\langle H_i + \sum_{j} H_{ij} \rangle (t+\tau) - \langle H_i + \sum_{j} H_{ij} \rangle (t)]}_{\Delta U} + \underbrace{\sum_{j} [\langle H_j \rangle (t+\tau) - \langle H_j \rangle (t)]}_{Q_j} \tag{S6}$$

If $\Delta U_S = 0$ and the interactions are **energy-conserving on average**, the standard first law is recovered: $W_S = \sum_j Q_j$.

---

## Derivation of the Generalized Second Law of Thermodynamics

The generalized second law (Eq. 2 of the main text) is expressed as:
$$\sum_{i} \Delta S(\rho_i) + \sum_{j} \frac{Q_j}{kT_j} = \Delta \Sigma \tag{S10}$$

**Definitions**:
*   **$k$**: Boltzmann constant.
*   **$S(\rho) = -tr[\rho \ln(\rho)]$**: von Neumann entropy.
*   **$\Delta \Sigma = \Sigma(t+\tau) - \Sigma(t)$**: Where $\Sigma = I(S,R) + C(S) + C(R) + \sum_j D(\rho_j || \rho^{th}_j)$.
*   **$I(S,R)$**: Mutual information between the system and all reservoirs.
*   **$C(S), C(R)$**: Total correlations between all subsystems and all reservoirs, respectively.
*   **$D(\rho_j || \rho^{th}_j)$**: Relative entropy between the state of $R_j$ and the thermal state $\rho^{th}_j = \exp(-H_j / kT_j) / Z_j$.

Using the identity:
$$\sum_i \Delta S(\rho_i) + \sum_j \Delta S(\rho_j) = \Delta I(S,R) + \Delta C(S) + \Delta C(R) \tag{S11}$$

The relative entropy is written explicitly as:
$$D(\rho_j(t) || \rho^{th}_j) = -S(\rho_j) + \frac{1}{kT_j} tr[\rho_j(t) H_j] + \ln(Z_j) \tag{S12}$$

Which implies:
$$\Delta S(\rho_j) = \frac{Q_j}{kT_j} - \Delta D(\rho_j || \rho^{th}_j) \tag{S13}$$

Replacing Eq. (S13) into Eq. (S11) results in the **generalized second law**:
$$\sum_i \Delta S(\rho_i) + \sum_j \frac{Q_j}{kT_j} = \underbrace{\Delta I(S,R) + \Delta C(S) + \Delta C(R) + \sum_j \Delta D(\rho_j || \rho^{th}_j)}_{\Delta \Sigma} \tag{S14}$$

---

## Derivation of the Generalized Efficiency

The efficiency $\eta$ for a quantum engine is derived as:
$$\eta = \gamma \left( \eta_{th} - \frac{kT_{min} \Delta \sigma}{Q_{in}} \right) \tag{S15}$$

**Parameters**:
*   **$\gamma = [1 + \Delta U_{in} / Q_{in}]^{-1}$**.
*   **$\eta_{th} = -\sum_j \eta_j Q_j / Q_{in}$**.
*   **$F_i = U_i - kT_{min} S(\rho_i)$**: Generalized free energy.
*   **$kT_{min} \Delta \sigma = kT_{min} \Delta \Sigma + \sum_i \Delta F_i$**.

Combining the first and second laws, work can be expressed as:
$$W = \sum_j \underbrace{\left( 1 - \frac{T_{min}}{T_j} \right)}_{\eta_j} Q_j + \underbrace{kT_{min} \Delta \Sigma + \sum_i \Delta F_i}_{kT_{min} \Delta \sigma} \tag{S18}$$

The final efficiency formula is determined by:
$$\eta = \gamma \left[ \underbrace{-\sum_j \frac{\eta_j Q_j}{Q_{in}}}_{\eta_{th}} - \frac{kT_{min} \Delta \sigma}{Q_{in}} \right] \tag{S19}$$

---

## Comparison with Previous Approaches

The sources compare these findings to the generalized **Clausius inequality** for correlated systems (Ref 13):
$$-\tilde{Q}_A (T_B - T_A) \ge kT_A T_B \Delta I(R_A, R_B) \tag{S20}$$

The sources extend this to an exact equality:
$$-\tilde{Q}_A (T_B - T_A) = kT_A T_B \Delta I(R_A, R_B) + kT_A [ T_A \Delta D(\rho_A || \rho^{th}_A) + T_B \Delta D(\rho_B || \rho^{th}_B) ] \tag{S22}$$
