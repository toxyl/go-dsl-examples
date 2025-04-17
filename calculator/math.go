package main

import (
	"fmt"

	"github.com/toxyl/math"
)

var (
	// @Name:  last
	// @Desc:  last result
	// @Range: -
	// @Unit:  -
	res = 0.0
)

// @Name: add
// @Desc: Adds two numbers
// @Param:      x       - -   0   First number
// @Param:      y       - -   0   Second number
// @Returns:    result  - -   0   Sum of the numbers
func add(x, y float64) (result float64, err error) {
	res = x + y
	return res, nil
}

// @Name: sub
// @Desc: Subtracts two numbers
// @Param:      x       - -   0   First number
// @Param:      y       - -   0   Second number
// @Returns:    result  - -   0   Difference of the numbers
func sub(x, y float64) (result float64, err error) {
	res = x - y
	return res, nil
}

// @Name: mul
// @Desc: Multiplies two numbers
// @Param:      x       - -   0   First number
// @Param:      y       - -   0   Second number
// @Returns:    result  - -   0   Product of the numbers
func mul(x, y float64) (result float64, err error) {
	res = x * y
	return res, nil
}

// @Name: div
// @Desc: Divides two numbers
// @Param:      x       - -   0   First number
// @Param:      y       - -   0   Second number
// @Returns:    result  - -   0   Quotient of the numbers
func div(x, y float64) (result float64, err error) {
	if y == 0 {
		res = 0
		return res, fmt.Errorf("division by zero")
	}
	res = x / y
	return res, nil
}

// @Name: pow
// @Desc: Raises a number to a pow
// @Param:      x       - -   1   Base number
// @Param:      y       - -   1   Exponent
// @Returns:    result  - -   1   Base raised to the pow of exp
func pow(x, y float64) (result float64, err error) {
	res = math.Pow(x, y)
	return res, nil
}

// @Name: sqrt
// @Desc: Calculates square root of a number
// @Param:      x       - -   0   Number to calculate square root of
// @Returns:    result  - -   0   Square root of the number
func sqrt(x float64) (result float64, err error) {
	if x < 0 {
		res = 0
		return res, fmt.Errorf("cannot calculate square root of negative number")
	}
	res = math.Sqrt(x)
	return res, nil
}

// @Name: abs
// @Desc: Returns the absolute value of a number
// @Param:      x       -	-   0   Number to get absolute value of
// @Returns:    result  -	-   0   Absolute value of the number
func abs(x float64) (result float64, err error) {
	res = math.Abs(x)
	return res, nil
}

// @Name: sin
// @Desc: Calculates the sine of a number (in radians)
// @Param:      x       rad	-   0   Angle in radians
// @Returns:    result  -	-   0   Sine of the angle
func sin(x float64) (result float64, err error) {
	res = math.Sin(x)
	return res, nil
}

// @Name: cos
// @Desc: Calculates the cosine of a number (in radians)
// @Param:      x       -	-   0   Angle in radians
// @Returns:    result  -	-   0   Cosine of the angle
func cos(x float64) (result float64, err error) {
	res = math.Cos(x)
	return res, nil
}

// @Name: tan
// @Desc: Calculates the tangent of a number (in radians)
// @Param:      x       -	-   0   Angle in radians
// @Returns:    result  -	-   0   Tangent of the angle
func tan(x float64) (result float64, err error) {
	res = math.Tan(x)
	return res, nil
}

// @Name: log
// @Desc: Calculates the natural logarithm of a number
// @Param:      x       -	-   0   Number to calculate logarithm of
// @Returns:    result  -	-   0   Natural logarithm of the number
func log(x float64) (result float64, err error) {
	if x <= 0 {
		res = 0
		return res, fmt.Errorf("cannot calculate logarithm of non-positive number")
	}
	res = math.Log(x)
	return res, nil
}

// @Name: exp
// @Desc: Calculates e raised to the power of x
// @Param:      x       -	-   0   Exponent
// @Returns:    result  -	-   0   e raised to the power of x
func exp(x float64) (result float64, err error) {
	res = math.Exp(x)
	return res, nil
}

// @Name: floor
// @Desc: Returns the greatest integer value less than or equal to x
// @Param:      x       -	-   0   Number to round down
// @Returns:    result  -	-   0   Greatest integer less than or equal to x
func floor(x float64) (result float64, err error) {
	res = math.Floor(x)
	return res, nil
}

// @Name: ceil
// @Desc: Returns the least integer value greater than or equal to x
// @Param:      x       -	-   0   Number to round up
// @Returns:    result  -	-   0   Least integer greater than or equal to x
func ceil(x float64) (result float64, err error) {
	res = math.Ceil(x)
	return res, nil
}
