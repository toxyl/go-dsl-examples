package main

var (
	// @Name:  v
	// @Desc:  Voltage
	// @Range: -
	// @Unit:  V
	voltage = 0.0

	// @Name:  r
	// @Desc:  Resistance
	// @Range: -
	// @Unit:  Ω
	resistance = 0.0

	// @Name:  i
	// @Desc:  Current
	// @Range: -
	// @Unit:  A
	current = 0.0

	// @Name:  p
	// @Desc:  Power
	// @Range: -
	// @Unit:  W
	power = 0.0
)

// @Name: V
// @Desc: Calculates the voltage
// @Param:      r       - -   0   Resistance
// @Param:      i       - -   0   Current
// @Returns:    result  - -   0   Voltage
func Voltage(r, i float64) (float64, error) {
	resistance = r
	current = i
	voltage = resistance * current
	power = voltage * current
	return voltage, nil
}

// @Name: I
// @Desc: Calculates the current
// @Param:      v       - -   0   Voltage
// @Param:      r       - -   0   Resistance
// @Returns:    result  - -   0   Current
func Current(v, r float64) (float64, error) {
	voltage = v
	resistance = r
	current = voltage / resistance
	power = voltage * current
	return current, nil
}

// @Name: R
// @Desc: Calculates the resistance
// @Param:      v       - -   0   Voltage
// @Param:      i       - -   0   Current
// @Returns:    result  - -   0   Resistance
func Resistance(v, i float64) (float64, error) {
	voltage = v
	current = i
	resistance = voltage / current
	power = voltage * current
	return resistance, nil
}

// @Name: P
// @Desc: Calculates the power
// @Param:      v       - -   0   Voltage
// @Param:      i       - -   0   Current
// @Returns:    result  - -   0   Power
func Power(v, i float64) (float64, error) {
	voltage = v
	current = i
	resistance = voltage / current
	power = voltage * current
	return power, nil
}
