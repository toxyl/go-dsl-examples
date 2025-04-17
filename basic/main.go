package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"strings"
	"time"

	"github.com/toxyl/math"

	"github.com/toxyl/flo"
)

func main() {
	r, err := dsl.run(strings.Join(os.Args[1:], " "), true)
	if err != nil {
		fmt.Println("\x1b[31mError:\x1b[0m", err)
		return
	}
	var res any
	if r.err != nil {
		fmt.Println("\x1b[31mError:\x1b[0m", r.err)
		return
	}
	res = r.value
	fmt.Printf("\x1b[32mResult:\x1b[0m %v\n\n", res)

	v := dsl.vars.get("gx")
	gx := v.get().(float64)
	fmt.Println("\x1b[32mgx\x1b[0m =", "gx", "*", 10)
	dsl.vars.set("gx", gx*10)
	fmt.Println("\x1b[32mgx\x1b[0m =", v.get())

	flo.File("doc.md").StoreString(dsl.docMarkdown())
	flo.File("doc.html").StoreString(dsl.docHTML())
}

///////////////////////////////////////////////////
// Below are all defined variables and functions //
///////////////////////////////////////////////////

var (
	// @Name:   gx
	// @Desc:   Global variable
	// @Range: 	0..10
	// @Unit:   px
	globalX = 20.0

	// @Name:   gy
	// @Desc:   Global variable
	// @Range: 	0..+Inf
	// @Unit:   px
	globalY = true
)

// @Name:    test-function-1
// @Desc:    This is a test function
// @Param:   x 		px 0..10 	0 		Position on the x axis
// @Param:   y 		px 0..10 	0 		Position on the y axis
// @Param:   str 				"hi" 	String to print
// @Returns: z 		px 0..20 	0 		Position on the z axis
func funcA(x, y int, str string) (z int, err error) {
	z = x + y
	fmt.Println(str, x, y, z)
	return
}

// @Name:    test-function-2
// @Desc:    This is a test function
// @Param:	 lat ° 		-90..90    	0 		Latitude
// @Param:   lon ° 		-180..180  	0 		Longitude
// @Returns: z   					false  	Is the point in the ocean?
// @Returns: err 							Something bad happened
func funcB(lat, lon float64) (z bool, err error) {
	z = lat+lon > 0
	return
}

// @Name:    flip
// @Desc:    Flips a boolean
// @Param:   x	false 	Boolean to flip
// @Returns: y	true 	Flipped boolean
func flip(x bool) (y bool, err error) {
	y = !x
	return
}

// @Name:    sign
// @Desc:    Returns the sign of a float
// @Param:   x 	-	- 	0.0 	Float to get the sign of
// @Returns: y 	-   -	0 		Sign of the float
func sign(x float64) (y int, err error) {
	if x > 0 {
		y = 1
	} else if x < 0 {
		y = -1
	} else {
		y = 0
	}
	return
}

// @Name:    +
// @Desc:    Adds two integers
// @Param:   x	-	-	0.0	First float
// @Param:   y 	-	-	0.0	Second float
// @Returns: z 	-	-	0.0	Sum of the two floats
func add(x, y float64) (z float64, err error) {
	z = x + y
	return
}

// @Name:    -
// @Desc:    Subtracts two floats
// @Param:   x 	-	-	0.0	First float
// @Param:   y 	-	-	0.0	Second float
// @Returns: z 	-	-	0.0	Difference of the two floats
func sub(x, y float64) (z float64, err error) {
	z = x - y
	return
}

// @Name:    *
// @Desc:    Multiplies two floats
// @Param:   x 	-	-	0.0	First float
// @Param:   y 	-	-	0.0	Second float
// @Returns: z 	-	-	0.0	Product of the two floats
func mul(x, y float64) (z float64, err error) {
	z = x * y
	return
}

// @Name:    /
// @Desc:    Divides two floats
// @Param:   x 	-	-	0.0	First float
// @Param:   y 	-	-	0.0	Second float
// @Returns: z 	-	-	0.0	Quotient of the two floats
func div(x, y float64) (z float64, err error) {
	z = x / y
	return
}

// @Name:    mod
// @Desc:    Modulus of two floats
// @Param:   x 	-	-	0.0	First float
// @Param:   y 	-	-	0.0	Second float
// @Returns: z 	-	-	0.0	Modulus of the two floats
func mod(x, y float64) (z float64, err error) {
	z = math.Mod(x, y)
	return
}

// @Name:    pow
// @Desc:    Raises a float to the power of another float
// @Param:   x 	-	-	0.0	Base
// @Param:   y 	-	-	0.0	Exponent
// @Returns: z 	-	-	0.0	Result of the power operation
func pow(x, y float64) (z float64, err error) {
	z = math.Pow(x, y)
	return
}

// @Name:    sqrt
// @Desc:    Square root of a float
// @Param:   x 	-	-	0.0	Float to get the square root of
// @Returns: z 	-	-	0.0	Result of the square root operation
func sqrt(x float64) (z float64, err error) {
	z = math.Sqrt(x)
	return
}

// @Name:    sin
// @Desc:    Sine of a float
// @Param:   x 	-	-	0.0	Float to get the sine of
// @Returns: z 	-	-	0.0	Result of the sine operation
func sin(x float64) (z float64, err error) {
	z = math.Sin(x)
	return
}

// @Name:    cos
// @Desc:    Cosine of a float
// @Param:   x 	-	-	0.0	Float to get the cosine of
// @Returns: z 	-	-	0.0	Result of the cosine operation
func cos(x float64) (z float64, err error) {
	z = math.Cos(x)
	return
}

// @Name:    tan
// @Desc:    Tangent of a float
// @Param:   x 	-	-	0.0	Float to get the tangent of
// @Returns: z 	-	-	0.0	Result of the tangent operation
func tan(x float64) (z float64, err error) {
	z = math.Tan(x)
	return
}

// @Name:    atan
// @Desc:    Arctangent of a float
// @Param:   x 	-	-	0.0	Float to get the arctangent of
// @Returns: z 	-	-	0.0	Result of the arctangent operation
func atan(x float64) (z float64, err error) {
	z = math.Atan(x)
	return
}

// @Name:    atan2
// @Desc:    Arctangent of two floats
// @Param:   x 	-	-	0.0	Float to get the arctangent of
// @Param:   y 	-	-	0.0	Float to get the arctangent of
// @Returns: z 	-	-	0.0	Result of the arctangent operation
func atan2(x, y float64) (z float64, err error) {
	z = math.Atan2(x, y)
	return
}

// @Name:    abs
// @Desc:    Absolute value of a float
// @Param:   x 	-	-	0.0	Float to get the absolute value of
// @Returns: z 	-	-	0.0	Result of the absolute value operation
func abs(x float64) (z float64, err error) {
	z = math.Abs(x)
	return
}

// @Name:    floor
// @Desc:    Floor of a float
// @Param:   x 	-	-	0.0	Float to get the floor of
// @Returns: z 	-	-	0.0	Result of the floor operation
func floor(x float64) (z float64, err error) {
	z = math.Floor(x)
	return
}

// @Name:    ceil
// @Desc:    Ceil of a float
// @Param:   x 	-	-	0.0	Float to get the ceil of
// @Returns: z 	-	-	0.0	Result of the ceil operation
func ceil(x float64) (z float64, err error) {
	z = math.Ceil(x)
	return
}

// @Name:    round
// @Desc:    Round a float
// @Param:   x 	-	-	0.0	Float to get the round of
// @Returns: z 	-	-	0.0	Result of the round operation
func round(x float64) (z float64, err error) {
	z = math.Round(x)
	return
}

// @Name:    exp
// @Desc:    Exponential of a float
// @Param:   x 	-	-	0.0	Float to get the exponential of
// @Returns: z 	-	-	0.0	Result of the exponential operation
func exp(x float64) (z float64, err error) {
	z = math.Exp(x)
	return
}

// @Name:    log
// @Desc:    Logarithm of a float
// @Param:   x 	-	-	0.0	Float to get the logarithm of
// @Returns: z 	-	-	0.0	Result of the logarithm operation
func log(x float64) (z float64, err error) {
	z = math.Log(x)
	return
}

// @Name:    log10
// @Desc:    Logarithm base 10 of a float
// @Param:   x 	-	-	0.0	Float to get the logarithm base 10 of
// @Returns: z 	-	-	0.0	Result of the logarithm base 10 operation
func log10(x float64) (z float64, err error) {
	z = math.Log10(x)
	return
}

// @Name:    log2
// @Desc:    Logarithm base 2 of a float
// @Param:   x 	-	-	0.0	Float to get the logarithm base 2 of
// @Returns: z 	-	-	0.0	Result of the logarithm base 2 operation
func log2(x float64) (z float64, err error) {
	z = math.Log2(x)
	return
}

// @Name:    logb
// @Desc:    Logarithm base b of a float
// @Param:   x 	-	-	0.0	Float to get the logarithm base b of
// @Returns: z 	-	-	0.0	Result of the logarithm base b operation
func logb(x float64) (z float64, err error) {
	z = math.Logb(x)
	return
}

// @Name:    log1p
// @Desc:    Logarithm of 1+x
// @Param:   x 	-	-	0.0	Float to get the logarithm of 1+x
// @Returns: z 	-	-	0.0	Result of the logarithm of 1+x operation
func log1p(x float64) (z float64, err error) {
	z = math.Log1p(x)
	return
}

// @Name:    millis
// @Desc:    Current time in milliseconds
// @Returns: z 	-	0..+Inf	0.0	Current time in milliseconds
func millis() (z int64, err error) {
	z = time.Now().UnixMilli()
	return
}

// @Name:    seconds
// @Desc:    Current time in seconds
// @Returns: z 	-	0..+Inf	0.0	Current time in seconds
func seconds() (z int64, err error) {
	z = time.Now().Unix()
	return
}

// @Name:    nanos
// @Desc:    Current time in nanoseconds
// @Returns: z 	-	0..+Inf	0.0	Current time in nanoseconds
func nanos() (z int64, err error) {
	z = time.Now().UnixNano()
	return
}

var t = time.Now()

// @Name:    runtime
// @Desc:    Runtime of the program
// @Returns: z 	-	0..+Inf	0.0	Runtime of the program
func runtime() (z int64, err error) {
	z = time.Since(t).Milliseconds()
	return
}

// @Name:    rand
// @Desc:    Random float between 0 and 1
// @Returns: z 	-	0..1	0.0	Random float between 0 and 1
func random() (z float64, err error) {
	z = rand.Float64()
	return
}

// @Name:    randInt
// @Desc:    Random integer between 0 and Inf
// @Returns: z 	-	0..Inf	0	Random integer between 0 and, not including, +Inf
func randInt() (z int, err error) {
	z = rand.Int()
	return
}
